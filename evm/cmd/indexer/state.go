package main

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// blockRef is a (number, hash) pair kept in the hashRing so we can detect
// reorgs by comparing our processed hash with what the chain reports now.
type blockRef struct {
	Number int64
	Hash   common.Hash
}

// hashRing holds the most recent N processed (number, hash) pairs sorted
// ascending by Number. It's the sole source of truth for reorg detection.
type hashRing struct {
	mu      sync.Mutex
	entries []blockRef
	limit   int
}

func newHashRing(limit int) *hashRing {
	return &hashRing{entries: make([]blockRef, 0, limit), limit: limit}
}

func (r *hashRing) Reset(entries []blockRef) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(entries) > r.limit {
		entries = entries[len(entries)-r.limit:]
	}
	r.entries = append(r.entries[:0], entries...)
}

func (r *hashRing) Append(refs ...blockRef) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, refs...)
	if len(r.entries) > r.limit {
		r.entries = r.entries[len(r.entries)-r.limit:]
	}
}

// TruncateAt drops every entry with Number >= block.
func (r *hashRing) TruncateAt(block int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cutoff := len(r.entries)
	for i, e := range r.entries {
		if e.Number >= block {
			cutoff = i
			break
		}
	}
	r.entries = r.entries[:cutoff]
}

func (r *hashRing) Last() (blockRef, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.entries) == 0 {
		return blockRef{}, false
	}
	return r.entries[len(r.entries)-1], true
}

func (r *hashRing) Entries() []blockRef {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]blockRef(nil), r.entries...)
}

func (r *hashRing) IsEmpty() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries) == 0
}

// indexerState is the shared counter used by all goroutines. The worker reads
// both fields via Snapshot; the subscription / poller advance `latest` and
// wake the worker via notify.
type indexerState struct {
	mu            sync.Mutex
	lastProcessed int64
	latest        int64
	notify        chan struct{}
	ring          *hashRing
}

func newIndexerState(lastProcessed, latest int64, ring *hashRing) *indexerState {
	return &indexerState{
		lastProcessed: lastProcessed,
		latest:        latest,
		notify:        make(chan struct{}, 1),
		ring:          ring,
	}
}

func (s *indexerState) UpdateLatest(n int64) {
	s.mu.Lock()
	updated := n > s.latest
	if updated {
		s.latest = n
	}
	s.mu.Unlock()
	if updated {
		s.wake()
	}
}

func (s *indexerState) SetLastProcessed(n int64) {
	s.mu.Lock()
	if n > s.lastProcessed {
		s.lastProcessed = n
	}
	s.mu.Unlock()
}

func (s *indexerState) ForceSetLastProcessed(n int64) {
	s.mu.Lock()
	s.lastProcessed = n
	s.mu.Unlock()
}

func (s *indexerState) Snapshot() (int64, int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastProcessed, s.latest
}

func (s *indexerState) Ring() *hashRing { return s.ring }

func (s *indexerState) wake() {
	select {
	case s.notify <- struct{}{}:
	default:
	}
}

func (s *indexerState) WaitForUpdate(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	case <-s.notify:
		return true
	}
}
