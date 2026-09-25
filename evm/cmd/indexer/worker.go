package main

import (
	"context"
	"errors"
	"math/big"
	"sync/atomic"
	"time"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/processors"
	"github.com/monolit-network/xlayer-indexer/util/iterators"
	"go.uber.org/zap"
)

var errShutdownDrain = errors.New("shutdown requested: drained in-flight window")

const throughputWindowSize = 10 * time.Second

type throughputSample struct {
	at     time.Time
	blocks int64
}

type throughputWindow struct {
	window  time.Duration
	samples []throughputSample
	total   int64
}

func newThroughputWindow(window time.Duration) *throughputWindow {
	return &throughputWindow{
		window:  window,
		samples: make([]throughputSample, 0, 128),
	}
}

func (w *throughputWindow) trim(now time.Time) {
	if w == nil {
		return
	}
	cutoff := now.Add(-w.window)
	drop := 0
	for drop < len(w.samples) && w.samples[drop].at.Before(cutoff) {
		w.total -= w.samples[drop].blocks
		drop++
	}
	if drop > 0 {
		w.samples = append(w.samples[:0], w.samples[drop:]...)
	}
}

func (w *throughputWindow) Record(now time.Time, blocks int64) (int64, float64) {
	if w == nil || blocks <= 0 {
		return 0, 0
	}
	w.trim(now)
	w.samples = append(w.samples, throughputSample{at: now, blocks: blocks})
	w.total += blocks
	return w.total, float64(w.total) / w.window.Seconds()
}

func (w *throughputWindow) Snapshot(now time.Time) (int64, float64) {
	if w == nil {
		return 0, 0
	}
	w.trim(now)
	return w.total, float64(w.total) / w.window.Seconds()
}

// routingClient is what we pass to ProcessorAggregator as EVMClient. It
// routes BatchGetBlocksInfo calls to the live or historical client based on
// an atomic flag toggled by the worker right before dispatch. Because the
// worker is the single writer and it waits for dispatch to complete before
// switching the flag, there is no race.
type routingClient struct {
	live          *evmclient.Client
	historical    *evmclient.Client
	useHistorical atomic.Bool
}

func (r *routingClient) BatchGetBlocksInfo(
	ctx context.Context,
	startBlock, endBlock *big.Int,
	dt evmclient.RequiredDataTypes,
) *evmclient.BlocksInfo {
	if r.useHistorical.Load() {
		return r.historical.BatchGetBlocksInfo(ctx, startBlock, endBlock, dt)
	}
	return r.live.BatchGetBlocksInfo(ctx, startBlock, endBlock, dt)
}

// workerDeps bundles what the worker needs; keeps runWorker's signature sane.
type workerDeps struct {
	logger       *zap.Logger
	agg          *processors.ProcessorAggregator
	router       *routingClient
	liveClient   *evmclient.Client
	histClient   *evmclient.Client
	state        *indexerState
	cfg          indexerConfig
	shuttingDown *atomic.Bool
	throughput   *throughputWindow
	// reorg dependencies (click/pg/notifier/chain) are captured in a closure
	// and invoked via handleReorgFn to keep this file free of DB imports.
	handleReorgFn func(ctx context.Context) bool
}

