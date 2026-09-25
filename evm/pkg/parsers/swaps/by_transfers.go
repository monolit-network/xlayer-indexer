package swaps

import (
	"fmt"
	"math/big"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type ByTransfersParser struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewByTransfersParser(evmClient *evmclient.Client, logger *zap.Logger) *ByTransfersParser {
	p := &ByTransfersParser{evmClient: evmClient}
	p.logger = logger.Named(p.Name())
	return p
}

func (p ByTransfersParser) Name() string {
	return "swaps_by_transfers_parser"
}

func (p *ByTransfersParser) ProducedEventsType() models.EventType {
	return models.EventTypeSwap
}

func (p *ByTransfersParser) RequiredData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesBlockHeader | evmclient.RequiredDataTypesTraces
}

func (p *ByTransfersParser) OptionalData() evmclient.RequiredDataTypes {
	return 0
}

func (p *ByTransfersParser) ParseTx(req parsers.ParseTxRequest) ([]models.Event, error) {
	if err := parsers.ValidateParseTxRequest(req, p.RequiredData()); err != nil {
		return nil, err
	}
	baseEvent, err := parsers.MakeBaseEvent(req, p.evmClient)
	if err != nil {
		return nil, err
	}

	transfers, err := models.FindAllTransfers(req.Receipt, req.Tx, req.Traces, p.evmClient.ChainID)
	if err != nil {
		return nil, err
	}

	if len(transfers) == 0 {
		return nil, fmt.Errorf("no transfers found, must be at least 1")
	}

	// 1) Normalize transfers (handle native wraps/unlocks to avoid double counting)
	normalizedTransfers := normalizeTransfersForSwapInference(transfers, req)
	p.logger.Debug("normalized_transfers", zap.Any("transfers", normalizedTransfers))

	// 2) Compute per-address per-token net deltas and basic flow stats
	deltas, flowStats := computeAddressTokenDeltas(normalizedTransfers)

	// 3) Classify nodes (extensible: keep visible for user to extend)
	nodeClass := classifyAddresses(deltas, flowStats)
	// Ensure tx sender is treated as non-intermediary; tx recipient (router) as intermediary
	delete(nodeClass.LikelyPools, baseEvent.TxFromAddress)
	delete(nodeClass.LikelyRouters, baseEvent.TxFromAddress)
	if (baseEvent.TxToAddress != common.Address{}) {
		nodeClass.LikelyRouters[baseEvent.TxToAddress] = struct{}{}
		delete(nodeClass.LikelyPools, baseEvent.TxToAddress)
	}
	p.logger.Debug("node_classification", zap.Any("pools", nodeClass.LikelyPools), zap.Any("routers", nodeClass.LikelyRouters))

	// 4) Choose tokenIn/tokenOut, payer and recipient using deltas on non-pool/non-router addresses
	tokenIn, payer := pickTokenInAndPayer(deltas, nodeClass, baseEvent.TxFromAddress)
	tokenOut, recipient := pickTokenOutAndRecipient(deltas, nodeClass, tokenIn, baseEvent.TxFromAddress, baseEvent.TxToAddress)

	// Post-selection normalization: if tokenOut is wrapped native but the sender received native, prefer native
	wrapped := models.ETHTokenAddress
	if req.Tx != nil && req.Tx.ChainId() != nil {
		if w, ok := models.ChainIDToWrappedNativeToken[models.ChainId(req.Tx.ChainId().Uint64())]; ok {
			wrapped = w
		}
	}
	if tokenOut == wrapped {
		if d, ok := deltas[baseEvent.TxFromAddress][models.ETHTokenAddress]; ok && d.Sign() > 0 {
			tokenOut = models.ETHTokenAddress
			recipient = baseEvent.TxFromAddress
		}
	}

	var amountIn *big.Int
	var amountOut *big.Int
	if (tokenIn == common.Address{}) || (tokenOut == common.Address{}) || (payer == common.Address{}) || (recipient == common.Address{}) {
		// Fallback: infer native-in swap if tx sender receives a token and a deposit_native occurred
		if fb, ok := fallbackNativeSwapInference(normalizedTransfers, deltas, baseEvent); ok {
			tokenIn = fb.tokenIn
			tokenOut = fb.tokenOut
			payer = fb.payer
			recipient = fb.recipient
			amountIn = fb.amountIn
			amountOut = fb.amountOut
			goto build_event
		}
		// Fallback: infer ERC20->native swap when sender sends ERC20 to pool-like and receives native
		if fb2, ok2 := fallbackNativeOutSwapInference(normalizedTransfers, deltas, baseEvent, nodeClass.LikelyPools); ok2 {
			tokenIn = fb2.tokenIn
			tokenOut = fb2.tokenOut
			payer = fb2.payer
			recipient = fb2.recipient
			amountIn = fb2.amountIn
			amountOut = fb2.amountOut
			goto build_event
		}
		return nil, fmt.Errorf("unable to infer swap endpoints (tokenIn/tokenOut/payer/recipient)")
	}

	// 5) Amounts are payer's negative delta for tokenIn, recipient's positive delta for tokenOut
	amountIn = new(big.Int).Abs(getDelta(deltas, payer, tokenIn))
	amountOut = new(big.Int).Set(getDelta(deltas, recipient, tokenOut))
	if amountIn.Sign() == 0 || amountOut.Sign() == 0 {
		return nil, fmt.Errorf("inferred zero amounts: in=%s out=%s", amountIn.String(), amountOut.String())
	}

build_event:
	// 6) Fetch decimals (fallback to 18 for native if needed)
	var baseDecimals *uint8
	var quoteDecimals *uint8
	if tokenIn != models.ETHTokenAddress {
		if d, derr := p.evmClient.GetDecimals(tokenIn); derr == nil {
			baseDecimals = &d
		} else {
		}
	}
	if tokenOut != models.ETHTokenAddress {
		if d, derr := p.evmClient.GetDecimals(tokenOut); derr == nil {
			quoteDecimals = &d
		} else {
		}
	}

	event := models.SwapEvent{
		BaseEvent:         baseEvent,
		BaseCoin:          tokenIn,
		QuoteCoin:         tokenOut,
		BaseCoinAmount:    amountIn,
		BaseCoinDecimals:  baseDecimals,
		QuoteCoinAmount:   amountOut,
		QuoteCoinDecimals: quoteDecimals,
		Sender:            payer,
		Receiver:          recipient,
	}

	return []models.Event{&event}, nil
}

