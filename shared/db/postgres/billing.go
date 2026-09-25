package postgres

import (
	"context"
	"time"

	evmmodels "github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/shared/db/postgres/queries"
	"github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/jackc/pgx/v5"
)

func (q *PostgresDB) CreateBillingUsersTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, queries.CreateTableBillingUsersSQL)
	return err
}

func (q *PostgresDB) CreateBillingTransactionsHistoryTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, queries.CreateTableBillingTransactionsHistorySQL)
	return err
}

func (q *PostgresDB) CreateBillingOrdersTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, queries.CreateTableBillingOrdersSQL)
	return err
}

func (q *PostgresDB) CreateBillingPromoCodesTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, queries.CreateTableBillingPromoCodesSQL)
	return err
}

func scanBillingUserState(row pgx.Row) (models.BillingUserState, error) {
	var state models.BillingUserState
	err := row.Scan(
		&state.UserID,
		&state.APIKeys,
		&state.Plan,
		&state.PlanCredits,
		&state.AdditionalCredits,
		&state.BillingPeriodPaidCents,
		&state.PlanStartedAt,
		&state.NextPlanRefillAt,
		&state.PlanEndingAt,
	)
	return state, err
}

func (q *PostgresDB) GetBillingUserState(ctx context.Context, userID string) (models.BillingUserState, error) {
	return scanBillingUserState(q.querier().QueryRow(ctx, queries.SelectBillingUserStateSQL, userID))
}

func (q *PostgresDB) GetBillingUserStateForUpdate(ctx context.Context, userID string) (models.BillingUserState, error) {
	return scanBillingUserState(q.querier().QueryRow(ctx, queries.SelectBillingUserStateForUpdateSQL, userID))
}

func (q *PostgresDB) GetBillingUserStateForUpdateNowait(ctx context.Context, userID string) (models.BillingUserState, error) {
	return scanBillingUserState(q.querier().QueryRow(ctx, queries.SelectBillingUserStateForUpdateNowaitSQL, userID))
}

func (q *PostgresDB) GetBillingUserPlan(ctx context.Context, userID string) (models.PlanCode, error) {
	row := q.querier().QueryRow(ctx, queries.SelectBillingUserPlanSQL, userID)
	var plan models.PlanCode
	if err := row.Scan(&plan); err != nil {
		return "", err
	}
	return plan, nil
}

func (q *PostgresDB) GetBillingUserPlanCredits(ctx context.Context, userID string) (int64, error) {
	row := q.querier().QueryRow(ctx, queries.SelectBillingUserPlanCreditsSQL, userID)
	var planCredits int64
	if err := row.Scan(&planCredits); err != nil {
		return 0, err
	}
	return planCredits, nil
}

func (q *PostgresDB) GetBillingUserAdditionalCredits(ctx context.Context, userID string) (int64, error) {
	row := q.querier().QueryRow(ctx, queries.SelectBillingUserAdditionalCreditsSQL, userID)
	var additionalCredits int64
	if err := row.Scan(&additionalCredits); err != nil {
		return 0, err
	}
	return additionalCredits, nil
}

func (q *PostgresDB) UpdateBillingUserState(ctx context.Context, state models.BillingUserState) error {
	_, err := q.querier().Exec(
		ctx,
		queries.UpdateBillingUserStateSQL,
		state.UserID,
		state.Plan,
		state.PlanCredits,
		state.AdditionalCredits,
		state.BillingPeriodPaidCents,
		state.PlanStartedAt,
		state.NextPlanRefillAt,
		state.PlanEndingAt,
	)
	return err
}

