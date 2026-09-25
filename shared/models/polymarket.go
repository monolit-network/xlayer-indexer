package models

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

// PolymarketMarketEventType represents the type of market event
type PolymarketMarketEventType string

const (
	PolymarketEventAdapterInitialized PolymarketMarketEventType = "ADAPTER_INITIALIZED"
	PolymarketEventAdapterReset       PolymarketMarketEventType = "ADAPTER_RESET"
	PolymarketEventAdapterResolved    PolymarketMarketEventType = "ADAPTER_RESOLVED"
	PolymarketEventAdapterPaused      PolymarketMarketEventType = "ADAPTER_PAUSED"
	PolymarketEventAdapterUnpaused    PolymarketMarketEventType = "ADAPTER_UNPAUSED"
	PolymarketEventAdapterFlagged     PolymarketMarketEventType = "ADAPTER_FLAGGED"
	PolymarketEventUmaRequested       PolymarketMarketEventType = "UMA_REQUESTED"
	PolymarketEventUmaProposed        PolymarketMarketEventType = "UMA_PROPOSED"
	PolymarketEventUmaDisputed        PolymarketMarketEventType = "UMA_DISPUTED"
	PolymarketEventUmaSettle          PolymarketMarketEventType = "UMA_SETTLE"
	PolymarketEventUmaSettleInvalid   PolymarketMarketEventType = "UMA_SETTLE_INVALID"
	PolymarketEventCondPrepared       PolymarketMarketEventType = "COND_PREPARED"
	PolymarketEventCondResolved       PolymarketMarketEventType = "COND_RESOLVED"
)

// PolymarketMarketEventData is an interface for event-specific data
type PolymarketMarketEventData interface {
	PolymarketMarketEventType() PolymarketMarketEventType
}

func UnmarshalPolymarketMarketEventData(eventType PolymarketMarketEventType, eventDataRaw []byte) (PolymarketMarketEventData, error) {
	switch eventType {
	case PolymarketEventAdapterInitialized:
		var eventData PolymarketMarketEventDataAdapterInitialized
		if err := json.Unmarshal(eventDataRaw, &eventData); err != nil {
			return nil, err
		}
		return &eventData, nil
	case PolymarketEventAdapterReset:
		var eventData PolymarketMarketEventDataAdapterReset
		if err := json.Unmarshal(eventDataRaw, &eventData); err != nil {
			return nil, err
		}
		return &eventData, nil
	case PolymarketEventAdapterResolved:
		var eventData PolymarketMarketEventDataAdapterResolved
		if err := json.Unmarshal(eventDataRaw, &eventData); err != nil {
			return nil, err
		}
		return &eventData, nil
	case PolymarketEventAdapterPaused:
		return nil, nil
	case PolymarketEventAdapterUnpaused:
		return nil, nil
	case PolymarketEventAdapterFlagged:
		return nil, nil
	case PolymarketEventUmaRequested:
		var eventData PolymarketMarketEventDataUmaRequested
		if err := json.Unmarshal(eventDataRaw, &eventData); err != nil {
			return nil, err
		}
		return &eventData, nil
	case PolymarketEventUmaProposed:
		var eventData PolymarketMarketEventDataUmaProposed
		if err := json.Unmarshal(eventDataRaw, &eventData); err != nil {
			return nil, err
		}
		return &eventData, nil
	case PolymarketEventUmaDisputed:
		var eventData PolymarketMarketEventDataUmaDisputed
		if err := json.Unmarshal(eventDataRaw, &eventData); err != nil {
			return nil, err
		}
		return &eventData, nil
	case PolymarketEventUmaSettle, PolymarketEventUmaSettleInvalid:
		var eventData PolymarketMarketEventDataUmaSettle
		if err := json.Unmarshal(eventDataRaw, &eventData); err != nil {
			return nil, err
		}
		return &eventData, nil
	case PolymarketEventCondPrepared:
		var eventData PolymarketMarketEventDataCondPrepared
		if err := json.Unmarshal(eventDataRaw, &eventData); err != nil {
			return nil, err
		}
		return &eventData, nil
	case PolymarketEventCondResolved:
		var eventData PolymarketMarketEventDataCondResolved
		if err := json.Unmarshal(eventDataRaw, &eventData); err != nil {
			return nil, err
		}
		return &eventData, nil
	}
	return nil, fmt.Errorf("unsupported event type: %s", eventType)
}

