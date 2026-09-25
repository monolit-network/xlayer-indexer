package polymarketProcessor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/shared/db"
	sharedmodels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/monolit-network/xlayer-indexer/util/notifier"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

func (c *Processor) processUmaCtfAdapterQuestionInitializedLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	questionInitialized, err := c.polymarketParser.ParseUmaCtfAdapterQuestionInitializedLog(log)
	if err != nil {
		return err
	}

	if questionInitialized == nil {
		c.logger.Warn("uma ctf adapter question initialized log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	eventData := sharedmodels.PolymarketMarketEventDataAdapterInitialized{
		RequestTimestamp: questionInitialized.RequestTimestamp.Uint64(),
		Creator:          strings.ToLower(questionInitialized.Creator.Hex()),
		RewardToken:      strings.ToLower(questionInitialized.RewardToken.Hex()),
		Reward:           questionInitialized.Reward.String(),
		ProposalBond:     questionInitialized.ProposalBond.String(),
	}

	event := sharedmodels.PolymarketMarketEvent{
		QuestionID:     strings.ToLower(questionInitialized.QuestionID.Hex()),
		EventType:      sharedmodels.PolymarketEventAdapterInitialized,
		BlockNumber:    log.BlockNumber,
		BlockHash:      strings.ToLower(log.BlockHash.Hex()),
		LogIndex:       uint32(log.Index),
		TxIndex:        uint32(log.TxIndex),
		TxHash:         strings.ToLower(log.TxHash.Hex()),
		BlockTimestamp: time.Unix(int64(log.BlockTimestamp), 0),
		EventData:      &eventData,
	}

	notifierEvent := event
	var parsedAncillaryData sharedmodels.ParsedAncillaryData
	if len(questionInitialized.AncillaryData) > 0 {
		parsedAncillaryData = c.polymarketParser.ParseAncillaryData(questionInitialized.AncillaryData)
		if !parsedAncillaryData.IsZero() {
			notifierEventData := eventData
			notifierEventData.ParsedAncillaryData = &parsedAncillaryData
			notifierEvent.EventData = &notifierEventData
		}
	}

	if c.notifier != nil {
		if err := notifier.Publish(c.notifier, context.Background(), string(sharedmodels.NotifierEventAdapterUMA), notifierEvent); err != nil {
			c.notifier.Logger.Error("failed to publish adapter initialized event", zap.Error(err))
		}
	}

	if err := c.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
		return c.insertPolymarketMarketEvent(ctx, q, event)
	}); err != nil {
		return err
	}

	questionID := strings.ToLower(questionInitialized.QuestionID.Hex())
	requestTimestamp := questionInitialized.RequestTimestamp.Uint64()
	c.setQuestionTimestamp(questionID, requestTimestamp)

	if len(questionInitialized.AncillaryData) > 0 {
		if err := c.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
			return q.UpdatePolymarketMarketAncillaryDataByQuestionID(ctx, questionID, questionInitialized.AncillaryData, parsedAncillaryData)
		}); err != nil {
			return err
		}
	}

	return nil
}

func (c *Processor) processUmaCtfAdapterQuestionResetLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	questionReset, err := c.polymarketParser.ParseUmaCtfAdapterQuestionResetLog(log)
	if err != nil {
		return err
	}

	if questionReset == nil {
		c.logger.Warn("uma ctf adapter question reset log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	event := sharedmodels.PolymarketMarketEvent{
		QuestionID:     strings.ToLower(questionReset.QuestionID.Hex()),
		EventType:      sharedmodels.PolymarketEventAdapterReset,
		BlockNumber:    log.BlockNumber,
		BlockHash:      strings.ToLower(log.BlockHash.Hex()),
		LogIndex:       uint32(log.Index),
		TxIndex:        uint32(log.TxIndex),
		TxHash:         strings.ToLower(log.TxHash.Hex()),
		BlockTimestamp: time.Unix(int64(log.BlockTimestamp), 0),
		EventData:      nil,
	}

	var requestPriceLog *types.Log
	for i := len(allLogs) - 1; i >= 0; i-- {
		candidateLog := allLogs[i]
		if int64(candidateLog.Index) > int64(log.Index) {
			continue
		}
		if len(candidateLog.Topics) > 0 && candidateLog.Topics[0] == models.UmaOptimisticOracleV2RequestPriceEventSelectorHash {
			requestPriceLog = candidateLog
			break
		}
	}
	if requestPriceLog == nil {
		return fmt.Errorf("request price log is not found during question reset tx")
	}

	requestPrice, err := c.polymarketParser.ParseUmaOptimisticOracleV2RequestPriceLog(requestPriceLog)
	if err != nil {
		return err
	}
	if requestPrice == nil {
		return fmt.Errorf("request price log is nil during question reset tx")
	}

	ts := requestPrice.Timestamp.Uint64()
	eventData := sharedmodels.PolymarketMarketEventDataAdapterReset{
		NewTimestamp: ts,
	}
	event.EventData = &eventData
	questionID := strings.ToLower(questionReset.QuestionID.Hex())
	c.setQuestionTimestamp(questionID, ts)

	if c.notifier != nil {
		if err := notifier.Publish(c.notifier, context.Background(), string(sharedmodels.NotifierEventAdapterUMA), event); err != nil {
			c.notifier.Logger.Error("failed to publish adapter reset event", zap.Error(err))
		}
	}
	if err := c.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
		return c.insertPolymarketMarketEvent(ctx, q, event)
	}); err != nil {
		return err
	}

	err = c.processUmaOptimisticOracleV2RequestPriceLog(requestPriceLog, allLogs, metadata)
	if err != nil {
		return err
	}

	return nil
}

