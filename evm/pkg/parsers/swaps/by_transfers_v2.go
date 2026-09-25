package swaps

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"go.uber.org/zap"
)

type settings struct {
	softPools bool
}

type ByTransfersParserV2Option func(*settings)

func WithSoftPools(softPools bool) ByTransfersParserV2Option {
	return func(s *settings) {
		s.softPools = softPools
	}
}

type ByTransfersParserV2 struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
	settings  *settings
}

func NewByTransfersParserV2(evmClient *evmclient.Client, logger *zap.Logger, options ...ByTransfersParserV2Option) *ByTransfersParserV2 {
	p := &ByTransfersParserV2{evmClient: evmClient, settings: &settings{}}
	p.logger = logger.Named(p.Name())
	for _, option := range options {
		option(p.settings)
	}
	return p
}

func (p ByTransfersParserV2) Name() string {
	return "swaps_by_transfers_parser_v2"
}

func (p *ByTransfersParserV2) ProducedEventsType() models.EventType {
	return models.EventTypeSwap
}

func (p *ByTransfersParserV2) RequiredData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesBlockHeader | evmclient.RequiredDataTypesTraces
}

func (p *ByTransfersParserV2) OptionalData() evmclient.RequiredDataTypes {
	return 0
}

func (p *ByTransfersParserV2) ParseTx(req parsers.ParseTxRequest) ([]models.Event, error) {
	events, err := p.parseTx(req, nil, nil, nil, nil)
	if err != nil && isTraceRecoverableError(err) {
		if sender, ok := p.recoverSenderFromTraces(req); ok {
			if events2, err2 := p.parseTx(req, &sender, nil, nil, nil); err2 == nil {
				return events2, nil
			}
		}
	}
	return events, err
}

// isTraceRecoverableError reports whether the failure is a sender/receiver
// attribution failure that call-tree descent may be able to fix.
func isTraceRecoverableError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "no suitable sender found") ||
		strings.Contains(msg, "sender has no net outflow tokens") ||
		strings.Contains(msg, "no suitable receiver candidates")
}

func (p *ByTransfersParserV2) parseTx(
	req parsers.ParseTxRequest,
	knownSender *common.Address,
	knownBaseToken *common.Address,
	knownReceiver *common.Address,
	knownQuoteToken *common.Address,
) ([]models.Event, error) {
	if err := parsers.ValidateParseTxRequest(req, p.RequiredData()); err != nil {
		return nil, err
	}
	baseEvent, err := parsers.MakeBaseEvent(req, p.evmClient)
	if err != nil {
		return nil, err
	}

	// this transfers contains both erc20 and native transfers in following order:
	// tx.value transfer (native, if value != 0)
	// erc20 transfers
	// all other native transfers

	// deposits/withdrawals of wrapped native token are included in the transfers as just erc20 transfers
	transfers, err := models.FindAllTransfers(req.Receipt, req.Tx, req.Traces, p.evmClient.ChainID)
	if err != nil {
		return nil, err
	}

	if len(transfers) < 2 {
		return nil, fmt.Errorf("no transfers found, must be at least 2")
	}

	// 1) Aggregate net inflow/outflow per address per token
	_, outf, net := aggregateFlows(transfers)

	// 3) Classify swap-like addresses (1 in, 1 out transfers) without knowing sender yet
	swapLike := p.classifySwapLike(net, transfers)
	swapLike[*req.Tx.To()] = struct{}{}
	delete(swapLike, baseEvent.TxFromAddress)
	p.logger.Debug("pools like", zap.Any("swap_like", swapLike))

	// 2) Determine sender (with fallback AFTER step 3 and excluding swap-like)
	var sender common.Address
	if knownSender != nil {
		// Sender recovered externally (robinhood call-tree descent); trust it.
		sender = *knownSender
	} else {
		sender, err = p.pickSender(baseEvent, outf, transfers, swapLike)
		if err != nil {
			return nil, err
		}
	}

	// Remove sender from swap-like per rule "except sender"
	delete(swapLike, sender)

	// 2.3) Determine base token and amount
	baseToken, baseAmount, err := p.pickBaseToken(sender, net, transfers)
	if err != nil {
		return nil, err
	}

	// 4) Determine receiver and quote token
	receiver, quoteToken, quoteAmount, err := p.pickReceiverAndQuote(net, sender, swapLike, transfers, baseToken)
	if err != nil {
		return nil, err
	}

	// Fetch decimals
	var baseDecimals *uint8
	var quoteDecimals *uint8
	if d, derr := p.evmClient.GetDecimals(baseToken); derr == nil {
		baseDecimals = &d
	} else {
		p.logger.Debug("get_decimals_failed", zap.String("token", baseToken.Hex()), zap.Error(derr))
	}

	if d, derr := p.evmClient.GetDecimals(quoteToken); derr == nil {
		quoteDecimals = &d
	} else {
		p.logger.Debug("get_decimals_failed", zap.String("token", quoteToken.Hex()), zap.Error(derr))
	}

	event := models.SwapEvent{
		BaseEvent:         baseEvent,
		BaseCoin:          baseToken,
		QuoteCoin:         quoteToken,
		BaseCoinAmount:    baseAmount,
		BaseCoinDecimals:  baseDecimals,
		QuoteCoinAmount:   quoteAmount,
		QuoteCoinDecimals: quoteDecimals,
		Sender:            sender,
		Receiver:          receiver,
	}

	result := []models.Event{&event}
	result = append(result, terminalTransferEvents(baseEvent, transfers, p.evmClient)...)
	return result, nil
}

