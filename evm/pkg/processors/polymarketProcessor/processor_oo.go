package polymarketProcessor

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/shared/db"
	sharedmodels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/monolit-network/xlayer-indexer/util/notifier"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

func (c *Processor) shouldProcessUmaOoLog(requester common.Address, questionID common.Hash, timestamp *big.Int, blockNumber uint64) (bool, error) {
	if _, ok := models.PolymarketCTFAdaptersAddressesMap[requester]; !ok {
		return false, nil
	}
	if timestamp == nil {
		return false, nil
	}
	currentTimestamp, err := c.getQuestionTimestampByQuestionID(questionID.Hex(), blockNumber)
	if err != nil {
		return false, err
	}
	if currentTimestamp != timestamp.Uint64() {
		return false, nil
	}
	return true, nil
}

func (c *Processor) processUmaOptimisticOracleV2ProposePriceLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	proposePrice, err := c.polymarketParser.ParseUmaOptimisticOracleV2ProposePriceLog(log)
	if err != nil {
		return err
	}

	if proposePrice == nil {
		c.logger.Warn("uma optimistic oracle v2 propose price log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	shouldProcess, err := c.shouldProcessUmaOoLog(proposePrice.Requester, proposePrice.QuestionID, proposePrice.Timestamp, log.BlockNumber)
	if err != nil {
		return err
	}
	if !shouldProcess {
		return nil
	}

	eventData := sharedmodels.PolymarketMarketEventDataUmaProposed{
		Requester:           strings.ToLower(proposePrice.Requester.Hex()),
		Proposer:            strings.ToLower(proposePrice.Proposer.Hex()),
		Identifier:          strings.ToLower(proposePrice.Identifier.Hex()),
		Timestamp:           proposePrice.Timestamp.Uint64(),
		ProposedPrice:       proposePrice.ProposedPrice.String(),
		ExpirationTimestamp: proposePrice.ExpirationTimestamp.Uint64(),
		Currency:            strings.ToLower(proposePrice.Currency.Hex()),
	}

	event := sharedmodels.PolymarketMarketEvent{
		QuestionID:     strings.ToLower(proposePrice.QuestionID.Hex()),
		EventType:      sharedmodels.PolymarketEventUmaProposed,
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
			c.notifier.Logger.Error("failed to publish adapter uma proposed event", zap.Error(err))
		}
	}
	return c.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
		return c.insertPolymarketMarketEvent(ctx, q, event)
	})
}

func (c *Processor) processUmaOptimisticOracleV2DisputePriceLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	disputePrice, err := c.polymarketParser.ParseUmaOptimisticOracleV2DisputePriceLog(log)
	if err != nil {
		return err
	}

	if disputePrice == nil {
		c.logger.Warn("uma optimistic oracle v2 dispute price log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	shouldProcess, err := c.shouldProcessUmaOoLog(disputePrice.Requester, disputePrice.QuestionID, disputePrice.Timestamp, log.BlockNumber)
	if err != nil {
		return err
	}
	if !shouldProcess {
		return nil
	}

	eventData := sharedmodels.PolymarketMarketEventDataUmaDisputed{
		Requester:     strings.ToLower(disputePrice.Requester.Hex()),
		Proposer:      strings.ToLower(disputePrice.Proposer.Hex()),
		Disputer:      strings.ToLower(disputePrice.Disputer.Hex()),
		Identifier:    strings.ToLower(disputePrice.Identifier.Hex()),
		Timestamp:     disputePrice.Timestamp.Uint64(),
		ProposedPrice: disputePrice.ProposedPrice.String(),
	}

	event := sharedmodels.PolymarketMarketEvent{
		QuestionID:     strings.ToLower(disputePrice.QuestionID.Hex()),
		EventType:      sharedmodels.PolymarketEventUmaDisputed,
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
			c.notifier.Logger.Error("failed to publish adapter uma disputed event", zap.Error(err))
		}
	}
	return c.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
		return c.insertPolymarketMarketEvent(ctx, q, event)
	})
}

func (c *Processor) processUmaOptimisticOracleV2SettleLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	settle, err := c.polymarketParser.ParseUmaOptimisticOracleV2SettleLog(log)
	if err != nil {
		return err
	}

	if settle == nil {
		c.logger.Warn("uma optimistic oracle v2 settle log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	shouldProcess, err := c.shouldProcessUmaOoLog(settle.Requester, settle.QuestionID, settle.Timestamp, log.BlockNumber)
	if err != nil {
		return err
	}
	if !shouldProcess {
		return nil
	}

	eventData := sharedmodels.PolymarketMarketEventDataUmaSettle{
		Requester:  strings.ToLower(settle.Requester.Hex()),
		Proposer:   strings.ToLower(settle.Proposer.Hex()),
		Disputer:   strings.ToLower(settle.Disputer.Hex()),
		Identifier: strings.ToLower(settle.Identifier.Hex()),
		Timestamp:  settle.Timestamp.Uint64(),
		Price:      settle.Price.String(),
		Payout:     settle.Payout.String(),
	}

	eventName := sharedmodels.PolymarketEventUmaSettle
	if eventData.Price == models.UmaOptimisticOracleV2InvalidPrice.String() {
		eventName = sharedmodels.PolymarketEventUmaSettleInvalid
	}

	event := sharedmodels.PolymarketMarketEvent{
		QuestionID:     strings.ToLower(settle.QuestionID.Hex()),
		EventType:      eventName,
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
			c.notifier.Logger.Error("failed to publish adapter uma settle event", zap.Error(err))
		}
	}
	return c.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
		return c.insertPolymarketMarketEvent(ctx, q, event)
	})
}

func (c *Processor) processUmaOptimisticOracleV2RequestPriceLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	requestPrice, err := c.polymarketParser.ParseUmaOptimisticOracleV2RequestPriceLog(log)
	if err != nil {
		return err
	}

	if requestPrice == nil {
		c.logger.Warn("uma optimistic oracle v2 request price log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	shouldProcess, err := c.shouldProcessUmaOoLog(requestPrice.Requester, requestPrice.QuestionID, requestPrice.Timestamp, log.BlockNumber)
	if err != nil {
		return err
	}
	if !shouldProcess {
		return nil
	}

	eventData := sharedmodels.PolymarketMarketEventDataUmaRequested{
		Requester:  strings.ToLower(requestPrice.Requester.Hex()),
		Identifier: strings.ToLower(requestPrice.Identifier.Hex()),
		Timestamp:  requestPrice.Timestamp.Uint64(),
		Currency:   strings.ToLower(requestPrice.Currency.Hex()),
		Reward:     requestPrice.Reward.String(),
		FinalFee:   requestPrice.FinalFee.String(),
	}

	event := sharedmodels.PolymarketMarketEvent{
		QuestionID:     strings.ToLower(requestPrice.QuestionID.Hex()),
		EventType:      sharedmodels.PolymarketEventUmaRequested,
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
			c.notifier.Logger.Error("failed to publish adapter uma requested event", zap.Error(err))
		}
	}
	return c.dbClient.WithRetry(context.Background(), -1, func(ctx context.Context, q db.DB) error {
		return c.insertPolymarketMarketEvent(ctx, q, event)
	})
}
