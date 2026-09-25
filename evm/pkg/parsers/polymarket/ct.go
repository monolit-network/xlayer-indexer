package polymarket

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/abis/polymarketCTFExchange"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/abis/polymarketCTFExchangeV2"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/abis/polymarketConditionalTokens"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/abis/polymarketFPMM"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/abis/polymarketFeeModule"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/abis/polymarketNegRiskAdapter"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/abis/polymarketNegRiskFeeModule"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/abis/polymarketUmaCtfAdapterV2"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/abis/umaOptimisticOracleV2"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/abis/umaOptimisticOracleV3"
	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	sharedModels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/ethereum/go-ethereum/common"

	"go.uber.org/zap"
)

type CTFParser struct {
	evmClient       *evmclient.Client
	logger          *zap.Logger
	tokenIdCache    *lru.Cache[string, *big.Int]
	marketDataCache *lru.Cache[common.Hash, MarketData]
}

func NewCTFParser(evmClient *evmclient.Client, logger *zap.Logger) *CTFParser {
	cache, _ := lru.New[string, *big.Int](65536)
	marketDataCache, _ := lru.New[common.Hash, MarketData](65536)
	return &CTFParser{
		evmClient:       evmClient,
		logger:          logger,
		tokenIdCache:    cache,
		marketDataCache: marketDataCache,
	}
}

func (p *CTFParser) Name() string {
	return ordersParserName
}

func (p *CTFParser) ProducedEventsType() models.EventType {
	return models.EventTypePolymarketOrder
}

func (p *CTFParser) RequiredData() evmclient.RequiredDataTypes {
	return evmclient.RequiredDataTypesReceipt
}

func (p *CTFParser) OptionalData() evmclient.RequiredDataTypes {
	return 0
}

func (p *CTFParser) ParseTx(req parsers.ParseTxRequest) ([]models.Event, error) {
	if err := parsers.ValidateParseTxRequest(req, p.RequiredData()); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("not implemented")
}

var (
	conditionalTokensContract, _          = polymarketConditionalTokens.NewPolymarketConditionalTokens(models.PolymarketConditionalTokensAddress, nil)
	polymarketCTFExchangeContract, _      = polymarketCTFExchange.NewPolymarketCTFExchange(models.PolymarketCTFExchangeAddress, nil)
	polymarketCTFExchangeV2Contract, _    = polymarketCTFExchangeV2.NewPolymarketCTFExchangeV2(models.PolymarketCTFExchangeAddressV2, nil)
	polymarketNegRiskAdapterContract, _   = polymarketNegRiskAdapter.NewPolymarketNegRiskAdapter(models.PolymarketNegRiskAdapterAddress, nil)
	polymarketUmaCtfAdapterContract, _    = polymarketUmaCtfAdapterV2.NewPolymarketUmaCtfAdapterV2(models.PolymarketUmaCtfAdapterV2Address, nil)
	umaOptimisticOracleV2Contract, _      = umaOptimisticOracleV2.NewUmaOptimisticOracleV2(models.UmaOptimisticOracleV2Address, nil)
	umaOptimisticOracleV3Contract, _      = umaOptimisticOracleV3.NewUmaOptimisticOracleV3(models.UmaOptimisticOracleV3Address, nil)
	polymarketFeeModuleContract, _        = polymarketFeeModule.NewPolymarketFeeModule(models.PolymarketFeeModuleAddress, nil)
	polymarketNegRiskFeeModuleContract, _ = polymarketNegRiskFeeModule.NewPolymarketNegRiskFeeModule(models.PolymarketNegRiskFeeModuleAddress, nil)
	polymarketFPMMContract, _             = polymarketFPMM.NewPolymarketFPMM(common.Address{}, nil)
)

func (p *CTFParser) IsPolymarketContractAddress(address common.Address) bool {
	return p.IsPolymarketFeeModuleAddress(address) ||
		p.IsPolymarketCTFExchangeAddress(address) ||
		p.IsPolymarketConditionalTokensAddress(address) ||
		p.IsPolymarketAdapterAddress(address)
}

func (p *CTFParser) IsPolymarketFeeModuleAddress(address common.Address) bool {
	return address == models.PolymarketFeeModuleAddress || address == models.PolymarketNegRiskFeeModuleAddress
}