func (c *Processor) processUmaCtfAdapterQuestionResolvedLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	questionResolved, err := c.polymarketParser.ParseUmaCtfAdapterQuestionResolvedLog(log)
	if err != nil {
		return err
	}

	if questionResolved == nil {
		c.logger.Warn("uma ctf adapter question resolved log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	payoutsStr := make([]string, len(questionResolved.Payouts))
	for i, p := range questionResolved.Payouts {
		payoutsStr[i] = p.String()
	}

	eventData := sharedmodels.PolymarketMarketEventDataAdapterResolved{
		SettledPrice: questionResolved.SettledPrice.String(),
		Payouts:      payoutsStr,
	}

	event := sharedmodels.PolymarketMarketEvent{
		QuestionID:     strings.ToLower(questionResolved.QuestionID.Hex()),
		EventType:      sharedmodels.PolymarketEventAdapterResolved,
		BlockNumber:    log.BlockNumber,
		BlockHash:      strings.ToLower(log.BlockHash.Hex()),
		LogIndex:       uint32(log.Index),
		TxIndex:        uint32(log.TxIndex),
		TxHash:         strings.ToLower(log.TxHash.Hex()),
		BlockTimestamp: time.Unix(int64(log.BlockTimestamp), 0),
		EventData:      &eventData,
	}

	if c.notifier != nil {
		if err := notifier.Publish(c.notifier, context.Background(), string(sharedmodels.NotifierEventAdapterUMA), event); err != nil {
			c.notifier.Logger.Error("failed to publish adapter resolved event", zap.Error(err))
		}
	}
	if err := c.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
		return c.insertPolymarketMarketEvent(ctx, q, event)
	}); err != nil {
		return err
	}

	questionID := strings.ToLower(questionResolved.QuestionID.Hex())
	c.clearQuestionTimestamp(questionID)

	return nil
}

func (c *Processor) processUmaCtfAdapterQuestionPausedLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	questionPaused, err := c.polymarketParser.ParseUmaCtfAdapterQuestionPausedLog(log)
	if err != nil {
		return err
	}

	if questionPaused == nil {
		c.logger.Warn("uma ctf adapter question paused log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	event := sharedmodels.PolymarketMarketEvent{
		QuestionID:     strings.ToLower(questionPaused.QuestionID.Hex()),
		EventType:      sharedmodels.PolymarketEventAdapterPaused,
		BlockNumber:    log.BlockNumber,
		BlockHash:      strings.ToLower(log.BlockHash.Hex()),
		LogIndex:       uint32(log.Index),
		TxIndex:        uint32(log.TxIndex),
		TxHash:         strings.ToLower(log.TxHash.Hex()),
		BlockTimestamp: time.Unix(int64(log.BlockTimestamp), 0),
		EventData:      nil,
	}

	if c.notifier != nil {
		if err := notifier.Publish(c.notifier, context.Background(), string(sharedmodels.NotifierEventAdapterUMA), event); err != nil {
			c.notifier.Logger.Error("failed to publish adapter paused event", zap.Error(err))
		}
	}
	return c.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
		return c.insertPolymarketMarketEvent(ctx, q, event)
	})
}

func (c *Processor) processUmaCtfAdapterQuestionFlaggedLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	questionFlagged, err := c.polymarketParser.ParseUmaCtfAdapterQuestionFlaggedLog(log)
	if err != nil {
		return err
	}

	if questionFlagged == nil {
		c.logger.Warn("uma ctf adapter question flagged log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	event := sharedmodels.PolymarketMarketEvent{
		QuestionID:     strings.ToLower(questionFlagged.QuestionID.Hex()),
		EventType:      sharedmodels.PolymarketEventAdapterFlagged,
		BlockNumber:    log.BlockNumber,
		BlockHash:      strings.ToLower(log.BlockHash.Hex()),
		LogIndex:       uint32(log.Index),
		TxIndex:        uint32(log.TxIndex),
		TxHash:         strings.ToLower(log.TxHash.Hex()),
		BlockTimestamp: time.Unix(int64(log.BlockTimestamp), 0),
		EventData:      nil,
	}

	if c.notifier != nil {
		if err := notifier.Publish(c.notifier, context.Background(), string(sharedmodels.NotifierEventAdapterUMA), event); err != nil {
			c.notifier.Logger.Error("failed to publish adapter flagged event", zap.Error(err))
		}
	}
	return c.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
		return c.insertPolymarketMarketEvent(ctx, q, event)
	})
}
