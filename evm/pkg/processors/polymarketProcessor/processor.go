package polymarketProcessor

import (
	"context"
	"fmt"
	"math/big"
	"runtime/debug"
	"slices"
	"sync"
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	evmmetrics "github.com/monolit-network/xlayer-indexer/evm/pkg/metrics"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers/polymarket"
	"github.com/monolit-network/xlayer-indexer/shared/db"
	sharedmodels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/monolit-network/xlayer-indexer/util/notifier"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru/v2"
	"go.uber.org/zap"
)

var polymarketContractsTopics = map[common.Hash]struct{}{
	models.PolymarketFeeRefundedEventSelectorHash:        struct{}{},
	models.PolymarketNegRiskFeeRefundedEventSelectorHash: struct{}{},
	// models.PolymarketAMMApprovalEventSelectorHash:        struct{}{},
	// models.PolymarketAMMTransferEventSelectorHash:        struct{}{},
	// models.PolymarketFPMMFundingAddedEventSelectorHash:   struct{}{},
	// models.PolymarketFPMMFundingRemovedEventSelectorHash: struct{}{},
	// models.PolymarketFPMMBuyEventSelectorHash:            struct{}{},
	// models.PolymarketFPMMSellEventSelectorHash:           struct{}{},
}

var polymarketParseRangeAddressFilteredAddresses = func() map[common.Address]struct{} {
	addresses := map[common.Address]struct{}{
		models.PolymarketCTFExchangeAddress:                  {},
		models.PolymarketCTFExchangeAddressV2:                {},
		models.PolymarketNegRiskCTFExchangeAddress:           {},
		models.PolymarketNegRiskCTFExchangeAddressV2:         {},
		models.PolymarketNegRiskAdapterAddress:               {},
		models.PolymarketConditionalTokensAddress:            {},
		models.UmaOptimisticOracleV2Address:                  {},
		models.UmaPolymarketManagedOptimisticOracleV2Address: {},
	}

	for _, addr := range models.PolymarketCTFAdaptersAddresses {
		addresses[addr] = struct{}{}
	}

	return addresses
}()

var polymarketParseRangeAddressFilteredTopics = map[common.Hash]struct{}{
	models.PolymarketOrderFilledEventSelectorHash:                      {},
	models.PolymarketOrderFilledEventV2SelectorHash:                    {},
	models.PolymarketOrdersMatchedEventV2SelectorHash:                  {},
	models.PolymarketTransferSingleEventSelectorHash:                   {},
	models.PolymarketTransferBatchEventSelectorHash:                    {},
	models.PolymarketConditionPreparationEventSelectorHash:             {},
	models.PolymarketConditionResolutionEventSelectorHash:              {},
	models.PolymarketPositionSplitEventSelectorHash:                    {},
	models.PolymarketPositionMergeEventSelectorHash:                    {},
	models.PolymarketPayoutRedemptionEventSelectorHash:                 {},
	models.PolymarketNegRiskPayoutRedemptionEventSelectorHash:          {},
	models.PolymarketUmaCtfAdapterQuestionInitializedEventSelectorHash: {},
	models.PolymarketUmaCtfAdapterQuestionResetEventSelectorHash:       {},
	models.PolymarketUmaCtfAdapterQuestionResolvedEventSelectorHash:    {},
	models.PolymarketUmaCtfAdapterQuestionPausedEventSelectorHash:      {},
	models.PolymarketUmaCtfAdapterQuestionFlaggedEventSelectorHash:     {},
	models.UmaOptimisticOracleV2RequestPriceEventSelectorHash:          {},
	models.UmaOptimisticOracleV2ProposePriceEventSelectorHash:          {},
	models.UmaOptimisticOracleV2DisputePriceEventSelectorHash:          {},
	models.UmaOptimisticOracleV2SettleEventSelectorHash:                {},
}

var polymarketParseRangeNoAddressTopics = map[common.Hash]struct{}{
	models.PolymarketFeeRefundedEventSelectorHash:               {},
	models.PolymarketNegRiskFeeRefundedEventSelectorHash:        {},
	models.PolymarketNegRiskPositionsConvertedEventSelectorHash: {},
	models.PolymarketNegRiskPositionSplitEventSelectorHash:      {},
	models.PolymarketNegRiskPositionMergeEventSelectorHash:      {},
}

