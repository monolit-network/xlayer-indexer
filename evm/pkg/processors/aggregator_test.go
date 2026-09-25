package processors

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/processors/mockProcessor"
	"github.com/monolit-network/xlayer-indexer/util/promise"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

type MockEVMClient struct {
	logger            *zap.Logger
	totalTxsGenerated int
}

func (c *MockEVMClient) BatchGetBlocksInfo(ctx context.Context, startBlock *big.Int, endBlock *big.Int, requiredDataTypes evmclient.RequiredDataTypes) *evmclient.BlocksInfo {
	result := &evmclient.BlocksInfo{
		Blocks:   promise.NewPromise[[]*types.Block](),
		Receipts: promise.NewPromise[[]evmclient.BlockReceipts](),
		Traces:   promise.NewPromise[[]evmclient.BlockTraces](),
		Headers:  promise.NewPromise[[]*types.Header](),
	}
	return result
}

func (c *MockEVMClient) GetTotalTxsGenerated() int {
	return c.totalTxsGenerated
}

func TestProcessorAggregator(t *testing.T) {
	processors := []IndexerProcessor{
		mockProcessor.NewMockProcessor(zap.NewNop(), evmclient.RequiredDataTypesTx|evmclient.RequiredDataTypesReceipt|evmclient.RequiredDataTypesBlockHeader|evmclient.RequiredDataTypesTraces),
	}

	loggerCfg := zap.NewProductionConfig()
	loggerCfg.OutputPaths = []string{"stdout", "stderr"}
	loggerCfg.ErrorOutputPaths = []string{"stdout", "stderr"}
	logger, _ := loggerCfg.Build()

	aggregator := NewProcessorAggregator(processors, &MockEVMClient{logger: logger}, logger)

	aggregator.ProcessBlockRange(context.Background(), big.NewInt(1), big.NewInt(10))
	aggregator.ProcessBlockRange(context.Background(), big.NewInt(11), big.NewInt(20))
}

type staticEVMClient struct {
	requiredData evmclient.RequiredDataTypes
}

func (c *staticEVMClient) BatchGetBlocksInfo(ctx context.Context, startBlock *big.Int, endBlock *big.Int, requiredDataTypes evmclient.RequiredDataTypes) *evmclient.BlocksInfo {
	info := &evmclient.BlocksInfo{
		Blocks:   promise.NewPromise[[]*types.Block](),
		Receipts: promise.NewPromise[[]evmclient.BlockReceipts](),
		Traces:   promise.NewPromise[[]evmclient.BlockTraces](),
		Headers:  promise.NewPromise[[]*types.Header](),
	}

	header := &types.Header{Number: new(big.Int).Set(startBlock)}
	tx := types.NewTx(&types.LegacyTx{Nonce: 1})
	block := types.NewBlockWithHeader(header).WithBody(types.Body{
		Transactions: []*types.Transaction{tx},
	})
	receipts := []evmclient.BlockReceipts{{tx.Hash(): {}}}
	traces := []evmclient.BlockTraces{{tx.Hash(): {}}}
	headers := []*types.Header{header}

	if requiredDataTypes.Has(evmclient.RequiredDataTypesTx) {
		info.Blocks.Resolve([]*types.Block{block})
	}
	if requiredDataTypes.Has(evmclient.RequiredDataTypesReceipt) {
		info.Receipts.Resolve(receipts)
	}
	if requiredDataTypes.Has(evmclient.RequiredDataTypesTraces) {
		info.Traces.Resolve(traces)
	}
	if requiredDataTypes.Has(evmclient.RequiredDataTypesBlockHeader) {
		info.Headers.Resolve(headers)
	}

	return info
}

type trackingProcessor struct {
	requiredDataTypes evmclient.RequiredDataTypes
	errorsCh          chan error
	processErr        error
	processCalls      atomic.Int32
	flushCalls        atomic.Int32
}

