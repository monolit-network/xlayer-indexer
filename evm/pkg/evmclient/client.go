package evmclient

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/util/promise"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/zap"

	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/sync/semaphore"
	"golang.org/x/time/rate"
)

type RequiredDataTypes uint8

const (
	RequiredDataTypesTx          RequiredDataTypes = 1 << iota
	RequiredDataTypesReceipt     RequiredDataTypes = 1 << iota
	RequiredDataTypesBlockHeader RequiredDataTypes = 1 << iota
	RequiredDataTypesTraces      RequiredDataTypes = 1 << iota
)

func (b RequiredDataTypes) Length() int {
	return bits.OnesCount8(uint8(b))
}

func (b RequiredDataTypes) Has(requiredDataType RequiredDataTypes) bool {
	return (b & requiredDataType) != 0
}

func (b RequiredDataTypes) Add(requiredDataType RequiredDataTypes) RequiredDataTypes {
	return b | requiredDataType
}

func (b RequiredDataTypes) Remove(requiredDataType RequiredDataTypes) RequiredDataTypes {
	return b & ^requiredDataType
}

type Client struct {
	logger *zap.Logger

	ethClient   *ethclient.Client
	wsEthClient *ethclient.Client
	rpcClient   *rpc.Client

	requestLimiter  *rate.Limiter
	requestCount    atomic.Int64
	rpcCallCount    atomic.Int64
	retryCount      atomic.Int64
	batchSplitCount atomic.Int64
	timings         *clientTimings

	decimalsCache      *lru.Cache[common.Address, uint8]
	decimalsErrorCache *lru.Cache[common.Address, error]
	factoryCache       *lru.Cache[common.Address, common.Address]
	isContractCache    *lru.Cache[common.Address, bool]

	maxBlockBatchConcurrentBlocks   int
	maxBlockBatchConcurrentTraces   int
	maxBlockBatchConcurrentReceipts int
	maxBlockBatchConcurrentHeaders  int

	maxBlockBatchSizeReceipts int
	maxBlockBatchSizeTraces   int
	maxBlockBatchSizeHeaders  int
	maxBlockBatchSizeBlocks   int

	blocksSem        *semaphore.Weighted
	blockTracesSem   *semaphore.Weighted
	blockHeadersSem  *semaphore.Weighted
	blockReceiptsSem *semaphore.Weighted

	filterLogsBatchSize int

	Chain                     models.Chain
	ChainID                   models.ChainId
	WrappedNativeTokenAddress common.Address
}

func (c *Client) GetEthClient() *ethclient.Client {
	return c.ethClient
}

func (c *Client) GetRpcClient() *rpc.Client {
	return c.rpcClient
}

func (c *Client) GetWsClient() *ethclient.Client {
	return c.wsEthClient
}

func (c *Client) GetMaxBlockBatchConcurrentBlocks() int {
	return c.maxBlockBatchConcurrentBlocks
}

func (c *Client) GetMaxBlockBatchConcurrentReceipts() int {
	return c.maxBlockBatchConcurrentReceipts
}

func (c *Client) GetMaxBlockBatchConcurrentHeaders() int {
	return c.maxBlockBatchConcurrentHeaders
}

func (c *Client) GetMaxBlockBatchConcurrentTraces() int {
	return c.maxBlockBatchConcurrentTraces
}

type Option func(*Client)

type StatsSnapshot struct {
	RequestCount    int64
	RPCCallCount    int64
	RetryCount      int64
	BatchSplitCount int64
}

func WithMaxBlockBatchConcurrentBlocks(maxConcurrent int) Option {
	return func(c *Client) {
		if maxConcurrent <= 0 {
			c.blocksSem = nil
			return
		}
		c.maxBlockBatchConcurrentBlocks = maxConcurrent
		c.blocksSem = semaphore.NewWeighted(int64(maxConcurrent))
	}
}

func WithMaxBlockBatchConcurrentReceipts(maxConcurrent int) Option {
	return func(c *Client) {
		if maxConcurrent <= 0 {
			c.blockReceiptsSem = nil
			return
		}
		c.maxBlockBatchConcurrentReceipts = maxConcurrent
		c.blockReceiptsSem = semaphore.NewWeighted(int64(maxConcurrent))
	}
}

