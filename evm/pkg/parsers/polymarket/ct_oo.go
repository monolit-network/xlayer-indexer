package polymarket

import (
	"fmt"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	sharedModels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (p *CTFParser) IsUmaOptimisticOracleV2Address(address common.Address) bool {
	return address == models.UmaOptimisticOracleV2Address || address == models.UmaPolymarketManagedOptimisticOracleV2Address
}

func (p *CTFParser) ParseUmaOptimisticOracleV2ProposePriceLog(log *types.Log) (*sharedModels.PolymarketUmaOptimisticOracleV2ProposePrice, error) {
	if !p.IsUmaOptimisticOracleV2Address(log.Address) {
		return nil, nil
	}
	if len(log.Topics) > 0 && log.Topics[0] != models.UmaOptimisticOracleV2ProposePriceEventSelectorHash {
		return nil, nil
	}

	umaOptimisticOracleV2ProposePrice, err := umaOptimisticOracleV2Contract.ParseProposePrice(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing uma optimistic oracle v2 propose price log: %w", err)
	}

	return &sharedModels.PolymarketUmaOptimisticOracleV2ProposePrice{
		Requester:           umaOptimisticOracleV2ProposePrice.Requester,
		Proposer:            umaOptimisticOracleV2ProposePrice.Proposer,
		Identifier:          umaOptimisticOracleV2ProposePrice.Identifier,
		Timestamp:           umaOptimisticOracleV2ProposePrice.Timestamp,
		AncillaryData:       umaOptimisticOracleV2ProposePrice.AncillaryData,
		ProposedPrice:       umaOptimisticOracleV2ProposePrice.ProposedPrice,
		ExpirationTimestamp: umaOptimisticOracleV2ProposePrice.ExpirationTimestamp,
		Currency:            umaOptimisticOracleV2ProposePrice.Currency,

		QuestionID: p.questionIDFromAncillaryData(umaOptimisticOracleV2ProposePrice.AncillaryData),
	}, nil
}

func (p *CTFParser) ParseUmaOptimisticOracleV2DisputePriceLog(log *types.Log) (*sharedModels.PolymarketUmaOptimisticOracleV2DisputePrice, error) {
	if !p.IsUmaOptimisticOracleV2Address(log.Address) {
		return nil, nil
	}
	if len(log.Topics) > 0 && log.Topics[0] != models.UmaOptimisticOracleV2DisputePriceEventSelectorHash {
		return nil, nil
	}

	umaOptimisticOracleV2DisputePrice, err := umaOptimisticOracleV2Contract.ParseDisputePrice(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing uma optimistic oracle v2 dispute price log: %w", err)
	}

	return &sharedModels.PolymarketUmaOptimisticOracleV2DisputePrice{
		Requester:     umaOptimisticOracleV2DisputePrice.Requester,
		Proposer:      umaOptimisticOracleV2DisputePrice.Proposer,
		Disputer:      umaOptimisticOracleV2DisputePrice.Disputer,
		Identifier:    umaOptimisticOracleV2DisputePrice.Identifier,
		Timestamp:     umaOptimisticOracleV2DisputePrice.Timestamp,
		AncillaryData: umaOptimisticOracleV2DisputePrice.AncillaryData,
		ProposedPrice: umaOptimisticOracleV2DisputePrice.ProposedPrice,

		QuestionID: p.questionIDFromAncillaryData(umaOptimisticOracleV2DisputePrice.AncillaryData),
	}, nil
}

func (p *CTFParser) ParseUmaOptimisticOracleV2SettleLog(log *types.Log) (*sharedModels.PolymarketUmaOptimisticOracleV2Settle, error) {
	if !p.IsUmaOptimisticOracleV2Address(log.Address) {
		return nil, nil
	}
	if len(log.Topics) > 0 && log.Topics[0] != models.UmaOptimisticOracleV2SettleEventSelectorHash {
		return nil, nil
	}

	umaOptimisticOracleV2Settle, err := umaOptimisticOracleV2Contract.ParseSettle(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing uma optimistic oracle v2 settle log: %w", err)
	}

	return &sharedModels.PolymarketUmaOptimisticOracleV2Settle{
		Requester:     umaOptimisticOracleV2Settle.Requester,
		Proposer:      umaOptimisticOracleV2Settle.Proposer,
		Disputer:      umaOptimisticOracleV2Settle.Disputer,
		Identifier:    umaOptimisticOracleV2Settle.Identifier,
		Timestamp:     umaOptimisticOracleV2Settle.Timestamp,
		AncillaryData: umaOptimisticOracleV2Settle.AncillaryData,
		Price:         umaOptimisticOracleV2Settle.Price,
		Payout:        umaOptimisticOracleV2Settle.Payout,

		QuestionID: p.questionIDFromAncillaryData(umaOptimisticOracleV2Settle.AncillaryData),
	}, nil
}

func (p *CTFParser) ParseUmaOptimisticOracleV2RequestPriceLog(log *types.Log) (*sharedModels.PolymarketUmaOptimisticOracleV2RequestPrice, error) {
	if !p.IsUmaOptimisticOracleV2Address(log.Address) {
		return nil, nil
	}
	if len(log.Topics) > 0 && log.Topics[0] != models.UmaOptimisticOracleV2RequestPriceEventSelectorHash {
		return nil, nil
	}

	umaOptimisticOracleV2RequestPrice, err := umaOptimisticOracleV2Contract.ParseRequestPrice(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing uma optimistic oracle v2 request price log: %w", err)
	}

	return &sharedModels.PolymarketUmaOptimisticOracleV2RequestPrice{
		Requester:     umaOptimisticOracleV2RequestPrice.Requester,
		Identifier:    umaOptimisticOracleV2RequestPrice.Identifier,
		Timestamp:     umaOptimisticOracleV2RequestPrice.Timestamp,
		AncillaryData: umaOptimisticOracleV2RequestPrice.AncillaryData,
		Currency:      umaOptimisticOracleV2RequestPrice.Currency,
		Reward:        umaOptimisticOracleV2RequestPrice.Reward,
		FinalFee:      umaOptimisticOracleV2RequestPrice.FinalFee,

		QuestionID: p.questionIDFromAncillaryData(umaOptimisticOracleV2RequestPrice.AncillaryData),
	}, nil
}
