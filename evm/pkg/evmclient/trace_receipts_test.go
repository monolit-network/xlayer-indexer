package evmclient

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestRecreatePolygonReceiptsAddsSyntheticFeeLogs(t *testing.T) {
	header := &types.Header{
		Number: big.NewInt(1),
		Time:   123,
	}
	tx0 := common.HexToHash("0x01")
	tx1 := common.HexToHash("0x02")
	tx2 := common.HexToHash("0x03")
	otherLogAddress := common.HexToAddress("0x000000000000000000000000000000000000beef")
	revertedLogAddress := common.HexToAddress("0x000000000000000000000000000000000000dead")

	receipts := recreateBlockReceiptsFromTraces(nil, header, BlockTraces{
		tx0: {
			SingleTraceResult: models.SingleTraceResult{
				GasUsed: "0x1",
				Error:   "execution reverted",
				Logs: []models.TraceLog{
					{Address: revertedLogAddress},
				},
			},
			TxIndex: 0,
		},
		tx1: {
			SingleTraceResult: models.SingleTraceResult{
				GasUsed: "0x1",
				Error:   "execution reverted",
				Logs: []models.TraceLog{
					{Address: revertedLogAddress},
				},
			},
			TxIndex: 1,
		},
		tx2: {
			SingleTraceResult: models.SingleTraceResult{
				GasUsed: "0x1",
				Logs: []models.TraceLog{
					{Address: otherLogAddress},
				},
			},
			TxIndex: 2,
		},
	}, models.ChainPolygon)

	if len(receipts[tx0].Logs) != 1 {
		t.Fatalf("tx0 log count mismatch: got %d, want 1", len(receipts[tx0].Logs))
	}
	if len(receipts[tx1].Logs) != 1 {
		t.Fatalf("tx1 log count mismatch: got %d, want 1", len(receipts[tx1].Logs))
	}
	assertSyntheticPolygonFeeLog(t, receipts[tx0].Logs[0], tx0, 0, 0)
	assertSyntheticPolygonFeeLog(t, receipts[tx1].Logs[0], tx1, 1, 1)

	tx2Logs := receipts[tx2].Logs
	if len(tx2Logs) != 1 {
		t.Fatalf("tx2 log count mismatch: got %d, want 1", len(tx2Logs))
	}
	if tx2Logs[0].Address != otherLogAddress {
		t.Fatalf("tx2 first log address mismatch: got %s, want %s", tx2Logs[0].Address.Hex(), otherLogAddress.Hex())
	}
	if tx2Logs[0].Index != 2 {
		t.Fatalf("tx2 other log index mismatch: got %d, want 2", tx2Logs[0].Index)
	}
}

func TestRecreateReceiptsOrdersTraceLogsByPosition(t *testing.T) {
	var trace models.FullTraceResult
	if err := json.Unmarshal([]byte(`{
		"gasUsed": "0x1",
		"logs": [
			{
				"position": "0x2",
				"address": "0x0000000000000000000000000000000000000002",
				"topics": ["0x0000000000000000000000000000000000000000000000000000000000000002"],
				"data": "0x02"
			}
		],
		"calls": [
			{
				"gasUsed": "0x1",
				"logs": [
					{
						"position": "0x1",
						"address": "0x0000000000000000000000000000000000000001",
						"topics": ["0x0000000000000000000000000000000000000000000000000000000000000001"],
						"data": "0x01"
					}
				]
			}
		]
	}`), &trace); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	trace.TxIndex = 7

	txHash := common.HexToHash("0x1234")
	logs, nextIndex := collectReceiptLogsFromTrace(&trace, txHash, trace.TxIndex, common.Hash{}, 1, 123, 850)
	if nextIndex != 852 {
		t.Fatalf("next log index mismatch: got %d, want 852", nextIndex)
	}
	if len(logs) != 2 {
		t.Fatalf("log count mismatch: got %d, want 2", len(logs))
	}
	if logs[0].Address != common.HexToAddress("0x0000000000000000000000000000000000000001") {
		t.Fatalf("first log address mismatch: got %s", logs[0].Address.Hex())
	}
	if logs[0].Index != 850 {
		t.Fatalf("first log index mismatch: got %d, want 850", logs[0].Index)
	}
	if logs[1].Address != common.HexToAddress("0x0000000000000000000000000000000000000002") {
		t.Fatalf("second log address mismatch: got %s", logs[1].Address.Hex())
	}
	if logs[1].Index != 851 {
		t.Fatalf("second log index mismatch: got %d, want 851", logs[1].Index)
	}
}