func (p *CTFParser) IsPolymarketCTFExchangeAddress(address common.Address) bool {
	return address == models.PolymarketCTFExchangeAddress ||
		address == models.PolymarketCTFExchangeAddressV2 ||
		address == models.PolymarketNegRiskCTFExchangeAddress ||
		address == models.PolymarketNegRiskCTFExchangeAddressV2
}

func (p *CTFParser) IsPolymarketAdapterAddress(address common.Address) bool {
	_, ok := models.PolymarketCTFAdaptersAddressesMap[address]
	return ok
}

func (p *CTFParser) IsPolymarketConditionalTokensAddress(address common.Address) bool {
	return address == models.PolymarketConditionalTokensAddress
}

func (p *CTFParser) ParseConditionPreparationLog(log *types.Log) (*sharedModels.PolymarketMarketOnchain, error) {
	if len(log.Topics) > 0 && log.Topics[0] != models.PolymarketConditionPreparationEventSelectorHash {
		return nil, nil
	}
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	conditionPreparation, err := conditionalTokensContract.ParseConditionPreparation(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing condition preparation log: %w", err)
	}

	outcomeSlotCount := conditionPreparation.OutcomeSlotCount.Int64()
	conditionId := common.BytesToHash(conditionPreparation.ConditionId[:])
	questionId := common.BytesToHash(conditionPreparation.QuestionId[:])
	return &sharedModels.PolymarketMarketOnchain{
		ConditionID:   conditionId,
		Oracle:        conditionPreparation.Oracle,
		QuestionID:    questionId,
		OutcomesCount: outcomeSlotCount,
	}, nil
}

func (p *CTFParser) ParseConditionResolutionLog(log *types.Log) (*sharedModels.PolymarketResolutionWinners, error) {
	if len(log.Topics) > 0 && log.Topics[0] != models.PolymarketConditionResolutionEventSelectorHash {
		return nil, nil
	}
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	conditionResolution, err := conditionalTokensContract.ParseConditionResolution(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing condition resolution log: %w", err)
	}
	conditionId := common.BytesToHash(conditionResolution.ConditionId[:])

	return &sharedModels.PolymarketResolutionWinners{
		Oracle:      conditionResolution.Oracle,
		ConditionID: conditionId,
		Numerators:  conditionResolution.PayoutNumerators,
	}, nil
}

func (p *CTFParser) ParsePositionSplitLog(log *types.Log) (*sharedModels.PolymarketPositionSplitMerge, error) {
	if len(log.Topics) > 0 && log.Topics[0] != models.PolymarketPositionSplitEventSelectorHash {
		return nil, nil
	}
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	positionSplit, err := conditionalTokensContract.ParsePositionSplit(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing position split log: %w", err)
	}

	tokenIds := make([]*big.Int, len(positionSplit.Partition))
	for i, partition := range positionSplit.Partition {
		tokenId := p.getTokenIdPartitionCached(positionSplit.ConditionId, partition, positionSplit.ParentCollectionId, positionSplit.CollateralToken)
		if tokenId == nil {
			return nil, fmt.Errorf("error getting token id: %w", err)
		}
		tokenIds[i] = tokenId
	}

	return &sharedModels.PolymarketPositionSplitMerge{
		Stakeholder:        positionSplit.Stakeholder,
		ConditionID:        positionSplit.ConditionId,
		TokenIDs:           tokenIds,
		CollateralToken:    positionSplit.CollateralToken,
		ParentCollectionId: positionSplit.ParentCollectionId,
		Partition:          positionSplit.Partition,
		Amount:             positionSplit.Amount,
		IsSplit:            true,
	}, nil
}

func (p *CTFParser) ParsePositionMergeLog(log *types.Log) (*sharedModels.PolymarketPositionSplitMerge, error) {
	if len(log.Topics) > 0 && log.Topics[0] != models.PolymarketPositionMergeEventSelectorHash {
		return nil, nil
	}
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	positionMerge, err := conditionalTokensContract.ParsePositionsMerge(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing position merge log: %w", err)
	}

	tokenIds := make([]*big.Int, len(positionMerge.Partition))
	for i, partition := range positionMerge.Partition {
		tokenId := p.getTokenIdPartitionCached(positionMerge.ConditionId, partition, positionMerge.ParentCollectionId, positionMerge.CollateralToken)
		if tokenId == nil {
			return nil, fmt.Errorf("error getting token id: %w", err)
		}
		tokenIds[i] = tokenId
	}

	return &sharedModels.PolymarketPositionSplitMerge{
		Stakeholder:        positionMerge.Stakeholder,
		ConditionID:        positionMerge.ConditionId,
		TokenIDs:           tokenIds,
		CollateralToken:    positionMerge.CollateralToken,
		ParentCollectionId: positionMerge.ParentCollectionId,
		Partition:          positionMerge.Partition,
		Amount:             positionMerge.Amount,
		IsSplit:            false,
	}, nil
}

