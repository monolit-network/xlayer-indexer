package polymarketProcessor

import (
	"errors"
	"math/big"
	"testing"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/abis/polymarketConditionalTokens"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers/polymarket"
	sharedmodels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPreProcessLogGroupFiltersByParseRangeRequirements(t *testing.T) {
	processor := &Processor{}
	metadata := map[string]any{}

	validAddressFiltered := testPolymarketLog(models.PolymarketOrderFilledEventSelectorHash, models.PolymarketCTFExchangeAddress)
	validV2AddressFiltered := testPolymarketLog(models.PolymarketOrderFilledEventV2SelectorHash, models.PolymarketCTFExchangeAddressV2)
	validNegRiskV2AddressFiltered := testPolymarketLog(models.PolymarketOrderFilledEventV2SelectorHash, models.PolymarketNegRiskCTFExchangeAddressV2)
	invalidAddressFiltered := testPolymarketLog(models.PolymarketOrderFilledEventSelectorHash, common.HexToAddress("0x00000000000000000000000000000000000000aa"))
	validNoAddressTopic := testPolymarketLog(models.PolymarketFeeRefundedEventSelectorHash, common.HexToAddress("0x00000000000000000000000000000000000000bb"))
	unrelatedTopic := testPolymarketLog(models.ERC20TransferEventSelectorHash, models.PolymarketCTFExchangeAddress)
	noTopics := &types.Log{Address: models.PolymarketCTFExchangeAddress}

	filtered, err := processor.preProcessLogGroup([]*types.Log{
		validAddressFiltered,
		validV2AddressFiltered,
		validNegRiskV2AddressFiltered,
		invalidAddressFiltered,
		validNoAddressTopic,
		unrelatedTopic,
		noTopics,
	}, metadata)
	require.NoError(t, err)

	require.Len(t, filtered, 4)
	require.Same(t, validNoAddressTopic, filtered[0])
	require.Same(t, validAddressFiltered, filtered[1])
	require.Same(t, validV2AddressFiltered, filtered[2])
	require.Same(t, validNegRiskV2AddressFiltered, filtered[3])
}

func TestProcessLogsFromTxPassesFilteredLogsToHandlers(t *testing.T) {
	validLog := testPolymarketLog(models.PolymarketOrderFilledEventSelectorHash, models.PolymarketCTFExchangeAddress)
	priorityLog := testPolymarketLog(models.PolymarketFeeRefundedEventSelectorHash, models.PolymarketConditionalTokensAddress)
	invalidLog := testPolymarketLog(models.PolymarketOrderFilledEventSelectorHash, common.HexToAddress("0x00000000000000000000000000000000000000aa"))
	noTopics := &types.Log{Address: models.PolymarketCTFExchangeAddress}

	var seenTxLogs []*types.Log
	processor := &Processor{
		logger:       zap.NewNop(),
		processStats: map[string]*processTiming{},
		processors: map[string]logProcessor{
			models.PolymarketOrderFilledEventSelectorHash.Hex(): {
				name: "test_order_filled",
				fn: func(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
					seenTxLogs = append([]*types.Log(nil), allLogs...)
					require.Same(t, validLog, log)
					return nil
				},
			},
		},
	}

	err := processor.ProcessLogsFromTx([]*types.Log{invalidLog, noTopics, validLog, priorityLog})
	require.NoError(t, err)

	require.Len(t, seenTxLogs, 2)
	require.Same(t, priorityLog, seenTxLogs[0])
	require.Same(t, validLog, seenTxLogs[1])
}

func TestFindOrderTakerUsesV2OrdersMatched(t *testing.T) {
	taker := common.HexToAddress("0x00000000000000000000000000000000000000ff")
	log := &types.Log{
		Address: models.PolymarketCTFExchangeAddressV2,
		Topics: []common.Hash{
			models.PolymarketOrdersMatchedEventV2SelectorHash,
			common.HexToHash("0x01"),
			common.BytesToHash(taker.Bytes()),
		},
	}

	got := (&Processor{}).findOrderTaker([]*types.Log{log}, map[string]any{})

	require.Equal(t, taker, got)
}

func TestFindCollateralAdapterTransferUserChoosesClosestExpectedAmount(t *testing.T) {
	user := common.HexToAddress("0x00000000000000000000000000000000000000ab")
	closerUser := common.HexToAddress("0x00000000000000000000000000000000000000ac")
	adapter := models.PolymarketCtfCollateralAdapterAddress
	tokenID := big.NewInt(55)
	expectedAmount := big.NewInt(100)
	processor := &Processor{
		polymarketParser: polymarket.NewCTFParser(nil, zap.NewNop()),
		logger:           zap.NewNop(),
	}

	got, ok, err := processor.findCollateralAdapterTransferUser(
		adapter,
		[]*types.Log{
			testConditionalTokensTransferLog(t, models.PolymarketTransferSingleEventSelectorHash, adapter, user, 1, tokenID, big.NewInt(80)),
			testConditionalTokensTransferLog(t, models.PolymarketTransferSingleEventSelectorHash, adapter, closerUser, 2, tokenID, big.NewInt(103)),
		},
		map[string]any{},
		[]collateralAdapterExpectedTransfer{{TokenID: tokenID, Amount: expectedAmount}},
	)

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, closerUser, got)
}

