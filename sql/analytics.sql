-- Starter analytics over swap_events.
-- Everything here is plain SQL over the indexed data: no extra services needed.

-- ============================================================================
-- 1. Attribution: resolve the economic user behind infrastructure addresses.
--    The MECHANISM ships here; the LABELS are yours to curate. Fill
--    infra_labels with routers / bundlers / custodial vaults you identify
--    (a good starting signal: receivers with thousands of unique senders).
-- ============================================================================
CREATE TABLE IF NOT EXISTS evm.infra_labels
(
    address FixedString(42),
    tag     LowCardinality(String),   -- e.g. 'router', 'bundler', 'custodial-vault', 'burn'
    note    String DEFAULT ''
)
ENGINE = ReplacingMergeTree
ORDER BY address;

INSERT INTO evm.infra_labels VALUES ('0x0000000000000000000000000000000000000000','burn',''), ('0x000000000000000000000000000000000000dead','burn','');

CREATE VIEW IF NOT EXISTS evm.v_swap_events_attributed AS
WITH infra AS (SELECT address FROM evm.infra_labels)
SELECT
    *,
    multiIf(
        toString(sender)   NOT IN (SELECT toString(address) FROM infra), toString(sender),
        toString(receiver) NOT IN (SELECT toString(address) FROM infra), toString(receiver),
        toString(tx_from_address) NOT IN (SELECT toString(address) FROM infra), toString(tx_from_address),
        'unknown') AS attributed_user
FROM evm.swap_events;

-- ============================================================================
-- 2. Token stats: volume, activity, breadth per token.
-- ============================================================================
CREATE VIEW IF NOT EXISTS evm.v_token_stats AS
SELECT
    chain,
    base_coin                        AS token,
    count()                          AS swaps,
    uniq(attributed_user)            AS traders,
    sum(base_coin_amount)            AS token_volume_raw,   -- divide by 10^decimals
    any(base_coin_decimals)          AS decimals,
    min(block_time)                  AS first_seen,
    max(block_time)                  AS last_seen
FROM evm.v_swap_events_attributed
GROUP BY chain, token;

-- ============================================================================
-- 3. Wallet stats: activity profile per attributed user.
-- ============================================================================
CREATE VIEW IF NOT EXISTS evm.v_wallet_stats AS
SELECT
    chain,
    attributed_user                  AS wallet,
    count()                          AS swaps,
    uniq(base_coin)                  AS tokens_traded,
    min(block_time)                  AS first_swap,
    max(block_time)                  AS last_swap,
    uniq(toDate(block_time))         AS active_days
FROM evm.v_swap_events_attributed
WHERE attributed_user != 'unknown'
GROUP BY chain, wallet;

-- Example: most active wallets on X Layer, last 7 days
--   SELECT wallet, count() s, uniq(base_coin) toks
--   FROM evm.v_swap_events_attributed
--   WHERE chain='xlayer' AND block_time > now() - INTERVAL 7 DAY
--   GROUP BY wallet ORDER BY s DESC LIMIT 20;
