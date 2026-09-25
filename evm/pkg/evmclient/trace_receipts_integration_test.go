package evmclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func TestRecreateReceiptsFromTraceMatchesRealNode(t *testing.T) {
	loadTestEnv()

	chain := models.Chain(os.Getenv("EVM_CHAIN"))
	if chain == "" {
		chain = models.ChainETH
	}
	if _, ok := models.ChainToID[chain]; !ok {
		t.Fatalf("unsupported EVM_CHAIN %q", chain)
	}

	rpcURL := testRPCURL(chain)
	if rpcURL == "" {
		t.Fatalf("EVM_RPC_URL_HIST_%s, EVM_RPC_URL_%s, EVM_RPC_URL_HIST, or EVM_RPC_URL is required for real-node trace receipt test", strings.ToUpper(string(chain)), strings.ToUpper(string(chain)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	rpcClient, err := dialTraceReceiptTestRPC(ctx, rpcURL)
	if err != nil {
		t.Fatalf("dialTraceReceiptTestRPC() error = %v", err)
	}
	defer rpcClient.Close()

	ethClient := ethclient.NewClient(rpcClient)
	client := NewClient(ethClient, nil, rpcClient, zap.NewNop(), chain)

	blockNumber, err := traceReceiptTestBlock(t, ctx, client)
	if err != nil {
		t.Fatalf("traceReceiptTestBlock() error = %v", err)
	}
	t.Logf("testing trace receipt reconstruction on block %s", blockNumber.String())

	required := RequiredDataTypesTx | RequiredDataTypesReceipt | RequiredDataTypesTraces | RequiredDataTypesBlockHeader
	start := time.Now()
	blocksInfo := client.BatchGetBlocksInfo(ctx, blockNumber, blockNumber, required)
	blocks, recreated, traces, _, err := blocksInfo.WaitContext(ctx, required)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("BatchGetBlocksInfo(%s) error = %v", blockNumber, err)
	}
	if len(blocks) != 1 || blocks[0] == nil {
		t.Fatalf("expected one block for %s", blockNumber)
	}
	if len(traces) != 1 {
		t.Fatalf("expected one trace block for %s, got %d", blockNumber, len(traces))
	}
	if len(recreated) != 1 {
		t.Fatalf("expected one recreated receipt block for %s, got %d", blockNumber, len(recreated))
	}
	t.Logf("BatchGetBlocksInfo() took %v | traces: %d | receipts: %d", elapsed, len(traces[0]), len(recreated[0]))

	start = time.Now()
	_, realReceipts, err := fetchBlockAndReceipts(t, ctx, client, blockNumber)
	t.Logf("fetchBlockAndReceipts() took %v", time.Since(start))
	if err != nil {
		t.Fatalf("fetchBlockAndReceipts(%s) error = %v", blockNumber, err)
	}
	if len(realReceipts) != 1 {
		t.Fatalf("expected one receipt block for %s, got %d", blockNumber, len(realReceipts))
	}

	saveReceiptLogsJSON(t, chain, blockNumber, realReceipts[0], recreated[0])

	assertReceiptBlocksMatch(t, chain, realReceipts[0], recreated[0])
	assertKnownTraceReceiptRegressions(t, chain, recreated[0])
}

func loadTestEnv() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../../../.env")
}

func dialTraceReceiptTestRPC(ctx context.Context, rpcURL string) (*rpc.Client, error) {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}
	httpClient := &http.Client{Transport: transport, Timeout: 300 * time.Second}
	return rpc.DialOptions(ctx, rpcURL, rpc.WithHTTPClient(httpClient))
}

