package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	evmmetrics "github.com/monolit-network/xlayer-indexer/evm/pkg/metrics"
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
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

const (
	ringBufferSize             = 128
	defaultLiveThreshold       = int64(96)
	defaultHistoricalBatch     = int64(200)
	defaultMaxInFlight         = 4
	defaultCatchupChunkBatches = 256
	defaultMetricsAddr         = ":9090"
	retryDelay                 = 5 * time.Second
	latestBlockPollInterval    = 5 * time.Second
	metricsPollInterval        = 5 * time.Second
)

var (
	excludeProcessor = pflag.StringSlice("exclude-processor", nil, "exclude processor: basic | polymarket")
	startBlockFlag   = pflag.String("start-block", "", "explicit start block (overrides db + env)")
	chainFlag        = pflag.String("chain", "", "chain name")
)

type indexerConfig struct {
	Chain               models.Chain
	LiveThreshold       int64
	HistoricalBatch     int64
	MaxInFlight         int
	CatchupChunkBatches int
	MetricsAddr         string
}

func main() {
	pflag.Parse()
	_ = godotenv.Overload()

	chain, ok := models.ValidateChain(*chainFlag)
	if !ok {
		log.Fatalf("unknown chain name: %s", chainFlag)
	}

	logger := initLogger(chain)

	cfg := loadConfig(logger, chain)
	evmmetrics.Register()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)
	var shuttingDown atomic.Bool
	go func() {
		for {
			<-sigCh
			if shuttingDown.CompareAndSwap(false, true) {
				logger.Info("received signal, draining current in-flight work and flushing")
				cancel()
				continue
			}
			logger.Warn("received second signal, forcing shutdown")
			os.Exit(1)
		}
	}()

	dbClient, err := postgres.NewPostgresClientFromEnv()
	if err != nil {
		logger.Fatal("error creating postgres client", zap.Error(err))
	}
	defer dbClient.Close()
	if err := dbClient.Querier().CreateEvmParserRegistryTable(ctx); err != nil {
		logger.Fatal("error creating evm parser registry table", zap.Error(err))
	}

	clickClient, err := clickhouse.NewClickhouseClientFromEnv(logger)
	if err != nil {
		logger.Fatal("error creating clickhouse client", zap.Error(err))
	}
	defer clickClient.Close()
	if err := clickClient.CreateSchemaEvm(ctx); err != nil {
		logger.Fatal("error creating evm schema", zap.Error(err))
	}
	if err := clickClient.CreateEvmSwapEventsTable(ctx); err != nil {
		logger.Fatal("error creating evm swap events table", zap.Error(err))
	}
	if err := clickClient.CreateEvmTransferEventsTable(ctx); err != nil {
		logger.Fatal("error creating evm transfer events table", zap.Error(err))
	}

	liveClient, liveCloser := buildLiveClient(logger, cfg.Chain)
	defer liveCloser()
	histClient, histCloser := buildHistoricalClient(logger, cfg.Chain)
	defer histCloser()
	router := &routingClient{live: liveClient, historical: histClient}

	eventNotifier := buildNotifier(logger, string(cfg.Chain))
	defer eventNotifier.Close()

	procs := buildProcessors(logger, clickClient, dbClient, liveClient, eventNotifier, cfg.Chain)
	for _, p := range procs {
		defer p.Stop()
	}

	// Sync processing keeps polymarketProcessor strictly ordered: with
	// parallel dispatch of multiple batches, BatchGetBlocksInfo fetches run
	// concurrently but the single aggregator worker applies batches in
	// enqueue order, so tx arrive at the processor in block-asc order.
	agg := processors.NewProcessorAggregator(procs, router, logger, processors.WithSyncProcessing(true))
	defer agg.Stop()

	latest := fetchLatestOnChainUntilSuccess(ctx, logger, liveClient)
	if ctx.Err() != nil {
		return
	}
	lastProcessed := determineStartBlock(ctx, logger, clickClient, dbClient, cfg.Chain, latest)

	ring := newHashRing(ringBufferSize)
	initialRing, err := fetchRingBufferUntilSuccess(ctx, logger, liveClient, lastProcessed, ringBufferSize)
	if err != nil {
		return
	}
	ring.Reset(initialRing)

	state := newIndexerState(lastProcessed, latest, ring)
	evmmetrics.ObserveProgress(string(cfg.Chain), lastProcessed, latest)

	logger.Info("indexer starting",
		zap.String("chain", string(cfg.Chain)),
		zap.Int64("last_processed", lastProcessed),
		zap.Int64("latest_on_chain", latest),
		zap.Int("ring_buffer_init", len(initialRing)),
		zap.Int64("live_threshold", cfg.LiveThreshold),
		zap.Int64("historical_batch", cfg.HistoricalBatch),
		zap.Int("max_in_flight", cfg.MaxInFlight),
		zap.Int("catchup_chunk_batches", cfg.CatchupChunkBatches),
		zap.String("metrics_addr", cfg.MetricsAddr),
	)

	deps := workerDeps{
		logger:       logger,
		agg:          agg,
		router:       router,
		liveClient:   liveClient,
		histClient:   histClient,
		state:        state,
		cfg:          cfg,
		shuttingDown: &shuttingDown,
		handleReorgFn: func(ctx context.Context) bool {
			return handleReorgUntilSuccess(ctx, logger, clickClient, dbClient, eventNotifier, liveClient, state, cfg.Chain)
		},
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); serveMetrics(ctx, logger, cfg, state) }()
	wg.Add(1)
	go func() { defer wg.Done(); observeIndexerProgress(ctx, cfg.Chain, state) }()
	wg.Add(1)
	go func() { defer wg.Done(); forwardErrors(ctx, logger, agg, procs) }()
	wg.Add(1)
	go func() { defer wg.Done(); subscribeHeadsLoop(ctx, logger, liveClient, state) }()
	wg.Add(1)
	go func() { defer wg.Done(); pollLatestLoop(ctx, logger, liveClient, state) }()
	wg.Add(1)
	go func() { defer wg.Done(); logClientTimingsLoop(ctx, liveClient, histClient) }()
	wg.Add(1)
	go func() { defer wg.Done(); runWorker(ctx, deps) }()

	wg.Wait()
	agg.FlushSync()
	logger.Info("indexer stopped")
}

