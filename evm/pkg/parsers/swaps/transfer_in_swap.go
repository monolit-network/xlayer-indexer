package swaps

import (
	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"

	"github.com/ethereum/go-ethereum/common"
)

// terminalTransferEvents emits one TransferEvent (source="transfer_in_swap")
// for every ERC20 leg whose recipient does NOT forward the same token later
// in the same tx — i.e. the final beneficiaries of a trading transaction
// (buyer, fee sink, airdropped influencer, settlement vault). Intermediate
// pool/router hops are skipped by construction, which keeps the volume sane.
// Burns (to 0x0) are skipped as well.
func terminalTransferEvents(
	base models.BaseEvent,
	transfers []models.ERC20TransferLog,
	cl *evmclient.Client,
) []models.Event {
	type edge struct{ token, addr common.Address }
	hasOut := make(map[edge]bool, len(transfers))
	for _, tr := range transfers {
		if tr.Source != "erc20" {
			continue
		}
		hasOut[edge{tr.TokenAddress, tr.FromAddress}] = true
	}
	seen := make(map[string]bool, 4)
	var out []models.Event
	for _, tr := range transfers {
		if tr.Source != "erc20" || tr.Value == nil || tr.Value.Sign() <= 0 {
			continue
		}
		if tr.ToAddress == (common.Address{}) {
			continue
		}
		if hasOut[edge{tr.TokenAddress, tr.ToAddress}] {
			continue // forwarded on — not terminal
		}
		k := tr.TokenAddress.Hex() + tr.FromAddress.Hex() + tr.ToAddress.Hex() + tr.Value.String()
		if seen[k] {
			continue
		}
		seen[k] = true
		be := base
		be.Source = "transfer_in_swap"
		ev := models.TransferEvent{
			BaseEvent:    be,
			FromAddress:  tr.FromAddress,
			ToAddress:    tr.ToAddress,
			TokenAddress: tr.TokenAddress,
			Amount:       tr.Value,
		}
		if d, err := cl.GetDecimals(tr.TokenAddress); err == nil {
			dd := d
			ev.TokenDecimals = &dd
		}
		out = append(out, &ev)
	}
	return out
}