type Processor struct {
	clickClient db.DBClick
	dbClient    db.Client

	notifier *notifier.Notifier

	evmClient        *evmclient.Client
	polymarketParser *polymarket.CTFParser

	logger *zap.Logger

	lastSubIndexCache     *lru.Cache[common.Hash, uint32]
	tokenToConditionCache *lru.Cache[string, common.Hash]

	processors map[string]logProcessor

	processStatsMu sync.Mutex
	processStats   map[string]*processTiming

	questionTimestampMu    sync.RWMutex
	questionTimestampCache map[string]uint64

	tokenMappingMu    sync.RWMutex
	conditionTokenIDs map[string]map[string]struct{}

	tokenIDToPayoutMu sync.RWMutex
	tokenIDToPayout   map[string]tokenIDToPayout

	conditionQuestionMappingMu sync.RWMutex
	conditionQuestionMapping   *lru.Cache[string, string]
	questionConditionMapping   *lru.Cache[string, string]

	errorsCh  chan error
	requests  chan parsers.ParseTxRequest
	workersWg *sync.WaitGroup
	writerWg  *sync.WaitGroup
	eventsCh  chan *models.PolymarketOrderEventNew

	pgWritersWg *sync.WaitGroup
	pgWritersCh chan *sharedmodels.PolymarketToken

	dbFlushCh       chan chan struct{}
	pgWriterFlushCh chan chan struct{}

	settings settings
	chain    string
}

type tokenIDToPayout struct {
	Numerator   *big.Int
	Denominator *big.Int
}

type settings struct {
	disableBuffer           bool
	logTimingsInterval      *time.Duration
	disableClickhouseWriter bool
}

type ProcessorOption func(*settings)

func WithDisableBuffer(disableBuffer bool) ProcessorOption {
	return func(s *settings) {
		s.disableBuffer = disableBuffer
	}
}

func WithLogTimingsInterval(interval time.Duration) ProcessorOption {
	return func(s *settings) {
		s.logTimingsInterval = &interval
	}
}

func WithDisableClickhouseWriter(disableClickhouseWriter bool) ProcessorOption {
	return func(s *settings) {
		s.disableClickhouseWriter = disableClickhouseWriter
	}
}

type logProcessor struct {
	name string
	fn   func(log *types.Log, allLogs []*types.Log, metadata map[string]any) error
}

type processTiming struct {
	total time.Duration
	count uint64

	currentTotal time.Duration
	currentCount uint64
}

func (c *Processor) RequiredDataTypes() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesReceipt
}

