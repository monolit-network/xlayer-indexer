package evmclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/util/iterators"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/zap"
)

func callsFirstError(calls []rpc.BatchElem) error {
	for _, call := range calls {
		if call.Error != nil {
			return call.Error
		}
	}
	return nil
}

// mergeFn must return a new, non-nil instance of T if no arguments are provided
func smartBatchCall[T any](
	ctx context.Context,
	logger *zap.Logger,
	onRetry func(),
	onSplit func(),
	minBatchSize int,
	maxBatchSize int,
	currentBatchSize int,
	startBlock *big.Int,
	endBlock *big.Int,
	maxRetries int,
	fn func(ctx context.Context, startBlock *big.Int, endBlock *big.Int) (T, error),
	mergeFn func(...T) T,
) (T, error) {
	if currentBatchSize <= 0 || currentBatchSize > maxBatchSize {
		currentBatchSize = maxBatchSize
	}

	if currentBatchSize < minBatchSize {
		return mergeFn(), fmt.Errorf("current batch size %d is less than min %d", currentBatchSize, minBatchSize)
	}

	var results []T

	for _, subRange := range iterators.BigRangeSubranges(startBlock, endBlock, currentBatchSize) {
		var res T
		var err error

		for attempt := 0; attempt < maxRetries || maxRetries < 0; attempt++ {
			if ctx.Err() != nil {
				return mergeFn(), ctx.Err()
			}

			res, err = fn(ctx, subRange.Start, subRange.End)
			if err == nil {
				break
			}

			if errors.Is(err, context.Canceled) {
				return mergeFn(), err
			}

			if isRetryableError(err) {
				if onRetry != nil {
					onRetry()
				}
				logger.Warn("retryable error in evm batch caller", zap.Error(err))
				time.Sleep(10 * time.Second)
				continue
			}

			break
		}

		if err != nil {
			if isBatchSplitError(err) {
				if onSplit != nil {
					onSplit()
				}
				nextBatchSize := currentBatchSize / 2
				subRes, subErr := smartBatchCall(
					ctx,
					logger,
					onRetry,
					onSplit,
					minBatchSize,
					maxBatchSize,
					nextBatchSize,
					subRange.Start,
					subRange.End,
					maxRetries,
					fn,
					mergeFn,
				)
				if subErr != nil {
					return mergeFn(), subErr
				}
				results = append(results, subRes)
				continue
			}

			logger.Error("non-retryable error in evm batch caller", zap.Error(err))
			return mergeFn(), err
		}

		results = append(results, res)
	}

	return mergeFn(results...), nil
}

