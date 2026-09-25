package swaps

import (
	"math/big"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type UniswapPoolV2Identifier struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewUniswapPoolV2Identifier(evmClient *evmclient.Client, logger *zap.Logger) *UniswapPoolV2Identifier {
	return &UniswapPoolV2Identifier{evmClient: evmClient, logger: logger}
}

func (i *UniswapPoolV2Identifier) Source() string {
	return "uniswap_pool_v2"
}

func (i *UniswapPoolV2Identifier) ShouldStop() bool {
	return true
}

func (i *UniswapPoolV2Identifier) CheckTx(req parsers.ParseTxRequest) (bool, error) {
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

	if factory != models.UniswapV2FactoryAddresses[chain] {
		i.logger.Info("factory is not uniswap v2 factory", zap.String("factory", factory.Hex()))
		return false, nil
	}

	if models.SelectorFromTx(req.Tx) != models.UniswapV2SwapFuncSelectorHash {
		i.logger.Info("data is not uniswap v2 swap event selector hash", zap.String("data", models.SelectorFromTx(req.Tx)), zap.String("selector hash", string(models.UniswapV2SwapFuncSelectorHash)))
		return false, nil
	}

	return true, nil
}

func (i *UniswapPoolV2Identifier) GetParserName() string {
	return "uniswap_pool_v2_parser"
}

type UniswapPoolV2Parser struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewUniswapPoolV2Parser(evmClient *evmclient.Client, logger *zap.Logger) *UniswapPoolV2Parser {
	return &UniswapPoolV2Parser{evmClient: evmClient, logger: logger}
}

func (p UniswapPoolV2Parser) Name() string {
	return "uniswap_pool_v2_parser"
}

func (p *UniswapPoolV2Parser) ProducedEventsType() models.EventType {
	return models.EventTypeSwap
}

func (p *UniswapPoolV2Parser) RequiredData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesBlockHeader
}

func (p *UniswapPoolV2Parser) OptionalData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesTraces
}

func (p *UniswapPoolV2Parser) ParseTx(req parsers.ParseTxRequest) ([]models.Event, error) {
	if err := parsers.ValidateParseTxRequest(req, p.RequiredData()); err != nil {
		return nil, err
	}

	inputParams := struct {
		Amount0Out *big.Int
		Amount1Out *big.Int
		To         common.Address
		Data       []byte
	}{}

	err := models.UnpackCalldata(models.UniswapV2SwapFuncABI, "swap", req.Tx.Data(), &inputParams)
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
