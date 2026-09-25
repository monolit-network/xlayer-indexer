package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/processors/polymarketProcessor"
	"github.com/monolit-network/xlayer-indexer/shared/db/clickhouse"
	"github.com/monolit-network/xlayer-indexer/shared/db/postgres"
	"github.com/monolit-network/xlayer-indexer/util/ctxlog"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Special binary to backfill Polymarket redemption events.

var (
	startBlock              = flag.String("start-block", "", "start block")
	endBlock                = flag.String("end-block", "", "end block")
	disableClickhouseWriter = flag.Bool("disable-click-writer", false, "disable clickhouse writer")
)

const receiptWorkers = 1000

func main() {
	godotenv.Load()
	flag.Parse()
	logger := ctxlog.NewCmdLogger()
	start = time.Now()

	dbClient, err := postgres.NewPostgresClientFromEnv()
	if err != nil {
		logger.Fatal("error creating database client: %v", zap.Error(err))
	}
	if err := dbClient.Querier().CreatePolymarketTokensTable(context.Background()); err != nil {
		logger.Fatal("error creating polymarket tokens table: %v", zap.Error(err))
	}
	if err := dbClient.Querier().CreatePolymarketMarketsNewTable(context.Background()); err != nil {
		logger.Fatal("error creating polymarket markets new table: %v", zap.Error(err))
	}
	if err := dbClient.Querier().CreatePolymarketMarketEventsTable(context.Background()); err != nil {
		logger.Fatal("error creating polymarket market events table: %v", zap.Error(err))
	}
	logger.Info("polymarket tables created successfully")
	defer dbClient.Close()
	logger.Info("database client initialized successfully")

	clickClient, err := clickhouse.NewClickhouseClientFromEnv(logger)
	if err != nil {
		logger.Fatal("error creating clickhouse client", zap.Error(err))
	}
	defer clickClient.Close()
	logger.Info("clickhouse client initialized successfully")

	rpcURL := os.Getenv("EVM_RPC_URL")
	if rpcURL == "" {
		logger.Fatal("EVM_RPC_URL is not set")
	}

	t := &http.Transport{
		MaxIdleConns:        2000,
		MaxIdleConnsPerHost: 2000,
		IdleConnTimeout:     90 * time.Second,
	}

	httpClient := &http.Client{
		Transport: t,
		Timeout:   300 * time.Second,
	}

	rpcClient, err := rpc.DialOptions(context.Background(), rpcURL, rpc.WithHTTPClient(httpClient))
	if err != nil {
		logger.Fatal("error dialing rpc", zap.Error(err))
	}
	ethClient := ethclient.NewClient(rpcClient)
	defer ethClient.Close()
	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		logger.Fatal("error getting chain id", zap.Error(err))
	}
	logger.Info("rpc client initialized successfully", zap.Int64("chain_id",
		chainID.Int64()),
		zap.String("chain", string(models.IdToChain[models.ChainId(chainID.Int64())])))

	chain := os.Getenv("EVM_CHAIN")
	if chain == "" {
		logger.Fatal("EVM_CHAIN is not set")
	}

	evmClientOptions := evmclient.EvmClientOptionsFromEnv(models.Chain(chain))
	evmClient := evmclient.NewClient(ethClient, ethClient, rpcClient, logger, models.Chain(chain), evmClientOptions...)

	startBlock, err := parseBlockNumber(*startBlock)
	if err != nil {
		logger.Fatal("invalid start-block", zap.Error(err))
	}
	startBlockNumber = new(big.Int).Set(startBlock)
	endBlock, err := parseBlockNumber(*endBlock)
	if err != nil {
		logger.Fatal("invalid end-block", zap.Error(err))
	}
	if endBlock.Cmp(startBlock) < 0 {
		logger.Fatal("end-block must be >= start-block")
	}

	latestBlockNumber, err := evmClient.GetLatestBlockNumber()
	if err != nil {
		logger.Fatal("error getting latest block number", zap.Error(err))
	}
	if endBlock.Cmp(latestBlockNumber) > 0 {
		endBlock = latestBlockNumber
	}

	logger.Info("parsing range", zap.Int64("start_block", startBlock.Int64()), zap.Int64("end_block", endBlock.Int64()), zap.Int64("total_blocks", big.NewInt(0).Sub(endBlock, startBlock).Int64()+1))

	processor := polymarketProcessor.NewProcessor(evmClient, clickClient, nil, dbClient, logger, polymarketProcessor.WithDisableClickhouseWriter(*disableClickhouseWriter))
	defer processor.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signalCh)

	doneCh := make(chan struct{})
	errCh := make(chan error, 1)

	go func() {
		errCh <- streamLogs(ctx, evmClient, logger, startBlock, endBlock, processor)
		close(doneCh)
	}()

	startTime := time.Now()
	defer func() {
		processor.LogTimings()
		logger.Info("time taken", zap.Duration("duration", time.Since(startTime)))
	}()

	select {
	case err := <-errCh:
		if err != nil {
			logger.Error("streaming logs stopped with error", zap.Error(err))
			return
		}
		logger.Info("completed parsing range")
	case sig := <-signalCh:
		logger.Info("received shutdown signal", zap.String("signal", sig.String()))
		cancel()
		logger.Info("waiting for streaming logs to finish")
		<-doneCh
		if err := <-errCh; err != nil {
			logger.Error("streaming logs stopped with error", zap.Error(err))
		}
	case <-ctx.Done():
		logger.Info("context canceled", zap.Error(ctx.Err()))
		<-doneCh
		if err := <-errCh; err != nil {
			logger.Error("streaming logs stopped with error", zap.Error(err))
		}
	}
}