func NewProcessor(
	evmClient *evmclient.Client,
	clickClient db.DBClick,
	notif *notifier.Notifier,
	dbClient db.Client,
	logger *zap.Logger,
	options ...ProcessorOption,
) *Processor {
	subIndexCache, _ := lru.New[common.Hash, uint32](2048)
	tokenToConditionCache, _ := lru.New[string, common.Hash](65536)
	conditionQuestionMapping, _ := lru.New[string, string](65536)
	questionConditionMapping, _ := lru.New[string, string](65536)

	settings := &settings{
		disableBuffer: false,
	}
	for _, option := range options {
		option(settings)
	}

	c := &Processor{
		evmClient:             evmClient,
		clickClient:           clickClient,
		dbClient:              dbClient,
		notifier:              notif,
		polymarketParser:      polymarket.NewCTFParser(evmClient, logger),
		logger:                logger,
		lastSubIndexCache:     subIndexCache,
		tokenToConditionCache: tokenToConditionCache,
		eventsCh:              make(chan *models.PolymarketOrderEventNew, 32768),

		processStats:             map[string]*processTiming{},
		questionTimestampCache:   map[string]uint64{},
		conditionTokenIDs:        map[string]map[string]struct{}{},
		tokenIDToPayout:          map[string]tokenIDToPayout{},
		conditionQuestionMapping: conditionQuestionMapping,
		questionConditionMapping: questionConditionMapping,

		requests:        make(chan parsers.ParseTxRequest, 1024),
		errorsCh:        make(chan error, 1024),
		pgWritersCh:     make(chan *sharedmodels.PolymarketToken),
		workersWg:       &sync.WaitGroup{},
		writerWg:        &sync.WaitGroup{},
		pgWritersWg:     &sync.WaitGroup{},
		dbFlushCh:       make(chan chan struct{}),
		pgWriterFlushCh: make(chan chan struct{}),
		settings:        *settings,
		chain:           string(models.ChainPolygon),
	}

	c.processors = map[string]logProcessor{
		models.PolymarketOrderFilledEventSelectorHash.Hex(): {
			name: "order_filled",
			fn:   c.processOrderFilledLog,
		},
		models.PolymarketOrderFilledEventV2SelectorHash.Hex(): {
			name: "order_filled_v2",
			fn:   c.processOrderFilledLog,
		},
		models.PolymarketFeeRefundedEventSelectorHash.Hex(): {
			name: "fee_refunded",
			fn:   c.processFeeRefundedLog,
		},
		models.PolymarketNegRiskFeeRefundedEventSelectorHash.Hex(): {
			name: "neg_risk_fee_refunded",
			fn:   c.processFeeRefundedLog,
		},
		models.PolymarketNegRiskPositionsConvertedEventSelectorHash.Hex(): {
			name: "neg_risk_positions_converted",
			fn:   c.processNegRiskPositionConvertedLog,
		},
		models.PolymarketPositionSplitEventSelectorHash.Hex(): {
			name: "position_split",
			fn:   c.processPositionSplitLog,
		},
		models.PolymarketNegRiskPositionSplitEventSelectorHash.Hex(): {
			name: "position_split_neg_risk",
			fn:   c.processNegRiskPositionSplitLog,
		},
		models.PolymarketPositionMergeEventSelectorHash.Hex(): {
			name: "position_merge",
			fn:   c.processPositionMergeLog,
		},
		models.PolymarketNegRiskPositionMergeEventSelectorHash.Hex(): {
			name: "position_merge_neg_risk",
			fn:   c.processNegRiskPositionMergeLog,
		},
		models.PolymarketPayoutRedemptionEventSelectorHash.Hex(): {
			name: "position_redeem",
			fn:   c.processPositionRedeemLog,
		},
		models.PolymarketNegRiskPayoutRedemptionEventSelectorHash.Hex(): {
			name: "payout_redeem_neg_risk",
			fn:   c.processNegRiskPayoutRedemptionLog,
		},
		models.PolymarketTransferSingleEventSelectorHash.Hex(): {
			name: "transfer_single",
			fn:   c.processTransferSingleLog,
		},
		models.PolymarketTransferBatchEventSelectorHash.Hex(): {
			name: "transfer_batch",
			fn:   c.processTransferBatchLog,
		},
		models.PolymarketConditionPreparationEventSelectorHash.Hex(): {
			name: "condition_preparation",
			fn:   c.processConditionPreparationLog,
		},
		models.PolymarketConditionResolutionEventSelectorHash.Hex(): {
			name: "condition_resolution",
			fn:   c.processConditionResolutionLog,
		},
		// UMA CTF Adapter events
		models.PolymarketUmaCtfAdapterQuestionInitializedEventSelectorHash.Hex(): {
			name: "adapter_question_initialized",
			fn:   c.processUmaCtfAdapterQuestionInitializedLog,
		},
		models.PolymarketUmaCtfAdapterQuestionResetEventSelectorHash.Hex(): {
			name: "adapter_question_reset",
			fn:   c.processUmaCtfAdapterQuestionResetLog,
		},
		models.PolymarketUmaCtfAdapterQuestionResolvedEventSelectorHash.Hex(): {
			name: "adapter_question_resolved",
			fn:   c.processUmaCtfAdapterQuestionResolvedLog,
		},
		models.PolymarketUmaCtfAdapterQuestionPausedEventSelectorHash.Hex(): {
			name: "adapter_question_paused",
			fn:   c.processUmaCtfAdapterQuestionPausedLog,
		},
		models.PolymarketUmaCtfAdapterQuestionFlaggedEventSelectorHash.Hex(): {
			name: "adapter_question_flagged",
			fn:   c.processUmaCtfAdapterQuestionFlaggedLog,
		},
		// UMA Optimistic Oracle V2 events
		models.UmaOptimisticOracleV2RequestPriceEventSelectorHash.Hex(): {
			name: "uma_request_price",
			fn:   c.processUmaOptimisticOracleV2RequestPriceLog,
		},
		models.UmaOptimisticOracleV2ProposePriceEventSelectorHash.Hex(): {
			name: "uma_propose_price",
			fn:   c.processUmaOptimisticOracleV2ProposePriceLog,
		},
		models.UmaOptimisticOracleV2DisputePriceEventSelectorHash.Hex(): {
			name: "uma_dispute_price",
			fn:   c.processUmaOptimisticOracleV2DisputePriceLog,
		},
		models.UmaOptimisticOracleV2SettleEventSelectorHash.Hex(): {
			name: "uma_settle",
			fn:   c.processUmaOptimisticOracleV2SettleLog,
		},
	}

	c.writerWg.Go(c.dbWriter)

	c.workersWg.Go(func() {
		for request := range c.requests {
			if err := c.ProcessLogsFromTx(request.Receipt.Logs); err != nil {
				c.errorsCh <- err
				continue
			}
		}
	})

	c.pgWritersWg.Go(c.tokenMappingWriter)

	if settings.logTimingsInterval != nil {
		go func() {
			ticker := time.NewTicker(*settings.logTimingsInterval)
			defer ticker.Stop()
			for range ticker.C {
				c.LogTimings()
			}
		}()
	}

	return c
}

