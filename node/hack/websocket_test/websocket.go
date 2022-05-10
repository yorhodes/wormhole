package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/rpc"

	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"log"

	"github.com/certusone/wormhole/node/pkg/common"
	eth "github.com/certusone/wormhole/node/pkg/ethereum"
	"github.com/certusone/wormhole/node/pkg/ethereum/abi"
	"github.com/certusone/wormhole/node/pkg/vaa"
)

type (
	Watcher struct {
		// Ethereum RPC url
		url string
		// Address of the Eth contract contract
		contract eth_common.Address
		// Human-readable name of the Eth network, for logging and monitoring.
		networkName string
		// VAA ChainID of the network we're connecting to.
		chainID vaa.ChainID

		pending   map[pendingKey]*pendingMessage
		pendingMu sync.Mutex

		// 0 is a valid guardian set, so we need a nil value here
		currentGuardianSet *uint32

		// Minimum number of confirmations to accept, regardless of what the contract specifies.
		minConfirmations uint64

		// Interface to the chain specific ethereum library.
		ethIntf             common.Ethish
	}

	pendingKey struct {
		TxHash         eth_common.Hash
		BlockHash      eth_common.Hash
		EmitterAddress vaa.Address
		Sequence       uint64
	}

	pendingMessage struct {
		message *common.MessagePublication
		height  uint64
	}
)

