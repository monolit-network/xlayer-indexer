package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/shared/db/postgres"
	"github.com/monolit-network/xlayer-indexer/util/ctxlog"
	"github.com/monolit-network/xlayer-indexer/util/iterators"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

var (
	startBlock = 4000000
	endBlock   = 85_000_000
	batchSize  = 512
	numWorkers = 30
)

func main() {
	godotenv.Load()
	logger := ctxlog.NewCmdLogger()

	dbClient, err := postgres.NewPostgresClientFromEnv()
	if err != nil {
		logger.Fatal("error creating database client", zap.Error(err))
	}
	defer dbClient.Close()
	logger.Info("database client initialized")

	rpcURL := os.Getenv("EVM_RPC_URL_POLYGON")
	if rpcURL == "" {
		logger.Fatal("EVM_RPC_URL_POLYGON is not set")
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
		logger.Fatal("error connecting to rpc", zap.Error(err))
	}
	ethClient := ethclient.NewClient(rpcClient)
	defer ethClient.Close()

	evmClient := evmclient.NewClient(ethClient, nil, rpcClient, logger, models.ChainPolygon, evmclient.EvmClientOptionsFromEnv(models.ChainPolygon)...)
	logger.Info("evm client initialized")

	latestBlockNumber, err := evmClient.GetLatestBlockNumber()
	if err != nil {
		logger.Fatal("error getting latest block number", zap.Error(err))
	}
	endBlock = int(latestBlockNumber.Int64())

	ctx := context.Background()

	if err := fillBlockHashes(ctx, logger, dbClient, evmClient); err != nil {
		logger.Fatal("error filling block hashes", zap.Error(err))
	}

	logger.Info("migration completed successfully")
}

func fillBlockHashes(ctx context.Context, logger *zap.Logger, dbClient *postgres.PostgresClient, evmClient *evmclient.Client) error {
	logger.Info("starting block hash migration",
		zap.Int("start_block", startBlock),
		zap.Int("end_block", endBlock),
		zap.Int("batch_size", batchSize),
		zap.Int("workers", numWorkers),
	)

	totalBatches := (endBlock - startBlock + batchSize - 1) / batchSize
	var processedBatches atomic.Int64
	var totalMarketsUpdated atomic.Int64
	var totalEventsUpdated atomic.Int64

	sem := semaphore.NewWeighted(int64(numWorkers))
	var wg sync.WaitGroup
	var firstErr error
	var errMu sync.Mutex

	startTime := time.Now()

	for _, batch := range iterators.BigRangeSubranges(big.NewInt(int64(startBlock)), big.NewInt(int64(endBlock)), batchSize) {
		if err := sem.Acquire(ctx, 1); err != nil {
			return fmt.Errorf("failed to acquire semaphore: %w", err)
		}

		wg.Add(1)
		go func(batchStart, batchEnd *big.Int) {
			defer wg.Done()
			defer sem.Release(1)

			batchStartTime := time.Now()

			headers, err := evmClient.BatchGetBlocksHeaders(ctx, batchStart, batchEnd)
			if err != nil {
				errMu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("failed to fetch headers for blocks %s-%s: %w", batchStart, batchEnd, err)
				}
				errMu.Unlock()
				logger.Error("failed to fetch headers", zap.String("start", batchStart.String()), zap.String("end", batchEnd.String()), zap.Error(err))
				return
			}

			if len(headers) == 0 {
				processedBatches.Add(1)
				logger.Info("no headers fetched", zap.String("start", batchStart.String()), zap.String("end", batchEnd.String()))
				return
			}

			// Build block number to hash map
			blockHashes := make(map[int64]string, len(headers))
			for _, header := range headers {
				if header != nil {
					blockHashes[header.Number.Int64()] = strings.ToLower(header.Hash().Hex())
				}
			}

			// Update both tables
			marketsUpdated, eventsUpdated, err := updateTables(dbClient, blockHashes, batchStart.Int64(), batchEnd.Int64())
			if err != nil {
				errMu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("failed to update tables for blocks %s-%s: %w", batchStart, batchEnd, err)
				}
				errMu.Unlock()
				logger.Error("failed to update tables", zap.String("start", batchStart.String()), zap.String("end", batchEnd.String()), zap.Error(err))
				return
			}

			totalMarketsUpdated.Add(marketsUpdated)
			totalEventsUpdated.Add(eventsUpdated)
			processed := processedBatches.Add(1)

			if processed%100 == 0 {
				logger.Info("progress",
					zap.Int64("batches_processed", processed),
					zap.Int("total_batches", totalBatches),
					zap.Int64("markets_updated", totalMarketsUpdated.Load()),
					zap.Int64("events_updated", totalEventsUpdated.Load()),
					zap.Duration("elapsed", time.Since(startTime)),
				)
			}

			logger.Debug("batch processed",
				zap.String("start", batchStart.String()),
				zap.String("end", batchEnd.String()),
				zap.Int("headers_fetched", len(headers)),
				zap.Int64("markets_updated", marketsUpdated),
				zap.Int64("events_updated", eventsUpdated),
				zap.Duration("duration", time.Since(batchStartTime)),
			)
		}(batch.Start, batch.End)
	}

	wg.Wait()

	if firstErr != nil {
		return firstErr
	}

	logger.Info("migration finished",
		zap.Int64("total_batches_processed", processedBatches.Load()),
		zap.Int64("total_markets_updated", totalMarketsUpdated.Load()),
		zap.Int64("total_events_updated", totalEventsUpdated.Load()),
		zap.Duration("total_duration", time.Since(startTime)),
	)
	return nil
}

func updateTables(dbClient *postgres.PostgresClient, blockHashes map[int64]string, batchStart, batchEnd int64) (int64, int64, error) {
	var marketsUpdated, eventsUpdated int64

	mapJson, err := json.Marshal(blockHashes)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to marshal block hashes: %w", err)
	}

	ctx := context.Background()

	// Update markets table
	result, err := dbClient.Querier().Exec(ctx, `
		UPDATE polymarket.markets m
		SET prepared_in_block_hash = v.block_hash
		FROM (
			SELECT key::bigint as block_num, value as block_hash
			FROM jsonb_each_text($1::jsonb)
		) v
		WHERE m.prepared_in_block = v.block_num
		AND m.prepared_in_block >= $2 AND m.prepared_in_block <= $3
		AND (m.prepared_in_block_hash IS NULL OR m.prepared_in_block_hash = '')
	`, string(mapJson), batchStart, batchEnd)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to update markets: %w", err)
	}
	marketsUpdated = result.RowsAffected()

	// Update market_events table
	result, err = dbClient.Querier().Exec(ctx, `
		UPDATE polymarket.market_events e
		SET block_hash = v.block_hash
		FROM (
			SELECT key::bigint as block_num, value as block_hash
			FROM jsonb_each_text($1::jsonb)
		) v
		WHERE e.block_number = v.block_num
		AND e.block_number >= $2 AND e.block_number <= $3
		AND (e.block_hash IS NULL OR e.block_hash = '')
	`, string(mapJson), batchStart, batchEnd)
	if err != nil {
		return marketsUpdated, 0, fmt.Errorf("failed to update market_events: %w", err)
	}
	eventsUpdated = result.RowsAffected()

	return marketsUpdated, eventsUpdated, nil
}
