package models

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var UmaOptimisticOracleV2InvalidPrice = new(big.Int).Lsh(big.NewInt(-1), 255)

func SelectorFromSignature(selector string) string {
	return fmt.Sprintf("%x", crypto.Keccak256([]byte(selector))[:4])
}

func SelectorFromTx(tx *types.Transaction) string {
	if len(tx.Data()) < 4 {
		return ""
	}
	return hex.EncodeToString(tx.Data()[:4])
}

func UnpackCalldata(parsedABI abi.ABI, methodName string, calldata []byte, out interface{}) error {
	if len(calldata) < 4 {
		return errors.New("calldata too short")
	}

	params := calldata[4:]

	var method *abi.Method
	if methodFound, ok := parsedABI.Methods[methodName]; ok {
		method = &methodFound
	} else {
		return fmt.Errorf("method %s not found in ABI", methodName)
	}

	args, err := method.Inputs.Unpack(params)
	if err != nil {
		return err
	}

	return method.Inputs.Copy(out, args)
}

// ABIs
var (
	ERC20ABI, _          = abi.JSON(strings.NewReader(erc20ABIRaw))
	WETHABI, _           = abi.JSON(strings.NewReader(wethABIRaw))
	CurveSwapEventABI, _ = abi.JSON(strings.NewReader(curveSwapEventABIRaw))

	UniswapV2SwapFuncABI, _  = abi.JSON(strings.NewReader(uniswapV2SwapFuncABIRaw))
	UniswapV2SwapEventABI, _ = abi.JSON(strings.NewReader(uniswapV2SwapEventABIRaw))

	UniswapV3SwapFuncABI, _  = abi.JSON(strings.NewReader(uniswapV3SwapFuncABIRaw))
	UniswapV3SwapEventABI, _ = abi.JSON(strings.NewReader(uniswapV3SwapEventABIRaw))

	UniswapUniversalRouterExecuteFuncABI, _ = abi.JSON(strings.NewReader(uniswapUniversalRouterExecuteFuncABIRaw))

	PolymarketOrderFilledEventABI, _ = abi.JSON(strings.NewReader(polymarketOrderFilledEventABIRaw))
)