func (q *PostgresDB) UpdateBillingUserPlanCreditsAndAdditionalCreditsBatch(ctx context.Context, users []models.UserCredits) error {
	tx, err := q.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, queries.CreateTempTableBillingUserPlanCreditsAndAdditionalCreditsBatchSQL)
	if err != nil {
		return err
	}

	rows := make([][]any, 0, len(users))
	for _, u := range users {
		rows = append(rows, []any{u.UserID, u.PlanCredits, u.AdditionalCredits})
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"temp_credits"},
		[]string{"user_id", "plan_credits", "additional_credits"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, queries.SetBillingUserPlanCreditsAndAdditionalCreditsBatchSQL)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (q *PostgresDB) GetBillingApiKeysByUserID(ctx context.Context, userID string) ([]string, error) {
	row := q.querier().QueryRow(ctx, queries.SelectBillingApiKeysByUserIDSQL, userID)
	var apiKeys []string
	if err := row.Scan(&apiKeys); err != nil {
		return nil, err
	}
	return apiKeys, nil
}

func (q *PostgresDB) GetBillingUserIDByAPIKey(ctx context.Context, APIKey string) (string, error) {
	row := q.querier().QueryRow(ctx, queries.SelectUserIDByAPIKeySQL, APIKey)
	var userID string
	if err := row.Scan(&userID); err != nil {
		return "", err
	}
	return userID, nil
}

func (q *PostgresDB) InsertBillingAPIKeyForUser(ctx context.Context, userID, APIKey string) error {
	_, err := q.querier().Exec(ctx, queries.InsertBillingAPIKeyForUserSQL, APIKey, userID)
	return err
}

func (q *PostgresDB) DeleteBillingAPIKeyFromUser(ctx context.Context, userID, APIKey string) error {
	_, err := q.querier().Exec(ctx, queries.DeleteBillingAPIKeyFromUserSQL, APIKey, userID)
	return err
}

func (q *PostgresDB) InsertBillingUser(ctx context.Context, userID string) error {
	_, err := q.querier().Exec(ctx, queries.InsertBillingUserSQL, userID)
	return err
}

