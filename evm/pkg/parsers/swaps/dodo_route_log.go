package swaps

import (
	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"

	"go.uber.org/zap"
)

// DodoRouteLogIdentifier matches any tx whose receipt contains an iZiSwap
// Swap event — topic-only, pool-address-agnostic (same philosophy as the
// uniswap identifiers). Amount/side reconstruction is delegated to the
// transfers heuristic.
type DodoRouteLogIdentifier struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewDodoRouteLogIdentifier(evmClient *evmclient.Client, logger *zap.Logger) *DodoRouteLogIdentifier {
	return &DodoRouteLogIdentifier{evmClient: evmClient, logger: logger.Named("iziswap_identifier")}
}

func (i *DodoRouteLogIdentifier) CheckTx(req parsers.ParseTxRequest) (bool, error) {
	if req.Receipt == nil {
		return false, nil
	}
	for _, lg := range req.Receipt.Logs {
		if lg != nil && len(lg.Topics) > 0 && lg.Topics[0] == models.DodoRouteEventTopic {
			return true, nil
		}
	}
	return false, nil
}

func (i *DodoRouteLogIdentifier) Source() string        { return "dodo_route" }
func (i *DodoRouteLogIdentifier) ShouldStop() bool      { return true }
func (i *DodoRouteLogIdentifier) GetParserName() string { return "swaps_by_transfers_parser_v2" }
