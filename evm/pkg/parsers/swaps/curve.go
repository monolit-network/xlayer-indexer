package swaps

import (
	"fmt"
	"math/big"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

type CurveParser struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewCurveParser(evmClient *evmclient.Client, logger *zap.Logger) *CurveParser {
	return &CurveParser{evmClient: evmClient, logger: logger}
}

func (p CurveParser) Name() string {
	return "curve_parser"
}

func (p *CurveParser) ProducedEventsType() models.EventType {
	return models.EventTypeSwap
}

func (p *CurveParser) RequiredData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesBlockHeader
}

func (p *CurveParser) OptionalData() evmclient.RequiredDataTypes {
	return 0
}

func (p *CurveParser) ParseTx(req parsers.ParseTxRequest) ([]models.Event, error) {
	if err := parsers.ValidateParseTxRequest(req, p.RequiredData()); err != nil {
		return nil, err
	}
	baseEvent, err := parsers.MakeBaseEvent(req, p.evmClient)
	if err != nil {
		return nil, err
	}

	var swapEvent models.SwapEvent
	found := false
	for _, log := range req.Receipt.Logs {
		if swapEvent, err = p.tryParseCurveSwapLog(log); err == nil {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("no curve swap log found")
	}

	swapEvent.BaseEvent = baseEvent

	return []models.Event{&swapEvent}, nil
}

func (p *CurveParser) tryParseCurveSwapLog(log *types.Log) (models.SwapEvent, error) {
	if log.Topics[0] != models.CurveSwapEventSelectorHash {
		return models.SwapEvent{}, fmt.Errorf("log is not an Curve swap log")
	}

	parsedLog := struct {
		Route      [11]common.Address
		SwapParams [5][5]*big.Int
		Pools      [5]common.Address
		InAmount   *big.Int
		OutAmount  *big.Int
	}{}

	resultLog := models.SwapEvent{}

	err := models.CurveSwapEventABI.UnpackIntoInterface(&parsedLog, "Exchange", log.Data)
	if err != nil {
		return models.SwapEvent{}, fmt.Errorf("error unpacking log: %w", err)
	}

	firstToken := common.HexToAddress(parsedLog.Route[0].Hex())
	var lastToken common.Address
	for i := len(parsedLog.Route) - 1; i >= 0; i-- {
		if parsedLog.Route[i] != (common.Address{}) {
			lastToken = parsedLog.Route[i]
			break
		}
	}

	if firstToken == models.CurveETHTokenAddress {
		firstToken = common.Address{}
	}
	if lastToken == models.CurveETHTokenAddress {
		lastToken = common.Address{}
	}

	resultLog.Sender = common.HexToAddress(log.Topics[1].Hex())
	resultLog.Receiver = common.HexToAddress(log.Topics[2].Hex())
	resultLog.BaseCoin = firstToken
	resultLog.QuoteCoin = lastToken
	resultLog.BaseCoinAmount = parsedLog.InAmount
	resultLog.QuoteCoinAmount = parsedLog.OutAmount

	if d, err := p.evmClient.GetDecimals(firstToken); err == nil {
		resultLog.BaseCoinDecimals = &d
	}
	if d, err := p.evmClient.GetDecimals(lastToken); err == nil {
		resultLog.QuoteCoinDecimals = &d
	}

	return resultLog, nil
}
