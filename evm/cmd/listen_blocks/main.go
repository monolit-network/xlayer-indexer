package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/processors"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/processors/basicProcessor"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/processors/polymarketProcessor"
	"github.com/monolit-network/xlayer-indexer/shared/db"
	"github.com/monolit-network/xlayer-indexer/shared/db/clickhouse"
	"github.com/monolit-network/xlayer-indexer/shared/db/postgres"
	sharedmodels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/monolit-network/xlayer-indexer/util/ctxlog"
	"github.com/monolit-network/xlayer-indexer/util/notifier"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

var skipFillGap = pflag.Bool("skip-fill-gap", false, "skip initial fill-gap from latest_block_in_db to current_block_on_chain")
var excludeProcessor = pflag.StringSlice("exclude-processor", []string{}, "exclude processor. available processors: basic, polymarket")

const retryDelay = 5 * time.Second
const latestBlockPollInterval = 5 * time.Second
const blockHashHistorySize = 200
const deepReorgLogThreshold = 64

type blockHashRef struct {
	Number int64
	Hash   common.Hash
}

type dbBlockHashState struct {
	Number int64
	Hashes []common.Hash
}

type processedBlockHistory struct {
	limit  int
	blocks []blockHashRef
}

func newProcessedBlockHistory(limit int) *processedBlockHistory {
	return &processedBlockHistory{
		limit:  limit,
		blocks: make([]blockHashRef, 0, limit),
	}
}

func (h *processedBlockHistory) Reset(blocks []blockHashRef) {
	if len(blocks) > h.limit {
		blocks = blocks[len(blocks)-h.limit:]
	}

	h.blocks = append(h.blocks[:0], blocks...)
}

func (h *processedBlockHistory) Push(block blockHashRef) {
	if len(h.blocks) > 0 {
		last := h.blocks[len(h.blocks)-1]
		switch {
		case block.Number < last.Number:
			return
		case block.Number == last.Number:
			h.blocks[len(h.blocks)-1] = block
			return
		}
	}

	h.blocks = append(h.blocks, block)
	if len(h.blocks) > h.limit {
		h.blocks = h.blocks[len(h.blocks)-h.limit:]
	}
}

func (h *processedBlockHistory) Latest() (blockHashRef, bool) {
	if len(h.blocks) == 0 {
		return blockHashRef{}, false
	}
	return h.blocks[len(h.blocks)-1], true
}

func (h *processedBlockHistory) Entries() []blockHashRef {
	return append([]blockHashRef(nil), h.blocks...)
}

type observedBlockQueue struct {
	mu        sync.Mutex
	hasValues bool
	minBlock  int64
	maxBlock  int64
	notifyCh  chan struct{}
}

func newObservedBlockQueue() *observedBlockQueue {
	return &observedBlockQueue{
		notifyCh: make(chan struct{}, 1),
	}
}

func (q *observedBlockQueue) Push(blockNumber int64) {
	q.mu.Lock()
	if !q.hasValues {
		q.minBlock = blockNumber
		q.maxBlock = blockNumber
		q.hasValues = true
	} else {
		if blockNumber < q.minBlock {
			q.minBlock = blockNumber
		}
		if blockNumber > q.maxBlock {
			q.maxBlock = blockNumber
		}
	}
	q.mu.Unlock()

	select {
	case q.notifyCh <- struct{}{}:
	default:
	}
}

func (q *observedBlockQueue) PopAll() (minBlock int64, maxBlock int64, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.hasValues {
		return 0, 0, false
	}

	minBlock = q.minBlock
	maxBlock = q.maxBlock
	q.hasValues = false
	return minBlock, maxBlock, true
}

func (q *observedBlockQueue) Wait(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	case <-q.notifyCh:
		return true
	}
}

