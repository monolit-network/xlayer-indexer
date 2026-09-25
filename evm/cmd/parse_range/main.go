package main

import (
	"sort"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
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
	"github.com/monolit-network/xlayer-indexer/util/ctxlog"
	"github.com/monolit-network/xlayer-indexer/util/iterators"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

var fillGap = pflag.Bool("fill-gap", false, "parse range from latest_block_in_db to current_block_on_chain")
var startBlock = pflag.String("start-block", "", "start block")
var blocksFile = pflag.String("blocks-file", "", "file with block numbers (one per line); overrides start/end block range")
var endBlock = pflag.String("end-block", "", "end block")
var excludeProcessor = pflag.StringSlice("exclude-processor", []string{}, "exclude processor. available processors: basic, polymarket")

var (
	defaultParseRangeSize = 1000
)

type parseRangeConfig struct {
	RangeSize int
}

const debugResponsePreviewLimit = 16 * 1024

type debugResponsesRoundTripper struct {
	base   http.RoundTripper
	logger *zap.Logger
}

func (t *debugResponsesRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBody, err := readRequestBody(req)
	if err != nil {
		return nil, err
	}

	req2 := req.Clone(req.Context())
	req2.Header = req.Header.Clone()
	req2.Header.Set("Accept-Encoding", "identity")

	resp, err := t.base.RoundTrip(req2)
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return resp, nil
	}

	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	resp.Body = io.NopCloser(bytes.NewReader(respBody))
	resp.ContentLength = int64(len(respBody))
	resp.Header.Set("Content-Length", strconv.Itoa(len(respBody)))

	if shouldLogDebugResponse(reqBody, resp.StatusCode, respBody) {
		t.logger.Warn(
			"rpc debug response",
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()),
			zap.Int("status_code", resp.StatusCode),
			zap.Bool("batch_request", isJSONArray(reqBody)),
			zap.ByteString("request_preview", truncateBytes(reqBody, debugResponsePreviewLimit)),
			zap.ByteString("response_preview", truncateBytes(respBody, debugResponsePreviewLimit)),
			zap.Bool("request_truncated", len(reqBody) > debugResponsePreviewLimit),
			zap.Bool("response_truncated", len(respBody) > debugResponsePreviewLimit),
		)
	}

	return resp, nil
}

func readRequestBody(req *http.Request) ([]byte, error) {
	if req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		defer body.Close()
		return io.ReadAll(body)
	}
	if req.Body == nil {
		return nil, nil
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}

func shouldLogDebugResponse(reqBody []byte, statusCode int, respBody []byte) bool {
	if statusCode < 200 || statusCode >= 300 {
		return true
	}
	if !isJSONArray(reqBody) {
		return false
	}
	if len(bytes.TrimSpace(respBody)) == 0 {
		return true
	}
	return !isJSONArray(respBody)
}

func isJSONArray(body []byte) bool {
	body = bytes.TrimSpace(body)
	return len(body) > 0 && body[0] == '['
}

func truncateBytes(body []byte, limit int) []byte {
	if len(body) <= limit {
		return body
	}
	return body[:limit]
}

func initRpcs(logger *zap.Logger, chain string) (*ethclient.Client, *rpc.Client, error) {
	rpcURL := os.Getenv(fmt.Sprintf("EVM_RPC_URL_%s", strings.ToUpper(chain)))
	if rpcURL == "" {
		return nil, nil, errors.New("EVM_RPC_URL is not set")
	}
	wsURL := os.Getenv(fmt.Sprintf("EVM_WS_URL_%s", strings.ToUpper(chain)))
	if wsURL == "" {
		return nil, nil, errors.New("EVM_WS_URL is not set")
	}
	t := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}
	transport := http.RoundTripper(t)
	if debugResponsesStr := getEnvForChain(chain, "EVM_DEBUG_RESPONSES"); debugResponsesStr != "" && debugResponsesStr != "0" && !strings.EqualFold(debugResponsesStr, "false") {
		logger.Info("enabling rpc debug responses logging")
		transport = &debugResponsesRoundTripper{
			base:   transport,
			logger: logger.Named("rpc-debug-responses"),
		}
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   300 * time.Second,
	}

	rpcClient, err := rpc.DialOptions(context.Background(), rpcURL, rpc.WithHTTPClient(httpClient))
	if err != nil {
		return nil, nil, err
	}
	ethClient := ethclient.NewClient(rpcClient)

	// wsEthClient, err := ethclient.Dial(wsURL)
	// if err != nil {
	// 	return nil, nil, nil, err
	// }

	return ethClient, rpcClient, nil
}

