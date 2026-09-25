package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/monolit-network/xlayer-indexer/shared/db/postgres/queries"
	"github.com/monolit-network/xlayer-indexer/shared/models"
)

func (q *PostgresDB) CreateUnverifiedTokensTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, queries.CreateUnverifiedTokensTableSQL)
	return err
}

func (q *PostgresDB) GetUnverifiedTokenByChainContractAddress(ctx context.Context, chain string, contractAddress string) (*models.Token, error) {
	row := q.querier().QueryRow(ctx, queries.SelectUnverifiedTokenByChainContractAddressSQL, chain, contractAddress)
	token, err := scanUnverifiedTokenRow(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (q *PostgresDB) GetUnverifiedTokensBySymbol(ctx context.Context, symbol string) ([]models.Token, error) {
	rows, err := q.querier().Query(ctx, queries.SelectUnverifiedTokensBySymbolSQL, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := make([]models.Token, 0, 8)
	for rows.Next() {
		token, err := scanUnverifiedTokenRows(rows)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

func (q *PostgresDB) GetUnverifiedTokensByContractAddress(ctx context.Context, contractAddress string) ([]models.Token, error) {
	rows, err := q.querier().Query(ctx, queries.SelectUnverifiedTokensByContractAddressSQL, contractAddress)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := make([]models.Token, 0, 4)
	for rows.Next() {
		token, err := scanUnverifiedTokenRows(rows)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

func (q *PostgresDB) GetUnverifiedTokensByContractAddresses(ctx context.Context, addresses []string) ([]models.Token, error) {
	if len(addresses) == 0 {
		return nil, nil
	}

	norm := make([]string, 0, len(addresses))
	seen := make(map[string]struct{}, len(addresses))
	for _, addr := range addresses {
		a := strings.ToLower(strings.TrimSpace(addr))
		if a == "" {
			continue
		}
		if _, ok := seen[a]; ok {
			continue
		}
		seen[a] = struct{}{}
		norm = append(norm, a)
	}
	if len(norm) == 0 {
		return nil, nil
	}

	rows, err := q.querier().Query(ctx, queries.SelectUnverifiedTokensByContractAddressesSQL, norm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := make([]models.Token, 0, len(norm))
	for rows.Next() {
		token, err := scanUnverifiedTokenRows(rows)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

func (q *PostgresDB) UpsertUnverifiedToken(ctx context.Context, token models.Token) error {
	_, err := q.querier().Exec(ctx, queries.UpsertUnverifiedTokenSQL,
		token.ContractAddress,
		token.Chain,
		token.Decimals,
		token.Symbol,
		token.Name,
		token.PriceUSD,
		token.MarketCapUSD,
		token.Supply,
		token.LargestLPPoolUSD,
		token.FirstTxDate,
		token.CreatedAt,
		token.UpdatedAt,
		token.ViewSource,
	)
	return err
}

func scanUnverifiedTokenRow(row pgx.Row) (*models.Token, error) {
	var token models.Token
	var firstTxDate time.Time
	var createdAt time.Time
	var updatedAt time.Time
	var name *string
	var symbol *string
	var decimals *int64

	if err := row.Scan(
		&token.Chain,
		&token.ContractAddress,
		&symbol,
		&name,
		&decimals,
		&token.PriceUSD,
		&token.MarketCapUSD,
		&token.Supply,
		&token.LargestLPPoolUSD,
		&firstTxDate,
		&createdAt,
		&updatedAt,
		&token.ViewSource,
	); err != nil {
		return nil, err
	}

	if symbol != nil {
		token.Symbol = *symbol
	}
	if decimals != nil {
		token.Decimals = int(*decimals)
	}
	if name != nil {
		token.Name = *name
	}

	return &token, nil
}

func scanUnverifiedTokenRows(rows pgx.Rows) (models.Token, error) {
	var token models.Token
	var firstTxDate time.Time
	var createdAt time.Time
	var updatedAt time.Time
	var name *string
	var symbol *string
	var decimals *int64

	if err := rows.Scan(
		&token.Chain,
		&token.ContractAddress,
		&symbol,
		&name,
		&decimals,
		&token.PriceUSD,
		&token.MarketCapUSD,
		&token.Supply,
		&token.LargestLPPoolUSD,
		&firstTxDate,
		&createdAt,
		&updatedAt,
		&token.ViewSource,
	); err != nil {
		return models.Token{}, err
	}

	if symbol != nil {
		token.Symbol = *symbol
	}
	if decimals != nil {
		token.Decimals = int(*decimals)
	}
	if name != nil {
		token.Name = *name
	}

	return token, nil
}