func parseBlockNumber(s string) (*big.Int, error) {
	if len(s) > 2 && (s[0:2] == "0x" || s[0:2] == "0X") {
		v, ok := new(big.Int).SetString(s[2:], 16)
		if !ok {
			return nil, fmt.Errorf("invalid hex block number: %s", s)
		}
		return v, nil
	}
	iv, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return big.NewInt(iv), nil
}

func streamLogs(
	ctx context.Context,
	evmClient *evmclient.Client,
	logger *zap.Logger,
	startBlock *big.Int,
	endBlock *big.Int,
	processor *polymarketProcessor.Processor,
) error {
	redeemAddresses := []common.Address{
		models.PolymarketConditionalTokensAddress,
	}
	redeemTopics := [][]common.Hash{
		{
			models.PolymarketPayoutRedemptionEventSelectorHash,
		},
	}

	streamCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	logsCh := evmClient.StreamGetLogs(
		streamCtx,
		[][]common.Address{redeemAddresses},
		[][][]common.Hash{redeemTopics},
		startBlock,
		endBlock,
	)

	jobs := make(chan groupingResult, receiptWorkers*4)
	receiptSem := make(chan struct{}, receiptWorkers)
	var workersWg sync.WaitGroup
	var errMu sync.Mutex
	var firstErr error

	setErr := func(err error) {
		if err == nil {
			return
		}
		errMu.Lock()
		if firstErr == nil {
			firstErr = err
			cancel()
		}
		errMu.Unlock()
	}

	for i := 0; i < receiptWorkers; i++ {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()
			for group := range jobs {
				if err := processRedeemTx(evmClient, processor, logger, receiptSem, group); err != nil {
					setErr(err)
				}
			}
		}()
	}

	latestBlockMark := uint64(0)
	prevTime := time.Now()
	blockMarkSize := uint64(10000)

	groupingCh := groupLogsIntoTxs(logsCh)

	first := true
discover:
	for group := range groupingCh {
		if group.Err != nil {
			if errors.Is(group.Err, context.Canceled) && streamCtx.Err() != nil {
				break
			}
			logger.Error("error streaming logs", zap.Error(group.Err))
			setErr(group.Err)
			break
		}
		if first {
			first = false
			latestBlockMark = group.BlockNumber
			prevTime = time.Now()
			logger.Info("first log", zap.String("hash", group.TxHash.Hex()), zap.Uint64("block_number", group.BlockNumber))
		}

		if group.BlockNumber >= latestBlockMark+blockMarkSize {
			timeTaken := time.Since(prevTime)
			blocksProcessed := group.BlockNumber - latestBlockMark

			logger.Info("progress",
				zap.Uint64("block_number", group.BlockNumber),
				zap.Duration("time_taken", timeTaken),
				zap.Float64("blocks_per_second", float64(blocksProcessed)/timeTaken.Seconds()),
				zap.Float64("hours_for_1mil_blocks", 1000000.0/float64(blocksProcessed)*float64(timeTaken.Seconds())/3600.0),
			)

			processor.LogTimingsCurrent()
			logger.Info("time taken iter", zap.Duration("duration", time.Since(prevTime)))

			prevTime = time.Now()
			latestBlockMark = group.BlockNumber
		}

		select {
		case jobs <- group:
		case <-streamCtx.Done():
			break discover
		}
	}

	close(jobs)
	workersWg.Wait()

	errMu.Lock()
	err := firstErr
	errMu.Unlock()
	if err != nil {
		return err
	}

	logger.Info("stopped streaming logs")
	return nil
}