func runWorker(ctx context.Context, d workerDeps) {
	if d.throughput == nil {
		d.throughput = newThroughputWindow(throughputWindowSize)
	}

	for {
		if d.shuttingDown != nil && d.shuttingDown.Load() {
			d.logger.Info("shutdown requested, worker exiting")
			return
		}

		last, latest := d.state.Snapshot()
		if last >= latest {
			if !d.state.WaitForUpdate(ctx) {
				if d.shuttingDown != nil && d.shuttingDown.Load() {
					d.logger.Info("shutdown requested while worker idle, exiting")
				}
				return
			}
			continue
		}

		gap := latest - last
		if gap > d.cfg.LiveThreshold {
			histLatest, err := historicalLatestUntilSuccess(ctx, d.logger, d.histClient)
			if err != nil {
				if errors.Is(err, context.Canceled) && d.shuttingDown != nil && d.shuttingDown.Load() {
					d.logger.Info("shutdown requested while waiting for historical head, exiting")
				}
				return
			}

			// Catchup: parallel fetch, ordered apply, chunked so each chunk
			// commits to state before we move on (bounds memory pressure).
			chunkEnd := last + d.cfg.HistoricalBatch*int64(d.cfg.CatchupChunkBatches)
			tipGuard := latest - 32
			if tipGuard < last+1 {
				tipGuard = last + 1
			}
			if tipGuard > histLatest {
				tipGuard = histLatest
			}
			if tipGuard < last+1 {
				d.logger.Info("historical node behind next required block, waiting",
					zap.Int64("next_block", last+1),
					zap.Int64("historical_latest", histLatest),
					zap.Int64("live_latest", latest),
				)
				if sleepCtx(ctx, retryDelay) != nil {
					if d.shuttingDown != nil && d.shuttingDown.Load() {
						d.logger.Info("shutdown requested while waiting for historical node, exiting")
					}
					return
				}
				continue
			}
			if chunkEnd > tipGuard {
				chunkEnd = tipGuard
			}
			// Keep current in-flight chunk drainable on shutdown: we stop
			// scheduling new chunks at loop top, but finish this one and flush.
			processCtx := context.WithoutCancel(ctx)
			err = runCatchupChunk(processCtx, d, last+1, chunkEnd)
			if errors.Is(err, errShutdownDrain) {
				d.logger.Info("shutdown drain complete, worker exiting")
				return
			}
			if err != nil {
				d.logger.Error("catchup chunk failed, worker exiting",
					zap.Int64("start", last+1),
					zap.Int64("end", chunkEnd),
					zap.Int64("last_processed", last),
					zap.Int64("latest", latest),
					zap.Error(err),
				)
				// runCatchupChunk retries internally; reaching here means an
				// unexpected failure surfaced from the chunk runner.
				return
			}
		} else {
			processCtx := context.WithoutCancel(ctx)
			err := runLiveBlock(processCtx, d, last+1)
			if errors.Is(err, errReorg) {
				if !d.handleReorgFn(processCtx) {
					return
				}
				continue
			}
			if err != nil {
				d.logger.Error("live block failed, worker exiting",
					zap.Int64("block", last+1),
					zap.Int64("last_processed", last),
					zap.Int64("latest", latest),
					zap.Error(err),
				)
				return
			}
		}
	}
}

func historicalLatestUntilSuccess(ctx context.Context, logger *zap.Logger, client *evmclient.Client) (int64, error) {
	for ctx.Err() == nil {
		n, err := client.GetEthClient().BlockNumber(ctx)
		if err == nil {
			return int64(n), nil
		}
		if errors.Is(err, context.Canceled) {
			return 0, err
		}
		logger.Warn("historical latest block fetch failed, retrying", zap.Error(err))
		if sleepCtx(ctx, retryDelay) != nil {
			return 0, ctx.Err()
		}
	}
	return 0, ctx.Err()
}

// runLiveBlock processes a single block through the live client with a
// reorg check against the ring buffer. Retries until success or ctx done.
func runLiveBlock(ctx context.Context, d workerDeps, blockNumber int64) error {
	for ctx.Err() == nil {
		if err := checkReorgForNext(ctx, d.logger, d.liveClient.GetEthClient(), d.state.Ring(), blockNumber-1); err != nil {
			if errors.Is(err, errReorg) {
				return errReorg
			}
			return err
		}

		header, err := headerByNumberUntilSuccess(ctx, d.logger, d.liveClient.GetEthClient(), blockNumber, true)
		if err != nil {
			return err
		}
		if header == nil {
			// Block not yet on chain; back off and let poller/subscription
			// push a fresh `latest` before we try again.
			if sleepCtx(ctx, retryDelay) != nil {
				return ctx.Err()
			}
			continue
		}

		d.router.useHistorical.Store(false)
		began := time.Now()
		resCh := d.agg.ProcessBlockRangeWithResult(ctx, big.NewInt(blockNumber), big.NewInt(blockNumber))
		err = <-resCh
		if err == nil {
			d.agg.FlushSync()
			d.state.Ring().Append(blockRef{Number: blockNumber, Hash: header.Hash()})
			d.state.SetLastProcessed(blockNumber)
			d.logger.Info("live block done",
				zap.Int64("block", blockNumber),
				zap.String("hash", header.Hash().Hex()),
				zap.Duration("duration", time.Since(began)),
			)
			return nil
		}
		if errors.Is(err, context.Canceled) {
			return ctx.Err()
		}
		d.logger.Error("live block failed, retrying",
			zap.Int64("block", blockNumber),
			zap.Duration("duration", time.Since(began)),
			zap.Error(err),
		)
		if sleepCtx(ctx, retryDelay) != nil {
			return ctx.Err()
		}
	}
	return ctx.Err()
}

