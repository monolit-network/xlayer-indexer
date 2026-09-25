CREATE DATABASE IF NOT EXISTS evm;

CREATE TABLE evm.swap_events
(
    `chain` LowCardinality(String),
    `block_time` DateTime('UTC'),
    `block_number` UInt256,
    `block_hash` FixedString(66),
    `tx_idx` UInt32,
    `tx_hash` FixedString(66),
    `gas_used` UInt64,
    `value` UInt256,
    `tx_from_address` FixedString(42),
    `tx_to_address` FixedString(42),
    `source` LowCardinality(String),
    `instruction_hash` FixedString(8),
    `base_coin` FixedString(42),
    `quote_coin` FixedString(42),
    `base_coin_amount` UInt256,
    `base_coin_decimals` UInt8,
    `quote_coin_amount` UInt256,
    `quote_coin_decimals` UInt8,
    `sender` FixedString(42),
    `receiver` FixedString(42),
    INDEX idx_block_time block_time TYPE minmax GRANULARITY 1,
    INDEX idx_swap_participants (sender, receiver, tx_from_address, tx_to_address) TYPE bloom_filter(0.001) GRANULARITY 1,
    INDEX idx_swap_coins (base_coin, quote_coin) TYPE bloom_filter(0.001) GRANULARITY 1,
    INDEX idx_swap_tx tx_hash TYPE bloom_filter(0.001) GRANULARITY 1
)
ENGINE = ReplacingMergeTree(block_time)
PARTITION BY (chain, toYYYYMM(block_time))
ORDER BY (chain, block_number, tx_hash, sender, receiver)
SETTINGS index_granularity = 8192;

CREATE TABLE evm.defi_events
(
    `chain` String,
    `block_time` DateTime('UTC'),
    `block_number` UInt256,
    `block_hash` FixedString(66),
    `tx_idx` UInt32,
    `tx_hash` FixedString(66),
    `gas_used` UInt64,
    `value` UInt256,
    `tx_from_address` FixedString(42),
    `tx_to_address` FixedString(42),
    `source` String,
    `instruction_hash` FixedString(8),
    `inputs.address` Array(FixedString(42)),
    `inputs.amount` Array(UInt256),
    `inputs.decimals` Array(UInt8),
    `outputs.address` Array(FixedString(42)),
    `outputs.amount` Array(UInt256),
    `outputs.decimals` Array(UInt8),
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
SETTINGS index_granularity = 8192;

CREATE TABLE evm.transfer_events
(
    `chain` LowCardinality(String),
    `block_time` DateTime('UTC'),
    `block_number` UInt256,
    `block_hash` FixedString(66),
    `tx_idx` UInt32,
    `tx_hash` FixedString(66),
    `gas_used` UInt64,
    `sender` FixedString(42),
    `receiver` FixedString(42),
    `token_address` FixedString(42),
    `token_decimals` UInt8,
    `amount` UInt256,
    `source` LowCardinality(String) DEFAULT 'erc20_or_native_transfers',
    INDEX idx_block_time block_time TYPE minmax GRANULARITY 1,
    INDEX idx_transfer_participants (sender, receiver, token_address) TYPE bloom_filter(0.001) GRANULARITY 1,
    INDEX idx_transfer_tx tx_hash TYPE bloom_filter(0.001) GRANULARITY 1
)
ENGINE = ReplacingMergeTree
PARTITION BY (chain, toYYYYMM(block_time))
ORDER BY (chain, block_number, tx_hash, sender, receiver)
SETTINGS index_granularity = 8192;

CREATE TABLE evm.error_events
(
    `chain` String,
    `block_number` UInt256,
    `tx_hash` FixedString(66),
    `error` String
)
ENGINE = ReplacingMergeTree
PARTITION BY chain
ORDER BY (chain, block_number, tx_hash)
SETTINGS index_granularity = 8192;

