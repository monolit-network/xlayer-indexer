package clickhouse

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	evmmodels "github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/ethereum/go-ethereum/common"
)

type evmSwapEventRow struct {
	Chain           string
	BlockTime       time.Time
	BlockNumber     *big.Int
	BlockHash       string
	TxIdx           uint32
	TxHash          string
	GasUsed         uint64
	Value           *big.Int
	TxFromAddress   string
	TxToAddress     string
	Source          string
	InstructionHash string

	BaseCoin          string
	QuoteCoin         string
	BaseCoinAmount    *big.Int
	BaseCoinDecimals  *uint8
	QuoteCoinAmount   *big.Int
	QuoteCoinDecimals *uint8
	Sender            string
	Receiver          string
}

type evmTransferEventRow struct {
	Chain       string
	BlockTime   time.Time
	BlockNumber *big.Int
	BlockHash   string
	TxIdx       uint32
	TxHash      string
	GasUsed     uint64

	Sender        string
	Receiver      string
	TokenAddress  string
	TokenDecimals *uint8
	Amount        *big.Int
	Source        string
}

func normalizeHex(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s)
}

func normalizeAddress(addr common.Address) string {
	return strings.ToLower(addr.Hex())
}

func normalizeHash(hash common.Hash) string {
	return strings.ToLower(hash.Hex())
}

func (evmSwapEventRow) FromSwapEvent(chain evmmodels.Chain, e *evmmodels.SwapEvent) evmSwapEventRow {
	return evmSwapEventRow{
		Chain:             string(chain),
		BlockTime:         e.BaseEvent.BlockTime,
		BlockNumber:       e.BaseEvent.BlockNumber,
		BlockHash:         normalizeHash(e.BaseEvent.BlockHash),
		TxIdx:             e.BaseEvent.TxIdx,
		TxHash:            normalizeHash(e.BaseEvent.TxHash),
		GasUsed:           e.BaseEvent.GasUsed,
		Value:             e.BaseEvent.Value,
		TxFromAddress:     normalizeAddress(e.BaseEvent.TxFromAddress),
		TxToAddress:       normalizeAddress(e.BaseEvent.TxToAddress),
		Source:            e.BaseEvent.Source,
		InstructionHash:   e.BaseEvent.InstructionHash,
		BaseCoin:          normalizeAddress(e.BaseCoin),
		QuoteCoin:         normalizeAddress(e.QuoteCoin),
		BaseCoinAmount:    e.BaseCoinAmount,
		BaseCoinDecimals:  e.BaseCoinDecimals,
		QuoteCoinAmount:   e.QuoteCoinAmount,
		QuoteCoinDecimals: e.QuoteCoinDecimals,
		Sender:            normalizeAddress(e.Sender),
		Receiver:          normalizeAddress(e.Receiver),
	}
}

func (evmTransferEventRow) FromTransferEvent(chain evmmodels.Chain, e *evmmodels.TransferEvent) evmTransferEventRow {
	return evmTransferEventRow{
		Chain:         string(chain),
		BlockTime:     e.BaseEvent.BlockTime,
		BlockNumber:   e.BaseEvent.BlockNumber,
		BlockHash:     normalizeHash(e.BaseEvent.BlockHash),
		TxIdx:         e.BaseEvent.TxIdx,
		TxHash:        normalizeHash(e.BaseEvent.TxHash),
		GasUsed:       e.BaseEvent.GasUsed,
		Sender:        normalizeAddress(e.FromAddress),
		Receiver:      normalizeAddress(e.ToAddress),
		TokenAddress:  normalizeAddress(e.TokenAddress),
		TokenDecimals: e.TokenDecimals,
		Amount:        e.Amount,
		Source:        e.BaseEvent.Source,
	}
}

func trimFixedString(s string) string {
	i := len(s)
	for i > 0 {
		ch := s[i-1]
		if ch != '\x00' && ch != ' ' {
			break
		}
		i--
	}
	return s[:i]
}

func (r evmSwapEventRow) ToEvent() evmmodels.Event {
	blockHash := trimFixedString(r.BlockHash)
	txHash := trimFixedString(r.TxHash)
	txFrom := trimFixedString(r.TxFromAddress)
	txTo := trimFixedString(r.TxToAddress)
	baseCoin := trimFixedString(r.BaseCoin)
	quoteCoin := trimFixedString(r.QuoteCoin)
	sender := trimFixedString(r.Sender)
	receiver := trimFixedString(r.Receiver)

	bn := r.BlockNumber
	if bn == nil {
		bn = big.NewInt(0)
	}
	val := r.Value
	if val == nil {
		val = big.NewInt(0)
	}
	bAmt := r.BaseCoinAmount
	if bAmt == nil {
		bAmt = big.NewInt(0)
	}
	qAmt := r.QuoteCoinAmount
	if qAmt == nil {
		qAmt = big.NewInt(0)
	}

	return &evmmodels.SwapEvent{
		BaseEvent: evmmodels.BaseEvent{
			BlockTime:       r.BlockTime,
			BlockNumber:     bn,
			BlockHash:       common.HexToHash(blockHash),
			TxIdx:           r.TxIdx,
			TxHash:          common.HexToHash(txHash),
			GasUsed:         r.GasUsed,
			Value:           val,
			TxFromAddress:   common.HexToAddress(txFrom),
			TxToAddress:     common.HexToAddress(txTo),
			Source:          r.Source,
			InstructionHash: r.InstructionHash,
		},
		BaseCoin:          common.HexToAddress(baseCoin),
		QuoteCoin:         common.HexToAddress(quoteCoin),
		BaseCoinAmount:    bAmt,
		BaseCoinDecimals:  r.BaseCoinDecimals,
		QuoteCoinAmount:   qAmt,
		QuoteCoinDecimals: r.QuoteCoinDecimals,
		Sender:            common.HexToAddress(sender),
		Receiver:          common.HexToAddress(receiver),
	}
}

