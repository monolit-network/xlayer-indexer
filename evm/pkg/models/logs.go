package models

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type ERC20TransferLog struct {
	FromAddress   common.Address
	ToAddress     common.Address
	TokenAddress  common.Address
	TokenDecimals uint8
	Value         *big.Int

	Source string
}

type SingleTraceResult struct {
	Type         string     `json:"type"`
	From         string     `json:"from"`
	To           string     `json:"to"`
	Value        string     `json:"value"`
	Gas          string     `json:"gas,omitempty"`
	GasUsed      string     `json:"gasUsed,omitempty"`
	Input        string     `json:"input,omitempty"`
	Output       string     `json:"output,omitempty"`
	Error        string     `json:"error,omitempty"`
	RevertReason string     `json:"revertReason,omitempty"`
	Logs         []TraceLog `json:"logs,omitempty"`
}

type FullTraceResult struct {
	SingleTraceResult
	Calls   []*FullTraceResult `json:"calls"`
	TxHash  common.Hash        `json:"-"`
	TxIndex uint               `json:"-"`
}

type TraceLogIndex uint64

func (i *TraceLogIndex) UnmarshalJSON(raw []byte) error {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	var text string
	if raw[0] == '"' {
		if err := json.Unmarshal(raw, &text); err != nil {
			return err
		}
	} else {
		text = string(raw)
	}

	var value uint64
	var err error
	if strings.HasPrefix(text, "0x") || strings.HasPrefix(text, "0X") {
		value, err = hexutil.DecodeUint64(text)
	} else {
		value, err = strconv.ParseUint(text, 10, 64)
	}
	if err != nil {
		return err
	}
	*i = TraceLogIndex(value)
	return nil
}

type TraceLog struct {
	Index    *TraceLogIndex `json:"index,omitempty"`
	Position *TraceLogIndex `json:"position,omitempty"`
	Address  common.Address `json:"address"`
	Topics   []common.Hash  `json:"topics"`
	Data     hexutil.Bytes  `json:"data"`
}

func (t *FullTraceResult) Flatten() []SingleTraceResult {
	results := []SingleTraceResult{t.SingleTraceResult}
	for _, call := range t.Calls {
		results = append(results, call.Flatten()...)
	}
	return results
}

func GetSenderAddress(tx *types.Transaction, chainId ChainId) (common.Address, error) {
	signer, err := types.Sender(types.LatestSignerForChainID(big.NewInt(int64(chainId))), tx)
	if err != nil {
		return common.Address{}, fmt.Errorf("error getting sender address: %w", err)
	}
	return signer, nil
}

func TryParseERC20TransferLog(log *types.Log) (ERC20TransferLog, error) {
	if len(log.Topics) < 3 {
		return ERC20TransferLog{}, fmt.Errorf("log has less than 3 topics")
	}

	if log.Topics[0] != ERC20TransferEventSelectorHash {
		return ERC20TransferLog{}, fmt.Errorf("log is not an ERC20 transfer log")
	}

	resLog := ERC20TransferLog{}

	err := ERC20ABI.UnpackIntoInterface(&resLog, "Transfer", log.Data)
	if err != nil {
		return ERC20TransferLog{}, fmt.Errorf("error unpacking log: %w", err)
	}

	resLog.FromAddress = common.HexToAddress(log.Topics[1].Hex())
	resLog.ToAddress = common.HexToAddress(log.Topics[2].Hex())
	resLog.TokenAddress = log.Address
	resLog.Source = "erc20"
	return resLog, nil
}

func TryParseWETHWithdrawLog(log *types.Log, chainId ChainId) (ERC20TransferLog, error) {
	if len(log.Topics) < 2 {
		return ERC20TransferLog{}, fmt.Errorf("log has less than 2 topics")
	}

	wrappedNativeToken, ok := ChainIDToWrappedNativeToken[chainId]
	if !ok {
		wrappedNativeToken = ETHTokenAddress
	}

	if log.Topics[0] != WETHWithdrawEventSelectorHash {
		return ERC20TransferLog{}, fmt.Errorf("log is not an WETH withdraw log")
	}

	tmp := struct {
		Wad *big.Int
	}{}

	err := WETHABI.UnpackIntoInterface(&tmp, "Withdraw", log.Data)
	if err != nil {
		return ERC20TransferLog{}, fmt.Errorf("error unpacking log: %w", err)
	}
	src := common.HexToAddress(log.Topics[1].Hex())

	outLog := ERC20TransferLog{
		FromAddress:  wrappedNativeToken,
		ToAddress:    src,
		TokenAddress: wrappedNativeToken,
		Value:        tmp.Wad,
		Source:       "withdraw_native",
	}
	return outLog, nil
}

