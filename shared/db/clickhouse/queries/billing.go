package queries

const (
	CreateSchemaBackendClickhouseSQL = `CREATE DATABASE IF NOT EXISTS backend ON CLUSTER main_cluster`

	CreateBillingLogsTableSQL = `
		CREATE TABLE IF NOT EXISTS backend.billing_logs ON CLUSTER main_cluster
		(
			created_at DateTime64(3, 'UTC') DEFAULT now64(3),
			user_id String,
			api_key String,
			credits_amount Int64,
			operation_type LowCardinality(String),
			source LowCardinality(String),
			additional_data JSON(max_dynamic_paths = 64)
		)
		ENGINE = MergeTree
		PARTITION BY toYYYYMM(created_at)
		ORDER BY (user_id, created_at)
	`

	AddBillingLogsAdditionalDataColumnSQL = `
		ALTER TABLE backend.billing_logs ON CLUSTER main_cluster
		ADD COLUMN IF NOT EXISTS additional_data JSON(max_dynamic_paths = 64)
	`

	InsertBillingLogsSQL = `
		INSERT INTO backend.billing_logs ON CLUSTER main_cluster
			(user_id, api_key, credits_amount, operation_type, source, additional_data)
	`
)
