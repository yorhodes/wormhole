package reactor

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/certusone/wormhole/node/pkg/supervisor"

	node_common "github.com/certusone/wormhole/node/pkg/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap"

	gossipv1 "github.com/certusone/wormhole/node/pkg/proto/gossip/v1"
	"github.com/wormhole-foundation/wormhole/sdk/vaa"
)

// ConsensusReactor implements the full consensus processor for a single Observation. It cannot be reused after being
// finalized.
type ConsensusReactor[K Observation] struct {
	// Current state of the reactor
	state      State
	stateMutex sync.Mutex

	// First time this digest was seen (possibly even before we observed it ourselves).
	firstSeen time.Time
	// Time of the last new observation received
	lastObservation time.Time
	// The most recent time that the signature / observation has been transmitted
	lastTransmission time.Time
	// Time quorum was reached
	timeQuorum time.Time

	// Copy of our observation.
	observation K
	// Map of signatures seen by guardian. During guardian set updates, this may contain signatures belonging
	// to either the old or new guardian set.
	signatures map[ethcommon.Address][]byte
	// Copy of the bytes we submitted (ourObservation, but signed and serialized). Used for retransmissions.
	localSignature []byte
	// Copy of the guardian set valid at observation/injection time.
	gs *node_common.GuardianSet

	// Channel for receiving local observations
	observationChannel chan K
	// Channel for receiving foreign observations
	foreignObservationChannel chan *gossipv1.SignedObservation

	// Hook to be called on a state transition
	stateTransitionHook StateTransitionHook[K]
	// Configuration of the reactor
	config Config

	logger *zap.Logger
}

// StateTransitionHook is a hook for handling state transitions. It is called as a goroutine after the state transition
// is completed.
type StateTransitionHook[K Observation] func(reactor *ConsensusReactor[K], oldState, newState State)

// Config allows to parametrize the consensus reactor
type Config struct {
	// RetransmitFrequency is the frequency of observation rebroadcasts
	RetransmitFrequency time.Duration
	// QuorumGracePeriod is the time to wait for more signatures after quorum before finalizing the reactor
	QuorumGracePeriod time.Duration
	// QuorumTimeout is the time to wait for quorum before finalizing the reactor
	QuorumTimeout time.Duration
	// UnobservedTimeout is the time to wait for either a completed VAA or a local observation before finalizing the
	// reactor after only having received remote observations.
	UnobservedTimeout time.Duration
	// Signer to use for local observations. If Signer is nil, the reactor will not participate in consensus.
	Signer Signer
	// NetworkAdapter used for broadcasting signatures. If NetworkAdapter is nil, the reactor will not broadcast local
	// signatures.
	NetworkAdapter NetworkAdapter
}

// Observation defines the interface for any events observed by the guardian.
type Observation interface {
	// GetEmitterChain returns the id of the chain where this event was observed.
	GetEmitterChain() vaa.ChainID
	// MessageID returns a human-readable emitter_chain/emitter_address/sequence tuple.
	MessageID() string
	// SigningMsg returns the hash of the signing body of the observation. This is used
	// for signature generation and verification.
	SigningMsg() ethcommon.Hash
}

// State of the reactor
type State string

const (
	// StateInitialized is used for a freshly created reactor. A reactor in the StateInitialized will wait for either
	// a local or foreign observation.
	StateInitialized State = "initialized"
	// StateObserved indicates that the reactor has seen =1 local observation and >= 0 foreign observations. It is able
	// to produce a full VAA upon reaching quorum.
	StateObserved State = "observed"
	// StateUnobserved indicates that the reactor has seen >= 1 foreign observations but no local observation. It is not
	// able to produce a full VAA without a local observation.
	StateUnobserved State = "unobserved"
	// StateQuorum indicates that the reactor has seen a local observation and a quorum of foreign observations. It has
	// all data to produce a full VAA.
	StateQuorum State = "quorum"
	// StateQuorumUnobserved indicates that the reactor has seen a quorum of foreign observations but no local observation.
	// It can only produce a full VAA after receiving a local observation.
	StateQuorumUnobserved State = "quorum_unobserved"
	// StateFinalized is a reactor that has gone through a full lifecycle. It holds all information required to produce
	// a full VAA.
	StateFinalized State = "finalized"
	// StateTimedOut is a reactor that has gone through a full lifecycle. It did not manage to achieve locally confirmed
	// locally confirmed quorum (i.e. both a local observation and quorum) within the configured timeouts.
	StateTimedOut State = "timed_out"
)

func NewReactor[K Observation](config Config, gs *node_common.GuardianSet, s StateTransitionHook[K]) *ConsensusReactor[K] {
	c := &ConsensusReactor[K]{
		state:                     StateInitialized,
		signatures:                map[ethcommon.Address][]byte{},
		gs:                        gs,
		stateTransitionHook:       s,
		config:                    config,
		foreignObservationChannel: make(chan *gossipv1.SignedObservation, 10),
		observationChannel:        make(chan K, 10),
	}

	return c
}