func (p *CTFParser) ParsePayoutRedemptionLog(log *types.Log) (*sharedModels.PolymarketPayoutRedemption, error) {
	if len(log.Topics) > 0 && log.Topics[0] != models.PolymarketPayoutRedemptionEventSelectorHash {
		return nil, nil
	}
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	payoutRedemption, err := conditionalTokensContract.ParsePayoutRedemption(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing payout redemption log: %w", err)
	}

	tokenIDs := make([]*big.Int, 0, len(payoutRedemption.IndexSets))
	parentCollectionID := common.BytesToHash(payoutRedemption.ParentCollectionId[:])
	for _, indexSet := range payoutRedemption.IndexSets {
		tokenIDs = append(tokenIDs, p.getTokenIdPartitionCached(common.BytesToHash(payoutRedemption.ConditionId[:]), indexSet, parentCollectionID, payoutRedemption.CollateralToken))
	}

	return &sharedModels.PolymarketPayoutRedemption{
		Redeemer:        payoutRedemption.Redeemer,
		ConditionID:     common.BytesToHash(payoutRedemption.ConditionId[:]),
		TokenIDs:        tokenIDs,
		Amounts:         nil,
		CollateralToken: payoutRedemption.CollateralToken,
		Payout:          payoutRedemption.Payout,
	}, nil
}

func (p *CTFParser) ParseNegRiskPayoutRedemptionLog(log *types.Log) (*sharedModels.PolymarketPayoutRedemption, error) {
	if len(log.Topics) > 0 && log.Topics[0] != models.PolymarketNegRiskPayoutRedemptionEventSelectorHash {
		return nil, nil
	}
	if len(log.Topics) != 3 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}
	collateralToken := models.PolymarketNegRiskWrappedCollateralAddress

	payoutRedemption, err := polymarketNegRiskAdapterContract.ParsePayoutRedemption(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing payout redemption log: %w", err)
	}

	tokenIds := make([]*big.Int, 0, 2)
	indexSets := []*big.Int{big.NewInt(1), big.NewInt(2)}
	for _, indexSet := range indexSets {
		tokenIds = append(tokenIds, p.getTokenIdPartitionCached(common.BytesToHash(payoutRedemption.ConditionId[:]), indexSet, common.Hash{}, collateralToken))
	}

	return &sharedModels.PolymarketPayoutRedemption{
		Redeemer:        payoutRedemption.Redeemer,
		ConditionID:     common.BytesToHash(payoutRedemption.ConditionId[:]),
		TokenIDs:        tokenIds,
		Amounts:         payoutRedemption.Amounts,
		CollateralToken: models.PolymarketNegRiskWrappedCollateralAddress,
		Payout:          payoutRedemption.Payout,
	}, nil
}

func (p *CTFParser) ParseOrderFilledLog(log *types.Log) (*sharedModels.PolymarketOrderFilled, error) {
	if len(log.Topics) > 0 &&
		log.Topics[0] != models.PolymarketOrderFilledEventSelectorHash &&
		log.Topics[0] != models.PolymarketOrderFilledEventV2SelectorHash {
		return nil, nil
	}

	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	if log.Topics[0] == models.PolymarketOrderFilledEventV2SelectorHash {
		return parseOrderFilledV2Log(log)
	}

	orderFilled, err := polymarketCTFExchangeContract.ParseOrderFilled(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing order filled log: %w", err)
	}

	return &sharedModels.PolymarketOrderFilled{
		OrderHash:         orderFilled.OrderHash,
		Maker:             orderFilled.Maker,
		Taker:             orderFilled.Taker,
		MakerAssetID:      orderFilled.MakerAssetId,
		TakerAssetID:      orderFilled.TakerAssetId,
		MakerAmountFilled: orderFilled.MakerAmountFilled,
		TakerAmountFilled: orderFilled.TakerAmountFilled,
		Fee:               orderFilled.Fee,
	}, nil
}

