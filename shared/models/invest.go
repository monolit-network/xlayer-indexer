package models

import "time"

type InvestorID int64

type Investor struct {
	ID              InvestorID     `json:"id"` // own id, not from dropstab
	Slug            string         `json:"slug"`
	Name            string         `json:"name"`
	Image           string         `json:"image"`
	Rank            int32          `json:"rank"`
	Country         string         `json:"country"`
	Description     string         `json:"description"`
	VentureType     string         `json:"venture_type"`
	Tier            string         `json:"tier"`
	LeadInvestments int32          `json:"lead_investments"`
	Links           []InvestorLink `json:"links"`

	PortfolioCoins      []VerifiedTokenAssetID `json:"portfolio_coins"`
	PortfolioCoinsSlugs []string               `json:"portfolio_coins_slugs"`

	RoundDistributionByCategory []InvestorRoundDistribution `json:"round_distribution_by_category"`
	RoundDistributionByStage    []InvestorRoundDistribution `json:"round_distribution_by_stage"`
}

type InvestorRoundDistribution struct {
	Name    string  `json:"name"`
	Amount  int32   `json:"amount"`
	Percent float32 `json:"percent"`
}

type InvestorLink struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

type InvestRoundID int64

type InvestRound struct {
	ID                     InvestRoundID        `json:"id"`
	AssetID                VerifiedTokenAssetID `json:"asset_id"`
	CoinSlug               string               `json:"coin_slug"`
	FundsRaised            int                  `json:"funds_raised"`
	PreValuation           int                  `json:"pre_valuation"`
	PreValuationInaccurate bool                 `json:"pre_valuation_inaccurate"`
	Stage                  string               `json:"stage"`
	Category               string               `json:"category"`
	Date                   time.Time            `json:"date"`
	Investors              []InvestorID         `json:"investors"`
	LeadInvestors          []InvestorID         `json:"lead_investors"`

	InvestorsSlugs     []string `json:"investors_slugs"`
	LeadInvestorsSlugs []string `json:"lead_investors_slugs"`
}
