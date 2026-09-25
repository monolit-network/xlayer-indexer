package polymarketProcessor

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	evmmetrics "github.com/monolit-network/xlayer-indexer/evm/pkg/metrics"
	"github.com/monolit-network/xlayer-indexer/shared/db"
	sharedmodels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

func (s *Processor) insertTokenMapping(conditionID common.Hash, tokens []sharedmodels.PolymarketToken, blockTime time.Time) error {
	tokenIDStrs := make([]string, len(tokens))
	for i, t := range tokens {
		tokenIDStrs[i] = t.TokenID.String()
	}

	conditionIDStr := strings.ToLower(conditionID.Hex())
	if s.isTokenMappingCached(conditionIDStr, tokenIDStrs) {
		return nil
	}

	var tokensInDB []sharedmodels.PolymarketToken
	err := s.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
		var dbErr error
		tokensInDB, dbErr = q.GetPolymarketTokensByConditionID(ctx, conditionIDStr)
		return dbErr
	})
	if err != nil {
		return fmt.Errorf("failed to get tokens from db: %w", err)
	}

	newTokens := make([]sharedmodels.PolymarketToken, 0, len(tokens))
	existingTokens := make([]sharedmodels.PolymarketToken, 0, len(tokensInDB))
	for _, t := range tokens {
		found := false
		for _, tInDB := range tokensInDB {
			if tInDB.TokenID.String() == t.TokenID.String() {
				found = true
				break
			}
		}
		if !found {
			newTokens = append(newTokens, t)
		} else {
			existingTokens = append(existingTokens, t)
		}
	}

	for _, t := range existingTokens {
		s.tokenToConditionCache.Add(t.TokenID.String(), conditionID)
	}

	for _, t := range newTokens {
		t.CreatedAt = blockTime
		s.pgWritersCh <- &t
	}

	for _, t := range newTokens {
		s.tokenToConditionCache.Add(t.TokenID.String(), conditionID)
	}

	s.cacheTokenMapping(conditionIDStr, tokenIDStrs)

	return nil
}

func (s *Processor) tokenMappingWriter() {
	flush := func(batch []sharedmodels.PolymarketToken) []sharedmodels.PolymarketToken {
		if len(batch) == 0 {
			return batch
		}
		if err := s.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
			return q.InsertPolymarketTokensBatch(ctx, batch)
		}); err != nil {
			s.logger.Error("failed to insert tokens into db", zap.Error(err), zap.Any("tokens", batch))
			evmmetrics.InsertErrors.WithLabelValues(s.chain, "polymarket", "postgres", "polymarket.tokens").Inc()
			return batch
		}
		evmmetrics.InsertedRows.WithLabelValues(s.chain, "polymarket", "postgres", "polymarket.tokens").Add(float64(len(batch)))
		return batch[:0]
	}

	batch := make([]sharedmodels.PolymarketToken, 0, 128)
	for {
		select {
		case request, ok := <-s.pgWritersCh:
			if !ok {
				flush(batch)
				return
			}
			batch = append(batch, *request)
			if len(batch) >= 128 {
				batch = flush(batch)
			}
		case done := <-s.pgWriterFlushCh:
			batch = flush(batch)
			close(done)
		}
	}
}

func (s *Processor) isTokenMappingCached(conditionID string, tokenIDs []string) bool {
	s.tokenMappingMu.RLock()
	defer s.tokenMappingMu.RUnlock()

	cachedTokens, ok := s.conditionTokenIDs[conditionID]
	if !ok {
		return false
	}
	if len(cachedTokens) != len(tokenIDs) {
		return false
	}
	for _, tID := range tokenIDs {
		if _, ok := cachedTokens[tID]; !ok {
			return false
		}
	}
	return len(tokenIDs) > 0
}

func (s *Processor) cacheTokenMapping(conditionID string, tokenIDs []string) {
	s.tokenMappingMu.Lock()
	defer s.tokenMappingMu.Unlock()

	if len(tokenIDs) == 0 {
		return
	}

	existing, ok := s.conditionTokenIDs[conditionID]
	if !ok {
		existing = make(map[string]struct{}, len(tokenIDs))
		s.conditionTokenIDs[conditionID] = existing
	}
	for _, tID := range tokenIDs {
		existing[tID] = struct{}{}
	}
}

func (s *Processor) getConditionIDByTokenID(tokenID *big.Int) (common.Hash, error) {
	tokenIDStr := tokenID.String()

	if conditionID, ok := s.tokenToConditionCache.Get(tokenIDStr); ok {
		return conditionID, nil
	}

	var token *sharedmodels.PolymarketToken
	err := s.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
		var dbErr error
		token, dbErr = q.GetPolymarketTokenByTokenID(ctx, tokenID)
		return dbErr
	})
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get token from db: %w", err)
	}
	if token == nil {
		return common.Hash{}, fmt.Errorf("token not found: %s", tokenIDStr)
	}

	conditionID := common.HexToHash(token.ConditionID)
	s.tokenToConditionCache.Add(tokenIDStr, conditionID)

	return conditionID, nil
}

