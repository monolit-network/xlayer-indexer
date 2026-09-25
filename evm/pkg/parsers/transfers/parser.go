package transfers

import (
	"fmt"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type ERC20OrNativeTransfersIdentifier struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewERC20OrNativeTransfersIdentifier(evmClient *evmclient.Client, logger *zap.Logger) *ERC20OrNativeTransfersIdentifier {
	return &ERC20OrNativeTransfersIdentifier{evmClient: evmClient, logger: logger}
}

func (i *ERC20OrNativeTransfersIdentifier) Source() string {
	return "erc20_or_native_transfers"
}

func (i *ERC20OrNativeTransfersIdentifier) ShouldStop() bool {
	return true
}

func (i *ERC20OrNativeTransfersIdentifier) CheckTx(req parsers.ParseTxRequest) (bool, error) {
	if req.Tx.To() == nil {
		return false, nil
	}

	if req.Receipt == nil {
		return false, fmt.Errorf("receipt is required")
	}

	// case 1: erc20 transfer
	if len(req.Receipt.Logs) == 1 {
		_, err := models.TryParseERC20TransferLog(req.Receipt.Logs[0])
		if err != nil {
			return false, nil
		}
		if req.Tx.Value() != nil && req.Tx.Value().Cmp(common.Big0) != 0 {
			// erc20 transfer + native transfer, for now should go to defi event
			return false, nil
		}
		if req.Tx.To() == nil || req.Tx.To().Hex() != req.Receipt.Logs[0].Address.Hex() {
			return false, nil
		}

		return true, nil
	}

	hasLogs := len(req.Receipt.Logs) > 0
	hasData := len(req.Tx.Data()) > 0

	// case 2: native transfer, mb to smart contract
	if !hasLogs && !hasData {
		return true, nil
	}

	// case 3: native transfer with data. if to contract -> defi, if to EOA -> transfer
	if !hasLogs && hasData {
		isContract, err := i.evmClient.GetIsContract(*req.Tx.To())
		if err != nil {
			return true, nil
		}
		return !isContract, nil
	}

	return false, nil
}

func (i *ERC20OrNativeTransfersIdentifier) GetParserName() string {
	return "erc20_or_native_transfers_parser"
}

func NewERC20OrNativeTransfersParser(evmClient *evmclient.Client, logger *zap.Logger) *ERC20OrNativeTransfersParser {
	p := &ERC20OrNativeTransfersParser{evmClient: evmClient}
	p.logger = logger.Named(p.Name())
	return p
}

type ERC20OrNativeTransfersParser struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func (p ERC20OrNativeTransfersParser) Name() string {
	return "erc20_or_native_transfers_parser"
}

func (p *ERC20OrNativeTransfersParser) ProducedEventsType() models.EventType {
	return models.EventTypeTransfer
}

func (p *ERC20OrNativeTransfersParser) RequiredData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesBlockHeader
}

func (p *ERC20OrNativeTransfersParser) OptionalData() evmclient.RequiredDataTypes {
	return 0
}

func (p *ERC20OrNativeTransfersParser) ParseTx(req parsers.ParseTxRequest) ([]models.Event, error) {
	if err := parsers.ValidateParseTxRequest(req, p.RequiredData()); err != nil {
		return nil, err
	}
	baseEvent, err := parsers.MakeBaseEvent(req, p.evmClient)
	if err != nil {
		return nil, err
	}

	// case 1: native transfer
	if len(req.Receipt.Logs) == 0 {
		if req.Tx.Value() == nil || req.Tx.Value().Cmp(common.Big0) == 0 {
			return nil, nil
		}
		dec := models.ETHTokensDecimals
		transfer := models.TransferEvent{
			BaseEvent:     baseEvent,
			FromAddress:   baseEvent.TxFromAddress,
			ToAddress:     baseEvent.TxToAddress,
			TokenAddress:  models.ETHTokenAddress,
			TokenDecimals: &dec,
			Amount:        req.Tx.Value(),
		}
		return []models.Event{&transfer}, nil
	}

	transfers, err := models.FindAllTransfers(req.Receipt, req.Tx, req.Traces, p.evmClient.ChainID)
	if err != nil {
		return nil, err
	}

	if len(transfers) != 1 {
		p.logger.Debug("found to many/less transfers, must be exactly 1", zap.Any("transfers", transfers))
		return nil, fmt.Errorf("found to many/less transfers, must be exactly 1")
	}

	transfer := transfers[0]

	var dec *uint8
	decimals, err := p.evmClient.GetDecimals(transfer.TokenAddress)
	if err != nil {
	} else {
		dec = &decimals
	}

	transferEvent := models.TransferEvent{
		BaseEvent:     baseEvent,
		FromAddress:   transfer.FromAddress,
		ToAddress:     transfer.ToAddress,
		TokenAddress:  transfer.TokenAddress,
		TokenDecimals: dec,
		Amount:        transfer.Value,
	}
	return []models.Event{&transferEvent}, nil
}
