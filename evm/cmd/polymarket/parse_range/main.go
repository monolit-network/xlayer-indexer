package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
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

// Special binary to parse only polymarket-related logs

var (
	startBlock = flag.String("start-block", "", "start block")
	endBlock   = flag.String("end-block", "", "end block")

	cpuProfilePath   = flag.String("cpu-profile", "", "write CPU profile to file")
	memProfilePath   = flag.String("mem-profile", "", "write heap profile to file on exit")
	traceProfilePath = flag.String("trace-profile", "", "write execution trace to file")

	blockProfilePath = flag.String("block-profile", "", "write block profile to file on exit")
	mutexProfilePath = flag.String("mutex-profile", "", "write mutex profile to file on exit")

	disableClickhouseWriter = flag.Bool("disable-click-writer", false, "disable clickhouse writer")
)

func main() {
	godotenv.Load()
	flag.Parse()
	logger := ctxlog.NewCmdLogger()

	stopProfiles, err := startProfiles()
	if err != nil {
		logger.Fatal("error starting profiles", zap.Error(err))
	}
	defer func() {
		if err := stopProfiles(); err != nil {
			logger.Error("error stopping profiles", zap.Error(err))
		}
	}()

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

func startProfiles() (func() error, error) {
	cleanupFns := make([]func() error, 0, 5)

	if *cpuProfilePath != "" {
		f, err := os.Create(*cpuProfilePath)
		if err != nil {
			return nil, fmt.Errorf("create cpu profile: %w", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			f.Close()
			return nil, fmt.Errorf("start cpu profile: %w", err)
		}
		cleanupFns = append(cleanupFns, func() error {
			pprof.StopCPUProfile()
			return f.Close()
		})
	}

	if *traceProfilePath != "" {
		f, err := os.Create(*traceProfilePath)
		if err != nil {
			return nil, fmt.Errorf("create trace profile: %w", err)
		}
		if err := trace.Start(f); err != nil {
			f.Close()
			return nil, fmt.Errorf("start trace profile: %w", err)
		}
		cleanupFns = append(cleanupFns, func() error {
			trace.Stop()
			return f.Close()
		})
	}

	if *blockProfilePath != "" {
		runtime.SetBlockProfileRate(1)
		cleanupFns = append(cleanupFns, func() error {
			runtime.SetBlockProfileRate(0)
			return writeNamedProfile("block", *blockProfilePath)
		})
	}

	if *mutexProfilePath != "" {
		runtime.SetMutexProfileFraction(1)
		cleanupFns = append(cleanupFns, func() error {
			runtime.SetMutexProfileFraction(0)
			return writeNamedProfile("mutex", *mutexProfilePath)
		})
	}

	if *memProfilePath != "" {
		cleanupFns = append(cleanupFns, func() error {
			runtime.GC()
			return writeNamedProfile("heap", *memProfilePath)
		})
	}

	return func() error {
		var firstErr error
		for i := len(cleanupFns) - 1; i >= 0; i-- {
			if err := cleanupFns[i](); err != nil && firstErr == nil {
				firstErr = err
			}
		}
		return firstErr
	}, nil
}

func writeNamedProfile(name string, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s profile: %w", name, err)
	}
	defer f.Close()

	profile := pprof.Lookup(name)
	if profile == nil {
		return fmt.Errorf("profile %q is not available", name)
	}

	if err := profile.WriteTo(f, 0); err != nil {
		return fmt.Errorf("write %s profile: %w", name, err)
	}
	return nil
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
	addressesMain := []common.Address{
		models.PolymarketCTFExchangeAddress,
		models.PolymarketCTFExchangeAddressV2,
		models.PolymarketNegRiskCTFExchangeAddress,
		models.PolymarketNegRiskCTFExchangeAddressV2,
		models.PolymarketNegRiskAdapterAddress,
		models.PolymarketConditionalTokensAddress,
	}
	topicsMain := [][]common.Hash{
		{
			models.PolymarketOrderFilledEventSelectorHash,
			models.PolymarketOrderFilledEventV2SelectorHash,
			models.PolymarketOrdersMatchedEventV2SelectorHash,
			models.PolymarketTransferSingleEventSelectorHash,
			models.PolymarketTransferBatchEventSelectorHash,
			models.PolymarketConditionPreparationEventSelectorHash,
			models.PolymarketConditionResolutionEventSelectorHash,
			models.PolymarketPositionSplitEventSelectorHash,
			models.PolymarketPositionMergeEventSelectorHash,
			models.PolymarketPayoutRedemptionEventSelectorHash,
			models.PolymarketNegRiskPayoutRedemptionEventSelectorHash,
		},
	}

	addressesUmaAdapter := models.PolymarketCTFAdaptersAddresses
	topicsUmaAdapter := [][]common.Hash{
		{
			models.PolymarketUmaCtfAdapterQuestionInitializedEventSelectorHash,
			models.PolymarketUmaCtfAdapterQuestionResetEventSelectorHash,
			models.PolymarketUmaCtfAdapterQuestionResolvedEventSelectorHash,
			models.PolymarketUmaCtfAdapterQuestionPausedEventSelectorHash,
			models.PolymarketUmaCtfAdapterQuestionFlaggedEventSelectorHash,
		},
	}

	addressesUmaOO := []common.Address{
		models.UmaOptimisticOracleV2Address,
		models.UmaPolymarketManagedOptimisticOracleV2Address,
	}
	topicsUmaOO := [][]common.Hash{
		{
			models.UmaOptimisticOracleV2RequestPriceEventSelectorHash,
			models.UmaOptimisticOracleV2ProposePriceEventSelectorHash,
			models.UmaOptimisticOracleV2DisputePriceEventSelectorHash,
			models.UmaOptimisticOracleV2SettleEventSelectorHash,
		},
	}

	topicsFeeModule := [][]common.Hash{
		{
			models.PolymarketFeeRefundedEventSelectorHash,
			models.PolymarketNegRiskFeeRefundedEventSelectorHash,
			models.PolymarketNegRiskPositionsConvertedEventSelectorHash,
		},
	}

	negRiskAdapterTopics := [][]common.Hash{
		{
			models.PolymarketNegRiskPositionSplitEventSelectorHash,
			models.PolymarketNegRiskPositionMergeEventSelectorHash,
		},
	}

	allAddrs := append(addressesMain, addressesUmaAdapter...)
	allAddrs = append(allAddrs, addressesUmaOO...)
	allTopics := append(topicsMain[0], topicsUmaAdapter[0]...)
	allTopics = append(allTopics, topicsUmaOO[0]...)

	noAddressTopics := append(topicsFeeModule[0], negRiskAdapterTopics[0]...)

	logsCh := evmClient.StreamGetLogs(
		ctx,
		[][]common.Address{allAddrs, {}},
		[][][]common.Hash{{allTopics}, {noAddressTopics}},
		startBlock,
		endBlock,
	)

	latestBlockMark := uint64(0)
	prevTime := time.Now()
	blockMarkSize := uint64(10000)

	groupingCh := groupLogsIntoTxs(logsCh)

	first := true
	for group := range groupingCh {
		if group.Err != nil {
			logger.Error("error streaming logs", zap.Error(group.Err))
			return group.Err
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

		if err := processor.ProcessLogsFromTx(group.Logs); err != nil {
			logger.Error("error processing logs from tx", zap.String("hash", group.TxHash.Hex()), zap.Uint64("block_number", group.BlockNumber), zap.Error(err))
			return fmt.Errorf("error processing logs from tx: %w", err)
		}
	}

	logger.Info("stopped streaming logs")
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