// PolymarketMarketEvent represents an event in the market lifecycle
type PolymarketMarketEvent struct {
	QuestionID     string                    `json:"question_id"`
	EventType      PolymarketMarketEventType `json:"event_type"`
	BlockNumber    uint64                    `json:"block_number"`
	BlockHash      string                    `json:"block_hash"`
	LogIndex       uint32                    `json:"log_index"`
	TxIndex        uint32                    `json:"tx_index"`
	TxHash         string                    `json:"tx_hash"`
	BlockTimestamp time.Time                 `json:"block_timestamp"`
	EventData      PolymarketMarketEventData `json:"event_data,omitempty"`
}

func (e *PolymarketMarketEvent) UnmarshalJSON(data []byte) error {
	type polymarketMarketEventJSON struct {
		QuestionID     string                    `json:"question_id"`
		EventType      PolymarketMarketEventType `json:"event_type"`
		BlockNumber    uint64                    `json:"block_number"`
		BlockHash      string                    `json:"block_hash"`
		LogIndex       uint32                    `json:"log_index"`
		TxIndex        uint32                    `json:"tx_index"`
		TxHash         string                    `json:"tx_hash"`
		BlockTimestamp time.Time                 `json:"block_timestamp"`
		EventData      json.RawMessage           `json:"event_data,omitempty"`
	}

	var raw polymarketMarketEventJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	e.QuestionID = raw.QuestionID
	e.EventType = raw.EventType
	e.BlockNumber = raw.BlockNumber
	e.BlockHash = raw.BlockHash
	e.LogIndex = raw.LogIndex
	e.TxIndex = raw.TxIndex
	e.TxHash = raw.TxHash
	e.BlockTimestamp = raw.BlockTimestamp

	if len(raw.EventData) == 0 || string(raw.EventData) == "null" {
		e.EventData = nil
		return nil
	}

	eventData, err := UnmarshalPolymarketMarketEventData(raw.EventType, raw.EventData)
	if err != nil {
		return err
	}
	e.EventData = eventData
	return nil
}

// Event data structs for each event type

// PolymarketMarketEventDataAdapterInitialized contains data for ADAPTER_INITIALIZED event
type PolymarketMarketEventDataAdapterInitialized struct {
	RequestTimestamp    uint64               `json:"request_timestamp"`
	Creator             string               `json:"creator"`
	RewardToken         string               `json:"reward_token"`
	Reward              string               `json:"reward"`
	ProposalBond        string               `json:"proposal_bond"`
	ParsedAncillaryData *ParsedAncillaryData `json:"parsed_ancillary_data,omitempty"` // Field only for notifier don't use in db
}

func (*PolymarketMarketEventDataAdapterInitialized) PolymarketMarketEventType() PolymarketMarketEventType {
	return PolymarketEventAdapterInitialized
}

// PolymarketMarketEventDataAdapterResolved contains data for ADAPTER_RESOLVED event
type PolymarketMarketEventDataAdapterResolved struct {
	SettledPrice string   `json:"settled_price"`
	Payouts      []string `json:"payouts"`
}

func (*PolymarketMarketEventDataAdapterResolved) PolymarketMarketEventType() PolymarketMarketEventType {
	return PolymarketEventAdapterResolved
}

// PolymarketMarketEventDataUmaProposed contains data for UMA_PROPOSED event
type PolymarketMarketEventDataUmaProposed struct {
	Requester           string `json:"requester"`
	Proposer            string `json:"proposer"`
	Identifier          string `json:"identifier"`
	Timestamp           uint64 `json:"timestamp"`
	ProposedPrice       string `json:"proposed_price"`
	ExpirationTimestamp uint64 `json:"expiration_timestamp"`
	Currency            string `json:"currency"`
}