func WithMaxBlockBatchConcurrentHeaders(maxConcurrent int) Option {
	return func(c *Client) {
		if maxConcurrent <= 0 {
			c.blockHeadersSem = nil
			return
		}
		c.maxBlockBatchConcurrentHeaders = maxConcurrent
		c.blockHeadersSem = semaphore.NewWeighted(int64(maxConcurrent))
	}
}

func WithMaxBlockBatchConcurrentTraces(maxConcurrent int) Option {
	return func(c *Client) {
		if maxConcurrent <= 0 {
			c.maxBlockBatchConcurrentTraces = maxConcurrent
			c.blockTracesSem = nil
			return
		}
		c.maxBlockBatchConcurrentTraces = maxConcurrent
		c.blockTracesSem = semaphore.NewWeighted(int64(maxConcurrent))
	}
}

func WithMaxBlockBatchSizeReceipts(maxBatchSize int) Option {
	return func(c *Client) {
		c.maxBlockBatchSizeReceipts = maxBatchSize
	}
}

func WithMaxBlockBatchSizeTraces(maxBatchSize int) Option {
	return func(c *Client) {
		c.maxBlockBatchSizeTraces = maxBatchSize
	}
}

func WithMaxBlockBatchSizeHeaders(maxBatchSize int) Option {
	return func(c *Client) {
		c.maxBlockBatchSizeHeaders = maxBatchSize
	}
}

func WithMaxBlockBatchSizeBlocks(maxBatchSize int) Option {
	return func(c *Client) {
		c.maxBlockBatchSizeBlocks = maxBatchSize
	}
}
func WithRPCLimiter(limiter *rate.Limiter) Option {
	return func(c *Client) {
		c.requestLimiter = limiter
	}
}

func WithFilterLogsBatchSize(size int) Option {
	return func(c *Client) {
		c.filterLogsBatchSize = size
	}
}

func getEnvForChain(chain models.Chain, key string) string {
	keyChain := fmt.Sprintf("%s_%s", key, strings.ToUpper(string(chain)))
	if os.Getenv(keyChain) != "" {
		return os.Getenv(keyChain)
	}
	return os.Getenv(key)
}

