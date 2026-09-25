package clickhouse

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"go.uber.org/zap"

	"github.com/monolit-network/xlayer-indexer/shared/db"
	"github.com/monolit-network/xlayer-indexer/shared/db/clickhouse/queries"
	"github.com/monolit-network/xlayer-indexer/shared/models"
)

type ClickhouseClient struct {
	Conn   clickhouse.Conn
	logger *zap.Logger
}

func NewClickhouseClient(host string, port int, username, password string, useCompression bool, logger *zap.Logger) (*ClickhouseClient, error) {
	options := &clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", host, port)},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: username,
			Password: password,
		},
		DialTimeout:  5 * time.Second,
		MaxOpenConns: 1000,
		MaxIdleConns: 100,
		Settings: clickhouse.Settings{
			"output_format_native_write_json_as_string": true,
		},
	}

	if useCompression {
		options.Compression = &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		}
	}
	if readTimeout := os.Getenv("CLICKHOUSE_READ_TIMEOUT"); readTimeout != "" {
		timeout, err := time.ParseDuration(readTimeout)
		if err != nil {
			return nil, fmt.Errorf("invalid CLICKHOUSE_READ_TIMEOUT: %w", err)
		}
		options.ReadTimeout = timeout
	}

	conn, err := clickhouse.Open(options)
	if err != nil {
		return nil, err
	}

	return &ClickhouseClient{Conn: conn, logger: logger.Named("clickhouse")}, nil
}

func NewClickhouseClientFromEnv(logger *zap.Logger) (*ClickhouseClient, error) {
	host := os.Getenv("CLICKHOUSE_HOST")
	port := os.Getenv("CLICKHOUSE_PORT")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}
	username := os.Getenv("CLICKHOUSE_USERNAME")
	password := os.Getenv("CLICKHOUSE_PASSWORD")
	useCompression := os.Getenv("CLICKHOUSE_USE_COMPRESSION") == "true"

	client, err := NewClickhouseClient(host, portInt, username, password, useCompression, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create clickhouse client: %w", err)
	}
	return client, nil
}

func (c *ClickhouseClient) Close() error {
	return c.Conn.Close()
}

func (c *ClickhouseClient) Ping(ctx context.Context) error {
	return c.Conn.Ping(ctx)
}

func (c *ClickhouseClient) CreateSchemaCex(ctx context.Context) error {
	return c.Conn.Exec(ctx, queries.CreateSchemaCexClickhouseSQL)
}

func (c *ClickhouseClient) CreateBybitKlineTable(ctx context.Context) error {
	return c.Conn.Exec(ctx, queries.CreateTableBybitKlineSQL)
}

func (c *ClickhouseClient) CreateBybitFundingRateTable(ctx context.Context) error {
	return c.Conn.Exec(ctx, queries.CreateTableBybitFundingRateSQL)
}

func (c *ClickhouseClient) CreateBybitTickerInfoTable(ctx context.Context) error {
	return c.Conn.Exec(ctx, queries.CreateTableBybitTickerInfoClickhouseSQL)
}

func (c *ClickhouseClient) Query(ctx context.Context, query string, args ...interface{}) (driver.Rows, error) {
	rows, err := c.Conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *ClickhouseClient) QueryRow(ctx context.Context, query string, args ...interface{}) driver.Row {
	row := c.Conn.QueryRow(ctx, query, args...)
	if row.Err() != nil {
		return nil
	}
	return row
}

func (c *ClickhouseClient) Exec(ctx context.Context, query string, args ...interface{}) error {
	return c.Conn.Exec(ctx, query, args...)
}

func (c *ClickhouseClient) InsertBybitKline(ctx context.Context, kline *models.KlineInfo) error {
	openTime := time.UnixMilli(int64(kline.OpenTime)).UTC()
	return c.Conn.Exec(ctx, queries.InsertBybitKlineSQL,
		kline.TickerID,
		kline.OpenTime,
		openTime,
		kline.HighPrice,
		kline.LowPrice,
		kline.ClosePrice,
		kline.Volume,
		kline.Turnover,
	)
}

