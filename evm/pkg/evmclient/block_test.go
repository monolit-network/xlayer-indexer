package evmclient

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

const zeroBloom = "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"

func TestBlockFromRpcJSONFiltersZeroPlaceholderTransactions(t *testing.T) {
	raw := json.RawMessage(`{
		"difficulty":"0x7",
		"extraData":"0xd58301090083626f7286676f312e3133856c696e7578000000000000000000003fd2069222b186c3cb4937de95cb605c4a201421ed9532e37e6d12784e33616f0b6b753f6d9a85dfe31328cab03badbb650076e49f6033863697b1c4937c1ecc00",
		"gasLimit":"0x1312d00",
		"gasUsed":"0x0",
		"hash":"0x2e3f40ddbdb9d5316356af2ac7ac94ba3e5ccbb6fc3b9a73671cdd5bb27570e5",
		"logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"miner":"0x0000000000000000000000000000000000000000",
		"mixHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
		"nonce":"0x0000000000000000",
		"number":"0x900",
		"parentHash":"0x9b863b8348e030fc6f2a566b7ad2914d4d9f39e93d0454e978e8509d3d14b91a",
		"receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"size":"0x261",
		"stateRoot":"0x47135251b4277da678ce402c2af63f3727f7c2297a19550ac4cc80f985073600",
		"timestamp":"0x5ed29c96",
		"transactions":[
			{
				"blockHash":"0x2e3f40ddbdb9d5316356af2ac7ac94ba3e5ccbb6fc3b9a73671cdd5bb27570e5",
				"blockNumber":"0x900",
				"from":"0x0000000000000000000000000000000000000000",
				"gas":"0x0",
				"gasPrice":"0x0",
				"hash":"0xef029df186fe80862876987b237d848734d534c1c857f49bb9606a1b38f2b6b8",
				"input":"0x",
				"nonce":"0x0",
				"to":"0x0000000000000000000000000000000000000000",
				"transactionIndex":"0x0",
				"value":"0x0",
				"type":"0x0",
				"v":"0x0",
				"r":"0x0",
				"s":"0x0"
			}
		],
		"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"uncles":[]
	}`)

	block, err := blockFromRpcJson(raw, true)
	if err != nil {
		t.Fatalf("blockFromRpcJson() error = %v", err)
	}
	if got := len(block.Transactions()); got != 0 {
		t.Fatalf("expected zero transactions after filtering placeholder tx, got %d", got)
	}
}

func TestBlockFromRpcJSONRejectsNonPlaceholderTransactionsWhenTxRootIsEmpty(t *testing.T) {
	raw := strings.ReplaceAll(`{
		"difficulty":"0x7",
		"extraData":"0xd58301090083626f7286676f312e3133856c696e7578000000000000000000003fd2069222b186c3cb4937de95cb605c4a201421ed9532e37e6d12784e33616f0b6b753f6d9a85dfe31328cab03badbb650076e49f6033863697b1c4937c1ecc00",
		"gasLimit":"0x1312d00",
		"gasUsed":"0x0",
		"hash":"0x2e3f40ddbdb9d5316356af2ac7ac94ba3e5ccbb6fc3b9a73671cdd5bb27570e5",
		"logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"miner":"0x0000000000000000000000000000000000000000",
		"mixHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
		"nonce":"0x0000000000000000",
		"number":"0x900",
		"parentHash":"0x9b863b8348e030fc6f2a566b7ad2914d4d9f39e93d0454e978e8509d3d14b91a",
		"receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"size":"0x261",
		"stateRoot":"0x47135251b4277da678ce402c2af63f3727f7c2297a19550ac4cc80f985073600",
		"timestamp":"0x5ed29c96",
		"transactions":[
			{
				"blockHash":"0x2e3f40ddbdb9d5316356af2ac7ac94ba3e5ccbb6fc3b9a73671cdd5bb27570e5",
				"blockNumber":"0x900",
				"from":"0x1111111111111111111111111111111111111111",
				"gas":"0x0",
				"gasPrice":"0x0",
				"hash":"0xef029df186fe80862876987b237d848734d534c1c857f49bb9606a1b38f2b6b8",
				"input":"0x",
				"nonce":"0x0",
				"to":"0x0000000000000000000000000000000000000000",
				"transactionIndex":"0x0",
				"value":"0x0",
				"type":"0x0",
				"v":"0x0",
				"r":"0x0",
				"s":"0x0"
			}
		],
		"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"uncles":[]
	}`, "\t", "")

	if _, err := blockFromRpcJson(json.RawMessage(raw), true); err == nil {
		t.Fatal("expected blockFromRpcJson() to reject non-placeholder transaction when tx root is empty")
	}
}