func initLogger(chain models.Chain) *zap.Logger {
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "./logs"
	}

	logDir = path.Join(logDir, string(chain))

	return ctxlog.NewCmdLogger(ctxlog.WithLogDir(logDir))
}

func logClientTimingsLoop(ctx context.Context, liveClient, histClient *evmclient.Client) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			liveClient.LogTimings()
			if histClient != liveClient {
				histClient.LogTimings()
			}
			return
		case <-ticker.C:
			liveClient.LogTimings()
			if histClient != liveClient {
				histClient.LogTimings()
			}
		}
	}
}

func loadConfig(logger *zap.Logger, chain models.Chain) indexerConfig {
	cfg := indexerConfig{
		Chain:               chain,
		LiveThreshold:       defaultLiveThreshold,
		HistoricalBatch:     defaultHistoricalBatch,
		MaxInFlight:         defaultMaxInFlight,
		CatchupChunkBatches: defaultCatchupChunkBatches,
		MetricsAddr:         defaultMetricsAddr,
	}
	if v := getEnvForChain(chain, "EVM_INDEXER_LIVE_THRESHOLD"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n <= 0 {
			logger.Fatal("invalid EVM_INDEXER_LIVE_THRESHOLD", zap.String("value", v))
		}
		cfg.LiveThreshold = n
	}
	if v := getEnvForChain(chain, "EVM_INDEXER_HISTORICAL_BATCH"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n <= 0 {
			logger.Fatal("invalid EVM_INDEXER_HISTORICAL_BATCH", zap.String("value", v))
		}
		cfg.HistoricalBatch = n
	}
	if v := getEnvForChain(chain, "EVM_INDEXER_MAX_IN_FLIGHT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			logger.Fatal("invalid EVM_INDEXER_MAX_IN_FLIGHT", zap.String("value", v))
		}
		cfg.MaxInFlight = n
	}
	if v := getEnvForChain(chain, "EVM_INDEXER_CATCHUP_CHUNK_BATCHES"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			logger.Fatal("invalid EVM_INDEXER_CATCHUP_CHUNK_BATCHES", zap.String("value", v))
		}
		cfg.CatchupChunkBatches = n
	}
	if v := getEnvForChain(chain, "EVM_INDEXER_METRICS_ADDR"); v != "" {
		cfg.MetricsAddr = v
	}
	return cfg
}