func TryParseWETHDepositLog(log *types.Log, chainId ChainId) (ERC20TransferLog, error) {
	if len(log.Topics) < 2 {
		return ERC20TransferLog{}, fmt.Errorf("log has less than 2 topics")
	}

	if log.Topics[0] != WETHDepositEventSelectorHash {
		return ERC20TransferLog{}, fmt.Errorf("log is not an WETH deposit log")
	}

	wrappedNativeToken, ok := ChainIDToWrappedNativeToken[chainId]
	if !ok {
		wrappedNativeToken = ETHTokenAddress
	}

	tmp := struct {
		Wad *big.Int
	}{}

	err := WETHABI.UnpackIntoInterface(&tmp, "Deposit", log.Data)
	if err != nil {
		return ERC20TransferLog{}, fmt.Errorf("error unpacking log: %w", err)
	}
	dst := common.HexToAddress(log.Topics[1].Hex())

	inLog := ERC20TransferLog{
		FromAddress:  dst,
		ToAddress:    wrappedNativeToken,
		TokenAddress: wrappedNativeToken,
		Value:        tmp.Wad,
		Source:       "deposit_native",
	}
	return inLog, nil
}

func LogsContainsDepositOrWithdrawNative(logs []*types.Log) bool {
	for _, log := range logs {
		if len(log.Topics) < 2 {
			continue
		}
		if log.Topics[0] == WETHDepositEventSelectorHash || log.Topics[0] == WETHWithdrawEventSelectorHash {
			return true
		}
	}
	return false
}

func FindERC20Transfers(receipt *types.Receipt, tx *types.Transaction, chainId ChainId) ([]ERC20TransferLog, error) {
	transfers := []ERC20TransferLog{}

	// fromAddress, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	// if err != nil {
	// 	return nil, fmt.Errorf("error getting sender address: %w", err)
	// }

	// if tx.Value() != nil && tx.Value().Cmp(common.Big0) != 0 {
	// 	transfers = append(transfers, ERC20TransferLog{
	// 		FromAddress:  fromAddress,
	// 		ToAddress:    *tx.To(),
	// 		TokenAddress: ETHTokenAddress,
	// 		Value:        tx.Value(),
	// 	})
	// }

	wrappedNativeToken, ok := ChainIDToWrappedNativeToken[chainId]
	if !ok {
		wrappedNativeToken = ETHTokenAddress
	}

	for _, log := range receipt.Logs {
		if len(log.Topics) == 0 {
			continue
		}
		if transfer, err := TryParseERC20TransferLog(log); err == nil {
			transfers = append(transfers, transfer)
		} else if log.Address == wrappedNativeToken {
			if transfer, err := TryParseWETHWithdrawLog(log, chainId); err == nil {
				transfers = append(transfers, transfer)
			} else if transfer, err := TryParseWETHDepositLog(log, chainId); err == nil {
				transfers = append(transfers, transfer)
			}
		}
	}
	return transfers, nil
}

func ExtractETHTransfersFromTrace(traces []SingleTraceResult) ([]ERC20TransferLog, error) {
	var transfers []ERC20TransferLog

	for _, t := range traces {
		if strings.ToLower(t.Type) != "call" {
			continue
		}

		val := new(big.Int)
		val.SetString(t.Value[2:], 16) // Value is hex string "0x..."
		if val.Sign() == 0 {
			continue
		}

		from := common.HexToAddress(t.From)
		to := common.HexToAddress(t.To)

		transfers = append(transfers, ERC20TransferLog{
			FromAddress:  from,
			ToAddress:    to,
			TokenAddress: ETHTokenAddress,
			Value:        val,
			Source:       "native",
		})
	}

	return transfers, nil
}

func FindAllTransfers(receipt *types.Receipt, tx *types.Transaction, traces *FullTraceResult, chainId ChainId) ([]ERC20TransferLog, error) {
	transfers, err := FindERC20Transfers(receipt, tx, chainId)
	if err != nil {
		return nil, err
	}

	hasInitNativeTransfer := tx.To() != nil && tx.Value() != nil && tx.Value().Cmp(common.Big0) != 0
	if hasInitNativeTransfer {
		sender, err := GetSenderAddress(tx, chainId)
		if err != nil {
			return nil, err
		}
		transfers = append([]ERC20TransferLog{{
			FromAddress:  sender,
			ToAddress:    *tx.To(),
			TokenAddress: ETHTokenAddress,
			Value:        tx.Value(),
			Source:       "native",
		}}, transfers...)
	}

	if traces != nil {
		ethTransfers, err := ExtractETHTransfersFromTrace(traces.Flatten())
		if err != nil {
			return nil, err
		}

		if hasInitNativeTransfer {
			transfers = append(transfers, ethTransfers[1:]...)
		} else {
			transfers = append(transfers, ethTransfers...)
		}
	}

	return transfers, nil
}
