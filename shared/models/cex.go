package models

type TickerInfo struct {
	ID             uint64
	BaseCoin       string
	Symbol         string
	Category       string
	BaseCoinParsed string
	Multiplier     uint64
	QuoteCoin      string
	Status         string
}

type KlineInfo struct {
	TickerID   uint64
	OpenTime   uint64
	OpenPrice  float64
	HighPrice  float64
	LowPrice   float64
	ClosePrice float64
	Volume     float64
	Turnover   float64
}

type FundingRateInfo struct {
	TickerID             uint64
	FundingRate          float64
	FundingRateTimestamp uint64
}

type LiquidationInfo struct {
	TickerID             uint64
	Size                 float64
	LiquidationPrice     float64
	LiquidationTimestamp uint64
}

type TradeInfo struct {
	TickerID       uint64
	Size           float64
	Price          float64
	TradeID        string
	TradeTimestamp uint64
}