// ===== Helpers =====

func aggregateFlows(transfers []models.ERC20TransferLog) (
	map[common.Address]map[common.Address]*big.Int,
	map[common.Address]map[common.Address]*big.Int,
	map[common.Address]map[common.Address]*big.Int,
) {
	inflow := make(map[common.Address]map[common.Address]*big.Int)
	outflow := make(map[common.Address]map[common.Address]*big.Int)
	net := make(map[common.Address]map[common.Address]*big.Int)

	addTo := func(m map[common.Address]map[common.Address]*big.Int, addr, token common.Address, amount *big.Int) {
		if _, ok := m[addr]; !ok {
			m[addr] = make(map[common.Address]*big.Int)
		}
		if _, ok := m[addr][token]; !ok {
			m[addr][token] = new(big.Int)
		}
		m[addr][token].Add(m[addr][token], amount)
	}

	for _, t := range transfers {
		if t.Value == nil || t.Value.Sign() == 0 {
			continue
		}
		addTo(outflow, t.FromAddress, t.TokenAddress, t.Value)
		addTo(inflow, t.ToAddress, t.TokenAddress, t.Value)
	}

	for addr, tin := range inflow {
		for token, inAmt := range tin {
			outAmt := new(big.Int)
			if touts, ok := outflow[addr]; ok {
				if v, ok2 := touts[token]; ok2 {
					outAmt = v
				}
			}
			delta := new(big.Int).Sub(inAmt, outAmt)
			if _, ok := net[addr]; !ok {
				net[addr] = make(map[common.Address]*big.Int)
			}
			net[addr][token] = delta
		}
	}
	for addr, tout := range outflow {
		for token, outAmt := range tout {
			if _, ok := inflow[addr]; ok {
				if _, ok2 := inflow[addr][token]; ok2 {
					continue
				}
			}
			delta := new(big.Int).Neg(outAmt)
			if _, ok := net[addr]; !ok {
				net[addr] = make(map[common.Address]*big.Int)
			}
			net[addr][token] = delta
		}
	}

	return inflow, outflow, net
}

func countPositive(net map[common.Address]map[common.Address]*big.Int, addr common.Address) int {
	count := 0
	for _, v := range net[addr] {
		if v.Sign() > 0 {
			count++
		}
	}
	return count
}

func countNegative(net map[common.Address]map[common.Address]*big.Int, addr common.Address) int {
	count := 0
	for _, v := range net[addr] {
		if v.Sign() < 0 {
			count++
		}
	}
	return count
}