// ===== Inference helpers =====

// normalizeTransfersForSwapInference filters out duplicated native <-> wrapped flows to avoid double counting.
//   - If a transfer is in native token (models.ETHTokenAddress) and the counterparty is the wrapped token contract,
//     skip it because a corresponding wrapped token event exists already.
func normalizeTransfersForSwapInference(transfers []models.ERC20TransferLog, req parsers.ParseTxRequest) []models.ERC20TransferLog {
	out := make([]models.ERC20TransferLog, 0, len(transfers))

	// Resolve wrapped native token for the chain; fallback to ETHTokenAddress if unknown
	wrapped := models.WETHTokenAddress
	if req.Tx != nil && req.Tx.ChainId() != nil {
		if w, ok := models.ChainIDToWrappedNativeToken[models.ChainId(req.Tx.ChainId().Uint64())]; ok {
			wrapped = w
		}
	}

	for _, t := range transfers {
		if t.TokenAddress == models.ETHTokenAddress {
			// Drop native token transfer when it is sent to/from the wrapped token contract address
			if t.ToAddress == wrapped || t.FromAddress == wrapped {
				continue
			}
		}
		out = append(out, t)
	}
	return out
}

// addressFlowStats captures which tokens an address received or sent during the tx.
type addressFlowStats struct {
	inTokens  map[common.Address]struct{}
	outTokens map[common.Address]struct{}
}

