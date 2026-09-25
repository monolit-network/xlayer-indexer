package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/shared/db"
	sharedmodels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/monolit-network/xlayer-indexer/util/notifier"
	"github.com/ethereum/go-ethereum"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// errReorg signals the worker that a reorg was detected and the reorg handler
// must run before processing resumes.
var errReorg = errors.New("reorg detected")

// checkReorgForNext returns errReorg if the header of (lastProcessed+1) on
// chain has a ParentHash that does not match our ring buffer's last entry.
// Returns nil if chain hasn't produced the next block yet.
func checkReorgForNext(
	ctx context.Context,
	logger *zap.Logger,
	client *ethclient.Client,
	ring *hashRing,
	lastProcessed int64,
) error {
	last, ok := ring.Last()
	if !ok {
		return nil
	}
	if last.Number != lastProcessed {
		logger.Warn("ring buffer out of sync",
			zap.Int64("ring_last", last.Number),
			zap.Int64("last_processed", lastProcessed),
		)
		return errReorg
	}

	header, err := headerByNumberUntilSuccess(ctx, logger, client, lastProcessed+1, true /*allowNotFound*/)
	if err != nil {
		return err
	}
	if header == nil {
		return nil
	}
	if header.ParentHash != last.Hash {
		logger.Warn("reorg: parent hash mismatch",
			zap.Int64("next_block", lastProcessed+1),
			zap.String("expected_parent", last.Hash.Hex()),
			zap.String("got_parent", header.ParentHash.Hex()),
		)
		return errReorg
	}
	return nil
}

// handleReorgUntilSuccess walks the ring from oldest to newest, finds the first
// block whose hash no longer matches the chain, rolls the DB back from that
// block, publishes a reorg event on NATS, rebuilds the ring buffer from chain
// headers, and rewinds state.lastProcessed. Retries the whole thing on error.
func handleReorgUntilSuccess(
	ctx context.Context,
	logger *zap.Logger,
	clickClient db.DBClick,
	dbClient db.Client,
	eventNotifier *notifier.Notifier,
	liveClient *evmclient.Client,
	state *indexerState,
	chain models.Chain,
) bool {
	for ctx.Err() == nil {
		if err := handleReorg(ctx, logger, clickClient, dbClient, eventNotifier, liveClient, state, chain); err == nil {
			return true
		} else if errors.Is(err, context.Canceled) {
			return false
		} else {
			logger.Error("reorg handling failed, retrying", zap.Error(err))
			if sleepCtx(ctx, retryDelay) != nil {
				return false
			}
		}
	}
	return false
}

func handleReorg(
	ctx context.Context,
	logger *zap.Logger,
	clickClient db.DBClick,
	dbClient db.Client,
	eventNotifier *notifier.Notifier,
	liveClient *evmclient.Client,
	state *indexerState,
	chain models.Chain,
) error {
	ring := state.Ring()
	entries := ring.Entries()
	if len(entries) == 0 {
		return nil
	}

	forkPoint, err := findForkPoint(ctx, logger, liveClient.GetEthClient(), entries)
	if err != nil {
		return err
	}
	lastGood := forkPoint - 1
	rollbackDepth := entries[len(entries)-1].Number - lastGood

	logger.Warn("rolling back for reorg",
		zap.Int64("fork_point", forkPoint),
		zap.Int64("last_good", lastGood),
		zap.Int64("rollback_depth", rollbackDepth),
		zap.Int("ring_len", len(entries)),
	)

	if eventNotifier != nil {
		if err := notifier.Publish(eventNotifier, ctx, string(sharedmodels.NotifierEventReorg), models.ReorgEvent{
			BlockNumber: big.NewInt(forkPoint),
			Chain:       string(chain),
		}); err != nil {
			logger.Warn("failed to publish reorg event", zap.Error(err))
		}
	}

	if err := rollbackStores(ctx, clickClient, dbClient, chain, forkPoint, logger); err != nil {
		return fmt.Errorf("rollback stores: %w", err)
	}

	// Rebuild ring buffer from chain for [lastGood-limit+1, lastGood].
	newRing, err := fetchRingBufferUntilSuccess(ctx, logger, liveClient, lastGood, ringBufferSize)
	if err != nil {
		return fmt.Errorf("rebuild ring buffer: %w", err)
	}
	ring.Reset(newRing)

	state.ForceSetLastProcessed(lastGood)
	logger.Info("reorg handled", zap.Int64("resume_from", lastGood+1))
	return nil
}

// findForkPoint walks the ring from oldest to newest, returns the Number of
// the first entry whose hash differs from the chain. If all entries still
// match (shouldn't happen if caller observed a mismatch) we rollback from the
// newest+1 — most conservative.
func findForkPoint(
	ctx context.Context,
	logger *zap.Logger,
	client *ethclient.Client,
	entries []blockRef,
) (int64, error) {
	for _, e := range entries {
		header, err := headerByNumberUntilSuccess(ctx, logger, client, e.Number, false)
		if err != nil {
			return 0, err
		}
		if header.Hash() != e.Hash {
			return e.Number, nil
		}
	}
	return entries[len(entries)-1].Number + 1, nil
}

