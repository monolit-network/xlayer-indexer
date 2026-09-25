package processors

import (
	"bytes"
	"context"
	"fmt"
	"iter"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

type EVMClient interface {
	BatchGetBlocksInfo(ctx context.Context, startBlock *big.Int, endBlock *big.Int, requiredDataTypes evmclient.RequiredDataTypes) *evmclient.BlocksInfo
}

type ProcessorAggregator struct {
	processors []IndexerProcessor
	workersWg  *sync.WaitGroup
	errorsCh   chan error
	cancelFunc context.CancelFunc
	evmClient  EVMClient

	blockInfosCh []chan processBlockRequest

	combinedRequiredDataTypes evmclient.RequiredDataTypes

	stopCh chan struct{}

	syncProcessing bool

	logger *zap.Logger
}

type AggregatorOption func(*ProcessorAggregator)

func WithSyncProcessing(sync bool) AggregatorOption {
	return func(a *ProcessorAggregator) {
		a.syncProcessing = sync
	}
}

type processBlockRequest struct {
	Ctx            context.Context
	Wg             *sync.WaitGroup
	BlockHeaders   []*types.Header
	BlocksInfo     *evmclient.BlocksInfo
	StartBlock     *big.Int
	EndBlock       *big.Int
	ExpectedBlocks int64
	Result         *requestResult
}

type requestResult struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

func newRequestResult(workers int) *requestResult {
	result := &requestResult{}
	if workers > 0 {
		result.wg.Add(workers)
	}
	return result
}

func (r *requestResult) complete(err error) {
	if err != nil {
		r.errOnce.Do(func() {
			r.err = err
		})
	}
	r.wg.Done()
}

func (r *requestResult) wait() error {
	r.wg.Wait()
	return r.err
}

func NewProcessorAggregator(processors []IndexerProcessor, evmClient EVMClient, logger *zap.Logger, opts ...AggregatorOption) *ProcessorAggregator {
	a := &ProcessorAggregator{
		processors: processors,
		workersWg:  &sync.WaitGroup{},
		errorsCh:   make(chan error, 1024),
		evmClient:  evmClient,
		stopCh:     make(chan struct{}),
		logger:     logger,
	}

	for _, opt := range opts {
		opt(a)
	}

	for i := 0; i < len(processors); i++ {
		ch := make(chan processBlockRequest, 1024)
		a.blockInfosCh = append(a.blockInfosCh, ch)
		a.workersWg.Add(1)
		go a.worker(ch, processors[i])
	}

	for _, processor := range processors {
		a.combinedRequiredDataTypes = a.combinedRequiredDataTypes.Add(processor.RequiredDataTypes())
	}

	return a
}

func (a *ProcessorAggregator) Stop() {
	close(a.stopCh)
}

func (a *ProcessorAggregator) worker(ch chan processBlockRequest, processor IndexerProcessor) {
	defer a.workersWg.Done()

	for {
		select {
		case <-a.stopCh:
			return
		case request := <-ch:
			complete := func(err error) {
				if request.Wg != nil {
					request.Wg.Done()
				}
				if request.Result != nil {
					request.Result.complete(err)
				}
			}

			var blocks []*types.Block
			var receipts []evmclient.BlockReceipts
			var traces []evmclient.BlockTraces
			var headers []*types.Header
			var expectedTxs [][]common.Hash
			var err error

			blocks, receipts, traces, headers, err = request.BlocksInfo.WaitContext(request.Ctx, processor.RequiredDataTypes())
			if err != nil {
				err = fmt.Errorf("error getting blocks info: %w", err)
				a.errorsCh <- err
				complete(err)
				continue
			}

			if request.Ctx.Err() != nil {
				complete(request.Ctx.Err())
				continue
			}

			expectedTxs, err = a.checkLengths(blocks, receipts, traces, headers, int(request.ExpectedBlocks), processor.RequiredDataTypes())
			if err != nil {
				err = fmt.Errorf("got unexpected number of blocks/receipts/traces/headers: %w", err)
				a.errorsCh <- err
				complete(err)
				continue
			}

			requests, err := a.constructRequests(blocks, receipts, traces, headers, expectedTxs, int(request.ExpectedBlocks), processor.RequiredDataTypes())
			if err != nil {
				err = fmt.Errorf("error constructing requests: %w", err)
				a.errorsCh <- err
				complete(err)
				continue
			}

			start := time.Now()
			a.logger.Debug("queueing requests", zap.Int64("start_block", request.StartBlock.Int64()), zap.Int64("end_block", request.EndBlock.Int64()))

			var processErr error
			for parseRequest := range requests {
				if parseRequest.Receipt != nil && parseRequest.Receipt.Status == types.ReceiptStatusFailed {
					continue
				}
				if a.syncProcessing {
					if err := processor.ProcessRequest(*parseRequest); err != nil {
						processErr = fmt.Errorf("error processing request: %w", err)
						a.errorsCh <- processErr
						break
					}
				} else {
					processor.QueueRequest(*parseRequest)
				}
			}

			complete(processErr)
			a.logger.Debug("requests queued", zap.Int64("start_block", request.StartBlock.Int64()), zap.Int64("end_block", request.EndBlock.Int64()), zap.Duration("duration", time.Since(start)))
		}
	}
}

func (a *ProcessorAggregator) checkLengths(blocks []*types.Block, receipts []evmclient.BlockReceipts, traces []evmclient.BlockTraces, headers []*types.Header, expectedBlocks int, requiredDataTypes evmclient.RequiredDataTypes) ([][]common.Hash, error) {
	if requiredDataTypes.Has(evmclient.RequiredDataTypesTx) {
		if len(blocks) != expectedBlocks {
			return nil, fmt.Errorf("number of blocks must be %d, got %d", expectedBlocks, len(blocks))
		}
	}
	if requiredDataTypes.Has(evmclient.RequiredDataTypesReceipt) {
		if len(receipts) != expectedBlocks {
			return nil, fmt.Errorf("number of receipts must be %d, got %d", expectedBlocks, len(receipts))
		}
	}
	if requiredDataTypes.Has(evmclient.RequiredDataTypesTraces) {
		if len(traces) != expectedBlocks {
			return nil, fmt.Errorf("number of traces must be %d, got %d", expectedBlocks, len(traces))
		}
	}
	if requiredDataTypes.Has(evmclient.RequiredDataTypesBlockHeader) && !requiredDataTypes.Has(evmclient.RequiredDataTypesTx) {
		if len(headers) != expectedBlocks {
			return nil, fmt.Errorf("number of headers must be %d, got %d", expectedBlocks, len(headers))
		}
	}

	expectedTxs := make([][]common.Hash, 0, expectedBlocks)
	for i := 0; i < expectedBlocks; i++ {
		var blockExpectedTxs []common.Hash
		if requiredDataTypes.Has(evmclient.RequiredDataTypesTx) {
			if blocks[i] == nil {
				return nil, fmt.Errorf("block %d is nil", i)
			}
			blockExpectedTxs = make([]common.Hash, 0, len(blocks[i].Transactions()))
			for _, tx := range blocks[i].Transactions() {
				blockExpectedTxs = append(blockExpectedTxs, tx.Hash())
			}
		}

		if requiredDataTypes.Has(evmclient.RequiredDataTypesReceipt) {
			if receipts[i] == nil {
				if requiredDataTypes.Has(evmclient.RequiredDataTypesTx) && len(blockExpectedTxs) == 0 {
					receipts[i] = evmclient.BlockReceipts{}
				} else {
					return nil, fmt.Errorf("receipt %d is nil", i)
				}
			}
			if len(blockExpectedTxs) == 0 {
				blockExpectedTxs = make([]common.Hash, 0, len(receipts[i]))
				for txHash := range receipts[i] {
					blockExpectedTxs = append(blockExpectedTxs, txHash)
				}
				sortTxHashesByReceiptIndex(blockExpectedTxs, receipts[i])
			}
		}

		if requiredDataTypes.Has(evmclient.RequiredDataTypesTraces) {
			if traces[i] == nil {
				if requiredDataTypes.Has(evmclient.RequiredDataTypesTx) && len(blockExpectedTxs) == 0 {
					traces[i] = evmclient.BlockTraces{}
				} else {
					return nil, fmt.Errorf("traces %d is nil", i)
				}
			}
			if len(blockExpectedTxs) == 0 {
				blockExpectedTxs = make([]common.Hash, 0, len(traces[i]))
				for txHash := range traces[i] {
					blockExpectedTxs = append(blockExpectedTxs, txHash)
				}
			}
		}

		expectedTxs = append(expectedTxs, blockExpectedTxs)
	}
	return expectedTxs, nil
}

func sortTxHashesByReceiptIndex(txHashes []common.Hash, receipts evmclient.BlockReceipts) {
	sort.SliceStable(txHashes, func(i, j int) bool {
		left := receipts[txHashes[i]]
		right := receipts[txHashes[j]]
		if left == nil || right == nil || left.TransactionIndex == right.TransactionIndex {
			return bytes.Compare(txHashes[i].Bytes(), txHashes[j].Bytes()) < 0
		}
		return left.TransactionIndex < right.TransactionIndex
	})
}

func (a *ProcessorAggregator) constructRequests(
	blocks []*types.Block,
	receipts []evmclient.BlockReceipts,
	traces []evmclient.BlockTraces,
	headers []*types.Header,
	expectedTxs [][]common.Hash,
	expectedBlocks int,
	requiredDataTypes evmclient.RequiredDataTypes,
) (iter.Seq[*parsers.ParseTxRequest], error) {
	txsRequired := requiredDataTypes.Has(evmclient.RequiredDataTypesTx)
	receiptsRequired := requiredDataTypes.Has(evmclient.RequiredDataTypesReceipt)
	tracesRequired := requiredDataTypes.Has(evmclient.RequiredDataTypesTraces)
	headersRequired := requiredDataTypes.Has(evmclient.RequiredDataTypesBlockHeader)
	if !txsRequired && !receiptsRequired && !tracesRequired {
		return nil, fmt.Errorf("required data types must include tx, receipt, or traces")
	}

	return func(yield func(*parsers.ParseTxRequest) bool) {
		for i := 0; i < expectedBlocks; i++ {
			txsByHash := make(map[common.Hash]*types.Transaction)
			blockHeader := (*types.Header)(nil)
			if i < len(blocks) && blocks[i] != nil {
				blockHeader = blocks[i].Header()
				for _, tx := range blocks[i].Transactions() {
					txsByHash[tx.Hash()] = tx
				}
			} else if headersRequired && !txsRequired && i < len(headers) {
				blockHeader = headers[i]
			}

			var blockExpectedTxs []common.Hash
			if i < len(expectedTxs) {
				blockExpectedTxs = expectedTxs[i]
			}

			for _, txHash := range blockExpectedTxs {
				req := &parsers.ParseTxRequest{}
				if txsRequired {
					tx, ok := txsByHash[txHash]
					if !ok {
						a.logger.Warn("expected tx not found in block transactions",
							zap.String("block", a.blockNumberForLog(blocks, headers, i)),
							zap.String("tx_hash", txHash.Hex()),
						)
						continue
					}
					req.Tx = tx
					req.BlockHeader = blockHeader
				}
				if receiptsRequired {
					receipt, ok := receipts[i][txHash]
					if !ok {
						a.logger.Warn("receipt not found for tx",
							zap.String("block", a.blockNumberForLog(blocks, headers, i)),
							zap.String("tx_hash", txHash.Hex()),
						)
						continue
					}
					req.Receipt = receipt
				}
				if tracesRequired {
					trace, ok := traces[i][txHash]
					if !ok {
						a.logger.Warn("trace not found for tx",
							zap.String("block", a.blockNumberForLog(blocks, headers, i)),
							zap.String("tx_hash", txHash.Hex()),
						)
						continue
					}
					req.Traces = trace
				}
				if headersRequired && !txsRequired {
					req.BlockHeader = blockHeader
				}
				if !yield(req) {
					return
				}
			}
		}
	}, nil
}

func (a *ProcessorAggregator) blockNumberForLog(blocks []*types.Block, headers []*types.Header, idx int) string {
	if idx < len(blocks) && blocks[idx] != nil && blocks[idx].Number() != nil {
		return blocks[idx].Number().String()
	}
	if idx < len(headers) && headers[idx] != nil && headers[idx].Number != nil {
		return headers[idx].Number.String()
	}
	return "unknown"
}

func (a *ProcessorAggregator) ProcessBlockNumber(ctx context.Context, blockNumber *big.Int) error {
	blocksInfo := a.evmClient.BatchGetBlocksInfo(ctx, blockNumber, blockNumber, a.combinedRequiredDataTypes)
	request := processBlockRequest{
		Ctx:            ctx,
		BlocksInfo:     blocksInfo,
		ExpectedBlocks: 1,
		StartBlock:     blockNumber,
		EndBlock:       blockNumber,
	}
	a.dispatchRequest(request)
	return nil
}

func (a *ProcessorAggregator) ProcessBlockNumberSync(ctx context.Context, blockNumber *big.Int) error {
	_, err := a.ProcessBlockNumberSyncWithHash(ctx, blockNumber)
	return err
}

func (a *ProcessorAggregator) ProcessBlockNumberSyncWithHash(ctx context.Context, blockNumber *big.Int) (common.Hash, error) {
	if !a.syncProcessing {
		return common.Hash{}, fmt.Errorf("sync processing is not enabled")
	}

	if len(a.processors) == 0 {
		return common.Hash{}, nil
	}

	requiredDataTypes := a.combinedRequiredDataTypes.Add(evmclient.RequiredDataTypesBlockHeader)
	blocksInfo := a.evmClient.BatchGetBlocksInfo(ctx, blockNumber, blockNumber, requiredDataTypes)
	result := newRequestResult(len(a.processors))
	request := processBlockRequest{
		Ctx:            ctx,
		BlocksInfo:     blocksInfo,
		ExpectedBlocks: 1,
		StartBlock:     blockNumber,
		EndBlock:       blockNumber,
		Result:         result,
	}
	a.dispatchRequest(request)

	if err := result.wait(); err != nil {
		return common.Hash{}, err
	}

	blocks, _, _, headers, err := blocksInfo.WaitContext(ctx, requiredDataTypes)
	if err != nil {
		return common.Hash{}, err
	}

	var blockHash common.Hash
	switch {
	case len(blocks) == 1 && blocks[0] != nil:
		blockHash = blocks[0].Hash()
	case len(headers) == 1 && headers[0] != nil:
		blockHash = headers[0].Hash()
	default:
		return common.Hash{}, fmt.Errorf("failed to resolve block hash for block %s", blockNumber.String())
	}

	a.FlushSync()
	return blockHash, nil
}

func (a *ProcessorAggregator) dispatchRequest(request processBlockRequest) {
	for _, ch := range a.blockInfosCh {
		ch <- request
	}
}

func (a *ProcessorAggregator) ProcessBlockHeader(ctx context.Context, blockHeader *types.Header) error {
	blocksInfo := a.evmClient.BatchGetBlocksInfo(ctx, blockHeader.Number, blockHeader.Number, a.combinedRequiredDataTypes)
	request := processBlockRequest{
		Ctx:            ctx,
		BlocksInfo:     blocksInfo,
		BlockHeaders:   []*types.Header{blockHeader},
		ExpectedBlocks: 1,
		StartBlock:     blockHeader.Number,
		EndBlock:       blockHeader.Number,
	}
	a.dispatchRequest(request)
	return nil
}

func (a *ProcessorAggregator) ProcessBlockRangeWithResult(ctx context.Context, startBlock *big.Int, endBlock *big.Int) <-chan error {
	blocksInfo := a.evmClient.BatchGetBlocksInfo(ctx, startBlock, endBlock, a.combinedRequiredDataTypes)
	numBlocks := big.NewInt(0).Sub(endBlock, startBlock).Int64() + 1
	result := newRequestResult(len(a.processors))
	request := processBlockRequest{
		Ctx:            ctx,
		BlocksInfo:     blocksInfo,
		ExpectedBlocks: numBlocks,
		StartBlock:     startBlock,
		EndBlock:       endBlock,
		Result:         result,
	}
	a.dispatchRequest(request)

	resultCh := make(chan error, 1)
	go func() {
		defer a.logger.Debug("block range queued to processors", zap.Int64("start_block", startBlock.Int64()), zap.Int64("end_block", endBlock.Int64()))
		defer close(resultCh)
		resultCh <- result.wait()
	}()
	return resultCh
}

func (a *ProcessorAggregator) ProcessBlockRange(ctx context.Context, startBlock *big.Int, endBlock *big.Int) chan struct{} {
	resultCh := a.ProcessBlockRangeWithResult(ctx, startBlock, endBlock)
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		<-resultCh
	}()
	return doneCh
}

func (a *ProcessorAggregator) Errors() chan error {
	return a.errorsCh
}

func (a *ProcessorAggregator) FlushSync() {
	var wg sync.WaitGroup
	for _, processor := range a.processors {
		wg.Go(func() {
			processor.FlushSync()
		})
	}
	wg.Wait()
}
