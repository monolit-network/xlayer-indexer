package models

import (
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
)

type ChainInfo struct {
	ChainName         models.Chain
	Tokens            []PaymentToken
	LatestBlockOffset int64
	URL               string
}

type PaymentToken struct {
	Address  string `json:"address"`
	Decimal  int32  `json:"decimal"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type Wallet struct {
	Address         string     `json:"address"`
	PrivateKey      []byte     `json:"private_key,omitempty"`
	ReservedOrderID *string    `json:"reserved_order_id,omitempty"`
	ReservedAt      *time.Time `json:"reserved_at,omitempty"`
	CooldownUntil   *time.Time `json:"cooldown_until,omitempty"`
}