func (q *PostgresDB) GetBillingPlans(ctx context.Context) ([]models.PlanConfig, error) {
	rows, err := q.querier().Query(ctx, queries.SelectBillingPlansSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.PlanConfig
	for rows.Next() {
		var plan models.PlanConfig
		if err := rows.Scan(&plan.RequestsPerMinute, &plan.Credits, &plan.Code, &plan.Features, &plan.PriceUSD, &plan.Name); err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

func scanBillingUsers(rows pgx.Rows) ([]models.BillingUserState, error) {
	defer rows.Close()

	var users []models.BillingUserState
	for rows.Next() {
		var user models.BillingUserState
		if err := rows.Scan(
			&user.UserID,
			&user.APIKeys,
			&user.Plan,
			&user.PlanCredits,
			&user.AdditionalCredits,
			&user.BillingPeriodPaidCents,
			&user.PlanStartedAt,
			&user.NextPlanRefillAt,
			&user.PlanEndingAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (q *PostgresDB) GetBillingUsersForMonthlyRefill(ctx context.Context) ([]models.BillingUserState, error) {
	rows, err := q.querier().Query(ctx, queries.SelectBillingUsersForMonthlyRefillSQL)
	if err != nil {
		return nil, err
	}
	return scanBillingUsers(rows)
}

func (q *PostgresDB) RefillBillingUsersDailyFreeBatch(ctx context.Context, planCredits int64, nextRefillAt, now time.Time, limit int) ([]string, error) {
	rows, err := q.querier().Query(ctx, queries.RefillBillingUsersDailyFreeBatchSQL, planCredits, nextRefillAt, now, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userIDs := make([]string, 0, limit)
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, rows.Err()
}

func (q *PostgresDB) GetBillingUsersForExpiration(ctx context.Context) ([]models.BillingUserState, error) {
	rows, err := q.querier().Query(ctx, queries.SelectBillingUsersForExpirationSQL)
	if err != nil {
		return nil, err
	}
	return scanBillingUsers(rows)
}

func (q *PostgresDB) GetBillingTransactionHistory(ctx context.Context, userID string, eventKind models.BillingEventKind) ([]models.BillingTransactionHistory, error) {
	rows, err := q.querier().Query(ctx, queries.SelectBillingTransactionHistorySQL, userID, eventKind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []models.BillingTransactionHistory{}
	for rows.Next() {
		var transaction models.BillingTransactionHistory
		err := rows.Scan(&transaction.OrderID, &transaction.UserID, &transaction.EventKind, &transaction.AmountUSDCents, &transaction.CreatedAt, &transaction.AdditionalInfo)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (q *PostgresDB) GetBillingTransactionHistoryByOrderID(ctx context.Context, orderID string) (models.BillingTransactionHistory, error) {
	row := q.querier().QueryRow(ctx, queries.SelectBillingTransactionHistoryByOrderIDSQL, orderID)
	var transaction models.BillingTransactionHistory
	err := row.Scan(&transaction.OrderID, &transaction.UserID, &transaction.EventKind, &transaction.AmountUSDCents, &transaction.CreatedAt, &transaction.AdditionalInfo)
	if err != nil {
		return models.BillingTransactionHistory{}, err
	}
	return transaction, nil
}

func (q *PostgresDB) InsertBillingTransactionHistory(ctx context.Context, transaction models.BillingTransactionHistory) error {
	if transaction.AdditionalInfo == nil {
		_, err := q.querier().Exec(ctx, queries.InsertBillingTransactionHistoryWithoutAdditionalInfoSQL, transaction.OrderID, transaction.UserID, transaction.EventKind, transaction.AmountUSDCents, transaction.CreatedAt)
		return err
	}
	_, err := q.querier().Exec(ctx, queries.InsertBillingTransactionHistorySQL, transaction.OrderID, transaction.UserID, transaction.EventKind, transaction.AmountUSDCents, transaction.CreatedAt, transaction.AdditionalInfo)
	return err
}

func scanBillingOrder(row pgx.Row) (models.BillingOrder, error) {
	var order models.BillingOrder
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.Type,
		&order.Chain,
		&order.TokenAddress,
		&order.AmountUSDCents,
		&order.WalletAddress,
		&order.PromoCode,
		&order.ProductData,
		&order.ExpiresAt,
		&order.ScanFromBlock,
		&order.PaidTxHash,
		&order.PaidBlockNumber,
		&order.PaymentProvider,
		&order.PaymentPayloadHash,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	return order, err
}

func scanBillingOrders(rows pgx.Rows) ([]models.BillingOrder, error) {
	defer rows.Close()

	var orders []models.BillingOrder
	for rows.Next() {
		order, err := scanBillingOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (q *PostgresDB) InsertBillingOrder(ctx context.Context, order models.BillingOrder) error {
	_, err := q.querier().Exec(
		ctx,
		queries.InsertBillingOrderSQL,
		order.ID,
		order.UserID,
		order.Status,
		order.Type,
		order.Chain,
		order.TokenAddress,
		order.AmountUSDCents,
		order.WalletAddress,
		order.PromoCode,
		order.ProductData,
		order.ExpiresAt,
		order.ScanFromBlock,
		order.PaidTxHash,
		order.PaidBlockNumber,
		order.PaymentProvider,
		order.PaymentPayloadHash,
		order.CreatedAt,
		order.UpdatedAt,
	)
	return err
}

func (q *PostgresDB) GetBillingOrdersByUserID(ctx context.Context, userID string) ([]models.BillingOrder, error) {
	rows, err := q.querier().Query(ctx, queries.SelectBillingOrdersByUserIDSQL, userID)
	if err != nil {
		return nil, err
	}
	return scanBillingOrders(rows)
}

func (q *PostgresDB) GetBillingOrderByID(ctx context.Context, orderID string) (models.BillingOrder, error) {
	return scanBillingOrder(q.querier().QueryRow(ctx, queries.SelectBillingOrderByIDSQL, orderID))
}

func (q *PostgresDB) GetBillingOrderByPaidTxHash(ctx context.Context, txHash string, provider models.BillingPaymentProvider) (models.BillingOrder, error) {
	return scanBillingOrder(q.querier().QueryRow(ctx, queries.SelectBillingOrderByPaidTxHashSQL, txHash, provider))
}

func (q *PostgresDB) GetBillingOrderByPaymentPayloadHash(ctx context.Context, paymentPayloadHash string, provider models.BillingPaymentProvider) (models.BillingOrder, error) {
	return scanBillingOrder(q.querier().QueryRow(ctx, queries.SelectBillingOrderByPaymentPayloadHashSQL, paymentPayloadHash, provider))
}

func (q *PostgresDB) GetBillingOrderByIDForUpdate(ctx context.Context, orderID string) (models.BillingOrder, error) {
	return scanBillingOrder(q.querier().QueryRow(ctx, queries.SelectBillingOrderByIDForUpdateSQL, orderID))
}

func (q *PostgresDB) GetX402PaymentSettledBillingOrders(ctx context.Context) ([]models.BillingOrder, error) {
	rows, err := q.querier().Query(ctx, queries.SelectX402PaymentSettledBillingOrdersSQL)
	if err != nil {
		return nil, err
	}
	return scanBillingOrders(rows)
}

func (q *PostgresDB) CountPendingBillingOrdersByUserAndTypes(ctx context.Context, userID string, types []models.BillingOrderType) (int64, error) {
	row := q.querier().QueryRow(ctx, queries.SelectPendingBillingOrderCountByUserAndTypesSQL, userID, types)
	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (q *PostgresDB) GetExpiredPendingBillingOrders(ctx context.Context) ([]models.BillingOrder, error) {
	rows, err := q.querier().Query(ctx, queries.SelectExpiredPendingBillingOrdersSQL)
	if err != nil {
		return nil, err
	}
	return scanBillingOrders(rows)
}

func (q *PostgresDB) GetPendingBillingOrderChains(ctx context.Context) ([]evmmodels.Chain, error) {
	rows, err := q.querier().Query(ctx, queries.SelectPendingBillingOrderChainsSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chains []evmmodels.Chain
	for rows.Next() {
		var chain evmmodels.Chain
		if err := rows.Scan(&chain); err != nil {
			return nil, err
		}
		chains = append(chains, chain)
	}
	return chains, rows.Err()
}

func (q *PostgresDB) GetPendingBillingOrdersBatchWithSameScanFromBlock(ctx context.Context, chainName evmmodels.Chain, safeBlock uint64, limit int64) ([]models.BillingOrder, error) {
	rows, err := q.querier().Query(ctx, queries.SelectPendingBillingOrderBatchWithSameScanFromBlockSQL, chainName, safeBlock, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.BillingOrder
	for rows.Next() {
		var order models.BillingOrder
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.Type,
			&order.Chain,
			&order.TokenAddress,
			&order.AmountUSDCents,
			&order.WalletAddress,
			&order.PromoCode,
			&order.ProductData,
			&order.ExpiresAt,
			&order.ScanFromBlock,
			&order.PaidTxHash,
			&order.PaidBlockNumber,
			&order.PaymentProvider,
			&order.PaymentPayloadHash,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (q *PostgresDB) GetFailedBillingOrders(ctx context.Context) ([]models.BillingOrder, error) {
	rows, err := q.querier().Query(ctx, queries.SelectFailedBillingOrdersSQL)
	if err != nil {
		return nil, err
	}
	return scanBillingOrders(rows)
}

func (q *PostgresDB) UpdatePendingBillingOrdersScanFromBlock(ctx context.Context, orderIDs []string, scanFromBlock uint64, updatedAt time.Time) error {
	if len(orderIDs) == 0 {
		return nil
	}
	_, err := q.querier().Exec(ctx, queries.UpdatePendingBillingOrdersScanFromBlockSQL, orderIDs, scanFromBlock, updatedAt)
	return err
}

func (q *PostgresDB) UpdateBillingOrderStatus(ctx context.Context, orderID string, status models.BillingOrderStatus, paidTxHash *string, paidBlockNumber *uint64, updatedAt time.Time) error {
	_, err := q.querier().Exec(ctx, queries.UpdateBillingOrderStatusSQL, orderID, status, paidTxHash, paidBlockNumber, updatedAt)
	return err
}
