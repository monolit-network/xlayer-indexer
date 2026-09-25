package parsers

import (
	"fmt"
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/ethereum/go-ethereum/core/types"
)

type Identifier interface {
	CheckTx(tx ParseTxRequest) (bool, error)
	GetParserName() string
	Source() string
	ShouldStop() bool
}

type Parser interface {
	Name() string
	ProducedEventsType() models.EventType
	RequiredData() evmclient.RequiredDataTypes
	OptionalData() evmclient.RequiredDataTypes
	ParseTx(req ParseTxRequest) ([]models.Event, error)
}

type ParseTxRequest struct {
	Tx          *types.Transaction
	Receipt     *types.Receipt
	BlockHeader *types.Header
	Traces      *models.FullTraceResult
	Source      string
}

func (r *ParseTxRequest) Copy() ParseTxRequest {
	return ParseTxRequest{
		Tx:          r.Tx,
		Receipt:     r.Receipt,
		BlockHeader: r.BlockHeader,
		Traces:      r.Traces,
		Source:      r.Source,
	}
}

// ValidateParseTxRequest: should be called in each parser
func ValidateParseTxRequest(request ParseTxRequest, requiredData evmclient.RequiredDataTypes) error {
	if request.Source == "" {
		return fmt.Errorf("source is required")
	}

	if requiredData.Has(evmclient.RequiredDataTypesTx) {
		if request.Tx == nil {
			return fmt.Errorf("tx is required")
		}
	}
	if requiredData.Has(evmclient.RequiredDataTypesReceipt) {
		if request.Receipt == nil {
			return fmt.Errorf("receipt is required")
		}
	}
	if requiredData.Has(evmclient.RequiredDataTypesBlockHeader) {
		if request.BlockHeader == nil {
			return fmt.Errorf("block header is required")
		}
	}
	if requiredData.Has(evmclient.RequiredDataTypesTraces) {
		if request.Traces == nil {
			return fmt.Errorf("traces are required")
		}
	}
	return nil
}

func MakeBaseEvent(request ParseTxRequest, evmClient *evmclient.Client) (models.BaseEvent, error) {
	event := models.BaseEvent{}
	event.BlockTime = time.Unix(int64(request.BlockHeader.Time), 0)
	event.Value = request.Tx.Value()
	event.BlockNumber = request.BlockHeader.Number
	event.BlockHash = request.BlockHeader.Hash()
	event.TxIdx = uint32(request.Receipt.TransactionIndex)
	event.TxHash = request.Tx.Hash()
	event.GasUsed = request.Receipt.GasUsed
	event.InstructionHash = models.SelectorFromTx(request.Tx)
	txFromAddress, err := models.GetSenderAddress(request.Tx, evmClient.ChainID)
	if err != nil {
		return event, fmt.Errorf("error getting sender address: %w", err)
	}
	event.TxFromAddress = txFromAddress

	if request.Tx.To() != nil {
		event.TxToAddress = *request.Tx.To()
	} else {
		return event, fmt.Errorf("to address is nil")
	}
	event.Source = request.Source
	return event, nil
}
