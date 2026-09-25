package queries

const (
	CreateSchemaCexClickhouseSQL = `CREATE DATABASE IF NOT EXISTS cex`

	// bybit kline
	CreateTableBybitKlineSQL = `
		CREATE TABLE IF NOT EXISTS cex.bybit_kline (
			ticker_id UInt64,
			open_time DateTime,
			open_price Float64,
			high_price Float64,
			low_price Float64,
			close_price Float64,
			volume Float64,
			turnover Float64
		)
		ENGINE = MergeTree
		PARTITION BY toYYYYMM(open_time)
		ORDER BY (ticker_id, open_time)`

	SelectBybitKlineLatestOpenTimeSQL = `
		SELECT open_time
		FROM cex.bybit_kline
		WHERE ticker_id = ?
		ORDER BY open_time DESC
		LIMIT 1;`

	InsertBybitKlineSQL = `
		INSERT INTO cex.bybit_kline (
			ticker_id, open_time, open_price, high_price, low_price, close_price, volume, turnover
		)
		VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		)`

	InsertBybitKlinesSQL = `
		INSERT INTO cex.bybit_kline (
			ticker_id, open_time, open_price, high_price, low_price, close_price, volume, turnover
		)`

	SelectBybitKlineDataRangeSQL = `
		SELECT ticker_id, open_time, open_price, high_price, low_price, close_price, volume, turnover
		FROM cex.bybit_kline
		WHERE ticker_id = ? AND open_time >= ? AND open_time <= ?
		ORDER BY open_time ASC`

	SelectBybitKlineDataAggregatedSQL = `
		WITH
			toInt32(?) AS interval_minutes,
			interval_minutes / 5 AS expected_rows
		SELECT
			ticker_id,
			bucket AS open_time,
			argMin(open_price, ts) AS open_price,
			max(high_price) AS high_price,
			min(low_price) AS low_price,
			argMax(close_price, ts) AS close_price,
			sum(volume) AS volume,
			sum(turnover) AS turnover
		FROM (
			SELECT
				ticker_id,
				open_time AS ts,
				toStartOfInterval(open_time, toIntervalMinute(interval_minutes)) AS bucket,
				open_price,
				high_price,
				low_price,
				close_price,
				volume,
				turnover
			FROM cex.bybit_kline
			WHERE ticker_id = ?
			  AND open_time >= ?
			  AND open_time < ?
		)
		GROUP BY ticker_id, bucket
		HAVING count() >= expected_rows
		ORDER BY ticker_id, bucket`

	SelectBybitKlineLatestCloseSQL = `
		SELECT close_price
		FROM cex.bybit_kline
		WHERE ticker_id = ?
		ORDER BY open_time DESC
		LIMIT 1;`

	// bybit funding rate
	CreateTableBybitFundingRateSQL = `
		CREATE TABLE IF NOT EXISTS cex.bybit_funding_rate (
			ticker_id UInt64,
			funding_rate Float64,
			funding_rate_timestamp DateTime
		)
		ENGINE = ReplacingMergeTree(funding_rate_timestamp)
		PARTITION BY toYYYYMM(funding_rate_timestamp)
		ORDER BY (ticker_id, funding_rate_timestamp)`

	SelectBybitFundingRateLatestFundingRateTimestampSQL = `
		SELECT funding_rate_timestamp
		FROM cex.bybit_funding_rate
		WHERE ticker_id = ?
		ORDER BY funding_rate_timestamp DESC
		LIMIT 1;`

	InsertBybitFundingRateSQL = `
		INSERT INTO cex.bybit_funding_rate (
			ticker_id, funding_rate, funding_rate_timestamp
		)
		VALUES (
			?, ?, ?
		)`

	InsertBybitFundingRatesSQL = `
		INSERT INTO cex.bybit_funding_rate (
			ticker_id, funding_rate, funding_rate_timestamp
		)`

	SelectBybitFundingRateDataRangeSQL = `
		SELECT ticker_id, funding_rate, funding_rate_timestamp
		FROM cex.bybit_funding_rate
		WHERE ticker_id = ? AND funding_rate_timestamp >= ? AND funding_rate_timestamp <= ?
		ORDER BY funding_rate_timestamp ASC`

	// bybit ticker info (ClickHouse copy)
	CreateTableBybitTickerInfoClickhouseSQL = `
		CREATE TABLE IF NOT EXISTS cex.bybit_ticker_info (
			id UInt64,
			base_coin String,
			symbol String,
			category String,
			base_coin_parsed Nullable(String),
			multiplier UInt64,
			quote_coin String,
			status String
		)
		ENGINE = ReplacingMergeTree()
		ORDER BY (id, symbol, category)`

	InsertBybitTickerInfoClickhouseSQL = `
		INSERT INTO cex.bybit_ticker_info (
			id, base_coin, symbol, category, base_coin_parsed, multiplier, quote_coin, status
		)
		VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		)`

	InsertBybitTickerInfosClickhouseSQL = `
		INSERT INTO cex.bybit_ticker_info (
			id, base_coin, symbol, category, base_coin_parsed, multiplier, quote_coin, status
		)`

	SelectAllBybitTickerInfoClickhouseSQL = `
		SELECT id, base_coin, symbol, category, base_coin_parsed, multiplier, quote_coin, status
		FROM cex.bybit_ticker_info
		ORDER BY id`

	// bybit liquidations
	CreateTableBybitLiquidationSQL = `
		CREATE TABLE IF NOT EXISTS cex.bybit_liquidation (
			ticker_id UInt64,
			size Float64,
			liquidation_price Float64,
			liquidation_timestamp DateTime
		)
		ENGINE = MergeTree
		PARTITION BY toYYYYMM(liquidation_timestamp)
		ORDER BY (ticker_id, liquidation_timestamp)`

	SelectBybitLiquidationLatestTimestampSQL = `
		SELECT liquidation_timestamp
		FROM cex.bybit_liquidation
		WHERE ticker_id = ?
		ORDER BY liquidation_timestamp DESC
		LIMIT 1;`

	InsertBybitLiquidationSQL = `
		INSERT INTO cex.bybit_liquidation (
			ticker_id, size, liquidation_price, liquidation_timestamp
		)
		VALUES (
			?, ?, ?, ?
		)`

	InsertBybitLiquidationsSQL = `
		INSERT INTO cex.bybit_liquidation (
			ticker_id, size, liquidation_price, liquidation_timestamp
		)`

	SelectBybitLiquidationDataRangeSQL = `
		SELECT ticker_id, size, liquidation_price, liquidation_timestamp
		FROM cex.bybit_liquidation
		WHERE ticker_id = ? AND liquidation_timestamp >= ? AND liquidation_timestamp <= ?
		ORDER BY liquidation_timestamp ASC`

	// bybit trades
	CreateTableBybitTradeSQL = `
		CREATE TABLE IF NOT EXISTS cex.bybit_trade (
			ticker_id UInt64,
			size Float64,
			price Float64,
			trade_id String,
			trade_timestamp DateTime
		)
		ENGINE = MergeTree
		PARTITION BY toYYYYMM(trade_timestamp)
		ORDER BY (ticker_id, trade_timestamp)`

	SelectBybitTradeLatestTimestampSQL = `
		SELECT trade_timestamp
		FROM cex.bybit_trade
		WHERE ticker_id = ?
		ORDER BY trade_timestamp DESC
		LIMIT 1;`

	InsertBybitTradeSQL = `
		INSERT INTO cex.bybit_trade (
			ticker_id, size, price, trade_id, trade_timestamp
		)
		VALUES (
			?, ?, ?, ?, ?
		)`

	InsertBybitTradesSQL = `
		INSERT INTO cex.bybit_trade (
			ticker_id, size, price, trade_id, trade_timestamp
		)`

	SelectBybitTradeDataRangeSQL = `
		SELECT ticker_id, size, price, trade_id, trade_timestamp
		FROM cex.bybit_trade
		WHERE ticker_id = ? AND trade_timestamp >= ? AND trade_timestamp <= ?
		ORDER BY trade_timestamp ASC`

	InsertListingSQL = `
		INSERT INTO cex.listings (
			cex_name, ticker, dt
		)
		VALUES (
			?, ?, ?
		)`
)