func (*PolymarketMarketEventDataUmaProposed) PolymarketMarketEventType() PolymarketMarketEventType {
	return PolymarketEventUmaProposed
}

// PolymarketMarketEventDataUmaDisputed contains data for UMA_DISPUTED event
type PolymarketMarketEventDataUmaDisputed struct {
	Requester     string `json:"requester"`
	Proposer      string `json:"proposer"`
	Disputer      string `json:"disputer"`
	Identifier    string `json:"identifier"`
	Timestamp     uint64 `json:"timestamp"`
	ProposedPrice string `json:"proposed_price"`
}

func (*PolymarketMarketEventDataUmaDisputed) PolymarketMarketEventType() PolymarketMarketEventType {
	return PolymarketEventUmaDisputed
}

// PolymarketMarketEventDataUmaSettle contains data for UMA_SETTLE event
type PolymarketMarketEventDataUmaSettle struct {
	Requester  string `json:"requester"`
	Proposer   string `json:"proposer"`
	Disputer   string `json:"disputer"`
	Identifier string `json:"identifier"`
	Timestamp  uint64 `json:"timestamp"`
	Price      string `json:"price"`
	Payout     string `json:"payout"`
}

func (*PolymarketMarketEventDataUmaSettle) PolymarketMarketEventType() PolymarketMarketEventType {
	return PolymarketEventUmaSettle
}

// PolymarketMarketEventDataUmaRequested contains data for UMA_REQUESTED event
type PolymarketMarketEventDataUmaRequested struct {
	Requester  string `json:"requester"`
	Identifier string `json:"identifier"`
	Timestamp  uint64 `json:"timestamp"`
	Currency   string `json:"currency"`
	Reward     string `json:"reward"`
	FinalFee   string `json:"final_fee"`
}

func (*PolymarketMarketEventDataUmaRequested) PolymarketMarketEventType() PolymarketMarketEventType {
	return PolymarketEventUmaRequested
}

// PolymarketMarketEventDataCondPrepared contains data for COND_PREPARED event
type PolymarketMarketEventDataCondPrepared struct {
	ConditionID   string `json:"condition_id"`
	Oracle        string `json:"oracle"`
	OutcomesCount int64  `json:"outcomes_count"`
}

func (*PolymarketMarketEventDataCondPrepared) PolymarketMarketEventType() PolymarketMarketEventType {
	return PolymarketEventCondPrepared
}

// PolymarketMarketEventDataCondResolved contains data for COND_RESOLVED event
type PolymarketMarketEventDataCondResolved struct {
	ConditionID string   `json:"condition_id"`
	Oracle      string   `json:"oracle"`
	Numerators  []string `json:"numerators"`
}

func (*PolymarketMarketEventDataCondResolved) PolymarketMarketEventType() PolymarketMarketEventType {
	return PolymarketEventCondResolved
}

type PolymarketMarketEventDataAdapterReset struct {
	NewTimestamp uint64 `json:"new_timestamp"`
}

func (*PolymarketMarketEventDataAdapterReset) PolymarketMarketEventType() PolymarketMarketEventType {
	return PolymarketEventAdapterReset
}

type PolymarketMarketOnchain struct {
	ConditionID   common.Hash    `json:"condition_id"`
	QuestionID    common.Hash    `json:"question_id"`
	Oracle        common.Address `json:"oracle_id"`
	OutcomesCount int64          `json:"outcomes_count"`
}

type PolymarketResolutionWinners struct {
	Oracle      common.Address `json:"oracle"`
	ConditionID common.Hash    `json:"condition_id"`
	Numerators  []*big.Int     `json:"numerators"`
}

type PolymarketPositionSplitMerge struct {
	Stakeholder        common.Address `json:"stakeholder"`
	CollateralToken    common.Address `json:"collateral_token"`
	ParentCollectionId common.Hash    `json:"parent_collection_id"`
	ConditionID        common.Hash    `json:"condition_id"`
	Partition          []*big.Int     `json:"partition"`
	TokenIDs           []*big.Int     `json:"token_ids"`
	Amount             *big.Int       `json:"amount"`
	IsSplit            bool           `json:"is_split"`
}

