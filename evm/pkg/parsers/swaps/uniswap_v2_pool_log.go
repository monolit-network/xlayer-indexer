package swaps

import (
	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"go.uber.org/zap"
)

// UniswapV2LogIdentifier fires when the tx receipt contains a Uniswap v2
// pool Swap event (any pool / v2-fork like swaphood). Delegates to
// swaps_by_transfers_parser_v2.
type UniswapV2LogIdentifier struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewUniswapV2LogIdentifier(evmClient *evmclient.Client, logger *zap.Logger) *UniswapV2LogIdentifier {
	return &UniswapV2LogIdentifier{evmClient: evmClient, logger: logger}
}

func (i *UniswapV2LogIdentifier) Source() string { return "uniswap_v2" }

func (i *UniswapV2LogIdentifier) ShouldStop() bool { return true }

func (i *UniswapV2LogIdentifier) GetParserName() string { return "swaps_by_transfers_parser_v2" }

func (i *UniswapV2LogIdentifier) CheckTx(req parsers.ParseTxRequest) (bool, error) {
	if req.Receipt == nil || req.Tx == nil {
		return false, nil
	}
	for _, lg := range req.Receipt.Logs {
		if lg == nil {
			continue
		}
		if len(lg.Topics) > 0 && lg.Topics[0] == models.UniswapV2SwapEventTopic {
			return true, nil
		}
	}
	return false, nil
}