func (c *ConsensusReactor[K]) Run(ctx context.Context) error {
	c.logger = supervisor.Logger(ctx).With(zap.String("state", string(c.state)))

	c.stateMutex.Lock()
	if c.state == StateFinalized {
		c.stateMutex.Unlock()
		return nil
	}
	c.stateMutex.Unlock()

	tickFrequency := time.Second
	ticker := time.NewTicker(tickFrequency)
	defer ticker.Stop()

	supervisor.Signal(ctx, supervisor.SignalHealthy)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			terminate, err := c.housekeeping(ctx)
			if err != nil {
				c.logger.Warn("housekeeping failed", zap.Error(err))
			}
			if terminate {
				c.logger.Debug("reactor concluded; shutting down processing loop")
				// Signal done such that the routine does not restart
				// TODO this might leave garbage in the supervisor tree
				supervisor.Signal(ctx, supervisor.SignalDone)
				return nil
			}
		case o := <-c.observationChannel:
			c.observed(ctx, o)
		case o := <-c.foreignObservationChannel:
			c.observationReceived(ctx, o)
		}
	}
}

func (c *ConsensusReactor[K]) ObservationChannel() chan<- K {
	return c.observationChannel
}

func (c *ConsensusReactor[K]) ForeignObservationChannel() chan<- *gossipv1.SignedObservation {
	return c.foreignObservationChannel
}

// State returns the current state of the reactor
func (c *ConsensusReactor[K]) State() State {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()
	return c.state
}

// VAASignatures returns the stored signatures in the order required by a VAA.
func (c *ConsensusReactor[K]) VAASignatures() []*vaa.Signature {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()
	agg := make([]bool, len(c.gs.Keys))
	var sigs []*vaa.Signature
	for i, a := range c.gs.Keys {
		s, ok := c.signatures[a]

		if ok {
			var bs [65]byte
			if n := copy(bs[:], s); n != 65 {
				panic(fmt.Sprintf("invalid sig len: %d", n))
			}

			sigs = append(sigs, &vaa.Signature{
				Index:     uint8(i),
				Signature: bs,
			})
		}

		agg[i] = ok
	}
	return sigs
}

func (c *ConsensusReactor[K]) housekeeping(ctx context.Context) (bool, error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()
	defer c.logger.Debug("housekeeping completed")

	c.logger.Debug("starting housekeeping")

	switch c.state {
	case StateUnobserved:
		if time.Since(c.firstSeen) > c.config.UnobservedTimeout {
			c.logger.Debug("timing out", zap.String("reason", "unobserved_timeout"))
			// Time out
			c.timeOut()
		}
	case StateObserved:
		if time.Since(c.lastObservation) > c.config.QuorumTimeout {
			c.logger.Debug("timing out", zap.String("reason", "quorum_timeout"))
			// Time out
			c.timeOut()
		}

		if time.Since(c.lastTransmission) > c.config.RetransmitFrequency {
			c.logger.Debug("retransmitting", zap.String("reason", "unobserved_timeout"))
			// Retransmit signature
			err := c.transmitSignature(ctx)
			if err != nil {
				c.logger.Error("failed to retransmit signature", zap.Error(err))
			}
		}
	case StateQuorum:
		if time.Since(c.timeQuorum) > c.config.QuorumGracePeriod || len(c.gs.Keys) == len(c.signatures) {
			c.logger.Debug("timing out", zap.String("reason", "quorum_grace"))
			// Time out
			c.timeOut()
		}
	case StateQuorumUnobserved:
		if time.Since(c.firstSeen) > c.config.UnobservedTimeout {
			c.logger.Debug("timing out", zap.String("reason", "quorum_unobserved_timeout"))
			// Time out
			c.timeOut()
		}
	case StateFinalized, StateTimedOut:
		// This is the final iteration. Do cleanup
		return true, nil
	}

	return false, nil
}

func (c *ConsensusReactor[K]) timeOut() {
	if c.state == StateQuorum {
		c.stateTransition(StateFinalized)
	} else {
		c.stateTransition(StateTimedOut)
	}
}

func (c *ConsensusReactor[K]) stateTransition(to State) {
	c.logger.Debug("state transition", zap.String("from", string(c.state)), zap.String("to", string(to)))
	previousState := c.state
	c.state = to

	if c.stateTransitionHook != nil {
		go c.stateTransitionHook(c, previousState, to)
	}
}

