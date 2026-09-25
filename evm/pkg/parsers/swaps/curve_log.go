package swaps

import (
	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"

	"go.uber.org/zap"
)

// CurvePoolLogIdentifier matches any tx whose receipt contains an iZiSwap
// Swap event — topic-only, pool-address-agnostic (same philosophy as the
// uniswap identifiers). Amount/side reconstruction is delegated to the
// transfers heuristic.
type CurvePoolLogIdentifier struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewCurvePoolLogIdentifier(evmClient *evmclient.Client, logger *zap.Logger) *CurvePoolLogIdentifier {
	return &CurvePoolLogIdentifier{evmClient: evmClient, logger: logger.Named("iziswap_identifier")}
}

func (i *CurvePoolLogIdentifier) CheckTx(req parsers.ParseTxRequest) (bool, error) {
	if req.Receipt == nil {
		return false, nil
	}
	for _, lg := range req.Receipt.Logs {
		if lg != nil && len(lg.Topics) > 0 && lg.Topics[0] == models.CurveExchangeEventTopic {
			return true, nil
		}
	}
	return false, nil
}

func (i *CurvePoolLogIdentifier) Source() string        { return "curve" }
func (i *CurvePoolLogIdentifier) ShouldStop() bool      { return true }
func (i *CurvePoolLogIdentifier) GetParserName() string { return "swaps_by_transfers_parser_v2" }
