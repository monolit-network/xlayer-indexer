package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	registerOnce sync.Once

	InsertedRows = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "evm_indexer_inserted_rows_total",
			Help: "Total number of rows successfully inserted by the EVM indexer.",
		},
		[]string{"chain", "processor", "database", "table"},
	)

	InsertErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "evm_indexer_insert_errors_total",
			Help: "Total number of failed insert attempts by the EVM indexer.",
		},
		[]string{"chain", "processor", "database", "table"},
	)

	LastProcessedBlock = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "evm_indexer_last_processed_block",
			Help: "Last block successfully processed by the EVM indexer.",
		},
		[]string{"chain"},
	)

	LatestBlock = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "evm_indexer_latest_block",
			Help: "Latest chain head known by the EVM indexer.",
		},
		[]string{"chain"},
	)

	BlockLag = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "evm_indexer_block_lag",
			Help: "Difference between latest known block and last processed block.",
		},
		[]string{"chain"},
	)
)

func Register() {
	registerOnce.Do(func() {
		prometheus.MustRegister(
			InsertedRows,
			InsertErrors,
			LastProcessedBlock,
			LatestBlock,
			BlockLag,
		)
	})
}

func ObserveProgress(chain string, lastProcessed, latest int64) {
	LastProcessedBlock.WithLabelValues(chain).Set(float64(lastProcessed))
	LatestBlock.WithLabelValues(chain).Set(float64(latest))
	lag := latest - lastProcessed
	if lag < 0 {
		lag = 0
	}
	BlockLag.WithLabelValues(chain).Set(float64(lag))
}
