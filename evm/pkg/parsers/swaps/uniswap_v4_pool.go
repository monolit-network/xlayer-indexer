package swaps

import (
	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"go.uber.org/zap"
)

// UniswapV4Identifier fires when the tx receipt contains a Uniswap v4
// PoolManager Swap event, regardless of which router/bot initiated it.
// It delegates amount extraction to swaps_by_transfers_parser_v2.
type UniswapV4Identifier struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewUniswapV4Identifier(evmClient *evmclient.Client, logger *zap.Logger) *UniswapV4Identifier {
	return &UniswapV4Identifier{evmClient: evmClient, logger: logger}
}

func (i *UniswapV4Identifier) Source() string { return "uniswap_v4" }

func (i *UniswapV4Identifier) ShouldStop() bool { return true }

func (i *UniswapV4Identifier) GetParserName() string { return "swaps_by_transfers_parser_v2" }

func (i *UniswapV4Identifier) CheckTx(req parsers.ParseTxRequest) (bool, error) {
	if req.Receipt == nil || req.Tx == nil {
		return false, nil
	}
	chain, ok := models.IdToChain[models.ChainId(req.Tx.ChainId().Uint64())]
	if !ok {
		return false, nil
	}
	pm, ok := models.UniswapV4PoolManagerAddresses[chain]
	if !ok {
		return false, nil
	}
	for _, lg := range req.Receipt.Logs {
		if lg == nil {
			continue
		}
		if lg.Address == pm && len(lg.Topics) > 0 && lg.Topics[0] == models.UniswapV4SwapEventTopic {
			return true, nil
		}
	}
	return false, nil
}
