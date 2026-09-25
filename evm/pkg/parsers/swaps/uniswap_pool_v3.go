package swaps

import (
	"math/big"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type UniswapPoolV3Identifier struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewUniswapPoolV3Identifier(evmClient *evmclient.Client, logger *zap.Logger) *UniswapPoolV3Identifier {
	return &UniswapPoolV3Identifier{evmClient: evmClient, logger: logger}
}

func (i *UniswapPoolV3Identifier) Source() string {
	return "uniswap_pool_v3"
}

func (i *UniswapPoolV3Identifier) ShouldStop() bool {
	return true
}

func (i *UniswapPoolV3Identifier) CheckTx(req parsers.ParseTxRequest) (bool, error) {
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

	if factory != models.UniswapV3FactoryAddresses[chain] {
		i.logger.Info("factory is not uniswap v3 factory", zap.String("factory", factory.Hex()))
		return false, nil
	}

	if models.SelectorFromTx(req.Tx) != models.UniswapV3SwapFuncSelectorHash {
		i.logger.Info("data is not uniswap v3 swap event selector hash", zap.String("data", models.SelectorFromTx(req.Tx)), zap.String("selector hash", string(models.UniswapV3SwapFuncSelectorHash)))
		return false, nil
	}

	return true, nil
}

func (i *UniswapPoolV3Identifier) GetParserName() string {
	return "uniswap_pool_v3_parser"
}

type UniswapPoolV3Parser struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewUniswapPoolV3Parser(evmClient *evmclient.Client, logger *zap.Logger) *UniswapPoolV3Parser {
	return &UniswapPoolV3Parser{evmClient: evmClient, logger: logger}
}

func (p UniswapPoolV3Parser) Name() string {
	return "uniswap_pool_v3_parser"
}

func (p *UniswapPoolV3Parser) ProducedEventsType() models.EventType {
	return models.EventTypeSwap
}

func (p *UniswapPoolV3Parser) RequiredData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesBlockHeader
}

func (p *UniswapPoolV3Parser) OptionalData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesTraces
}

func (p *UniswapPoolV3Parser) ParseTx(req parsers.ParseTxRequest) ([]models.Event, error) {
	if err := parsers.ValidateParseTxRequest(req, p.RequiredData()); err != nil {
		return nil, err
	}

	inputParams := struct {
		Token             common.Address
		ZeroForOne        bool
		AmountSpecified   *big.Int
		SqrtPriceLimitX96 *big.Int
		Data              []byte
	}{}

	err := models.UnpackCalldata(models.UniswapV3SwapFuncABI, "swap", req.Tx.Data(), &inputParams)
	if err != nil {
		return nil, err
	}

	if len(inputParams.Data) != 0 {
		// skip flash loans
		return []models.Event{}, nil
	}

	swapsByTransfersParser := NewByTransfersParser(p.evmClient, p.logger)
	return swapsByTransfersParser.ParseTx(req)
}
