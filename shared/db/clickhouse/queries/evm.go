package queries

const (
	CreateSchemaEvmClickhouseSQL = `CREATE DATABASE IF NOT EXISTS evm`

	CreateTableEvmSwapEventsSQL = `
        CREATE TABLE IF NOT EXISTS evm.swap_events
        (
            chain String,
            block_time DateTime('UTC'),
            block_number UInt256,
            block_hash FixedString(66),
            tx_idx UInt32,
            tx_hash FixedString(66),
            gas_used UInt64,
            value UInt256,
            tx_from_address FixedString(42),
            tx_to_address FixedString(42),
            source String,
            instruction_hash FixedString(8),
            base_coin FixedString(42),
            quote_coin FixedString(42),
            base_coin_amount UInt256,
            base_coin_decimals UInt8,
            quote_coin_amount UInt256,
            quote_coin_decimals UInt8,
            sender FixedString(42),
            receiver FixedString(42),
            INDEX idx_swap_source source TYPE set(0) GRANULARITY 1,
            INDEX idx_swap_participants (sender, receiver, tx_from_address, tx_to_address) TYPE bloom_filter(0.001) GRANULARITY 1,
            INDEX idx_swap_coins (base_coin, quote_coin) TYPE bloom_filter(0.001) GRANULARITY 1,
            INDEX idx_swap_tx tx_hash TYPE bloom_filter(0.001) GRANULARITY 1
        )
        ENGINE = ReplacingMergeTree
        PARTITION BY (chain, toYYYYMM(block_time))
        ORDER BY (chain, block_number, tx_hash, sender, receiver)
        SETTINGS index_granularity = 8192
    `

	SelectEvmSwapEventsSQL = `
        SELECT chain, block_time, block_number, block_hash, tx_idx, tx_hash, gas_used, value, tx_from_address, tx_to_address, source, instruction_hash,
            base_coin, quote_coin, base_coin_amount, base_coin_decimals, quote_coin_amount, quote_coin_decimals, sender, receiver
        FROM evm.swap_events
        WHERE chain = ? AND block_number >= ? AND block_number <= ?
        ORDER BY block_number ASC
    `

	InsertEvmSwapEventsSQL = `
        INSERT INTO evm.swap_events (
            chain, block_time, block_number, block_hash, tx_idx, tx_hash, gas_used, value, tx_from_address, tx_to_address, source, instruction_hash,
            base_coin, quote_coin, base_coin_amount, base_coin_decimals, quote_coin_amount, quote_coin_decimals, sender, receiver
        )`

	CreateTableEvmTransferEventsSQL = `
        CREATE TABLE IF NOT EXISTS evm.transfer_events
        (
            chain String,
            block_time DateTime('UTC'),
            block_number UInt256,
            block_hash FixedString(66),
            tx_idx UInt32,
            tx_hash FixedString(66),
            gas_used UInt64,
            sender FixedString(42),
            receiver FixedString(42),
            token_address FixedString(42),
            token_decimals UInt8,
            amount UInt256,
            source LowCardinality(String) DEFAULT 'erc20_or_native_transfers',
            INDEX idx_tr_amount amount TYPE minmax GRANULARITY 4,
            INDEX idx_tr_addresses (sender, receiver, token_address) TYPE bloom_filter(0.001) GRANULARITY 1,
            INDEX idx_tr_tx_hash tx_hash TYPE minmax GRANULARITY 4
        )
        ENGINE = ReplacingMergeTree
        PARTITION BY (chain, toYYYYMM(block_time))
        ORDER BY (chain, block_number, tx_hash, sender, receiver)
        SETTINGS index_granularity = 8192
        `

	SelectEvmTransferEventsSQL = `
        SELECT chain, block_time, block_number, block_hash, tx_idx, tx_hash, gas_used, sender, receiver, token_address, token_decimals, amount
        FROM evm.transfer_events
        WHERE chain = ? AND block_number >= ? AND block_number <= ?
        ORDER BY block_number ASC
    `

	InsertEvmTransferEventsSQL = `
        INSERT INTO evm.transfer_events (
            chain, block_time, block_number, block_hash, tx_idx, tx_hash, gas_used,
            sender, receiver, token_address, token_decimals, amount, source
        )`

	CreateTableEvmDefiEventsSQL = `
        CREATE TABLE IF NOT EXISTS evm.defi_events
        (
            chain String,
            block_time DateTime('UTC'),
            block_number UInt256,
            block_hash FixedString(66),
            tx_idx UInt32,
            tx_hash FixedString(66),
            gas_used UInt64,
            value UInt256,
            tx_from_address FixedString(42),
            tx_to_address FixedString(42),
            source String,
            instruction_hash FixedString(8),
            
            inputs Nested (
                address FixedString(42),
                amount UInt256,
                decimals UInt8
            ),
            outputs Nested (
                address FixedString(42),
                amount UInt256,
                decimals UInt8
            ),

            INDEX idx_addresses (tx_from_address, tx_to_address) TYPE bloom_filter(0.01) GRANULARITY 1,
            INDEX idx_io_addresses (inputs.address, outputs.address) TYPE bloom_filter(0.01) GRANULARITY 1,
            INDEX idx_inst_hash instruction_hash TYPE bloom_filter(0.01) GRANULARITY 1,
            INDEX idx_source source TYPE set(0) GRANULARITY 1,
            INDEX idx_value value TYPE minmax GRANULARITY 1,
            INDEX idx_defi_tx tx_hash TYPE bloom_filter(0.001) GRANULARITY 1
        )
        ENGINE = ReplacingMergeTree
        PARTITION BY (chain, toYYYYMM(block_time))
        ORDER BY (chain, block_number, tx_hash)
        SETTINGS index_granularity = 8192
        `

	InsertEvmDefiEventsSQL = `
        INSERT INTO evm.defi_events (
            chain, block_time, block_number, block_hash, tx_idx, tx_hash, gas_used, value, tx_from_address, tx_to_address, source, instruction_hash,
            inputs.address, inputs.amount, inputs.decimals, outputs.address, outputs.amount, outputs.decimals
        )`

	SelectEvmDefiEventsSQL = `
        SELECT chain, block_time, block_number, block_hash, tx_idx, tx_hash, gas_used, value, tx_from_address, tx_to_address, source, instruction_hash,
            inputs.address, inputs.amount, inputs.decimals, outputs.address, outputs.amount, outputs.decimals
        FROM evm.defi_events
        WHERE chain = ? AND block_number >= ? AND block_number <= ? AND tx_to_address = ? AND instruction_hash = ?
        ORDER BY block_number ASC
    `

	DeleteEvmDefiEventsSQL = `
        ALTER TABLE evm.defi_events DELETE
        WHERE chain = ? 
        AND (block_number, tx_hash) IN arrayZip(
            ?::Array(UInt256), 
            ?::Array(FixedString(66))
        )
    `

	DeleteEvmDefiEventsByContractAddressSQL = `
    ALTER TABLE evm.defi_events DELETE
    WHERE chain = ?
    AND tx_to_address = ?
    AND instruction_hash = ?
    `

	CreateTableEvmErrorEventsSQL = `
        CREATE TABLE IF NOT EXISTS evm.error_events (
            chain String,
            block_number UInt256,
            tx_hash FixedString(66),
            error String
        )
        ENGINE = ReplacingMergeTree()
        PARTITION BY chain
        ORDER BY (chain, block_number, tx_hash);
        `

	InsertEvmErrorEventsSQL = `
        INSERT INTO evm.error_events (
            chain, block_number, tx_hash, error
        )`

	SelectEvmErrorEventsSQL = `
        SELECT chain, block_number, tx_hash, error
        FROM evm.error_events
        WHERE chain = ? AND block_number >= ? AND block_number <= ?
        ORDER BY block_number ASC
    `

	CreateTableEvmPolymarketOrderEventsSQL = `
        CREATE TABLE IF NOT EXISTS evm.polymarket_order_events
        (
            block_time DateTime,
            block_number UInt256,
            block_hash FixedString(66),
            tx_idx UInt32,
            tx_hash FixedString(66),
            tx_from_address FixedString(42),
            tx_to_address FixedString(42),
            instruction_hash FixedString(8),

            id FixedString(133), -- ({txHashWith0x}_{orderHashWith0x}
            transaction_hash FixedString(66),
            order_hash FixedString(66),
            maker FixedString(42),
            taker FixedString(42),
            maker_asset_id String,
            taker_asset_id String,
            maker_amount_filled UInt256,
            taker_amount_filled UInt256,
            fee Decimal(76, 18),
            is_deleted UInt8 DEFAULT 0,

            INDEX idx_maker (maker) TYPE bloom_filter GRANULARITY 1024,
            INDEX idx_taker (taker) TYPE bloom_filter GRANULARITY 1024,
            INDEX idx_maker_asset_id (maker_asset_id) TYPE bloom_filter GRANULARITY 1024,
            INDEX idx_taker_asset_id (taker_asset_id) TYPE bloom_filter GRANULARITY 1024,
            INDEX idx_transaction_hash (transaction_hash) TYPE bloom_filter GRANULARITY 1024,
        ) 
        ENGINE = ReplacingMergeTree
        ORDER BY (block_number, id)
        PARTITION BY (toYYYYMM(block_time))
        SETTINGS index_granularity = 8192`

	InsertEvmPolymarketOrderEventsSQL = `
        INSERT INTO evm.polymarket_order_events (
            block_time, block_number, block_hash, tx_idx, tx_hash, tx_from_address, tx_to_address, instruction_hash,
            id, transaction_hash, order_hash, maker, taker, maker_asset_id, taker_asset_id, maker_amount_filled, taker_amount_filled, fee, is_deleted
        )`

	// Delete by block_hash for reorg handling
	DeleteEvmSwapEventsByBlockHashSQL = `
		ALTER TABLE evm.swap_events DELETE WHERE chain = ? AND block_hash = lower(?)
	`

	DeleteEvmTransferEventsByBlockHashSQL = `
		ALTER TABLE evm.transfer_events DELETE WHERE chain = ? AND block_hash = lower(?)
	`

	DeleteEvmDefiEventsByBlockHashSQL = `
		ALTER TABLE evm.defi_events DELETE WHERE chain = ? AND block_hash = lower(?)
	`

	DeleteEvmPolymarketOrderEventsByBlockHashSQL = `
		ALTER TABLE evm.polymarket_order_events DELETE WHERE block_hash = lower(?)
	`

	// Delete by block_number for reorg handling (>= block_number)
	DeleteEvmSwapEventsFromBlockNumberSQL = `
		ALTER TABLE evm.swap_events DELETE WHERE chain = ? AND block_number >= ?
	`

	DeleteEvmTransferEventsFromBlockNumberSQL = `
		ALTER TABLE evm.transfer_events DELETE WHERE chain = ? AND block_number >= ?
	`

	DeleteEvmDefiEventsFromBlockNumberSQL = `
		ALTER TABLE evm.defi_events DELETE WHERE chain = ? AND block_number >= ?
	`

	DeleteEvmPolymarketOrderEventsFromBlockNumberSQL = `
		ALTER TABLE evm.polymarket_order_events DELETE WHERE block_number >= ?
	`
)