func testRPCURL(chain models.Chain) string {
	suffix := strings.ToUpper(string(chain))
	for _, key := range []string{
		"EVM_RPC_URL_HIST_" + suffix,
		"EVM_RPC_URL_" + suffix,
	} {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	if os.Getenv("EVM_TRACE_RECEIPTS_REQUIRE_RPC") == "1" {
		return ""
	}
	for _, key := range []string{"EVM_RPC_URL_HIST", "EVM_RPC_URL"} {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return ""
}

func traceReceiptTestBlock(t *testing.T, ctx context.Context, client *Client) (*big.Int, error) {
	if raw := os.Getenv("EVM_TRACE_RECEIPTS_TEST_TX"); raw != "" {
		txHash := common.HexToHash(raw)
		var tx struct {
			BlockNumber string `json:"blockNumber"`
		}
		if err := client.GetRpcClient().CallContext(ctx, &tx, "eth_getTransactionByHash", txHash.Hex()); err != nil {
			return nil, err
		}
		if tx.BlockNumber == "" {
			return nil, fmt.Errorf("transaction %s is not found or is pending", txHash.Hex())
		}
		return hexutil.DecodeBig(tx.BlockNumber)
	}

	if raw := os.Getenv("EVM_TRACE_RECEIPTS_TEST_BLOCK"); raw != "" {
		block, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || block <= 0 {
			return nil, fmt.Errorf("invalid EVM_TRACE_RECEIPTS_TEST_BLOCK %q", raw)
		}
		return big.NewInt(block), nil
	}

	latest, err := client.GetEthClient().BlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	if latest <= 20 {
		return nil, fmt.Errorf("latest block %d is too low", latest)
	}

	for block := int64(latest) - 20; block > int64(latest)-80 && block > 0; block-- {
		blockNumber := big.NewInt(block)
		blocks, receipts, err := fetchBlockAndReceipts(t, ctx, client, blockNumber)
		if err != nil || len(blocks) != 1 || blocks[0] == nil || len(receipts) != 1 {
			continue
		}
		if len(blocks[0].Transactions()) == 0 {
			continue
		}
		for _, receipt := range receipts[0] {
			if len(receipt.Logs) > 0 {
				return blockNumber, nil
			}
		}
	}

	return nil, fmt.Errorf("failed to find recent non-empty block with logs near %d", latest)
}

func fetchBlockAndReceipts(t *testing.T, ctx context.Context, client *Client, blockNumber *big.Int) ([]*types.Block, []BlockReceipts, error) {
	start := time.Now()
	t.Logf("BatchGetBlocksWithTransactions(%s) started", blockNumber.String())
	blocks, err := client.BatchGetBlocksWithTransactions(ctx, blockNumber, blockNumber)
	if err != nil {
		return nil, nil, err
	}
	t.Logf("BatchGetBlocksWithTransactions(%s) took %v", blockNumber.String(), time.Since(start))

	start = time.Now()
	t.Logf("BatchGetBlocksReceipts() started")
	receipts, err := client.BatchGetBlocksReceipts(ctx, blockNumber, blockNumber)
	if err != nil {
		return nil, nil, err
	}
	t.Logf("BatchGetBlocksReceipts() took %v", time.Since(start))
	return blocks, receipts, nil
}

type receiptLogJSON struct {
	TxHash             string   `json:"tx_hash"`
	TxIndex            uint     `json:"tx_index"`
	ReceiptLogPosition int      `json:"receipt_log_position"`
	LogIndex           uint     `json:"log_index"`
	Address            string   `json:"address"`
	Topic0             string   `json:"topic0,omitempty"`
	Topics             []string `json:"topics"`
	Data               string   `json:"data"`
}

func saveReceiptLogsJSON(t *testing.T, chain models.Chain, blockNumber *big.Int, real BlockReceipts, recreated BlockReceipts) {
	t.Helper()

	root := os.Getenv("EVM_TRACE_RECEIPTS_SAVE_JSON_DIR")
	if root == "" {
		return
	}

	caseName := os.Getenv("EVM_TRACE_RECEIPTS_SAVE_JSON_CASE")
	if caseName == "" {
		caseName = fmt.Sprintf("%s_block_%s", chain, blockNumber.String())
	}
	caseName = strings.NewReplacer("/", "_", ":", "_", " ", "_").Replace(caseName)
	dir := filepath.Join(root, caseName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir trace receipt json dir %s: %v", dir, err)
	}

	writeReceiptLogsJSON(t, filepath.Join(dir, "receipt_logs.json"), flattenReceiptLogsForJSON(real))
	writeReceiptLogsJSON(t, filepath.Join(dir, "reconstructed_logs.json"), flattenReceiptLogsForJSON(recreated))
	t.Logf("saved receipt log json to %s", dir)
}

func writeReceiptLogsJSON(t *testing.T, path string, logs []receiptLogJSON) {
	t.Helper()

	raw, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		t.Fatalf("marshal receipt log json %s: %v", path, err)
	}
	raw = append(raw, '\n')
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatalf("write receipt log json %s: %v", path, err)
	}
}