// removed: getNegatives no longer used

func (p *ByTransfersParserV2) classifySwapLike(
	net map[common.Address]map[common.Address]*big.Int,
	transfers []models.ERC20TransferLog,
) map[common.Address]struct{} {
	res := make(map[common.Address]struct{})

	// Count number of incoming and outgoing transfers per address (by transfers count, not net tokens)
	inCount := make(map[common.Address]int)
	outCount := make(map[common.Address]int)
	for _, tr := range transfers {
		if tr.Value == nil || tr.Value.Sign() == 0 {
			continue
		}
		outCount[tr.FromAddress]++
		inCount[tr.ToAddress]++
	}

	for addr := range net {
		positive := countPositive(net, addr)
		negative := countNegative(net, addr)
		if p.settings.softPools {
			if positive != 0 && negative != 0 {
				res[addr] = struct{}{}
				continue
			}
		}

		if positive != 0 && negative == positive && inCount[addr] != 0 && outCount[addr] == inCount[addr] {
			res[addr] = struct{}{}
		}
	}
	return res
}

func (p *ByTransfersParserV2) pickSender(
	baseEvent models.BaseEvent,
	outflow map[common.Address]map[common.Address]*big.Int,
	transfers []models.ERC20TransferLog,
	swapLike map[common.Address]struct{},
) (common.Address, error) {
	// Prefer tx_from if it sent anything
	txFrom := baseEvent.TxFromAddress
	if m, ok := outflow[txFrom]; ok {
		for _, amt := range m {
			if amt.Sign() > 0 {
				return txFrom, nil
			}
		}
	}
	// Fallback AFTER step 3: choose sender of the first transfer that is not swap-like
	for _, tr := range transfers {
		if _, isSwap := swapLike[tr.FromAddress]; isSwap {
			continue
		}
		return tr.FromAddress, nil
	}
	return common.Address{}, fmt.Errorf("no suitable sender found (all candidates are swap-like)")
}

func (p *ByTransfersParserV2) pickBaseToken(
	sender common.Address,
	net map[common.Address]map[common.Address]*big.Int,
	transfers []models.ERC20TransferLog,
) (common.Address, *big.Int, error) {
	// Silence potential unused param warning if not referenced elsewhere
	_ = net

	// Recompute sender net ignoring burns (transfers to zero address)
	senderNet := make(map[common.Address]*big.Int)
	var zero common.Address
	for _, tr := range transfers {
		if tr.Value == nil || tr.Value.Sign() == 0 {
			continue
		}
		if tr.FromAddress == sender {
			if tr.ToAddress == zero {
				// Ignore burns
				continue
			}
			if _, ok := senderNet[tr.TokenAddress]; !ok {
				senderNet[tr.TokenAddress] = new(big.Int)
			}
			senderNet[tr.TokenAddress].Sub(senderNet[tr.TokenAddress], tr.Value)
		}
		if tr.ToAddress == sender {
			if _, ok := senderNet[tr.TokenAddress]; !ok {
				senderNet[tr.TokenAddress] = new(big.Int)
			}
			senderNet[tr.TokenAddress].Add(senderNet[tr.TokenAddress], tr.Value)
		}
	}

	// Count negatives and determine base token
	negatives := make([]common.Address, 0)
	for tok, d := range senderNet {
		if d.Sign() < 0 {
			negatives = append(negatives, tok)
		}
	}

	if len(negatives) == 1 {
		tok := negatives[0]
		amt := new(big.Int).Abs(senderNet[tok])
		return tok, amt, nil
	}
	if len(negatives) > 1 {
		negativesStrs := make([]string, 0, len(negatives))
		for _, tok := range negatives {
			negativesStrs = append(negativesStrs, tok.Hex())
		}
		return common.Address{}, nil, fmt.Errorf("multiple base tokens spent by sender: %s", strings.Join(negativesStrs, ", "))
	}
	return common.Address{}, nil, fmt.Errorf("sender has no net outflow tokens")
}