func (c *Processor) Errors() chan error {
	return c.errorsCh
}

func (c *Processor) LogTimings() {
	c.processStatsMu.Lock()
	defer c.processStatsMu.Unlock()

	totalProcessLogsFromTx := time.Duration(0)
	totalProcessLogsFromTxCount := uint64(0)
	total := time.Duration(0)
	totalCount := uint64(0)
	for name, stats := range c.processStats {
		if stats.count == 0 {
			continue
		}
		if name != "process_logs_from_tx" {
			total += stats.total
			totalCount += stats.count
			c.logger.Info("log timing", zap.String("name", name), zap.Duration("total", stats.total), zap.Uint64("count", stats.count), zap.Duration("avg", stats.total/time.Duration(stats.count)))
		} else {
			totalProcessLogsFromTx += stats.total
			totalProcessLogsFromTxCount += stats.count
		}
	}
	if totalCount == 0 {
		return
	}
	c.logger.Info("log timing total",
		zap.Duration("total", total),
		zap.Uint64("count", totalCount),
		zap.Duration("avg", total/time.Duration(totalCount)),
		zap.Duration("avg_process_logs_from_tx", totalProcessLogsFromTx/time.Duration(totalProcessLogsFromTxCount)),
		zap.Duration("total_process_logs_from_tx", totalProcessLogsFromTx),
		zap.Uint64("count_process_logs_from_tx", totalProcessLogsFromTxCount),
	)
}

func (c *Processor) LogTimingsCurrent() {
	c.processStatsMu.Lock()
	defer c.processStatsMu.Unlock()

	totalProcessLogsFromTx := time.Duration(0)
	totalProcessLogsFromTxCount := uint64(0)
	total := time.Duration(0)
	totalCount := uint64(0)
	for name, stats := range c.processStats {
		if stats.currentCount == 0 {
			continue
		}
		if name != "process_logs_from_tx" {
			total += stats.currentTotal
			totalCount += stats.currentCount
			c.logger.Info("log timing", zap.String("name", name), zap.Duration("total", stats.currentTotal), zap.Uint64("count", stats.currentCount), zap.Duration("avg", stats.currentTotal/time.Duration(stats.currentCount)))
		} else {
			totalProcessLogsFromTx += stats.currentTotal
			totalProcessLogsFromTxCount += stats.currentCount
		}
		stats.currentTotal = 0
		stats.currentCount = 0
	}
	if totalCount == 0 {
		return
	}
	var avg time.Duration
	var avgProcessLogsFromTx time.Duration
	if totalCount > 0 {
		avg = total / time.Duration(totalCount)
	}
	if totalProcessLogsFromTxCount > 0 {
		avgProcessLogsFromTx = totalProcessLogsFromTx / time.Duration(totalProcessLogsFromTxCount)
	}
	c.logger.Info("log timing current",
		zap.Duration("total", total),
		zap.Uint64("count", totalCount),
		zap.Duration("avg", avg),
		zap.Duration("avg_process_logs_from_tx", avgProcessLogsFromTx),
		zap.Duration("total_process_logs_from_tx", totalProcessLogsFromTx),
		zap.Uint64("count_process_logs_from_tx", totalProcessLogsFromTxCount),
	)
}