func EvmClientOptionsFromEnv(chain models.Chain) []Option {
	opts := []Option{}

	blocksMaxConcurrentStr := getEnvForChain(chain, "EVM_BLOCK_BATCH_MAX_CONCURRENT_BLOCKS")
	if blocksMaxConcurrentStr != "" {
		blocksMaxConcurrent, err := strconv.Atoi(blocksMaxConcurrentStr)
		if err != nil {
			panic(fmt.Sprintf("error parsing EVM_BLOCK_BATCH_MAX_CONCURRENT_BLOCKS: %v", err))
		}
		opts = append(opts, WithMaxBlockBatchConcurrentBlocks(blocksMaxConcurrent))
	}

	blockReceiptsMaxConcurrentStr := getEnvForChain(chain, "EVM_BLOCK_BATCH_MAX_CONCURRENT_RECEIPTS")
	if blockReceiptsMaxConcurrentStr != "" {
		blockReceiptsMaxConcurrent, err := strconv.Atoi(blockReceiptsMaxConcurrentStr)
		if err != nil {
			panic(fmt.Sprintf("error parsing EVM_BLOCK_BATCH_MAX_CONCURRENT_RECEIPTS: %v", err))
		}
		opts = append(opts, WithMaxBlockBatchConcurrentReceipts(blockReceiptsMaxConcurrent))
	}

	blockTracesMaxConcurrentStr := getEnvForChain(chain, "EVM_BLOCK_BATCH_MAX_CONCURRENT_TRACES")
	if blockTracesMaxConcurrentStr != "" {
		blockTracesMaxConcurrent, err := strconv.Atoi(blockTracesMaxConcurrentStr)
		if err != nil {
			panic(fmt.Sprintf("error parsing EVM_BLOCK_BATCH_MAX_CONCURRENT_TRACES: %v", err))
		}
		opts = append(opts, WithMaxBlockBatchConcurrentTraces(blockTracesMaxConcurrent))
	}

	blockHeadersMaxConcurrentStr := getEnvForChain(chain, "EVM_BLOCK_BATCH_MAX_CONCURRENT_HEADERS")
	if blockHeadersMaxConcurrentStr != "" {
		blockHeadersMaxConcurrent, err := strconv.Atoi(blockHeadersMaxConcurrentStr)
		if err != nil {
			panic(fmt.Sprintf("error parsing EVM_BLOCK_BATCH_MAX_CONCURRENT_HEADERS: %v", err))
		}
		opts = append(opts, WithMaxBlockBatchConcurrentHeaders(blockHeadersMaxConcurrent))
	}

	rpcRpsStr := getEnvForChain(chain, "EVM_RPC_RPS")
	if rpcRpsStr != "" {
		rps, err := strconv.Atoi(rpcRpsStr)
		if err != nil {
			panic(fmt.Sprintf("error parsing EVM_RPC_RPS: %v", err))
		}
		burst := rps
		rpcBurstStr := getEnvForChain(chain, "EVM_RPC_BURST")
		if rpcBurstStr != "" && rpcBurstStr != "-1" {
			burst, err = strconv.Atoi(rpcBurstStr)
			if err != nil {
				panic(fmt.Sprintf("error parsing EVM_RPC_BURST: %v", err))
			}
		}
		if rps > 0 {
			limiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(rps)), burst)
			opts = append(opts, WithRPCLimiter(limiter))
		} else if rps != -1 {
			panic("EVM_RPC_RPS must be positive or -1")
		}
	}

	maxBlockBatchSizeReceiptsStr := getEnvForChain(chain, "EVM_MAX_BLOCK_BATCH_SIZE_RECEIPTS")
	if maxBlockBatchSizeReceiptsStr != "" {
		maxBlockBatchSizeReceipts, err := strconv.Atoi(maxBlockBatchSizeReceiptsStr)
		if err != nil {
			panic(fmt.Sprintf("error parsing EVM_MAX_BLOCK_BATCH_SIZE_RECEIPTS: %v", err))
		}
		opts = append(opts, WithMaxBlockBatchSizeReceipts(maxBlockBatchSizeReceipts))
	}

	maxBlockBatchSizeTracesStr := getEnvForChain(chain, "EVM_MAX_BLOCK_BATCH_SIZE_TRACES")
	if maxBlockBatchSizeTracesStr != "" {
		maxBlockBatchSizeTraces, err := strconv.Atoi(maxBlockBatchSizeTracesStr)
		if err != nil {
			panic(fmt.Sprintf("error parsing EVM_MAX_BLOCK_BATCH_SIZE_TRACES: %v", err))
		}
		opts = append(opts, WithMaxBlockBatchSizeTraces(maxBlockBatchSizeTraces))
	}

	maxBlockBatchSizeHeadersStr := getEnvForChain(chain, "EVM_MAX_BLOCK_BATCH_SIZE_HEADERS")
	if maxBlockBatchSizeHeadersStr != "" {
		maxBlockBatchSizeHeaders, err := strconv.Atoi(maxBlockBatchSizeHeadersStr)
		if err != nil {
			panic(fmt.Sprintf("error parsing EVM_MAX_BLOCK_BATCH_SIZE_HEADERS: %v", err))
		}
		opts = append(opts, WithMaxBlockBatchSizeHeaders(maxBlockBatchSizeHeaders))
	}

	maxBlockBatchSizeBlocksStr := getEnvForChain(chain, "EVM_MAX_BLOCK_BATCH_SIZE_BLOCKS")
	if maxBlockBatchSizeBlocksStr != "" {
		maxBlockBatchSizeBlocks, err := strconv.Atoi(maxBlockBatchSizeBlocksStr)
		if err != nil {
			panic(fmt.Sprintf("error parsing EVM_MAX_BLOCK_BATCH_SIZE_BLOCKS: %v", err))
		}
		opts = append(opts, WithMaxBlockBatchSizeBlocks(maxBlockBatchSizeBlocks))
	}

	if filterLogsBatchSizeStr := getEnvForChain(chain, "EVM_FILTER_LOGS_BATCH_SIZE"); filterLogsBatchSizeStr != "" {
		filterLogsBatchSize, err := strconv.Atoi(filterLogsBatchSizeStr)
		if err != nil {
			panic(fmt.Sprintf("error parsing EVM_FILTER_LOGS_BATCH_SIZE: %v", err))
		}
		opts = append(opts, WithFilterLogsBatchSize(filterLogsBatchSize))
	}

	return opts
}