// computeAddressTokenDeltas builds per-address per-token deltas and simple in/out token sets.
func computeAddressTokenDeltas(transfers []models.ERC20TransferLog) (map[common.Address]map[common.Address]*big.Int, map[common.Address]*addressFlowStats) {
	deltas := make(map[common.Address]map[common.Address]*big.Int)
	stats := make(map[common.Address]*addressFlowStats)

	ensure := func(addr common.Address) {
		if _, ok := deltas[addr]; !ok {
			deltas[addr] = make(map[common.Address]*big.Int)
		}
		if _, ok := stats[addr]; !ok {
			stats[addr] = &addressFlowStats{inTokens: make(map[common.Address]struct{}), outTokens: make(map[common.Address]struct{})}
		}
	}

	for _, tr := range transfers {
		ensure(tr.FromAddress)
		ensure(tr.ToAddress)

		// From loses token
		if _, ok := deltas[tr.FromAddress][tr.TokenAddress]; !ok {
			deltas[tr.FromAddress][tr.TokenAddress] = new(big.Int)
		}
		deltas[tr.FromAddress][tr.TokenAddress].Sub(deltas[tr.FromAddress][tr.TokenAddress], tr.Value)
		stats[tr.FromAddress].outTokens[tr.TokenAddress] = struct{}{}

		// To gains token
		if _, ok := deltas[tr.ToAddress][tr.TokenAddress]; !ok {
			deltas[tr.ToAddress][tr.TokenAddress] = new(big.Int)
		}
		deltas[tr.ToAddress][tr.TokenAddress].Add(deltas[tr.ToAddress][tr.TokenAddress], tr.Value)
		stats[tr.ToAddress].inTokens[tr.TokenAddress] = struct{}{}
	}

	return deltas, stats
}

// NodeClassification keeps pools/routers sets; tweak heuristics or fill Known* sets to improve quality.
type NodeClassification struct {
	LikelyPools   map[common.Address]struct{}
	LikelyRouters map[common.Address]struct{}
}

// KnownPools and KnownRouters are placeholders for curated registries. Extend from callers or config in future.
var (
	KnownPools   = map[common.Address]struct{}{}
	KnownRouters = map[common.Address]struct{}{}
)

// classifyAddresses uses simple heuristics:
// - Pool-like: address that both receives at least one token and sends at least one different token
// - Router-like: currently empty unless populated in KnownRouters; kept separate for future logic
func classifyAddresses(deltas map[common.Address]map[common.Address]*big.Int, stats map[common.Address]*addressFlowStats) NodeClassification {
	pools := make(map[common.Address]struct{})
	routers := make(map[common.Address]struct{})

	// Seed with known sets
	for a := range KnownPools {
		pools[a] = struct{}{}
	}
	for a := range KnownRouters {
		routers[a] = struct{}{}
	}

	for addr, st := range stats {
		// Skip if already tagged
		if _, ok := pools[addr]; ok {
			continue
		}
		if _, ok := routers[addr]; ok {
			continue
		}
		// Heuristic: pool-like if it has inTokens and outTokens with at least two distinct tokens across directions
		if len(st.inTokens) > 0 && len(st.outTokens) > 0 {
			// Check if there exists a token in inTokens that is not in outTokens or vice versa
			different := false
			for tin := range st.inTokens {
				if _, ok := st.outTokens[tin]; !ok {
					different = true
					break
				}
			}
			if !different {
				for tout := range st.outTokens {
					if _, ok := st.inTokens[tout]; !ok {
						different = true
						break
					}
				}
			}
			if different {
				pools[addr] = struct{}{}
			}
		}
	}

	return NodeClassification{LikelyPools: pools, LikelyRouters: routers}
}

func isNonIntermediary(addr common.Address, cls NodeClassification) bool {
	if _, ok := cls.LikelyPools[addr]; ok {
		return false
	}
	if _, ok := cls.LikelyRouters[addr]; ok {
		return false
	}
	return true
}

