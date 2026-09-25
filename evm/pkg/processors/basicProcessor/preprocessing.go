package basicProcessor

import (
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"github.com/ethereum/go-ethereum/core/types"
)

type requestPreprocessor func(req parsers.ParseTxRequest) (parsers.ParseTxRequest, error)

var allPreprocessFuncs = map[models.Chain][]requestPreprocessor{
	models.ChainPolygon: {preprocessPolygonRequest},
}

func preprocessPolygonRequest(req parsers.ParseTxRequest) (parsers.ParseTxRequest, error) {
	logs := make([]*types.Log, 0)
	for _, log := range req.Receipt.Logs {
		if log == nil {
			continue
		}
		if log.Address == models.POLLogTokenAddress {
			continue
		}

		logs = append(logs, log)
	}
	req.Receipt.Logs = logs
	return req, nil
}