func rollbackStores(
	ctx context.Context,
	clickClient db.DBClick,
	dbClient db.Client,
	chain models.Chain,
	forkPoint int64,
	logger *zap.Logger,
) error {
	if chain == models.ChainPolygon {
		if err := dbClient.Transaction(ctx, func(ctx context.Context, tx db.DB) error {
			logger.Info("deleting polymarket markets", zap.Int64("fork_point", forkPoint))
			if _, err := tx.DeletePolymarketMarketsFromBlockNumber(ctx, forkPoint); err != nil {
				return fmt.Errorf("delete polymarket markets: %w", err)
			}
			logger.Info("deleted polymarket markets", zap.Int64("fork_point", forkPoint))
			if _, err := tx.DeletePolymarketMarketEventsFromBlockNumber(ctx, forkPoint); err != nil {
				return fmt.Errorf("delete polymarket market events: %w", err)
			}
			logger.Info("deleted polymarket market events", zap.Int64("fork_point", forkPoint))
			if _, err := tx.ResetPolymarketTokensResolutionFromBlockNumber(ctx, forkPoint); err != nil {
				return fmt.Errorf("reset polymarket tokens: %w", err)
			}
			logger.Info("reset polymarket tokens", zap.Int64("fork_point", forkPoint))
			if _, err := tx.ResetPolymarketMarketsResolutionFromBlockNumber(ctx, forkPoint); err != nil {
				return fmt.Errorf("reset polymarket markets: %w", err)
			}
			logger.Info("reset polymarket markets", zap.Int64("fork_point", forkPoint))
			return nil
		}); err != nil {
			return err
		}
	}
	logger.Info("deleted evm events", zap.Int64("fork_point", forkPoint))
	if err := clickClient.DeleteEvmEventsFromBlockNumber(ctx, chain, uint64(forkPoint)); err != nil {
		return fmt.Errorf("delete evm events from click: %w", err)
	}
	return nil
}

// fetchRingBufferUntilSuccess builds a ring buffer by pulling the last `limit`
// canonical headers from chain ending at lastProcessed (inclusive).
func fetchRingBufferUntilSuccess(
	ctx context.Context,
	logger *zap.Logger,
	client *evmclient.Client,
	lastProcessed int64,
	limit int,
) ([]blockRef, error) {
	if lastProcessed <= 0 {
		return nil, nil
	}
	start := lastProcessed - int64(limit) + 1
	if start < 1 {
		start = 1
	}
	for ctx.Err() == nil {
		headers, err := client.BatchGetBlocksHeaders(ctx, big.NewInt(start), big.NewInt(lastProcessed))
		if err == nil {
			return headersToRefs(headers), nil
		}
		if errors.Is(err, context.Canceled) {
			return nil, err
		}
		logger.Warn("fetch ring buffer failed, retrying",
			zap.Int64("start", start),
			zap.Int64("end", lastProcessed),
			zap.Error(err),
		)
		if sleepCtx(ctx, retryDelay) != nil {
			return nil, ctx.Err()
		}
	}
	return nil, ctx.Err()
}

func headersToRefs(headers []*ethtypes.Header) []blockRef {
	refs := make([]blockRef, 0, len(headers))
	for _, h := range headers {
		if h == nil || h.Number == nil {
			continue
		}
		refs = append(refs, blockRef{Number: h.Number.Int64(), Hash: h.Hash()})
	}
	return refs
}

// headerByNumberUntilSuccess retries ethClient.HeaderByNumber until success
// or ctx cancellation. If allowNotFound is true, a "not found" (block not
// yet mined) returns (nil, nil); otherwise it's retried like any other error.
func headerByNumberUntilSuccess(
	ctx context.Context,
	logger *zap.Logger,
	client *ethclient.Client,
	blockNumber int64,
	allowNotFound bool,
) (*ethtypes.Header, error) {
	bn := big.NewInt(blockNumber)
	for ctx.Err() == nil {
		header, err := client.HeaderByNumber(ctx, bn)
		if err == nil {
			return header, nil
		}
		if errors.Is(err, context.Canceled) {
			return nil, err
		}
		// ethereum.NotFound — allow caller to treat as "no block yet".
		if allowNotFound && isNotFoundError(err) {
			return nil, nil
		}
		logger.Warn("header by number failed, retrying",
			zap.Int64("block", blockNumber),
			zap.Error(err),
		)
		if sleepCtx(ctx, retryDelay) != nil {
			return nil, ctx.Err()
		}
	}
	return nil, ctx.Err()
}

func isNotFoundError(err error) bool {
	return err != nil && errors.Is(err, ethereum.NotFound)
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