// pickTokenInAndPayer selects the token with the largest total negative outflow among non-intermediaries
// and the address that contributed the largest negative delta for that token.
func pickTokenInAndPayer(deltas map[common.Address]map[common.Address]*big.Int, cls NodeClassification, txFrom common.Address) (common.Address, common.Address) {
	// Aggregate negative totals per token among non-intermediary addresses
	tokenNegTotals := make(map[common.Address]*big.Int)
	for addr, perToken := range deltas {
		if addr != txFrom {
			if !isNonIntermediary(addr, cls) {
				continue
			}
		}
		for token, delta := range perToken {
			if delta.Sign() < 0 {
				if _, ok := tokenNegTotals[token]; !ok {
					tokenNegTotals[token] = new(big.Int)
				}
				tokenNegTotals[token].Add(tokenNegTotals[token], new(big.Int).Abs(delta))
			}
		}
	}
	var tokenIn common.Address
	maxNeg := new(big.Int)
	for token, tot := range tokenNegTotals {
		if tot.Cmp(maxNeg) > 0 {
			tokenIn = token
			maxNeg = tot
		}
	}
	if (tokenIn == common.Address{}) {
		return common.Address{}, common.Address{}
	}
	// Payer: address with most negative delta on tokenIn
	var payer common.Address
	worst := new(big.Int)
	for addr, perToken := range deltas {
		if addr != txFrom {
			if !isNonIntermediary(addr, cls) {
				continue
			}
		}
		if d, ok := perToken[tokenIn]; ok && d.Sign() < 0 {
			abs := new(big.Int).Abs(d)
			if abs.Cmp(worst) > 0 {
				worst = abs
				payer = addr
			}
		}
	}
	// Prefer tx sender if it contributed negative flow of tokenIn
	if per, ok := deltas[txFrom]; ok {
		if d, ok2 := per[tokenIn]; ok2 && d.Sign() < 0 {
			payer = txFrom
		}
	}
	return tokenIn, payer
}

// pickTokenOutAndRecipient selects the token with the largest total positive inflow among non-intermediaries
// (excluding tokenIn) and the address with the largest positive delta for that token.
func pickTokenOutAndRecipient(deltas map[common.Address]map[common.Address]*big.Int, cls NodeClassification, tokenIn common.Address, txFrom, txTo common.Address) (common.Address, common.Address) {
	tokenPosTotals := make(map[common.Address]*big.Int)
	for addr, perToken := range deltas {
		if addr != txFrom && addr != txTo {
			if !isNonIntermediary(addr, cls) {
				continue
			}
		}
		for token, delta := range perToken {
			if token == tokenIn {
				continue
			}
			if delta.Sign() > 0 {
				if _, ok := tokenPosTotals[token]; !ok {
					tokenPosTotals[token] = new(big.Int)
				}
				tokenPosTotals[token].Add(tokenPosTotals[token], delta)
			}
		}
	}
	var tokenOut common.Address
	best := new(big.Int)
	for token, tot := range tokenPosTotals {
		if tot.Cmp(best) > 0 {
			tokenOut = token
			best = tot
		}
	}
	if (tokenOut == common.Address{}) {
		return common.Address{}, common.Address{}
	}
	var recipient common.Address
	most := new(big.Int)
	for addr, perToken := range deltas {
		if addr != txFrom && addr != txTo {
			if !isNonIntermediary(addr, cls) {
				continue
			}
		}
		if d, ok := perToken[tokenOut]; ok && d.Sign() > 0 {
			if d.Cmp(most) > 0 {
				most = d
				recipient = addr
			}
		}
	}
	// Prefer the tx sender as recipient if it ends up with positive tokenOut
	if per, ok := deltas[txFrom]; ok {
		if d, ok2 := per[tokenOut]; ok2 && d.Sign() > 0 {
			recipient = txFrom
		}
	}
	return tokenOut, recipient
}

func getDelta(deltas map[common.Address]map[common.Address]*big.Int, addr, token common.Address) *big.Int {
	if per, ok := deltas[addr]; ok {
		if d, ok2 := per[token]; ok2 {
			return d
		}
	}
	return new(big.Int)
}

// ===== Fallbacks =====

type fallbackResult struct {
	tokenIn   common.Address
	tokenOut  common.Address
	payer     common.Address
	recipient common.Address
	amountIn  *big.Int
	amountOut *big.Int
}