// runCatchupChunk processes [start, end] in parallel fetch / ordered apply
// batches of size cfg.HistoricalBatch with up to cfg.MaxInFlight concurrent
// BlocksInfo fetches. Returns errReorg if the boundary check fails BEFORE we
// touch aggregator state (so no partial apply); non-reorg errors trigger
// an inline retry until ctx cancellation.
func runCatchupChunk(ctx context.Context, d workerDeps, start, end int64) error {
	d.router.useHistorical.Store(true)

	// Boundary reorg check: fetch parent of `start` and compare with ring.Last().
	if !d.state.Ring().IsEmpty() {
		if err := checkReorgForNext(ctx, d.logger, d.liveClient.GetEthClient(), d.state.Ring(), start-1); err != nil {
			return err
		}
	}

	subranges := iterators.BigRangeSubranges(big.NewInt(start), big.NewInt(end), int(d.cfg.HistoricalBatch))
	ranges := make([]iterators.BigSubrange, 0)
	for _, r := range subranges {
		ranges = append(ranges, r)
	}
	if len(ranges) == 0 {
		return nil
	}

	processedEnd, err := runParallelBatchesUntilSuccess(ctx, d, ranges)
	if err != nil && !errors.Is(err, errShutdownDrain) {
		return err
	}

	if processedEnd >= start {
		// The ring only keeps the most recent ringBufferSize hashes, so fetch
		// just the processed tail instead of the whole chunk.
		tailStart := processedEnd - int64(ringBufferSize) + 1
		if tailStart < start {
			tailStart = start
		}
		refs, ferr := fetchCanonicalHashes(ctx, d.logger, d.liveClient, tailStart, processedEnd)
		if ferr != nil {
			return ferr
		}
		d.state.Ring().Append(refs...)
		d.agg.FlushSync()
		d.state.SetLastProcessed(processedEnd)
		blocksLast10s, blocksPerSec10s := d.throughput.Snapshot(time.Now())
		if errors.Is(err, errShutdownDrain) {
			d.logger.Info("catchup chunk drained in-flight window",
				zap.Int64("start", start),
				zap.Int64("processed_end", processedEnd),
				zap.Int64("requested_end", end),
				zap.Int64("blocks_last_10s", blocksLast10s),
				zap.Float64("blocks_per_sec_10s", blocksPerSec10s),
			)
		} else {
			d.logger.Info("catchup chunk done",
				zap.Int64("start", start),
				zap.Int64("end", processedEnd),
				zap.Int("num_batches", len(ranges)),
				zap.Int64("blocks_last_10s", blocksLast10s),
				zap.Float64("blocks_per_sec_10s", blocksPerSec10s),
			)
		}
	}
	return err
}