func TestBlockFromRpcJSONFiltersPlaceholderTransactionInNonEmptyBlock(t *testing.T) {
	raw := json.RawMessage(strings.ReplaceAll(`{
		"difficulty":"0x13",
		"extraData":"0x00",
		"gasLimit":"0x13a8613",
		"gasUsed":"0x5208",
		"hash":"0x6fadf094c2274a059e958bb2f4b46ebd5af2b3c51f7a9d2462c6a146e70207dc",
		"logsBloom":"__ZERO_BLOOM__",
		"miner":"0x0000000000000000000000000000000000000000",
		"mixHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
		"nonce":"0x0000000000000000",
		"number":"0x10ce400",
		"parentHash":"0xc2d4d0a5f231d62b1553913bb821d4a4d39a94db667b5bbfe5817ef98eaa9863",
		"receiptsRoot":"0x5495bffa0c7afe1f1e5342e86708356afe20bea8ffefd77e459f5e9f6a7bd9c4",
		"sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"size":"0x65ff",
		"stateRoot":"0x49a6404f1da6a90dd02cba217dc086373c262bc1c0d7b930b78decabef1839a8",
		"timestamp":"0x610ae517",
		"transactions":[
			{
				"blockHash":"0x6fadf094c2274a059e958bb2f4b46ebd5af2b3c51f7a9d2462c6a146e70207dc",
				"blockNumber":"0x10ce400",
				"from":"0xbc9f62c241f9c1273022709a0e37a6a2caf4cf2f",
				"gas":"0x5208",
				"gasPrice":"0xa190c761",
				"hash":"0x45070b344a09dc9a78208b9760edf98ad82311bbb235a455cad15ee6e510a450",
				"input":"0x",
				"nonce":"0x2321",
				"to":"0xbc9f62c241f9c1273022709a0e37a6a2caf4cf2f",
				"transactionIndex":"0x30",
				"value":"0x0",
				"type":"0x0",
				"v":"0x136",
				"r":"0xc12f01afb1fc02bf7ad453cab7321bf33840dab3dca15c78a82790ecf4655d2",
				"s":"0x1baec021d1f6843601a18b2689cac5407f9e232cdf0ccec85f4b3ab3c82f8772"
			},
			{
				"blockHash":"0x6fadf094c2274a059e958bb2f4b46ebd5af2b3c51f7a9d2462c6a146e70207dc",
				"blockNumber":"0x10ce400",
				"from":"0x0000000000000000000000000000000000000000",
				"gas":"0x0",
				"gasPrice":"0x0",
				"hash":"0xa5d962cb5e45e858c0715a6b0eb3140ff89b3a9925c257ade5f9d2afef3d55e0",
				"input":"0x",
				"nonce":"0x0",
				"to":"0x0000000000000000000000000000000000000000",
				"transactionIndex":"0x34",
				"value":"0x0",
				"type":"0x0",
				"v":"0x0",
				"r":"0x0",
				"s":"0x0"
			}
		],
		"transactionsRoot":"0xbdf65fca258b47a90df30b70c9a404a297cbef7a3d756b9fdde37b82099f93f9",
		"uncles":[]
	}`, "__ZERO_BLOOM__", zeroBloom))

	block, err := blockFromRpcJson(raw, true)
	if err != nil {
		t.Fatalf("blockFromRpcJson() error = %v", err)
	}
	if got := len(block.Transactions()); got != 1 {
		t.Fatalf("expected placeholder tx to be filtered from non-empty block, got %d transactions", got)
	}
}

func TestFilterZeroPlaceholderReceipts(t *testing.T) {
	raw := []byte(strings.ReplaceAll(`[
		{
			"blockHash":"0x6fadf094c2274a059e958bb2f4b46ebd5af2b3c51f7a9d2462c6a146e70207dc",
			"blockNumber":"0x10ce400",
			"cumulativeGasUsed":"0x5208",
			"effectiveGasPrice":"0xa190c761",
			"gasUsed":"0x5208",
			"logs":[],
			"logsBloom":"__ZERO_BLOOM__",
			"status":"0x1",
			"transactionHash":"0x45070b344a09dc9a78208b9760edf98ad82311bbb235a455cad15ee6e510a450",
			"transactionIndex":"0x30",
			"from":"0xbc9f62c241f9c1273022709a0e37a6a2caf4cf2f",
			"to":"0xbc9f62c241f9c1273022709a0e37a6a2caf4cf2f",
			"type":"0x0"
		},
		{
			"blockHash":"0x6fadf094c2274a059e958bb2f4b46ebd5af2b3c51f7a9d2462c6a146e70207dc",
			"blockNumber":"0x10ce400",
			"cumulativeGasUsed":"0x0",
			"effectiveGasPrice":"0x0",
			"gasUsed":"0x0",
			"logs":[{"address":"0x1111111111111111111111111111111111111111","topics":[],"data":"0x","blockNumber":"0x10ce400","transactionHash":"0xa5d962cb5e45e858c0715a6b0eb3140ff89b3a9925c257ade5f9d2afef3d55e0","transactionIndex":"0x34","blockHash":"0x6fadf094c2274a059e958bb2f4b46ebd5af2b3c51f7a9d2462c6a146e70207dc","logIndex":"0x0","removed":false}],
			"logsBloom":"__ZERO_BLOOM__",
			"status":"0x1",
			"transactionHash":"0xa5d962cb5e45e858c0715a6b0eb3140ff89b3a9925c257ade5f9d2afef3d55e0",
			"transactionIndex":"0x34",
			"from":"0x0000000000000000000000000000000000000000",
			"to":"0x0000000000000000000000000000000000000000",
			"type":"0x0"
		}
	]`, "__ZERO_BLOOM__", zeroBloom))

	var receipts []rpcReceipt
	if err := json.Unmarshal(raw, &receipts); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	filtered := filterZeroPlaceholderReceipts(receipts)
	if got := len(filtered); got != 1 {
		t.Fatalf("expected placeholder receipt to be filtered, got %d receipts", got)
	}
	if _, ok := filtered[common.HexToHash("0x45070b344a09dc9a78208b9760edf98ad82311bbb235a455cad15ee6e510a450")]; !ok {
		t.Fatal("expected receipt to be indexed by tx hash")
	}
}