// fallbackNativeSwapInference tries to infer a simple native->ERC20 swap when:
// - there is a deposit_native (native -> wrapped) by a router/contract
// - an ERC20 is transferred to the tx sender in the same tx
// It chooses the largest positive ERC20 balance to the tx sender as tokenOut
// and uses tx value if amountIn cannot be read from deltas.
func fallbackNativeSwapInference(transfers []models.ERC20TransferLog, deltas map[common.Address]map[common.Address]*big.Int, base models.BaseEvent) (fallbackResult, bool) {
	// Detect any deposit_native or native->wrapped like transfer presence
	hasDeposit := false
	for _, t := range transfers {
		if t.Source == "deposit_native" {
			hasDeposit = true
			break
		}
	}
	if !hasDeposit {
		// Also allow presence of any native CALL in traces reflected as native transfers
		// but we only proceed if sender later receives some ERC20.
	}

	// Find the token with the largest positive balance for the tx sender (excluding native)
	var tokenOut common.Address
	most := new(big.Int)
	for token, d := range deltas[base.TxFromAddress] {
		if token == models.ETHTokenAddress {
			continue
		}
		if d.Sign() > 0 && d.Cmp(most) > 0 {
			tokenOut = token
			most = d
		}
	}
	if (tokenOut == common.Address{}) || most.Sign() == 0 {
		return fallbackResult{}, false
	}

	// Amount in: prefer negative native delta of tx sender; else tx value
	amountIn := new(big.Int)
	if d, ok := deltas[base.TxFromAddress][models.ETHTokenAddress]; ok && d.Sign() < 0 {
		amountIn = new(big.Int).Abs(d)
	} else if base.Value != nil {
		amountIn = new(big.Int).Set(base.Value)
	} else {
		return fallbackResult{}, false
	}

	return fallbackResult{
		tokenIn:   models.ETHTokenAddress,
		tokenOut:  tokenOut,
		payer:     base.TxFromAddress,
		recipient: base.TxFromAddress,
		amountIn:  amountIn,
		amountOut: new(big.Int).Set(most),
	}, true
}

func fallbackNativeOutSwapInference(transfers []models.ERC20TransferLog, deltas map[common.Address]map[common.Address]*big.Int, base models.BaseEvent, pools map[common.Address]struct{}) (fallbackResult, bool) {
	// Detect: sender sends some ERC20 to a pool-like address and later gets native back
	// tokenIn: ERC20 with largest negative delta for sender (non-native)
	var tokenIn common.Address
	maxNeg := new(big.Int)
	for tok, d := range deltas[base.TxFromAddress] {
		if tok == models.ETHTokenAddress {
			continue
		}
		if d.Sign() < 0 {
			abs := new(big.Int).Abs(d)
			if abs.Cmp(maxNeg) > 0 {
				maxNeg = abs
				tokenIn = tok
			}
		}
	}
	if (tokenIn == common.Address{}) || maxNeg.Sign() == 0 {
		return fallbackResult{}, false
	}
	// Require sender to have positive native delta (received native)
	if dn, ok := deltas[base.TxFromAddress][models.ETHTokenAddress]; !ok || dn.Sign() <= 0 {
		return fallbackResult{}, false
	}
	// Soft heuristic: ensure there exists a transfer of tokenIn from sender within this tx
	hasSenderOut := false
	for _, t := range transfers {
		if t.TokenAddress == tokenIn && t.FromAddress == base.TxFromAddress {
			hasSenderOut = true
			break
		}
	}
	if !hasSenderOut {
		return fallbackResult{}, false
	}
	return fallbackResult{
		tokenIn:   tokenIn,
		tokenOut:  models.ETHTokenAddress,
		payer:     base.TxFromAddress,
		recipient: base.TxFromAddress,
		amountIn:  new(big.Int).Set(maxNeg),
		amountOut: new(big.Int).Set(deltas[base.TxFromAddress][models.ETHTokenAddress]),
	}, true
}
