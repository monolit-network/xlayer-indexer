package swaps

import (
	"slices"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"go.uber.org/zap"
)

type UniswapUniversalRouterParser struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewUniswapUniversalRouterParser(evmClient *evmclient.Client, logger *zap.Logger) *UniswapUniversalRouterParser {
	return &UniswapUniversalRouterParser{evmClient: evmClient, logger: logger}
}

func (p UniswapUniversalRouterParser) Name() string {
	return "uniswap_universal_router_parser"
}

func (p *UniswapUniversalRouterParser) ProducedEventsType() models.EventType {
	return models.EventTypeSwap
}

func (p *UniswapUniversalRouterParser) RequiredData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesBlockHeader | evmclient.RequiredDataTypesTraces
}

func (p *UniswapUniversalRouterParser) OptionalData() evmclient.RequiredDataTypes {
	return 0
}

func (p *UniswapUniversalRouterParser) ParseTx(req parsers.ParseTxRequest) ([]models.Event, error) {
	if err := parsers.ValidateParseTxRequest(req, p.RequiredData()); err != nil {
		return nil, err
	}

	inputParams := struct {
		Commands []byte
		Inputs   [][]byte
	}{}

	err := models.UnpackCalldata(models.UniswapUniversalRouterExecuteFuncABI, "execute", req.Tx.Data(), &inputParams)
	if err != nil {
		return nil, err
	}

	swapFound := false
	allowedCommands := []byte{
		0x00, // V3_SWAP_EXACT_IN
		0x01, // V3_SWAP_EXACT_OUT
		0x08, // V2_SWAP_EXACT_IN
		0x09, // V2_SWAP_EXACT_OUT
		0x10, // V4_SWAP
	}
	for _, command := range inputParams.Commands {
		if slices.Contains(allowedCommands, command) {
			swapFound = true
			break
		}
	}
	if !swapFound {
		p.logger.Debug("swap not found in uniswap universal router", zap.Any("commands", inputParams.Commands))
		return nil, nil
	}

	swapsByTransfersParser := NewByTransfersParserV2(p.evmClient, p.logger)
	return swapsByTransfersParser.ParseTx(req)
}
