package queries

const (
	// Unverified tokens
	CreateUnverifiedTokensTableSQL = `
                CREATE TABLE IF NOT EXISTS unverified_tokens (
                    id BIGSERIAL PRIMARY KEY,
                    contract_address VARCHAR(48) NOT NULL,
                    chain VARCHAR(50) NOT NULL,
                    decimals INT,
                    symbol VARCHAR(20),
                    name VARCHAR(255),
                    price_usd DOUBLE PRECISION NOT NULL DEFAULT 0,
                    market_cap_usd DOUBLE PRECISION NOT NULL DEFAULT 0,
                    supply DOUBLE PRECISION NOT NULL DEFAULT 0,
                    largest_lp_pool_usd DOUBLE PRECISION NOT NULL DEFAULT 0,
                    first_tx_date TIMESTAMP NULL,
                    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                    view_source VARCHAR(100),
                    UNIQUE (contract_address, chain)
                );

                CREATE INDEX IF NOT EXISTS idx_unverified_contract_address ON unverified_tokens (contract_address);
                CREATE INDEX IF NOT EXISTS idx_unverified_chain ON unverified_tokens (chain);
                CREATE INDEX IF NOT EXISTS idx_unverified_contract_chain ON unverified_tokens (contract_address, chain);
                CREATE INDEX IF NOT EXISTS idx_unverified_contract_chain_decimals ON unverified_tokens (contract_address, chain, decimals, updated_at DESC) WHERE decimals IS NOT NULL;
                CREATE INDEX IF NOT EXISTS idx_unverified_market_cap ON unverified_tokens (market_cap_usd DESC);
                CREATE INDEX IF NOT EXISTS idx_unverified_price ON unverified_tokens (price_usd DESC);
                CREATE INDEX IF NOT EXISTS idx_unverified_updated_at ON unverified_tokens (updated_at DESC);`

	SelectUnverifiedTokenByChainContractAddressSQL = `
                SELECT chain, contract_address, symbol, name, decimals, price_usd, market_cap_usd, supply, largest_lp_pool_usd, first_tx_date, created_at, updated_at, view_source
                FROM unverified_tokens
                WHERE LOWER(chain) = LOWER($1) AND LOWER(contract_address) = LOWER($2)`

	SelectUnverifiedTokensBySymbolSQL = `
                SELECT chain, contract_address, symbol, name, decimals, price_usd, market_cap_usd, supply, largest_lp_pool_usd, first_tx_date, created_at, updated_at, view_source
                FROM unverified_tokens
                WHERE LOWER(symbol) = LOWER($1)`

	SelectUnverifiedTokensByContractAddressSQL = `
                SELECT chain, contract_address, symbol, name, decimals, price_usd, market_cap_usd, supply, largest_lp_pool_usd, first_tx_date, created_at, updated_at, view_source
                FROM unverified_tokens
                WHERE LOWER(contract_address) = LOWER($1)`

	SelectUnverifiedTokensByContractAddressesSQL = `
                SELECT chain, contract_address, symbol, name, decimals, price_usd, market_cap_usd, supply, largest_lp_pool_usd, first_tx_date, created_at, updated_at, view_source
                FROM unverified_tokens
                WHERE LOWER(contract_address) = ANY($1)`

	UpsertUnverifiedTokenSQL = `
                INSERT INTO unverified_tokens (
                    contract_address,
                    chain,
                    decimals,
                    symbol,
                    name,
                    price_usd,
                    market_cap_usd,
                    supply,
                    largest_lp_pool_usd,
                    first_tx_date,
                    created_at,
                    updated_at,
                    view_source
                )
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
                ON CONFLICT (contract_address, chain)
                DO UPDATE SET
                    decimals = EXCLUDED.decimals,
                    symbol = EXCLUDED.symbol,
                    name = EXCLUDED.name,
                    price_usd = EXCLUDED.price_usd,
                    market_cap_usd = EXCLUDED.market_cap_usd,
                    supply = EXCLUDED.supply,
                    largest_lp_pool_usd = EXCLUDED.largest_lp_pool_usd,
                    first_tx_date = EXCLUDED.first_tx_date,
                    updated_at = CURRENT_TIMESTAMP,
                    view_source = EXCLUDED.view_source`

	// Verified tokens
	CreateVerifiedTokensTableSQL = `
               CREATE TABLE IF NOT EXISTS verified_tokens (
                    asset_id INT NOT NULL,
                    chain VARCHAR(255) NOT NULL,
                    contract_address VARCHAR(255) NOT NULL,
                    symbol VARCHAR(255) NOT NULL,
                    name VARCHAR(255) NOT NULL,
                    decimals INT,
                    created_at TIMESTAMP,
                    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                    PRIMARY KEY (asset_id, chain),
                    CONSTRAINT fk_verified_tokens_asset FOREIGN KEY (asset_id) REFERENCES dropstab_info(asset_id)
                );

                CREATE INDEX IF NOT EXISTS verified_tokens_chain_contract_address_idx 
                    ON verified_tokens (chain, contract_address);

                CREATE INDEX IF NOT EXISTS verified_tokens_contract_address_idx 
                    ON verified_tokens (contract_address);`

	SelectVerifiedTokenByAssetIDChainSQL = `
                SELECT asset_id, chain, LOWER(contract_address), symbol, name, decimals, created_at, updated_at
                FROM verified_tokens
                WHERE asset_id = $1 AND chain != 'pls' AND LOWER(chain) = LOWER($2)`

	SelectVerifiedTokenByChainContractAddressSQL = `
                SELECT asset_id, chain, LOWER(contract_address), symbol, name, decimals, created_at, updated_at
                FROM verified_tokens
                WHERE chain != 'pls' AND LOWER(chain) = LOWER($1) AND LOWER(contract_address) = LOWER($2)`

	SelectVerifiedTokensByAssetIDSQL = `
                SELECT asset_id, chain, LOWER(contract_address), symbol, name, decimals, created_at, updated_at
                FROM verified_tokens
                WHERE asset_id = $1 AND chain != 'pls'`

	SelectVerifiedTokensByContractAddressSQL = `
                SELECT asset_id, chain, LOWER(contract_address), symbol, name, decimals, created_at, updated_at
                FROM verified_tokens
                WHERE chain != 'pls' AND LOWER(contract_address) = LOWER($1)`

	SelectVerifiedTokensBySymbolSQL = `
                SELECT asset_id, chain, LOWER(contract_address), symbol, name, decimals, created_at, updated_at
                FROM verified_tokens
                WHERE LOWER(symbol) = LOWER($1) AND chain != 'pls'`

	UpsertVerifiedTokenSQL = `
                INSERT INTO verified_tokens (asset_id, chain, contract_address, symbol, name, decimals, created_at, updated_at)
                VALUES ($1, LOWER($2), LOWER($3), $4, $5, $6, $7, $8)
                ON CONFLICT (asset_id, chain)
                DO UPDATE SET 
                contract_address = LOWER(EXCLUDED.contract_address),
                symbol = EXCLUDED.symbol,
                name = EXCLUDED.name,
                decimals = EXCLUDED.decimals,
                created_at = EXCLUDED.created_at,
                updated_at = EXCLUDED.updated_at
                `

	SelectNextDropstabInfoAssetIDSQL = `
                SELECT nextval('dropstab_info_asset_id_seq')`

	DeleteVerifiedTokenSQL = `
                DELETE FROM verified_tokens 
                WHERE asset_id = $1 AND LOWER(chain) = LOWER($2)`

	DeleteVerifiedTokensByAssetIDSQL = `
                DELETE FROM verified_tokens 
                WHERE asset_id = $1`

	// Dropstab
	CreateDropstabInfoTableSQL = `
                CREATE TABLE IF NOT EXISTS dropstab_info (
                    asset_id SERIAL PRIMARY KEY,
                    slug VARCHAR(255) NOT NULL,
                    status VARCHAR(255) NOT NULL,
                    symbol VARCHAR(255) NOT NULL,
                    name VARCHAR(255) NOT NULL,
                    socials JSONB NOT NULL,
                    description TEXT,
                    image_url TEXT,
                    price_usd DOUBLE PRECISION,
                    max_supply DOUBLE PRECISION,
                    circulating_supply DOUBLE PRECISION,
                    total_supply DOUBLE PRECISION,
                    categories TEXT[] NOT NULL DEFAULT '{}'::TEXT[],
                    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
                );

                CREATE INDEX IF NOT EXISTS dropstab_info_symbol_idx 
                    ON dropstab_info (symbol);

                CREATE INDEX IF NOT EXISTS dropstab_info_slug_idx 
                    ON dropstab_info (slug);

                CREATE INDEX IF NOT EXISTS dropstab_info_description_trgm_idx 
                    ON dropstab_info USING gin (description gin_trgm_ops);

                CREATE INDEX IF NOT EXISTS dropstab_info_categories_gin_idx
                    ON dropstab_info USING gin (categories);`

	SelectDropstabInfoByAssetIDSQL = `
                SELECT asset_id, slug, status, symbol, name, socials, description, image_url, price_usd, max_supply, circulating_supply, total_supply, categories, updated_at
                FROM dropstab_info
                WHERE asset_id = $1`

	SelectDropstabInfosByAssetIDsSQL = `
                SELECT asset_id, slug, status, symbol, name, socials, description, image_url, price_usd, max_supply, circulating_supply, total_supply, categories, updated_at
                FROM dropstab_info
                WHERE asset_id = ANY($1)`

	SelectDropstabInfoBySymbolSQL = `
                SELECT asset_id, slug, status, symbol, name, socials, description, image_url, price_usd, max_supply, circulating_supply, total_supply, categories, updated_at
                FROM dropstab_info
                WHERE LOWER(symbol) = LOWER($1)`

	SelectDropstabInfoBySlugSQL = `
                SELECT asset_id, slug, status, symbol, name, socials, description, image_url, price_usd, max_supply, circulating_supply, total_supply, categories, updated_at
                FROM dropstab_info
                WHERE LOWER(slug) = LOWER($1)`

	SelectExistingDropstabInfoSlugsSQL = `
                SELECT slug, updated_at
                FROM dropstab_info
                WHERE slug = ANY($1)`

	SelectAssetIDsBySlugsSQL = `
                SELECT slug, asset_id
                FROM dropstab_info
                WHERE slug = ANY($1)`

	UpsertDropstabInfoSQL = `
                INSERT INTO dropstab_info (asset_id, slug, status, symbol, name, socials, description, image_url, price_usd, max_supply, circulating_supply, total_supply, categories, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
                ON CONFLICT (asset_id) 
                DO UPDATE SET 
                slug = EXCLUDED.slug,
                status = EXCLUDED.status,
                symbol = EXCLUDED.symbol,
                name = EXCLUDED.name,
                socials = EXCLUDED.socials,
                description = EXCLUDED.description,
                image_url = EXCLUDED.image_url,
                price_usd = EXCLUDED.price_usd,
                max_supply = EXCLUDED.max_supply,
                circulating_supply = EXCLUDED.circulating_supply,
                total_supply = EXCLUDED.total_supply,
                categories = EXCLUDED.categories,
                updated_at = EXCLUDED.updated_at
                RETURNING asset_id`

	DeleteDropstabInfoSQL = `
                DELETE FROM dropstab_info 
                WHERE asset_id = $1`
)