// Event selectors
var (
	ERC20TransferEventSelectorHash = common.BytesToHash(crypto.Keccak256([]byte("Transfer(address,address,uint256)")))
	WETHWithdrawEventSelectorHash  = common.BytesToHash(crypto.Keccak256([]byte("Withdraw(address,uint256)")))
	WETHDepositEventSelectorHash   = common.BytesToHash(crypto.Keccak256([]byte("Deposit(address,uint256)")))

	CurveSwapEventSelectorHash = common.BytesToHash(crypto.Keccak256([]byte("Exchange(address,address,address[11],uint256[5][5],address[5],uint256,uint256)")))

	// CTF Exchange
	PolymarketOrderFilledEventSelectorHash        = common.BytesToHash(crypto.Keccak256([]byte("OrderFilled(bytes32,address,address,uint256,uint256,uint256,uint256,uint256)")))
	PolymarketOrderFilledEventV2SelectorHash      = common.BytesToHash(crypto.Keccak256([]byte("OrderFilled(bytes32,address,address,uint8,uint256,uint256,uint256,uint256,bytes32,bytes32)")))
	PolymarketOrdersMatchedEventV2SelectorHash    = common.BytesToHash(crypto.Keccak256([]byte("OrdersMatched(bytes32,address,uint8,uint256,uint256,uint256)")))
	PolymarketFeeRefundedEventSelectorHash        = common.BytesToHash(crypto.Keccak256([]byte("FeeRefunded(bytes32,address,uint256,uint256,uint256)")))
	PolymarketNegRiskFeeRefundedEventSelectorHash = common.BytesToHash(crypto.Keccak256([]byte("FeeRefunded(address,address,uint256,uint256)")))

	// NegRisk adapter

	PolymarketNegRiskMarketPreparedEventSelectorHash     = common.BytesToHash(crypto.Keccak256([]byte("MarketPrepared(bytes32,address,uint256,bytes)")))
	PolymarketNegRiskPositionsConvertedEventSelectorHash = common.BytesToHash(crypto.Keccak256([]byte("PositionsConverted(address,bytes32,uint256,uint256)")))
	PolymarketNegRiskPayoutRedemptionEventSelectorHash   = common.BytesToHash(crypto.Keccak256([]byte("PayoutRedemption(address,bytes32,uint256[],uint256)")))

	// Conditional Tokens
	PolymarketTransferSingleEventSelectorHash       = common.BytesToHash(crypto.Keccak256([]byte("TransferSingle(address,address,address,uint256,uint256)")))
	PolymarketTransferBatchEventSelectorHash        = common.BytesToHash(crypto.Keccak256([]byte("TransferBatch(address,address,address,uint256[],uint256[])")))
	PolymarketConditionPreparationEventSelectorHash = common.BytesToHash(crypto.Keccak256([]byte("ConditionPreparation(bytes32,address,bytes32,uint256)")))
	PolymarketConditionResolutionEventSelectorHash  = common.BytesToHash(crypto.Keccak256([]byte("ConditionResolution(bytes32,address,bytes32,uint256,uint256[])")))
	PolymarketPositionSplitEventSelectorHash        = common.BytesToHash(crypto.Keccak256([]byte("PositionSplit(address,address,bytes32,bytes32,uint256[],uint256)")))
	PolymarketNegRiskPositionSplitEventSelectorHash = common.BytesToHash(crypto.Keccak256([]byte("PositionSplit(address,bytes32,uint256)")))
	PolymarketPositionMergeEventSelectorHash        = common.BytesToHash(crypto.Keccak256([]byte("PositionsMerge(address,address,bytes32,bytes32,uint256[],uint256)")))
	PolymarketNegRiskPositionMergeEventSelectorHash = common.BytesToHash(crypto.Keccak256([]byte("PositionsMerge(address,bytes32,uint256)")))
	PolymarketPayoutRedemptionEventSelectorHash     = common.BytesToHash(crypto.Keccak256([]byte("PayoutRedemption(address,address,bytes32,bytes32,uint256[],uint256)")))

	// AMM
	PolymarketFPMMFundingAddedEventSelectorHash   = common.BytesToHash(crypto.Keccak256([]byte("FPMMFundingAdded(address,uint256[],uint256)")))
	PolymarketFPMMFundingRemovedEventSelectorHash = common.BytesToHash(crypto.Keccak256([]byte("FPMMFundingRemoved(address,uint256[],uint256,uint256)")))
	PolymarketFPMMBuyEventSelectorHash            = common.BytesToHash(crypto.Keccak256([]byte("FPMMBuy(address,uint256,uint256,uint256,uint256)")))
	PolymarketFPMMSellEventSelectorHash           = common.BytesToHash(crypto.Keccak256([]byte("FPMMSell(address,uint256,uint256,uint256,uint256)")))
	PolymarketAMMTransferEventSelectorHash        = common.BytesToHash(crypto.Keccak256([]byte("Transfer(address,address,uint256)")))
	PolymarketAMMApprovalEventSelectorHash        = common.BytesToHash(crypto.Keccak256([]byte("Approval(address,address,uint256)")))

	// UMA CTF Adapter
	PolymarketUmaCtfAdapterQuestionInitializedEventSelectorHash = common.BytesToHash(crypto.Keccak256([]byte("QuestionInitialized(bytes32,uint256,address,bytes,address,uint256,uint256)")))
	PolymarketUmaCtfAdapterQuestionResetEventSelectorHash       = common.BytesToHash(crypto.Keccak256([]byte("QuestionReset(bytes32)")))
	PolymarketUmaCtfAdapterQuestionResolvedEventSelectorHash    = common.BytesToHash(crypto.Keccak256([]byte("QuestionResolved(bytes32,int256,uint256[])")))
	PolymarketUmaCtfAdapterQuestionPausedEventSelectorHash      = common.BytesToHash(crypto.Keccak256([]byte("QuestionPaused(bytes32)")))
	PolymarketUmaCtfAdapterQuestionFlaggedEventSelectorHash     = common.BytesToHash(crypto.Keccak256([]byte("QuestionFlagged(bytes32)")))

	// Uma Optimistic Oracle
	UmaOptimisticOracleV2ProposePriceEventSelectorHash = common.BytesToHash(crypto.Keccak256([]byte("ProposePrice(address,address,bytes32,uint256,bytes,int256,uint256,address)")))
	UmaOptimisticOracleV2DisputePriceEventSelectorHash = common.BytesToHash(crypto.Keccak256([]byte("DisputePrice(address,address,address,bytes32,uint256,bytes,int256)")))
	UmaOptimisticOracleV2SettleEventSelectorHash       = common.BytesToHash(crypto.Keccak256([]byte("Settle(address,address,address,bytes32,uint256,bytes,int256,uint256)")))
	UmaOptimisticOracleV2RequestPriceEventSelectorHash = common.BytesToHash(crypto.Keccak256([]byte("RequestPrice(address,bytes32,uint256,bytes,address,uint256,uint256)")))
)