func (p *ByTransfersParserV2) pickReceiverAndQuote(
	net map[common.Address]map[common.Address]*big.Int,
	sender common.Address,
	swapLike map[common.Address]struct{},
	transfers []models.ERC20TransferLog,
	baseToken common.Address,
) (common.Address, common.Address, *big.Int, error) {
	wrappedAddr := p.evmClient.WrappedNativeTokenAddress
	// If sender received something -> receiver = sender
	if countPositive(net, sender) > 0 {
		// Prefer token received by sender from a router/pool (address in swapLike)
		// Aggregate amounts per token for transfers where From in swapLike and To == sender
		perToken := make(map[common.Address]*big.Int)
		for _, tr := range transfers {
			if tr.Value == nil || tr.Value.Sign() == 0 {
				continue
			}
			if tr.ToAddress == sender {
				if _, isSwap := swapLike[tr.FromAddress]; isSwap {
					if tr.TokenAddress == baseToken || ((tr.TokenAddress == models.ETHTokenAddress && baseToken == wrappedAddr) ||
						(tr.TokenAddress == wrappedAddr && baseToken == models.ETHTokenAddress)) {
						// Ignore base token returned from router/pool to sender
						continue
					}
					if _, ok := perToken[tr.TokenAddress]; !ok {
						perToken[tr.TokenAddress] = new(big.Int)
					}
					perToken[tr.TokenAddress].Add(perToken[tr.TokenAddress], tr.Value)
				}
			}
		}

		if len(perToken) == 0 {
			// Fallback: pick the largest single transfer to sender (excluding base token)
			var bestAmt *big.Int
			var bestToken common.Address
			for _, tr := range transfers {
				if tr.Value == nil || tr.Value.Sign() == 0 {
					continue
				}
				if tr.ToAddress == sender {
					if tr.TokenAddress == baseToken {
						// Ignore base token returned to sender
						continue
					}
					if bestAmt == nil || tr.Value.Cmp(bestAmt) > 0 {
						bestAmt = tr.Value
						bestToken = tr.TokenAddress
					}
				}
			}
			if bestAmt != nil {
				return sender, bestToken, new(big.Int).Set(bestAmt), nil
			}
			p.logger.Debug("sender received no tokens from router/pool", zap.String("sender", sender.Hex()), zap.Any("pools", swapLike))
			return common.Address{}, common.Address{}, nil, fmt.Errorf("sender received no tokens to infer quote token: sender: %s", sender.Hex())
		}

		if len(perToken) > 1 {
			return common.Address{}, common.Address{}, nil, fmt.Errorf("multiple tokens received from router/pool: sender: %s", sender.Hex())
		}

		for tok, amt := range perToken {
			return sender, tok, amt, nil
		}
	}

	// Else: pick address with exactly 1 net inflow token, excluding swap-like and sender; pick the largest inflow
	var bestAddr common.Address
	bestAmt := new(big.Int)
	var quoteToken common.Address
	for addr := range net {
		if addr == sender {
			continue
		}
		if _, isSwap := swapLike[addr]; isSwap {
			continue
		}
		if countPositive(net, addr) == 1 {
			for tok, d := range net[addr] {
				if d.Sign() > 0 {
					if bestAddr == (common.Address{}) || d.Cmp(bestAmt) > 0 {
						bestAddr = addr
						bestAmt = new(big.Int).Set(d)
						quoteToken = tok
					}
				}
			}
		}
	}
	if bestAddr == (common.Address{}) {
		// Last-resort fallback: custodial swaps where the pool-chain output stays
		// inside the router (tx.to) or another swap-like contract. Accept such a
		// candidate only if it has exactly one net-inflow token AND that token
		// differs from the base token (the sold token accumulating in a pool is
		// not proceeds). Receiver then records WHERE the proceeds are custodied;
		// sender/PnL semantics are unchanged.
		for addr := range net {
			if addr == sender {
				continue
			}
			if countPositive(net, addr) != 1 {
				continue
			}
			for tok, d := range net[addr] {
				if d.Sign() > 0 {
					if tok == baseToken || ((tok == models.ETHTokenAddress && baseToken == wrappedAddr) ||
						(tok == wrappedAddr && baseToken == models.ETHTokenAddress)) {
						continue
					}
					if bestAddr == (common.Address{}) || d.Cmp(bestAmt) > 0 {
						bestAddr = addr
						bestAmt = new(big.Int).Set(d)
						quoteToken = tok
					}
				}
			}
		}
	}
	if bestAddr == (common.Address{}) {
		return common.Address{}, common.Address{}, nil, fmt.Errorf("no suitable receiver candidates with single net inflow")
	}
	return bestAddr, quoteToken, bestAmt, nil
}

