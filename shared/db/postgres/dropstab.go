package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/monolit-network/xlayer-indexer/shared/db/postgres/queries"
	"github.com/monolit-network/xlayer-indexer/shared/models"
)

func (q *PostgresDB) CreateVerifiedTokensTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, queries.CreateVerifiedTokensTableSQL)
	return err
}

func (q *PostgresDB) CreateDropstabInfoTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, queries.CreateDropstabInfoTableSQL)
	return err
}

func (q *PostgresDB) GetVerifiedTokenByAssetIDChain(ctx context.Context, assetID models.VerifiedTokenAssetID, chain string) (*models.VerifiedToken, error) {
	var token models.VerifiedToken

	err := q.querier().QueryRow(ctx, queries.SelectVerifiedTokenByAssetIDChainSQL, assetID, chain).Scan(
		&token.AssetID,
		&token.Chain,
		&token.ContractAddress,
		&token.Symbol,
		&token.Name,
		&token.Decimals,
		&token.CreatedAt,
		&token.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (q *PostgresDB) GetVerifiedTokenByChainContractAddress(ctx context.Context, chain string, contractAddress string) (*models.VerifiedToken, error) {
	var token models.VerifiedToken

	err := q.querier().QueryRow(ctx, queries.SelectVerifiedTokenByChainContractAddressSQL, chain, contractAddress).Scan(
		&token.AssetID,
		&token.Chain,
		&token.ContractAddress,
		&token.Symbol,
		&token.Name,
		&token.Decimals,
		&token.CreatedAt,
		&token.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (q *PostgresDB) GetVerifiedTokensByAssetID(ctx context.Context, assetID models.VerifiedTokenAssetID) ([]models.VerifiedToken, error) {
	rows, err := q.querier().Query(ctx, queries.SelectVerifiedTokensByAssetIDSQL, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []models.VerifiedToken
	for rows.Next() {
		var token models.VerifiedToken
		if err := rows.Scan(
			&token.AssetID,
			&token.Chain,
			&token.ContractAddress,
			&token.Symbol,
			&token.Name,
			&token.Decimals,
			&token.CreatedAt,
			&token.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

func (q *PostgresDB) GetVerifiedTokensByContractAddress(ctx context.Context, contractAddress string) ([]models.VerifiedToken, error) {
	rows, err := q.querier().Query(ctx, queries.SelectVerifiedTokensByContractAddressSQL, contractAddress)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []models.VerifiedToken
	for rows.Next() {
		var token models.VerifiedToken
		if err := rows.Scan(
			&token.AssetID,
			&token.Chain,
			&token.ContractAddress,
			&token.Symbol,
			&token.Name,
			&token.Decimals,
			&token.CreatedAt,
			&token.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

func (q *PostgresDB) GetVerifiedTokensBySymbol(ctx context.Context, symbol string) ([]models.VerifiedToken, error) {
	rows, err := q.querier().Query(ctx, queries.SelectVerifiedTokensBySymbolSQL, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []models.VerifiedToken

	for rows.Next() {
		var token models.VerifiedToken
		if err := rows.Scan(
			&token.AssetID,
			&token.Chain,
			&token.ContractAddress,
			&token.Symbol,
			&token.Name,
			&token.Decimals,
			&token.CreatedAt,
			&token.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

func (q *PostgresDB) UpsertVerifiedToken(ctx context.Context, token models.VerifiedToken) error {
	var err error
	_, err = q.querier().Exec(ctx, queries.UpsertVerifiedTokenSQL,
		token.AssetID,
		token.Chain,
		token.ContractAddress,
		token.Symbol,
		token.Name,
		token.Decimals,
		token.CreatedAt,
		token.UpdatedAt,
	)

	return err
}

func (q *PostgresDB) DeleteVerifiedToken(ctx context.Context, assetID models.VerifiedTokenAssetID, chain string) error {
	result, err := q.querier().Exec(ctx, queries.DeleteVerifiedTokenSQL, assetID, chain)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("no verified token found with asset_id %d and chain %s", assetID, chain)
	}

	return nil
}

func (q *PostgresDB) DeleteVerifiedTokensByAssetID(ctx context.Context, assetID models.VerifiedTokenAssetID) error {
	_, err := q.querier().Exec(ctx, queries.DeleteVerifiedTokensByAssetIDSQL, assetID)
	if err != nil {
		return err
	}
	return nil
}

func (q *PostgresDB) scanDropstabInfoRow(row pgx.Row) (*models.DropstabInfo, error) {
	var info models.DropstabInfo
	var socialsRaw []byte
	var categories []string
	err := row.Scan(
		&info.AssetID,
		&info.Slug,
		&info.Status,
		&info.Symbol,
		&info.Name,
		&socialsRaw,
		&info.Description,
		&info.ImageURL,
		&info.PriceUSD,
		&info.MaxSupply,
		&info.CirculatingSupply,
		&info.TotalSupply,
		&categories,
		&info.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(socialsRaw, &info.Socials); err != nil {
		return nil, err
	}

	info.Categories = categories

	return &info, nil
}

func (q *PostgresDB) scanDropstabInfoRows(rows pgx.Rows) (map[models.VerifiedTokenAssetID]*models.DropstabInfo, error) {
	dropstabInfos := make(map[models.VerifiedTokenAssetID]*models.DropstabInfo)

	for rows.Next() {
		var info models.DropstabInfo
		var socialsRaw []byte
		var categories []string
		if err := rows.Scan(
			&info.AssetID,
			&info.Slug,
			&info.Status,
			&info.Symbol,
			&info.Name,
			&socialsRaw,
			&info.Description,
			&info.ImageURL,
			&info.PriceUSD,
			&info.MaxSupply,
			&info.CirculatingSupply,
			&info.TotalSupply,
			&categories,
			&info.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(socialsRaw, &info.Socials); err != nil {
			return nil, err
		}

		info.Categories = categories

		dropstabInfos[info.AssetID] = &info
	}

	return dropstabInfos, nil
}

func (q *PostgresDB) GetDropstabInfoByAssetID(ctx context.Context, assetID models.VerifiedTokenAssetID) (*models.DropstabInfo, error) {
	row := q.querier().QueryRow(ctx, queries.SelectDropstabInfoByAssetIDSQL, assetID)
	info, err := q.scanDropstabInfoRow(row)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (q *PostgresDB) GetDropstabInfosByAssetIDs(ctx context.Context, assetIDs []models.VerifiedTokenAssetID) (map[models.VerifiedTokenAssetID]*models.DropstabInfo, error) {
	dropstabInfos := make(map[models.VerifiedTokenAssetID]*models.DropstabInfo)
	for _, assetID := range assetIDs {
		dropstabInfos[assetID] = nil
	}

	assetIDsInt64 := make([]int64, len(assetIDs))
	for i, assetID := range assetIDs {
		assetIDsInt64[i] = int64(assetID)
	}

	rows, err := q.querier().Query(ctx, queries.SelectDropstabInfosByAssetIDsSQL, assetIDsInt64)
	switch err {
	case pgx.ErrNoRows:
		return dropstabInfos, nil
	case nil:
	default:
		return nil, err
	}
	defer rows.Close()

	scannedInfos, err := q.scanDropstabInfoRows(rows)
	if err != nil {
		return nil, err
	}

	for assetID, info := range scannedInfos {
		dropstabInfos[assetID] = info
	}

	return dropstabInfos, nil
}

func (q *PostgresDB) GetDropstabInfoBySymbol(ctx context.Context, symbol string) (map[models.VerifiedTokenAssetID]*models.DropstabInfo, error) {
	rows, err := q.querier().Query(ctx, queries.SelectDropstabInfoBySymbolSQL, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scannedInfos, err := q.scanDropstabInfoRows(rows)
	if err != nil {
		return nil, err
	}

	return scannedInfos, nil
}

func (q *PostgresDB) GetDropstabInfoBySlug(ctx context.Context, slug string) (*models.DropstabInfo, error) {
	row := q.querier().QueryRow(ctx, queries.SelectDropstabInfoBySlugSQL, slug)
	info, err := q.scanDropstabInfoRow(row)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (q *PostgresDB) GetAssetIDsBySlugs(ctx context.Context, slugs []string) (map[string]models.VerifiedTokenAssetID, error) {
	rows, err := q.querier().Query(ctx, queries.SelectAssetIDsBySlugsSQL, slugs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assetIDs := make(map[string]models.VerifiedTokenAssetID)
	for rows.Next() {
		var slug string
		var assetID models.VerifiedTokenAssetID
		if err := rows.Scan(&slug, &assetID); err != nil {
			return nil, err
		}
		assetIDs[slug] = assetID
	}

	return assetIDs, nil
}

func (q *PostgresDB) GetExistingDropstabInfoSlugs(ctx context.Context, slugs []string) ([]models.DropstabShortCoinDescription, error) {
	rows, err := q.querier().Query(ctx, queries.SelectExistingDropstabInfoSlugsSQL, slugs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		if err == pgx.ErrNoRows {
			return []models.DropstabShortCoinDescription{}, nil
		}
		return nil, err
	}

	var existingSlugs []models.DropstabShortCoinDescription
	for rows.Next() {
		var slug string
		var updatedAt time.Time
		if err := rows.Scan(&slug, &updatedAt); err != nil {
			return nil, err
		}
		existingSlugs = append(existingSlugs, models.DropstabShortCoinDescription{
			Slug:      slug,
			UpdatedAt: updatedAt,
		})
	}

	return existingSlugs, nil
}

func (q *PostgresDB) UpsertDropstabInfo(ctx context.Context, info models.DropstabInfo) (models.VerifiedTokenAssetID, error) {
	socialsRaw, err := json.Marshal(info.Socials)
	if err != nil {
		return 0, err
	}

	if info.AssetID == 0 {
		info.AssetID, err = q.getNextAssetID(ctx)
		if err != nil {
			return 0, err
		}
	}

	_, err = q.querier().Exec(ctx, queries.UpsertDropstabInfoSQL,
		info.AssetID,
		info.Slug,
		info.Status,
		info.Symbol,
		info.Name,
		socialsRaw,
		info.Description,
		info.ImageURL,
		info.PriceUSD,
		info.MaxSupply,
		info.CirculatingSupply,
		info.TotalSupply,
		info.Categories,
		info.UpdatedAt,
	)

	if err != nil {
		return 0, err
	}

	return info.AssetID, nil
}

func (q *PostgresDB) DeleteDropstabInfo(ctx context.Context, assetID models.VerifiedTokenAssetID) error {
	result, err := q.querier().Exec(ctx, queries.DeleteDropstabInfoSQL, assetID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("no dropstab info found with asset_id %d", assetID)
	}

	return nil
}