func flattenReceiptLogsForJSON(receipts BlockReceipts) []receiptLogJSON {
	out := make([]receiptLogJSON, 0, receiptBlockLogCount(receipts))
	txHashByHex := make(map[string]common.Hash)
	for txHash := range receipts {
		txHashByHex[txHash.Hex()] = txHash
	}
	for _, txHashHex := range sortedReceiptTxHashes(receipts) {
		txHash := txHashByHex[txHashHex]
		receipt := receipts[txHash]
		if receipt == nil {
			continue
		}
		for i, log := range receipt.Logs {
			if log == nil {
				continue
			}
			topics := make([]string, len(log.Topics))
			for topicIndex, topic := range log.Topics {
				topics[topicIndex] = topic.Hex()
			}
			var topic0 string
			if len(topics) > 0 {
				topic0 = topics[0]
			}
			out = append(out, receiptLogJSON{
				TxHash:             log.TxHash.Hex(),
				TxIndex:            log.TxIndex,
				ReceiptLogPosition: i,
				LogIndex:           log.Index,
				Address:            log.Address.Hex(),
				Topic0:             topic0,
				Topics:             topics,
				Data:               hexutil.Encode(log.Data),
			})
		}
	}
	return out
}

func assertReceiptBlocksMatch(t *testing.T, chain models.Chain, real BlockReceipts, recreated BlockReceipts) {
	t.Helper()

	if len(real) != len(recreated) {
		t.Fatalf("receipt count mismatch: real=%d recreated=%d", len(real), len(recreated))
	}
	assertReceiptBlockLogIndexes(t, "real", real)
	assertReceiptBlockLogIndexes(t, "recreated", recreated)

	txHashByHex := make(map[string]common.Hash)
	for txHash := range real {
		txHashByHex[txHash.Hex()] = txHash
	}
	realLogs := receiptBlockLogCount(real)
	recreatedLogs := receiptBlockLogCount(recreated)
	if realLogs != recreatedLogs {
		for _, txHashHex := range sortedReceiptTxHashes(real) {
			txHash := txHashByHex[txHashHex]
			if len(real[txHash].Logs) != len(recreated[txHash].Logs) {
				t.Fatalf("block log count mismatch: real=%d recreated=%d; first tx log count mismatch tx=%s tx_index=%d status(real/recreated)=%d/%d logs(real/recreated)=%d/%d",
					realLogs,
					recreatedLogs,
					txHash.Hex(),
					real[txHash].TransactionIndex,
					real[txHash].Status,
					recreated[txHash].Status,
					len(real[txHash].Logs),
					len(recreated[txHash].Logs),
				)
			}
		}
		t.Fatalf("block log count mismatch: real=%d recreated=%d", realLogs, recreatedLogs)
	}

	for _, txHashHex := range sortedReceiptTxHashes(real) {
		txHash := txHashByHex[txHashHex]
		want := real[txHash]
		got, ok := recreated[txHash]
		if !ok {
			t.Fatalf("missing recreated receipt for tx %s", txHash.Hex())
		}
		if got.TxHash != want.TxHash {
			t.Fatalf("tx hash mismatch for %s: real=%s recreated=%s", txHash.Hex(), want.TxHash.Hex(), got.TxHash.Hex())
		}
		if got.TransactionIndex != want.TransactionIndex {
			t.Fatalf("tx index mismatch for %s: real=%d recreated=%d", txHash.Hex(), want.TransactionIndex, got.TransactionIndex)
		}
		if got.Status != want.Status {
			t.Fatalf("status mismatch for %s: real=%d recreated=%d", txHash.Hex(), want.Status, got.Status)
		}
		if got.GasUsed != want.GasUsed {
			t.Fatalf("gas used mismatch for %s: real=%d recreated=%d", txHash.Hex(), want.GasUsed, got.GasUsed)
		}
		if got.CumulativeGasUsed != want.CumulativeGasUsed {
			t.Fatalf("cumulative gas used mismatch for %s: real=%d recreated=%d", txHash.Hex(), want.CumulativeGasUsed, got.CumulativeGasUsed)
		}
		if got.BlockHash != want.BlockHash {
			t.Fatalf("block hash mismatch for %s: real=%s recreated=%s", txHash.Hex(), want.BlockHash.Hex(), got.BlockHash.Hex())
		}
		if got.BlockNumber == nil || want.BlockNumber == nil || got.BlockNumber.Cmp(want.BlockNumber) != 0 {
			t.Fatalf("block number mismatch for %s: real=%v recreated=%v", txHash.Hex(), want.BlockNumber, got.BlockNumber)
		}
		assertReceiptLogsMatch(t, txHash.Hex(), want.Logs, got.Logs)
	}
}