func (c *Processor) Stop() {
	close(c.requests)
	c.workersWg.Wait()
	close(c.eventsCh)
	c.writerWg.Wait()
	close(c.errorsCh)

	close(c.pgWritersCh)
	c.pgWritersWg.Wait()
}

func basicEventFromLog(log *types.Log) models.PolymarketOrderEventNew {
	return models.PolymarketOrderEventNew{
		BlockTime:   time.Unix(int64(log.BlockTimestamp), 0),
		BlockNumber: new(big.Int).SetUint64(log.BlockNumber),
		BlockHash:   log.BlockHash,
		TxIdx:       uint32(log.TxIndex),
		TxHash:      log.TxHash,
		LogIndex:    uint32(log.Index),
	}
}

func (c *Processor) ProcessLog(log *types.Log, txLogs []*types.Log) error {
	return c.processLog(log, txLogs, map[string]any{})
}

func (c *Processor) processLog(log *types.Log, txLogs []*types.Log, metadata map[string]any) error {
	// if log.Removed {
	// 	return c.rollbackLog(ctx, log)
	// }
	if len(log.Topics) == 0 {
		return nil
	}

	processor, ok := c.processors[log.Topics[0].Hex()]
	if !ok {
		// c.logger.Warn("no processor found for log topic", zap.String("topic", log.Topics[0].Hex()))
		return nil
	}

	start := time.Now()
	err := processor.fn(log, txLogs, metadata)
	c.recordProcessDuration(processor.name, time.Since(start))

	if err != nil {
		return fmt.Errorf("error processing log (processor: %s): %w", processor.name, err)
	}

	if r := recover(); r != nil {
		c.logger.Error("panic processing log",
			zap.String("topic", log.Topics[0].Hex()),
			zap.Any("recover", r),
			zap.String("tx_hash", log.TxHash.Hex()),
			zap.Uint64("log_index", uint64(log.Index)),
			zap.Uint64("tx_index", uint64(log.TxIndex)),
		)
		if stack := debug.Stack(); stack != nil {
			c.logger.Error("stack trace", zap.ByteString("stack", stack))
		}
		return fmt.Errorf("panic processing log: %s | %s | %d", log.TxHash.Hex(), log.Topics[0].Hex(), log.Index)
	}
	return nil
}

func (c *Processor) ProcessLogsFromTx(logs []*types.Log) error {
	start := time.Now()
	metadata := map[string]any{}

	filteredLogs, err := c.preProcessLogGroup(logs, metadata)
	if err != nil {
		return fmt.Errorf("error pre processing log group: %w", err)
	}

	for _, log := range filteredLogs {
		err := c.processLog(log, filteredLogs, metadata)
		if err != nil {
			for i, log := range filteredLogs {
				c.logger.Info("error tx trace", zap.String("tx_hash", log.TxHash.Hex()), zap.String("topic", log.Topics[0].Hex()), zap.Int("i", i), zap.Uint64("log_index", uint64(log.Index)), zap.Uint64("tx_index", uint64(log.TxIndex)))
			}
			return fmt.Errorf("error processing log (tx_size: %d): %w", len(logs), err)
		}
	}

	if err := c.postProcessLogGroup(filteredLogs, filteredLogs, metadata); err != nil {
		return fmt.Errorf("error post processing log group: %w", err)
	}

	c.recordProcessDuration("process_logs_from_tx", time.Since(start))
	return nil
}

func (c *Processor) QueueRequest(request parsers.ParseTxRequest) {
	c.requests <- request
}

func (c *Processor) ProcessRequest(request parsers.ParseTxRequest) error {
	if err := c.ProcessLogsFromTx(request.Receipt.Logs); err != nil {
		return err
	}
	return nil
}