type PolymarketOrderFilled struct {
	OrderHash         common.Hash    `json:"order_hash"`
	Maker             common.Address `json:"maker"`
	Taker             common.Address `json:"taker"`
	MakerAssetID      *big.Int       `json:"maker_asset_id"`
	TakerAssetID      *big.Int       `json:"taker_asset_id"`
	MakerAmountFilled *big.Int       `json:"maker_amount_filled"`
	TakerAmountFilled *big.Int       `json:"taker_amount_filled"`
	Fee               *big.Int       `json:"fee"`
}

type PolymarketTransfer struct {
	From    common.Address `json:"from"`
	To      common.Address `json:"to"`
	Amount  *big.Int       `json:"amount"`
	TokenID *big.Int       `json:"token_id"`
}

type PolymarketMarketNew struct {
	ConditionID         string      `json:"condition_id"`
	QuestionID          string      `json:"question_id"`
	Question            *string     `json:"question,omitempty"`
	Description         *string     `json:"description,omitempty"`
	AncillaryData       []byte      `json:"ancillary_data,omitempty"`
	Oracle              string      `json:"oracle"`
	PreparedAt          time.Time   `json:"prepared_at"`
	PreparedInBlock     int         `json:"prepared_in_block"`
	PreparedInBlockHash common.Hash `json:"prepared_in_block_hash"`
	PreparedInTxHash    common.Hash `json:"prepared_in_tx_hash"`
	TotalOutcomes       int         `json:"total_outcomes"`
	IsResolved          bool        `json:"is_resolved"`
	ResolvedAt          *time.Time  `json:"resolved_at"`
	ResolvedInBlock     int64       `json:"resolved_in_block"`
	ResolvedInBlockHash string      `json:"resolved_in_block_hash"`
	PayoutNumerators    []*big.Int  `json:"payout_numerators"`
}

type PolymarketMarketNewMinimal struct {
	ConditionID string `json:"condition_id"`
	QuestionID  string `json:"question_id"`
}

type PolymarketToken struct {
	TokenID             *big.Int `json:"token_id"`
	ConditionID         string   `json:"condition_id"`
	CollateralToken     string   `json:"collateral_token"`
	ParentCollectionID  string   `json:"parent_collection_id"`
	Partition           *big.Int `json:"partition"`
	IsResolved          bool     `json:"is_resolved"`
	ResolvedInBlock     *int64   `json:"resolved_in_block"`
	ResolvedInBlockHash *string  `json:"resolved_in_block_hash"`
	Numerator           *big.Int `json:"numerator"`
	Denominator         *big.Int `json:"denominator"`

	CreatedAt time.Time `json:"created_at"`
}

type PolymarketUmaCtfAdapterQuestionInitialized struct {
	QuestionID       common.Hash    `json:"question_id"`
	RequestTimestamp *big.Int       `json:"request_timestamp"`
	Creator          common.Address `json:"creator"`
	AncillaryData    []byte         `json:"ancillary_data"`
	RewardToken      common.Address `json:"reward_token"`
	Reward           *big.Int       `json:"reward"`
	ProposalBond     *big.Int       `json:"proposal_bond"`
}

type PolymarketUmaCtfAdapterQuestionReset struct {
	QuestionID common.Hash `json:"question_id"`
}

type PolymarketUmaCtfAdapterQuestionResolved struct {
	QuestionID   common.Hash `json:"question_id"`
	SettledPrice *big.Int    `json:"settled_price"`
	Payouts      []*big.Int  `json:"payouts"`
}

type PolymarketUmaCtfAdapterQuestionPaused struct {
	QuestionID common.Hash `json:"question_id"`
}

type PolymarketUmaCtfAdapterQuestionFlagged struct {
	QuestionID common.Hash `json:"question_id"`
}

type PolymarketFeeRefunded struct {
	To      common.Address `json:"to"`
	TokenID *big.Int       `json:"token_id"`
	Refund  *big.Int       `json:"refund"`
}