func (c *ConsensusReactor[K]) observationReceived(ctx context.Context, m *gossipv1.SignedObservation) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	c.logger.Debug("received foreign observation", zap.Any("observation", m))

	if c.state != StateObserved && c.state != StateUnobserved && c.state != StateQuorum && c.state != StateQuorumUnobserved && c.state != StateInitialized {
		return
	}

	hash := hex.EncodeToString(m.Hash)

	// Hooray! Now, we have verified all fields on SignedObservation and know that it includes
	// a valid signature by an active guardian. We still don't fully trust them, as they may be
	// byzantine, but now we know who we're dealing with.
	err := verifySignedObservation(m, c.gs)
	if err != nil {
		c.logger.Info("failed to verify signed observation",
			zap.Error(err),
			zap.String("digest", hash),
			zap.String("signature", hex.EncodeToString(m.Signature)),
			zap.String("addr", hex.EncodeToString(m.Addr)),
		)
		return
	}
	theirAddr := ethcommon.BytesToAddress(m.Addr)

	c.logger.Debug("accepted foreign observation", zap.Any("observation", m), zap.Stringer("address", theirAddr))

	// Have we already received this observation
	if _, has := c.signatures[theirAddr]; has {
		// TODO log duplicate
		return
	}

	// Store their signature
	c.signatures[theirAddr] = m.Signature
	c.lastObservation = time.Now()

	if c.state == StateInitialized {
		c.logger.Debug("received observation before own observation", zap.Any("observation", m))
		c.firstSeen = time.Now()
		c.stateTransition(StateUnobserved)
	}

	// If we haven't reached quorum yet, there is nothing more to do
	if !c.hasQuorum() {
		return
	}

	// Transition to quorum states
	switch c.state {
	case StateObserved:
		c.stateTransition(StateQuorum)
	case StateUnobserved:
		c.stateTransition(StateQuorumUnobserved)
	}
}

func (c *ConsensusReactor[K]) observed(ctx context.Context, o K) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	c.logger.Debug("observed", zap.Any("observation", o))

	if c.state != StateInitialized && c.state != StateUnobserved && c.state != StateQuorumUnobserved {
		return
	}

	c.observation = o

	if c.config.Signer != nil {
		// Generate digest of the unsigned VAA.
		digest := o.SigningMsg()

		// Sign the digest using our node's guardian key.
		timeout, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		localAddr, err := c.config.Signer.Address(timeout)
		if err != nil {
			panic(err)
		}
		s, err := c.config.Signer.Sign(timeout, digest.Bytes())
		if err != nil {
			panic(err)
		}
		c.localSignature = s

		// Store in signatures array
		c.signatures[localAddr] = s

		err = c.transmitSignature(ctx)
		if err != nil {
			c.logger.Error("failed to transmit signature on observation", zap.Error(err))
		}

	}

	// Transition to quorum states
	switch c.state {
	case StateInitialized:
		c.firstSeen = time.Now()
		c.lastObservation = time.Now()
		c.stateTransition(StateObserved)
	case StateUnobserved:
		c.stateTransition(StateObserved)
	case StateQuorumUnobserved:
		c.stateTransition(StateQuorum)
	}

	// If we haven't reached quorum yet, there is nothing more to do
	if !c.hasQuorum() {
		return
	}

	// We immediately reached quorum
	switch c.state {
	case StateObserved:
		c.stateTransition(StateQuorum)
	}
}

func (c *ConsensusReactor[K]) hasQuorum() bool {
	return len(c.signatures) >= vaa.CalculateQuorum(len(c.gs.Keys))
}

// State returns the current state of the reactor
func (c *ConsensusReactor[K]) transmitSignature(ctx context.Context) error {
	if c.config.Signer == nil {
		return fmt.Errorf("can't broadcast signature without signer")
	}
	if c.config.NetworkAdapter == nil {
		return fmt.Errorf("can't broadcast signature without network adapter")
	}

	timeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	addr, err := c.config.Signer.Address(timeout)
	if err != nil {
		return fmt.Errorf("failed to get signer address for signature broadcast: %w", err)
	}
	signedO := &gossipv1.SignedObservation{
		Addr:      addr.Bytes(),
		Hash:      c.observation.SigningMsg().Bytes(),
		Signature: c.localSignature,
		TxHash:    []byte{},
		MessageId: c.observation.MessageID(),
	}
	timeout, cancel = context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err = c.config.NetworkAdapter.BroadcastObservation(timeout, signedO)
	if err != nil {
		return fmt.Errorf("failed to broadcast observation: %w", err)
	}

	return nil
}

func verifySignedObservation(m *gossipv1.SignedObservation, gs *node_common.GuardianSet) error {
	// Verify the Guardian's signature. This verifies that m.Signature matches m.Hash and recovers
	// the public key that was used to sign the payload.
	pk, err := crypto.Ecrecover(m.Hash, m.Signature)
	if err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	// Verify that m.Addr matches the public key that signed m.Hash.
	theirAddr := ethcommon.BytesToAddress(m.Addr)
	sigPub := ethcommon.BytesToAddress(crypto.Keccak256(pk[1:])[12:])

	if theirAddr != sigPub {
		return fmt.Errorf("address does not match pubkey: %s != %s", theirAddr.Hex(), sigPub.Hex())
	}

	// Verify that m.Addr is included in the guardian set. If it's not, drop the message. In case it's us
	// who have the outdated guardian set, we'll just wait for the message to be retransmitted eventually.
	_, ok := gs.KeyIndex(theirAddr)
	if !ok {
		return fmt.Errorf("unknown guardian: %s", theirAddr.Hex())
	}

	return nil
}
