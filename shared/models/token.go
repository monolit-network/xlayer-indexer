package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type TokenType string

const (
	TokenTypeRegular  TokenType = "regular"
	TokenTypeVerified TokenType = "verified"
)

type TokenLike interface {
	TokenType() TokenType

	GetChain() string
	GetContractAddress() string
	GetSymbol() string
	GetName() string
}

type Token struct {
	Chain            string    `json:"chain"`
	ContractAddress  string    `json:"contract_address"`
	Symbol           string    `json:"symbol"`
	Name             string    `json:"name"`
	Decimals         int       `json:"decimals"`
	PriceUSD         float64   `json:"price_usd"`
	MarketCapUSD     float64   `json:"market_cap_usd"`
	Supply           float64   `json:"supply"`
	LargestLPPoolUSD float64   `json:"largest_lp_pool_usd"`
	FirstTxDate      time.Time `json:"first_tx_date"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ViewSource       string    `json:"view_source"`
}

func (t Token) TokenType() TokenType {
	return TokenTypeRegular
}

func (t Token) GetChain() string {
	return t.Chain
}

func (t Token) GetContractAddress() string {
	return t.ContractAddress
}

func (t Token) GetSymbol() string {
	return t.Symbol
}

func (t Token) GetName() string {
	return t.Name
}

type VerifiedTokenAssetID int

type VerifiedToken struct {
	AssetID         VerifiedTokenAssetID `json:"asset_id"` // own id, not from dropstab
	Chain           string               `json:"chain"`
	ContractAddress string               `json:"contract_address"`
	Symbol          string               `json:"symbol"`
	Name            string               `json:"name"`
	Decimals        *int                 `json:"decimals,omitempty"`
	CreatedAt       *time.Time           `json:"created_at,omitempty"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

func (t VerifiedToken) TokenType() TokenType {
	return TokenTypeVerified
}

func (t VerifiedToken) GetChain() string {
	return t.Chain
}

func (t VerifiedToken) GetContractAddress() string {
	return t.ContractAddress
}

func (t VerifiedToken) GetSymbol() string {
	return t.Symbol
}

func (t VerifiedToken) GetName() string {
	return t.Name
}

type TokenMessage struct {
	Token     TokenLike `json:"token"`
	TokenType TokenType `json:"type"`
}

func TokenMessageFromTokenLike(t TokenLike) TokenMessage {
	return TokenMessage{
		Token:     t,
		TokenType: t.TokenType(),
	}
}

func (msg *TokenMessage) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Type TokenType `json:"type"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	var t TokenLike
	switch tmp.Type {
	case TokenTypeVerified:
		t = &VerifiedToken{}
	case TokenTypeRegular:
		t = &Token{}
	default:
		return fmt.Errorf("invalid token type: %s", tmp.Type)
	}

	if err := json.Unmarshal(data, t); err != nil {
		return err
	}

	msg.Token = t
	msg.TokenType = tmp.Type

	return nil
}

func (msg *TokenMessage) TokenLike() TokenLike {
	return msg.Token
}

type DropstabInfo struct {
	AssetID           VerifiedTokenAssetID `json:"asset_id"` // own id, not from dropstab
	Slug              string               `json:"slug"`
	Status            string               `json:"status"`
	Symbol            string               `json:"symbol"`
	Name              string               `json:"name"`
	PriceUSD          *float64             `json:"-"`
	MaxSupply         *float64             `json:"max_supply"`
	CirculatingSupply *float64             `json:"circulating_supply"`
	TotalSupply       *float64             `json:"total_supply"`
	Categories        []string             `json:"categories"`
	Socials           TokenSocials         `json:"socials"`
	ImageURL          *string              `json:"image_url"`

	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DropstabInfoWithScore struct {
	DropstabInfo DropstabInfo `json:"dropstab_info"`
	Score        float32      `json:"score"`
}

type DropstabShortCoinDescription struct {
	Slug      string    `json:"slug"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TokenSocials struct {
	Website      []string `json:"website"`
	Explorers    []string `json:"explorers"`
	Twitter      *string  `json:"twitter,omitempty"`
	Reddit       *string  `json:"reddit,omitempty"`
	Linkedin     *string  `json:"linkedin,omitempty"`
	Github       *string  `json:"github,omitempty"`
	TelegramChat *string  `json:"telegramChat,omitempty"`
	TelegramANN  *string  `json:"telegramAnn,omitempty"`
	Facebook     *string  `json:"facebook,omitempty"`
	Medium       *string  `json:"medium,omitempty"`
	Blog         *string  `json:"blog,omitempty"`
	Discord      *string  `json:"discord,omitempty"`
	Youtube      *string  `json:"youtube,omitempty"`
	Bridge       []string `json:"bridge"`
}