func assertReceiptBlockLogIndexes(t *testing.T, name string, receipts BlockReceipts) {
	t.Helper()

	txHashByHex := make(map[string]common.Hash)
	expectedIndex := uint(0)
	for _, txHashHex := range sortedReceiptTxHashes(receipts) {
		txHash := txHashByHex[txHashHex]
		if txHash == (common.Hash{}) {
			txHash = common.HexToHash(txHashHex)
		}
		receipt := receipts[txHash]
		if receipt == nil {
			t.Fatalf("%s receipt is nil for tx %s", name, txHashHex)
		}
		for logPosition, log := range receipt.Logs {
			if log == nil {
				t.Fatalf("%s log is nil: tx=%s tx_index=%d pos=%d expected_index=%d", name, txHash.Hex(), receipt.TransactionIndex, logPosition, expectedIndex)
			}
			if log.Index != expectedIndex {
				t.Fatalf("%s log index/order mismatch: tx=%s tx_index=%d pos=%d topic=%s got_index=%d want_index=%d",
					name,
					txHash.Hex(),
					receipt.TransactionIndex,
					logPosition,
					logTopicHex(log),
					log.Index,
					expectedIndex,
				)
			}
			expectedIndex++
		}
	}
}

func sortedReceiptTxHashes(receipts BlockReceipts) []string {
	txHashByHex := make(map[string]common.Hash)
	txHashes := make([]string, 0, len(receipts))
	for txHash := range receipts {
		txHashHex := txHash.Hex()
		txHashes = append(txHashes, txHashHex)
		txHashByHex[txHashHex] = txHash
	}
	sort.Slice(txHashes, func(i, j int) bool {
		left := receipts[txHashByHex[txHashes[i]]]
		right := receipts[txHashByHex[txHashes[j]]]
		if left.TransactionIndex == right.TransactionIndex {
			return txHashes[i] < txHashes[j]
		}
		return left.TransactionIndex < right.TransactionIndex
	})
	return txHashes
}

func logTopicHex(log *types.Log) string {
	if log == nil || len(log.Topics) == 0 {
		return ""
	}
	return log.Topics[0].Hex()
}

func receiptBlockLogCount(receipts BlockReceipts) int {
	count := 0
	for _, receipt := range receipts {
		count += len(receipt.Logs)
	}
	return count
}

