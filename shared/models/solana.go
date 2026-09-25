package models

import "time"

type SolanaSwapEvent struct {
	BlockTime                    time.Time `json:"block_time"`
	Slot                         uint32    `json:"slot"`
	TxIdx                        uint16    `json:"tx_idx"`
	Ix                           uint16    `json:"ix"`
	SigningWallet                string    `json:"signing_wallet"`
	TokenSold                    string    `json:"token_sold"`
	TokenBought                  string    `json:"token_bought"`
	AmountSold                   uint64    `json:"amount_sold"`
	AmountBought                 uint64    `json:"amount_bought"`
	PoolTokenSoldBalanceBefore   uint64    `json:"pool_token_sold_balance_before"`
	PoolTokenBoughtBalanceBefore uint64    `json:"pool_token_bought_balance_before"`
	PoolTokenSoldBalanceAfter    uint64    `json:"pool_token_sold_balance_after"`
	PoolTokenBoughtBalanceAfter  uint64    `json:"pool_token_bought_balance_after"`
	Signature                    string    `json:"signature"`
	Source                       string    `json:"source"`
}

type SolanaTransferEvent struct {
	BlockTime time.Time `json:"block_time"`
	Slot      uint32    `json:"slot"`
	TxIdx     uint16    `json:"tx_idx"`
	Ix        uint16    `json:"ix"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Token     string    `json:"token"`
	Amount    uint64    `json:"amount"`
	Signature string    `json:"signature"`
	Type      string    `json:"type"`
}

type SolanaMintEvent struct {
	BlockTime time.Time `json:"block_time"`
	Slot      uint32    `json:"slot"`
	TxIdx     uint16    `json:"tx_idx"`
	Ix        uint16    `json:"ix"`
	Mint      string    `json:"mint"`
	To        string    `json:"to"`
	Amount    uint64    `json:"amount"`
	Authority string    `json:"authority"`
	Signature string    `json:"signature"`
	Type      string    `json:"type"`
}

type SolanaBurnEvent struct {
	BlockTime time.Time `json:"block_time"`
	Slot      uint32    `json:"slot"`
	TxIdx     uint16    `json:"tx_idx"`
	Ix        uint16    `json:"ix"`
	Mint      string    `json:"mint"`
	From      string    `json:"from"`
	Amount    uint64    `json:"amount"`
	Authority string    `json:"authority"`
	Signature string    `json:"signature"`
	Type      string    `json:"type"`
}