func (c *Processor) preProcessLogGroup(logs []*types.Log, metadata map[string]any) ([]*types.Log, error) {
	filteredLogs := make([]*types.Log, 0, len(logs))
	for _, log := range logs {
		if !isLogAllowedByPolymarketParseRangeRequirements(log) {
			continue
		}
		filteredLogs = append(filteredLogs, log)
	}

	polymarketAddrs := c.collectPolymarketContractAddrs(filteredLogs)
	if len(polymarketAddrs) > 0 {
		metadata["polymarket_contract_addrs"] = polymarketAddrs
	}

	// hasFPMMLog := false
	// for _, log := range logs {
	// 	if c.polymarketParser.IsFPMMLog(log) {
	// 		// c.logger.Info("skipping fpmm log", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
	// 		hasFPMMLog = true
	// 		break
	// 	}
	// }

	// if hasFPMMLog {
	// 	// split logs are essential for mapping tokens to conditions
	// 	splitLogs := make([]*types.Log, 0)
	// 	for _, log := range filteredLogs {
	// 		if len(log.Topics) > 0 && log.Topics[0] == models.PolymarketPositionSplitEventSelectorHash {
	// 			splitLogs = append(splitLogs, log)
	// 		}
	// 	}
	// 	filteredLogs = splitLogs
	// }

	filteredLogs = sortLogsByPriority(filteredLogs)

	return filteredLogs, nil
}

func (c *Processor) collectPolymarketContractAddrs(allLogs []*types.Log) map[common.Address]struct{} {
	polymarketAddrs := make(map[common.Address]struct{})
	for _, log := range allLogs {
		if len(log.Topics) == 0 {
			continue
		}
		if _, ok := polymarketContractsTopics[log.Topics[0]]; ok {
			polymarketAddrs[log.Address] = struct{}{}
		}
	}
	polymarketAddrs[common.HexToAddress("0xa5Ef39C3D3e10d0B270233af41CaC69796B12966")] = struct{}{} // neg risk adapter NO_TOKENS_BURN_ADDRESS

	return polymarketAddrs
}

func (c *Processor) isPolymarketContractAddress(addr common.Address, metadata map[string]any) bool {
	if c.polymarketParser.IsPolymarketContractAddress(addr) {
		return true
	}
	if polymarketAddrs, ok := metadata["polymarket_contract_addrs"].(map[common.Address]struct{}); ok {
		if _, ok := polymarketAddrs[addr]; ok {
			return true
		}
	}
	return false
}

func (c *Processor) postProcessLogGroup(originalLogs, modifiedLogs []*types.Log, metadata map[string]any) error {
	if orders, ok := metadata["current_tx_orders_filled"]; ok {
		orders := orders.([]*models.PolymarketOrderEventNew)
		if len(orders) > 0 {
			var refunds []*sharedmodels.PolymarketFeeRefunded
			if refundsRaw, ok := metadata["current_tx_fee_refunds"]; ok {
				refunds = refundsRaw.([]*sharedmodels.PolymarketFeeRefunded)
			}
			if err := c.postProcessOrderFilledLog(orders, refunds); err != nil {
				return fmt.Errorf("error post processing order filled log: %w", err)
			}
		}
	}

	// if err := c.processFPMMLogs(originalLogs); err != nil {
	// 	return fmt.Errorf("error processing fpmm logs: %w", err)
	// }

	return nil
}

var logsPriority = map[common.Hash]int{
	models.PolymarketConditionPreparationEventSelectorHash:             -3,
	models.PolymarketUmaCtfAdapterQuestionInitializedEventSelectorHash: -2,
	models.PolymarketPositionSplitEventSelectorHash:                    -2,
	models.PolymarketFeeRefundedEventSelectorHash:                      -1,
	models.PolymarketNegRiskFeeRefundedEventSelectorHash:               -1,
}

func sortLogsByPriority(logs []*types.Log) []*types.Log {
	slices.SortStableFunc(logs, func(a, b *types.Log) int {
		return logPriority(a) - logPriority(b)
	})
	return logs
}

func logPriority(log *types.Log) int {
	if len(log.Topics) == 0 {
		return 0
	}
	return logsPriority[log.Topics[0]]
}

func isLogAllowedByPolymarketParseRangeRequirements(log *types.Log) bool {
	if len(log.Topics) == 0 {
		return false
	}

	topic0 := log.Topics[0]
	if _, ok := polymarketParseRangeNoAddressTopics[topic0]; ok {
		return true
	}
	if _, ok := polymarketParseRangeAddressFilteredTopics[topic0]; !ok {
		return false
	}

	_, ok := polymarketParseRangeAddressFilteredAddresses[log.Address]
	return ok
}