func (c *ClickhouseClient) InsertBybitKlineBatch(ctx context.Context, klines []*models.KlineInfo) error {
	if len(klines) == 0 {
		return nil
	}

	batch, err := c.Conn.PrepareBatch(ctx, queries.InsertBybitKlinesSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, kline := range klines {
		openTime := time.UnixMilli(int64(kline.OpenTime)).UTC()

		err := batch.Append(
			kline.TickerID,
			openTime,
			kline.OpenPrice,
			kline.HighPrice,
			kline.LowPrice,
			kline.ClosePrice,
			kline.Volume,
			kline.Turnover,
		)
		if err != nil {
			return fmt.Errorf("failed to append to batch: %w", err)
		}
	}

	return batch.Send()
}

func (c *ClickhouseClient) GetBybitKlineLatestOpenTime(ctx context.Context, tickerID uint64) (time.Time, error) {
	var openTime time.Time

	err := c.Conn.QueryRow(ctx, queries.SelectBybitKlineLatestOpenTimeSQL, tickerID).Scan(&openTime)
	if errors.Is(err, sql.ErrNoRows) {
		return time.Time{}, db.ErrNoRows
	}

	if err != nil {
		return time.Time{}, err
	}

	return openTime, nil
}

func (c *ClickhouseClient) GetBybitKlineData(ctx context.Context, tickerID uint64, startTime, endTime time.Time) ([]*models.KlineInfo, error) {
	rows, err := c.Conn.Query(ctx, queries.SelectBybitKlineDataRangeSQL, tickerID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var klines []*models.KlineInfo
	for rows.Next() {
		var kline models.KlineInfo
		var openTime time.Time

		err := rows.Scan(
			&kline.TickerID,
			&openTime,
			&kline.OpenPrice,
			&kline.HighPrice,
			&kline.LowPrice,
			&kline.ClosePrice,
			&kline.Volume,
			&kline.Turnover,
		)
		if err != nil {
			return nil, err
		}

		kline.OpenTime = uint64(openTime.UnixMilli())
		klines = append(klines, &kline)
	}

	return klines, rows.Err()
}

func (c *ClickhouseClient) GetBybitKlineDataAggregated(ctx context.Context, tickerID uint64, intervalMinutes int, startTime, endTime time.Time) ([]*models.KlineInfo, error) {
	rows, err := c.Conn.Query(
		ctx,
		queries.SelectBybitKlineDataAggregatedSQL,
		intervalMinutes,
		tickerID,
		startTime,
		endTime,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var klines []*models.KlineInfo
	for rows.Next() {
		var kline models.KlineInfo
		var openTime time.Time

		err := rows.Scan(
			&kline.TickerID,
			&openTime,
			&kline.OpenPrice,
			&kline.HighPrice,
			&kline.LowPrice,
			&kline.ClosePrice,
			&kline.Volume,
			&kline.Turnover,
		)
		if err != nil {
			return nil, err
		}

		kline.OpenTime = uint64(openTime.UnixMilli())
		klines = append(klines, &kline)
	}

	return klines, rows.Err()
}

func (c *ClickhouseClient) GetBybitKlineLatestClosePrice(ctx context.Context, tickerID uint64) (float64, error) {
	var price float64
	err := c.Conn.QueryRow(ctx, queries.SelectBybitKlineLatestCloseSQL, tickerID).Scan(&price)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (c *ClickhouseClient) InsertBybitFundingRate(ctx context.Context, fundingRate *models.FundingRateInfo) error {
	return c.Conn.Exec(ctx, queries.InsertBybitFundingRateSQL,
		fundingRate.TickerID,
		fundingRate.FundingRate,
		time.UnixMilli(int64(fundingRate.FundingRateTimestamp)).UTC(),
	)
}

func (c *ClickhouseClient) InsertBybitFundingRateBatch(ctx context.Context, fundingRates []*models.FundingRateInfo) error {
	if len(fundingRates) == 0 {
		return nil
	}

	batch, err := c.Conn.PrepareBatch(ctx, queries.InsertBybitFundingRatesSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, fundingRate := range fundingRates {
		err := batch.Append(
			fundingRate.TickerID,
			fundingRate.FundingRate,
			time.UnixMilli(int64(fundingRate.FundingRateTimestamp)).UTC(),
		)
		if err != nil {
			return fmt.Errorf("failed to append to batch: %w", err)
		}
	}

	return batch.Send()
}

func (c *ClickhouseClient) GetBybitFundingRateLatestFundingRateTimestamp(ctx context.Context, tickerID uint64) (uint64, error) {
	var timestamp time.Time

	err := c.Conn.QueryRow(ctx, queries.SelectBybitFundingRateLatestFundingRateTimestampSQL, tickerID).Scan(&timestamp)
	if err != nil {
		return 0, err
	}

	return uint64(timestamp.UnixMilli()), nil
}

func (c *ClickhouseClient) GetBybitFundingRateData(ctx context.Context, tickerID uint64, startTimestamp, endTimestamp uint64) ([]*models.FundingRateInfo, error) {
	rows, err := c.Conn.Query(ctx, queries.SelectBybitFundingRateDataRangeSQL, tickerID, startTimestamp, endTimestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fundingRates []*models.FundingRateInfo
	for rows.Next() {
		var fundingRate models.FundingRateInfo
		var fundingRateTimestamp time.Time

		err := rows.Scan(
			&fundingRate.TickerID,
			&fundingRate.FundingRate,
			&fundingRateTimestamp,
		)
		if err != nil {
			return nil, err
		}

		fundingRate.FundingRateTimestamp = uint64(fundingRateTimestamp.UnixMilli())

		fundingRates = append(fundingRates, &fundingRate)
	}

	return fundingRates, rows.Err()
}

func (c *ClickhouseClient) InsertBybitTickerInfo(ctx context.Context, ticker *models.TickerInfo) error {
	return c.Conn.Exec(ctx, queries.InsertBybitTickerInfoClickhouseSQL,
		ticker.ID,
		ticker.BaseCoin,
		ticker.Symbol,
		ticker.Category,
		ticker.BaseCoinParsed,
		ticker.Multiplier,
		ticker.QuoteCoin,
		ticker.Status,
	)
}

func (c *ClickhouseClient) InsertBybitTickerInfoBatch(ctx context.Context, tickers []models.TickerInfo) error {
	if len(tickers) == 0 {
		return nil
	}

	batch, err := c.Conn.PrepareBatch(ctx, queries.InsertBybitTickerInfosClickhouseSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, ticker := range tickers {
		err := batch.Append(
			ticker.ID,
			ticker.BaseCoin,
			ticker.Symbol,
			ticker.Category,
			ticker.BaseCoinParsed,
			ticker.Multiplier,
			ticker.QuoteCoin,
			ticker.Status,
		)
		if err != nil {
			return fmt.Errorf("failed to append to batch: %w", err)
		}
	}

	return batch.Send()
}

func (c *ClickhouseClient) GetAllBybitTickerInfo(ctx context.Context) ([]models.TickerInfo, error) {
	rows, err := c.Conn.Query(ctx, queries.SelectAllBybitTickerInfoClickhouseSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickers []models.TickerInfo
	for rows.Next() {
		var ticker models.TickerInfo
		err := rows.Scan(
			&ticker.ID,
			&ticker.BaseCoin,
			&ticker.Symbol,
			&ticker.Category,
			&ticker.BaseCoinParsed,
			&ticker.Multiplier,
			&ticker.QuoteCoin,
			&ticker.Status,
		)
		if err != nil {
			return nil, err
		}
		tickers = append(tickers, ticker)
	}

	return tickers, rows.Err()
}

func (c *ClickhouseClient) CreateBybitLiquidationTable(ctx context.Context) error {
	return c.Conn.Exec(ctx, queries.CreateTableBybitLiquidationSQL)
}

func (c *ClickhouseClient) CreateBybitTradeTable(ctx context.Context) error {
	return c.Conn.Exec(ctx, queries.CreateTableBybitTradeSQL)
}

func (c *ClickhouseClient) InsertBybitLiquidation(ctx context.Context, liq *models.LiquidationInfo) error {
	return c.Conn.Exec(ctx, queries.InsertBybitLiquidationSQL,
		liq.TickerID,
		liq.Size,
		liq.LiquidationPrice,
		time.UnixMilli(int64(liq.LiquidationTimestamp)).UTC(),
	)
}

func (c *ClickhouseClient) InsertBybitLiquidationBatch(ctx context.Context, liqs []*models.LiquidationInfo) error {
	if len(liqs) == 0 {
		return nil
	}

	batch, err := c.Conn.PrepareBatch(ctx, queries.InsertBybitLiquidationsSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, liq := range liqs {
		err := batch.Append(
			liq.TickerID,
			liq.Size,
			liq.LiquidationPrice,
			time.UnixMilli(int64(liq.LiquidationTimestamp)).UTC(),
		)
		if err != nil {
			return fmt.Errorf("failed to append to batch: %w", err)
		}
	}

	return batch.Send()
}

func (c *ClickhouseClient) GetBybitLiquidationLatestTimestamp(ctx context.Context, tickerID uint64) (uint64, error) {
	var ts time.Time
	if err := c.Conn.QueryRow(ctx, queries.SelectBybitLiquidationLatestTimestampSQL, tickerID).Scan(&ts); err != nil {
		return 0, err
	}
	return uint64(ts.UnixMilli()), nil
}

func (c *ClickhouseClient) GetBybitLiquidationData(ctx context.Context, tickerID uint64, startTimestamp, endTimestamp uint64) ([]*models.LiquidationInfo, error) {
	rows, err := c.Conn.Query(ctx, queries.SelectBybitLiquidationDataRangeSQL, tickerID, startTimestamp, endTimestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var liquidations []*models.LiquidationInfo
	for rows.Next() {
		var liq models.LiquidationInfo
		var ts time.Time
		if err := rows.Scan(
			&liq.TickerID,
			&liq.Size,
			&liq.LiquidationPrice,
			&ts,
		); err != nil {
			return nil, err
		}
		liq.LiquidationTimestamp = uint64(ts.UnixMilli())
		liquidations = append(liquidations, &liq)
	}

	return liquidations, rows.Err()
}

func (c *ClickhouseClient) InsertBybitTrade(ctx context.Context, trade *models.TradeInfo) error {
	return c.Conn.Exec(ctx, queries.InsertBybitTradeSQL,
		trade.TickerID,
		trade.Size,
		trade.Price,
		trade.TradeID,
		time.UnixMilli(int64(trade.TradeTimestamp)).UTC(),
	)
}

func (c *ClickhouseClient) InsertBybitTradeBatch(ctx context.Context, trades []*models.TradeInfo) error {
	if len(trades) == 0 {
		return nil
	}

	batch, err := c.Conn.PrepareBatch(ctx, queries.InsertBybitTradesSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, trade := range trades {
		err := batch.Append(
			trade.TickerID,
			trade.Size,
			trade.Price,
			trade.TradeID,
			time.UnixMilli(int64(trade.TradeTimestamp)).UTC(),
		)
		if err != nil {
			return fmt.Errorf("failed to append to batch: %w", err)
		}
	}

	return batch.Send()
}

func (c *ClickhouseClient) GetBybitTradeLatestTimestamp(ctx context.Context, tickerID uint64) (uint64, error) {
	var ts time.Time
	if err := c.Conn.QueryRow(ctx, queries.SelectBybitTradeLatestTimestampSQL, tickerID).Scan(&ts); err != nil {
		return 0, err
	}
	return uint64(ts.UnixMilli()), nil
}

func (c *ClickhouseClient) GetBybitTradeData(ctx context.Context, tickerID uint64, startTimestamp, endTimestamp uint64) ([]*models.TradeInfo, error) {
	rows, err := c.Conn.Query(ctx, queries.SelectBybitTradeDataRangeSQL, tickerID, startTimestamp, endTimestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []*models.TradeInfo
	for rows.Next() {
		var tr models.TradeInfo
		var ts time.Time
		if err := rows.Scan(
			&tr.TickerID,
			&tr.Size,
			&tr.Price,
			&tr.TradeID,
			&ts,
		); err != nil {
			return nil, err
		}
		tr.TradeTimestamp = uint64(ts.UnixMilli())
		trades = append(trades, &tr)
	}

	return trades, rows.Err()
}

func (c *ClickhouseClient) InsertListing(ctx context.Context, cexName string, ticker string, dt time.Time) error {
	return c.Conn.Exec(ctx, queries.InsertListingSQL, cexName, ticker, dt)
}

var _ db.DBClick = (*ClickhouseClient)(nil)