func initProcessors(logger *zap.Logger, dbClickClient db.DBClick, dbClient db.Client, evmClient *evmclient.Client, chain string) ([]processors.IndexerProcessor, chan error) {
	var processors = []processors.IndexerProcessor{}
	errCh := make(chan error, 1024)

	basicProc := basicProcessor.NewProcessor(logger, dbClient, dbClickClient, evmClient, nil, chain)

	if !slices.Contains(*excludeProcessor, "basic") {
		logger.Info("initializing basic processor")
		processors = append(processors, basicProc)

		go func() {
			for err := range basicProc.Errors() {
				errCh <- fmt.Errorf("basic processor error: %w", err)
			}
		}()
	}

	if strings.EqualFold(chain, string(models.ChainPolygon)) && !slices.Contains(*excludeProcessor, "polymarket") {
		logger.Info("initializing polymarket processor")
		polymarketProc := polymarketProcessor.NewProcessor(evmClient, dbClickClient, nil, dbClient, logger)
		processors = append(processors, polymarketProc)

		go func() {
			for err := range polymarketProc.Errors() {
				errCh <- fmt.Errorf("polymarket processor error: %w", err)
			}
		}()
	}

	return processors, errCh
}

func getEnvForChain(chain string, key string) string {
	keyChain := fmt.Sprintf("%s_%s", key, strings.ToUpper(chain))
	if value := os.Getenv(keyChain); value != "" {
		return value
	}
	return os.Getenv(key)
}

func parseRangeConfigFromEnv(chain string) (parseRangeConfig, error) {
	cfg := parseRangeConfig{
		RangeSize: defaultParseRangeSize,
	}

	if value := getEnvForChain(chain, "EVM_PARSE_RANGE_SIZE"); value != "" {
		rangeSize, err := strconv.Atoi(value)
		if err != nil {
			return parseRangeConfig{}, fmt.Errorf("invalid EVM_PARSE_RANGE_SIZE: %w", err)
		}
		if rangeSize <= 0 {
			return parseRangeConfig{}, errors.New("EVM_PARSE_RANGE_SIZE must be > 0")
		}
		cfg.RangeSize = rangeSize
	}

	return cfg, nil
}

func deriveMaxInFlightRanges(evmClient *evmclient.Client) int {
	maxConcurrent := evmClient.GetMaxBlockBatchConcurrentBlocks()
	if evmClient.GetMaxBlockBatchConcurrentReceipts() > maxConcurrent {
		maxConcurrent = evmClient.GetMaxBlockBatchConcurrentReceipts()
	}
	if evmClient.GetMaxBlockBatchConcurrentHeaders() > maxConcurrent {
		maxConcurrent = evmClient.GetMaxBlockBatchConcurrentHeaders()
	}
	if evmClient.GetMaxBlockBatchConcurrentTraces() > maxConcurrent {
		maxConcurrent = evmClient.GetMaxBlockBatchConcurrentTraces()
	}
	if maxConcurrent <= 0 {
		return 1
	}
	return maxConcurrent * 2
}

