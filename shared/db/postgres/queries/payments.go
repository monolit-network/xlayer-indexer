package queries

const (
	CreateTablePaymentWalletsSQL = `
	CREATE TABLE IF NOT EXISTS backend.payment_wallets (
		address VARCHAR(255) PRIMARY KEY NOT NULL,
		private_key BYTEA NOT NULL,
		reserved_order_id UUID,
		reserved_at TIMESTAMP,
		cooldown_until TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_backend_payment_wallets_available ON backend.payment_wallets (reserved_order_id, cooldown_until);
	`

	CreateTablePaymentChainInfoSQL = `
	CREATE TABLE IF NOT EXISTS backend.chain_info (
		chain VARCHAR(255) PRIMARY KEY NOT NULL,
		token_addresses VARCHAR(255)[] NOT NULL,
		latest_block_offset BIGINT NOT NULL,
		decimals INTEGER[] NOT NULL,
		token_names TEXT[] NOT NULL,
		token_image_urls TEXT[] NOT NULL,
		url VARCHAR(255) NOT NULL,
		CONSTRAINT chain_info_tokens_len_chk CHECK (
			cardinality(token_addresses) = cardinality(decimals)
			AND cardinality(token_addresses) = cardinality(token_names)
			AND cardinality(token_addresses) = cardinality(token_image_urls)
		)
	);`

	SelectAvailablePaymentWalletFromCursorForUpdateSQL = `
	SELECT address, private_key, reserved_order_id, reserved_at, cooldown_until
	FROM backend.payment_wallets
	WHERE reserved_order_id IS NULL
		AND (cooldown_until IS NULL OR cooldown_until <= NOW())
		AND address >= $1
	ORDER BY address
	FOR UPDATE SKIP LOCKED
	LIMIT 1;
	`

	SelectAvailablePaymentWalletBeforeCursorForUpdateSQL = `
	SELECT address, private_key, reserved_order_id, reserved_at, cooldown_until
	FROM backend.payment_wallets
	WHERE reserved_order_id IS NULL
		AND (cooldown_until IS NULL OR cooldown_until <= NOW())
		AND address < $1
	ORDER BY address
	FOR UPDATE SKIP LOCKED
	LIMIT 1;
	`

	SelectPaymentWalletsSQL = `
	SELECT address, private_key, reserved_order_id, reserved_at, cooldown_until
	FROM backend.payment_wallets
	ORDER BY address;
	`

	InsertPaymentWalletSQL             = `INSERT INTO backend.payment_wallets (address, private_key, reserved_order_id, reserved_at, cooldown_until) VALUES ($1, $2, $3, $4, $5);`
	ReservePaymentWalletSQL            = `UPDATE backend.payment_wallets SET reserved_order_id = $2, reserved_at = $3, cooldown_until = NULL WHERE address = $1;`
	ReleasePaymentWalletReservationSQL = `UPDATE backend.payment_wallets SET reserved_order_id = NULL, reserved_at = NULL, cooldown_until = NOW() + INTERVAL '30 minutes' WHERE address = $1;`

	SelectPaymentChainInfosSQL = `SELECT chain, token_addresses, latest_block_offset, decimals, token_names, token_image_urls, url FROM backend.chain_info;`
)