func NewClient(
	ethClient *ethclient.Client,
	wsEthClient *ethclient.Client,
	rpcClient *rpc.Client,
	logger *zap.Logger,
	chain models.Chain,
	opts ...Option,
) *Client {
	decimalsCache, _ := lru.New[common.Address, uint8](8192 * 2 * 2 * 2)
	decimalsErrorCache, _ := lru.New[common.Address, error](8192 * 2 * 2 * 2)
	factoryCache, _ := lru.New[common.Address, common.Address](8192 * 2 * 2)
	isContractCache, _ := lru.New[common.Address, bool](8192 * 2 * 2 * 2)

	chainID, ok := models.ChainToID[chain]
	if !ok {
		panic(fmt.Sprintf("chain id %d not found", chainID))
	}
	nativeTokenAddress, ok := models.ChainIDToWrappedNativeToken[chainID]
	if !ok {
		panic(fmt.Sprintf("native token address for chain id %d not found", chainID))
	}
	c := &Client{
		ethClient:       ethClient,
		wsEthClient:     wsEthClient,
		rpcClient:       rpcClient,
		requestLimiter:  nil,
		requestCount:    atomic.Int64{},
		rpcCallCount:    atomic.Int64{},
		retryCount:      atomic.Int64{},
		batchSplitCount: atomic.Int64{},
		timings:         newClientTimings(),

		decimalsCache:      decimalsCache,
		decimalsErrorCache: decimalsErrorCache,
		factoryCache:       factoryCache,
		isContractCache:    isContractCache,
		logger:             logger,

		Chain:                     chain,
		ChainID:                   chainID,
		WrappedNativeTokenAddress: nativeTokenAddress,

		// defaults
		maxBlockBatchConcurrentBlocks:   1000,
		maxBlockBatchConcurrentTraces:   10,
		maxBlockBatchConcurrentReceipts: 100,
		maxBlockBatchConcurrentHeaders:  1000,

		blocksSem:        semaphore.NewWeighted(1000),
		blockTracesSem:   semaphore.NewWeighted(10),
		blockReceiptsSem: semaphore.NewWeighted(100),
		blockHeadersSem:  semaphore.NewWeighted(1000),

		maxBlockBatchSizeReceipts: 20,
		maxBlockBatchSizeTraces:   5,
		maxBlockBatchSizeHeaders:  500,
		maxBlockBatchSizeBlocks:   100,

		filterLogsBatchSize: 3000,
	}

	for _, opt := range opts {
		opt(c)
	}
	logger.Info("evm client options",
		zap.Any("max_concurrent_block_receipts", c.maxBlockBatchConcurrentReceipts),
		zap.Any("max_concurrent_blocks", c.maxBlockBatchConcurrentBlocks),
		zap.Any("max_concurrent_block_traces", c.maxBlockBatchConcurrentTraces),
		zap.Any("max_concurrent_block_headers", c.maxBlockBatchConcurrentHeaders),
		zap.Any("filter_logs_batch_size", c.filterLogsBatchSize),
		zap.Bool("rpc_rps_limited", c.requestLimiter != nil))

	return c
}

func (c *Client) waitRequests(ctx context.Context, n int64) error {
	if n <= 0 {
		return nil
	}
	if c.requestLimiter != nil {
		if err := c.requestLimiter.WaitN(ctx, int(n)); err != nil {
			return err
		}
	}
	c.requestCount.Add(int64(n))
	return nil
}

func (c *Client) recordRPCCall() {
	c.rpcCallCount.Add(1)
}

func (c *Client) recordRetry() {
	c.retryCount.Add(1)
}