func (s *Processor) cacheConditionQuestionMapping(conditionID, questionID string) error {
	conditionIDLower := strings.ToLower(conditionID)
	questionIDLower := strings.ToLower(questionID)

	s.conditionQuestionMappingMu.Lock()
	defer s.conditionQuestionMappingMu.Unlock()
	s.conditionQuestionMapping.Add(conditionIDLower, questionIDLower)
	s.questionConditionMapping.Add(questionIDLower, conditionIDLower)

	return nil
}

func (s *Processor) getQuestionIDByConditionID(conditionID string) (string, error) {
	conditionIDLower := strings.ToLower(conditionID)

	s.conditionQuestionMappingMu.RLock()
	questionID, ok := s.conditionQuestionMapping.Get(conditionIDLower)
	s.conditionQuestionMappingMu.RUnlock()

	if !ok {
		var market *sharedmodels.PolymarketMarketNewMinimal
		err := s.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
			var dbErr error
			market, dbErr = q.GetPolymarketMarketMinimalByConditionID(ctx, conditionIDLower)
			return dbErr
		})
		if err != nil {
			return "", fmt.Errorf("failed to get market from db: %w", err)
		}
		if market == nil {
			return "", nil
		}

		questionID := strings.ToLower(market.QuestionID)
		if cacheErr := s.cacheConditionQuestionMapping(conditionIDLower, questionID); cacheErr != nil {
			return "", fmt.Errorf("failed to cache mapping: %w", cacheErr)
		}

		return questionID, nil
	}
	return questionID, nil
}

func (s *Processor) setQuestionTimestamp(questionID string, timestamp uint64) {
	questionIDLower := strings.ToLower(questionID)
	s.questionTimestampMu.Lock()
	s.questionTimestampCache[questionIDLower] = timestamp
	s.questionTimestampMu.Unlock()
}

func (s *Processor) clearQuestionTimestamp(questionID string) {
	questionIDLower := strings.ToLower(questionID)
	s.questionTimestampMu.Lock()
	delete(s.questionTimestampCache, questionIDLower)
	s.questionTimestampMu.Unlock()
}

func (s *Processor) getQuestionTimestampByQuestionID(questionID string, maxBlockNumber uint64) (uint64, error) {
	questionIDLower := strings.ToLower(questionID)

	s.questionTimestampMu.RLock()
	if cachedTimestamp, ok := s.questionTimestampCache[questionIDLower]; ok {
		s.questionTimestampMu.RUnlock()
		return cachedTimestamp, nil
	}
	s.questionTimestampMu.RUnlock()

	var timestamp uint64
	err := s.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
		var dbErr error
		timestamp, dbErr = q.GetPolymarketQuestionTimestampFromEvents(ctx, questionIDLower, maxBlockNumber)
		return dbErr
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get question timestamp from events: %w", err)
	}

	s.questionTimestampMu.Lock()
	s.questionTimestampCache[questionIDLower] = timestamp
	s.questionTimestampMu.Unlock()

	return timestamp, nil
}

func (s *Processor) setTokenIDToPayout(tokenID *big.Int, payout tokenIDToPayout) {
	tokenIDStr := tokenID.String()
	s.tokenIDToPayoutMu.Lock()
	s.tokenIDToPayout[tokenIDStr] = payout
	s.tokenIDToPayoutMu.Unlock()
}

func (s *Processor) getTokenIDToPayout(tokenID *big.Int) (tokenIDToPayout, bool, error) {
	tokenIDStr := tokenID.String()
	s.tokenIDToPayoutMu.RLock()
	payout, ok := s.tokenIDToPayout[tokenIDStr]
	s.tokenIDToPayoutMu.RUnlock()
	if ok {
		return payout, true, nil
	}

	tokenInfo, err := s.dbClient.Querier().GetPolymarketTokenByTokenID(context.Background(), tokenID)
	if err != nil {
		return tokenIDToPayout{}, false, fmt.Errorf("failed to get token info from db: %w", err)
	}
	if tokenInfo == nil {
		s.flushPgSync()
		tokenInfo, err = s.dbClient.Querier().GetPolymarketTokenByTokenID(context.Background(), tokenID)
	}
	if tokenInfo == nil {
		return tokenIDToPayout{}, false, fmt.Errorf("token info not found: %s", tokenIDStr)
	}
	if tokenInfo.Denominator == nil || tokenInfo.Numerator == nil {
		return tokenIDToPayout{}, false, fmt.Errorf("numerator or denominator is nil: %s", tokenIDStr)
	}
	payout = tokenIDToPayout{
		Numerator:   tokenInfo.Numerator,
		Denominator: tokenInfo.Denominator,
	}
	s.setTokenIDToPayout(tokenID, payout)
	return payout, true, nil
}
