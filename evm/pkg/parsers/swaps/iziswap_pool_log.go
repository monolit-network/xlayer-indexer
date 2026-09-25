package swaps

import (
	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"

	"go.uber.org/zap"
)

// IziSwapLogIdentifier matches any tx whose receipt contains an iZiSwap
// Swap event — topic-only, pool-address-agnostic (same philosophy as the
// uniswap identifiers). Amount/side reconstruction is delegated to the
// transfers heuristic.
type IziSwapLogIdentifier struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewIziSwapLogIdentifier(evmClient *evmclient.Client, logger *zap.Logger) *IziSwapLogIdentifier {
	return &IziSwapLogIdentifier{evmClient: evmClient, logger: logger.Named("iziswap_identifier")}
}

func (i *IziSwapLogIdentifier) CheckTx(req parsers.ParseTxRequest) (bool, error) {
	if req.Receipt == nil {
		return false, nil
	}
	for _, lg := range req.Receipt.Logs {
		if lg != nil && len(lg.Topics) > 0 && lg.Topics[0] == models.IziSwapEventTopic {
			return true, nil
		}
	}
	return false, nil
}

func (i *IziSwapLogIdentifier) Source() string        { return "iziswap" }
func (i *IziSwapLogIdentifier) ShouldStop() bool      { return true }
func (i *IziSwapLogIdentifier) GetParserName() string { return "swaps_by_transfers_parser_v2" }
