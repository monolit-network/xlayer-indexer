package polymarket

import (
	"fmt"

	sharedModels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/ethereum/go-ethereum/core/types"
)

func (p *CTFParser) ParseFPMMBuyLog(log *types.Log) (*sharedModels.PolymarketFPMMBuy, error) {
	if len(log.Topics) != 3 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	buy, err := polymarketFPMMContract.ParseFPMMBuy(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing buy log: %w", err)
	}

	return &sharedModels.PolymarketFPMMBuy{
		FPMMAddress:         log.Address,
		Buyer:               buy.Buyer,
		InvestmentAmount:    buy.InvestmentAmount,
		FeeAmount:           buy.FeeAmount,
		OutcomeIndex:        buy.OutcomeIndex,
		OutcomeTokensBought: buy.OutcomeTokensBought,
	}, nil
}

func (p *CTFParser) ParseFPMMSellLog(log *types.Log) (*sharedModels.PolymarketFPMMSell, error) {
	if len(log.Topics) != 3 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	sell, err := polymarketFPMMContract.ParseFPMMSell(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing sell log: %w", err)
	}

	return &sharedModels.PolymarketFPMMSell{
		FPMMAddress:       log.Address,
		Seller:            sell.Seller,
		ReturnAmount:      sell.ReturnAmount,
		FeeAmount:         sell.FeeAmount,
		OutcomeIndex:      sell.OutcomeIndex,
		OutcomeTokensSold: sell.OutcomeTokensSold,
	}, nil
}
