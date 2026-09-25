package defi

import (
	"fmt"
	"math/big"
	"sort"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type DefiParser struct {
	evmClient *evmclient.Client
	logger    *zap.Logger
}

func NewDefiParser(evmClient *evmclient.Client, logger *zap.Logger) *DefiParser {
	return &DefiParser{evmClient: evmClient, logger: logger}
}

func (p DefiParser) Name() string {
	return "defi_parser"
}

func (p *DefiParser) ProducedEventsType() models.EventType {
	return models.EventTypeDefi
}

func (p *DefiParser) RequiredData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesBlockHeader
}

func (p *DefiParser) OptionalData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesTraces
}

func (p *DefiParser) ParseTx(req parsers.ParseTxRequest) ([]models.Event, error) {
	if err := parsers.ValidateParseTxRequest(req, p.RequiredData()); err != nil {
		return nil, err
	}

	if req.Traces == nil {
		if (req.Tx.Value() != nil && req.Tx.Value().Cmp(common.Big0) != 0) ||
			models.LogsContainsDepositOrWithdrawNative(req.Receipt.Logs) {
			return nil, fmt.Errorf("native transfer or deposit/withdraw native log found")
		}
	}

	baseEvent, err := parsers.MakeBaseEvent(req, p.evmClient)
	if err != nil {
		return nil, err
	}

	transfers, err := models.FindAllTransfers(req.Receipt, req.Tx, req.Traces, p.evmClient.ChainID)
	if err != nil {
		return nil, err
	}

	inputs := make(map[common.Address]*big.Int)  // transfers to user
	outputs := make(map[common.Address]*big.Int) // transfers from user
	for _, transfer := range transfers {
		if transfer.FromAddress == baseEvent.TxFromAddress {
			if _, ok := outputs[transfer.TokenAddress]; !ok {
				outputs[transfer.TokenAddress] = big.NewInt(0)
				inputs[transfer.TokenAddress] = big.NewInt(0)
			}
			outputs[transfer.TokenAddress].Add(outputs[transfer.TokenAddress], transfer.Value)
			inputs[transfer.TokenAddress].Sub(inputs[transfer.TokenAddress], transfer.Value)
		}
		if transfer.ToAddress == baseEvent.TxFromAddress {
			if _, ok := inputs[transfer.TokenAddress]; !ok {
				outputs[transfer.TokenAddress] = big.NewInt(0)
				inputs[transfer.TokenAddress] = big.NewInt(0)
			}
			inputs[transfer.TokenAddress].Add(inputs[transfer.TokenAddress], transfer.Value)
			outputs[transfer.TokenAddress].Sub(outputs[transfer.TokenAddress], transfer.Value)
		}
	}

	defiEvent := models.DefiEvent{
		BaseEvent: baseEvent,
	}

	// set 2 biggest positive inputs/outputs to event
	type tokenAmt struct {
		addr   common.Address
		amount *big.Int
	}

	positiveInputs := make([]tokenAmt, 0)
	for addr, amt := range inputs {
		if amt != nil && amt.Sign() > 0 {
			positiveInputs = append(positiveInputs, tokenAmt{addr: addr, amount: new(big.Int).Set(amt)})
		}
	}
	positiveOutputs := make([]tokenAmt, 0)
	for addr, amt := range outputs {
		if amt != nil && amt.Sign() > 0 {
			positiveOutputs = append(positiveOutputs, tokenAmt{addr: addr, amount: new(big.Int).Set(amt)})
		}
	}

	sort.Slice(positiveInputs, func(i, j int) bool { return positiveInputs[i].amount.Cmp(positiveInputs[j].amount) > 0 })
	sort.Slice(positiveOutputs, func(i, j int) bool { return positiveOutputs[i].amount.Cmp(positiveOutputs[j].amount) > 0 })

	getDecimals := func(addr common.Address) *uint8 {
		var dec *uint8
		if d, derr := p.evmClient.GetDecimals(addr); derr == nil {
			dec = &d
		} else {
			dec = nil
		}
		return dec
	}

	makeDescriptions := func(tokens []tokenAmt) []models.TokenTransferDescription {
		if len(tokens) == 0 {
			return nil
		}
		descs := make([]models.TokenTransferDescription, 0, len(tokens))
		for _, token := range tokens {
			dec := getDecimals(token.addr)
			descs = append(descs, models.TokenTransferDescription{
				TokenAddress: token.addr,
				Amount:       new(big.Int).Set(token.amount),
				Decimals:     dec,
			})
		}
		return descs
	}

	defiEvent.Inputs = makeDescriptions(positiveInputs)
	defiEvent.Outputs = makeDescriptions(positiveOutputs)

	return []models.Event{&defiEvent}, nil
}
