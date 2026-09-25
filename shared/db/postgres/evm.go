package postgres

import (
	"context"
	"fmt"

	"github.com/monolit-network/xlayer-indexer/shared/db/postgres/queries"
)

func (q *PostgresDB) CreateEvmParserRegistryTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, queries.CreateEvmParserRegistryTableSQL)
	return err
}

func (q *PostgresDB) UpsertEvmParser(ctx context.Context, chain string, contractAddress string, instructionHash string, parserName string, source string) error {
	_, err := q.querier().Exec(ctx, queries.UpsertEvmParserSQL, chain, contractAddress, instructionHash, parserName, source)
	return err
}

func (q *PostgresDB) GetEvmParser(ctx context.Context, chain string, contractAddress string, instructionHash string) (string, string, error) {
	row := q.querier().QueryRow(ctx, queries.GetEvmParserSQL, chain, contractAddress, instructionHash)
	var parserName string
	var source string
	if err := row.Scan(&parserName, &source); err != nil {
		return "", "", fmt.Errorf("failed to get evm parser: %w", err)
	}
	return parserName, source, nil
}

func (q *PostgresDB) GetAllEvmParsersForChain(ctx context.Context, chain string) (map[string]map[string][2]string, error) {
	rows, err := q.querier().Query(ctx, queries.GetAllEvmParsersForChainSQL, chain)
	if err != nil {
		return nil, fmt.Errorf("failed to get all evm parsers for chain: %w", err)
	}
	defer rows.Close()

	parsers := make(map[string]map[string][2]string)
	for rows.Next() {
		var contractAddress string
		var instructionHash string
		var parserName string
		var source string
		if err := rows.Scan(&contractAddress, &instructionHash, &parserName, &source); err != nil {
			return nil, fmt.Errorf("failed to get evm parser: %w", err)
		}
		if _, ok := parsers[contractAddress]; !ok {
			parsers[contractAddress] = make(map[string][2]string)
		}
		parsers[contractAddress][instructionHash] = [2]string{parserName, source}
	}
	return parsers, nil
}
