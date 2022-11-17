package processor

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/certusone/wormhole/node/pkg/common"
	"github.com/certusone/wormhole/node/pkg/db"
	"github.com/certusone/wormhole/node/pkg/governor"
	"github.com/certusone/wormhole/node/pkg/processor/reactor"
	gossipv1 "github.com/certusone/wormhole/node/pkg/proto/gossip/v1"
	"github.com/certusone/wormhole/node/pkg/reporter"
	"github.com/certusone/wormhole/node/pkg/supervisor"
	"github.com/wormhole-foundation/wormhole/sdk/vaa"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type VAAReactor struct {
	// lockC is a channel of observed emitted messages
	lockC chan *common.MessagePublication
	// obsC is a channel of VAAs used to forward to the reactor manager
	obsC chan *vaa.VAA
	// obsvReqSendC is a send-only channel of outbound re-observation requests to broadcast on p2p
	obsvReqSendC chan<- *gossipv1.ObservationRequest
	// signedInC is a channel of inbound signed VAA observations from p2p
	signedInC chan *gossipv1.SignedVAAWithQuorum
	// sendC is a channel of outbound messages to broadcast on p2p
	sendC chan []byte

	attestationEvents *reporter.AttestationEventReporter
	logger            *zap.Logger
	db                *db.Database

	manager *reactor.Manager[*vaa.VAA]

	// gst is managed by the processor and allows concurrent access to the
	// guardian set by other components.
	gst *common.GuardianSetState

	governor    *governor.ChainGovernor
	pythnetVaas map[string]PythNetVaaEntry
}

func NewVAAReactor(
	db *db.Database,
	lockC chan *common.MessagePublication,
	setC chan *common.GuardianSet,
	sendC chan []byte,
	obsvC chan *gossipv1.SignedObservation,
	obsvReqSendC chan<- *gossipv1.ObservationRequest,
	signedInC chan *gossipv1.SignedVAAWithQuorum,
	gk *ecdsa.PrivateKey,
	gst *common.GuardianSetState,
	attestationEvents *reporter.AttestationEventReporter,
	g *governor.ChainGovernor,
) *VAAReactor {
	obsC := make(chan *vaa.VAA, 10)

	r := &VAAReactor{
		obsvReqSendC:      obsvReqSendC,
		signedInC:         signedInC,
		lockC:             lockC,
		obsC:              obsC,
		sendC:             sendC,
		gst:               gst,
		db:                db,
		attestationEvents: attestationEvents,
		governor:          g,
	}

	manager := reactor.NewManager[*vaa.VAA](obsC, obsvC, setC, gst, reactor.Config{
		RetransmitFrequency: 30 * time.Second,
		QuorumGracePeriod:   30 * time.Second,
		QuorumTimeout:       5 * time.Minute,
		UnobservedTimeout:   5 * time.Minute,
		Signer:              reactor.NewEcdsaKeySigner(gk),
		NetworkAdapter:      reactor.NewChannelNetworkAdapter(sendC),
	}, r)
	r.manager = manager

	return r
}

func (p *VAAReactor) Run(ctx context.Context) error {
	p.logger = supervisor.Logger(ctx)

	err := supervisor.Run(ctx, "vaa-reactor-manager", p.manager.Run)
	if err != nil {
		return fmt.Errorf("failed to start reactor manager: %w", err)
	}

	cleanup := time.NewTicker(30 * time.Second)
	defer cleanup.Stop()

	// Always initialize the timer so don't have a nil pointer in the case below. It won't get rearmed after that.
	govTimer := time.NewTicker(time.Minute)
	defer govTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case k := <-p.lockC:
			p.logger.Info("saw new VAA")
			if p.governor != nil {
				if !p.governor.ProcessMsg(k) {
					continue
				}
			}
			gs := p.manager.GuardianSet()
			if gs == nil {
				p.logger.Warn("received observation before guardian set was known - skipping")
				continue
			}
			v := MessagePublicationToVAA(k, gs)
			p.attestationEvents.ReportMessagePublication(&reporter.MessagePublication{
				VAA:            *v,
				InitiatingTxID: k.TxHash,
			})
			p.obsC <- v
		case <-cleanup.C:
		// TODO p.handleCleanup(ctx)
		case s := <-p.signedInC:
			p.handleSignedVAA(s)
		case <-govTimer.C:
			if p.governor != nil {
				toBePublished, err := p.governor.CheckPending()
				if err != nil {
					return err
				}
				if len(toBePublished) != 0 {
					for _, k := range toBePublished {
						gs := p.manager.GuardianSet()
						if gs == nil {
							p.logger.Warn("received observation before guardian set was known - skipping")
							continue
						}
						p.obsC <- MessagePublicationToVAA(k, gs)
					}
				}
			}
		}
	}
}

