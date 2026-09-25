# xlayer-indexer

**The first complete DEX index of X Layer** — OKX's L2 ("The New Money Chain").

Indexes every swap on the network in real time and across full history, and — unlike
generic indexers — resolves the **real user** behind routers, aggregators, and
ERC-4337 smart accounts.

Built during [OKX Dev Day 2026](https://luma.com/l4aq8vii) as the data layer for our
AI agent pipeline; released as a standalone open-source contribution to the X Layer
ecosystem.

## What it covers

| | |
|---|---|
| Swaps indexed | **24.4M** (full OP-era history, block 45,000,000 → tip) |
| Unique trading wallets | **508k** |
| Live indexing lag | **0 blocks** |
| Protocols | Uniswap v2 / v3 / v4, Curve, DODO, iZiSwap — incl. all forks (PotatoSwap, OkieSwap, ...) |
| Error rate | every unparsed tx recorded in `error_events` with a reason — nothing silently dropped |

## Philosophy: truth is in token movements

Most indexers decode protocol-specific event payloads. That breaks the moment a tx
goes through an aggregator, a custom router, a bundler, or a contract the indexer
has never seen.

This indexer treats **event topics only as triggers**. Amounts, direction, and the
actual user are derived from the *net token flows* of the transaction:

1. **Identifiers** match known swap-event topics anywhere in the receipt
   (topic-agnostic: no pool/router address lists to maintain — forks work day one).
2. The **by-transfers heuristic** reconstructs who *paid* what and who *received*
   what from ERC-20 transfers + native flows.
3. Special forms are handled explicitly:
   - **ERC-4337 bundles** are segmented per `UserOperationEvent` — each user op in a
     bundle becomes its own swap, attributed to the smart-account owner, not the bundler;
   - **relayed / solver-executed swaps** recover the sender via call-trace descent;
   - **custodial routers** (proceeds kept by the platform contract) are recognized,
     recorded, and the seller still gets correct attribution.

## Architecture

```
X Layer node (RPC/WS)
        │ blocks, receipts, traces
        ▼
  indexer (live head, lag 0)          parse_range (backfill / targeted repair)
        │                                   │  --start/--end or --blocks-file
        └────────────► parser registry ◄────┘
                (identifiers → by-transfers heuristic)
                            │
                            ▼
                       ClickHouse
        swap_events · defi_events · transfer_events · error_events
```

- `evm/cmd/indexer` — live indexer with reorg handling and persistent state.
- `evm/cmd/parse_range` — batch backfill; `--blocks-file` mode parses an exact
  block list (9× faster for targeted repairs than range sweeps).
- `sql/schema.sql` — ClickHouse DDL (ReplacingMergeTree → idempotent re-parses;
  re-running any range is always safe).

## Quickstart

```bash
cp .env.example .env        # point to your X Layer node + ClickHouse + Postgres
clickhouse-client < sql/schema.sql
go build -o bin/indexer ./evm/cmd/indexer && ./bin/indexer          # live
go build -o bin/parse_range ./evm/cmd/parse_range
./bin/parse_range --start-block 45000000 --end-block 45100000       # backfill
```

Note: X Layer migrated from zkEVM to the OP stack at block **45,000,000** — public
OP-stack nodes serve history from that block. Earlier history requires a zkEVM
archive node.

## Attribution in queries

Raw rows keep the mechanical sender/receiver. To resolve the human behind
infrastructure, join your infra-address labels and pick the first non-infra party —
example view in [`sql/schema.sql`](sql/schema.sql) comments:

```sql
SELECT attributed_user, count() AS swaps
FROM v_swap_events_attributed
WHERE chain = 'xlayer'
GROUP BY attributed_user ORDER BY swaps DESC
```

## License

Apache-2.0