// recoverSenderFromTraces (robinhood-only) handles AA/ERC-4337/relayer/
// contract-mediated swaps where tx.from (bundler/relayer) holds none of the
// swapped tokens, so flat net-flow attribution fails. It descends the call
// tree to the frame(s) that call into a Uniswap pool (an address that emitted
// a Swap event in this receipt) and walks the caller chain upwards from the
// pool, picking the deepest caller that actually has net token outflow in
// this tx — i.e. the real swap initiator sitting below the relayer/EntryPoint.
// Attribution is only returned when it is unambiguous (exactly one candidate
// across all swap frames); otherwise we do nothing and the tx stays an error.
func (p *ByTransfersParserV2) recoverSenderFromTraces(req parsers.ParseTxRequest) (common.Address, bool) {
	if req.Traces == nil || req.Receipt == nil || req.Tx == nil {
		return common.Address{}, false
	}

	// Addresses that emitted a Uniswap Swap event = pools / v4 PoolManager.
	swapEmitters := make(map[common.Address]struct{})
	for _, lg := range req.Receipt.Logs {
		if len(lg.Topics) == 0 {
			continue
		}
		if lg.Topics[0] == models.UniswapV2SwapEventTopic ||
			lg.Topics[0] == models.UniswapV3SwapEventTopic ||
			lg.Topics[0] == models.UniswapV4SwapEventTopic {
			swapEmitters[lg.Address] = struct{}{}
		}
	}
	if len(swapEmitters) == 0 {
		return common.Address{}, false
	}

	transfers, err := models.FindAllTransfers(req.Receipt, req.Tx, req.Traces, p.evmClient.ChainID)
	if err != nil || len(transfers) < 2 {
		return common.Address{}, false
	}
	_, _, net := aggregateFlows(transfers)

	candidates := make(map[common.Address]struct{})

	var walk func(frame *models.FullTraceResult, callers []common.Address)
	walk = func(frame *models.FullTraceResult, callers []common.Address) {
		if frame == nil {
			return
		}
		callee := common.HexToAddress(frame.To)
		if _, isEmitter := swapEmitters[callee]; isEmitter {
			// Caller chain from root down to the immediate caller of the pool.
			chain := make([]common.Address, 0, len(callers)+1)
			chain = append(chain, callers...)
			chain = append(chain, common.HexToAddress(frame.From))
			// Deepest-first: the most specific account below relayers/EntryPoint
			// that actually spends a token is the swap initiator.
			for i := len(chain) - 1; i >= 0; i-- {
				addr := chain[i]
				if _, e := swapEmitters[addr]; e {
					continue
				}
				if countNegative(net, addr) >= 1 {
					candidates[addr] = struct{}{}
					break
				}
			}
		}
		next := make([]common.Address, len(callers), len(callers)+1)
		copy(next, callers)
		next = append(next, common.HexToAddress(frame.From))
		for _, sub := range frame.Calls {
			walk(sub, next)
		}
	}
	walk(req.Traces, nil)

	if len(candidates) != 1 {
		return common.Address{}, false
	}
	for addr := range candidates {
		p.logger.Debug("trace_sender_recovered",
			zap.String("sender", addr.Hex()),
			zap.String("tx", req.Tx.Hash().Hex()))
		return addr, true
	}
	return common.Address{}, false
}