func initRpcs(chain string) (*ethclient.Client, *ethclient.Client, *rpc.Client, error) {
	rpcURL := os.Getenv(fmt.Sprintf("EVM_RPC_URL_%s", strings.ToUpper(chain)))
	if rpcURL == "" {
		return nil, nil, nil, errors.New("EVM_RPC_URL is not set")
	}
	wsURL := os.Getenv(fmt.Sprintf("EVM_WS_URL_%s", strings.ToUpper(chain)))
	if wsURL == "" {
		return nil, nil, nil, errors.New("EVM_WS_URL is not set")
	}
	t := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}
	httpClient := &http.Client{
		Transport: t,
		Timeout:   300 * time.Second,
	}

	rpcClient, err := rpc.DialOptions(context.Background(), rpcURL, rpc.WithHTTPClient(httpClient))
	if err != nil {
		return nil, nil, nil, err
	}
	ethClient := ethclient.NewClient(rpcClient)

	wsEthClient, err := ethclient.Dial(wsURL)
	if err != nil {
		return nil, nil, nil, err
	}

	return ethClient, wsEthClient, rpcClient, nil
}

func initProcessors(logger *zap.Logger, dbClickClient db.DBClick, dbClient db.Client, evmClient *evmclient.Client, chain string) ([]processors.IndexerProcessor, chan error, *notifier.Notifier) {
	var processors = []processors.IndexerProcessor{}
	errCh := make(chan error)
	var eventNotifier *notifier.Notifier

	notifierTopics := []notifier.Topic{
		notifier.TopicOf[models.SwapEvent](string(sharedmodels.NotifierEventEVMSwap)),
		notifier.TopicOf[models.TransferEvent](string(sharedmodels.NotifierEventEVMTransfer)),
		notifier.TopicOf[models.DefiEvent](string(sharedmodels.NotifierEventEVMDefi)),
		notifier.TopicOf[sharedmodels.PolymarketMarketEvent](string(sharedmodels.NotifierEventCondPrepared)),
		notifier.TopicOf[sharedmodels.PolymarketMarketEvent](string(sharedmodels.NotifierEventCondResolved)),
		notifier.TopicOf[sharedmodels.PolymarketMarketEvent](string(sharedmodels.NotifierEventAdapterUMA)),
		notifier.TopicOf[models.PolymarketOrderEventNew](string(sharedmodels.NotifierEventOrder)),
		notifier.TopicOf[models.ReorgEvent](string(sharedmodels.NotifierEventReorg)),
	}
	eventNotifier, err := notifier.NewNATSFromEnv(logger, notifierTopics...)
	if err != nil {
		logger.Fatal("error initializing notifier", zap.Error(err))
	}

	if !slices.Contains(*excludeProcessor, "basic") {
		basicProc := basicProcessor.NewProcessor(logger, dbClient, dbClickClient, evmClient, eventNotifier, chain)
		processors = append(processors, basicProc)

		go func() {
			for err := range basicProc.Errors() {
				errCh <- fmt.Errorf("basic processor error: %w", err)
			}
		}()
	}

	if strings.EqualFold(chain, string(models.ChainPolygon)) && !slices.Contains(*excludeProcessor, "polymarket") {
		polymarketProc := polymarketProcessor.NewProcessor(evmClient, dbClickClient, eventNotifier, dbClient, logger)
		processors = append(processors, polymarketProc)

		go func() {
			for err := range polymarketProc.Errors() {
				errCh <- fmt.Errorf("polymarket processor error: %w", err)
			}
		}()
	}

	return processors, errCh, eventNotifier
}

