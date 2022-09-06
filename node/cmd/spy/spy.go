package spy

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/certusone/wormhole/node/pkg/common"
	"github.com/certusone/wormhole/node/pkg/p2p"
	gossipv1 "github.com/certusone/wormhole/node/pkg/proto/gossip/v1"
	spyv1 "github.com/certusone/wormhole/node/pkg/proto/spy/v1"
	"github.com/certusone/wormhole/node/pkg/supervisor"
	"github.com/certusone/wormhole/node/pkg/vaa"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	ipfslog "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	rootCtx       context.Context
	rootCtxCancel context.CancelFunc
)

var (
	p2pNetworkID *string
	p2pPort      *uint
	p2pBootstrap *string

	statusAddr *string

	nodeKeyPath *string

	logLevel *string

	spyRPC *string
)

func init() {
	p2pNetworkID = SpyCmd.Flags().String("network", "/wormhole/dev", "P2P network identifier")
	p2pPort = SpyCmd.Flags().Uint("port", 8999, "P2P UDP listener port")
	p2pBootstrap = SpyCmd.Flags().String("bootstrap", "", "P2P bootstrap peers (comma-separated)")

	statusAddr = SpyCmd.Flags().String("statusAddr", "[::]:6060", "Listen address for status server (disabled if blank)")

	nodeKeyPath = SpyCmd.Flags().String("nodeKey", "", "Path to node key (will be generated if it doesn't exist)")

	logLevel = SpyCmd.Flags().String("logLevel", "info", "Logging level (debug, info, warn, error, dpanic, panic, fatal)")

	spyRPC = SpyCmd.Flags().String("spyRPC", "", "Listen address for gRPC interface")
}

// SpyCmd represents the node command
var SpyCmd = &cobra.Command{
	Use:   "spy",
	Short: "Run gossip spy client",
	Run:   runSpy,
}

type spyServer struct {
	spyv1.UnimplementedSpyRPCServiceServer
	logger *zap.Logger
	subs   map[string]*subscription
	subsMu sync.Mutex
}

type message struct {
	vaaBytes []byte
}

type filter struct {
	chainId     vaa.ChainID
	emitterAddr vaa.Address
}

type subscription struct {
	filters []filter
	ch      chan message
}

func subscriptionId() string {
	return uuid.New().String()
}

func decodeEmitterAddr(hexAddr string) (vaa.Address, error) {
	address, err := hex.DecodeString(hexAddr)
	if err != nil {
		return vaa.Address{}, status.Error(codes.InvalidArgument, fmt.Sprintf("failed to decode address: %v", err))
	}
	if len(address) != 32 {
		return vaa.Address{}, status.Error(codes.InvalidArgument, "address must be 32 bytes")
	}

	addr := vaa.Address{}
	copy(addr[:], address)

	return addr, nil
}

func (s *spyServer) Publish(vaaBytes []byte) error {
	s.subsMu.Lock()
	defer s.subsMu.Unlock()

	var v *vaa.VAA

	for _, sub := range s.subs {
		if len(sub.filters) == 0 {
			sub.ch <- message{vaaBytes: vaaBytes}
		} else {
			if v == nil {
				var err error
				v, err = vaa.Unmarshal(vaaBytes)
				if err != nil {
					return err
				}
			}

			for _, fi := range sub.filters {
				if fi.chainId == v.EmitterChain && fi.emitterAddr == v.EmitterAddress {
					sub.ch <- message{vaaBytes: vaaBytes}
				}
			}
		}
	}

	return nil
}