func (r evmTransferEventRow) ToEvent() evmmodels.Event {
	blockHash := trimFixedString(r.BlockHash)
	txHash := trimFixedString(r.TxHash)
	sender := trimFixedString(r.Sender)
	receiver := trimFixedString(r.Receiver)
	tokenAddr := trimFixedString(r.TokenAddress)

	bn := r.BlockNumber
	if bn == nil {
		bn = big.NewInt(0)
	}
	amt := r.Amount
	if amt == nil {
		amt = big.NewInt(0)
	}

	return &evmmodels.TransferEvent{
		BaseEvent: evmmodels.BaseEvent{
			BlockTime:       r.BlockTime,
			BlockNumber:     bn,
			BlockHash:       common.HexToHash(blockHash),
			TxIdx:           r.TxIdx,
			TxHash:          common.HexToHash(txHash),
			GasUsed:         r.GasUsed,
			Value:           nil,
			TxFromAddress:   common.HexToAddress(sender),
			TxToAddress:     common.HexToAddress(receiver),
			Source:          "",
			InstructionHash: "",
		},
		FromAddress:   common.HexToAddress(sender),
		ToAddress:     common.HexToAddress(receiver),
		TokenAddress:  common.HexToAddress(tokenAddr),
		TokenDecimals: r.TokenDecimals,
		Amount:        amt,
	}
}

func (r evmDefiEventRow) ToEvent() evmmodels.Event {
	blockHash := trimFixedString(r.BlockHash)
	txHash := trimFixedString(r.TxHash)
	txFrom := trimFixedString(r.TxFromAddress)
	txTo := trimFixedString(r.TxToAddress)

	bn := r.BlockNumber
	if bn == nil {
		bn = big.NewInt(0)
	}
	val := r.Value
	if val == nil {
		val = big.NewInt(0)
	}

	makeTokens := func(addresses []string, amounts []*big.Int, decimals []*uint8) []evmmodels.TokenTransferDescription {
		minLen := len(addresses)
		if len(amounts) < minLen {
			minLen = len(amounts)
		}
		if len(decimals) < minLen {
			minLen = len(decimals)
		}
		if minLen == 0 {
			return nil
		}

		results := make([]evmmodels.TokenTransferDescription, 0, minLen)
		for i := 0; i < minLen; i++ {
			addr := trimFixedString(addresses[i])
			if addr == "" {
				continue
			}

			amount := big.NewInt(0)
			if amounts[i] != nil {
				amount = new(big.Int).Set(amounts[i])
			}

			results = append(results, evmmodels.TokenTransferDescription{
				TokenAddress: common.HexToAddress(addr),
				Amount:       amount,
				Decimals:     decimals[i],
			})
		}
		return results
	}

	return &evmmodels.DefiEvent{
		BaseEvent: evmmodels.BaseEvent{
			BlockTime:       r.BlockTime,
			BlockNumber:     bn,
			BlockHash:       common.HexToHash(blockHash),
			TxIdx:           r.TxIdx,
			TxHash:          common.HexToHash(txHash),
			GasUsed:         r.GasUsed,
			Value:           val,
			TxFromAddress:   common.HexToAddress(txFrom),
			TxToAddress:     common.HexToAddress(txTo),
			Source:          r.Source,
			InstructionHash: r.InstructionHash,
		},
		Inputs:  makeTokens(r.InputAddresses, r.InputAmounts, r.InputDecimals),
		Outputs: makeTokens(r.OutputAddresses, r.OutputAmounts, r.OutputDecimals),
	}
}

func (r evmErrorEventRow) ToEvent() evmmodels.Event {
	bn := r.BlockNumber
	if bn == nil {
		bn = big.NewInt(0)
	}
	return &evmmodels.ErrorEvent{
		Chain:       r.Chain,
		BlockNumber: bn,
		TxHash:      common.HexToHash(trimFixedString(r.TxHash)),
		Error:       r.Error,
	}
}