// Function selectors
var (
	UniswapV2SwapFuncSelectorHash = SelectorFromSignature("swap(uint256,uint256,address,bytes)")
	UniswapV3SwapFuncSelectorHash = SelectorFromSignature("swap(address,bool,int256,uint160,bytes)")
)

var (
	UniswapV1EthToTokenSwapInput  = SelectorFromSignature("ethToTokenSwapInput(uint256,uint256)")
	UniswapV1EthToTokenSwapOutput = SelectorFromSignature("ethToTokenSwapOutput(uint256,uint256)")

	UniswapV1TokenToEthSwapInput  = SelectorFromSignature("tokenToEthSwapInput(uint256,uint256,uint256)")
	UniswapV1TokenToEthSwapOutput = SelectorFromSignature("tokenToEthSwapOutput(uint256,uint256,uint256)")

	UniswapV1TokenToTokenSwapInput  = SelectorFromSignature("tokenToTokenSwapInput(uint256,uint256,uint256,uint256,address)")
	UniswapV1TokenToTokenSwapOutput = SelectorFromSignature("tokenToTokenSwapOutput(uint256,uint256,uint256,uint256,address)")

	UniswapV1SwapSelectors = map[string]struct{}{
		UniswapV1EthToTokenSwapInput:    {},
		UniswapV1EthToTokenSwapOutput:   {},
		UniswapV1TokenToEthSwapInput:    {},
		UniswapV1TokenToEthSwapOutput:   {},
		UniswapV1TokenToTokenSwapInput:  {},
		UniswapV1TokenToTokenSwapOutput: {},
	}
)

var (
	CurveSwapSelector1 = SelectorFromSignature("exchange(address[11],uint256[5][5],uint256,uint256)")
	CurveSwapSelector2 = SelectorFromSignature("exchange(address[11],uint256[5][5],uint256,uint256,address[5])")
	CurveSwapSelector3 = SelectorFromSignature("exchange(address[11],uint256[5][5],uint256,uint256,address[5],address)")
	CurveSwapSelectors = map[string]struct{}{
		CurveSwapSelector1: {},
		CurveSwapSelector2: {},
		CurveSwapSelector3: {},
	}
)

