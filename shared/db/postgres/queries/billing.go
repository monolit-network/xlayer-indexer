package queries

const (
	CreateTableBillingUsersSQL = `
	CREATE TABLE IF NOT EXISTS backend.billing_plans (
		plan_code VARCHAR(255) PRIMARY KEY,
		plan_price INTEGER NOT NULL,
		plan_name VARCHAR(255) NOT NULL,
		plan_credits BIGINT NOT NULL,
		plan_requests_per_minute INTEGER NOT NULL,
		plan_features TEXT[] NOT NULL
	);

	CREATE TABLE IF NOT EXISTS backend.billing_users (
		user_id uuid PRIMARY KEY,
		api_keys uuid[],
		plan_credits BIGINT NOT NULL DEFAULT 0,
		additional_credits BIGINT NOT NULL DEFAULT 0,
		billing_period_paid_cents INTEGER NOT NULL DEFAULT 0,
		plan VARCHAR(255) NOT NULL,
		plan_started_at TIMESTAMP,
		next_plan_refill_at TIMESTAMP,
		plan_ending_at TIMESTAMP,
		CONSTRAINT fk_backend_billing_users_user_id FOREIGN KEY (user_id) REFERENCES backend.users (id) ON DELETE CASCADE,
		CONSTRAINT fk_backend_billing_users_plan FOREIGN KEY (plan) REFERENCES backend.billing_plans (plan_code)
	);

	CREATE INDEX IF NOT EXISTS idx_backend_billing_users_api_keys ON backend.billing_users USING GIN (api_keys);
	CREATE INDEX IF NOT EXISTS idx_backend_billing_users_plan_ending_at ON backend.billing_users (plan_ending_at);
	CREATE INDEX IF NOT EXISTS idx_backend_billing_users_next_plan_refill_at ON backend.billing_users (next_plan_refill_at);
	`

	CreateTableBillingTransactionsHistorySQL = `
	CREATE TABLE IF NOT EXISTS backend.billing_transactions_history (
		order_id VARCHAR(255) NOT NULL,
		user_id UUID NOT NULL,
		event_kind VARCHAR(255) NOT NULL,
		amount_usd_cents INTEGER,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		additional_info JSONB NOT NULL DEFAULT '{}',
		CONSTRAINT fk_backend_billing_transactions_history_user_id FOREIGN KEY (user_id) REFERENCES backend.users(id),
		PRIMARY KEY (order_id, event_kind)
	);

	CREATE INDEX IF NOT EXISTS idx_backend_billing_transactions_history_user_id_kind ON backend.billing_transactions_history (user_id, event_kind);
	`

	CreateTableBillingOrdersSQL = `
	CREATE TABLE IF NOT EXISTS backend.billing_orders (
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL,
		status VARCHAR(255) NOT NULL,
		type VARCHAR(255) NOT NULL,
		chain VARCHAR(255) NOT NULL,
		token_address VARCHAR(255) NOT NULL,
		amount_usd_cents INTEGER NOT NULL,
		wallet_address VARCHAR(255) NOT NULL,
		promo_code VARCHAR(255),
		product_data JSONB NOT NULL DEFAULT '{}'::jsonb,
		expires_at TIMESTAMP NOT NULL,
		scan_from_block BIGINT NOT NULL DEFAULT 0,
		paid_tx_hash VARCHAR(255),
		paid_block_number BIGINT,
		payment_provider VARCHAR(32) NOT NULL DEFAULT 'wallet_transfer',
		payment_payload_hash VARCHAR(64),
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_backend_billing_orders_user_id FOREIGN KEY (user_id) REFERENCES backend.users(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_backend_billing_orders_user_id_created_at ON backend.billing_orders (user_id, created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_backend_billing_orders_wallet_transfer_expiring ON backend.billing_orders (expires_at) WHERE status = 'pending' AND payment_provider = 'wallet_transfer';
	CREATE INDEX IF NOT EXISTS idx_backend_billing_orders_wallet_pending_scan ON backend.billing_orders (chain, scan_from_block, created_at) WHERE status = 'pending' AND payment_provider = 'wallet_transfer';
	CREATE UNIQUE INDEX IF NOT EXISTS idx_backend_billing_orders_pending_wallet_address ON backend.billing_orders (wallet_address) WHERE status = 'pending' AND payment_provider = 'wallet_transfer';
	CREATE UNIQUE INDEX IF NOT EXISTS idx_backend_billing_orders_x402_paid_tx_hash ON backend.billing_orders (paid_tx_hash) WHERE payment_provider = 'x402' AND paid_tx_hash IS NOT NULL;
	CREATE UNIQUE INDEX IF NOT EXISTS idx_backend_billing_orders_x402_payment_payload_hash ON backend.billing_orders (payment_payload_hash) WHERE payment_provider = 'x402' AND payment_payload_hash IS NOT NULL;
	CREATE INDEX IF NOT EXISTS idx_backend_billing_orders_x402_payment_settled ON backend.billing_orders (updated_at) WHERE status = 'payment_settled' AND payment_provider = 'x402';
	`

	CreateTableBillingPromoCodesSQL = `
	CREATE TABLE IF NOT EXISTS backend.billing_promo_codes (
		code VARCHAR(255) PRIMARY KEY,
		active BOOLEAN NOT NULL DEFAULT TRUE,
		scopes TEXT[] NOT NULL DEFAULT '{}'::text[],
		discount_type VARCHAR(255) NOT NULL,
		discount_value INTEGER NOT NULL,
		max_uses INTEGER,
		uses_count INTEGER NOT NULL DEFAULT 0,
		target_plans TEXT[] NOT NULL DEFAULT '{}'::text[],
		min_months INTEGER,
		max_months INTEGER,
		min_credits BIGINT,
		max_credits BIGINT,
		expires_at TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	SelectBillingPlansSQL = `SELECT plan_requests_per_minute, plan_credits, plan_code, plan_features, plan_price, plan_name FROM backend.billing_plans;`

	SelectBillingUserStateSQL                = `SELECT user_id, COALESCE(api_keys, '{}'::uuid[])::text[], plan, plan_credits, additional_credits, billing_period_paid_cents, plan_started_at, next_plan_refill_at, plan_ending_at FROM backend.billing_users WHERE user_id = $1;`
	SelectBillingUserStateForUpdateSQL       = `SELECT user_id, COALESCE(api_keys, '{}'::uuid[])::text[], plan, plan_credits, additional_credits, billing_period_paid_cents, plan_started_at, next_plan_refill_at, plan_ending_at FROM backend.billing_users WHERE user_id = $1 FOR UPDATE;`
	SelectBillingUserStateForUpdateNowaitSQL = `SELECT user_id, COALESCE(api_keys, '{}'::uuid[])::text[], plan, plan_credits, additional_credits, billing_period_paid_cents, plan_started_at, next_plan_refill_at, plan_ending_at FROM backend.billing_users WHERE user_id = $1 FOR UPDATE NOWAIT;`
	SelectBillingUserPlanSQL                 = `SELECT plan FROM backend.billing_users WHERE user_id = $1;`
	SelectBillingUserPlanCreditsSQL          = `SELECT plan_credits FROM backend.billing_users WHERE user_id = $1;`
	SelectBillingUserAdditionalCreditsSQL    = `SELECT additional_credits FROM backend.billing_users WHERE user_id = $1;`
	UpdateBillingUserStateSQL                = `UPDATE backend.billing_users SET plan = $2, plan_credits = $3, additional_credits = $4, billing_period_paid_cents = $5, plan_started_at = $6, next_plan_refill_at = $7, plan_ending_at = $8 WHERE user_id = $1;`

	SelectBillingApiKeysByUserIDSQL = `SELECT api_keys FROM backend.billing_users WHERE user_id = $1;`
	SelectUserIDByAPIKeySQL         = `SELECT user_id FROM backend.billing_users WHERE COALESCE(api_keys, '{}'::uuid[]) @> ARRAY[$1::uuid];`

	InsertBillingAPIKeyForUserSQL  = `UPDATE backend.billing_users SET api_keys = array_append(COALESCE(api_keys, '{}'::uuid[]), $1) WHERE user_id = $2 AND NOT COALESCE(api_keys, '{}'::uuid[]) @> ARRAY[$1::uuid];`
	DeleteBillingAPIKeyFromUserSQL = `UPDATE backend.billing_users SET api_keys = array_remove(COALESCE(api_keys, '{}'::uuid[]), $1) WHERE user_id = $2 AND COALESCE(api_keys, '{}'::uuid[]) @> ARRAY[$1::uuid];`
	InsertBillingUserSQL           = `WITH plan_credits AS (SELECT plan_credits FROM backend.billing_plans WHERE plan_code = 'free') INSERT INTO backend.billing_users (user_id, plan_credits, plan, plan_started_at, next_plan_refill_at, plan_ending_at) VALUES ($1, (SELECT plan_credits FROM plan_credits), 'free', NOW(), date_trunc('day', NOW() AT TIME ZONE 'UTC') + interval '1 day', NULL) ON CONFLICT (user_id) DO NOTHING;`

	SelectBillingUsersForMonthlyRefillSQL = `
	SELECT user_id, COALESCE(api_keys, '{}'::uuid[])::text[], plan, plan_credits, additional_credits, billing_period_paid_cents, plan_started_at, next_plan_refill_at, plan_ending_at
	FROM backend.billing_users
	WHERE plan != 'free' AND next_plan_refill_at IS NOT NULL AND next_plan_refill_at <= NOW() AND plan_ending_at > NOW();
	`

	RefillBillingUsersDailyFreeBatchSQL = `
	WITH due_users AS MATERIALIZED (
		SELECT user_id
		FROM backend.billing_users
		WHERE plan = 'free'
			AND next_plan_refill_at IS NOT NULL
			AND next_plan_refill_at <= $3
		ORDER BY next_plan_refill_at, user_id
		LIMIT $4
		FOR UPDATE SKIP LOCKED
	), updated_users AS (
		UPDATE backend.billing_users AS billing_user
		SET plan_credits = $1, next_plan_refill_at = $2
		FROM due_users
		WHERE billing_user.user_id = due_users.user_id
		RETURNING billing_user.user_id
	), inserted_history AS (
		INSERT INTO backend.billing_transactions_history (order_id, user_id, event_kind, amount_usd_cents, created_at)
		SELECT updated_users.user_id::text || ':' || ((EXTRACT(EPOCH FROM $2::timestamp))::bigint)::text,
			updated_users.user_id, 'free_daily_refill', 0, $3
		FROM updated_users
		RETURNING user_id
	)
	SELECT updated_users.user_id::text
	FROM updated_users
	JOIN inserted_history USING (user_id);
	`

	SelectBillingUsersForExpirationSQL = `
	SELECT user_id, COALESCE(api_keys, '{}'::uuid[])::text[], plan, plan_credits, additional_credits, billing_period_paid_cents, plan_started_at, next_plan_refill_at, plan_ending_at
	FROM backend.billing_users
	WHERE plan != 'free' AND plan_ending_at IS NOT NULL AND plan_ending_at <= NOW();
	`

	CreateTempTableBillingUserPlanCreditsAndAdditionalCreditsBatchSQL = `CREATE TEMP TABLE temp_credits (
		user_id uuid,
		plan_credits BIGINT NOT NULL DEFAULT 0,
		additional_credits BIGINT NOT NULL DEFAULT 0
	) ON COMMIT DROP;`

	SetBillingUserPlanCreditsAndAdditionalCreditsBatchSQL = `
        UPDATE backend.billing_users u
        SET plan_credits = t.plan_credits, additional_credits = t.additional_credits
        FROM temp_credits t
        WHERE u.user_id = t.user_id;
    `

	SelectBillingTransactionHistorySQL                      = `SELECT order_id, user_id, event_kind, amount_usd_cents, created_at, additional_info FROM backend.billing_transactions_history WHERE user_id = $1 AND event_kind = $2 ORDER BY created_at DESC;`
	SelectBillingTransactionHistoryByOrderIDSQL             = `SELECT order_id, user_id, event_kind, amount_usd_cents, created_at, additional_info FROM backend.billing_transactions_history WHERE order_id = $1 ORDER BY created_at DESC, event_kind DESC LIMIT 1;`
	InsertBillingTransactionHistorySQL                      = `INSERT INTO backend.billing_transactions_history (order_id, user_id, event_kind, amount_usd_cents, created_at, additional_info) VALUES ($1, $2, $3, $4, $5, $6);`
	InsertBillingTransactionHistoryWithoutAdditionalInfoSQL = `INSERT INTO backend.billing_transactions_history (order_id, user_id, event_kind, amount_usd_cents, created_at) VALUES ($1, $2, $3, $4, $5);`

	SelectBillingOrdersByUserIDSQL                  = `SELECT id, user_id, status, type, chain, token_address, amount_usd_cents, wallet_address, promo_code, product_data, expires_at, scan_from_block, paid_tx_hash, paid_block_number, payment_provider, payment_payload_hash, created_at, updated_at FROM backend.billing_orders WHERE user_id = $1 ORDER BY created_at DESC;`
	SelectBillingOrderByIDSQL                       = `SELECT id, user_id, status, type, chain, token_address, amount_usd_cents, wallet_address, promo_code, product_data, expires_at, scan_from_block, paid_tx_hash, paid_block_number, payment_provider, payment_payload_hash, created_at, updated_at FROM backend.billing_orders WHERE id = $1;`
	SelectBillingOrderByPaidTxHashSQL               = `SELECT id, user_id, status, type, chain, token_address, amount_usd_cents, wallet_address, promo_code, product_data, expires_at, scan_from_block, paid_tx_hash, paid_block_number, payment_provider, payment_payload_hash, created_at, updated_at FROM backend.billing_orders WHERE paid_tx_hash = $1 AND payment_provider = $2;`
	SelectBillingOrderByPaymentPayloadHashSQL       = `SELECT id, user_id, status, type, chain, token_address, amount_usd_cents, wallet_address, promo_code, product_data, expires_at, scan_from_block, paid_tx_hash, paid_block_number, payment_provider, payment_payload_hash, created_at, updated_at FROM backend.billing_orders WHERE payment_payload_hash = $1 AND payment_provider = $2;`
	SelectBillingOrderByIDForUpdateSQL              = `SELECT id, user_id, status, type, chain, token_address, amount_usd_cents, wallet_address, promo_code, product_data, expires_at, scan_from_block, paid_tx_hash, paid_block_number, payment_provider, payment_payload_hash, created_at, updated_at FROM backend.billing_orders WHERE id = $1 FOR UPDATE;`
	SelectX402PaymentSettledBillingOrdersSQL        = `SELECT id, user_id, status, type, chain, token_address, amount_usd_cents, wallet_address, promo_code, product_data, expires_at, scan_from_block, paid_tx_hash, paid_block_number, payment_provider, payment_payload_hash, created_at, updated_at FROM backend.billing_orders WHERE status = 'payment_settled' AND payment_provider = 'x402' ORDER BY updated_at LIMIT 100;`
	SelectPendingBillingOrderCountByUserAndTypesSQL = `SELECT COUNT(*) FROM backend.billing_orders WHERE user_id = $1 AND status = 'pending' AND type = ANY($2);`
	SelectExpiredPendingBillingOrdersSQL            = `SELECT id, user_id, status, type, chain, token_address, amount_usd_cents, wallet_address, promo_code, product_data, expires_at, scan_from_block, paid_tx_hash, paid_block_number, payment_provider, payment_payload_hash, created_at, updated_at FROM backend.billing_orders WHERE status = 'pending' AND payment_provider = 'wallet_transfer' AND expires_at <= NOW() ORDER BY expires_at;`
	SelectPendingBillingOrderChainsSQL              = `SELECT DISTINCT chain FROM backend.billing_orders WHERE status = 'pending' AND payment_provider = 'wallet_transfer'`

	InsertBillingOrderSQL = `
	INSERT INTO backend.billing_orders (
		id, user_id, status, type, chain, token_address, amount_usd_cents, wallet_address, promo_code, product_data, expires_at, scan_from_block, paid_tx_hash, paid_block_number, payment_provider, payment_payload_hash, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18);
	`
	UpdateBillingOrderStatusSQL                = `UPDATE backend.billing_orders SET status = $2, paid_tx_hash = $3, paid_block_number = $4, updated_at = $5 WHERE id = $1;`
	UpdatePendingBillingOrdersScanFromBlockSQL = `UPDATE backend.billing_orders SET scan_from_block = $2, updated_at = $3 WHERE id = ANY($1) AND status = 'pending' AND scan_from_block < $2;`

	SelectPendingBillingOrderBatchWithSameScanFromBlockSQL = `
	SELECT id, user_id, status, type, chain, token_address, amount_usd_cents, wallet_address, promo_code, product_data, expires_at, scan_from_block, paid_tx_hash, paid_block_number, payment_provider, payment_payload_hash, created_at, updated_at
	FROM backend.billing_orders o
	WHERE o.status = 'pending'
		AND o.payment_provider = 'wallet_transfer'
		AND o.chain = $1
		AND o.scan_from_block < $2
		AND o.scan_from_block = (
			SELECT MIN(o2.scan_from_block)
			FROM backend.billing_orders o2
			WHERE o2.status = 'pending' AND o2.payment_provider = 'wallet_transfer' AND o2.chain = $1 AND o2.scan_from_block < $2
		)
	ORDER BY o.created_at
	LIMIT $3;
	`

	SelectFailedBillingOrdersSQL = `
	SELECT id, user_id, status, type, chain, token_address, amount_usd_cents, wallet_address, promo_code, product_data, expires_at, scan_from_block, paid_tx_hash, paid_block_number, payment_provider, payment_payload_hash, created_at, updated_at
	FROM backend.billing_orders
	WHERE status = 'fulfillment_failed' AND payment_provider = 'wallet_transfer';
	`

	SelectBillingPromoCodeByCodeSQL          = `SELECT code, active, scopes, discount_type, discount_value, max_uses, uses_count, target_plans, min_months, max_months, min_credits, max_credits, expires_at, created_at, updated_at FROM backend.billing_promo_codes WHERE code = $1;`
	SelectBillingPromoCodeByCodeForUpdateSQL = `SELECT code, active, scopes, discount_type, discount_value, max_uses, uses_count, target_plans, min_months, max_months, min_credits, max_credits, expires_at, created_at, updated_at FROM backend.billing_promo_codes WHERE code = $1 FOR UPDATE;`
	SelectBillingPromoCodesSQL               = `SELECT code, active, scopes, discount_type, discount_value, max_uses, uses_count, target_plans, min_months, max_months, min_credits, max_credits, expires_at, created_at, updated_at FROM backend.billing_promo_codes ORDER BY code;`
	InsertBillingPromoCodeSQL                = `INSERT INTO backend.billing_promo_codes (code, active, scopes, discount_type, discount_value, max_uses, uses_count, target_plans, min_months, max_months, min_credits, max_credits, expires_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);`
	UpdateBillingPromoCodeSQL                = `UPDATE backend.billing_promo_codes SET active = $2, scopes = $3, discount_type = $4, discount_value = $5, max_uses = $6, uses_count = $7, target_plans = $8, min_months = $9, max_months = $10, min_credits = $11, max_credits = $12, expires_at = $13, updated_at = $14 WHERE code = $1;`
)