func parseOrderFilledV2Log(log *types.Log) (*sharedModels.PolymarketOrderFilled, error) {
	orderFilled, err := polymarketCTFExchangeV2Contract.ParseOrderFilled(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing order filled v2 log: %w", err)
	}

	makerAssetID, takerAssetID, err := polymarketV2AssetIDs(orderFilled.Side, orderFilled.TokenId)
	if err != nil {
		return nil, err
	}

	return &sharedModels.PolymarketOrderFilled{
		OrderHash:         common.BytesToHash(orderFilled.OrderHash[:]),
		Maker:             orderFilled.Maker,
		Taker:             orderFilled.Taker,
		MakerAssetID:      makerAssetID,
		TakerAssetID:      takerAssetID,
		MakerAmountFilled: orderFilled.MakerAmountFilled,
		TakerAmountFilled: orderFilled.TakerAmountFilled,
		Fee:               orderFilled.Fee,
	}, nil
}

func polymarketV2AssetIDs(side uint8, tokenID *big.Int) (*big.Int, *big.Int, error) {
	if tokenID == nil {
		return nil, nil, fmt.Errorf("v2 order filled token id is nil")
	}

	switch side {
	case 0:
		return big.NewInt(0), new(big.Int).Set(tokenID), nil
	case 1:
		return new(big.Int).Set(tokenID), big.NewInt(0), nil
	default:
		return nil, nil, fmt.Errorf("unexpected v2 order side: %d", side)
	}
}

func (p *CTFParser) ParseTransferSingleLog(log *types.Log) (*sharedModels.PolymarketTransfer, error) {
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	transfer, err := conditionalTokensContract.ParseTransferSingle(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing transfer log: %w", err)
	}

	return &sharedModels.PolymarketTransfer{
		From:    transfer.From,
		To:      transfer.To,
		Amount:  transfer.Value,
		TokenID: transfer.Id,
	}, nil
}

func (p *CTFParser) ParseTransferBatchLog(log *types.Log) ([]*sharedModels.PolymarketTransfer, error) {
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	transferBatch, err := conditionalTokensContract.ParseTransferBatch(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing transfer batch log: %w", err)
	}

	transfers := make([]*sharedModels.PolymarketTransfer, len(transferBatch.Ids))
	for i, id := range transferBatch.Ids {
		transfers[i] = &sharedModels.PolymarketTransfer{
			From:    transferBatch.From,
			To:      transferBatch.To,
			Amount:  transferBatch.Values[i],
			TokenID: id,
		}
	}

	return transfers, nil
}

func (p *CTFParser) IsFPMMLog(log *types.Log) bool {
	if len(log.Topics) == 0 {
		return false
	}
	topic0 := log.Topics[0]
	return topic0 == models.PolymarketFPMMFundingAddedEventSelectorHash ||
		topic0 == models.PolymarketFPMMFundingRemovedEventSelectorHash ||
		topic0 == models.PolymarketFPMMBuyEventSelectorHash ||
		topic0 == models.PolymarketFPMMSellEventSelectorHash
}

func (p *CTFParser) questionIDFromAncillaryData(ancillaryData []byte) common.Hash {
	return crypto.Keccak256Hash(ancillaryData)
}

func (p *CTFParser) getTokenIdPartitionCached(conditionId common.Hash, partition *big.Int, parentCollectionId common.Hash, collateralToken common.Address) *big.Int {
	cacheKey := fmt.Sprintf("%s:%s:%s:%s", conditionId.Hex(), partition.String(), parentCollectionId.Hex(), collateralToken.Hex())

	if tokenId, ok := p.tokenIdCache.Get(cacheKey); ok {
		return tokenId
	}

	tokenId := getTokenIdPartition(conditionId, partition, parentCollectionId, collateralToken)
	if tokenId != nil {
		p.tokenIdCache.Add(cacheKey, tokenId)
	}

	return tokenId
}

