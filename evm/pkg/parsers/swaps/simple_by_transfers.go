package swaps

import (
	"fmt"
	"math/big"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type SimpleByTransfersParser struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewSimpleByTransfersParser(evmClient *evmclient.Client, logger *zap.Logger) *SimpleByTransfersParser {
	p := &SimpleByTransfersParser{evmClient: evmClient}
	p.logger = logger.Named(p.Name())
	return p
}

func (p SimpleByTransfersParser) Name() string {
	return "simple_by_transfers_parser"
}

func (p *SimpleByTransfersParser) ProducedEventsType() models.EventType {
	return models.EventTypeSwap
}

func (p *SimpleByTransfersParser) RequiredData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesBlockHeader | evmclient.RequiredDataTypesTraces
}

func (p *SimpleByTransfersParser) OptionalData() evmclient.RequiredDataTypes {
	return 0
}

func (p *SimpleByTransfersParser) ParseTx(req parsers.ParseTxRequest) ([]models.Event, error) {
	if err := parsers.ValidateParseTxRequest(req, p.RequiredData()); err != nil {
		return nil, err
	}
	baseEvent, err := parsers.MakeBaseEvent(req, p.evmClient)
	if err != nil {
		return nil, err
	}

	transfers, err := models.FindAllTransfers(req.Receipt, req.Tx, req.Traces, p.evmClient.ChainID)
	if err != nil {
		return nil, err
	}

	if len(transfers) != 2 {
		return nil, fmt.Errorf("expected 2 transfers, got %d", len(transfers))
	}

	_, _, netFlows := aggregateFlows(transfers)

	senderFlows, ok := netFlows[baseEvent.TxFromAddress]
	if !ok {
		return nil, fmt.Errorf("sender not found in net flows")
	}

	if len(senderFlows) != 2 {
		return nil, fmt.Errorf("sender has %d flows, expected 2", len(senderFlows))
	}

	var baseToken, quoteToken common.Address
	var baseAmount, quoteAmount *big.Int
	for token, amount := range senderFlows {
		sign := amount.Sign()
		if sign == 0 {
			continue
		}
		if sign < 0 {
			baseToken = token
			baseAmount = big.NewInt(0).Abs(amount)
		} else {
			quoteToken = token
			quoteAmount = amount
		}
	}

	if baseToken == (common.Address{}) || quoteToken == (common.Address{}) {
		return nil, fmt.Errorf("base token or quote token not found")
	}

	if baseAmount == nil {
		return nil, fmt.Errorf("base amount not found")
	}

	if quoteAmount == nil {
		return nil, fmt.Errorf("quote amount not found")
	}

	var baseDecimals *uint8
	var quoteDecimals *uint8

	if d, derr := p.evmClient.GetDecimals(baseToken); derr == nil {
		baseDecimals = &d
	} else {
		p.logger.Debug("get decimals failed for base token", zap.String("token", baseToken.Hex()), zap.Error(derr))
	}

	if d, derr := p.evmClient.GetDecimals(quoteToken); derr == nil {
		quoteDecimals = &d
	} else {
		p.logger.Debug("get decimals failed for quote token", zap.String("token", quoteToken.Hex()), zap.Error(derr))
	}

	return []models.Event{&models.SwapEvent{
		BaseEvent:         baseEvent,
		BaseCoin:          baseToken,
		QuoteCoin:         quoteToken,
		BaseCoinAmount:    baseAmount,
		BaseCoinDecimals:  baseDecimals,
		QuoteCoinAmount:   quoteAmount,
		QuoteCoinDecimals: quoteDecimals,
		Sender:            baseEvent.TxFromAddress,
		Receiver:          baseEvent.TxFromAddress,
	}}, nil
}