func (c *Processor) recordProcessDuration(name string, duration time.Duration) {
	c.processStatsMu.Lock()
	defer c.processStatsMu.Unlock()

	stats, ok := c.processStats[name]
	if !ok {
		stats = &processTiming{}
		c.processStats[name] = stats
	}
	stats.total += duration
	stats.currentTotal += duration
	stats.count++
	stats.currentCount++
}

func (c *Processor) insertPolymarketMarketEvent(ctx context.Context, q db.DB, event sharedmodels.PolymarketMarketEvent) error {
	if err := q.InsertPolymarketMarketEvent(ctx, event); err != nil {
		c.recordPolymarketMarketEventInsertError()
		return err
	}
	c.recordPolymarketMarketEventInserted()
	return nil
}

func (c *Processor) recordPolymarketMarketEventInserted() {
	evmmetrics.InsertedRows.WithLabelValues(c.chain, "polymarket", "postgres", "polymarket.market_events").Inc()
}

func (c *Processor) recordPolymarketMarketEventInsertError() {
	evmmetrics.InsertErrors.WithLabelValues(c.chain, "polymarket", "postgres", "polymarket.market_events").Inc()
}

func (c *Processor) completeEvent(event *models.PolymarketOrderEventNew) error {
	if event.Source == "" {
		return fmt.Errorf("source is required")
	}

	if event.TokenID == nil {
		return fmt.Errorf("token id is required")
	}

	if event.ConditionID == (common.Hash{}) {
		conditionID, err := c.getConditionIDByTokenID(event.TokenID)
		if err != nil {
			return fmt.Errorf("error getting condition id by token id: %w", err)
		}
		event.ConditionID = conditionID
	}

	return nil
}

func (c *Processor) dbWriter() {
	ctx := context.Background()
	flush := func(batch []*models.PolymarketOrderEventNew) []*models.PolymarketOrderEventNew {
		if len(batch) == 0 {
			return batch
		}
		for {
			if c.settings.disableClickhouseWriter {
				break
			}
			err := c.clickClient.InsertEvmPolymarketOrderEventsBatchNew(ctx, batch)
			if err != nil {
				c.logger.Error("error inserting polymarket order events batch", zap.Error(err))
				evmmetrics.InsertErrors.WithLabelValues(c.chain, "polymarket", "clickhouse", "evm.polymarket_order_events").Inc()
				time.Sleep(5 * time.Second)
				continue
			}
			evmmetrics.InsertedRows.WithLabelValues(c.chain, "polymarket", "clickhouse", "evm.polymarket_order_events").Add(float64(len(batch)))
			break
		}
		return batch[:0]
	}

	batch := make([]*models.PolymarketOrderEventNew, 0, 4096)
	for {
		select {
		case event, ok := <-c.eventsCh:
			if !ok {
				flush(batch)
				return
			}
			if err := c.completeEvent(event); err != nil {
				c.logger.Error("error completing event", zap.Error(err), zap.String("tx_hash", event.TxHash.Hex()), zap.Uint64("log_index", uint64(event.LogIndex)))
				continue
			}
			if c.notifier != nil {
				if err := notifier.Publish(c.notifier, ctx, string(sharedmodels.NotifierEventOrder), *event); err != nil {
					c.notifier.Logger.Warn("failed to publish order event", zap.String("topic", string(sharedmodels.NotifierEventOrder)), zap.Error(err))
				}
			}
			batch = append(batch, event)
			if len(batch) >= 4096 {
				batch = flush(batch)
			}
		case done := <-c.dbFlushCh:
			batch = flush(batch)
			close(done)
		}
	}
}

func (c *Processor) FlushSync() {
	dbDone := make(chan struct{})
	pgDone := make(chan struct{})
	c.dbFlushCh <- dbDone
	c.pgWriterFlushCh <- pgDone
	<-dbDone
	<-pgDone
}

func (c *Processor) flushPgSync() {
	pgDone := make(chan struct{})
	c.pgWriterFlushCh <- pgDone
	<-pgDone
}
