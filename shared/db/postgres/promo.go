package postgres

import (
	"context"

	"github.com/monolit-network/xlayer-indexer/shared/db/postgres/queries"
	"github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/jackc/pgx/v5"
)

func scanPromoCode(row pgx.Row) (models.PromoCode, error) {
	var promo models.PromoCode
	var scopes []models.PromoScope
	var targetPlans []models.PlanCode
	err := row.Scan(
		&promo.Code,
		&promo.Active,
		&scopes,
		&promo.DiscountType,
		&promo.DiscountValue,
		&promo.MaxUses,
		&promo.UsesCount,
		&targetPlans,
		&promo.MinMonths,
		&promo.MaxMonths,
		&promo.MinCredits,
		&promo.MaxCredits,
		&promo.ExpiresAt,
		&promo.CreatedAt,
		&promo.UpdatedAt,
	)
	if err != nil {
		return models.PromoCode{}, err
	}
	promo.Scopes = scopes
	promo.TargetPlans = targetPlans
	return promo, nil
}

func scanPromoCodes(rows pgx.Rows) ([]models.PromoCode, error) {
	defer rows.Close()

	var promos []models.PromoCode
	for rows.Next() {
		promo, err := scanPromoCode(rows)
		if err != nil {
			return nil, err
		}
		promos = append(promos, promo)
	}
	return promos, rows.Err()
}

func (q *PostgresDB) GetBillingPromoCodeByCode(ctx context.Context, code string) (models.PromoCode, error) {
	return scanPromoCode(q.querier().QueryRow(ctx, queries.SelectBillingPromoCodeByCodeSQL, code))
}

func (q *PostgresDB) GetBillingPromoCodeByCodeForUpdate(ctx context.Context, code string) (models.PromoCode, error) {
	return scanPromoCode(q.querier().QueryRow(ctx, queries.SelectBillingPromoCodeByCodeForUpdateSQL, code))
}

func (q *PostgresDB) GetBillingPromoCodes(ctx context.Context) ([]models.PromoCode, error) {
	rows, err := q.querier().Query(ctx, queries.SelectBillingPromoCodesSQL)
	if err != nil {
		return nil, err
	}
	return scanPromoCodes(rows)
}

func (q *PostgresDB) InsertBillingPromoCode(ctx context.Context, promo models.PromoCode) error {
	_, err := q.querier().Exec(
		ctx,
		queries.InsertBillingPromoCodeSQL,
		promo.Code,
		promo.Active,
		promo.Scopes,
		promo.DiscountType,
		promo.DiscountValue,
		promo.MaxUses,
		promo.UsesCount,
		promo.TargetPlans,
		promo.MinMonths,
		promo.MaxMonths,
		promo.MinCredits,
		promo.MaxCredits,
		promo.ExpiresAt,
		promo.CreatedAt,
		promo.UpdatedAt,
	)
	return err
}

func (q *PostgresDB) UpdateBillingPromoCode(ctx context.Context, promo models.PromoCode) error {
	_, err := q.querier().Exec(
		ctx,
		queries.UpdateBillingPromoCodeSQL,
		promo.Code,
		promo.Active,
		promo.Scopes,
		promo.DiscountType,
		promo.DiscountValue,
		promo.MaxUses,
		promo.UsesCount,
		promo.TargetPlans,
		promo.MinMonths,
		promo.MaxMonths,
		promo.MinCredits,
		promo.MaxCredits,
		promo.ExpiresAt,
		promo.UpdatedAt,
	)
	return err
}