func (s *spyServer) SubscribeSignedVAA(req *spyv1.SubscribeSignedVAARequest, resp spyv1.SpyRPCService_SubscribeSignedVAAServer) error {
	var fi []filter
	if req.Filters != nil {
		for _, f := range req.Filters {
			switch t := f.Filter.(type) {
			case *spyv1.FilterEntry_EmitterFilter:
				addr, err := decodeEmitterAddr(t.EmitterFilter.EmitterAddress)
				if err != nil {
					return status.Error(codes.InvalidArgument, fmt.Sprintf("failed to decode emitter address: %v", err))
				}
				fi = append(fi, filter{
					chainId:     vaa.ChainID(t.EmitterFilter.ChainId),
					emitterAddr: addr,
				})
			default:
				return status.Error(codes.InvalidArgument, "unsupported filter type")
			}
		}
	}

	s.subsMu.Lock()
	id := subscriptionId()
	sub := &subscription{
		ch:      make(chan message, 1),
		filters: fi,
	}
	s.subs[id] = sub
	s.subsMu.Unlock()

	defer func() {
		s.subsMu.Lock()
		defer s.subsMu.Unlock()
		delete(s.subs, id)
	}()

	for {
		select {
		case <-resp.Context().Done():
			return resp.Context().Err()
		case msg := <-sub.ch:
			if err := resp.Send(&spyv1.SubscribeSignedVAAResponse{
				VaaBytes: msg.vaaBytes,
			}); err != nil {
				return err
			}
		}
	}
}

func newSpyServer(logger *zap.Logger) *spyServer {
	return &spyServer{
		logger: logger.Named("spyserver"),
		subs:   make(map[string]*subscription),
	}
}

func spyServerRunnable(s *spyServer, logger *zap.Logger, listenAddr string) (supervisor.Runnable, *grpc.Server, error) {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen: %w", err)
	}

	logger.Info("publicrpc server listening", zap.String("addr", l.Addr().String()))

	grpcServer := common.NewInstrumentedGRPCServer(logger)
	spyv1.RegisterSpyRPCServiceServer(grpcServer, s)

	return supervisor.GRPCServer(grpcServer, l, false), grpcServer, nil
}