func TestFindCollateralAdapterTransferUserSupportsBatchToAdapter(t *testing.T) {
	user := common.HexToAddress("0x00000000000000000000000000000000000000cd")
	adapter := models.PolymarketNegRiskCtfCollateralAdapterAddress
	tokenID := big.NewInt(77)
	amount := big.NewInt(200)
	processor := &Processor{
		polymarketParser: polymarket.NewCTFParser(nil, zap.NewNop()),
		logger:           zap.NewNop(),
	}

	got, ok, err := processor.findCollateralAdapterTransferUser(
		adapter,
		[]*types.Log{testConditionalTokensTransferLog(t, models.PolymarketTransferBatchEventSelectorHash, user, adapter, 9, tokenID, amount)},
		map[string]any{},
		[]collateralAdapterExpectedTransfer{{TokenID: tokenID, Amount: amount}},
	)

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, user, got)
}

func TestFindCollateralAdapterTransferUserIgnoresNonConditionalTokensTransfers(t *testing.T) {
	user := common.HexToAddress("0x00000000000000000000000000000000000000ef")
	adapter := models.PolymarketCtfCollateralAdapterAddress
	tokenID := big.NewInt(88)
	amount := big.NewInt(300)
	processor := &Processor{
		polymarketParser: polymarket.NewCTFParser(nil, zap.NewNop()),
		logger:           zap.NewNop(),
	}
	usdcLikeLog := testConditionalTokensTransferLog(t, models.PolymarketTransferSingleEventSelectorHash, user, adapter, 9, tokenID, amount)
	usdcLikeLog.Address = common.HexToAddress("0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174")

	_, ok, err := processor.findCollateralAdapterTransferUser(
		adapter,
		[]*types.Log{usdcLikeLog},
		map[string]any{},
		[]collateralAdapterExpectedTransfer{{TokenID: tokenID, Amount: amount}},
	)

	require.NoError(t, err)
	require.False(t, ok)
}

func TestFindCollateralAdapterTransferUserIgnoresPolymarketContracts(t *testing.T) {
	user := common.HexToAddress("0x00000000000000000000000000000000000000de")
	adapter := models.PolymarketCtfCollateralAdapterAddress
	tokenID := big.NewInt(99)
	expectedAmount := big.NewInt(100)
	processor := &Processor{
		polymarketParser: polymarket.NewCTFParser(nil, zap.NewNop()),
		logger:           zap.NewNop(),
	}

	got, ok, err := processor.findCollateralAdapterTransferUser(
		adapter,
		[]*types.Log{
			testConditionalTokensTransferLog(t, models.PolymarketTransferSingleEventSelectorHash, adapter, models.PolymarketCTFExchangeAddress, 1, tokenID, expectedAmount),
			testConditionalTokensTransferLog(t, models.PolymarketTransferSingleEventSelectorHash, adapter, user, 2, tokenID, big.NewInt(101)),
		},
		map[string]any{},
		[]collateralAdapterExpectedTransfer{{TokenID: tokenID, Amount: expectedAmount}},
	)

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, user, got)
}

func TestFindPositionRedeemBurnAmountsUsesOriginalArrayOrderAndStopsAtPreviousRedeem(t *testing.T) {
	redeemer := common.HexToAddress("0x00000000000000000000000000000000000000ab")
	tokenID := big.NewInt(55)
	processor := &Processor{
		polymarketParser: polymarket.NewCTFParser(nil, zap.NewNop()),
		logger:           zap.NewNop(),
	}
	previousRedeem := testPolymarketLog(models.PolymarketPayoutRedemptionEventSelectorHash, models.PolymarketConditionalTokensAddress)
	currentRedeem := testPolymarketLog(models.PolymarketPayoutRedemptionEventSelectorHash, models.PolymarketConditionalTokensAddress)

	amounts, err := processor.findPositionRedeemBurnAmounts(
		currentRedeem,
		[]*types.Log{
			testConditionalTokensTransferLog(t, models.PolymarketTransferSingleEventSelectorHash, redeemer, models.ZeroAddress, 100, tokenID, big.NewInt(999)),
			previousRedeem,
			testConditionalTokensTransferLog(t, models.PolymarketTransferSingleEventSelectorHash, redeemer, models.ZeroAddress, 1, tokenID, big.NewInt(123)),
			currentRedeem,
		},
		&sharedmodels.PolymarketPayoutRedemption{
			Redeemer: redeemer,
			TokenIDs: []*big.Int{tokenID},
		},
	)

	require.NoError(t, err)
	require.Len(t, amounts, 1)
	require.Equal(t, "123", amounts[0].String())
}