func main() {
	pflag.Parse()
	godotenv.Load()
	logger := ctxlog.NewCmdLogger()

	dbClient, err := postgres.NewPostgresClientFromEnv()
	if err != nil {
		logger.Fatal("error creating database client: %v", zap.Error(err))
	}
	defer dbClient.Close()
	logger.Info("database client initialized successfully")

	clickClient, err := clickhouse.NewClickhouseClientFromEnv(logger)
	if err != nil {
		logger.Fatal("error creating clickhouse client", zap.Error(err))
	}
	defer clickClient.Close()

	chain := os.Getenv("EVM_CHAIN")
	if chain == "" {
		logger.Fatal("EVM_CHAIN is not set")
	}

	ethClient, wsEthClient, rpcClient, err := initRpcs(chain)
	if err != nil {
		logger.Fatal("error initializing rpcs", zap.Error(err))
	}
	defer ethClient.Close()
	defer wsEthClient.Close()

	evmClient := evmclient.NewClient(ethClient, wsEthClient, rpcClient, logger, models.Chain(chain), evmclient.EvmClientOptionsFromEnv(models.Chain(chain))...)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	procs, procErrCh, eventNotifier := initProcessors(logger, clickClient, dbClient, evmClient, chain)
	if eventNotifier != nil {
		defer eventNotifier.Close()
	}
	for _, proc := range procs {
		defer proc.Stop()
	}

	aggregator := processors.NewProcessorAggregator(procs, evmClient, logger, processors.WithSyncProcessing(true))
	defer aggregator.Stop()

	sub, err := evmClient.SubscribeHeads(ctx)
	if err != nil {
		logger.Fatal("failed to subscribe to new heads", zap.Error(err))
	}
	logger.Info("subscribed to new heads")

	dbBlockHistory, err := getRecentDBBlockHashes(clickClient, chain, blockHashHistorySize)
	if err != nil {
		logger.Fatal("error getting recent block hashes from db", zap.Error(err))
	}

	latestBlockOnChain, err := ethClient.BlockNumber(ctx)
	if err != nil {
		logger.Fatal("error getting latest block on chain", zap.Error(err))
	}

	blockQueue := newObservedBlockQueue()
	blockHistory := newProcessedBlockHistory(blockHashHistorySize)

	currentLatestBlock := latestDBBlockNumber(dbBlockHistory)
	initialTargetBlock := currentLatestBlock
	startupReorgHandled := false

	if mismatchBlock, foundMismatch, err := findFirstDBHashMismatchWithRetry(ctx, logger, ethClient, dbBlockHistory); err != nil {
		logger.Fatal("failed to verify db block hashes against node", zap.Error(err))
	} else if foundMismatch {
		startupReorgHandled = true
		var rebuiltHistory *processedBlockHistory
		currentLatestBlock, initialTargetBlock, rebuiltHistory, err = handleDetectedReorg(
			ctx,
			logger,
			ethClient,
			evmClient,
			clickClient,
			dbClient,
			eventNotifier,
			chain,
			currentLatestBlock,
			mismatchBlock,
			"startup db hash mismatch",
		)
		if err != nil {
			logger.Fatal("failed to handle startup reorg", zap.Error(err))
		}
		blockHistory = rebuiltHistory
	}

	if *skipFillGap && startupReorgHandled {
		logger.Warn("skip fill gap ignored after startup reorg",
			zap.Int64("current_latest_block", currentLatestBlock),
			zap.Int64("target_block", initialTargetBlock),
		)
	} else if *skipFillGap {
		currentLatestBlock = int64(latestBlockOnChain)
		initialTargetBlock = currentLatestBlock
		logger.Warn("skip fill gap enabled",
			zap.Int64("latest_block_in_db", latestDBBlockNumber(dbBlockHistory)),
			zap.Uint64("latest_block_on_chain", latestBlockOnChain),
			zap.Int64("starting_from_block", currentLatestBlock+1),
		)
	} else if currentLatestBlock < int64(latestBlockOnChain) {
		initialTargetBlock = int64(latestBlockOnChain)
		logger.Info("queueing initial target",
			zap.Int64("start_block", currentLatestBlock+1),
			zap.Uint64("end_block", latestBlockOnChain),
		)
	}

	if !startupReorgHandled || *skipFillGap {
		rebuiltHistory, err := buildCanonicalBlockHashHistoryUntilSuccess(ctx, logger, evmClient, currentLatestBlock, blockHashHistorySize)
		if err != nil {
			logger.Fatal("failed to build startup block hash history", zap.Error(err))
		}
		blockHistory = rebuiltHistory
	}

	if initialTargetBlock > currentLatestBlock {
		blockQueue.Push(initialTargetBlock)
	}

	logger.Info("current latest block", zap.Int64("block_number", currentLatestBlock))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	var workersWg sync.WaitGroup

	workersWg.Add(1)
	go func() {
		defer workersWg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-procErrCh:
				logger.Error("error in processor", zap.Error(err))
			case err := <-aggregator.Errors():
				logger.Error("error in aggregator", zap.Error(err))
			}
		}
	}()

	workersWg.Add(1)
	go func() {
		defer workersWg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case head, ok := <-sub:
				if !ok {
					return
				}
				if head == nil {
					logger.Error("head is nil")
					continue
				}
				blockQueue.Push(head.Number.Int64())
			}
		}
	}()

	workersWg.Add(1)
	go func() {
		defer workersWg.Done()
		pollLatestBlockTargets(ctx, logger, ethClient, initialTargetBlock, blockQueue)
	}()

	workersWg.Add(1)
	go func() {
		defer workersWg.Done()
		runSequentialBlockWorker(
			ctx,
			logger,
			aggregator,
			blockQueue,
			ethClient,
			evmClient,
			clickClient,
			dbClient,
			eventNotifier,
			chain,
			currentLatestBlock,
			blockHistory,
		)
	}()

	<-sigCh
	logger.Info("received signal, shutting down")
	cancel()
	workersWg.Wait()
	aggregator.FlushSync()
	logger.Info("flush sync completed")
}