var totalTxs atomic.Int64
var start time.Time
var startBlockNumber *big.Int

func processRedeemTx(
	evmClient *evmclient.Client,
	processor *polymarketProcessor.Processor,
	logger *zap.Logger,
	receiptSem chan struct{},
	group groupingResult,
) error {
	receiptSem <- struct{}{}
	defer func() {
		<-receiptSem
	}()
	receipt, err := evmClient.GetReceipt(group.TxHash)
	if err != nil {
		logger.Error("error getting receipt", zap.String("hash", group.TxHash.Hex()), zap.Uint64("block_number", group.BlockNumber), zap.Error(err))
		return fmt.Errorf("error getting receipt: %w", err)
	}
	if receipt == nil {
		logger.Warn("receipt is nil", zap.String("hash", group.TxHash.Hex()), zap.Uint64("block_number", group.BlockNumber))
		return nil
	}

	logs := make([]*types.Log, 0, len(receipt.Logs))
	for _, log := range receipt.Logs {
		if len(log.Topics) == 0 {
			continue
		}
		if log.Topics[0] == models.PolymarketTransferBatchEventSelectorHash ||
			log.Topics[0] == models.PolymarketTransferSingleEventSelectorHash ||
			log.Topics[0] == models.PolymarketPayoutRedemptionEventSelectorHash ||
			log.Topics[0] == models.PolymarketNegRiskPayoutRedemptionEventSelectorHash {
			logs = append(logs, log)
		}
	}

	if err := processor.ProcessLogsFromTx(receipt.Logs); err != nil {
		logger.Error("error processing logs from tx", zap.String("hash", group.TxHash.Hex()), zap.Uint64("block_number", group.BlockNumber), zap.Error(err))
		return fmt.Errorf("error processing logs from tx: %w", err)
	}

	currentTotal := totalTxs.Add(1)
	if currentTotal%1000 == 0 {
		speed := float64(currentTotal) / time.Since(start).Seconds()
		bps := new(big.Int).Sub(big.NewInt(int64(group.BlockNumber)), startBlockNumber)
		bpsFloat := float64(bps.Int64()) / float64(time.Since(start).Seconds())
		logger.Info("processed txs", zap.Int64("total_txs", currentTotal), zap.Duration("elapsed", time.Since(start)), zap.Float64("txps", speed), zap.Float64("bps", bpsFloat))
	}
	return nil
}

type groupingResult struct {
	Logs        []*types.Log
	TxHash      common.Hash
	BlockNumber uint64
	Err         error
}

func groupLogsIntoTxs(logsCh <-chan evmclient.LogResult) chan groupingResult {
	groupingCh := make(chan groupingResult)
	go func() {
		defer close(groupingCh)
		logs := make([]*types.Log, 0, 10)
		prevTxHash := common.Hash{}
		for log := range logsCh {
			if log.Err != nil {
				groupingCh <- groupingResult{Err: log.Err}
				return
			}

			if log.Log.TxHash != prevTxHash {
				if len(logs) > 0 {
					groupingCh <- groupingResult{TxHash: prevTxHash, Logs: logs, BlockNumber: logs[0].BlockNumber}
					logs = make([]*types.Log, 0, 10)
				}
				prevTxHash = log.Log.TxHash
			}

			logs = append(logs, &log.Log)
		}
		if len(logs) > 0 {
			groupingCh <- groupingResult{TxHash: prevTxHash, Logs: logs, BlockNumber: logs[0].BlockNumber}
		}
	}()
	return groupingCh
}