func TestFindPositionRedeemBurnAmountsSkipsPayoutMintBeforeRedeem(t *testing.T) {
	redeemer := common.HexToAddress("0x00000000000000000000000000000000000000bc")
	parentTokenID := big.NewInt(77)
	redeemedTokenID := big.NewInt(88)
	processor := &Processor{
		polymarketParser: polymarket.NewCTFParser(nil, zap.NewNop()),
		logger:           zap.NewNop(),
	}
	currentRedeem := testPolymarketLog(models.PolymarketPayoutRedemptionEventSelectorHash, models.PolymarketConditionalTokensAddress)

	amounts, err := processor.findPositionRedeemBurnAmounts(
		currentRedeem,
		[]*types.Log{
			testConditionalTokensTransferLog(t, models.PolymarketTransferSingleEventSelectorHash, redeemer, models.ZeroAddress, 1, redeemedTokenID, big.NewInt(321)),
			testConditionalTokensTransferLog(t, models.PolymarketTransferSingleEventSelectorHash, models.ZeroAddress, redeemer, 2, parentTokenID, big.NewInt(111)),
			currentRedeem,
		},
		&sharedmodels.PolymarketPayoutRedemption{
			Redeemer: redeemer,
			TokenIDs: []*big.Int{redeemedTokenID},
		},
	)

	require.NoError(t, err)
	require.Len(t, amounts, 1)
	require.Equal(t, "321", amounts[0].String())
}

func TestRetryGetMarketDataCallRetriesRetryableErrors(t *testing.T) {
	t.Parallel()

	want := polymarket.MarketData{QuestionCount: 2, FeeBips: 15}
	attempts := 0

	got, err := retryGetMarketDataCall(zap.NewNop(), func() (polymarket.MarketData, error) {
		attempts++
		if attempts == 1 {
			return polymarket.MarketData{}, errors.New(`Post "http://rpc": dial tcp 127.0.0.1:8545: connect: connection refused`)
		}
		return want, nil
	})

	require.NoError(t, err)
	require.Equal(t, 2, attempts)
	require.Equal(t, want, got)
}

func TestRetryGetMarketDataCallDoesNotRetryNonRetryableErrors(t *testing.T) {
	t.Parallel()

	attempts := 0

	_, err := retryGetMarketDataCall(zap.NewNop(), func() (polymarket.MarketData, error) {
		attempts++
		return polymarket.MarketData{}, errors.New("market not found")
	})

	require.EqualError(t, err, "market not found")
	require.Equal(t, 1, attempts)
}

func testPolymarketLog(topic common.Hash, address common.Address) *types.Log {
	return &types.Log{
		Address: address,
		Topics:  []common.Hash{topic},
	}
}

func testConditionalTokensTransferLog(t *testing.T, topic common.Hash, from common.Address, to common.Address, index uint, tokenID *big.Int, amount *big.Int) *types.Log {
	t.Helper()

	contractABI, err := polymarketConditionalTokens.PolymarketConditionalTokensMetaData.GetAbi()
	require.NoError(t, err)

	var data []byte
	switch topic {
	case models.PolymarketTransferSingleEventSelectorHash:
		data, err = contractABI.Events["TransferSingle"].Inputs.NonIndexed().Pack(tokenID, amount)
	case models.PolymarketTransferBatchEventSelectorHash:
		data, err = contractABI.Events["TransferBatch"].Inputs.NonIndexed().Pack([]*big.Int{tokenID}, []*big.Int{amount})
	default:
		t.Fatalf("unexpected transfer topic: %s", topic.Hex())
	}
	require.NoError(t, err)

	return &types.Log{
		Address: models.PolymarketConditionalTokensAddress,
		Index:   index,
		Topics: []common.Hash{
			topic,
			common.HexToHash("0x01"),
			common.BytesToHash(from.Bytes()),
			common.BytesToHash(to.Bytes()),
		},
		Data: data,
	}
}