func TestRecreateReceiptsOrdersRealLiveTraceLogsByPosition(t *testing.T) {
	raw, err := os.ReadFile("testdata/polygon_live_trace_position.json")
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	var response struct {
		Result models.FullTraceResult `json:"result"`
	}
	if err := json.Unmarshal(raw, &response); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	logs, nextIndex := collectReceiptLogsFromTrace(&response.Result, common.HexToHash("0x73e64a8932439b878a220c40d84e5521980fde7437be39452e724f580a398c02"), 0, common.Hash{}, 1, 123, 0)
	expectedTopics := []common.Hash{
		common.HexToHash("0xbb47ee3e183a558b1a2ff0874b079f3fc5478b7454eacf2bfc5af2ff5878f972"),
		common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"),
		common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
		common.HexToHash("0x0b0c36936a3b6852af5796e1c928e945f9b7d331d8be5de3055541899136d451"),
		common.HexToHash("0x49628fd1471006c1482da88028e9ce4dbb080b815c9b0344d39e5a8e6ec1419f"),
		common.HexToHash("0xe6497e3ee548a3372136af2fcb0696db31fc6cf20260707645068bd3fe97f3c4"),
		common.HexToHash("0x4dfe1bbbcf077ddc3e01291eea2d5c70c2b422b415d95645b9adcfd678cb1d63"),
	}
	if len(logs) != len(expectedTopics) {
		t.Fatalf("log count mismatch: got %d, want %d", len(logs), len(expectedTopics))
	}
	if nextIndex != uint(len(expectedTopics)) {
		t.Fatalf("next log index mismatch: got %d, want %d", nextIndex, len(expectedTopics))
	}
	for i, expectedTopic := range expectedTopics {
		if len(logs[i].Topics) == 0 || logs[i].Topics[0] != expectedTopic {
			t.Fatalf("log %d topic mismatch: got %v, want %s", i, logs[i].Topics, expectedTopic.Hex())
		}
		if logs[i].Index != uint(i) {
			t.Fatalf("log %d index mismatch: got %d, want %d", i, logs[i].Index, i)
		}
	}
}

func TestRecreateReceiptsPrefersCompleteTraceLogIndexes(t *testing.T) {
	var trace models.FullTraceResult
	if err := json.Unmarshal([]byte(`{
		"logs": [{
			"index": 1,
			"position": "0x0",
			"address": "0x0000000000000000000000000000000000000002"
		}],
		"calls": [{
			"logs": [{
				"index": 0,
				"position": "0x0",
				"address": "0x0000000000000000000000000000000000000001"
			}]
		}]
	}`), &trace); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	logs, _ := collectReceiptLogsFromTrace(&trace, common.Hash{}, 0, common.Hash{}, 1, 123, 0)
	if len(logs) != 2 {
		t.Fatalf("log count mismatch: got %d, want 2", len(logs))
	}
	if logs[0].Address != common.HexToAddress("0x0000000000000000000000000000000000000001") {
		t.Fatalf("first log address mismatch: got %s", logs[0].Address.Hex())
	}
	if logs[1].Address != common.HexToAddress("0x0000000000000000000000000000000000000002") {
		t.Fatalf("second log address mismatch: got %s", logs[1].Address.Hex())
	}
}

func assertSyntheticPolygonFeeLog(t *testing.T, log *types.Log, txHash common.Hash, txIndex uint, logIndex uint) {
	t.Helper()
	if log == nil {
		t.Fatalf("missing synthetic polygon fee log")
	}
	if log.Address != models.POLLogTokenAddress {
		t.Fatalf("synthetic log address mismatch: got %s, want %s", log.Address.Hex(), models.POLLogTokenAddress.Hex())
	}
	if log.TxHash != txHash {
		t.Fatalf("synthetic log tx hash mismatch: got %s, want %s", log.TxHash.Hex(), txHash.Hex())
	}
	if log.TxIndex != txIndex {
		t.Fatalf("synthetic log tx index mismatch: got %d, want %d", log.TxIndex, txIndex)
	}
	if log.Index != logIndex {
		t.Fatalf("synthetic log index mismatch: got %d, want %d", log.Index, logIndex)
	}
	if len(log.Topics) != 0 || len(log.Data) != 0 {
		t.Fatalf("synthetic log must not invent topics/data")
	}
}
