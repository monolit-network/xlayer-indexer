package queries

const (
	// Investors
	CreateInvestorsTableSQL = `
                CREATE TABLE IF NOT EXISTS investors (
                id SERIAL PRIMARY KEY,
                slug TEXT NOT NULL,
                name TEXT NOT NULL,
                image TEXT,
                rank INT,
                country TEXT,
                description TEXT,
                venture_type TEXT,
                tier TEXT,
                lead_investments INT NOT NULL DEFAULT 0,
                links JSONB NOT NULL DEFAULT '[]'::jsonb,
                portfolio_coins INT[] NOT NULL DEFAULT '{}',
                portfolio_coins_slugs TEXT[] NOT NULL DEFAULT '{}',
                round_distribution_by_category JSONB NOT NULL DEFAULT '[]'::jsonb,
                round_distribution_by_stage JSONB NOT NULL DEFAULT '[]'::jsonb
                );

                CREATE INDEX IF NOT EXISTS investors_slug_idx ON investors (slug);
                CREATE INDEX IF NOT EXISTS investors_portfolio_coins_gin_idx ON investors USING GIN (portfolio_coins);
                CREATE INDEX IF NOT EXISTS investors_portfolio_coins_slugs_gin_idx ON investors USING GIN (portfolio_coins_slugs);
                CREATE INDEX IF NOT EXISTS investors_links_gin_idx ON investors USING GIN (links jsonb_path_ops);
                CREATE INDEX IF NOT EXISTS investors_round_dist_category_gin_idx ON investors USING GIN (round_distribution_by_category jsonb_path_ops);
                CREATE INDEX IF NOT EXISTS investors_round_dist_stage_gin_idx ON investors USING GIN (round_distribution_by_stage jsonb_path_ops);`

	UpsertInvestorSQL = `
                INSERT INTO investors (
                    id, slug, name, image, rank, country, description, venture_type, tier, lead_investments,
                    links, portfolio_coins, portfolio_coins_slugs, round_distribution_by_category, round_distribution_by_stage
                )
                VALUES (
                    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
                    $11, $12, $13, $14, $15
                )
                ON CONFLICT (id)
                DO UPDATE SET
                    slug = EXCLUDED.slug,
                    name = EXCLUDED.name,
                    image = EXCLUDED.image,
                    rank = EXCLUDED.rank,
                    country = EXCLUDED.country,
                    description = EXCLUDED.description,
                    venture_type = EXCLUDED.venture_type,
                    tier = EXCLUDED.tier,
                    lead_investments = EXCLUDED.lead_investments,
                    links = EXCLUDED.links,
                    portfolio_coins = EXCLUDED.portfolio_coins,
                    portfolio_coins_slugs = EXCLUDED.portfolio_coins_slugs,
                    round_distribution_by_category = EXCLUDED.round_distribution_by_category,
                    round_distribution_by_stage = EXCLUDED.round_distribution_by_stage`

	SelectExistingInvestorSlugsSQL = `
                SELECT slug
                FROM investors
                WHERE slug = ANY($1)`

	SelectInvestorIdsBySlugsSQL = `
                SELECT id, slug
                FROM investors
                WHERE slug = ANY($1)`

	// Invest rounds
	CreateInvestRoundsTableSQL = `
                CREATE TABLE IF NOT EXISTS invest_rounds (
                id INT PRIMARY KEY,
                asset_id INT NOT NULL,
                funds_raised BIGINT,
                pre_valuation BIGINT,
                pre_valuation_inaccurate BOOLEAN NOT NULL DEFAULT FALSE,
                stage TEXT,
                category TEXT,
                date TIMESTAMP NOT NULL,
                investors INT[] NOT NULL DEFAULT '{}',
                lead_investors INT[] NOT NULL DEFAULT '{}',
                investors_slugs TEXT[] NOT NULL DEFAULT '{}',
                lead_investors_slugs TEXT[] NOT NULL DEFAULT '{}'
                );

                CREATE INDEX IF NOT EXISTS invest_rounds_asset_id_idx ON invest_rounds (asset_id);
                CREATE INDEX IF NOT EXISTS invest_rounds_date_idx ON invest_rounds (date DESC);
                CREATE INDEX IF NOT EXISTS invest_rounds_stage_idx ON invest_rounds (stage);
                CREATE INDEX IF NOT EXISTS invest_rounds_category_idx ON invest_rounds (category);
                CREATE INDEX IF NOT EXISTS invest_rounds_investors_gin_idx ON invest_rounds USING GIN (investors);
                CREATE INDEX IF NOT EXISTS invest_rounds_lead_investors_gin_idx ON invest_rounds USING GIN (lead_investors);
                CREATE INDEX IF NOT EXISTS invest_rounds_investors_slugs_gin_idx ON invest_rounds USING GIN (investors_slugs);
                CREATE INDEX IF NOT EXISTS invest_rounds_lead_investors_slugs_gin_idx ON invest_rounds USING GIN (lead_investors_slugs);`

	UpsertInvestRoundSQL = `
                INSERT INTO invest_rounds (
                    id, asset_id, funds_raised, pre_valuation, pre_valuation_inaccurate,
                    stage, category, date, investors, lead_investors, investors_slugs, lead_investors_slugs
                )
                VALUES (
                    $1, $2, $3, $4, $5,
                    $6, $7, $8, $9, $10, $11, $12
                )
                ON CONFLICT (id)
                DO UPDATE SET
                    asset_id = EXCLUDED.asset_id,
                    funds_raised = EXCLUDED.funds_raised,
                    pre_valuation = EXCLUDED.pre_valuation,
                    pre_valuation_inaccurate = EXCLUDED.pre_valuation_inaccurate,
                    stage = EXCLUDED.stage,
                    category = EXCLUDED.category,
                    date = EXCLUDED.date,
                    investors = EXCLUDED.investors,
                    lead_investors = EXCLUDED.lead_investors,
                    investors_slugs = EXCLUDED.investors_slugs,
                    lead_investors_slugs = EXCLUDED.lead_investors_slugs`

	SelectExistingInvestRoundIDsSQL = `
                SELECT id
                FROM invest_rounds
                WHERE id = ANY($1) AND asset_id > 0`

	SelectNextInvestorIDSQL = `
                SELECT nextval(pg_get_serial_sequence('investors','id'))`
)