func newTrackingProcessor(requiredDataTypes evmclient.RequiredDataTypes, processErr error) *trackingProcessor {
	return &trackingProcessor{
		requiredDataTypes: requiredDataTypes,
		errorsCh:          make(chan error, 1),
		processErr:        processErr,
	}
}

func (p *trackingProcessor) QueueRequest(request parsers.ParseTxRequest) {}

func (p *trackingProcessor) ProcessRequest(request parsers.ParseTxRequest) error {
	p.processCalls.Add(1)
	return p.processErr
}

func (p *trackingProcessor) Stop() {}

func (p *trackingProcessor) Errors() chan error {
	return p.errorsCh
}

func (p *trackingProcessor) RequiredDataTypes() evmclient.RequiredDataTypes {
	return p.requiredDataTypes
}

func (p *trackingProcessor) FlushSync() {
	p.flushCalls.Add(1)
}

func TestProcessorAggregatorProcessBlockNumberSyncReturnsProcessorError(t *testing.T) {
	required := evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesBlockHeader | evmclient.RequiredDataTypesTraces
	processor := newTrackingProcessor(required, errors.New("boom"))
	aggregator := NewProcessorAggregator([]IndexerProcessor{processor}, &staticEVMClient{}, zap.NewNop(), WithSyncProcessing(true))
	t.Cleanup(aggregator.Stop)

	err := aggregator.ProcessBlockNumberSync(context.Background(), big.NewInt(1))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected processor error, got %v", err)
	}
	if processor.flushCalls.Load() != 0 {
		t.Fatalf("expected flush not to run on error, got %d", processor.flushCalls.Load())
	}
}

func TestProcessorAggregatorProcessBlockRangeWaitsForAllProcessors(t *testing.T) {
	required := evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesBlockHeader | evmclient.RequiredDataTypesTraces
	processorA := newTrackingProcessor(required, nil)
	processorB := newTrackingProcessor(required, nil)
	aggregator := NewProcessorAggregator([]IndexerProcessor{processorA, processorB}, &staticEVMClient{}, zap.NewNop(), WithSyncProcessing(true))
	t.Cleanup(aggregator.Stop)

	doneCh := aggregator.ProcessBlockRange(context.Background(), big.NewInt(1), big.NewInt(1))

	select {
	case <-doneCh:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for block range processing")
	}

	if processorA.processCalls.Load() != 1 {
		t.Fatalf("expected processor A to process 1 request, got %d", processorA.processCalls.Load())
	}
	if processorB.processCalls.Load() != 1 {
		t.Fatalf("expected processor B to process 1 request, got %d", processorB.processCalls.Load())
	}
}

func TestProcessorAggregatorSkipsMissingReceiptOrTraceByTxHash(t *testing.T) {
	required := evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesTraces
	processor := newTrackingProcessor(required, nil)
	aggregator := NewProcessorAggregator([]IndexerProcessor{processor}, &staticMissingTxDataEVMClient{}, zap.NewNop(), WithSyncProcessing(true))
	t.Cleanup(aggregator.Stop)

	resultCh := aggregator.ProcessBlockRangeWithResult(context.Background(), big.NewInt(1), big.NewInt(1))

	select {
	case err := <-resultCh:
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for block range result")
	}

	if processor.processCalls.Load() != 1 {
		t.Fatalf("expected only fully matched tx to be processed, got %d", processor.processCalls.Load())
	}
}

func TestProcessorAggregatorProcessBlockRangeWithReceiptOnlyProcessor(t *testing.T) {
	required := evmclient.RequiredDataTypesReceipt
	processor := newTrackingProcessor(required, nil)
	aggregator := NewProcessorAggregator([]IndexerProcessor{processor}, &staticEVMClient{}, zap.NewNop(), WithSyncProcessing(true))
	t.Cleanup(aggregator.Stop)

	resultCh := aggregator.ProcessBlockRangeWithResult(context.Background(), big.NewInt(1), big.NewInt(1))

	select {
	case err := <-resultCh:
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for receipt-only block range result")
	}

	if processor.processCalls.Load() != 1 {
		t.Fatalf("expected receipt-only processor to process 1 request, got %d", processor.processCalls.Load())
	}
}

