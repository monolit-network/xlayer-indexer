package swaps

import (
	"fmt"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

// ERC4337Identifier matches any tx whose receipt contains a
// UserOperationEvent — an ERC-4337 bundle. Matched by topic only, so it is
// agnostic to EntryPoint addresses and versions (same philosophy as the
// uniswap/pons identifiers). It is registered as an EXCLUSIVE priority
// identifier: the whole-tx by_transfers path must not run for bundles,
// otherwise N independent user operations get collapsed into one bogus
// net-swap attributed to a bundler/pool.
type ERC4337Identifier struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewERC4337Identifier(evmClient *evmclient.Client, logger *zap.Logger) *ERC4337Identifier {
	return &ERC4337Identifier{evmClient: evmClient, logger: logger.Named("erc4337_identifier")}
}

func (i *ERC4337Identifier) CheckTx(req parsers.ParseTxRequest) (bool, error) {
	if req.Receipt == nil {
		return false, nil
	}
	for _, lg := range req.Receipt.Logs {
		if lg != nil && len(lg.Topics) > 0 && lg.Topics[0] == models.UserOperationEventTopic {
			return true, nil
		}
	}
	return false, nil
}

func (i *ERC4337Identifier) Source() string        { return "erc4337" }
func (i *ERC4337Identifier) ShouldStop() bool      { return true }
func (i *ERC4337Identifier) GetParserName() string { return "erc4337_bundle_parser" }

// ERC4337BundleParser slices a bundle into per-user-operation segments by
// log position (ops execute sequentially; each UserOperationEvent closes
// its window) and runs the EXISTING by_transfers heuristic on every
// segment with the op sender pinned. The receiver is then forced to the op
// sender as well: proceeds of 4337 sells physically land on app settlement
// vaults that keep an internal balance per user, so the user IS the
// economic receiver. Curve segments are skipped — the pons parser (a
// non-exclusive priority identifier) already emits those per event.
type ERC4337BundleParser struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
	bt        *ByTransfersParserV2
}

func NewERC4337BundleParser(evmClient *evmclient.Client, logger *zap.Logger) *ERC4337BundleParser {
	return &ERC4337BundleParser{
		evmClient: evmClient,
		logger:    logger.Named("erc4337_bundle_parser"),
		bt:        NewByTransfersParserV2(evmClient, logger),
	}
}

func (p *ERC4337BundleParser) Name() string { return "erc4337_bundle_parser" }

func (p *ERC4337BundleParser) ProducedEventsType() models.EventType { return models.EventTypeSwap }

func (p *ERC4337BundleParser) RequiredData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesBlockHeader | evmclient.RequiredDataTypesTraces
}

func (p *ERC4337BundleParser) OptionalData() evmclient.RequiredDataTypes { return 0 }

func (p *ERC4337BundleParser) ParseTx(req parsers.ParseTxRequest) ([]models.Event, error) {
	if err := parsers.ValidateParseTxRequest(req, p.RequiredData()); err != nil {
		return nil, err
	}
	baseEvent, err := parsers.MakeBaseEvent(req, p.evmClient)
	if err != nil {
		return nil, err
	}

	type opWindow struct {
		endPos int
		sender common.Address
	}
	var ops []opWindow
	for pos, lg := range req.Receipt.Logs {
		if lg == nil || len(lg.Topics) < 3 || lg.Topics[0] != models.UserOperationEventTopic {
			continue
		}
		ops = append(ops, opWindow{endPos: pos, sender: common.BytesToAddress(lg.Topics[2].Bytes())})
	}
	if len(ops) == 0 {
		return nil, fmt.Errorf("erc4337: no UserOperationEvent in tx %s", req.Tx.Hash().Hex())
	}

	events := make([]models.Event, 0, len(ops))
	prev := -1
	for _, op := range ops {
		segment := req.Receipt.Logs[prev+1 : op.endPos]
		prev = op.endPos
		if len(segment) == 0 {
			continue
		}
		hasCurve := false
		protocol := ""
		for _, lg := range segment {
			if lg == nil || len(lg.Topics) == 0 {
				continue
			}
			switch lg.Topics[0] {
			case models.UniswapV4SwapEventTopic:
				protocol = "uniswap_v4"
			case models.UniswapV3SwapEventTopic:
				if protocol == "" {
					protocol = "uniswap_v3"
				}
			case models.UniswapV2SwapEventTopic:
				if protocol == "" {
					protocol = "uniswap_v2"
				}
			}
		}
		if hasCurve || protocol == "" {
			// curve ops are emitted by the pons parser; non-swap ops
			// (plain transfers, approvals) carry no swap to record.
			continue
		}

		segReceipt := *req.Receipt
		segReceipt.Logs = segment
		segReq := req
		segReq.Receipt = &segReceipt
		segReq.Source = protocol

		opSender := op.sender
		segEvents, segErr := p.bt.parseTx(segReq, &opSender, nil, nil, nil)
		if segErr != nil {
			p.logger.Debug("erc4337_segment_unparsed",
				zap.String("tx", req.Tx.Hash().Hex()),
				zap.String("op_sender", opSender.Hex()),
				zap.Error(segErr))
			continue
		}
		for _, ev := range segEvents {
			swap, ok := ev.(*models.SwapEvent)
			if !ok {
				continue // per-segment transfer emission is skipped; done once below
			}
			swap.Sender = opSender
			swap.Receiver = opSender
			swap.BaseEvent.Source = protocol
			events = append(events, swap)
		}
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("erc4337: no swap segments parsed in tx %s", req.Tx.Hash().Hex())
	}

	transfers, terr := models.FindAllTransfers(req.Receipt, req.Tx, req.Traces, p.evmClient.ChainID)
	if terr == nil {
		events = append(events, terminalTransferEvents(baseEvent, transfers, p.evmClient)...)
	}
	return events, nil
}
