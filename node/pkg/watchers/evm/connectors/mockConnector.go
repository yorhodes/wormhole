package connectors

import (
	"context"
	"encoding/json"
	"fmt"

	ethAbi "github.com/certusone/wormhole/node/pkg/watchers/evm/connectors/ethabi"

	ethereum "github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	ethClient "github.com/ethereum/go-ethereum/ethclient"
	ethEvent "github.com/ethereum/go-ethereum/event"

	"go.uber.org/zap"
)

// MockConnector implements the connector interface for testing purposes.
type MockConnector struct {
	address ethCommon.Address
	client  *ethClient.Client
	results []string
	err     error
}

// SetResults takes an array of json results strings. Each time a test makes an RPC call, it uses the first
// element in the array as the response, and the discards it. If the array is empty, an error will be returned.
func (m *MockConnector) SetResults(results []string) {
	m.results = results
}

// SetError takes an error which will be returned on the next RPC call. After that, the error is reset to nil.
func (m *MockConnector) SetError(err error) {
	m.err = err
}

func NewMockConnector(ctx context.Context, networkName, rawUrl string, address ethCommon.Address, logger *zap.Logger) (*MockConnector, error) {
	return nil, fmt.Errorf("not implemented")
}

func (e *MockConnector) NetworkName() string {
	return "MockConnector"
}

func (e *MockConnector) ContractAddress() ethCommon.Address {
	return e.address
}

func (e *MockConnector) GetCurrentGuardianSetIndex(ctx context.Context) (uint32, error) {
	return 0, fmt.Errorf("not implemented")
}

func (e *MockConnector) GetGuardianSet(ctx context.Context, index uint32) (ethAbi.StructsGuardianSet, error) {
	return ethAbi.StructsGuardianSet{}, fmt.Errorf("not implemented")
}

func (e *MockConnector) WatchLogMessagePublished(ctx context.Context, sink chan<- *ethAbi.AbiLogMessagePublished) (ethEvent.Subscription, error) {
	var s ethEvent.Subscription
	return s, fmt.Errorf("not implemented")
}

func (e *MockConnector) TransactionReceipt(ctx context.Context, txHash ethCommon.Hash) (*ethTypes.Receipt, error) {
	return nil, fmt.Errorf("not implemented")
}

func (e *MockConnector) TimeOfBlockByHash(ctx context.Context, hash ethCommon.Hash) (uint64, error) {
	return 0, fmt.Errorf("not implemented")
}

func (e *MockConnector) ParseLogMessagePublished(log ethTypes.Log) (*ethAbi.AbiLogMessagePublished, error) {
	return nil, fmt.Errorf("not implemented")
}

func (e *MockConnector) SubscribeForBlocks(ctx context.Context, sink chan<- *NewBlock) (ethereum.Subscription, error) {
	var s ethEvent.Subscription
	return s, fmt.Errorf("not implemented")
}

func (e *MockConnector) RawCallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	if e.err != nil {
		err := e.err
		e.err = nil
		return err
	}
	if len(e.results) == 0 {
		return fmt.Errorf("no more results")
	}

	str := e.results[0]
	e.results = e.results[1:]
	return json.Unmarshal([]byte(str), &result)
}

func (e *MockConnector) Client() *ethClient.Client {
	return e.client
}