func serveMetrics(ctx context.Context, logger *zap.Logger, cfg indexerConfig, state *indexerState) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		lastProcessed, latest := state.Snapshot()
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"status":"running","chain":%q,"last_processed":%d,"latest":%d,"lag":%d,"timestamp":%q}`+"\n",
			cfg.Chain,
			lastProcessed,
			latest,
			max(latest-lastProcessed, 0),
			time.Now().UTC().Format(time.RFC3339),
		)
	})

	server := &http.Server{Addr: cfg.MetricsAddr, Handler: mux}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Warn("metrics server shutdown failed", zap.Error(err))
		}
	}()

	logger.Info("starting metrics server", zap.String("addr", cfg.MetricsAddr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("metrics server stopped with error", zap.Error(err))
	}
}

func observeIndexerProgress(ctx context.Context, chain models.Chain, state *indexerState) {
	ticker := time.NewTicker(metricsPollInterval)
	defer ticker.Stop()
	for {
		lastProcessed, latest := state.Snapshot()
		evmmetrics.ObserveProgress(string(chain), lastProcessed, latest)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func buildLiveClient(logger *zap.Logger, chain models.Chain) (*evmclient.Client, func()) {
	logger = logger.Named("live-client")
	rpcURL := envWithFallback(chain, "EVM_RPC_URL")
	wsURL := envWithFallback(chain, "EVM_WS_URL")
	if rpcURL == "" {
		logger.Fatal("EVM_INDEXER_LIVE_RPC_URL / EVM_RPC_URL is not set")
	}
	if wsURL == "" {
		logger.Fatal("EVM_INDEXER_LIVE_WS_URL / EVM_WS_URL is not set")
	}
	logger.Info("live rpc url", zap.String("url", rpcURL), zap.String("ws url", wsURL))
	ethClient, rpcClient := dialHTTP(logger, rpcURL)
	wsClient, err := ethclient.Dial(wsURL)
	if err != nil {
		logger.Fatal("error dialing live ws", zap.Error(err))
	}
	opts := evmclient.EvmClientOptionsFromEnv(models.Chain(chain))
	client := evmclient.NewClient(ethClient, wsClient, rpcClient, logger, models.Chain(chain), opts...)
	return client, func() { ethClient.Close(); wsClient.Close() }
}

func buildHistoricalClient(logger *zap.Logger, chain models.Chain) (*evmclient.Client, func()) {
	logger = logger.Named("historical-client")
	rpcURL := envWithFallback(chain, "EVM_RPC_URL_HIST", "EVM_RPC_URL")
	if rpcURL == "" {
		logger.Fatal("EVM_RPC_URL_HIST / EVM_RPC_URL is not set")
	}
	logger.Info("historical rpc url", zap.String("url", rpcURL))
	ethClient, rpcClient := dialHTTP(logger, rpcURL)
	opts := evmclient.EvmClientOptionsFromEnv(models.Chain(chain))
	client := evmclient.NewClient(ethClient, nil, rpcClient, logger, models.Chain(chain), opts...)
	return client, func() { ethClient.Close() }
}

func dialHTTP(logger *zap.Logger, url string) (*ethclient.Client, *rpc.Client) {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}
	httpClient := &http.Client{Transport: transport, Timeout: 300 * time.Second}
	rpcClient, err := rpc.DialOptions(context.Background(), url, rpc.WithHTTPClient(httpClient))
	if err != nil {
		logger.Fatal("error dialing rpc", zap.String("url", url), zap.Error(err))
	}
	return ethclient.NewClient(rpcClient), rpcClient
}

func buildNotifier(logger *zap.Logger, chain string) *notifier.Notifier {
	topics := []notifier.Topic{
		notifier.TopicOf[models.SwapEvent](string(sharedmodels.NotifierEventEVMSwap)),
		notifier.TopicOf[models.TransferEvent](string(sharedmodels.NotifierEventEVMTransfer)),
		notifier.TopicOf[models.DefiEvent](string(sharedmodels.NotifierEventEVMDefi)),
		notifier.TopicOf[sharedmodels.PolymarketMarketEvent](string(sharedmodels.NotifierEventCondPrepared)),
		notifier.TopicOf[sharedmodels.PolymarketMarketEvent](string(sharedmodels.NotifierEventCondResolved)),
		notifier.TopicOf[sharedmodels.PolymarketMarketEvent](string(sharedmodels.NotifierEventAdapterUMA)),
		notifier.TopicOf[models.PolymarketOrderEventNew](string(sharedmodels.NotifierEventOrder)),
		notifier.TopicOf[models.ReorgEvent](string(sharedmodels.NotifierEventReorg)),
	}
	n, err := notifier.NewNATSFromEnv(logger, topics...)
	if err != nil {
		logger.Fatal("error initializing nats notifier", zap.Error(err))
	}
	return n
}

func buildProcessors(
	logger *zap.Logger,
	clickClient db.DBClick,
	dbClient db.Client,
	evmClient *evmclient.Client,
	eventNotifier *notifier.Notifier,
	chain models.Chain,
) []processors.IndexerProcessor {
	procs := make([]processors.IndexerProcessor, 0, 2)
	if !slices.Contains(*excludeProcessor, "basic") {
		procs = append(procs, basicProcessor.NewProcessor(logger, dbClient, clickClient, evmClient, eventNotifier, string(chain)))
	}
	if strings.EqualFold(string(chain), string(models.ChainPolygon)) && !slices.Contains(*excludeProcessor, "polymarket") {
		procs = append(procs, polymarketProcessor.NewProcessor(evmClient, clickClient, eventNotifier, dbClient, logger))
	}
	if len(procs) == 0 {
		logger.Fatal("no processors enabled")
	}
	return procs
}

func forwardErrors(ctx context.Context, logger *zap.Logger, agg *processors.ProcessorAggregator, procs []processors.IndexerProcessor) {
	for _, p := range procs {
		p := p
		go func() {
			for err := range p.Errors() {
				if ctx.Err() != nil {
					return
				}
				logger.Error("processor error", zap.Error(err))
			}
		}()
	}
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-agg.Errors():
			logger.Fatal("aggregator error", zap.Error(err))
		}
	}
}

func subscribeHeadsLoop(ctx context.Context, logger *zap.Logger, liveClient *evmclient.Client, state *indexerState) {
	for ctx.Err() == nil {
		sub, err := liveClient.SubscribeHeads(ctx)
		if err != nil {
			logger.Warn("subscribe heads failed, retrying", zap.Error(err))
			if sleepCtx(ctx, retryDelay) != nil {
				return
			}
			continue
		}
		logger.Info("subscribed to new heads")
		consumeHeads(ctx, sub, state)
		if ctx.Err() != nil {
			return
		}
		logger.Warn("heads channel closed, resubscribing")
	}
}

func consumeHeads(ctx context.Context, sub <-chan *ethtypes.Header, state *indexerState) {
	for {
		select {
		case <-ctx.Done():
			return
		case head, ok := <-sub:
			if !ok {
				return
			}
			if head == nil || head.Number == nil {
				continue
			}
			state.UpdateLatest(head.Number.Int64())
		}
	}
}

func pollLatestLoop(ctx context.Context, logger *zap.Logger, liveClient *evmclient.Client, state *indexerState) {
	ticker := time.NewTicker(latestBlockPollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			n, err := liveClient.GetEthClient().BlockNumber(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				logger.Warn("poll latest block failed", zap.Error(err))
				continue
			}
			state.UpdateLatest(int64(n))
		}
	}
}

func determineStartBlock(
	ctx context.Context,
	logger *zap.Logger,
	clickClient db.DBClick,
	dbClient db.Client,
	chain models.Chain,
	latestOnChain int64,
) int64 {
	if *startBlockFlag != "" {
		n, err := strconv.ParseInt(*startBlockFlag, 10, 64)
		if err != nil || n < 0 {
			logger.Fatal("invalid --start-block", zap.String("value", *startBlockFlag))
		}
		return n
	}
	if v := getEnvForChain(chain, "EVM_INDEXER_START_BLOCK"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n < 0 {
			logger.Fatal("invalid EVM_INDEXER_START_BLOCK", zap.String("value", v))
		}
		return n
	}
	n, err := latestDBBlock(ctx, clickClient, chain)
	if err != nil {
		logger.Fatal("error reading latest block from db", zap.Error(err))
	}
	if n == 0 {
		tip := latestOnChain - 1
		if tip < 0 {
			tip = 0
		}
		logger.Info("no prior progress in db, starting from tip", zap.Int64("block", tip))
		return tip
	}
	if err := rollbackStores(ctx, clickClient, dbClient, chain, n, logger); err != nil {
		logger.Fatal("error rolling back last db block before startup", zap.Int64("block", n), zap.Error(err))
	}
	lastProcessed := n - 1
	logger.Info("rolled back last db block before startup",
		zap.Int64("rolled_back_from", n),
		zap.Int64("resume_from", n),
		zap.Int64("last_processed", lastProcessed),
	)
	return lastProcessed
}

func latestDBBlock(ctx context.Context, clickClient db.DBClick, chain models.Chain) (int64, error) {
	sources := []string{
		"SELECT ifNull(max(block_number), 0) AS mx FROM evm.swap_events WHERE chain = ?",
		"SELECT ifNull(max(block_number), 0) AS mx FROM evm.transfer_events WHERE chain = ?",
		"SELECT ifNull(max(block_number), 0) AS mx FROM evm.defi_events WHERE chain = ?",
	}
	args := []any{chain, chain, chain}
	if strings.EqualFold(string(chain), string(models.ChainPolygon)) {
		sources = append(sources, "SELECT ifNull(max(block_number), 0) AS mx FROM evm.polymarket_order_events")
	}
	query := fmt.Sprintf("SELECT toInt64(ifNull(max(mx), 0)) FROM (%s)", strings.Join(sources, " UNION ALL "))

	rows, err := clickClient.Query(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, rows.Err()
	}
	var maxBlock int64
	if err := rows.Scan(&maxBlock); err != nil {
		return 0, err
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return maxBlock, nil
}

func fetchLatestOnChainUntilSuccess(ctx context.Context, logger *zap.Logger, client *evmclient.Client) int64 {
	for ctx.Err() == nil {
		n, err := client.GetEthClient().BlockNumber(ctx)
		if err == nil {
			return int64(n)
		}
		logger.Warn("initial latest block fetch failed, retrying", zap.Error(err))
		if sleepCtx(ctx, retryDelay) != nil {
			return 0
		}
	}
	return 0
}

func envWithFallback(chain models.Chain, keys ...string) string {
	for _, key := range keys {
		if v := getEnvForChain(chain, key); v != "" {
			return v
		}
	}
	return ""
}

func getEnvForChain(chain models.Chain, key string) string {
	if v := os.Getenv(fmt.Sprintf("%s_%s", key, strings.ToUpper(string(chain)))); v != "" {
		return v
	}
	return os.Getenv(key)
}