func main() {
	ctx := context.Background()
	e := &Watcher{
		url:                 "wss://ws-matic-mainnet.chainstacklabs.com",
		contract:            eth_common.HexToAddress("0x7A4B5a56256163F07b2C80A7cA55aBE66c4ec4d7"),
		networkName:         "polygon",
		minConfirmations:    512,
		chainID:             5,
		pending:             map[pendingKey]*pendingMessage{},
		ethIntf:             &eth.EthImpl{NetworkName: "polygon"}}

	timeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	err := e.ethIntf.DialContext(timeout, e.url)
	if err != nil {
		log.Print(fmt.Errorf("dialing eth client failed: %w", err))
		return
	}

	err = e.ethIntf.NewAbiFilterer(e.contract)
	if err != nil {
		log.Print(fmt.Errorf("could not create wormhole contract filter: %w", err))
		return
	}

	err = e.ethIntf.NewAbiCaller(e.contract)
	if err != nil {
		panic(err)
	}

	// Timeout for initializing subscriptions
	timeout, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Subscribe to new message publications
	messageC := make(chan *abi.AbiLogMessagePublished, 2)
	messageSub, err := e.ethIntf.WatchLogMessagePublished(ctx, timeout, messageC)
	if err != nil {
		log.Print(fmt.Errorf("failed to subscribe to message publication events: %w", err))
		return
	}

	// Fetch initial guardian set
	if err := e.fetchAndUpdateGuardianSet(ctx, e.ethIntf); err != nil {
		log.Print(fmt.Errorf("failed to request guardian set: %v", err))
		return
	}

	errC := make(chan error)

	// Poll for guardian set.
	go func() {
		t := time.NewTicker(15 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if err := e.fetchAndUpdateGuardianSet(ctx, e.ethIntf); err != nil {
					errC <- fmt.Errorf("failed to request guardian set: %v", err)
					return
				}
			}
		}
	}()

	// Track the current block number so we can compare it to the block number of
	// the message publication for observation requests.
	var currentBlockNumber uint64

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-messageSub.Err():
				errC <- fmt.Errorf("error while processing message publication subscription: %w", err)
				return
			case ev := <-messageC:
				// Request timestamp for block
				timeout, cancel := context.WithTimeout(ctx, 15*time.Second)
				blockTime, err := e.ethIntf.TimeOfBlockByHash(timeout, ev.Raw.BlockHash)
				cancel()
				
				if err != nil {
					errC <- fmt.Errorf("failed to request timestamp for block %d: %w", ev.Raw.BlockNumber, err)
					return
				}

				message := &common.MessagePublication{
					TxHash:           ev.Raw.TxHash,
					Timestamp:        time.Unix(int64(blockTime), 0),
					Nonce:            ev.Nonce,
					Sequence:         ev.Sequence,
					EmitterChain:     e.chainID,
					EmitterAddress:   eth.PadAddress(ev.Sender),
					Payload:          ev.Payload,
					ConsistencyLevel: ev.ConsistencyLevel,
				}

				log.Print("found new message publication transaction",
					zap.Stringer("tx", ev.Raw.TxHash),
					zap.Uint64("block", ev.Raw.BlockNumber),
					zap.Stringer("blockhash", ev.Raw.BlockHash),
					zap.String("eth_network", e.networkName))

				key := pendingKey{
					TxHash:         message.TxHash,
					BlockHash:      ev.Raw.BlockHash,
					EmitterAddress: message.EmitterAddress,
					Sequence:       message.Sequence,
				}

				e.pendingMu.Lock()
				e.pending[key] = &pendingMessage{
					message: message,
					height:  ev.Raw.BlockNumber,
				}
				e.pendingMu.Unlock()
			}
		}
	}()

	// Watch headers
	headSink := make(chan *types.Header, 2)
	headerSubscription, err := e.ethIntf.SubscribeNewHead(ctx, headSink)
	if err != nil {
		log.Print(fmt.Errorf("failed to subscribe to header events: %w", err))
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-headerSubscription.Err():
				errC <- fmt.Errorf("error while processing header subscription: %w", err)
				return
			case ev := <-headSink:
				if ev == nil {
					log.Print("new header event is nil", zap.String("eth_network", e.networkName))
					continue
				}

				start := time.Now()
				currentHash := ev.Hash()
				log.Print("processing new header",
					zap.Stringer("current_block", ev.Number),
					zap.Stringer("current_blockhash", currentHash),
					zap.String("eth_network", e.networkName))
				
				e.pendingMu.Lock()

				blockNumberU := ev.Number.Uint64()
				atomic.StoreUint64(&currentBlockNumber, blockNumberU)

				for key, pLock := range e.pending {
					expectedConfirmations := uint64(pLock.message.ConsistencyLevel)
					if expectedConfirmations < e.minConfirmations {
						expectedConfirmations = e.minConfirmations
					}

					// Transaction was dropped and never picked up again
					if pLock.height+4*uint64(expectedConfirmations) <= blockNumberU {
						log.Print("observation timed out",
							zap.Stringer("tx", pLock.message.TxHash),
							zap.Stringer("blockhash", key.BlockHash),
							zap.Stringer("emitter_address", key.EmitterAddress),
							zap.Uint64("sequence", key.Sequence),
							zap.Stringer("current_block", ev.Number),
							zap.Stringer("current_blockhash", currentHash),
							zap.String("eth_network", e.networkName),
						)
						delete(e.pending, key)
						continue
					}

					// Transaction is now ready
					if pLock.height+uint64(expectedConfirmations) <= blockNumberU {
						timeout, cancel := context.WithTimeout(ctx, 5*time.Second)
						tx, err := e.ethIntf.TransactionReceipt(timeout, pLock.message.TxHash)
						cancel()

						// If the node returns an error after waiting expectedConfirmation blocks,
						// it means the chain reorged and the transaction was orphaned. The
						// TransactionReceipt call is using the same websocket connection than the
						// head notifications, so it's guaranteed to be atomic.
						//
						// Check multiple possible error cases - the node seems to return a
						// "not found" error most of the time, but it could conceivably also
						// return a nil tx or rpc.ErrNoResult.
						if tx == nil || err == rpc.ErrNoResult || (err != nil && err.Error() == "not found") {
							log.Print("tx was orphaned",
								zap.Stringer("tx", pLock.message.TxHash),
								zap.Stringer("blockhash", key.BlockHash),
								zap.Stringer("emitter_address", key.EmitterAddress),
								zap.Uint64("sequence", key.Sequence),
								zap.Stringer("current_block", ev.Number),
								zap.Stringer("current_blockhash", currentHash),
								zap.String("eth_network", e.networkName),
								zap.Error(err))
							delete(e.pending, key)
							continue
						}

						// This should never happen - if we got this far, it means that logs were emitted,
						// which is only possible if the transaction succeeded. We check it anyway just
						// in case the EVM implementation is buggy.
						if tx.Status != 1 {
							log.Print("transaction receipt with non-success status",
								zap.Stringer("tx", pLock.message.TxHash),
								zap.Stringer("blockhash", key.BlockHash),
								zap.Stringer("emitter_address", key.EmitterAddress),
								zap.Uint64("sequence", key.Sequence),
								zap.Stringer("current_block", ev.Number),
								zap.Stringer("current_blockhash", currentHash),
								zap.String("eth_network", e.networkName),
								zap.Error(err))
							delete(e.pending, key)
							continue
						}

						// Any error other than "not found" is likely transient - we retry next block.
						if err != nil {
							log.Print("transaction could not be fetched",
								zap.Stringer("tx", pLock.message.TxHash),
								zap.Stringer("blockhash", key.BlockHash),
								zap.Stringer("emitter_address", key.EmitterAddress),
								zap.Uint64("sequence", key.Sequence),
								zap.Stringer("current_block", ev.Number),
								zap.Stringer("current_blockhash", currentHash),
								zap.String("eth_network", e.networkName),
								zap.Error(err))
							continue
						}

						// It's possible for a transaction to be orphaned and then included in a different block
						// but with the same tx hash. Drop the observation (it will be re-observed and needs to
						// wait for the full confirmation time again).
						if tx.BlockHash != key.BlockHash {
							log.Print("tx got dropped and mined in a different block; the message should have been reobserved",
								zap.Stringer("tx", pLock.message.TxHash),
								zap.Stringer("blockhash", key.BlockHash),
								zap.Stringer("emitter_address", key.EmitterAddress),
								zap.Uint64("sequence", key.Sequence),
								zap.Stringer("current_block", ev.Number),
								zap.Stringer("current_blockhash", currentHash),
								zap.String("eth_network", e.networkName))
							delete(e.pending, key)
							continue
						}

						log.Print("observation confirmed",
							zap.Stringer("tx", pLock.message.TxHash),
							zap.Stringer("blockhash", key.BlockHash),
							zap.Stringer("emitter_address", key.EmitterAddress),
							zap.Uint64("sequence", key.Sequence),
							zap.Stringer("current_block", ev.Number),
							zap.Stringer("current_blockhash", currentHash),
							zap.String("eth_network", e.networkName))
						delete(e.pending, key)
						// e.msgChan <- pLock.message
					}
				}

				e.pendingMu.Unlock()
				log.Print("processed new header",
					zap.Stringer("current_block", ev.Number),
					zap.Stringer("current_blockhash", currentHash),
					zap.Duration("took", time.Since(start)),
					zap.String("eth_network", e.networkName))
			}
		}
	}()

	select {
	case <-ctx.Done():
		ctx.Err()
		return
	case err := <-errC:
		log.Print(err.Error())
		log.Printf("Dumping %d pending messages", len(e.pending))
		return
	}
}