type evmDefiEventRow struct {
	Chain           string
	BlockTime       time.Time
	BlockNumber     *big.Int
	BlockHash       string
	TxIdx           uint32
	TxHash          string
	GasUsed         uint64
	Value           *big.Int
	TxFromAddress   string
	TxToAddress     string
	Source          string
	InstructionHash string

	InputAddresses []string
	InputAmounts   []*big.Int
	InputDecimals  []*uint8

	OutputAddresses []string
	OutputAmounts   []*big.Int
	OutputDecimals  []*uint8
}

func (evmDefiEventRow) FromDefiEvent(chain evmmodels.Chain, e *evmmodels.DefiEvent) evmDefiEventRow {
	getTokens := func(descs []evmmodels.TokenTransferDescription) ([]string, []*big.Int, []*uint8) {
		if len(descs) == 0 {
			return nil, nil, nil
		}
		addresses := make([]string, len(descs))
		amounts := make([]*big.Int, len(descs))
		decimals := make([]*uint8, len(descs))

		for i, desc := range descs {
			addresses[i] = normalizeAddress(desc.TokenAddress)
			if desc.Amount != nil {
				amounts[i] = new(big.Int).Set(desc.Amount)
			} else {
				amounts[i] = big.NewInt(0)
			}
			decimals[i] = desc.Decimals
		}

		return addresses, amounts, decimals
	}

	inAddr, inAmt, inDec := getTokens(e.Inputs)
	outAddr, outAmt, outDec := getTokens(e.Outputs)

	return evmDefiEventRow{
		Chain:           string(chain),
		BlockTime:       e.BaseEvent.BlockTime,
		BlockNumber:     e.BaseEvent.BlockNumber,
		BlockHash:       normalizeHash(e.BaseEvent.BlockHash),
		TxIdx:           e.BaseEvent.TxIdx,
		TxHash:          normalizeHash(e.BaseEvent.TxHash),
		GasUsed:         e.BaseEvent.GasUsed,
		Value:           e.BaseEvent.Value,
		TxFromAddress:   normalizeAddress(e.BaseEvent.TxFromAddress),
		TxToAddress:     normalizeAddress(e.BaseEvent.TxToAddress),
		Source:          e.BaseEvent.Source,
		InstructionHash: e.BaseEvent.InstructionHash,
		InputAddresses:  inAddr,
		InputAmounts:    inAmt,
		InputDecimals:   inDec,
		OutputAddresses: outAddr,
		OutputAmounts:   outAmt,
		OutputDecimals:  outDec,
	}
}

type evmErrorEventRow struct {
	Chain       string
	BlockNumber *big.Int
	TxHash      string
	Error       string
}

func (evmErrorEventRow) FromErrorEvent(chain evmmodels.Chain, e *evmmodels.ErrorEvent) evmErrorEventRow {
	return evmErrorEventRow{
		Chain:       string(chain),
		BlockNumber: e.BlockNumber,
		TxHash:      normalizeHash(e.TxHash),
		Error:       e.Error,
	}
}

type evmPolymarketOrderEventRow struct {
	BlockTime       time.Time
	BlockNumber     *big.Int
	BlockHash       string
	TxIdx           uint32
	TxHash          string
	TxFromAddress   string
	TxToAddress     string
	InstructionHash string

	ID                string
	TransactionHash   string
	OrderHash         string
	Maker             string
	Taker             string
	MakerAssetID      string
	TakerAssetID      string
	MakerAmountFilled *big.Int
	TakerAmountFilled *big.Int
	Fee               *big.Int
	IsDeleted         uint8
}

func (evmPolymarketOrderEventRow) FromPolymarketOrderEvent(e *evmmodels.PolymarketOrderEvent) evmPolymarketOrderEventRow {
	txHash := normalizeHash(e.BaseEvent.TxHash)
	orderHash := normalizeHash(e.OrderHash)
	id := fmt.Sprintf("%s_%s", txHash, orderHash)

	return evmPolymarketOrderEventRow{
		BlockTime:       e.BaseEvent.BlockTime,
		BlockNumber:     e.BaseEvent.BlockNumber,
		BlockHash:       normalizeHash(e.BaseEvent.BlockHash),
		TxIdx:           e.BaseEvent.TxIdx,
		TxHash:          txHash,
		TxFromAddress:   normalizeAddress(e.BaseEvent.TxFromAddress),
		TxToAddress:     normalizeAddress(e.BaseEvent.TxToAddress),
		InstructionHash: e.BaseEvent.InstructionHash,

		ID:                id,
		TransactionHash:   txHash,
		OrderHash:         orderHash,
		Maker:             normalizeAddress(e.Maker),
		Taker:             normalizeAddress(e.Taker),
		MakerAssetID:      e.MakerAssetID.String(),
		TakerAssetID:      e.TakerAssetID.String(),
		MakerAmountFilled: e.MakerAmountFilled,
		TakerAmountFilled: e.TakerAmountFilled,
		Fee:               e.Fee,
		IsDeleted:         0,
	}
}