func (p *CTFParser) ParseFeeRefundedLog(log *types.Log) (*sharedModels.PolymarketFeeRefunded, error) {
	if len(log.Topics) == 0 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	switch log.Topics[0] {
	case models.PolymarketFeeRefundedEventSelectorHash:
		feeRefunded, err := polymarketFeeModuleContract.ParseFeeRefunded(*log)
		if err != nil {
			return nil, fmt.Errorf("error parsing fee refunded log: %w", err)
		}

		return &sharedModels.PolymarketFeeRefunded{
			To:      feeRefunded.To,
			TokenID: feeRefunded.Id,
			Refund:  feeRefunded.Refund,
		}, nil
	case models.PolymarketNegRiskFeeRefundedEventSelectorHash:
		feeRefunded, err := polymarketNegRiskFeeModuleContract.ParseFeeRefunded(*log)
		if err != nil {
			return nil, fmt.Errorf("error parsing fee refunded log: %w", err)
		}

		return &sharedModels.PolymarketFeeRefunded{
			To:      feeRefunded.To,
			TokenID: feeRefunded.Id,
			Refund:  feeRefunded.Amount,
		}, nil
	}

	return nil, fmt.Errorf("unexpected fee refunded log: %s", log.Topics[0].Hex())
}

func (p *CTFParser) ParseNegRiskPositionConvertedLogRaw(log *types.Log) (*sharedModels.PolymarketNegRiskPositionsConvertedRaw, error) {
	if len(log.Topics) == 0 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	parsedLog, err := polymarketNegRiskAdapterContract.ParsePositionsConverted(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing positions converted log: %w", err)
	}

	return &sharedModels.PolymarketNegRiskPositionsConvertedRaw{
		LogAddress:  log.Address,
		Stakeholder: parsedLog.Stakeholder,
		MarketId:    parsedLog.MarketId,
		IndexSet:    parsedLog.IndexSet,
		Amount:      parsedLog.Amount,
	}, nil
}

var (
	partitionYes   = big.NewInt(1) // index set bit 0
	partitionNo    = big.NewInt(2) // index set bit 1
	feeDenominator = big.NewInt(10000)
)

func (p *CTFParser) ParseNegRiskPositionConverted(
	raw *sharedModels.PolymarketNegRiskPositionsConvertedRaw,
	marketData MarketData,
) *sharedModels.PolymarketNegRiskPositionConverted {
	collateralToken := models.PolymarketNegRiskWrappedCollateralAddress
	marketId := common.BytesToHash(raw.MarketId[:])
	parentCollectionID := common.Hash{}
	amount := raw.Amount
	indexSet := raw.IndexSet
	oracle := raw.LogAddress
	questionCount := marketData.QuestionCount
	feeBips := big.NewInt(int64(marketData.FeeBips))

	noTokens := make([]sharedModels.PolymarketToken, 0)
	yesTokens := make([]sharedModels.PolymarketToken, 0)

	for i := uint8(0); i < questionCount; i++ {
		questionId := getQuestionId(marketId, i)
		conditionId := getConditionId(oracle, questionId)

		if indexSet.Bit(int(i)) == 1 {
			tokenId := getTokenIdPartition(conditionId, partitionNo, parentCollectionID, collateralToken)
			noTokens = append(noTokens, sharedModels.PolymarketToken{
				TokenID:            tokenId,
				ConditionID:        conditionId.Hex(),
				CollateralToken:    collateralToken.Hex(),
				ParentCollectionID: parentCollectionID.Hex(),
				Partition:          new(big.Int).Set(partitionNo),
			})
		} else {
			tokenId := getTokenIdPartition(conditionId, partitionYes, parentCollectionID, collateralToken)
			yesTokens = append(yesTokens, sharedModels.PolymarketToken{
				TokenID:            tokenId,
				ConditionID:        conditionId.Hex(),
				CollateralToken:    collateralToken.Hex(),
				ParentCollectionID: parentCollectionID.Hex(),
				Partition:          new(big.Int).Set(partitionYes),
			})
		}
	}

	feeAmount := new(big.Int).Div(new(big.Int).Mul(amount, feeBips), feeDenominator)
	amountOut := new(big.Int).Sub(amount, feeAmount)

	collateralOut := big.NewInt(0)
	if len(noTokens) > 1 {
		collateralOut = new(big.Int).Mul(big.NewInt(int64(len(noTokens)-1)), amountOut)
	}

	return &sharedModels.PolymarketNegRiskPositionConverted{
		NoTokensGiven:     noTokens,
		YesTokensReceived: yesTokens,
		AmountIn:          new(big.Int).Set(amount),
		AmountOut:         amountOut,
		CollateralOut:     collateralOut,
	}
}