func (e *Watcher) fetchAndUpdateGuardianSet(
	ctx context.Context,
	ethIntf common.Ethish,
) error {
	log.Printf("fetching guardian set")
	timeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	idx, gs, err := fetchCurrentGuardianSet(timeout, ethIntf)
	if err != nil {
		return err
	}

	if e.currentGuardianSet != nil && *(e.currentGuardianSet) == idx {
		return nil
	}

	log.Print("updated guardian set found",
		zap.Any("value", gs), zap.Uint32("index", idx),
		zap.String("eth_network", e.networkName))

	e.currentGuardianSet = &idx

	// if e.setChan != nil {
	// 	e.setChan <- &common.GuardianSet{
	// 		Keys:  gs.Keys,
	// 		Index: idx,
	// 	}
	// }

	return nil
}

// Fetch the current guardian set ID and guardian set from the chain.
func fetchCurrentGuardianSet(ctx context.Context, ethIntf common.Ethish) (uint32, *abi.StructsGuardianSet, error) {
	currentIndex, err := ethIntf.GetCurrentGuardianSetIndex(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("error requesting current guardian set index: %w", err)
	}

	gs, err := ethIntf.GetGuardianSet(ctx, currentIndex)
	if err != nil {
		return 0, nil, fmt.Errorf("error requesting current guardian set value: %w", err)
	}

	return currentIndex, &gs, nil
}

func (e *Watcher) checkForSafeMode(ctx context.Context) error {
	timeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	c, err := rpc.DialContext(timeout, e.url)
	if err != nil {
		return fmt.Errorf("failed to connect to url %s to check for safe mode: %w", e.url, err)
	}

	var safe bool
	err = c.CallContext(ctx, &safe, "net_isSafeMode")
	if err != nil {
		return fmt.Errorf("check for safe mode for url %s failed: %w", e.url, err)
	}

	if !safe {
		return fmt.Errorf("url %s is not using safe mode", e.url)
	}

	return nil
}
