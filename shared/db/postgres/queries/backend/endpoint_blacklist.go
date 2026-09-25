package queries

const (
	CreateTableAPIEndpointBlacklistSQL = `CREATE TABLE IF NOT EXISTS backend.api_endpoint_blacklist (
		id BIGSERIAL PRIMARY KEY,
		enabled BOOLEAN NOT NULL DEFAULT TRUE,
		method TEXT,
		match_type TEXT NOT NULL CHECK (match_type IN ('exact', 'prefix', 'regexp')),
		pattern TEXT NOT NULL,
		reason TEXT,
		disabled_until TIMESTAMPTZ,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_backend_api_endpoint_blacklist_active
		ON backend.api_endpoint_blacklist (enabled, disabled_until);
	CREATE INDEX IF NOT EXISTS idx_backend_api_endpoint_blacklist_method
		ON backend.api_endpoint_blacklist (method);`

	SelectActiveAPIEndpointBlacklistRulesSQL = `SELECT id, method, match_type, pattern, reason, disabled_until
	FROM backend.api_endpoint_blacklist
	WHERE enabled = TRUE
		AND (disabled_until IS NULL OR disabled_until > NOW())
	ORDER BY id;`
)