func main() {
	pflag.Parse()
	logger := ctxlog.NewCmdLogger()
	if err := godotenv.Overload(); err != nil {
		logger.Warn("failed to overload .env", zap.Error(err))
	}
	rangeSizeStr := os.Getenv("EVM_INDEXER_HISTORICAL_BATCH")
	if rangeSizeStr != "" {
		rangeSize, err := strconv.Atoi(rangeSizeStr)
		if err != nil {
			logger.Fatal("error parsing EVM_INDEXER_HISTORICAL_BATCH", zap.Error(err))
		}
		defaultParseRangeSize = rangeSize
	}

	dbClient, err := postgres.NewPostgresClientFromEnv()
	if err != nil {
		logger.Fatal("error creating database client: %v", zap.Error(err))
	}
	if err := dbClient.Querier().CreateEvmParserRegistryTable(context.Background()); err != nil {
		logger.Fatal("error creating evm parser registry table: %v", zap.Error(err))
	}
	defer dbClient.Close()
	logger.Info("database client initialized successfully")

	clickClient, err := clickhouse.NewClickhouseClientFromEnv(logger)
	if err != nil {
		logger.Fatal("error creating clickhouse client", zap.Error(err))
	}
	defer clickClient.Close()
	if err := clickClient.CreateSchemaEvm(context.Background()); err != nil {
		logger.Fatal("error creating evm schema in clickhouse", zap.Error(err))
	}
	if err := clickClient.CreateEvmSwapEventsTable(context.Background()); err != nil {
		logger.Fatal("error creating swap_events table in clickhouse", zap.Error(err))
	}
	if err := clickClient.CreateEvmTransferEventsTable(context.Background()); err != nil {
		logger.Fatal("error creating transfer_events table in clickhouse", zap.Error(err))
	}

	chain := os.Getenv("EVM_CHAIN")
	if chain == "" {
		logger.Fatal("EVM_CHAIN is not set")
	}
	parseCfg, err := parseRangeConfigFromEnv(chain)
	if err != nil {
		logger.Fatal("invalid parse range env config", zap.Error(err))
	}

	ethClient, rpcClient, err := initRpcs(logger, chain)
	if err != nil {
		logger.Fatal("error initializing rpcs", zap.Error(err))
	}
	defer ethClient.Close()

	evmClient := evmclient.NewClient(ethClient, nil, rpcClient, logger, models.Chain(chain), evmclient.EvmClientOptionsFromEnv(models.Chain(chain))...)
	maxInFlightRanges := deriveMaxInFlightRanges(evmClient)
	logger.Info(
		"parse range config",
		zap.Int("range_size", parseCfg.RangeSize),
		zap.Int("max_in_flight_ranges", maxInFlightRanges),
	)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	procs, procErrCh := initProcessors(logger, clickClient, dbClient, evmClient, chain)
	for _, proc := range procs {
		defer proc.Stop()
	}

	aggregator := processors.NewProcessorAggregator(procs, evmClient, logger)
	defer aggregator.Stop()

	latestBlockOnChain, err := ethClient.BlockNumber(ctx)
	if err != nil {
		logger.Fatal("error getting latest block on chain", zap.Error(err))
	}

	var preset []iterators.BigSubrange
	var fromBlock, toBlock *big.Int
	if *blocksFile != "" {
		data, ferr := os.ReadFile(*blocksFile)
		if ferr != nil {
			logger.Fatal("cannot read blocks-file", zap.Error(ferr))
		}
		nums := make([]int64, 0, 16384)
		for _, ln := range strings.Split(string(data), "\n") {
			ln = strings.TrimSpace(ln)
			if ln == "" {
				continue
			}
			v, perr := strconv.ParseInt(ln, 10, 64)
			if perr != nil {
				logger.Fatal("bad block number in blocks-file", zap.String("line", ln), zap.Error(perr))
			}
			nums = append(nums, v)
		}
		if len(nums) == 0 {
			logger.Fatal("blocks-file is empty")
		}
		sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
		lo := nums[0]
		prev := nums[0]
		for _, v := range nums[1:] {
			if v == prev || v == prev+1 {
				prev = v
				continue
			}
			preset = append(preset, iterators.BigSubrange{Start: big.NewInt(lo), End: big.NewInt(prev)})
			lo, prev = v, v
		}
		preset = append(preset, iterators.BigSubrange{Start: big.NewInt(lo), End: big.NewInt(prev)})
		fromBlock = big.NewInt(nums[0])
		toBlock = big.NewInt(nums[len(nums)-1])
		logger.Info("blocks-file mode", zap.Int("blocks", len(nums)), zap.Int("subranges", len(preset)))
	} else {
		var err error
		fromBlock, err = parseBlockNumber(*startBlock)
		if err != nil {
			logger.Fatal("invalid start-block", zap.Error(err))
		}
		toBlock, err = parseBlockNumber(*endBlock)
		if err != nil {
			logger.Fatal("invalid end-block", zap.Error(err))
		}
		if toBlock.Cmp(fromBlock) < 0 {
			logger.Fatal("end_block must be >= start_block")
		}
	}

	if toBlock.Cmp(new(big.Int).SetUint64(latestBlockOnChain)) > 0 {
		toBlock = new(big.Int).SetUint64(latestBlockOnChain)
	}

	go func() {
		parseRange(ctx, logger, aggregator, fromBlock, toBlock, parseCfg, maxInFlightRanges, preset)
		logger.Info("range parsed successfully")
		cancel()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	go func() {
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

	select {
	case <-sigCh:
		logger.Info("received signal, shutting down gracefully")
		cancel()
	case <-ctx.Done():
		logger.Info("context canceled, shutting down gracefully")
		cancel()
	}
	logger.Info("received signal, shutting down gracefully")
	cancel()
}

func parseRange(ctx context.Context, logger *zap.Logger, aggregator *processors.ProcessorAggregator, startBlock *big.Int, endBlock *big.Int, cfg parseRangeConfig, maxInFlightRanges int, preset []iterators.BigSubrange) {
	start := time.Now()
	var processedBlocks atomic.Int64

	type rangeProcessResult struct {
		batch    iterators.BigSubrange
		duration time.Duration
		err      error
	}

	pending := make([]iterators.BigSubrange, 0)
	if len(preset) > 0 {
		pending = preset
	} else {
		for _, batch := range iterators.BigRangeSubranges(startBlock, endBlock, cfg.RangeSize) {
			pending = append(pending, batch)
		}
	}

	resultCh := make(chan rangeProcessResult, maxInFlightRanges)
	inFlight := 0

	logger.Info(
		"parsing range",
		zap.Int64("start_block", startBlock.Int64()),
		zap.Int64("end_block", endBlock.Int64()),
		zap.Int("range_size", cfg.RangeSize),
		zap.Int("max_in_flight_ranges", maxInFlightRanges),
		zap.Int("num_ranges", len(pending)),
	)
	startTime := time.Now()

	dispatch := func(batch iterators.BigSubrange) {
		start := time.Now()
		doneCh := aggregator.ProcessBlockRangeWithResult(ctx, batch.Start, batch.End)
		logger.Info("block range queued", zap.Int64("start_block", batch.Start.Int64()), zap.Int64("end_block", batch.End.Int64()), zap.Duration("duration", time.Since(start)))

		inFlight++
		go func(batch iterators.BigSubrange, queuedAt time.Time, doneCh <-chan error) {
			err := <-doneCh
			resultCh <- rangeProcessResult{
				batch:    batch,
				duration: time.Since(queuedAt),
				err:      err,
			}
		}(batch, start, doneCh)
	}

	for len(pending) > 0 || inFlight > 0 {
		for len(pending) > 0 && inFlight < maxInFlightRanges {
			batch := pending[0]
			pending = pending[1:]
			dispatch(batch)
		}

		if inFlight == 0 {
			break
		}

		select {
		case <-ctx.Done():
			logger.Info("range parsing canceled", zap.Error(ctx.Err()))
			return
		case res := <-resultCh:
			inFlight--
			if res.err != nil {
				if errors.Is(res.err, context.Canceled) || ctx.Err() != nil {
					logger.Info("block range canceled", zap.Int64("start_block", res.batch.Start.Int64()), zap.Int64("end_block", res.batch.End.Int64()), zap.Error(res.err))
					return
				}

				logger.Warn("block range failed, requeueing at front",
					zap.Int64("start_block", res.batch.Start.Int64()),
					zap.Int64("end_block", res.batch.End.Int64()),
					zap.Duration("duration", res.duration),
					zap.Error(res.err),
				)
				pending = append([]iterators.BigSubrange{res.batch}, pending...)
				continue
			}

			numberOfBlocks := big.NewInt(0).Sub(res.batch.End, res.batch.Start).Int64() + 1
			totalProcessed := processedBlocks.Add(numberOfBlocks)

			speed := float64(totalProcessed) / time.Since(start).Seconds()

			logger.Info("block range processed", zap.Int64("start_block", res.batch.Start.Int64()), zap.Int64("end_block", res.batch.End.Int64()), zap.Duration("duration", res.duration), zap.Float64("bps", speed))
		}
	}

	logger.Info("full range parsed successfully", zap.Int64("start_block", startBlock.Int64()), zap.Int64("end_block", endBlock.Int64()), zap.Duration("duration", time.Since(startTime)))
}

func parseBlockNumber(s string) (*big.Int, error) {
	// Support decimal and 0x-prefixed hex
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