func runSequentialBlockWorker(
	ctx context.Context,
	logger *zap.Logger,
	aggregator *processors.ProcessorAggregator,
	blockQueue *observedBlockQueue,
	ethClient *ethclient.Client,
	evmClient *evmclient.Client,
	clickClient db.DBClick,
	dbClient db.Client,
	eventNotifier *notifier.Notifier,
	chain string,
	currentLatestBlock int64,
	blockHistory *processedBlockHistory,
) {
	targetBlock := currentLatestBlock
	reorgCheckRequested := false

	for {
		if minObservedBlock, maxObservedBlock, ok := blockQueue.PopAll(); ok {
			if minObservedBlock <= currentLatestBlock {
				reorgCheckRequested = true
			}

			if maxObservedBlock > targetBlock {
				targetBlock = maxObservedBlock
			}
		}

		if reorgCheckRequested {
			reorgStartBlock, foundMismatch, err := findFirstHashMismatchWithRetry(ctx, logger, ethClient, blockHistory.Entries())
			if err != nil {
				return
			}
			if foundMismatch {
				var rebuiltHistory *processedBlockHistory
				currentLatestBlock, targetBlock, rebuiltHistory, err = handleDetectedReorg(
					ctx,
					logger,
					ethClient,
					evmClient,
					clickClient,
					dbClient,
					eventNotifier,
					chain,
					currentLatestBlock,
					reorgStartBlock,
					"observed head at or below current height",
				)
				if err != nil {
					return
				}
				blockHistory = rebuiltHistory
			}
			reorgCheckRequested = false
			continue
		}

		if currentLatestBlock < targetBlock {
			reorgStartBlock, foundMismatch, err := detectReorgBeforeNextBlock(ctx, logger, ethClient, currentLatestBlock, blockHistory)
			if err != nil {
				return
			}
			if foundMismatch {
				var rebuiltHistory *processedBlockHistory
				currentLatestBlock, targetBlock, rebuiltHistory, err = handleDetectedReorg(
					ctx,
					logger,
					ethClient,
					evmClient,
					clickClient,
					dbClient,
					eventNotifier,
					chain,
					currentLatestBlock,
					reorgStartBlock,
					"next block parent hash mismatch",
				)
				if err != nil {
					return
				}
				blockHistory = rebuiltHistory
				continue
			}

			nextBlock := currentLatestBlock + 1
			blockHash, err := processBlockUntilSuccess(ctx, logger, aggregator, nextBlock)
			if err != nil {
				return
			}
			currentLatestBlock = nextBlock
			blockHistory.Push(blockHashRef{
				Number: nextBlock,
				Hash:   blockHash,
			})
			continue
		}

		if !blockQueue.Wait(ctx) {
			return
		}
	}
}

