# Deployment guide

Four moving parts: an X Layer node you can query, ClickHouse, PostgreSQL, and
the two binaries from this repo. All wiring happens through environment
variables (`.env.example` is the complete list).

## 1. X Layer node

- JSON-RPC over HTTP **and** WebSocket.
- The `debug` namespace **must** be enabled: the parsers require
  `debug_traceBlockByNumber` (callTracer) to see internal calls and native
  transfers — this is what makes relayed/AA swap attribution work.
  For op-geth/reth style nodes: `--http.api eth,net,web3,debug --ws.api eth,net,web3,debug`.
- History: standard OP-stack nodes serve X Layer from the migration block
  **45,000,000**. Blocks before that require a legacy zkEVM archive node.

```
EVM_RPC_URL_XLAYER=http://node:8545
EVM_WS_URL_XLAYER=ws://node:8546
```

## 2. ClickHouse (the event store)

Any recent ClickHouse (tested on 25.x/26.x). One database, four tables:

```bash
clickhouse-client < sql/schema.sql   # includes CREATE DATABASE evm
```

```
CLICKHOUSE_HOST=clickhouse    CLICKHOUSE_PORT=9000   # native protocol port
CLICKHOUSE_USER=default      CLICKHOUSE_PASSWORD=...
CLICKHOUSE_DATABASE=evm
```

All event tables are ReplacingMergeTree → any re-parse is idempotent.

## 3. PostgreSQL (tiny, but required)

A stock Postgres with one empty database is enough:

```
DATABASE_URL=postgres://user:pass@postgres:5432/indexer
```

What actually lives here (verified against the code, so you don't wonder):

- **one small table, auto-created at startup** — an optional override registry
  mapping `contract + method → parser/source` for special-casing;
- that's it. The **resume point does *not* live in Postgres** — on start the
  indexer takes `max(block_number)` from the ClickHouse event tables and
  continues from there. Reorg detection is an in-memory ring of recent block
  hashes, rebuilt from the chain on start. You can lose the Postgres volume and
  nothing of value is gone.

## 4. Local ClickHouse + Postgres via docker compose

```yaml
services:
  clickhouse:
    image: clickhouse/clickhouse-server:26.1
    ports: ["9000:9000", "8123:8123"]
    volumes: ["ch-data:/var/lib/clickhouse"]
  postgres:
    image: postgres:16
    environment: { POSTGRES_USER: indexer, POSTGRES_PASSWORD: indexer, POSTGRES_DB: indexer }
    ports: ["5432:5432"]
    volumes: ["pg-data:/var/lib/postgresql/data"]
volumes: { ch-data: {}, pg-data: {} }
```

## 5. First run

```bash
cp .env.example .env            # fill in the three connection groups above
clickhouse-client < sql/schema.sql
go build -o bin/indexer ./evm/cmd/indexer
go build -o bin/parse_range ./evm/cmd/parse_range

# fresh install: backfill history first (from the migration block, or any
# later block you care about), in one or many parallel disjoint ranges:
./bin/parse_range --start-block 45000000 --end-block 46000000

# then start live indexing — it resumes from the highest indexed block:
./bin/indexer
```

## 6. Verify it works

```bash
curl localhost:19096/status
# {"status":"running","chain":"xlayer","last_processed":N,"latest":N,"lag":0,...}
```

```sql
SELECT source, count() FROM evm.swap_events WHERE chain='xlayer' GROUP BY source;
SELECT substring(error,1,60) e, count() FROM evm.error_events WHERE chain='xlayer' GROUP BY e ORDER BY 2 DESC LIMIT 5;
```

If `lag` stays near 0 and swap sources are filling up, you're live.

## Sizing hints

- The parser is CPU-light; **the node's trace endpoint is the throughput
  ceiling**. Concurrency knobs: `EVM_BLOCK_BATCH_MAX_CONCURRENT_{BLOCKS,RECEIPTS,TRACES}_XLAYER`.
- Full OP-era backfill (≈26M blocks) fits in hours with 5–10 parallel
  `parse_range` workers against a dedicated node.
- ClickHouse footprint for full X Layer history: single-digit GB compressed.