func (p *VAAReactor) handleSignedVAA(m *gossipv1.SignedVAAWithQuorum) {
	v, err := vaa.Unmarshal(m.Vaa)
	if err != nil {
		p.logger.Warn("received invalid VAA in SignedVAAWithQuorum message",
			zap.Error(err), zap.Any("message", m))
		return
	}

	// Calculate digest for logging
	digest := v.SigningMsg()
	hash := hex.EncodeToString(digest.Bytes())

	gs := p.manager.GuardianSet()
	if gs == nil {
		p.logger.Warn("dropping SignedVAAWithQuorum message since we haven't initialized our guardian set yet",
			zap.String("digest", hash),
			zap.Any("message", m),
		)
		return
	}

	if err := v.Verify(gs.Keys); err != nil {
		p.logger.Warn("dropping SignedVAAWithQuorum message because it failed verification", zap.Error(err))
		return
	}

	// We now established that:
	//  - all signatures on the VAA are valid
	//  - the signature's addresses match the node's current guardian set
	//  - enough signatures are present for the VAA to reach quorum

	// Check if we already store this VAA
	_, err = p.getSignedVAA(*db.VaaIDFromVAA(v))
	if err == nil {
		p.logger.Debug("ignored SignedVAAWithQuorum message for VAA we already store",
			zap.String("digest", hash),
		)
		return
	} else if err != db.ErrVAANotFound {
		p.logger.Error("failed to look up VAA in database",
			zap.String("digest", hash),
			zap.Error(err),
		)
		return
	}

	// Store signed VAA in database.
	p.logger.Info("storing inbound signed VAA with quorum",
		zap.String("digest", hash),
		zap.Any("vaa", v),
		zap.String("bytes", hex.EncodeToString(m.Vaa)),
		zap.String("message_id", v.MessageID()))

	if err := p.storeSignedVAA(v); err != nil {
		p.logger.Error("failed to store signed VAA", zap.Error(err))
		return
	}
	p.attestationEvents.ReportVAAQuorum(v)
}

func (p *VAAReactor) HandleQuorum(observation *vaa.VAA, signatures []*vaa.Signature) {
	observation.Signatures = signatures

	vaaBytes, err := observation.Marshal()
	if err != nil {
		panic(err)
	}

	// Store signed VAA in database.
	p.logger.Info("signed VAA with quorum",
		zap.String("digest", observation.SigningMsg().Hex()),
		zap.Any("vaa", observation),
		zap.String("bytes", hex.EncodeToString(vaaBytes)),
		zap.String("message_id", observation.MessageID()))

	if err := p.storeSignedVAA(observation); err != nil {
		p.logger.Error("failed to store signed VAA", zap.Error(err))
	}

	p.broadcastSignedVAA(observation)
	p.attestationEvents.ReportVAAQuorum(observation)
}

func (p *VAAReactor) HandleFinalization(observation *vaa.VAA, signatures []*vaa.Signature) {
	// TODO store with full signatures
}

func (p *VAAReactor) HandleTimeout(previousState reactor.State, observation *vaa.VAA, signatures []*vaa.Signature) {
	// TODO store all signatures still?
}

type PythNetVaaEntry struct {
	v          *vaa.VAA
	updateTime time.Time // Used for determining when to delete entries
}

func (p *VAAReactor) storeSignedVAA(v *vaa.VAA) error {
	if v.EmitterChain == vaa.ChainIDPythNet {
		key := fmt.Sprintf("%v/%v", v.EmitterAddress, v.Sequence)
		p.pythnetVaas[key] = PythNetVaaEntry{v: v, updateTime: time.Now()}
		return nil
	}
	return p.db.StoreSignedVAA(v)
}

func (p *VAAReactor) getSignedVAA(id db.VAAID) (*vaa.VAA, error) {
	if id.EmitterChain == vaa.ChainIDPythNet {
		key := fmt.Sprintf("%v/%v", id.EmitterAddress, id.Sequence)
		ret, exists := p.pythnetVaas[key]
		if exists {
			return ret.v, nil
		}

		return nil, db.ErrVAANotFound
	}

	vb, err := p.db.GetSignedVAABytes(id)
	if err != nil {
		return nil, err
	}

	v, err := vaa.Unmarshal(vb)
	if err != nil {
		panic("failed to unmarshal VAA from db")
	}

	return v, err
}

func (p *VAAReactor) broadcastSignedVAA(v *vaa.VAA) {
	b, err := v.Marshal()
	if err != nil {
		panic(err)
	}

	w := gossipv1.GossipMessage{Message: &gossipv1.GossipMessage_SignedVaaWithQuorum{
		SignedVaaWithQuorum: &gossipv1.SignedVAAWithQuorum{Vaa: b},
	}}

	msg, err := proto.Marshal(&w)
	if err != nil {
		panic(err)
	}

	p.sendC <- msg
}

func MessagePublicationToVAA(k *common.MessagePublication, gs *common.GuardianSet) *vaa.VAA {
	return &vaa.VAA{
		Version:          vaa.SupportedVAAVersion,
		GuardianSetIndex: gs.Index,
		Signatures:       nil,
		Timestamp:        k.Timestamp,
		Nonce:            k.Nonce,
		EmitterChain:     k.EmitterChain,
		EmitterAddress:   k.EmitterAddress,
		Payload:          k.Payload,
		Sequence:         k.Sequence,
		ConsistencyLevel: k.ConsistencyLevel,
	}
}