func (c *Client) batchGetBlocksReceipts(ctx context.Context, startBlock *big.Int, endBlock *big.Int) (results []BlockReceipts, err error) {
	totalBlocks := big.NewInt(0).Sub(endBlock, startBlock).Int64() + 1
	started := time.Now()
	batch := newBatchTiming("receipts", startBlock.String(), endBlock.String(), totalBlocks)
	defer func() {
		c.finishBatchTiming(batch, started, err)
	}()

	if err := c.measureBatchStage(batch, "semaphore_wait", func() error {
		return c.blockReceiptsSem.Acquire(ctx, 1)
	}); err != nil {
		return nil, err
	}
	defer c.blockReceiptsSem.Release(1)
	if err := c.measureBatchStage(batch, "rate_limit_wait", func() error {
		return c.waitRequests(ctx, 1)
	}); err != nil {
		return nil, err
	}

	resultsRaw := make([]json.RawMessage, totalBlocks)
	calls := make([]rpc.BatchElem, 0, totalBlocks)
	for i, blockNumber := range iterators.BigRange(startBlock, endBlock) {
		resultsRaw[i] = json.RawMessage{}
		calls = append(calls, rpc.BatchElem{
			Method: "eth_getBlockReceipts",
			Args:   []interface{}{hexutil.EncodeBig(blockNumber)},
			Result: &resultsRaw[i],
		})
	}

	if err := c.measureBatchStage(batch, "rate_limit_wait", func() error {
		return c.waitRequests(ctx, 1)
	}); err != nil {
		return nil, err
	}

	if err := c.measureBatchStage(batch, "rpc", func() error {
		c.recordRPCCall()
		return c.rpcClient.BatchCallContext(ctx, calls)
	}); err != nil {
		return nil, err
	}

	firstErr := callsFirstError(calls)
	if firstErr != nil {
		return nil, firstErr
	}

	rawBytes := uint64(0)
	for _, resultRaw := range resultsRaw {
		rawBytes += uint64(len(resultRaw))
	}
	batch.addRawBytes(rawBytes)
	c.timings.recordRawBytes(batch.dataType, rawBytes)

	receiptsRaw := make([][]rpcReceipt, 0, totalBlocks)
	if err := c.measureBatchStage(batch, "unmarshal", func() error {
		for _, resultRaw := range resultsRaw {
			receipts := make([]rpcReceipt, 0)
			if err := json.Unmarshal(resultRaw, &receipts); err != nil {
				return err
			}
			receiptsRaw = append(receiptsRaw, receipts)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	results = make([]BlockReceipts, 0, totalBlocks)
	receiptCount := uint64(0)
	if err := c.measureBatchStage(batch, "postprocess", func() error {
		for _, receipts := range receiptsRaw {
			receiptCount += uint64(len(receipts))
			results = append(results, filterZeroPlaceholderReceipts(receipts))
		}
		return nil
	}); err != nil {
		return nil, err
	}
	batch.addItems("receipts", receiptCount)
	c.timings.recordItems(batch.dataType, "receipts", receiptCount)
	batch.addItems("blocks", uint64(totalBlocks))
	c.timings.recordItems(batch.dataType, "blocks", uint64(totalBlocks))

	return results, nil
}

func (c *Client) batchGetBlocksDebugTraces(ctx context.Context, startBlock *big.Int, endBlock *big.Int, withLogs bool) (results []BlockTraces, err error) {
	totalBlocks := big.NewInt(0).Sub(endBlock, startBlock).Int64() + 1
	started := time.Now()
	dataType := "traces"
	if withLogs {
		dataType = "traces_with_logs"
	}
	batch := newBatchTiming(dataType, startBlock.String(), endBlock.String(), totalBlocks)
	defer func() {
		c.finishBatchTiming(batch, started, err)
	}()

	if err := c.measureBatchStage(batch, "semaphore_wait", func() error {
		return c.blockTracesSem.Acquire(ctx, 1)
	}); err != nil {
		return nil, err
	}
	defer c.blockTracesSem.Release(1)
	if err := c.measureBatchStage(batch, "rate_limit_wait", func() error {
		return c.waitRequests(ctx, 1)
	}); err != nil {
		return nil, err
	}

	resultsRaw := make([]json.RawMessage, 0, totalBlocks)
	calls := make([]rpc.BatchElem, 0, totalBlocks)
	for i, blockNumber := range iterators.BigRange(startBlock, endBlock) {
		resultsRaw = append(resultsRaw, json.RawMessage{})
		traceConfig := map[string]interface{}{"tracer": "callTracer"}
		if withLogs {
			traceConfig["tracerConfig"] = map[string]interface{}{"withLog": true}
		}
		calls = append(calls, rpc.BatchElem{
			Method: "debug_traceBlockByNumber",
			Args:   []interface{}{hexutil.EncodeBig(blockNumber), traceConfig},
			Result: &resultsRaw[i],
		})
	}
	if err := c.measureBatchStage(batch, "rpc", func() error {
		c.recordRPCCall()
		return c.rpcClient.BatchCallContext(ctx, calls)
	}); err != nil {
		return nil, err
	}

	firstErr := callsFirstError(calls)
	if firstErr != nil {
		return nil, firstErr
	}

	type blockTraceResultItem struct {
		TxHash common.Hash             `json:"txHash"`
		Trace  *models.FullTraceResult `json:"result"`
	}

	rawBytes := uint64(0)
	for _, resultRaw := range resultsRaw {
		rawBytes += uint64(len(resultRaw))
	}
	batch.addRawBytes(rawBytes)
	c.timings.recordRawBytes(batch.dataType, rawBytes)

	fullResults := make([][]*blockTraceResultItem, 0, totalBlocks)
	if err := c.measureBatchStage(batch, "unmarshal", func() error {
		for _, resultRaw := range resultsRaw {
			fullResult := make([]*blockTraceResultItem, 0)
			if err := json.Unmarshal(resultRaw, &fullResult); err != nil {
				return err
			}
			fullResults = append(fullResults, fullResult)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	results = make([]BlockTraces, 0, totalBlocks)
	traceTxCount := uint64(0)
	if err := c.measureBatchStage(batch, "postprocess", func() error {
		for _, fullResult := range fullResults {
			traceTxCount += uint64(len(fullResult))
			result := make(BlockTraces, len(fullResult))
			for txIndex, item := range fullResult {
				if item.Trace != nil {
					item.Trace.TxHash = item.TxHash
					item.Trace.TxIndex = uint(txIndex)
				}
				result[item.TxHash] = item.Trace
			}
			results = append(results, filterZeroPlaceholderTraces(result))
		}
		return nil
	}); err != nil {
		return nil, err
	}
	batch.addItems("trace_txs", traceTxCount)
	c.timings.recordItems(batch.dataType, "trace_txs", traceTxCount)
	batch.addItems("blocks", uint64(totalBlocks))
	c.timings.recordItems(batch.dataType, "blocks", uint64(totalBlocks))
	return results, nil
}

func (c *Client) batchGetBlocksWithTransactions(ctx context.Context, startBlock *big.Int, endBlock *big.Int) (results []*types.Block, err error) {
	totalBlocks := big.NewInt(0).Sub(endBlock, startBlock).Int64() + 1
	started := time.Now()
	batch := newBatchTiming("blocks", startBlock.String(), endBlock.String(), totalBlocks)
	defer func() {
		c.finishBatchTiming(batch, started, err)
	}()

	if err := c.measureBatchStage(batch, "semaphore_wait", func() error {
		return c.blocksSem.Acquire(ctx, 1)
	}); err != nil {
		return nil, err
	}
	defer c.blocksSem.Release(1)
	if err := c.measureBatchStage(batch, "rate_limit_wait", func() error {
		return c.waitRequests(ctx, 1)
	}); err != nil {
		return nil, err
	}

	calls := make([]rpc.BatchElem, 0, totalBlocks)
	responses := make([]*json.RawMessage, 0, totalBlocks)
	for i, blockNumber := range iterators.BigRange(startBlock, endBlock) {
		responses = append(responses, new(json.RawMessage))
		calls = append(calls, rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{hexutil.EncodeBig(blockNumber), true},
			Result: responses[i],
		})
	}

	if err := c.measureBatchStage(batch, "rpc", func() error {
		c.recordRPCCall()
		return c.rpcClient.BatchCallContext(ctx, calls)
	}); err != nil {
		return nil, err
	}

	var firstErr error
	rawBytes := uint64(0)
	for _, response := range responses {
		if response != nil {
			rawBytes += uint64(len(*response))
		}
	}
	batch.addRawBytes(rawBytes)
	c.timings.recordRawBytes(batch.dataType, rawBytes)

	results = make([]*types.Block, totalBlocks)
	txCount := uint64(0)
	if err := c.measureBatchStage(batch, "unmarshal", func() error {
		for i, response := range responses {
			block, err := blockFromRpcJson(*response, true)
			if err != nil && firstErr == nil {
				firstErr = err
			}
			if block != nil {
				txCount += uint64(len(block.Transactions()))
			}
			results[i] = block
		}
		return firstErr
	}); err != nil {
		return results, err
	}
	batch.addItems("txs", txCount)
	c.timings.recordItems(batch.dataType, "txs", txCount)
	batch.addItems("blocks", uint64(totalBlocks))
	c.timings.recordItems(batch.dataType, "blocks", uint64(totalBlocks))

	return results, firstErr
}

func (c *Client) batchGetBlocksHeaders(ctx context.Context, startBlock *big.Int, endBlock *big.Int) (results []*types.Header, err error) {
	totalBlocks := big.NewInt(0).Sub(endBlock, startBlock).Int64() + 1
	started := time.Now()
	batch := newBatchTiming("headers", startBlock.String(), endBlock.String(), totalBlocks)
	defer func() {
		c.finishBatchTiming(batch, started, err)
	}()

	if err := c.measureBatchStage(batch, "semaphore_wait", func() error {
		return c.blockHeadersSem.Acquire(ctx, 1)
	}); err != nil {
		return nil, err
	}
	defer c.blockHeadersSem.Release(1)
	if err := c.measureBatchStage(batch, "rate_limit_wait", func() error {
		return c.waitRequests(ctx, 1)
	}); err != nil {
		return nil, err
	}

	calls := make([]rpc.BatchElem, 0, totalBlocks)
	responses := make([]*json.RawMessage, 0, totalBlocks)
	for i, blockNumber := range iterators.BigRange(startBlock, endBlock) {
		responses = append(responses, new(json.RawMessage))
		calls = append(calls, rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{hexutil.EncodeBig(blockNumber), false},
			Result: responses[i],
		})
	}

	if err := c.measureBatchStage(batch, "rpc", func() error {
		c.recordRPCCall()
		return c.rpcClient.BatchCallContext(ctx, calls)
	}); err != nil {
		return nil, err
	}

	var firstErr error
	rawBytes := uint64(0)
	for _, response := range responses {
		if response != nil {
			rawBytes += uint64(len(*response))
		}
	}
	batch.addRawBytes(rawBytes)
	c.timings.recordRawBytes(batch.dataType, rawBytes)

	results = make([]*types.Header, totalBlocks)
	headerCount := uint64(0)
	if err := c.measureBatchStage(batch, "unmarshal", func() error {
		for i, response := range responses {
			block, err := blockFromRpcJson(*response, false)
			if err != nil && firstErr == nil {
				firstErr = err
			}
			if err == nil {
				results[i] = block.Header()
				headerCount++
			}
		}
		return firstErr
	}); err != nil {
		return results, err
	}
	batch.addItems("headers", headerCount)
	c.timings.recordItems(batch.dataType, "headers", headerCount)
	batch.addItems("blocks", uint64(totalBlocks))
	c.timings.recordItems(batch.dataType, "blocks", uint64(totalBlocks))

	return results, firstErr
}

func (c *Client) BatchGetBlocksReceipts(ctx context.Context, startBlock *big.Int, endBlock *big.Int) ([]BlockReceipts, error) {
	return smartBatchCall(
		ctx,
		c.logger,
		c.recordRetry,
		c.recordBatchSplit,
		1,
		c.maxBlockBatchSizeReceipts,
		0,
		startBlock,
		endBlock,
		-1,
		c.batchGetBlocksReceipts,
		func(receipts ...[]BlockReceipts) []BlockReceipts {
			if len(receipts) == 0 {
				return []BlockReceipts{}
			}
			totalReceipts := 0
			for _, receipts := range receipts {
				totalReceipts += len(receipts)
			}
			result := make([]BlockReceipts, 0, totalReceipts)
			for _, receipts := range receipts {
				result = append(result, receipts...)
			}
			return result
		},
	)
}

func (c *Client) BatchGetBlocksDebugTraces(ctx context.Context, startBlock *big.Int, endBlock *big.Int) ([]BlockTraces, error) {
	return c.batchGetBlocksDebugTracesSmart(ctx, startBlock, endBlock, false)
}

func (c *Client) BatchGetBlocksDebugTracesWithLogs(ctx context.Context, startBlock *big.Int, endBlock *big.Int) ([]BlockTraces, error) {
	return c.batchGetBlocksDebugTracesSmart(ctx, startBlock, endBlock, true)
}

func (c *Client) batchGetBlocksDebugTracesSmart(ctx context.Context, startBlock *big.Int, endBlock *big.Int, withLogs bool) ([]BlockTraces, error) {
	return smartBatchCall(
		ctx,
		c.logger,
		c.recordRetry,
		c.recordBatchSplit,
		1,
		c.maxBlockBatchSizeTraces,
		0,
		startBlock,
		endBlock,
		-1,
		func(ctx context.Context, startBlock *big.Int, endBlock *big.Int) ([]BlockTraces, error) {
			return c.batchGetBlocksDebugTraces(ctx, startBlock, endBlock, withLogs)
		},
		func(allBlockTraces ...[]BlockTraces) []BlockTraces {
			if len(allBlockTraces) == 0 {
				return []BlockTraces{}
			}
			totalTraceBlocks := 0
			for _, traces := range allBlockTraces {
				totalTraceBlocks += len(traces)
			}
			result := make([]BlockTraces, 0, totalTraceBlocks)
			for _, blockTraces := range allBlockTraces {
				result = append(result, blockTraces...)
			}
			return result
		},
	)
}

func (c *Client) BatchGetBlocksWithTransactions(ctx context.Context, startBlock *big.Int, endBlock *big.Int) ([]*types.Block, error) {
	return smartBatchCall(
		ctx,
		c.logger,
		c.recordRetry,
		c.recordBatchSplit,
		1,
		c.maxBlockBatchSizeBlocks,
		0,
		startBlock,
		endBlock,
		-1,
		c.batchGetBlocksWithTransactions,
		func(blocks ...[]*types.Block) []*types.Block {
			if len(blocks) == 0 {
				return []*types.Block{}
			}
			totalBlocks := 0
			for _, blocks := range blocks {
				totalBlocks += len(blocks)
			}
			result := make([]*types.Block, 0, totalBlocks)
			for _, blocks := range blocks {
				result = append(result, blocks...)
			}
			return result
		},
	)
}

func (c *Client) BatchGetBlocksHeaders(ctx context.Context, startBlock *big.Int, endBlock *big.Int) ([]*types.Header, error) {
	return smartBatchCall(
		ctx,
		c.logger,
		c.recordRetry,
		c.recordBatchSplit,
		1,
		c.maxBlockBatchSizeHeaders,
		0,
		startBlock,
		endBlock,
		-1,
		c.batchGetBlocksHeaders,
		func(headers ...[]*types.Header) []*types.Header {
			if len(headers) == 0 {
				return []*types.Header{}
			}
			totalHeaders := 0
			for _, headers := range headers {
				totalHeaders += len(headers)
			}
			result := make([]*types.Header, 0, totalHeaders)
			for _, headers := range headers {
				result = append(result, headers...)
			}
			return result
		},
	)
}