func getRecentDBBlockHashes(clickClient db.DBClick, chain string, limit int) ([]dbBlockHashState, error) {
	sources := []string{
		"SELECT DISTINCT block_number, lower(block_hash) AS block_hash FROM evm.swap_events WHERE chain = ?",
		"SELECT DISTINCT block_number, lower(block_hash) AS block_hash FROM evm.transfer_events WHERE chain = ?",
		"SELECT DISTINCT block_number, lower(block_hash) AS block_hash FROM evm.defi_events WHERE chain = ?",
	}
	args := []any{chain, chain, chain}

	if strings.EqualFold(chain, string(models.ChainPolygon)) {
		sources = append(sources, "SELECT DISTINCT block_number, lower(block_hash) AS block_hash FROM evm.polymarket_order_events")
	}

	query := fmt.Sprintf(`
		SELECT block_number, groupUniqArray(block_hash) AS block_hashes
		FROM (%s)
		GROUP BY block_number
		ORDER BY block_number DESC
		LIMIT ?
	`, strings.Join(sources, " UNION ALL "))
	args = append(args, limit)

	rows, err := clickClient.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blocks := make([]dbBlockHashState, 0, limit)
	for rows.Next() {
		var blockNumber big.Int
		var blockHashes []string
		if err := rows.Scan(&blockNumber, &blockHashes); err != nil {
			return nil, err
		}

		hashes := make([]common.Hash, 0, len(blockHashes))
		for _, blockHash := range blockHashes {
			normalizedHash := strings.TrimRight(strings.TrimSpace(blockHash), "\x00")
			if normalizedHash == "" {
				continue
			}
			hashes = append(hashes, common.HexToHash(normalizedHash))
		}

		blocks = append(blocks, dbBlockHashState{
			Number: blockNumber.Int64(),
			Hashes: hashes,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	slices.Reverse(blocks)
	return blocks, nil
}

func latestDBBlockNumber(blocks []dbBlockHashState) int64 {
	if len(blocks) == 0 {
		return 0
	}
	return blocks[len(blocks)-1].Number
}

func findFirstDBHashMismatchWithRetry(ctx context.Context, logger *zap.Logger, ethClient *ethclient.Client, dbBlocks []dbBlockHashState) (int64, bool, error) {
	for _, block := range dbBlocks {
		if len(block.Hashes) != 1 {
			return block.Number, true, nil
		}

		header, err := getBlockHeaderByNumberWithRetry(ctx, logger, ethClient, block.Number)
		if err != nil {
			return 0, false, err
		}
		if header.Hash() != block.Hashes[0] {
			return block.Number, true, nil
		}
	}

	return 0, false, nil
}

func findFirstHashMismatchWithRetry(ctx context.Context, logger *zap.Logger, ethClient *ethclient.Client, blocks []blockHashRef) (int64, bool, error) {
	for _, block := range blocks {
		header, err := getBlockHeaderByNumberWithRetry(ctx, logger, ethClient, block.Number)
		if err != nil {
			return 0, false, err
		}
		if header.Hash() != block.Hash {
			return block.Number, true, nil
		}
	}

	return 0, false, nil
}

func getBlockHeaderByNumberWithRetry(ctx context.Context, logger *zap.Logger, ethClient *ethclient.Client, blockNumber int64) (*types.Header, error) {
	block := big.NewInt(blockNumber)
	for {
		header, err := ethClient.HeaderByNumber(ctx, block)
		if err == nil {
			return header, nil
		}

		if errors.Is(err, context.Canceled) {
			return nil, err
		}

		logger.Error("failed to get block header by number, retrying",
			zap.Int64("block_number", blockNumber),
			zap.Error(err),
		)
		if err := sleepWithContext(ctx, retryDelay); err != nil {
			return nil, err
		}
	}
}

func buildCanonicalBlockHashHistory(ctx context.Context, evmClient *evmclient.Client, currentLatestBlock int64, limit int) (*processedBlockHistory, error) {
	history := newProcessedBlockHistory(limit)
	if currentLatestBlock <= 0 {
		return history, nil
	}

	startBlock := currentLatestBlock - int64(limit) + 1
	if startBlock < 1 {
		startBlock = 1
	}

	headers, err := evmClient.BatchGetBlocksHeaders(ctx, big.NewInt(startBlock), big.NewInt(currentLatestBlock))
	if err != nil {
		return nil, err
	}

	blocks := make([]blockHashRef, 0, len(headers))
	for _, header := range headers {
		if header == nil || header.Number == nil {
			return nil, fmt.Errorf("got nil header while building block hash history")
		}
		blocks = append(blocks, blockHashRef{
			Number: header.Number.Int64(),
			Hash:   header.Hash(),
		})
	}

	history.Reset(blocks)
	return history, nil
}

func buildCanonicalBlockHashHistoryUntilSuccess(ctx context.Context, logger *zap.Logger, evmClient *evmclient.Client, currentLatestBlock int64, limit int) (*processedBlockHistory, error) {
	for {
		history, err := buildCanonicalBlockHashHistory(ctx, evmClient, currentLatestBlock, limit)
		if err == nil {
			return history, nil
		}

		if errors.Is(err, context.Canceled) {
			return nil, err
		}

		logger.Error("failed to build canonical block hash history, retrying",
			zap.Int64("current_latest_block", currentLatestBlock),
			zap.Error(err),
		)
		if err := sleepWithContext(ctx, retryDelay); err != nil {
			return nil, err
		}
	}
}

func detectReorgBeforeNextBlock(ctx context.Context, logger *zap.Logger, ethClient *ethclient.Client, currentLatestBlock int64, blockHistory *processedBlockHistory) (int64, bool, error) {
	if currentLatestBlock <= 0 {
		return 0, false, nil
	}

	latestBlock, ok := blockHistory.Latest()
	if !ok {
		return 0, false, nil
	}
	if latestBlock.Number != currentLatestBlock {
		logger.Warn("block hash history is not aligned with current latest block",
			zap.Int64("history_latest_block", latestBlock.Number),
			zap.Int64("current_latest_block", currentLatestBlock),
		)
		return findFirstHashMismatchWithRetry(ctx, logger, ethClient, blockHistory.Entries())
	}

	nextHeader, err := getBlockHeaderByNumberWithRetry(ctx, logger, ethClient, currentLatestBlock+1)
	if err != nil {
		return 0, false, err
	}
	if nextHeader.ParentHash == latestBlock.Hash {
		return 0, false, nil
	}

	return findFirstHashMismatchWithRetry(ctx, logger, ethClient, blockHistory.Entries())
}

func handleDetectedReorg(
	ctx context.Context,
	logger *zap.Logger,
	ethClient *ethclient.Client,
	evmClient *evmclient.Client,
	clickClient db.DBClick,
	dbClient db.Client,
	eventNotifier *notifier.Notifier,
	chain string,
	currentLatestBlock int64,
	reorgStartBlock int64,
	reason string,
) (int64, int64, *processedBlockHistory, error) {
	rollbackDepth := currentLatestBlock - reorgStartBlock + 1
	if rollbackDepth < 1 {
		rollbackDepth = 1
	}

	if rollbackDepth > deepReorgLogThreshold {
		logger.Error("deep reorg detected",
			zap.String("reason", reason),
			zap.Int64("reorg_start_block", reorgStartBlock),
			zap.Int64("current_latest_block", currentLatestBlock),
			zap.Int64("rollback_depth", rollbackDepth),
		)
	}

	logger.Warn("reorg detected",
		zap.String("reason", reason),
		zap.Int64("reorg_start_block", reorgStartBlock),
		zap.Int64("current_latest_block", currentLatestBlock),
		zap.Int64("rollback_depth", rollbackDepth),
	)

	if err := handleReorgUntilSuccess(ctx, logger, clickClient, dbClient, eventNotifier, chain, reorgStartBlock); err != nil {
		return 0, 0, nil, err
	}

	currentLatestBlock = reorgStartBlock - 1
	blockHistory, err := buildCanonicalBlockHashHistoryUntilSuccess(ctx, logger, evmClient, currentLatestBlock, blockHashHistorySize)
	if err != nil {
		return 0, 0, nil, err
	}

	latestBlockOnChain, err := latestBlockOnChainWithRetry(ctx, logger, ethClient)
	if err != nil {
		return 0, 0, nil, err
	}

	targetBlock := int64(latestBlockOnChain)
	logger.Info("reorg handled",
		zap.Int64("restart_block", currentLatestBlock+1),
		zap.Int64("target_block", targetBlock),
	)

	return currentLatestBlock, targetBlock, blockHistory, nil
}

func pollLatestBlockTargets(ctx context.Context, logger *zap.Logger, ethClient *ethclient.Client, initialLatestBlock int64, blockQueue *observedBlockQueue) {
	ticker := time.NewTicker(latestBlockPollInterval)
	defer ticker.Stop()

	lastObservedBlock := initialLatestBlock
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			latestBlockOnChain, err := ethClient.BlockNumber(ctx)
			if err != nil {
				logger.Warn("failed to poll latest block on chain", zap.Error(err))
				continue
			}

			latestBlock := int64(latestBlockOnChain)
			if latestBlock > lastObservedBlock {
				blockQueue.Push(latestBlock)
			}
			lastObservedBlock = latestBlock
		}
	}
}

func latestBlockOnChainWithRetry(ctx context.Context, logger *zap.Logger, ethClient *ethclient.Client) (uint64, error) {
	for {
		latestBlockOnChain, err := ethClient.BlockNumber(ctx)
		if err == nil {
			return latestBlockOnChain, nil
		}

		logger.Error("failed to get latest block on chain, retrying", zap.Error(err))
		if err := sleepWithContext(ctx, retryDelay); err != nil {
			return 0, err
		}
	}
}

func processBlockUntilSuccess(ctx context.Context, logger *zap.Logger, aggregator *processors.ProcessorAggregator, blockNumber int64) (common.Hash, error) {
	block := big.NewInt(blockNumber)
	for {
		start := time.Now()
		blockHash, err := aggregator.ProcessBlockNumberSyncWithHash(ctx, block)
		if err == nil {
			logger.Info("block processed",
				zap.Int64("block_number", blockNumber),
				zap.Duration("duration", time.Since(start)),
			)
			return blockHash, nil
		}

		if errors.Is(err, context.Canceled) {
			return common.Hash{}, err
		}

		logger.Error("failed to process block, retrying",
			zap.Int64("block_number", blockNumber),
			zap.Error(err),
		)
		if err := sleepWithContext(ctx, retryDelay); err != nil {
			return common.Hash{}, err
		}
	}
}

func handleReorgUntilSuccess(ctx context.Context, logger *zap.Logger, dbClick db.DBClick, dbPsql db.Client, eventNotifier *notifier.Notifier, chain string, blockNumber int64) error {
	for {
		err := handleReorg(ctx, logger, dbClick, dbPsql, eventNotifier, chain, blockNumber)
		if err == nil {
			return nil
		}

		if errors.Is(err, context.Canceled) {
			return err
		}

		logger.Error("failed to handle reorg, retrying",
			zap.Int64("block_number", blockNumber),
			zap.Error(err),
		)
		if err := sleepWithContext(ctx, retryDelay); err != nil {
			return err
		}
	}
}

func sleepWithContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func handleReorg(ctx context.Context, logger *zap.Logger, dbClick db.DBClick, dbPsql db.Client, eventNotifier *notifier.Notifier, chain string, blockNumber int64) error {
	if eventNotifier != nil {
		if err := notifier.Publish(eventNotifier, ctx, string(sharedmodels.NotifierEventReorg), models.ReorgEvent{
			BlockNumber: big.NewInt(blockNumber),
			Chain:       chain,
		}); err != nil {
			eventNotifier.Logger.Error("failed to publish reorg event", zap.Error(err))
		}
	}
	if chain == string(models.ChainPolygon) {
		err := dbPsql.Transaction(ctx, func(ctx context.Context, tx db.DB) error {
			if _, err := tx.DeletePolymarketMarketsFromBlockNumber(ctx, blockNumber); err != nil {
				return fmt.Errorf("failed to delete polymarket markets from block number: %w", err)
			}
			if _, err := tx.DeletePolymarketMarketEventsFromBlockNumber(ctx, blockNumber); err != nil {
				return fmt.Errorf("failed to delete polymarket market events from block number: %w", err)
			}
			if _, err := tx.ResetPolymarketTokensResolutionFromBlockNumber(ctx, blockNumber); err != nil {
				return fmt.Errorf("failed to reset polymarket tokens resolution from block number: %w", err)
			}
			if _, err := tx.ResetPolymarketMarketsResolutionFromBlockNumber(ctx, blockNumber); err != nil {
				return fmt.Errorf("failed to reset polymarket markets resolution from block number: %w", err)
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	if err := dbClick.DeleteEvmEventsFromBlockNumber(ctx, models.Chain(chain), uint64(blockNumber)); err != nil {
		return fmt.Errorf("failed to delete evm events by block number: %w", err)
	}

	return nil
}