func (c *Client) recordBatchSplit() {
	c.batchSplitCount.Add(1)
}

func (c *Client) withRetry(ctx context.Context, op string, fn func() error) error {
	const maxRetries = math.MaxInt
	baseDelay := 200 * time.Millisecond
	for attempt := 0; ; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		if !isRetryableError(err) || attempt >= maxRetries {
			return err
		}
		c.recordRetry()
		delay := baseDelay << attempt
		if delay > 5*time.Second {
			delay = 5 * time.Second
		}
		c.logger.Warn("retrying operation", zap.String("op", op), zap.Int("attempt", attempt+1), zap.Duration("backoff", delay), zap.Error(err))
		select {
		case <-time.After(delay):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "connection reset by peer") || strings.Contains(msg, "connection refused") {
		return true
	}
	if strings.Contains(msg, "rate limit") || strings.Contains(msg, "too many requests") || strings.Contains(msg, " 429") || strings.Contains(msg, "status 429") {
		return true
	}
	if strings.Contains(msg, "internal server error") || strings.Contains(msg, " 500") || strings.Contains(msg, "status 500") {
		return true
	}
	if strings.Contains(msg, "retryable") || strings.Contains(msg, "unexpected end of json") || strings.Contains(msg, "timed out") {
		return true
	}
	if strings.Contains(msg, "database error") || strings.Contains(msg, "not found") || strings.Contains(msg, "db closed") {
		return true
	}
	return false
}

func IsRetryableError(err error) bool {
	return isRetryableError(err)
}

var possibleBatchErrorMessages = []string{
	"exceeded",
	"too large",
	"eof",
	"unexpected end of json",
	"timed out",
	"cannot unmarshal object into go value of type",
}

func isBatchSplitError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, errMsg := range possibleBatchErrorMessages {
		if strings.Contains(msg, errMsg) {
			return true
		}
	}
	return false
}

func (c *Client) GetRequestCount() int64 {
	return c.requestCount.Load()
}

func (c *Client) GetStatsSnapshot() StatsSnapshot {
	return StatsSnapshot{
		RequestCount:    c.requestCount.Load(),
		RPCCallCount:    c.rpcCallCount.Load(),
		RetryCount:      c.retryCount.Load(),
		BatchSplitCount: c.batchSplitCount.Load(),
	}
}

type BlocksInfo struct {
	Blocks   promise.Promise[[]*types.Block]
	Receipts promise.Promise[[]BlockReceipts]
	Traces   promise.Promise[[]BlockTraces]
	Headers  promise.Promise[[]*types.Header]
}

func (b *BlocksInfo) WaitContext(ctx context.Context, requiredDataTypes RequiredDataTypes) ([]*types.Block, []BlockReceipts, []BlockTraces, []*types.Header, error) {
	var blocks []*types.Block
	var receipts []BlockReceipts
	var traces []BlockTraces
	var headers []*types.Header
	var err error

	if requiredDataTypes.Has(RequiredDataTypesTx) {
		blocks, err = b.Blocks.WaitContext(ctx)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("error getting blocks: %w", err)
		}
	}

	if requiredDataTypes.Has(RequiredDataTypesReceipt) {
		receipts, err = b.Receipts.WaitContext(ctx)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("error getting receipts: %w", err)
		}
	}

	if requiredDataTypes.Has(RequiredDataTypesTraces) {
		traces, err = b.Traces.WaitContext(ctx)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("error getting traces: %w", err)
		}

	}

	if requiredDataTypes.Has(RequiredDataTypesBlockHeader) {
		headers, err = b.Headers.WaitContext(ctx)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("error getting headers: %w", err)
		}
	}
	return blocks, receipts, traces, headers, nil
}