// ABIs raw
const (
	// base
	erc20ABIRaw = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`
	wethABIRaw  = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"src","type":"address"},{"indexed":false,"name":"wad","type":"uint256"}],"name":"Withdraw","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"dst","type":"address"},{"indexed":false,"name":"wad","type":"uint256"}],"name":"Deposit","type":"event"}]`

	// curve
	curveSwapEventABIRaw = `[{"name":"Exchange","inputs":[{"name":"sender","type":"address","indexed":true},{"name":"receiver","type":"address","indexed":true},{"name":"route","type":"address[11]","indexed":false},{"name":"swap_params","type":"uint256[5][5]","indexed":false},{"name":"pools","type":"address[5]","indexed":false},{"name":"in_amount","type":"uint256","indexed":false},{"name":"out_amount","type":"uint256","indexed":false}],"anonymous":false,"type":"event"}]`

	// uniswap
	uniswapV2SwapFuncABIRaw  = `[{"constant":false,"inputs":[{"internalType":"uint256","name":"amount0Out","type":"uint256"},{"internalType":"uint256","name":"amount1Out","type":"uint256"},{"internalType":"address","name":"to","type":"address"},{"internalType":"bytes","name":"data","type":"bytes"}],"name":"swap","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	uniswapV2SwapEventABIRaw = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount0In","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1In","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount0Out","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1Out","type":"uint256"},{"indexed":true,"internalType":"address","name":"to","type":"address"}],"name":"Swap","type":"event"}]`

	uniswapV3SwapFuncABIRaw  = `[{"inputs":[{"internalType":"address","name":"recipient","type":"address"},{"internalType":"bool","name":"zeroForOne","type":"bool"},{"internalType":"int256","name":"amountSpecified","type":"int256"},{"internalType":"uint160","name":"sqrtPriceLimitX96","type":"uint160"},{"internalType":"bytes","name":"data","type":"bytes"}],"name":"swap","outputs":[{"internalType":"int256","name":"amount0","type":"int256"},{"internalType":"int256","name":"amount1","type":"int256"}],"stateMutability":"nonpayable","type":"function"}]`
	uniswapV3SwapEventABIRaw = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":true,"internalType":"address","name":"recipient","type":"address"},{"indexed":false,"internalType":"int256","name":"amount0","type":"int256"},{"indexed":false,"internalType":"int256","name":"amount1","type":"int256"},{"indexed":false,"internalType":"uint160","name":"sqrtPriceX96","type":"uint160"},{"indexed":false,"internalType":"uint128","name":"liquidity","type":"uint128"},{"indexed":false,"internalType":"int24","name":"tick","type":"int24"}],"name":"Swap","type":"event"}]`

	uniswapUniversalRouterExecuteFuncABIRaw = `[{"inputs":[{"internalType":"bytes","name":"commands","type":"bytes"},{"internalType":"bytes[]","name":"inputs","type":"bytes[]"}],"name":"execute","outputs":[],"stateMutability":"payable","type":"function"}]`

	// CTF Exchange
	polymarketOrderFilledEventABIRaw = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"orderHash","type":"bytes32"},{"indexed":true,"internalType":"address","name":"maker","type":"address"},{"indexed":true,"internalType":"address","name":"taker","type":"address"},{"indexed":false,"internalType":"uint256","name":"makerAssetId","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"takerAssetId","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"makerAmountFilled","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"takerAmountFilled","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"fee","type":"uint256"}],"name":"OrderFilled","type":"event"}]`

	// Conditional Tokens
	polymarketConditionPreparationEventABIRaw = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"conditionId","type":"bytes32"},{"indexed":true,"name":"oracle","type":"address"},{"indexed":true,"name":"questionId","type":"bytes32"},{"indexed":false,"name":"outcomeSlotCount","type":"uint256"}],"name":"ConditionPreparation","type":"event"}]`
	polymarketConditionResolutionEventABIRaw  = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"conditionId","type":"bytes32"},{"indexed":true,"name":"oracle","type":"address"},{"indexed":true,"name":"questionId","type":"bytes32"},{"indexed":false,"name":"outcomeSlotCount","type":"uint256"},{"indexed":false,"name":"payoutNumerators","type":"uint256[]"}],"name":"ConditionResolution","type":"event"}]`

	polymarketPositionSplitEventABIRaw    = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"stakeholder","type":"address"},{"indexed":false,"name":"collateralToken","type":"address"},{"indexed":true,"name":"parentCollectionId","type":"bytes32"},{"indexed":true,"name":"conditionId","type":"bytes32"},{"indexed":false,"name":"partition","type":"uint256[]"},{"indexed":false,"name":"amount","type":"uint256"}],"name":"PositionSplit","type":"event"}]`
	polymarketMergePositionsEventABIRaw   = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"stakeholder","type":"address"},{"indexed":false,"name":"collateralToken","type":"address"},{"indexed":true,"name":"parentCollectionId","type":"bytes32"},{"indexed":true,"name":"conditionId","type":"bytes32"},{"indexed":false,"name":"partition","type":"uint256[]"},{"indexed":false,"name":"amount","type":"uint256"}],"name":"PositionsMerge","type":"event"}]`
	polymarketPayoutRedemptionEventABIRaw = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"redeemer","type":"address"},{"indexed":true,"name":"collateralToken","type":"address"},{"indexed":true,"name":"parentCollectionId","type":"bytes32"},{"indexed":false,"name":"conditionId","type":"bytes32"},{"indexed":false,"name":"indexSets","type":"uint256[]"},{"indexed":false,"name":"payout","type":"uint256"}],"name":"PayoutRedemption","type":"event"}]`

	polymarketTransferSingleEventABIRaw = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"id","type":"uint256"},{"indexed":false,"name":"value","type":"uint256"}],"name":"TransferSingle","type":"event"}]`
	polymarketTransferBatchEventABIRaw  = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"ids","type":"uint256[]"},{"indexed":false,"name":"values","type":"uint256[]"}],"name":"TransferBatch","type":"event"}]`
)