// runParallelBatchesUntilSuccess dispatches up to cfg.MaxInFlight
// BlocksInfo fetches in parallel, waits for them in enqueue order (so the
// processor applies tx strictly in block order), and retries any failed
// batch inline until success. Because ordered waiting matches the single
// aggregator worker's FIFO dispatch, polymarket ordering is preserved even
// across retries.
func runParallelBatchesUntilSuccess(
	ctx context.Context,
	d workerDeps,
	ranges []iterators.BigSubrange,
) (int64, error) {
	type pending struct {
		idx   int
		batch iterators.BigSubrange
		resCh <-chan error
	}

	maxInFlight := d.cfg.MaxInFlight
	if maxInFlight < 1 {
		maxInFlight = 1
	}

	nextIdx := 0
	queue := make([]pending, 0, maxInFlight)
	stopDispatch := false
	lastCompletedEnd := ranges[0].Start.Int64() - 1

	dispatch := func() {
		if stopDispatch {
			return
		}
		for nextIdx < len(ranges) && len(queue) < maxInFlight {
			if d.shuttingDown != nil && d.shuttingDown.Load() {
				stopDispatch = true
				return
			}
			r := ranges[nextIdx]
			d.logger.Info("catchup batch queued",
				zap.Int("batch_idx", nextIdx),
				zap.Int64("start", r.Start.Int64()),
				zap.Int64("end", r.End.Int64()),
				zap.Int("in_flight", len(queue)+1),
				zap.Int("max_in_flight", maxInFlight),
			)
			resCh := d.agg.ProcessBlockRangeWithResult(ctx, r.Start, r.End)
			queue = append(queue, pending{idx: nextIdx, batch: r, resCh: resCh})
			nextIdx++
		}
	}

	dispatch()

	for len(queue) > 0 {
		if ctx.Err() != nil {
			// Drain in-flight to avoid leaking goroutines inside aggregator.
			for _, p := range queue {
				<-p.resCh
			}
			return lastCompletedEnd, ctx.Err()
		}

		head := queue[0]
		began := time.Now()
		d.logger.Info("catchup batch waiting",
			zap.Int("batch_idx", head.idx),
			zap.Int64("start", head.batch.Start.Int64()),
			zap.Int64("end", head.batch.End.Int64()),
			zap.Int("queued", len(queue)),
		)
		err := <-head.resCh
		if err == nil {
			lastCompletedEnd = head.batch.End.Int64()
			now := time.Now()
			blocksLast10s, blocksPerSec10s := d.throughput.Record(now, head.batch.End.Int64()-head.batch.Start.Int64()+1)
			d.logger.Info("catchup batch done",
				zap.Int64("start", head.batch.Start.Int64()),
				zap.Int64("end", head.batch.End.Int64()),
				zap.Duration("duration", now.Sub(began)),
				zap.Int64("blocks_last_10s", blocksLast10s),
				zap.Float64("blocks_per_sec_10s", blocksPerSec10s),
			)
			queue = queue[1:]
			dispatch()
			continue
		}
		if errors.Is(err, context.Canceled) {
			for _, p := range queue[1:] {
				<-p.resCh
			}
			if ctx.Err() != nil {
				return lastCompletedEnd, ctx.Err()
			}
			return lastCompletedEnd, err
		}

		// Batch failed. Drain the rest of the window (already fetched, values
		// may have been queued to processors ahead of us — but with sync
		// processing the aggregator worker is still waiting for THIS batch's
		// tx loop). With parallel dispatch + sync=true, the aggregator worker
		// processes requests in enqueue order: it won't start the next until
		// this one returns. So returning err here means subsequent batches
		// were not applied yet — but we already dispatched them, so we must
		// let them drain or they'll leak.
		d.logger.Error("catchup batch failed, will retry",
			zap.Int64("start", head.batch.Start.Int64()),
			zap.Int64("end", head.batch.End.Int64()),
			zap.Duration("duration", time.Since(began)),
			zap.Error(err),
		)

		// Wait for the rest of the queue (they may succeed or fail — either
		// way we have to drain), then retry from the failed batch.
		for _, p := range queue[1:] {
			<-p.resCh
		}
		queue = queue[:0]

		// Rewind to the failed batch index.
		nextIdx = head.idx

		if sleepCtx(ctx, retryDelay) != nil {
			return lastCompletedEnd, ctx.Err()
		}
		dispatch()
	}
	if stopDispatch && nextIdx < len(ranges) {
		return lastCompletedEnd, errShutdownDrain
	}
	return lastCompletedEnd, nil
}

func fetchCanonicalHashes(
	ctx context.Context,
	logger *zap.Logger,
	client *evmclient.Client,
	start, end int64,
) ([]blockRef, error) {
	for ctx.Err() == nil {
		headers, err := client.BatchGetBlocksHeaders(ctx, big.NewInt(start), big.NewInt(end))
		if err == nil {
			return headersToRefs(headers), nil
		}
		if errors.Is(err, context.Canceled) {
			return nil, err
		}
		logger.Warn("fetch canonical hashes failed, retrying",
			zap.Int64("start", start),
			zap.Int64("end", end),
			zap.Error(err),
		)
		if sleepCtx(ctx, retryDelay) != nil {
			return nil, ctx.Err()
		}
	}
	return nil, ctx.Err()
}
