package queries

const (
	InsertEvmPolymarketOrderEventsNewSQL = `
		INSERT INTO evm.polymarket_order_events (
			block_time, block_number, block_hash, tx_idx, tx_hash, condition_id, log_index, sub_index,
			user_address, source, source_address, token_id, token_amount_diff, usdc_amount_diff, is_taker, fee
		)`
)
