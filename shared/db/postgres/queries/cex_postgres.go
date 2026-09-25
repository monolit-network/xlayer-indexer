package queries

const (
	CreateSchemaCexSQL = `CREATE SCHEMA IF NOT EXISTS cex`

	CreateTableBybitTickerInfoSQL = `
		CREATE TABLE IF NOT EXISTS cex.bybit_ticker_info (
			id              SERIAL PRIMARY KEY,
			base_coin       TEXT NOT NULL,
			symbol          TEXT NOT NULL,
			category        TEXT NOT NULL,
			base_coin_parsed TEXT,
			multiplier      BIGINT NOT NULL,
			quote_coin      TEXT NOT NULL,
			status          TEXT NOT NULL,
			UNIQUE (symbol, category)
		);`

	UpsertBybitTickerInfoSQL = `
		INSERT INTO cex.bybit_ticker_info (
			base_coin, symbol, category, base_coin_parsed, multiplier, quote_coin, status
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		ON CONFLICT (symbol, category)
		DO UPDATE SET status = EXCLUDED.status
		RETURNING id;`

	SelectBybitTickerInfoBySymbolAndCategorySQL = `
		SELECT id, base_coin, symbol, category, base_coin_parsed, multiplier, quote_coin, status
		FROM cex.bybit_ticker_info
		WHERE symbol = $1 AND category = $2;`

	SelectBybitTickerInfoByIDSQL = `
		SELECT id, base_coin, symbol, category, base_coin_parsed, multiplier, quote_coin, status
		FROM cex.bybit_ticker_info
		WHERE id = $1;`

	SelectAllBybitTickerInfoPostgresSQL = `
		SELECT id, base_coin, symbol, category, base_coin_parsed, multiplier, quote_coin, status
		FROM cex.bybit_ticker_info
		ORDER BY id;`
)
