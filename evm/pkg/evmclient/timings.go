package evmclient

import (
	"sort"
	"sync"
	"time"

	"go.uber.org/zap"
)

type timingKey struct {
	dataType string
	stage    string
	status   string
}

type timingStats struct {
	count uint64
	total time.Duration
	max   time.Duration
}

type itemKey struct {
	dataType string
	item     string
}

type clientTimings struct {
	mu       sync.Mutex
	stages   map[timingKey]timingStats
	rawBytes map[string]uint64
	items    map[itemKey]uint64
}

func newClientTimings() *clientTimings {
	return &clientTimings{
		stages:   make(map[timingKey]timingStats),
		rawBytes: make(map[string]uint64),
		items:    make(map[itemKey]uint64),
	}
}

func (t *clientTimings) recordStage(dataType, stage string, duration time.Duration, ok bool) {
	if t == nil {
		return
	}
	status := "ok"
	if !ok {
		status = "error"
	}
	key := timingKey{dataType: dataType, stage: stage, status: status}
	t.mu.Lock()
	stats := t.stages[key]
	stats.count++
	stats.total += duration
	if duration > stats.max {
		stats.max = duration
	}
	t.stages[key] = stats
	t.mu.Unlock()
}

func (t *clientTimings) recordRawBytes(dataType string, bytes uint64) {
	if t == nil || bytes == 0 {
		return
	}
	t.mu.Lock()
	t.rawBytes[dataType] += bytes
	t.mu.Unlock()
}

func (t *clientTimings) recordItems(dataType, item string, count uint64) {
	if t == nil || count == 0 {
		return
	}
	t.mu.Lock()
	t.items[itemKey{dataType: dataType, item: item}] += count
	t.mu.Unlock()
}

func (t *clientTimings) snapshotAndReset() (map[timingKey]timingStats, map[string]uint64, map[itemKey]uint64) {
	if t == nil {
		return nil, nil, nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	stages := t.stages
	rawBytes := t.rawBytes
	items := t.items

	t.stages = make(map[timingKey]timingStats)
	t.rawBytes = make(map[string]uint64)
	t.items = make(map[itemKey]uint64)

	return stages, rawBytes, items
}

type batchTiming struct {
	dataType string
}

func newBatchTiming(dataType string, _ string, _ string, _ int64) *batchTiming {
	return &batchTiming{
		dataType: dataType,
	}
}

func (b *batchTiming) addStage(stage string, duration time.Duration) {
}

func (b *batchTiming) addRawBytes(bytes uint64) {
}

func (b *batchTiming) addItems(item string, count uint64) {
}

func (c *Client) measureBatchStage(batch *batchTiming, stage string, fn func() error) error {
	start := time.Now()
	err := fn()
	duration := time.Since(start)
	batch.addStage(stage, duration)
	c.timings.recordStage(batch.dataType, stage, duration, err == nil)
	return err
}

func (c *Client) finishBatchTiming(batch *batchTiming, started time.Time, err error) {
	if batch == nil {
		return
	}
	c.timings.recordStage(batch.dataType, "total", time.Since(started), err == nil)
}

func (c *Client) LogTimings() {
	stages, rawBytes, items := c.timings.snapshotAndReset()
	if len(stages) == 0 && len(rawBytes) == 0 && len(items) == 0 {
		c.logger.Info("evm client timings", zap.String("status", "empty"))
		return
	}

	stageKeys := make([]timingKey, 0, len(stages))
	for key := range stages {
		stageKeys = append(stageKeys, key)
	}
	sort.Slice(stageKeys, func(i, j int) bool {
		if stageKeys[i].dataType != stageKeys[j].dataType {
			return stageKeys[i].dataType < stageKeys[j].dataType
		}
		if stageKeys[i].stage != stageKeys[j].stage {
			return stageKeys[i].stage < stageKeys[j].stage
		}
		return stageKeys[i].status < stageKeys[j].status
	})

	for _, key := range stageKeys {
		stats := stages[key]
		avg := time.Duration(0)
		if stats.count > 0 {
			avg = stats.total / time.Duration(stats.count)
		}
		c.logger.Info("evm client timing",
			zap.String("data_type", key.dataType),
			zap.String("stage", key.stage),
			zap.String("status", key.status),
			zap.Uint64("count", stats.count),
			zap.Duration("total", stats.total),
			zap.Duration("avg", avg),
			zap.Duration("max", stats.max),
		)
	}

	rawKeys := sortedStringKeys(rawBytes)
	for _, dataType := range rawKeys {
		c.logger.Info("evm client timing bytes",
			zap.String("data_type", dataType),
			zap.Uint64("raw_bytes", rawBytes[dataType]),
		)
	}

	itemKeys := make([]itemKey, 0, len(items))
	for key := range items {
		itemKeys = append(itemKeys, key)
	}
	sort.Slice(itemKeys, func(i, j int) bool {
		if itemKeys[i].dataType != itemKeys[j].dataType {
			return itemKeys[i].dataType < itemKeys[j].dataType
		}
		return itemKeys[i].item < itemKeys[j].item
	})
	for _, key := range itemKeys {
		c.logger.Info("evm client timing items",
			zap.String("data_type", key.dataType),
			zap.String("item", key.item),
			zap.Uint64("count", items[key]),
		)
	}
}

func sortedStringKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
