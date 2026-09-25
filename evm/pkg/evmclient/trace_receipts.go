package evmclient

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type traceReceiptTx struct {
	hash  common.Hash
	trace *models.FullTraceResult
	tx    *types.Transaction
}

func recreateReceiptsFromTraces(blocks []*types.Block, headers []*types.Header, traces []BlockTraces, chain models.Chain) ([]BlockReceipts, error) {
	if len(traces) == 0 {
		return []BlockReceipts{}, nil
	}

	results := make([]BlockReceipts, 0, len(traces))
	for i, blockTraces := range traces {
		var block *types.Block
		if i < len(blocks) {
			block = blocks[i]
		}

		var header *types.Header
		if block != nil {
			header = block.Header()
		} else if i < len(headers) {
			header = headers[i]
		}
		if header == nil {
			return nil, fmt.Errorf("block header is required to recreate receipts for trace block %d", i)
		}

		results = append(results, recreateBlockReceiptsFromTraces(block, header, blockTraces, chain))
	}
	return results, nil
}

func recreateBlockReceiptsFromTraces(block *types.Block, header *types.Header, traces BlockTraces, chain models.Chain) BlockReceipts {
	if len(traces) == 0 {
		return BlockReceipts{}
	}

	blockHash := header.Hash()
	blockNumber := header.Number
	blockTimestamp := header.Time

	ordered := make([]traceReceiptTx, 0, len(traces))
	txsByHash := make(map[common.Hash]*types.Transaction)
	if block != nil {
		for _, tx := range block.Transactions() {
			txsByHash[tx.Hash()] = tx
		}
	}
	for txHash, trace := range traces {
		if trace == nil {
			continue
		}
		trace.TxHash = txHash
		ordered = append(ordered, traceReceiptTx{hash: txHash, trace: trace, tx: txsByHash[txHash]})
	}
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].trace.TxIndex == ordered[j].trace.TxIndex {
			return bytes.Compare(ordered[i].hash.Bytes(), ordered[j].hash.Bytes()) < 0
		}
		return ordered[i].trace.TxIndex < ordered[j].trace.TxIndex
	})

	receipts := make(BlockReceipts, len(ordered))
	cumulativeGasUsed := uint64(0)
	nextLogIndex := uint(0)
	for _, item := range ordered {
		gasUsed := traceGasUsed(item.trace)
		cumulativeGasUsed += gasUsed

		status := types.ReceiptStatusSuccessful
		var logs []*types.Log
		nextIndex := nextLogIndex
		if traceFailed(item.trace) {
			status = types.ReceiptStatusFailed
			if chain == models.ChainPolygon {
				logs, nextIndex = collectReceiptLogsFromTrace(item.trace, item.hash, item.trace.TxIndex, blockHash, blockNumber.Uint64(), blockTimestamp, nextLogIndex, models.POLLogTokenAddress)
				if !traceContainsLogAddress(item.trace, models.POLLogTokenAddress) {
					logs = append(logs, &types.Log{
						Address:        models.POLLogTokenAddress,
						BlockNumber:    blockNumber.Uint64(),
						TxHash:         item.hash,
						TxIndex:        item.trace.TxIndex,
						BlockHash:      blockHash,
						BlockTimestamp: blockTimestamp,
						Index:          nextIndex,
					})
					nextIndex++
				}
			}
		}
		if status == types.ReceiptStatusSuccessful {
			logs, nextIndex = collectReceiptLogsFromTrace(item.trace, item.hash, item.trace.TxIndex, blockHash, blockNumber.Uint64(), blockTimestamp, nextLogIndex)
		}
		nextLogIndex = nextIndex

		receipt := &types.Receipt{
			Type:              receiptType(item.tx),
			Status:            status,
			CumulativeGasUsed: cumulativeGasUsed,
			Logs:              logs,
			TxHash:            item.hash,
			GasUsed:           gasUsed,
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			TransactionIndex:  item.trace.TxIndex,
		}
		receipt.Bloom = types.CreateBloom(receipt)
		receipts[item.hash] = receipt
	}

	return receipts
}

func traceContainsLogAddress(trace *models.FullTraceResult, address common.Address) bool {
	if trace == nil {
		return false
	}
	for _, log := range trace.Logs {
		if log.Address == address {
			return true
		}
	}
	for _, call := range trace.Calls {
		if traceContainsLogAddress(call, address) {
			return true
		}
	}
	return false
}

