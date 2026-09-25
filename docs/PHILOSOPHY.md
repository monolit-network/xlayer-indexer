# Parser philosophy

## The problem with decoding events

The standard way to index DEX activity is to *decode protocol events*: know the
ABI, parse `Swap(amount0In, amount1Out, ...)`, trust the payload. It works on the
happy path and silently falls apart everywhere else:

- **Aggregators and routers** emit the pool's event, but the pool only saw the
  router. The "trader" you index is a contract shared by thousands of users.
- **Account abstraction (ERC-4337)** breaks it twice: the tx signer is a bundler,
  and one tx carries many users' operations.
- **Relayed / solver-executed trades**: the signer is the operator's EOA; the
  economic party appears nowhere in the event.
- **Forks**: every new Uniswap fork needs a new adapter, a new address list, a
  redeploy. Coverage decays weekly on a young chain.
- **Lying payloads**: fee-on-transfer tokens, non-standard decimals, events whose
  amounts are denominated in units you didn't expect. The event says one thing;
  the tokens moved another way.

## The principle: index what the user experienced

A swap *is* a movement of tokens: someone's balance of A went down, their (or
their designated recipient's) balance of B went up. Events can lie or be missing;
**net token flows cannot** — they are the chain's ground truth.

Just as important: we index the **top-level economic action**, not the plumbing.
Intermediate hops inside a route are analytically meaningless — no human decided
to hold the middle token for 400 milliseconds. They cancel out in net flows, so
the record collapses to the one thing the user actually did.

**Concrete example** (real X Layer tx
[`0x4f0bce34…`](https://web3.okx.com/explorer/x-layer/tx/0x4f0bce34afc86799ff497ce7f96e96b2fc475da515a817db3bc4c6c13dc27681)):
the receipt contains two `V3.Swap` pool events, four ERC-20 transfers and two
router events. Per-event indexers record two pool swaps whose "trader" is a
router, denominated partly in an intermediate stable leg the user never held.
We record **one row**:

> `0xbac6…` sold **1.93075925 wSPYx** (tokenized S&P 500) and received
> **1500.494795 USDC**.

That row matches Debank's independent parse of the same tx digit-for-digit —
and it is the row an analyst, a PnL model, or an AI agent actually wants.

How it works mechanically:

1. **Events are only triggers.** An *identifier* watches for a known swap-event
   topic anywhere in the receipt — topic-agnostic, no pool/router address lists.
   A fork of Uniswap emits the same topic → indexed from day one, zero code.
2. **All facts come from flows.** A generic heuristic aggregates every ERC-20
   transfer (plus wrapped-native deposits/withdrawals) in the tx into per-address
   net balances, then answers three questions:
   - *Who paid?* — the clean address with a net outflow of exactly one token
     (that token = base, that address = sender).
   - *Who received?* — preference order: the sender themselves; else the single
     clean net-inflow third party; else (last resort) the custodial contract that
     kept the proceeds — recorded as such, with the seller still correctly
     attributed.
   - *How much?* — net amounts, not event payloads.
3. **Special forms get structure, not special amounts.** ERC-4337 bundles are cut
   into per-`UserOperationEvent` segments and each segment runs through the same
   heuristic with the smart account as the known sender. Relayed swaps recover
   the sender by descending call traces. In every case the *arithmetic* stays in
   one place — the flow heuristic — so a fix there fixes every venue at once.

## Every transaction is supported from day one

The table layout is a funnel, not a filter:

- **`swap_events`** holds *only* fully-shaped swaps — the graduated, high-trust
  rows.
- **`defi_events`** is the universal catch-all: **every smart-contract
  interaction we did not recognize** still gets a row with the signer's net
  inflows/outputs. Nothing about the chain is invisible — you can already run
  analytics on protocols we never heard of (volumes, users, token flows).
- When a category becomes worth first-class treatment (lending, bridges,
  launchpads), you *graduate* it: add an identifier + a dedicated heuristic and a
  sibling table (`lending_events`, `bridge_events`, ...) by the same recipe as
  `swap_events` — see [ADDING_A_DEX.md](ADDING_A_DEX.md). Rows migrate from the
  catch-all into shaped form, and coverage ratchets up monotonically.
- **`transfer_events`** holds plain sends (no contract logic), and
  **`error_events`** holds parse attempts that failed, with reasons.

## Failure is data

Any tx the heuristic cannot resolve is written to `error_events` with a reason —
never dropped. This turns parser gaps into a queryable backlog:

- error reasons form a taxonomy ("no suitable receiver", "multiple base tokens",
  ...), so a whole *class* of misses is one `GROUP BY` away from being visible;
- after a parser improvement, `--blocks-file` re-parses exactly the affected
  blocks; `ReplacingMergeTree` makes the re-run idempotent.

This loop is how the indexer grew on X Layer: the "no suitable receiver" class
turned out to be custodial routers (platforms keeping user proceeds in the
contract) — one fallback later, millions of previously-lost trades were recovered
and the platforms themselves became mappable infrastructure.

## Attribution as a view, not a rewrite

Raw rows keep the mechanical sender/receiver. Resolving "the human behind it"
(router → signer, bundler → smart-account owner) happens in a SQL view over an
infra-address label set. Labels improve → history re-attributes instantly, with
no re-parse and no destroyed evidence. The raw record stays forensically honest.

## Trade-offs (honest section)

- **Traces cost.** Call-trace descent needs `debug_trace*`; on trace-heavy blocks
  the node, not the parser, is the throughput ceiling.
- **Heuristics have edge cases.** Wash-loops with net-zero flows, dust arbitrage
  with no clean party, burn-mechanic tokens — these land in `error_events` by
  design rather than being guessed at. We prefer an honest gap to a plausible lie.
- **LP operations are not swaps.** Add/remove liquidity produces two-token flows
  and is deliberately rejected by the swap heuristic (it lands in `defi_events`).
