package queries

const (
	CreateSchemaPolymarketSQL = `
		CREATE SCHEMA IF NOT EXISTS polymarket;
	`

	CreatePolymarketMarketsNewTableSQL = `
		CREATE TABLE IF NOT EXISTS polymarket.markets (
			-- general info
			condition_id TEXT PRIMARY KEY,
			question_id TEXT NOT NULL,
			oracle TEXT NOT NULL,
			prepared_at TIMESTAMP NOT NULL,
			prepared_in_block BIGINT NOT NULL,
			prepared_in_block_hash TEXT NOT NULL,

			-- priority: umactf adapter contract -> gamma
			anc_question TEXT,
			anc_description TEXT,
			anc_res_data TEXT,
			anc_market_id INT,
			anc_initializer TEXT,

			-- uma info
			uma_current_status TEXT,

			-- technical info
			ancillary_data BYTEA,

			-- outcomes info
			total_outcomes INT NOT NULL,

			-- resolution info
			is_resolved BOOLEAN NOT NULL,
			resolved_at TIMESTAMP,
			resolved_in_block BIGINT,
			resolved_in_block_hash TEXT,
			payout_numerators NUMERIC(78, 0)[] NOT NULL, 

			-- info from gamma api
			is_presented_in_gamma BOOLEAN,
			gamma_id INT,
			gamma_question TEXT,
			gamma_description TEXT,
			gamma_slug TEXT,
			gamma_event_slug TEXT,
			gamma_event_link TEXT,
			gamma_resolution_source TEXT,
			gamma_start_date TIMESTAMP,
			gamma_end_date TIMESTAMP,
			gamma_event_image_url TEXT,
			gamma_event_icon_url TEXT,
			gamma_market_image_url TEXT,
			gamma_market_icon_url TEXT,
			gamma_outcomes TEXT[],
			gamma_created_at TIMESTAMP,
			gamma_updated_at TIMESTAMP,
			gamma_closed_at TIMESTAMP,
			gamma_group_item_title TEXT,
			
			gamma_order_price_min_tick_size NUMERIC(78, 18), 
			gamma_order_min_size NUMERIC(78, 18),
			gamma_accepting_orders BOOLEAN,

			gamma_tag_slugs TEXT[],

			gamma_neg_risk BOOLEAN,
			gamma_neg_risk_request_id TEXT,
			gamma_neg_risk_other BOOLEAN,

			gamma_uma_resolution_status TEXT,
			gamma_uma_resolution_statuses TEXT[],
			gamma_uma_bond NUMERIC(78, 18),
			gamma_uma_reward NUMERIC(78, 18),

			gamma_fees_enabled BOOLEAN,
			gamma_fee NUMERIC(78, 18),

			gamma_raw_response BYTEA,

			-- system fields
			system_created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			system_updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE OR REPLACE FUNCTION update_system_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.system_updated_at = NOW();
			RETURN NEW;
		END;
		$$ language 'plpgsql';

		DROP TRIGGER IF EXISTS update_markets_modtime ON polymarket.markets;
		CREATE TRIGGER update_markets_modtime
			BEFORE UPDATE ON polymarket.markets
			FOR EACH ROW
			EXECUTE PROCEDURE update_system_updated_at_column();

		CREATE INDEX IF NOT EXISTS idx_polymarket_markets_condition_id ON polymarket.markets (condition_id);
		CREATE INDEX IF NOT EXISTS idx_polymarket_markets_question_id ON polymarket.markets (question_id);
		CREATE INDEX IF NOT EXISTS idx_polymarket_markets_resolved_in_block ON polymarket.markets (resolved_in_block);

		CREATE OR REPLACE FUNCTION polymarket.sync_market_tokens_outcome_name()
		RETURNS TRIGGER AS $$
		BEGIN
			IF to_regclass('polymarket.tokens') IS NULL THEN
				RETURN NEW;
			END IF;

			IF TG_OP = 'UPDATE' AND NEW.gamma_outcomes IS NOT DISTINCT FROM OLD.gamma_outcomes THEN
				RETURN NEW;
			END IF;

			UPDATE polymarket.tokens t
			SET outcome_name = polymarket.compute_token_outcome_name(t.partition, NEW.gamma_outcomes)
			WHERE t.condition_id = NEW.condition_id
			  AND t.outcome_name IS DISTINCT FROM polymarket.compute_token_outcome_name(t.partition, NEW.gamma_outcomes);

			RETURN NEW;
		END;
		$$ language 'plpgsql';

		DROP TRIGGER IF EXISTS trg_polymarket_markets_sync_token_outcome_name ON polymarket.markets;
		CREATE TRIGGER trg_polymarket_markets_sync_token_outcome_name
			AFTER INSERT OR UPDATE OF gamma_outcomes ON polymarket.markets
			FOR EACH ROW
			EXECUTE PROCEDURE polymarket.sync_market_tokens_outcome_name();
	`

	CreatePolymarketTokensTableSQL = `
		CREATE TABLE IF NOT EXISTS polymarket.tokens (
			token_id NUMERIC(78, 0) PRIMARY KEY,
			condition_id TEXT NOT NULL,
			collateral_token TEXT NOT NULL,
			parent_collection_id TEXT NOT NULL,
			partition NUMERIC(78, 0) NOT NULL,
			created_at TIMESTAMP NOT NULL,

			is_resolved BOOLEAN NOT NULL,
			resolved_at TIMESTAMP,
			resolved_in_block BIGINT,
			resolved_in_block_hash TEXT,
			numerator NUMERIC(78, 0),
			denominator NUMERIC(78, 0),
			payout NUMERIC(38, 18) GENERATED ALWAYS AS (
				CASE
					WHEN denominator IS NULL OR denominator = 0 THEN 0
					ELSE ( (numerator * 1.000000000000000000) / denominator )::NUMERIC(38, 18)
				END
			) STORED,
			outcome_name TEXT
		);
		ALTER TABLE polymarket.tokens ADD COLUMN IF NOT EXISTS outcome_name TEXT;
		CREATE INDEX IF NOT EXISTS idx_polymarket_tokens_condition_id ON polymarket.tokens (condition_id);

		CREATE OR REPLACE FUNCTION polymarket.compute_token_outcome_name(partition_value NUMERIC, outcomes TEXT[])
		RETURNS TEXT
		LANGUAGE SQL
		IMMUTABLE
		RETURNS NULL ON NULL INPUT
		AS $$
			SELECT CASE
				WHEN array_length(outcomes, 1) IS NULL THEN NULL
				ELSE (
					SELECT string_agg(outcomes[s.idx + 1], ' OR ' ORDER BY s.idx)
					FROM generate_series(0, array_length(outcomes, 1) - 1) AS s(idx)
					WHERE floor((partition_value / power(2::numeric, s.idx::numeric))) % 2::numeric = 1::numeric
				)
			END
		$$;

		CREATE OR REPLACE FUNCTION polymarket.sync_token_outcome_name()
		RETURNS TRIGGER AS $$
		BEGIN
			IF to_regclass('polymarket.markets') IS NULL THEN
				NEW.outcome_name := NULL;
				RETURN NEW;
			END IF;

			SELECT polymarket.compute_token_outcome_name(NEW.partition, m.gamma_outcomes)
			INTO NEW.outcome_name
			FROM polymarket.markets m
			WHERE m.condition_id = NEW.condition_id;

			IF NOT FOUND THEN
				NEW.outcome_name := NULL;
			END IF;

			RETURN NEW;
		END;
		$$ language 'plpgsql';

		DROP TRIGGER IF EXISTS trg_polymarket_tokens_sync_outcome_name ON polymarket.tokens;
		CREATE TRIGGER trg_polymarket_tokens_sync_outcome_name
			BEFORE INSERT OR UPDATE OF condition_id, partition ON polymarket.tokens
			FOR EACH ROW
			EXECUTE PROCEDURE polymarket.sync_token_outcome_name();
	`

	InsertPolymarketMarketNewSQL = `
		INSERT INTO polymarket.markets (condition_id, question_id, oracle, prepared_at, prepared_in_block, prepared_in_block_hash, prepared_in_tx_hash, total_outcomes, is_resolved, payout_numerators)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, false, $9)
		ON CONFLICT (condition_id) DO UPDATE SET
			question_id = EXCLUDED.question_id,
			oracle = EXCLUDED.oracle,
			prepared_at = EXCLUDED.prepared_at,
			prepared_in_block = EXCLUDED.prepared_in_block,
			prepared_in_block_hash = EXCLUDED.prepared_in_block_hash,
			prepared_in_tx_hash = EXCLUDED.prepared_in_tx_hash,
			total_outcomes = EXCLUDED.total_outcomes,
			is_resolved = EXCLUDED.is_resolved,
			payout_numerators = EXCLUDED.payout_numerators
	`

	UpdatePolymarketMarketResolutionSQL = `
		UPDATE polymarket.markets
		SET is_resolved = true, resolved_at = $2, payout_numerators = $3, resolved_in_block = $4, resolved_in_block_hash = LOWER($5)
		WHERE condition_id = LOWER($1)
	`

	UpdatePolymarketMarketResolutionFullSQL = `
		WITH input_params AS (
			SELECT
				$1 AS raw_cond_id,
				LOWER($1) AS lower_cond_id,
				$2::TIMESTAMP AS res_at,
				$3::NUMERIC(78,0)[] AS payouts,
				$4::BIGINT AS block_num,
				LOWER($5) AS block_hash
		),
		payout_stats AS (
			SELECT
				COALESCE(SUM(val), 0) as total_sum
			FROM input_params, unnest(payouts) as val
		),
		unnested_payouts AS (
			SELECT n, idx
			FROM input_params, unnest(payouts) WITH ORDINALITY AS arr(n, idx)
		),
		update_markets AS (
			UPDATE polymarket.markets m
			SET
				is_resolved = true,
				resolved_at = i.res_at,
				payout_numerators = i.payouts,
				resolved_in_block = i.block_num,
				resolved_in_block_hash = i.block_hash
			FROM input_params i
			WHERE m.condition_id = i.lower_cond_id
		)
		UPDATE polymarket.tokens t
		SET
			is_resolved = true,
			resolved_at = i.res_at,
			resolved_in_block = i.block_num,
			resolved_in_block_hash = i.block_hash,
			denominator = ps.total_sum,
			numerator = (
				SELECT COALESCE(SUM(up.n), 0)
				FROM unnested_payouts up
				WHERE MOD(DIV(t.partition, POWER(2::numeric, up.idx - 1)), 2) = 1
			)
		FROM input_params i, payout_stats ps
		WHERE t.condition_id = i.lower_cond_id
		RETURNING t.token_id, t.condition_id, t.collateral_token, t.parent_collection_id, t.partition, t.numerator, t.denominator, t.is_resolved, t.resolved_in_block, t.resolved_in_block_hash
	`

	InsertPolymarketTokenSQL = `
		WITH input_tokens (token_id, condition_id, collateral_token, parent_collection_id, partition, created_at) AS (
			VALUES ($1::NUMERIC(78,0), $2::TEXT, $3::TEXT, $4::TEXT, $5::NUMERIC(78,0), $6::TIMESTAMP)
		)
		INSERT INTO polymarket.tokens AS existing (
			token_id,
			condition_id,
			collateral_token,
			parent_collection_id,
			partition,
			created_at,
			is_resolved,
			resolved_at,
			resolved_in_block,
			resolved_in_block_hash,
			numerator,
			denominator
		)
		SELECT
			t.token_id,
			LOWER(t.condition_id),
			t.collateral_token,
			t.parent_collection_id,
			t.partition,
			t.created_at,
			COALESCE(m.is_resolved, false),
			CASE WHEN m.is_resolved THEN m.resolved_at ELSE NULL END,
			CASE WHEN m.is_resolved THEN m.resolved_in_block ELSE NULL END,
			CASE WHEN m.is_resolved THEN m.resolved_in_block_hash ELSE NULL END,
			CASE WHEN m.is_resolved THEN (
				SELECT COALESCE(SUM(up.n), 0)
				FROM unnest(m.payout_numerators) WITH ORDINALITY AS up(n, idx)
				WHERE MOD(DIV(t.partition, POWER(2::numeric, up.idx - 1)), 2) = 1
			) ELSE NULL END,
			CASE WHEN m.is_resolved THEN (
				SELECT COALESCE(SUM(n), 0)
				FROM unnest(m.payout_numerators) AS n
			) ELSE NULL END
		FROM input_tokens t
		LEFT JOIN polymarket.markets m ON m.condition_id = LOWER(t.condition_id)
		ON CONFLICT (token_id) DO UPDATE
		SET
			is_resolved = EXCLUDED.is_resolved,
			resolved_at = EXCLUDED.resolved_at,
			resolved_in_block = EXCLUDED.resolved_in_block,
			resolved_in_block_hash = EXCLUDED.resolved_in_block_hash,
			numerator = EXCLUDED.numerator,
			denominator = EXCLUDED.denominator
		WHERE COALESCE(existing.is_resolved, false) = false
			AND EXCLUDED.is_resolved = true
	`

	InsertPolymarketTokensBatchSQL = `
		WITH input_tokens (token_id, condition_id, collateral_token, parent_collection_id, partition, created_at) AS (
			VALUES %s
		)
		INSERT INTO polymarket.tokens AS existing (
			token_id,
			condition_id,
			collateral_token,
			parent_collection_id,
			partition,
			created_at,
			is_resolved,
			resolved_at,
			resolved_in_block,
			resolved_in_block_hash,
			numerator,
			denominator
		)
		SELECT
			t.token_id,
			LOWER(t.condition_id),
			t.collateral_token,
			t.parent_collection_id,
			t.partition,
			t.created_at,
			COALESCE(m.is_resolved, false),
			CASE WHEN m.is_resolved THEN m.resolved_at ELSE NULL END,
			CASE WHEN m.is_resolved THEN m.resolved_in_block ELSE NULL END,
			CASE WHEN m.is_resolved THEN m.resolved_in_block_hash ELSE NULL END,
			CASE WHEN m.is_resolved THEN (
				SELECT COALESCE(SUM(up.n), 0)
				FROM unnest(m.payout_numerators) WITH ORDINALITY AS up(n, idx)
				WHERE MOD(DIV(t.partition, POWER(2::numeric, up.idx - 1)), 2) = 1
			) ELSE NULL END,
			CASE WHEN m.is_resolved THEN (
				SELECT COALESCE(SUM(n), 0)
				FROM unnest(m.payout_numerators) AS n
			) ELSE NULL END
		FROM input_tokens t
		LEFT JOIN polymarket.markets m ON m.condition_id = LOWER(t.condition_id)
		ON CONFLICT (token_id) DO UPDATE
		SET
			is_resolved = EXCLUDED.is_resolved,
			resolved_at = EXCLUDED.resolved_at,
			resolved_in_block = EXCLUDED.resolved_in_block,
			resolved_in_block_hash = EXCLUDED.resolved_in_block_hash,
			numerator = EXCLUDED.numerator,
			denominator = EXCLUDED.denominator
		WHERE COALESCE(existing.is_resolved, false) = false
			AND EXCLUDED.is_resolved = true
	`

	GetPolymarketTokensByConditionIDSQL = `
		SELECT token_id, condition_id, collateral_token, parent_collection_id, partition
		FROM polymarket.tokens
		WHERE condition_id = LOWER($1)
	`

	GetPolymarketTokenByTokenIDSQL = `
		SELECT token_id, condition_id, collateral_token, parent_collection_id, partition, numerator, denominator, is_resolved, resolved_in_block, resolved_in_block_hash
		FROM polymarket.tokens
		WHERE token_id = $1
	`

	GetPolymarketMarketByConditionIDSQL = `
		SELECT condition_id, question_id, anc_question, anc_description, ancillary_data, oracle, total_outcomes, is_resolved, resolved_at, payout_numerators
		FROM polymarket.markets
		WHERE condition_id = LOWER($1)
	`

	GetPolymarketMarketMinimalByConditionIDSQL = `
		SELECT condition_id, question_id
		FROM polymarket.markets
		WHERE condition_id = LOWER($1)
	`

	GetPolymarketQuestionTimestampFromEventsSQL = `
		SELECT event_data, event_type
		FROM polymarket.market_events
		WHERE question_id = LOWER($1)
		  AND event_type IN ('ADAPTER_INITIALIZED', 'ADAPTER_RESET')
		  AND block_number <= $2
		ORDER BY block_number DESC, tx_index DESC, log_index DESC
		LIMIT 1
	`

	UpdatePolymarketMarketAncillaryDataByQuestionIDSQL = `
		UPDATE polymarket.markets
		SET ancillary_data = $2,
			anc_question = $3,
			anc_description = $4,
			anc_res_data = $5,
			anc_market_id = $6,
			anc_initializer = $7
		WHERE question_id = LOWER($1)
	`

	UpdatePolymarketMarketParsedAncillaryDataByQuestionIDSQL = `
		UPDATE polymarket.markets
		SET
			anc_question = $2,
			anc_description = $3,
			anc_res_data = $4,
			anc_market_id = $5,
			anc_initializer = $6
		WHERE question_id = LOWER($1)
	`

	// UpdatePolymarketTokensResolutionSQL updates all tokens for a condition with resolution data
	// $1 = condition_id, $2 = resolved_at, $3 = payout_numerators array, $4 = resolved_in_block, $5 = resolved_in_block_hash
	// numerator is calculated based on partition bitmask and payout_numerators
	// denominator is sum of all payout_numerators
	UpdatePolymarketTokensResolutionSQL = `
		UPDATE polymarket.tokens t
		SET
			is_resolved = true,
			resolved_at = $2,
			resolved_in_block = $4,
			resolved_in_block_hash = LOWER($5),
			numerator = (
				SELECT COALESCE(SUM(n), 0)
				FROM unnest($3::NUMERIC(78,0)[]) WITH ORDINALITY AS arr(n, idx)
				WHERE MOD(DIV(t.partition, POWER(2::numeric, idx - 1)), 2) = 1
			),
			denominator = (
				SELECT COALESCE(SUM(n), 0)
				FROM unnest($3::NUMERIC(78,0)[]) AS n
			)
		WHERE condition_id = $1
	`

	CreatePolymarketMarketEventsTableSQL = `
		CREATE TABLE IF NOT EXISTS polymarket.market_events (
			question_id TEXT NOT NULL,
			event_type TEXT NOT NULL,
			block_number BIGINT NOT NULL,
			block_hash TEXT NOT NULL,
			log_index BIGINT NOT NULL,
			tx_index BIGINT NOT NULL,
			tx_hash TEXT NOT NULL,
			block_timestamp TIMESTAMP NOT NULL,
			event_data JSONB,
			PRIMARY KEY(block_number, tx_index, log_index)
		);
		CREATE INDEX IF NOT EXISTS idx_polymarket_market_events_question_id ON polymarket.market_events (question_id);
		CREATE INDEX IF NOT EXISTS idx_polymarket_market_events_event_type ON polymarket.market_events (event_type);
	`

	InsertPolymarketMarketEventSQL = `
		INSERT INTO polymarket.market_events (question_id, event_type, block_number, block_hash, log_index, tx_index, tx_hash, block_timestamp, event_data)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (block_number, tx_index, log_index) DO UPDATE SET
			question_id = EXCLUDED.question_id,
			event_type = EXCLUDED.event_type,
			block_hash = EXCLUDED.block_hash,
			tx_hash = EXCLUDED.tx_hash,
			block_timestamp = EXCLUDED.block_timestamp,
			event_data = EXCLUDED.event_data
	`

	GetPolymarketMarketEventsByQuestionIDSQL = `
		SELECT question_id, event_type, block_number, block_hash, log_index, tx_index, tx_hash, block_timestamp, event_data
		FROM polymarket.market_events
		WHERE question_id = $1
		ORDER BY block_number, tx_index, log_index
	`

	UpdatePolymarketMarketGammaDataSQL = `
		UPDATE polymarket.markets
		SET
			is_presented_in_gamma = $2,
			gamma_id = $3,
			gamma_question = $4,
			gamma_description = $5,
			gamma_slug = $6,
			gamma_event_slug = $7,
			gamma_event_link = $8,
			gamma_resolution_source = $9,
			gamma_start_date = $10,
			gamma_end_date = $11,
			gamma_created_at = $12,
			gamma_updated_at = $13,
			gamma_closed_at = $14,
			gamma_event_image_url = $15,
			gamma_event_icon_url = $16,
			gamma_market_image_url = $17,
			gamma_market_icon_url = $18,
			gamma_outcomes = $19,
			gamma_group_item_title = $20,
			gamma_accepting_orders = $21,
			gamma_order_price_min_tick_size = $22,
			gamma_order_min_size = $23,
			gamma_tag_slugs = $24,
			gamma_neg_risk = $25,
			gamma_neg_risk_request_id = $26,
			gamma_neg_risk_other = $27,
			gamma_uma_resolution_status = $28,
			gamma_uma_resolution_statuses = $29,
			gamma_uma_bond = $30,
			gamma_uma_reward = $31,
			gamma_fees_enabled = $32,
			gamma_fee = $33
			-- gamma_raw_response = $34
		WHERE condition_id = $1
			AND (gamma_updated_at IS NULL OR gamma_updated_at < $13)
	`

	// Reorg handling: delete/reset by block_hash

	DeletePolymarketMarketsByBlockHashSQL = `
		DELETE FROM polymarket.markets
		WHERE prepared_in_block_hash = LOWER($1)
	`

	DeletePolymarketMarketEventsByBlockHashSQL = `
		DELETE FROM polymarket.market_events
		WHERE block_hash = LOWER($1)
	`

	ResetPolymarketTokensResolutionByBlockHashSQL = `
		UPDATE polymarket.tokens
		SET
			is_resolved = false,
			resolved_at = NULL,
			resolved_in_block = NULL,
			resolved_in_block_hash = NULL,
			numerator = NULL,
			denominator = NULL
		WHERE resolved_in_block_hash = LOWER($1)
	`

	ResetPolymarketMarketsResolutionByBlockHashSQL = `
		UPDATE polymarket.markets
		SET
			is_resolved = false,
			resolved_at = NULL,
			resolved_in_block = NULL,
			resolved_in_block_hash = NULL,
			payout_numerators = '{}'
		WHERE resolved_in_block_hash = LOWER($1)
	`

	// Reorg handling: delete/reset by block_number (>= $1)

	DeletePolymarketMarketsFromBlockNumberSQL = `
		DELETE FROM polymarket.markets
		WHERE prepared_in_block >= $1
	`

	DeletePolymarketMarketEventsFromBlockNumberSQL = `
		DELETE FROM polymarket.market_events
		WHERE block_number >= $1
	`

	ResetPolymarketTokensResolutionFromBlockNumberSQL = `
		UPDATE polymarket.tokens
		SET
			is_resolved = false,
			resolved_at = NULL,
			resolved_in_block = NULL,
			resolved_in_block_hash = NULL,
			numerator = NULL,
			denominator = NULL
		WHERE resolved_in_block >= $1
	`

	ResetPolymarketMarketsResolutionFromBlockNumberSQL = `
		UPDATE polymarket.markets
		SET
			is_resolved = false,
			resolved_at = NULL,
			resolved_in_block = NULL,
			resolved_in_block_hash = NULL,
			payout_numerators = '{}'
		WHERE resolved_in_block >= $1
	`
)