func receiptType(tx *types.Transaction) uint8 {
	if tx == nil {
		return types.LegacyTxType
	}
	return tx.Type()
}

func traceFailed(trace *models.FullTraceResult) bool {
	return trace != nil && (trace.Error != "" || trace.RevertReason != "")
}

func traceGasUsed(trace *models.FullTraceResult) uint64 {
	if trace == nil || trace.GasUsed == "" {
		return 0
	}
	gasUsed, err := hexutil.DecodeUint64(trace.GasUsed)
	if err != nil {
		return 0
	}
	return gasUsed
}

type orderedTraceLog struct {
	log models.TraceLog
}

func collectReceiptLogsFromTrace(trace *models.FullTraceResult, txHash common.Hash, txIndex uint, blockHash common.Hash, blockNumber uint64, blockTimestamp uint64, firstLogIndex uint, onlyAddresses ...common.Address) ([]*types.Log, uint) {
	traceLogs := collectTraceLogsInReceiptOrder(trace)
	if len(onlyAddresses) > 0 {
		allowed := make(map[common.Address]struct{}, len(onlyAddresses))
		for _, address := range onlyAddresses {
			allowed[address] = struct{}{}
		}
		filtered := traceLogs[:0]
		for _, traceLog := range traceLogs {
			if _, ok := allowed[traceLog.log.Address]; ok {
				filtered = append(filtered, traceLog)
			}
		}
		traceLogs = filtered
	}
	if len(traceLogs) == 0 {
		return nil, firstLogIndex
	}
	logs := make([]*types.Log, 0, len(traceLogs))
	for i, traceLog := range traceLogs {
		logIndex := firstLogIndex + uint(i)
		topics := make([]common.Hash, len(traceLog.log.Topics))
		copy(topics, traceLog.log.Topics)
		data := make([]byte, len(traceLog.log.Data))
		copy(data, traceLog.log.Data)

		logs = append(logs, &types.Log{
			Address:        traceLog.log.Address,
			Topics:         topics,
			Data:           data,
			BlockNumber:    blockNumber,
			TxHash:         txHash,
			TxIndex:        txIndex,
			BlockHash:      blockHash,
			BlockTimestamp: blockTimestamp,
			Index:          logIndex,
		})
	}

	return logs, firstLogIndex + uint(len(logs))
}

func collectTraceLogsInReceiptOrder(trace *models.FullTraceResult) []orderedTraceLog {
	traceLogs := make([]orderedTraceLog, 0)
	collectTraceLogsDepthFirst(trace, &traceLogs)
	if traceLogsHaveIndexes(traceLogs) {
		sort.SliceStable(traceLogs, func(i, j int) bool {
			return *traceLogs[i].log.Index < *traceLogs[j].log.Index
		})
		return traceLogs
	}

	traceLogs = traceLogs[:0]
	collectTraceLogsByPosition(trace, &traceLogs)
	return traceLogs
}

func traceLogsHaveIndexes(logs []orderedTraceLog) bool {
	for _, log := range logs {
		if log.log.Index == nil {
			return false
		}
	}
	return len(logs) > 0
}

func collectTraceLogsDepthFirst(trace *models.FullTraceResult, out *[]orderedTraceLog) {
	if trace == nil {
		return
	}
	for _, log := range trace.Logs {
		*out = append(*out, orderedTraceLog{log: log})
	}
	for _, call := range trace.Calls {
		collectTraceLogsDepthFirst(call, out)
	}
}

func collectTraceLogsByPosition(trace *models.FullTraceResult, out *[]orderedTraceLog) {
	if trace == nil {
		return
	}

	nextLog := 0
	for callPosition, call := range trace.Calls {
		for nextLog < len(trace.Logs) {
			position := trace.Logs[nextLog].Position
			if position != nil && uint64(*position) > uint64(callPosition) {
				break
			}
			*out = append(*out, orderedTraceLog{log: trace.Logs[nextLog]})
			nextLog++
		}
		collectTraceLogsByPosition(call, out)
	}
	for ; nextLog < len(trace.Logs); nextLog++ {
		*out = append(*out, orderedTraceLog{log: trace.Logs[nextLog]})
	}
}