func (c *Client) BatchGetBlocksInfo(ctx context.Context, startBlock *big.Int, endBlock *big.Int, requiredDataTypes RequiredDataTypes) *BlocksInfo {
	txsRequired := requiredDataTypes.Has(RequiredDataTypesTx)
	receiptsRequired := requiredDataTypes.Has(RequiredDataTypesReceipt)
	tracesRequired := requiredDataTypes.Has(RequiredDataTypesTraces)
	headersRequired := requiredDataTypes.Has(RequiredDataTypesBlockHeader)
	recreateReceiptsRequired := receiptsRequired && tracesRequired && (txsRequired || headersRequired)

	blocksInfo := &BlocksInfo{
		Blocks:   promise.NewPromise[[]*types.Block](),
		Receipts: promise.NewPromise[[]BlockReceipts](),
		Traces:   promise.NewPromise[[]BlockTraces](),
		Headers:  promise.NewPromise[[]*types.Header](),
	}

	if txsRequired {
		go func() {
			blocks, err := c.BatchGetBlocksWithTransactions(ctx, startBlock, endBlock)
			if err != nil {
				c.logger.Warn("error getting blocks", zap.Error(err))
				blocksInfo.Blocks.Reject(err)
				return
			}

			blocksInfo.Blocks.Resolve(blocks)
		}()
	} else {
		blocksInfo.Blocks.Resolve(nil)
	}

	if !txsRequired && headersRequired {
		go func() {
			headers, err := c.BatchGetBlocksHeaders(ctx, startBlock, endBlock)
			if err != nil {
				c.logger.Warn("error getting headers", zap.Error(err))
				blocksInfo.Headers.Reject(err)
				return
			}
			blocksInfo.Headers.Resolve(headers)
		}()
	} else {
		blocksInfo.Headers.Resolve(nil)
	}

	if recreateReceiptsRequired {
		go func() {
			var blocks []*types.Block
			var headers []*types.Header
			var traces []BlockTraces
			var err error

			if txsRequired {
				blocks, err = blocksInfo.Blocks.WaitContext(ctx)
				if err != nil {
					c.logger.Warn("error getting blocks for recreated receipts", zap.Error(err))
					blocksInfo.Receipts.Reject(err)
					return
				}
			} else {
				headers, err = blocksInfo.Headers.WaitContext(ctx)
				if err != nil {
					c.logger.Warn("error getting headers for recreated receipts", zap.Error(err))
					blocksInfo.Receipts.Reject(err)
					return
				}
			}

			traces, err = blocksInfo.Traces.WaitContext(ctx)
			if err != nil {
				c.logger.Warn("error getting traces for recreated receipts", zap.Error(err))
				blocksInfo.Receipts.Reject(err)
				return
			}

			recreateStarted := time.Now()
			receipts, err := recreateReceiptsFromTraces(blocks, headers, traces, c.Chain)
			c.timings.recordStage("receipts", "recreate", time.Since(recreateStarted), err == nil)
			if err != nil {
				c.logger.Warn("error recreating receipts from traces", zap.Error(err))
				blocksInfo.Receipts.Reject(err)
				return
			}
			var receiptCount uint64
			for _, blockReceipts := range receipts {
				receiptCount += uint64(len(blockReceipts))
			}
			c.timings.recordItems("receipts", "recreated_receipts", receiptCount)
			blocksInfo.Receipts.Resolve(receipts)
		}()
	} else if receiptsRequired {
		go func() {
			receipts, err := c.BatchGetBlocksReceipts(ctx, startBlock, endBlock)
			if err != nil {
				c.logger.Warn("error getting receipts", zap.Error(err))
				blocksInfo.Receipts.Reject(err)
				return
			}

			blocksInfo.Receipts.Resolve(receipts)
		}()
	} else {
		blocksInfo.Receipts.Resolve(nil)
	}

	if tracesRequired {
		go func() {
			var traces []BlockTraces
			var err error
			if recreateReceiptsRequired {
				traces, err = c.BatchGetBlocksDebugTracesWithLogs(ctx, startBlock, endBlock)
			} else {
				traces, err = c.BatchGetBlocksDebugTraces(ctx, startBlock, endBlock)
			}
			if err != nil {
				c.logger.Warn("error getting traces", zap.Error(err))
				blocksInfo.Traces.Reject(err)
				return
			}

			blocksInfo.Traces.Resolve(traces)
		}()
	} else {
		blocksInfo.Traces.Resolve(nil)
	}

	return blocksInfo
}
