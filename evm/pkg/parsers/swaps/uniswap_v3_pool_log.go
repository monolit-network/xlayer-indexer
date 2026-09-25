package swaps

import (
	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"go.uber.org/zap"
)

// UniswapV3LogIdentifier fires when the tx receipt contains a Uniswap v3
// pool Swap event (any pool), regardless of initiator. Covers standalone
// Uniswap v3 AND Pons V1 (which trades in locked Uniswap v3 pools).
// Delegates amount extraction to swaps_by_transfers_parser_v2.
type UniswapV3LogIdentifier struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewUniswapV3LogIdentifier(evmClient *evmclient.Client, logger *zap.Logger) *UniswapV3LogIdentifier {
	return &UniswapV3LogIdentifier{evmClient: evmClient, logger: logger}
}

func (i *UniswapV3LogIdentifier) Source() string { return "uniswap_v3" }

func (i *UniswapV3LogIdentifier) ShouldStop() bool { return true }

func (i *UniswapV3LogIdentifier) GetParserName() string { return "swaps_by_transfers_parser_v2" }

func (i *UniswapV3LogIdentifier) CheckTx(req parsers.ParseTxRequest) (bool, error) {
	if req.Receipt == nil || req.Tx == nil {
		return false, nil
	}
	for _, lg := range req.Receipt.Logs {
		if lg == nil {
			continue
		}
		if len(lg.Topics) > 0 && lg.Topics[0] == models.UniswapV3SwapEventTopic {
			return true, nil
		}
	}
	return false, nil
}
