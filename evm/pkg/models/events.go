package models

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type EventType string

const (
	EventTypeSwap               EventType = "swap"
	EventTypeTransfer           EventType = "transfer"
	EventTypeDefi               EventType = "defi"
	EventTypePolymarketOrder    EventType = "polymarket_order"
	EventTypePolymarketCTF      EventType = "polymarket_ctf_general"
	EventTypePolymarketCTFOrder EventType = "polymarket_ctf_order"
	EventTypeError              EventType = "error"
)

type Event interface {
	EventType() EventType
}

type BaseEvent struct {
	BlockTime       time.Time      `json:"block_time"`
	BlockNumber     *big.Int       `json:"block_number"`
	BlockHash       common.Hash    `json:"block_hash"`
	TxIdx           uint32         `json:"tx_idx"`
	TxHash          common.Hash    `json:"tx_hash"`
	GasUsed         uint64         `json:"gas_used"`
	Value           *big.Int       `json:"value"`
	TxFromAddress   common.Address `json:"tx_from_address"`
	TxToAddress     common.Address `json:"tx_to_address"`
	Source          string         `json:"source"`
	InstructionHash string         `json:"instruction_hash"`
}

type SwapEvent struct {
	BaseEvent         BaseEvent
	BaseCoin          common.Address `json:"base_coin"`
	QuoteCoin         common.Address `json:"quote_coin"`
	BaseCoinAmount    *big.Int       `json:"base_coin_amount"`
	BaseCoinDecimals  *uint8         `json:"base_coin_decimals"`
	QuoteCoinAmount   *big.Int       `json:"quote_coin_amount"`
	QuoteCoinDecimals *uint8         `json:"quote_coin_decimals"`
	Sender            common.Address `json:"sender"`
	Receiver          common.Address `json:"receiver"`
}

func (e *SwapEvent) EventType() EventType {
	return EventTypeSwap
}

type TransferEvent struct {
	BaseEvent     BaseEvent
	FromAddress   common.Address `json:"from_address"`
	ToAddress     common.Address `json:"to_address"`
	TokenAddress  common.Address `json:"token_address"`
	TokenDecimals *uint8         `json:"token_decimals"`
	Amount        *big.Int       `json:"amount"`
}

func (e *TransferEvent) EventType() EventType {
	return EventTypeTransfer
}

type DefiEvent struct {
	BaseEvent BaseEvent

	Inputs  []TokenTransferDescription `json:"inputs"`
	Outputs []TokenTransferDescription `json:"outputs"`
}

type TokenTransferDescription struct {
	TokenAddress common.Address `json:"token_address"`
	Amount       *big.Int       `json:"amount"`
	Decimals     *uint8         `json:"decimals"`
}

func (e *DefiEvent) EventType() EventType {
	return EventTypeDefi
}

type PolymarketOrderEvent struct {
	BaseEvent BaseEvent

	OrderHash         common.Hash    `json:"order_hash"`
	Maker             common.Address `json:"maker"`
	Taker             common.Address `json:"taker"`
	MakerAssetID      *big.Int       `json:"maker_asset_id"`
	TakerAssetID      *big.Int       `json:"taker_asset_id"`
	MakerAmountFilled *big.Int       `json:"maker_amount_filled"`
	TakerAmountFilled *big.Int       `json:"taker_amount_filled"`
	Fee               *big.Int       `json:"fee"`
	LogIndex          uint           `json:"log_index"`
}

func (e *PolymarketOrderEvent) EventType() EventType {
	return EventTypePolymarketOrder
}

type ErrorEvent struct {
	Chain       string
	BlockNumber *big.Int
	TxHash      common.Hash
	Error       string
}

func (e *ErrorEvent) EventType() EventType {
	return EventTypeError
}

type PolymarketOrderEventNew struct {
	BlockTime   time.Time   `json:"block_time"`
	BlockNumber *big.Int    `json:"block_number"`
	BlockHash   common.Hash `json:"block_hash"`
	TxIdx       uint32      `json:"tx_idx"`
	TxHash      common.Hash `json:"tx_hash"`
	LogIndex    uint32      `json:"log_index"`
	SubIndex    uint32      `json:"sub_index"`

	UserAddress   common.Address                `json:"user_address"`
	SourceAddress *common.Address               `json:"source_address"`
	Source        PolymarketOrderEventNewSource `json:"source"`

	ConditionID     common.Hash `json:"condition_id"`
	TokenID         *big.Int    `json:"token_id"`
	TokenAmountDiff *big.Int    `json:"token_amount_diff"`
	UsdcAmountDiff  *big.Int    `json:"usdc_amount_diff"`

	IsTaker bool     `json:"is_taker"`
	Fee     *big.Int `json:"fee"`
}

func (e *PolymarketOrderEventNew) EventType() EventType {
	return EventTypePolymarketCTFOrder
}

func (e *PolymarketOrderEventNew) Copy(subIndex uint32) *PolymarketOrderEventNew {
	return &PolymarketOrderEventNew{
		BlockTime:   e.BlockTime,
		BlockNumber: e.BlockNumber,
		BlockHash:   e.BlockHash,
		TxIdx:       e.TxIdx,
		TxHash:      e.TxHash,
		LogIndex:    e.LogIndex,
		SubIndex:    subIndex,

		ConditionID:     e.ConditionID,
		UserAddress:     e.UserAddress,
		SourceAddress:   e.SourceAddress,
		Source:          e.Source,
		TokenID:         e.TokenID,
		TokenAmountDiff: e.TokenAmountDiff,
		UsdcAmountDiff:  e.UsdcAmountDiff,

		IsTaker: e.IsTaker,
		Fee:     e.Fee,
	}
}

func (e *PolymarketOrderEventNew) Validate() error {
	if e.BlockNumber == nil {
		return fmt.Errorf("block number is required")
	}
	if e.BlockHash == (common.Hash{}) {
		return fmt.Errorf("block hash is required")
	}

	if e.SourceAddress != nil {
		if !common.IsHexAddress(e.SourceAddress.Hex()) {
			return fmt.Errorf("source address is invalid")
		}
	}

	if e.Source == "" {
		return fmt.Errorf("source is required")
	}

	if e.TokenID == nil {
		return fmt.Errorf("token id is required")
	}
	if e.TokenAmountDiff == nil {
		return fmt.Errorf("token amount diff is required")
	}
	if e.UsdcAmountDiff == nil {
		return fmt.Errorf("usdc amount diff is required")
	}

	return nil
}

type PolymarketOrderEventNewSource string

const (
	PolymarketOrderEventNewSourceOrderFilled              PolymarketOrderEventNewSource = "order_filled"
	PolymarketOrderEventNewSourcePositionSplit            PolymarketOrderEventNewSource = "position_split"
	PolymarketOrderEventNewSourcePositionMerge            PolymarketOrderEventNewSource = "position_merge"
	PolymarketOrderEventNewSourcePositionRedeem           PolymarketOrderEventNewSource = "position_redeem"
	PolymarketOrderEventNewSourcePositionTransfer         PolymarketOrderEventNewSource = "position_transfer"
	PolymarketOrderEventNewSourcePositionTransferBatch    PolymarketOrderEventNewSource = "position_transfer_batch"
	PolymarketOrderEventNewSourceNegRiskPositionConverted PolymarketOrderEventNewSource = "convert"
)

type ReorgEvent struct {
	BlockNumber *big.Int `json:"block_number"`
	Chain       string   `json:"chain"`
}