func (p *CTFParser) GetMarketData(marketId common.Hash, adapterAddress common.Address) (MarketData, error) {
	if marketData, ok := p.marketDataCache.Get(marketId); ok {
		return marketData, nil
	}

	ethClient := p.evmClient.GetEthClient()
	contract, err := polymarketNegRiskAdapter.NewPolymarketNegRiskAdapter(adapterAddress, ethClient)
	if err != nil {
		return MarketData{}, fmt.Errorf("error creating neg risk adapter contract: %w", err)
	}

	raw, err := contract.GetMarketData(nil, marketId)
	if err != nil {
		return MarketData{}, fmt.Errorf("error calling getMarketData: %w", err)
	}

	md := ParseMarketData(common.BytesToHash(raw[:]))
	p.marketDataCache.Add(marketId, md)
	return md, nil
}

type MarketData struct {
	QuestionCount uint8
	Determined    bool
	Result        uint8
	FeeBips       uint16
	Oracle        common.Address
}

func ParseMarketData(data common.Hash) MarketData {
	return MarketData{
		QuestionCount: data[0],
		Determined:    data[1] != 0,
		Result:        data[2],
		FeeBips:       binary.BigEndian.Uint16(data[3:5]),
		Oracle:        common.BytesToAddress(data[12:32]),
	}
}

func (p *CTFParser) ParseNegRiskPositionSplitLog(log *types.Log) (*sharedModels.PolymarketPositionSplitMerge, error) {
	if len(log.Topics) > 0 && log.Topics[0] != models.PolymarketNegRiskPositionSplitEventSelectorHash {
		return nil, nil
	}
	if len(log.Topics) != 3 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	positionSplit, err := polymarketNegRiskAdapterContract.ParsePositionSplit(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing position split log: %w", err)
	}

	tokenIDs := make([]*big.Int, 0, 2)
	tokenIDs = append(tokenIDs, getTokenIdPartition(common.BytesToHash(positionSplit.ConditionId[:]), partitionNo, common.Hash{}, models.PolymarketNegRiskWrappedCollateralAddress))
	tokenIDs = append(tokenIDs, getTokenIdPartition(common.BytesToHash(positionSplit.ConditionId[:]), partitionYes, common.Hash{}, models.PolymarketNegRiskWrappedCollateralAddress))

	return &sharedModels.PolymarketPositionSplitMerge{
		Stakeholder:        positionSplit.Stakeholder,
		ConditionID:        common.BytesToHash(positionSplit.ConditionId[:]),
		TokenIDs:           tokenIDs,
		CollateralToken:    models.PolymarketNegRiskWrappedCollateralAddress,
		ParentCollectionId: common.Hash{},
		Partition:          []*big.Int{partitionNo, partitionYes},
		Amount:             positionSplit.Amount,
		IsSplit:            true,
	}, nil
}

func (p *CTFParser) ParseNegRiskPositionMergeLog(log *types.Log) (*sharedModels.PolymarketPositionSplitMerge, error) {
	if len(log.Topics) > 0 && log.Topics[0] != models.PolymarketNegRiskPositionMergeEventSelectorHash {
		return nil, nil
	}
	if len(log.Topics) != 3 {
		return nil, fmt.Errorf("unexpected topics count: %d", len(log.Topics))
	}

	positionMerge, err := polymarketNegRiskAdapterContract.ParsePositionsMerge(*log)
	if err != nil {
		return nil, fmt.Errorf("error parsing position merge log: %w", err)
	}

	tokenIDs := make([]*big.Int, 0, 2)
	tokenIDs = append(tokenIDs, getTokenIdPartition(common.BytesToHash(positionMerge.ConditionId[:]), partitionNo, common.Hash{}, models.PolymarketNegRiskWrappedCollateralAddress))
	tokenIDs = append(tokenIDs, getTokenIdPartition(common.BytesToHash(positionMerge.ConditionId[:]), partitionYes, common.Hash{}, models.PolymarketNegRiskWrappedCollateralAddress))

	return &sharedModels.PolymarketPositionSplitMerge{
		Stakeholder:        positionMerge.Stakeholder,
		ConditionID:        common.BytesToHash(positionMerge.ConditionId[:]),
		TokenIDs:           tokenIDs,
		CollateralToken:    models.PolymarketNegRiskWrappedCollateralAddress,
		ParentCollectionId: common.Hash{},
		Partition:          []*big.Int{partitionNo, partitionYes},
		Amount:             positionMerge.Amount,
		IsSplit:            false,
	}, nil
}