func runSpy(cmd *cobra.Command, args []string) {
	common.SetRestrictiveUmask()

	lvl, err := ipfslog.LevelFromString(*logLevel)
	if err != nil {
		fmt.Println("Invalid log level")
		os.Exit(1)
	}

	logger := ipfslog.Logger("wormhole-spy").Desugar()

	ipfslog.SetAllLoggers(lvl)

	// Status server
	if *statusAddr != "" {
		router := mux.NewRouter()

		router.Handle("/metrics", promhttp.Handler())

		go func() {
			logger.Info("status server listening on [::]:6060")
			logger.Error("status server crashed", zap.Error(http.ListenAndServe(*statusAddr, router)))
		}()
	}

	// Verify flags

	if *nodeKeyPath == "" {
		logger.Fatal("Please specify --nodeKey")
	}
	if *p2pBootstrap == "" {
		logger.Fatal("Please specify --bootstrap")
	}

	// Node's main lifecycle context.
	rootCtx, rootCtxCancel = context.WithCancel(context.Background())
	defer rootCtxCancel()

	// Outbound gossip message queue
	sendC := make(chan []byte)

	// Inbound observations
	obsvC := make(chan *gossipv1.SignedObservation, 50)

	// Inbound signed VAAs
	signedInC := make(chan *gossipv1.SignedVAAWithQuorum, 50)

	// Heartbeat updates
	heartbeatC := make(chan *gossipv1.Heartbeat, 50)

	// Guardian set state managed by processor
	gst := common.NewGuardianSetState(heartbeatC)

	// Bootstrap guardian set, otherwise heartbeats would be skipped
	// gst.Set(&common.GuardianSet{
	// 	Index: 2,
	// 	Keys: []eth_common.Address{
	// 		eth_common.HexToAddress("0x58CC3AE5C097b213cE3c81979e1B9f9570746AA5"),
	// 		eth_common.HexToAddress("0xfF6CB952589BDE862c25Ef4392132fb9D4A42157"),
	// 		eth_common.HexToAddress("0x114De8460193bdf3A2fCf81f86a09765F4762fD1"),
	// 		eth_common.HexToAddress("0x107A0086b32d7A0977926A205131d8731D39cbEB"),
	// 		eth_common.HexToAddress("0x8C82B2fd82FaeD2711d59AF0F2499D16e726f6b2"),
	// 		eth_common.HexToAddress("0x11b39756C042441BE6D8650b69b54EbE715E2343"),
	// 		eth_common.HexToAddress("0x54Ce5B4D348fb74B958e8966e2ec3dBd4958a7cd"),
	// 		eth_common.HexToAddress("0x66B9590e1c41e0B226937bf9217D1d67Fd4E91F5"),
	// 		eth_common.HexToAddress("0x74a3bf913953D695260D88BC1aA25A4eeE363ef0"),
	// 		eth_common.HexToAddress("0x000aC0076727b35FBea2dAc28fEE5cCB0fEA768e"),
	// 		eth_common.HexToAddress("0xAF45Ced136b9D9e24903464AE889F5C8a723FC14"),
	// 		eth_common.HexToAddress("0xf93124b7c738843CBB89E864c862c38cddCccF95"),
	// 		eth_common.HexToAddress("0xD2CC37A4dc036a8D232b48f62cDD4731412f4890"),
	// 		eth_common.HexToAddress("0xDA798F6896A3331F64b48c12D1D57Fd9cbe70811"),
	// 		eth_common.HexToAddress("0x71AA1BE1D36CaFE3867910F99C09e347899C19C3"),
	// 		eth_common.HexToAddress("0x8192b6E7387CCd768277c17DAb1b7a5027c0b3Cf"),
	// 		eth_common.HexToAddress("0x178e21ad2E77AE06711549CFBB1f9c7a9d8096e8"),
	// 		eth_common.HexToAddress("0x5E1487F35515d02A92753504a8D75471b9f49EdB"),
	// 		eth_common.HexToAddress("0x6FbEBc898F403E4773E95feB15E80C9A99c8348d"),
	// 	},
	// })

	// RPC server
	s := newSpyServer(logger)
	rpcSvc, _, err := spyServerRunnable(s, logger, *spyRPC)
	if err != nil {
		logger.Fatal("failed to start RPC server", zap.Error(err))
	}

	// Ignore observations
	go func() {
		for {
			select {
			case <-rootCtx.Done():
				return
			case <-obsvC:
			}
		}
	}()

	// Log signed VAAs
	go func() {
		for {
			select {
			case <-rootCtx.Done():
				return
			case v := <-signedInC:
				logger.Info("Received signed VAA",
					zap.Any("vaa", v.Vaa))
				if err := s.Publish(v.Vaa); err != nil {
					logger.Error("failed to publish signed VAA", zap.Error(err))
				}
			}
		}
	}()

	// Ignore heartbeats
	go func() {
		for {
			select {
			case <-rootCtx.Done():
				return
			case <- heartbeatC:
			}
		}
	}()

	// Load p2p private key
	var priv crypto.PrivKey
	priv, err = common.GetOrCreateNodeKey(logger, *nodeKeyPath)
	if err != nil {
		logger.Fatal("Failed to load node key", zap.Error(err))
	}

	// Run supervisor.
	supervisor.New(rootCtx, logger, func(ctx context.Context) error {
		if err := supervisor.Run(ctx, "p2p", p2p.Run(obsvC, nil, nil, sendC, signedInC, priv, nil, gst, *p2pPort, *p2pNetworkID, *p2pBootstrap, "", false, rootCtxCancel, nil)); err != nil {
			return err
		}

		if err := supervisor.Run(ctx, "spyrpc", rpcSvc); err != nil {
			return err
		}

		logger.Info("Started internal services")

		<-ctx.Done()
		return nil
	},
		// It's safer to crash and restart the process in case we encounter a panic,
		// rather than attempting to reschedule the runnable.
		supervisor.WithPropagatePanic)

	<-rootCtx.Done()
	logger.Info("root context cancelled, exiting...")
	// TODO: wait for things to shut down gracefully
}
