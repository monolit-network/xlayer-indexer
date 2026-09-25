package swaps

import (
	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"go.uber.org/zap"
)

type UniswapPoolV1Identifier struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewUniswapPoolV1Identifier(evmClient *evmclient.Client, logger *zap.Logger) *UniswapPoolV1Identifier {
	return &UniswapPoolV1Identifier{evmClient: evmClient, logger: logger}
}

func (i *UniswapPoolV1Identifier) Source() string {
	return "uniswap_pool_v1"
}

func (i *UniswapPoolV1Identifier) ShouldStop() bool {
	return true
}

func (i *UniswapPoolV1Identifier) CheckTx(req parsers.ParseTxRequest) (bool, error) {
	if req.Tx.To() == nil {
		return false, nil
	}

	factory, err := i.evmClient.GetFactory(*req.Tx.To())
	if err != nil {
		return false, err
	}

	chain, ok := models.IdToChain[models.ChainId(req.Tx.ChainId().Uint64())]
	if !ok {
		return false, nil
	}
	if chain != models.ChainETH {
		return false, nil
	}

	if models.UniswapV1FactoryAddressEth != factory {
		i.logger.Info("factory is not uniswap v1 factory", zap.String("factory", factory.Hex()))
		return false, nil
	}

	if _, ok := models.UniswapV1SwapSelectors[models.SelectorFromTx(req.Tx)]; !ok {
		i.logger.Info("data is not uniswap v1 swap event selector hash", zap.String("data", models.SelectorFromTx(req.Tx)))
		return false, nil
	}

	return true, nil
}

func (i *UniswapPoolV1Identifier) GetParserName() string {
	return ByTransfersParser{}.Name()
}
