# Adding a new DEX / venue

The indexer is **topic-agnostic by design**: to support a new venue you do NOT
maintain pool or router address lists. You register one *identifier* — a matcher
for the venue's swap-event topic — and delegate amount/user extraction to the
built-in by-transfers heuristic. Forks of the venue start working automatically.

Adding a venue is ~40 lines and three steps. Real example: iZiSwap.

## 1. Add the event topic

`evm/pkg/models/addrs.go`:

```go
// iZiSwap Swap(address,address,bool,uint256,uint256) — X Layer deployment.
var IziSwapEventTopic = common.HexToHash("0xcd3829a3813dc3cdd188fd3d01dcf3268c16be2fdd2dd21d0665418816e46062")
```

How to find the topic: take any tx on the venue, open its receipt, and look at
`topics[0]` of the swap log (`eth_getTransactionReceipt`, or any explorer).

## 2. Create the identifier

`evm/pkg/parsers/swaps/iziswap_pool_log.go` — copy this file for your venue and
change three things: the type name, the topic constant, and the source label.

```go
type IziSwapLogIdentifier struct { /* evmClient, logger */ }

// CheckTx: match if ANY log in the receipt carries the venue's swap topic.
func (i *IziSwapLogIdentifier) CheckTx(req parsers.ParseTxRequest) (bool, error) {
    for _, lg := range req.Receipt.Logs {
        if len(lg.Topics) > 0 && lg.Topics[0] == models.IziSwapEventTopic {
            return true, nil
        }
    }
    return false, nil
}

func (i *IziSwapLogIdentifier) Source() string        { return "iziswap" } // lands in swap_events.source
func (i *IziSwapLogIdentifier) ShouldStop() bool      { return true }
func (i *IziSwapLogIdentifier) GetParserName() string { return "swaps_by_transfers_parser_v2" }
```

Note `GetParserName()`: you delegate to the generic heuristic. You almost never
write amount-parsing code — the heuristic reconstructs base/quote/sender/receiver
from the transaction's net token flows, which also makes aggregator-wrapped and
relayed swaps work for free.

## 3. Register it

`evm/pkg/processors/basicProcessor/registry.go`, inside `IdentifierMap`:

```go
swaps.NewIziSwapLogIdentifier(evmClient, logger): nil,
```

## 4. Rebuild, verify, backfill

```bash
go build ./evm/...
# spot-check on a block that contains the venue's swaps:
./bin/parse_range --start-block N --end-block N+1
# then backfill history for the venue (idempotent — safe over already-parsed ranges):
./bin/parse_range --start-block <deploy_block> --end-block <tip>
```

Verify in ClickHouse:

```sql
SELECT source, count() FROM evm.swap_events WHERE source = 'iziswap' GROUP BY source;
```

## When the venue is NOT uniswap-shaped

If the venue emits no per-pool swap event at all (aggregator routers, custom
settlement), match its router-level event instead — see
`dodo_route_log.go` (DODO `OrderHistory`) for a router-topic example. The
delegation stays the same.
