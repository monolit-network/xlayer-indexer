package polymarket

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	sharedModels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var ancillaryDataKeys = []string{
	"q:", "title:", "description:", "desc:", "market_id:", "res_data:", "initializer:",
}

func (p *CTFParser) ParseAncillaryData(data []byte) sharedModels.ParsedAncillaryData {
	if len(data) == 0 {
		return sharedModels.ParsedAncillaryData{}
	}

	trimmed := strings.TrimSpace(strings.Trim(string(data), "\x00"))
	if trimmed == "" {
		return sharedModels.ParsedAncillaryData{}
	}

	parsed := sharedModels.ParsedAncillaryData{}

	question := extractAncillaryField(trimmed, []string{"title:"})
	if question == "" {
		question = extractAncillaryField(trimmed, []string{"q:"})
	}

	initializer := extractAncillaryField(trimmed, []string{"initializer:"})
	description := extractAncillaryField(trimmed, []string{"description:", "desc:"})
	resData := extractAncillaryField(trimmed, []string{"res_data:"})
	marketIDStr := extractAncillaryField(trimmed, []string{"market_id:"})

	if question != "" {
		parsed.Question = &question
	}
	if description != "" {
		parsed.Description = &description
	}
	if resData != "" {
		parsed.ResData = &resData
	}
	if marketIDStr != "" {
		if marketIDInt, err := strconv.ParseInt(marketIDStr, 10, 64); err == nil {
			parsed.MarketID = &marketIDInt
		}
	}
	if initializer != "" {
		parsed.Initializer = &initializer
	}

	return parsed
}
func extractAncillaryField(data string, keys []string) string {
	lowerData := strings.ToLower(data)

	for _, key := range keys {
		idx := strings.Index(lowerData, key)
		if idx == -1 {
			continue
		}

		start := idx + len(key)
		rest := data[start:]
		lowerRest := lowerData[start:]

		finalEnd := -1
		for _, k := range ancillaryDataKeys {
			kIdx := strings.Index(lowerRest, k)
			if kIdx != -1 {
				if finalEnd == -1 || kIdx < finalEnd {
					finalEnd = kIdx
				}
			}
		}

		var result string
		if finalEnd == -1 {
			result = rest
		} else {
			result = rest[:finalEnd]
		}

		return cleanValue(result)
	}
	return ""
}

func cleanValue(s string) string {
	s = strings.ReplaceAll(s, "\x00", "")
	s = strings.TrimSpace(s)
	for strings.HasSuffix(s, ",") || strings.HasSuffix(s, ".") {
		s = strings.TrimSuffix(s, ",")
		s = strings.TrimSuffix(s, ".")
		s = strings.TrimSpace(s)
	}
	return s
}

func (p *CTFParser) IsPolymarketUmaCtfAdapterAddress(address common.Address) bool {
	_, ok := models.PolymarketCTFAdaptersAddressesMap[address]
	return ok
}

func (p *CTFParser) ParseUmaCtfAdapterQuestionInitializedLog(log *types.Log) (*sharedModels.PolymarketUmaCtfAdapterQuestionInitialized, error) {
	if !p.IsPolymarketUmaCtfAdapterAddress(log.Address) {
		return nil, nil
	}
	if len(log.Topics) > 0 && log.Topics[0] != models.PolymarketUmaCtfAdapterQuestionInitializedEventSelectorHash {
		return nil, nil
	}

	umaCtfAdapterQuestionInitialized, err := polymarketUmaCtfAdapterContract.ParseQuestionInitialized(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing uma ctf adapter question initialized log: %w", err)
	}

	return &sharedModels.PolymarketUmaCtfAdapterQuestionInitialized{
		QuestionID:       umaCtfAdapterQuestionInitialized.QuestionID,
		RequestTimestamp: umaCtfAdapterQuestionInitialized.RequestTimestamp,
		Creator:          umaCtfAdapterQuestionInitialized.Creator,
		AncillaryData:    umaCtfAdapterQuestionInitialized.AncillaryData,
		RewardToken:      umaCtfAdapterQuestionInitialized.RewardToken,
		Reward:           umaCtfAdapterQuestionInitialized.Reward,
		ProposalBond:     umaCtfAdapterQuestionInitialized.ProposalBond,
	}, nil
}