func assertReceiptLogsMatch(t *testing.T, txHash string, real []*types.Log, recreated []*types.Log) {
	t.Helper()

	if len(real) != len(recreated) {
		t.Fatalf("log count mismatch for %s: real=%d recreated=%d", txHash, len(real), len(recreated))
	}
	for i := range real {
		want := real[i]
		got := recreated[i]
		if got.Address != want.Address {
			t.Fatalf("log address mismatch for %s log %d: real=%s recreated=%s", txHash, i, want.Address.Hex(), got.Address.Hex())
		}
		if got.Index != want.Index {
			t.Fatalf("log index mismatch for %s log %d: real=%d recreated=%d", txHash, i, want.Index, got.Index)
		}
		if got.TxHash != want.TxHash {
			t.Fatalf("log tx hash mismatch for %s log %d: real=%s recreated=%s", txHash, i, want.TxHash.Hex(), got.TxHash.Hex())
		}
		if got.TxIndex != want.TxIndex {
			t.Fatalf("log tx index mismatch for %s log %d: real=%d recreated=%d", txHash, i, want.TxIndex, got.TxIndex)
		}
		if got.BlockHash != want.BlockHash {
			t.Fatalf("log block hash mismatch for %s log %d: real=%s recreated=%s", txHash, i, want.BlockHash.Hex(), got.BlockHash.Hex())
		}
		if got.BlockNumber != want.BlockNumber {
			t.Fatalf("log block number mismatch for %s log %d: real=%d recreated=%d", txHash, i, want.BlockNumber, got.BlockNumber)
		}
		if want.BlockTimestamp != 0 && got.BlockTimestamp != want.BlockTimestamp {
			t.Fatalf("log block timestamp mismatch for %s log %d: real=%d recreated=%d", txHash, i, want.BlockTimestamp, got.BlockTimestamp)
		}
		if got.Address == models.POLLogTokenAddress && len(got.Topics) == 0 && len(got.Data) == 0 {
			continue
		}
		if !bytes.Equal(got.Data, want.Data) {
			t.Fatalf("log data mismatch for %s log %d index %d: real=%x recreated=%x", txHash, i, want.Index, want.Data, got.Data)
		}
		if len(got.Topics) != len(want.Topics) {
			t.Fatalf("log topic count mismatch for %s log %d index %d: real=%d recreated=%d", txHash, i, want.Index, len(want.Topics), len(got.Topics))
		}
		for topicIndex := range want.Topics {
			if got.Topics[topicIndex] != want.Topics[topicIndex] {
				t.Fatalf("log topic mismatch for %s log %d index %d topic %d: real=%s recreated=%s", txHash, i, want.Index, topicIndex, want.Topics[topicIndex].Hex(), got.Topics[topicIndex].Hex())
			}
		}
	}
}

func assertKnownTraceReceiptRegressions(t *testing.T, chain models.Chain, recreated BlockReceipts) {
	t.Helper()
	if chain != models.ChainPolygon {
		return
	}

	txHash := common.HexToHash("0x1af2512df44baa8f93ac12ee34ec153ac3de9c86b05e5ea45b123d6b63461195")
	receipt, ok := recreated[txHash]
	if !ok {
		return
	}

	requestPricePos := -1
	questionResetPos := -1
	for i, log := range receipt.Logs {
		if log == nil || len(log.Topics) == 0 {
			continue
		}
		switch log.Topics[0] {
		case models.UmaOptimisticOracleV2RequestPriceEventSelectorHash:
			requestPricePos = i
		case models.PolymarketUmaCtfAdapterQuestionResetEventSelectorHash:
			questionResetPos = i
		}
	}
	if requestPricePos < 0 {
		t.Fatalf("polygon regression tx %s: request price log is missing", txHash.Hex())
	}
	if questionResetPos < 0 {
		t.Fatalf("polygon regression tx %s: question reset log is missing", txHash.Hex())
	}
	if requestPricePos+1 != questionResetPos {
		t.Fatalf("polygon regression tx %s: request price log must be immediately before question reset log; request_pos=%d reset_pos=%d request_index=%d reset_index=%d",
			txHash.Hex(),
			requestPricePos,
			questionResetPos,
			receipt.Logs[requestPricePos].Index,
			receipt.Logs[questionResetPos].Index,
		)
	}
}