type PolymarketNegRiskPositionsConvertedRaw struct {
	LogAddress  common.Address
	Stakeholder common.Address
	MarketId    [32]byte
	IndexSet    *big.Int
	Amount      *big.Int
}

type PolymarketTokenIDWithConditionID struct {
	TokenID     *big.Int
	ConditionID common.Hash
}

type PolymarketNegRiskPositionConverted struct {
	NoTokensGiven     []PolymarketToken
	YesTokensReceived []PolymarketToken
	AmountIn          *big.Int
	AmountOut         *big.Int
	CollateralOut     *big.Int // -fee included
}

type PolymarketUmaOptimisticOracleV2Settle struct {
	Requester     common.Address `json:"requester"`
	Proposer      common.Address `json:"proposer"`
	Disputer      common.Address `json:"disputer"`
	Identifier    common.Hash    `json:"identifier"`
	Timestamp     *big.Int       `json:"timestamp"`
	AncillaryData []byte         `json:"ancillary_data"`
	Price         *big.Int       `json:"proposed_price"`
	Payout        *big.Int       `json:"payout"`

	QuestionID common.Hash `json:"question_id"`
}

type PolymarketMarketEssentialGammaInfo struct {
	IsPresentedInGamma bool       `db:"is_presented_in_gamma" json:"is_presented_in_gamma"`
	ID                 int32      `db:"gamma_id" json:"gamma_id"`
	Question           string     `db:"gamma_question" json:"gamma_question"`
	Description        string     `db:"gamma_description" json:"gamma_description"`
	Slug               string     `db:"gamma_slug" json:"gamma_slug"`
	EventSlug          string     `db:"gamma_event_slug" json:"gamma_event_slug"`
	EventLink          string     `db:"gamma_event_link" json:"gamma_event_link"`
	ResolutionSource   *string    `db:"gamma_resolution_source" json:"gamma_resolution_source"`
	StartDate          *time.Time `db:"gamma_start_date" json:"gamma_start_date"`
	EndDate            *time.Time `db:"gamma_end_date" json:"gamma_end_date"`
	CreatedAt          *time.Time `db:"gamma_created_at" json:"gamma_created_at"`
	UpdatedAt          *time.Time `db:"gamma_updated_at" json:"gamma_updated_at"`
	ClosedAt           *time.Time `db:"gamma_closed_at" json:"gamma_closed_at"`

	EventImageURL  *string `db:"gamma_event_image_url" json:"gamma_event_image_url"`
	EventIconURL   *string `db:"gamma_event_icon_url" json:"gamma_event_icon_url"`
	MarketImageURL *string `db:"gamma_market_image_url" json:"gamma_market_image_url"`
	MarketIconURL  *string `db:"gamma_market_icon_url" json:"gamma_market_icon_url"`

	Outcomes        []string `db:"gamma_outcomes" json:"gamma_outcomes"`
	GroupItemTitle  *string  `db:"gamma_group_item_title" json:"gamma_group_item_title"`
	AcceptingOrders bool     `db:"gamma_accepting_orders" json:"gamma_accepting_orders"`

	OrderPriceMinTickSize *decimal.Decimal `db:"gamma_order_price_min_tick_size" json:"gamma_order_price_min_tick_size"`
	OrderMinSize          *decimal.Decimal `db:"gamma_order_min_size" json:"gamma_order_min_size"`

	TagSlugs []string `db:"gamma_tag_slugs" json:"gamma_tag_slugs"`

	NegRisk               *bool            `db:"gamma_neg_risk" json:"gamma_neg_risk"`
	NegRiskRequestID      *string          `db:"gamma_neg_risk_request_id" json:"gamma_neg_risk_request_id"`
	NegRiskOther          *bool            `db:"gamma_neg_risk_other" json:"gamma_neg_risk_other"`
	UmaResolutionStatus   *string          `db:"gamma_uma_resolution_status" json:"gamma_uma_resolution_status"`
	UmaResolutionStatuses []string         `db:"gamma_uma_resolution_statuses" json:"gamma_uma_resolution_statuses"`
	UmaBond               *decimal.Decimal `db:"gamma_uma_bond" json:"gamma_uma_bond"`
	UmaReward             *decimal.Decimal `db:"gamma_uma_reward" json:"gamma_uma_reward"`

	FeesEnabled *bool            `db:"gamma_fees_enabled" json:"gamma_fees_enabled"`
	Fee         *decimal.Decimal `db:"gamma_fee" json:"gamma_fee"`

	RawResponse []byte `db:"gamma_raw_response" json:"gamma_raw_response"`

	SystemUpdatedAt time.Time `db:"system_updated_at" json:"system_updated_at"`
}

