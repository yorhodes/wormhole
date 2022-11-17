package reactor

import (
	"context"
	"fmt"
	"sync"

	"github.com/certusone/wormhole/node/pkg/common"
	gossipv1 "github.com/certusone/wormhole/node/pkg/proto/gossip/v1"
	"github.com/certusone/wormhole/node/pkg/supervisor"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/wormhole-foundation/wormhole/sdk/vaa"
	"go.uber.org/zap"
)

// Manager handles the creation and maintenance of reactors for incoming local and foreign observations
type Manager[K Observation] struct {
	// observationC is a channel of observed emitted messages
	observationC <-chan K

	// confirmationC is a channel of inbound decoded observations from p2p
	confirmationC <-chan *gossipv1.SignedObservation

	// setC is a channel of guardian set updates
	setC <-chan *common.GuardianSet

	// reactors are the current live reactors. This list may contain prepared reactors in StateInitialized
	reactors     map[ethcommon.Hash]*ConsensusReactor[K]
	reactorsLock sync.Mutex

	// gs is the currently valid guardian set
	gs *common.GuardianSet
	// gst is managed by the processor and allows concurrent access to the
	// guardian set by other components.
	gst *common.GuardianSetState

	// config template for the reactors
	config Config
	// handler for manager events
	handler ManagerEventHandler[K]

	logger *zap.Logger
}

// ManagerEventHandler handles significant consensus event from reactors
type ManagerEventHandler[K Observation] interface {
	HandleQuorum(observation K, signatures []*vaa.Signature)
	HandleFinalization(observation K, signatures []*vaa.Signature)
	HandleTimeout(previousState State, observation K, signatures []*vaa.Signature)
}

// NewManager creates a new reactor manager
func NewManager[K Observation](observationC <-chan K, confirmationC <-chan *gossipv1.SignedObservation, setC <-chan *common.GuardianSet, gst *common.GuardianSetState, config Config, handler ManagerEventHandler[K]) *Manager[K] {

	return &Manager[K]{
		observationC:  observationC,
		confirmationC: confirmationC,
		setC:          setC,
		reactors:      map[ethcommon.Hash]*ConsensusReactor[K]{},
		gs:            gst.Get(),
		gst:           gst,
		config:        config,
		handler:       handler,
	}
}

func (p *Manager[K]) Run(ctx context.Context) error {
	p.logger = supervisor.Logger(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case p.gs = <-p.setC:
			p.logger.Info("guardian set updated",
				zap.Strings("set", p.gs.KeysAsHexStrings()),
				zap.Uint32("index", p.gs.Index))
			p.gst.Set(p.gs)
		case k := <-p.observationC:
			digest := k.SigningMsg()
			p.logger.Debug("received observation", zap.Stringer("digest", digest), zap.Any("observation", k))

			var r *ConsensusReactor[K]
			p.reactorsLock.Lock()
			if reactor, exists := p.reactors[digest]; exists {
				p.logger.Debug("sending observation to existing reactor", zap.Stringer("digest", digest))
				r = reactor
			} else {
				p.logger.Debug("creating new reactor for observation", zap.Stringer("digest", digest))
				r = NewReactor[K](p.config, p.gs, func(reactor *ConsensusReactor[K], oldState, newState State) {
					p.transitionHook(digest, reactor, oldState, newState)
				})
				err := supervisor.Run(ctx, fmt.Sprintf("reactor-%s", digest.String()), r.Run)
				if err != nil {
					p.logger.Error("failed to spawn reactor routine", zap.Error(err), zap.Stringer("digest", digest))
				}
				p.reactors[digest] = r
			}
			p.reactorsLock.Unlock()
			r.ObservationChannel() <- k
		case m := <-p.confirmationC:
			digest := ethcommon.BytesToHash(m.TxHash)
			var r *ConsensusReactor[K]
			p.reactorsLock.Lock()
			if reactor, exists := p.reactors[digest]; exists {
				r = reactor
			} else {
				r = NewReactor[K](p.config, p.gs, func(reactor *ConsensusReactor[K], oldState, newState State) {
					p.transitionHook(digest, reactor, oldState, newState)
				})
				err := supervisor.Run(ctx, fmt.Sprintf("reactor-%s", digest.String()), r.Run)
				if err != nil {
					p.logger.Error("failed to spawn reactor routine", zap.Error(err), zap.Stringer("digest", digest))
				}
				p.reactors[digest] = r
			}
			p.reactorsLock.Unlock()
			r.ForeignObservationChannel() <- m
		}
	}
}

func (p *Manager[K]) transitionHook(digest ethcommon.Hash, reactor *ConsensusReactor[K], oldState, newState State) {
	switch newState {
	case StateQuorum:
		p.logger.Debug("reactor reached quorum", zap.Stringer("digest", digest))
		// Handle consensus
		go p.handler.HandleQuorum(reactor.observation, reactor.VAASignatures())
	case StateFinalized:
		// Remove from the reactors list
		p.reactorsLock.Lock()
		delete(p.reactors, digest)
		p.reactorsLock.Unlock()

		p.logger.Debug("reactor finalized and removed from manager", zap.Stringer("digest", digest))
		go p.handler.HandleFinalization(reactor.observation, reactor.VAASignatures())
	case StateTimedOut:
		// Remove from the reactors list
		p.reactorsLock.Lock()
		delete(p.reactors, digest)
		p.reactorsLock.Unlock()

		p.logger.Debug("reactor timed out and removed from manager", zap.Stringer("digest", digest))
		go p.handler.HandleTimeout(oldState, reactor.observation, reactor.VAASignatures())
	}
}

func (p *Manager[K]) GuardianSet() *common.GuardianSet {
	return p.gs
}