func (p *CTFParser) ParseUmaCtfAdapterQuestionResetLog(log *types.Log) (*sharedModels.PolymarketUmaCtfAdapterQuestionReset, error) {
	if !p.IsPolymarketUmaCtfAdapterAddress(log.Address) {
		return nil, nil
	}
	if len(log.Topics) > 0 && log.Topics[0] != models.PolymarketUmaCtfAdapterQuestionResetEventSelectorHash {
		return nil, nil
	}

	umaCtfAdapterQuestionReset, err := polymarketUmaCtfAdapterContract.ParseQuestionReset(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing uma ctf adapter question reset log: %w", err)
	}

	return &sharedModels.PolymarketUmaCtfAdapterQuestionReset{
		QuestionID: umaCtfAdapterQuestionReset.QuestionID,
	}, nil
}

func (p *CTFParser) ParseUmaCtfAdapterQuestionResolvedLog(log *types.Log) (*sharedModels.PolymarketUmaCtfAdapterQuestionResolved, error) {
	if !p.IsPolymarketUmaCtfAdapterAddress(log.Address) {
		return nil, nil
	}
	if len(log.Topics) > 0 && log.Topics[0] != models.PolymarketUmaCtfAdapterQuestionResolvedEventSelectorHash {
		return nil, nil
	}

	umaCtfAdapterQuestionResolved, err := polymarketUmaCtfAdapterContract.ParseQuestionResolved(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing uma ctf adapter question resolved log: %w", err)
	}

	return &sharedModels.PolymarketUmaCtfAdapterQuestionResolved{
		QuestionID:   umaCtfAdapterQuestionResolved.QuestionID,
		SettledPrice: umaCtfAdapterQuestionResolved.SettledPrice,
		Payouts:      umaCtfAdapterQuestionResolved.Payouts,
	}, nil
}

func (p *CTFParser) ParseUmaCtfAdapterQuestionPausedLog(log *types.Log) (*sharedModels.PolymarketUmaCtfAdapterQuestionPaused, error) {
	if !p.IsPolymarketUmaCtfAdapterAddress(log.Address) {
		return nil, nil
	}
	if len(log.Topics) > 0 && log.Topics[0] != models.PolymarketUmaCtfAdapterQuestionPausedEventSelectorHash {
		return nil, nil
	}

	umaCtfAdapterQuestionPaused, err := polymarketUmaCtfAdapterContract.ParseQuestionPaused(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing uma ctf adapter question paused log: %w", err)
	}

	return &sharedModels.PolymarketUmaCtfAdapterQuestionPaused{
		QuestionID: umaCtfAdapterQuestionPaused.QuestionID,
	}, nil
}

func (p *CTFParser) ParseUmaCtfAdapterQuestionFlaggedLog(log *types.Log) (*sharedModels.PolymarketUmaCtfAdapterQuestionFlagged, error) {
	if !p.IsPolymarketUmaCtfAdapterAddress(log.Address) {
		return nil, nil
	}
	if len(log.Topics) > 0 && log.Topics[0] != models.PolymarketUmaCtfAdapterQuestionFlaggedEventSelectorHash {
		return nil, nil
	}

	umaCtfAdapterQuestionFlagged, err := polymarketUmaCtfAdapterContract.ParseQuestionFlagged(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing uma ctf adapter question flagged log: %w", err)
	}

	return &sharedModels.PolymarketUmaCtfAdapterQuestionFlagged{
		QuestionID: umaCtfAdapterQuestionFlagged.QuestionID,
	}, nil
}