type PolymarketUmaOptimisticOracleV2ProposePrice struct {
	Requester           common.Address `json:"requester"`
	Proposer            common.Address `json:"proposer"`
	Identifier          common.Hash    `json:"identifier"`
	Timestamp           *big.Int       `json:"timestamp"`
	AncillaryData       []byte         `json:"ancillary_data"`
	ProposedPrice       *big.Int       `json:"proposed_price"`
	ExpirationTimestamp *big.Int       `json:"expiration_timestamp"`
	Currency            common.Address `json:"currency"`

	QuestionID common.Hash `json:"question_id"`
}

type PolymarketUmaOptimisticOracleV2DisputePrice struct {
	Requester     common.Address `json:"requester"`
	Proposer      common.Address `json:"proposer"`
	Disputer      common.Address `json:"disputer"`
	Identifier    common.Hash    `json:"identifier"`
	Timestamp     *big.Int       `json:"timestamp"`
	AncillaryData []byte         `json:"ancillary_data"`
	ProposedPrice *big.Int       `json:"proposed_price"`

	QuestionID common.Hash `json:"question_id"`
}

type PolymarketUmaOptimisticOracleV2RequestPrice struct {
	Requester     common.Address `json:"requester"`
	Identifier    common.Hash    `json:"identifier"`
	Timestamp     *big.Int       `json:"timestamp"`
	AncillaryData []byte         `json:"ancillary_data"`
	Currency      common.Address `json:"currency"`
	Reward        *big.Int       `json:"reward"`
	FinalFee      *big.Int       `json:"final_fee"`

	QuestionID common.Hash `json:"question_id"`
}

type ParsedAncillaryData struct {
	Question    *string `json:"question,omitempty"`
	Description *string `json:"description,omitempty"`
	ResData     *string `json:"res_data,omitempty"`
	MarketID    *int64  `json:"market_id,omitempty"`
	Initializer *string `json:"initializer,omitempty"`
}

func (d ParsedAncillaryData) IsZero() bool {
	return d.Question == nil &&
		d.Description == nil &&
		d.ResData == nil &&
		d.MarketID == nil &&
		d.Initializer == nil
}

type PolymarketFPMMBuy struct {
	FPMMAddress         common.Address `json:"fpmm_address"`
	Buyer               common.Address `json:"buyer"`
	InvestmentAmount    *big.Int       `json:"investment_amount"`
	FeeAmount           *big.Int       `json:"fee_amount"`
	OutcomeIndex        *big.Int       `json:"outcome_index"`
	OutcomeTokensBought *big.Int       `json:"outcome_tokens_bought"`
}

type PolymarketFPMMSell struct {
	FPMMAddress       common.Address `json:"fpmm_address"`
	Seller            common.Address `json:"seller"`
	ReturnAmount      *big.Int       `json:"return_amount"`
	FeeAmount         *big.Int       `json:"fee_amount"`
	OutcomeIndex      *big.Int       `json:"outcome_index"`
	OutcomeTokensSold *big.Int       `json:"outcome_tokens_sold"`
}

type PolymarketPayoutRedemption struct {
	Redeemer        common.Address `json:"redeemer"`
	ConditionID     common.Hash    `json:"condition_id"`
	CollateralToken common.Address `json:"collateral_token"`
	TokenIDs        []*big.Int     `json:"token_ids"`
	Amounts         []*big.Int     `json:"amounts"`
	Payout          *big.Int       `json:"payout"`
}
