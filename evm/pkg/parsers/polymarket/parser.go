package polymarket

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

const (
	ordersParserName = "polymarket_orders_parser"
)

type PolymarketOrderIdentifier struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewPolymarketOrderIdentifier(evmClient *evmclient.Client, logger *zap.Logger) *PolymarketOrderIdentifier {
	return &PolymarketOrderIdentifier{evmClient: evmClient, logger: logger}
}

func (i *PolymarketOrderIdentifier) Source() string {
	return "polymarket_order"
}

func (i *PolymarketOrderIdentifier) ShouldStop() bool {
	return false
}

func (i *PolymarketOrderIdentifier) GetParserName() string {
	return ordersParserName
}

func (i *PolymarketOrderIdentifier) CheckTx(req parsers.ParseTxRequest) (bool, error) {
	if req.Tx.To() == nil {
		return false, nil
	}

	for _, log := range req.Receipt.Logs {
		if log == nil {
			continue
		}
		if !isAllowedPolymarketOrderLog(log) {
			continue
		}
		return true, nil
	}

	return false, nil
}

type OrdersParser struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewOrdersParser(evmClient *evmclient.Client, logger *zap.Logger) *OrdersParser {
	return &OrdersParser{
		evmClient: evmClient,
		logger:    logger.Named("polymarket-orders-parser"),
	}
}

func (p *OrdersParser) Name() string {
	return ordersParserName
}

func (p *OrdersParser) ProducedEventsType() models.EventType {
	return models.EventTypePolymarketOrder
}

func (p *OrdersParser) RequiredData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesBlockHeader
}

func (p *OrdersParser) OptionalData() evmclient.RequiredDataTypes {
	return 0
}

func (p *OrdersParser) ParseTx(req parsers.ParseTxRequest) ([]models.Event, error) {
	if err := parsers.ValidateParseTxRequest(req, p.RequiredData()); err != nil {
		return nil, err
	}

	baseEvent, err := parsers.MakeBaseEvent(req, p.evmClient)
	if err != nil {
		return nil, err
	}

	var events []models.Event
	for _, log := range req.Receipt.Logs {
		if log == nil {
			continue
		}
		if !isAllowedPolymarketOrderLog(log) {
			continue
		}
		event, err := p.parseOrderFilledLog(baseEvent, log)
		if err != nil {
			p.logger.Warn("polymarket_order_log_decode_failed",
				zap.Error(err),
				zap.String("tx_hash", req.Tx.Hash().Hex()),
				zap.Uint("log_index", log.Index),
			)
			continue
		}
		events = append(events, event)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("no polymarket order filled events found")
	}

	return events, nil
}

type orderFilledIndexedArgs struct {
	OrderHash common.Hash    `abi:"orderHash"`
	Maker     common.Address `abi:"maker"`
	Taker     common.Address `abi:"taker"`
}

type orderFilledDataArgs struct {
	MakerAssetID      *big.Int `abi:"makerAssetId"`
	TakerAssetID      *big.Int `abi:"takerAssetId"`
	MakerAmountFilled *big.Int `abi:"makerAmountFilled"`
	TakerAmountFilled *big.Int `abi:"takerAmountFilled"`
	Fee               *big.Int `abi:"fee"`
}

func (p *OrdersParser) parseOrderFilledLog(base models.BaseEvent, log *types.Log) (*models.PolymarketOrderEvent, error) {
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	if log.Topics[0] == models.PolymarketOrderFilledEventV2SelectorHash {
		orderFilled, err := parseOrderFilledV2Log(log)
		if err != nil {
			return nil, err
		}

		return &models.PolymarketOrderEvent{
			BaseEvent:         base,
			OrderHash:         orderFilled.OrderHash,
			Maker:             orderFilled.Maker,
			Taker:             orderFilled.Taker,
			MakerAssetID:      orderFilled.MakerAssetID,
			TakerAssetID:      orderFilled.TakerAssetID,
			MakerAmountFilled: orderFilled.MakerAmountFilled,
			TakerAmountFilled: orderFilled.TakerAmountFilled,
			Fee:               orderFilled.Fee,
			LogIndex:          log.Index,
		}, nil
	}

	indexedArgs := orderFilledIndexedArgs{
		OrderHash: log.Topics[1],
		Maker:     common.BytesToAddress(log.Topics[2].Bytes()),
		Taker:     common.BytesToAddress(log.Topics[3].Bytes()),
	}

	dataArgs := orderFilledDataArgs{}
	if err := models.PolymarketOrderFilledEventABI.UnpackIntoInterface(&dataArgs, "OrderFilled", log.Data); err != nil {
		return nil, fmt.Errorf("unpack data: %w", err)
	}

	return &models.PolymarketOrderEvent{
		BaseEvent:         base,
		OrderHash:         indexedArgs.OrderHash,
		Maker:             indexedArgs.Maker,
		Taker:             indexedArgs.Taker,
		MakerAssetID:      dataArgs.MakerAssetID,
		TakerAssetID:      dataArgs.TakerAssetID,
		MakerAmountFilled: dataArgs.MakerAmountFilled,
		TakerAmountFilled: dataArgs.TakerAmountFilled,
		Fee:               dataArgs.Fee,
		LogIndex:          log.Index,
	}, nil
}

func isAllowedPolymarketOrderLog(log *types.Log) bool {
	isExchangeAddress := log.Address == models.PolymarketCTFExchangeAddress ||
		log.Address == models.PolymarketCTFExchangeAddressV2 ||
		log.Address == models.PolymarketNegRiskCTFExchangeAddress ||
		log.Address == models.PolymarketNegRiskCTFExchangeAddressV2
	if !isExchangeAddress || len(log.Topics) != 4 {
		return false
	}

	return log.Topics[0] == models.PolymarketOrderFilledEventSelectorHash ||
		log.Topics[0] == models.PolymarketOrderFilledEventV2SelectorHash
}
