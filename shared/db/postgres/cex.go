package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/monolit-network/xlayer-indexer/shared/db/postgres/queries"
	"github.com/monolit-network/xlayer-indexer/shared/models"
)

func (q *PostgresDB) CreateCexSchema(ctx context.Context) error {
	if _, err := q.querier().Exec(ctx, queries.CreateSchemaCexSQL); err != nil {
		return fmt.Errorf("failed to create cex schema: %w", err)
	}
	if _, err := q.querier().Exec(ctx, queries.CreateTableBybitTickerInfoSQL); err != nil {
		return fmt.Errorf("failed to create bybit_ticker_info table: %w", err)
	}
	return nil
}

func (q *PostgresDB) UpsertBybitTickerInfo(ctx context.Context, ticker models.TickerInfo) (uint64, error) {
	row := q.querier().QueryRow(ctx, queries.UpsertBybitTickerInfoSQL, ticker.BaseCoin, ticker.Symbol, ticker.Category, ticker.BaseCoinParsed, ticker.Multiplier, ticker.QuoteCoin, ticker.Status)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to upsert ticker info: %w", err)
	}
	return id, nil
}

func (q *PostgresDB) SelectBybitTickerInfoBySymbolAndCategory(ctx context.Context, symbol string, category string) (*models.TickerInfo, error) {
	row := q.querier().QueryRow(ctx, queries.SelectBybitTickerInfoBySymbolAndCategorySQL, symbol, category)
	var ticker models.TickerInfo
	if err := row.Scan(&ticker.ID, &ticker.BaseCoin, &ticker.Symbol, &ticker.Category, &ticker.BaseCoinParsed, &ticker.Multiplier, &ticker.QuoteCoin, &ticker.Status); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to select ticker info: %w", err)
	}
	return &ticker, nil
}

func (q *PostgresDB) SelectBybitTickerInfoByID(ctx context.Context, id uint64) (*models.TickerInfo, error) {
	row := q.querier().QueryRow(ctx, queries.SelectBybitTickerInfoByIDSQL, id)
	var ticker models.TickerInfo
	if err := row.Scan(&ticker.ID, &ticker.BaseCoin, &ticker.Symbol, &ticker.Category, &ticker.BaseCoinParsed, &ticker.Multiplier, &ticker.QuoteCoin, &ticker.Status); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to select ticker info: %w", err)
	}
	return &ticker, nil
}

func (q *PostgresDB) SelectAllBybitTickerInfo(ctx context.Context) ([]models.TickerInfo, error) {
	rows, err := q.querier().Query(ctx, queries.SelectAllBybitTickerInfoPostgresSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to query all ticker info: %w", err)
	}
	defer rows.Close()

	var tickers []models.TickerInfo
	for rows.Next() {
		var ticker models.TickerInfo
		if err := rows.Scan(&ticker.ID, &ticker.BaseCoin, &ticker.Symbol, &ticker.Category, &ticker.BaseCoinParsed, &ticker.Multiplier, &ticker.QuoteCoin, &ticker.Status); err != nil {
			return nil, fmt.Errorf("failed to scan ticker info: %w", err)
		}
		tickers = append(tickers, ticker)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ticker info rows: %w", err)
	}

	return tickers, nil
}