func TestProcessorAggregatorOrdersReceiptOnlyRequestsByTransactionIndex(t *testing.T) {
	aggregator := NewProcessorAggregator(nil, &staticEVMClient{}, zap.NewNop())
	t.Cleanup(aggregator.Stop)

	tx0 := types.NewTx(&types.LegacyTx{Nonce: 0})
	tx1 := types.NewTx(&types.LegacyTx{Nonce: 1})
	tx2 := types.NewTx(&types.LegacyTx{Nonce: 2})
	receipts := evmclient.BlockReceipts{
		tx2.Hash(): &types.Receipt{TxHash: tx2.Hash(), TransactionIndex: 2},
		tx0.Hash(): &types.Receipt{TxHash: tx0.Hash(), TransactionIndex: 0},
		tx1.Hash(): &types.Receipt{TxHash: tx1.Hash(), TransactionIndex: 1},
	}

	expectedTxs, err := aggregator.checkLengths(
		nil,
		[]evmclient.BlockReceipts{receipts},
		nil,
		nil,
		1,
		evmclient.RequiredDataTypesReceipt,
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	got := expectedTxs[0]
	want := []string{tx0.Hash().Hex(), tx1.Hash().Hex(), tx2.Hash().Hex()}
	for i, txHash := range got {
		if txHash.Hex() != want[i] {
			t.Fatalf("expected tx %d to be %s, got %s", i, want[i], txHash.Hex())
		}
	}
}

func TestProcessorAggregatorProcessBlockRangeWithResultReturnsError(t *testing.T) {
	required := evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesBlockHeader | evmclient.RequiredDataTypesTraces
	processor := newTrackingProcessor(required, errors.New("boom"))
	aggregator := NewProcessorAggregator([]IndexerProcessor{processor}, &staticEVMClient{}, zap.NewNop(), WithSyncProcessing(true))
	t.Cleanup(aggregator.Stop)

	resultCh := aggregator.ProcessBlockRangeWithResult(context.Background(), big.NewInt(1), big.NewInt(1))

	select {
	case err := <-resultCh:
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "boom") {
			t.Fatalf("expected processor error, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for block range result")
	}
}

type staticMissingTxDataEVMClient struct{}

func (c *staticMissingTxDataEVMClient) BatchGetBlocksInfo(ctx context.Context, startBlock *big.Int, endBlock *big.Int, requiredDataTypes evmclient.RequiredDataTypes) *evmclient.BlocksInfo {
	info := &evmclient.BlocksInfo{
		Blocks:   promise.NewPromise[[]*types.Block](),
		Receipts: promise.NewPromise[[]evmclient.BlockReceipts](),
		Traces:   promise.NewPromise[[]evmclient.BlockTraces](),
		Headers:  promise.NewPromise[[]*types.Header](),
	}

	header := &types.Header{Number: new(big.Int).Set(startBlock)}
	tx1 := types.NewTx(&types.LegacyTx{Nonce: 1})
	tx2 := types.NewTx(&types.LegacyTx{Nonce: 2})
	block := types.NewBlockWithHeader(header).WithBody(types.Body{
		Transactions: []*types.Transaction{tx1, tx2},
	})

	info.Blocks.Resolve([]*types.Block{block})
	info.Receipts.Resolve([]evmclient.BlockReceipts{{
		tx1.Hash(): &types.Receipt{TxHash: tx1.Hash()},
	}})
	info.Traces.Resolve([]evmclient.BlockTraces{{
		tx1.Hash(): {},
		tx2.Hash(): {},
	}})
	info.Headers.Resolve(nil)

	return info
}
