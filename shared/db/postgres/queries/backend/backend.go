package queries

const (
	CreateBackendDatabaseSQL = `CREATE SCHEMA IF NOT EXISTS backend`

	CreateTableBackendUsersSQL = `CREATE TABLE IF NOT EXISTS backend.users (
		id uuid PRIMARY KEY,
		thirdweb_user_id VARCHAR(255) NOT NULL UNIQUE,
		wallet_address VARCHAR(64),
		source VARCHAR(32) NOT NULL DEFAULT 'third_web',
		profiles BYTEA NOT NULL,
		referral_code text,
		referrer_code text,
		referral_wallet text,
		profile_extra JSONB NOT NULL DEFAULT '{"schema_version":1,"wallets":[],"notes":[],"preferences":{}}'::jsonb,
		created_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
		CONSTRAINT backend_users_referral_code_format CHECK (referral_code IS NULL OR referral_code ~ '^[A-Za-z0-9.-]{5,20}$'),
		CONSTRAINT backend_users_referrer_code_format CHECK (referrer_code IS NULL OR referrer_code ~ '^[A-Za-z0-9.-]{5,20}$')
	);

	CREATE UNIQUE INDEX IF NOT EXISTS idx_backend_users_wallet_address ON backend.users (wallet_address) WHERE wallet_address IS NOT NULL;
	CREATE INDEX IF NOT EXISTS idx_backend_users_thirdweb_user_id ON backend.users (thirdweb_user_id);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_backend_users_referral_code ON backend.users (referral_code) WHERE referral_code IS NOT NULL;
	CREATE INDEX IF NOT EXISTS idx_backend_users_referrer_code ON backend.users (referrer_code);
	CREATE INDEX IF NOT EXISTS idx_backend_users_profile_wallets ON backend.users USING GIN ((profile_extra -> 'wallets'));`

	CreateTableBackendReferralRewardsSQL = `
	CREATE TABLE IF NOT EXISTS backend.referral_rewards (
		id uuid PRIMARY KEY,
		referrer_user_id uuid NOT NULL,
		referred_user_id uuid NOT NULL,
		billing_order_id uuid NOT NULL,
		amount_usd_cents bigint NOT NULL,
		status VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		paid_out_at TIMESTAMP,
		CONSTRAINT fk_backend_referral_rewards_referrer_user_id FOREIGN KEY (referrer_user_id) REFERENCES backend.users(id) ON DELETE CASCADE,
		CONSTRAINT fk_backend_referral_rewards_referred_user_id FOREIGN KEY (referred_user_id) REFERENCES backend.users(id) ON DELETE CASCADE,
		CONSTRAINT fk_backend_referral_rewards_billing_order_id FOREIGN KEY (billing_order_id) REFERENCES backend.billing_orders(id) ON DELETE CASCADE,
		CONSTRAINT backend_referral_rewards_amount_nonnegative CHECK (amount_usd_cents >= 0),
		CONSTRAINT backend_referral_rewards_status CHECK (status IN ('available', 'paid_out'))
	);

	CREATE UNIQUE INDEX IF NOT EXISTS idx_backend_referral_rewards_billing_order_id ON backend.referral_rewards (billing_order_id);
	CREATE INDEX IF NOT EXISTS idx_backend_referral_rewards_referrer_status ON backend.referral_rewards (referrer_user_id, status);
	CREATE INDEX IF NOT EXISTS idx_backend_referral_rewards_referred_user_id ON backend.referral_rewards (referred_user_id);`

	InsertBackendUserWithWalletSQL = `INSERT INTO backend.users (id, thirdweb_user_id, profiles, wallet_address, source) VALUES ($1, $3::text, $2, lower($3::text), 'x402') ON CONFLICT (wallet_address) WHERE wallet_address IS NOT NULL DO NOTHING;`

	SelectBackendUserIDByWalletSQL = `SELECT id FROM backend.users WHERE wallet_address = lower($1);`

	SelectBackendUserByThirdwebUserIDSQL = `SELECT id, thirdweb_user_id, source, profiles, referral_code, referrer_code, referral_wallet, profile_extra FROM backend.users WHERE thirdweb_user_id = $1;`

	SelectBackendUserByIDSQL = `SELECT id, thirdweb_user_id, source, profiles, referral_code, referrer_code, referral_wallet, profile_extra FROM backend.users WHERE id = $1;`

	SelectBackendUserByReferralCodeSQL = `SELECT id, thirdweb_user_id, source, profiles, referral_code, referrer_code, referral_wallet, profile_extra FROM backend.users WHERE referral_code = $1;`

	UpdateBackendUserSQL = `UPDATE backend.users SET thirdweb_user_id = $1::text, wallet_address = lower($1::text), profiles = $2 WHERE id = $3;`

	SetBackendUserReferralCodeSQL = `UPDATE backend.users SET referral_code = $1 WHERE id = $2 AND referral_code IS NULL;`

	SetBackendUserReferrerCodeSQL = `UPDATE backend.users SET referrer_code = $1 WHERE id = $2 AND referrer_code IS NULL;`

	SetBackendUserReferralWalletSQL = `UPDATE backend.users SET referral_wallet = $1 WHERE id = $2;`

	InsertBackendReferralRewardForPaidOrderSQL = `
		INSERT INTO backend.referral_rewards (
			id,
			referrer_user_id,
			referred_user_id,
			billing_order_id,
			amount_usd_cents,
			status,
			created_at,
			updated_at
		)
		SELECT
			$1::uuid,
			referrer.id,
			referred.id,
			$2::uuid,
			(($3::bigint * $4::bigint) / 10000)::bigint,
			$5,
			$6,
			$6
		FROM backend.users referred
		JOIN backend.users referrer
			ON referred.referrer_code IS NOT NULL
			AND referrer.referral_code = referred.referrer_code
		WHERE referred.id = $7::uuid
		  AND referrer.id <> referred.id
		  AND $3::bigint > 0
		  AND (($3::bigint * $4::bigint) / 10000)::bigint > 0
		ON CONFLICT (billing_order_id) DO NOTHING;`

	SelectBackendUserReferralInfoSQL = `
		SELECT
			u.referral_code,
			u.referral_wallet,
			COALESCE(referrals.referral_count, 0)::bigint AS referral_count,
			COALESCE(rewards.total_usd_cents, 0)::bigint AS total_usd_cents,
			COALESCE(rewards.paid_out_usd_cents, 0)::bigint AS paid_out_usd_cents,
			COALESCE(rewards.available_usd_cents, 0)::bigint AS available_usd_cents
		FROM backend.users u
		LEFT JOIN LATERAL (
			SELECT COUNT(*)::bigint AS referral_count
			FROM backend.users referrals
			WHERE u.referral_code IS NOT NULL
			  AND referrals.referrer_code = u.referral_code
		) referrals ON TRUE
		LEFT JOIN LATERAL (
			SELECT
				SUM(amount_usd_cents)::bigint AS total_usd_cents,
				SUM(amount_usd_cents) FILTER (WHERE status = $2)::bigint AS paid_out_usd_cents,
				SUM(amount_usd_cents) FILTER (WHERE status = $3)::bigint AS available_usd_cents
			FROM backend.referral_rewards
			WHERE referrer_user_id = u.id
		) rewards ON TRUE
		WHERE u.id = $1
		LIMIT 1;`

	SelectBackendUserReferralsSQL = `
		WITH referred AS (
			SELECT
				referrals.id AS referral_id,
				LEFT(referrals.thirdweb_user_id, 7) AS wallet_prefix,
				COALESCE(SUM(rewards.amount_usd_cents), 0) AS referral_earnings_usd_cents,
				COUNT(*) OVER() AS total_count
			FROM backend.users owner
			JOIN backend.users referrals
				ON owner.referral_code IS NOT NULL
				AND referrals.referrer_code = owner.referral_code
			LEFT JOIN backend.referral_rewards rewards
				ON rewards.referrer_user_id = owner.id
				AND rewards.referred_user_id = referrals.id
			WHERE owner.id = $1
			GROUP BY referrals.id, referrals.thirdweb_user_id
			ORDER BY referrals.id
			LIMIT $2 OFFSET $3
		)
		SELECT
			wallet_prefix,
			referral_earnings_usd_cents::bigint,
			total_count
		FROM referred
		ORDER BY referral_id;`

	DeleteBackendUserByIDSQL = `DELETE FROM backend.users WHERE id = $1;`

	SelectProfileExtraByUserIDSQL = `SELECT profile_extra FROM backend.users WHERE id = $1;`

	SelectProfileExtraByUserIDForUpdateSQL = `SELECT profile_extra FROM backend.users WHERE id = $1 FOR UPDATE;`

	UpdateProfileExtraSQL = `UPDATE backend.users SET profile_extra = $2::jsonb WHERE id = $1;`

	CreateTableBackendSignalsSQL = `CREATE TABLE IF NOT EXISTS backend.signals (
		id uuid PRIMARY KEY,
		user_id uuid NOT NULL,
		kind VARCHAR(16) NOT NULL,
		status VARCHAR(16) NOT NULL DEFAULT 'active',
		nl_text TEXT NOT NULL,
		spec JSONB NOT NULL,
		human_readable TEXT NOT NULL,
		cadence_minutes INT NOT NULL,
		cooldown_minutes INT NOT NULL DEFAULT 360,
		last_value JSONB,
		last_evaluated_at TIMESTAMP,
		last_fired_at TIMESTAMP,
		next_eval_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		error_count INT NOT NULL DEFAULT 0,
		last_error TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_backend_signals_user_id FOREIGN KEY (user_id) REFERENCES backend.users(id) ON DELETE CASCADE,
		CONSTRAINT backend_signals_kind CHECK (kind IN ('alert','interest')),
		CONSTRAINT backend_signals_status CHECK (status IN ('active','paused','error')),
		CONSTRAINT backend_signals_cadence CHECK (cadence_minutes IN (5,15,60,360,720,1440))
	);

	CREATE INDEX IF NOT EXISTS idx_backend_signals_due ON backend.signals (status, next_eval_at);
	CREATE INDEX IF NOT EXISTS idx_backend_signals_user ON backend.signals (user_id);`

	CreateTableBackendSignalEventsSQL = `CREATE TABLE IF NOT EXISTS backend.signal_events (
		id uuid PRIMARY KEY,
		signal_id uuid NOT NULL,
		user_id uuid NOT NULL,
		fired_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		title TEXT NOT NULL,
		body TEXT NOT NULL,
		value_snapshot JSONB,
		seed_query TEXT,
		seen_at TIMESTAMP,
		channel VARCHAR(16) NOT NULL DEFAULT 'in_app',
		CONSTRAINT fk_backend_signal_events_signal_id FOREIGN KEY (signal_id) REFERENCES backend.signals(id) ON DELETE CASCADE,
		CONSTRAINT backend_signal_events_channel CHECK (channel IN ('in_app'))
	);

	CREATE INDEX IF NOT EXISTS idx_backend_signal_events_inbox ON backend.signal_events (user_id, fired_at DESC);`

	InsertSignalSQL = `INSERT INTO backend.signals
		(id, user_id, kind, status, nl_text, spec, human_readable, cadence_minutes, cooldown_minutes, next_eval_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$11);`

	// AcquireSignalCapLockSQL serializes concurrent signal creates for the same
	// (user_id, kind) pair only — the advisory xact lock is released
	// automatically at transaction end. Used by InsertSignalWithCap to make the
	// count+insert cap check atomic (TOCTOU-safe).
	AcquireSignalCapLockSQL = `SELECT pg_advisory_xact_lock(hashtextextended($1 || ':' || $2, 0));`

	// LeaseSignalSQL pushes next_eval_at forward as a claim lease. Unlike
	// UpdateSignalAfterEvalSQL it does NOT touch error_count or
	// last_evaluated_at — claiming is not an evaluation.
	LeaseSignalSQL = `UPDATE backend.signals SET next_eval_at = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1;`

	SelectSignalsByUserSQL = `SELECT id, user_id, kind, status, nl_text, spec, human_readable, cadence_minutes, cooldown_minutes,
		last_value, last_evaluated_at, last_fired_at, next_eval_at, error_count, last_error, created_at, updated_at
		FROM backend.signals WHERE user_id = $1 ORDER BY created_at DESC;`

	SelectSignalByIDSQL = `SELECT id, user_id, kind, status, nl_text, spec, human_readable, cadence_minutes, cooldown_minutes,
		last_value, last_evaluated_at, last_fired_at, next_eval_at, error_count, last_error, created_at, updated_at
		FROM backend.signals WHERE id = $1 AND user_id = $2;`

	UpdateSignalSQL = `UPDATE backend.signals SET
		status = COALESCE($3, status),
		cadence_minutes = COALESCE($4, cadence_minutes),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2;`

	DeleteSignalSQL = `DELETE FROM backend.signals WHERE id = $1 AND user_id = $2;`

	CountActiveSignalsSQL = `SELECT COUNT(*) FROM backend.signals WHERE user_id = $1 AND kind = $2 AND status = 'active';`

	SelectDueSignalsSQL = `SELECT id, user_id, kind, status, nl_text, spec, human_readable, cadence_minutes, cooldown_minutes,
		last_value, last_evaluated_at, last_fired_at, next_eval_at, error_count, last_error, created_at, updated_at
		FROM backend.signals
		WHERE status = 'active' AND next_eval_at <= CURRENT_TIMESTAMP
		ORDER BY next_eval_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED;`

	// last_value = COALESCE($3, last_value): a NULL new value (Python's quiet
	// missing-data tick returns fired=false, value=null) PRESERVES the stored
	// baseline instead of clearing it. crossed_above/crossed_below/
	// abs_change_gt/pct_change_gt compare against the persisted last_value, so
	// overwriting it to NULL on a data gap would re-arm (or never fire) those
	// signals. No caller has a legitimate need to clear last_value to NULL.
	UpdateSignalAfterEvalSQL = `UPDATE backend.signals SET
		last_evaluated_at = CURRENT_TIMESTAMP,
		last_fired_at = CASE WHEN $2 THEN CURRENT_TIMESTAMP ELSE last_fired_at END,
		last_value = COALESCE($3, last_value),
		next_eval_at = $4,
		error_count = CASE WHEN $5::text IS NULL THEN 0 ELSE error_count + 1 END,
		last_error = $5,
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status = 'active';`

	// status = 'active' guards (here and in UpdateSignalAfterEvalSQL /
	// InsertSignalEventSQL): with the lease pattern a user can pause or
	// delete the signal between claim and apply — a stale evaluation must
	// not stamp, escalate, or notify a signal that is no longer active.
	MarkSignalErrorSQL = `UPDATE backend.signals SET status = 'error', last_error = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND status = 'active';`

	// INSERT ... SELECT WHERE EXISTS: with the lease-pattern evaluator the
	// signal can be deleted between claim and apply; a plain VALUES insert
	// would violate the signal_id FK and roll back the whole apply tx.
	InsertSignalEventSQL = `INSERT INTO backend.signal_events
		(id, signal_id, user_id, fired_at, title, body, value_snapshot, seed_query, channel)
		SELECT $1,$2,$3,$4,$5,$6,$7,$8,$9
		WHERE EXISTS (SELECT 1 FROM backend.signals WHERE id = $2 AND status = 'active');`

	SelectSignalEventsSQL = `SELECT id, signal_id, user_id, fired_at, title, body, value_snapshot, seed_query, seen_at, channel
		FROM backend.signal_events
		WHERE user_id = $1 AND ($2::timestamp IS NULL OR fired_at > $2)
		ORDER BY fired_at DESC LIMIT $3;`

	MarkSignalEventSeenSQL = `UPDATE backend.signal_events SET seen_at = CURRENT_TIMESTAMP WHERE id = $1 AND user_id = $2 AND seen_at IS NULL;`

	PruneSignalEventsSQL = `DELETE FROM backend.signal_events WHERE user_id = $1 AND id NOT IN (
		SELECT id FROM backend.signal_events WHERE user_id = $1 ORDER BY fired_at DESC LIMIT $2);`

	ResetSignalErrorStateSQL = `UPDATE backend.signals SET error_count = 0, last_error = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND user_id = $2;`

	CountDueSignalsSQL = `SELECT COUNT(*) FROM backend.signals WHERE status = 'active' AND next_eval_at <= CURRENT_TIMESTAMP;`
)
