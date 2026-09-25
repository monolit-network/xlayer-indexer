package polymarketProcessor

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers/polymarket"
	"github.com/monolit-network/xlayer-indexer/shared/db"
	sharedmodels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/monolit-network/xlayer-indexer/util/notifier"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

func (c *Processor) processOrderFilledLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	orderFilled, err := c.polymarketParser.ParseOrderFilledLog(log)
	if err != nil {
		return err
	}

	if orderFilled == nil {
		c.logger.Warn("order filled log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	zero := big.NewInt(0)
	makerAssetIsZero := orderFilled.MakerAssetID.Cmp(zero) == 0
	takerAssetIsZero := orderFilled.TakerAssetID.Cmp(zero) == 0

	if makerAssetIsZero && takerAssetIsZero {
		c.logger.Warn("both maker and taker asset IDs are zero", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	if !makerAssetIsZero && !takerAssetIsZero {
		c.logger.Warn("both maker and taker asset IDs are not zero", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	event, err := c.processOrderFilledLogOneNotZero(log, orderFilled, allLogs, metadata)
	if err != nil {
		return err
	}

	if err := c.completeEvent(event); err != nil {
		return fmt.Errorf("error completing event: %w", err)
	}

	// c.eventsCh <- event

	return nil
}

func (c *Processor) processOrderFilledLogOneNotZero(log *types.Log, orderFilled *sharedmodels.PolymarketOrderFilled, allLogs []*types.Log, metadata map[string]any) (*models.PolymarketOrderEventNew, error) {
	event := basicEventFromLog(log)
	event.Source = models.PolymarketOrderEventNewSourceOrderFilled
	event.SourceAddress = &orderFilled.Taker
	event.UserAddress = orderFilled.Maker

	assetID := orderFilled.MakerAssetID
	assetAmount := new(big.Int).Neg(orderFilled.MakerAmountFilled)
	usdcAmount := orderFilled.TakerAmountFilled
	if assetID.Cmp(big.NewInt(0)) == 0 {
		assetID = orderFilled.TakerAssetID
		assetAmount, usdcAmount = usdcAmount, assetAmount
	}

	event.TokenID = assetID
	event.TokenAmountDiff = assetAmount
	event.UsdcAmountDiff = usdcAmount
	event.Fee = orderFilled.Fee

	taker := c.findOrderTaker(allLogs, metadata)
	event.IsTaker = taker != (common.Address{}) && event.UserAddress == taker

	if _, ok := metadata["current_tx_orders_filled"]; !ok {
		metadata["current_tx_orders_filled"] = make([]*models.PolymarketOrderEventNew, 0, 1)
	}
	metadata["current_tx_orders_filled"] = append(metadata["current_tx_orders_filled"].([]*models.PolymarketOrderEventNew), &event)

	return &event, nil
}

func (c *Processor) postProcessOrderFilledLog(orderEvents []*models.PolymarketOrderEventNew, refunds []*sharedmodels.PolymarketFeeRefunded) error {
	orderEventsMap := make(map[common.Address]map[string][]*models.PolymarketOrderEventNew) // user address -> token id -> order event
	for _, orderEvent := range orderEvents {
		if _, ok := orderEventsMap[orderEvent.UserAddress]; !ok {
			orderEventsMap[orderEvent.UserAddress] = make(map[string][]*models.PolymarketOrderEventNew)
		}
		orderEventsMap[orderEvent.UserAddress][orderEvent.TokenID.String()] = append(orderEventsMap[orderEvent.UserAddress][orderEvent.TokenID.String()], orderEvent)
	}

	refundsMap := make(map[common.Address]map[string]*big.Int) // user address -> token id -> refund
	for _, refund := range refunds {
		if _, ok := refundsMap[refund.To]; !ok {
			refundsMap[refund.To] = make(map[string]*big.Int)
		}

		if _, ok := refundsMap[refund.To][refund.TokenID.String()]; !ok {
			refundsMap[refund.To][refund.TokenID.String()] = big.NewInt(0)
		}

		if refund.Refund != nil {
			refundsMap[refund.To][refund.TokenID.String()] = new(big.Int).Add(refundsMap[refund.To][refund.TokenID.String()], refund.Refund)
		}
	}

	for userAddress, tokensToOrderEvents := range orderEventsMap {
		for tokenID, orderEvents := range tokensToOrderEvents {
			refundAmountToken := refundsMap[userAddress][tokenID]
			refundAmountUsdc := refundsMap[userAddress]["0"]

			slices.SortFunc(orderEvents, func(a, b *models.PolymarketOrderEventNew) int {
				diff := int(int32(a.LogIndex) - int32(b.LogIndex))
				if diff == 0 {
					return int(int32(a.SubIndex) - int32(b.SubIndex))
				}
				return diff
			})

			resEvent1 := c.constructFinalOrderEvent(orderEvents, refundAmountToken, true)
			resEvent2 := c.constructFinalOrderEvent(orderEvents, refundAmountUsdc, false)

			if resEvent1 != nil {
				if err := c.completeEvent(resEvent1); err != nil {
					return fmt.Errorf("error completing event: %w", err)
				}
				c.eventsCh <- resEvent1
			}
			if resEvent2 != nil {
				if err := c.completeEvent(resEvent2); err != nil {
					return fmt.Errorf("error completing event: %w", err)
				}
				c.eventsCh <- resEvent2
			}
		}
	}

	return nil
}

func (c *Processor) constructFinalOrderEvent(orderEvents []*models.PolymarketOrderEventNew, refundAmount *big.Int, tokenDiffPositive bool) *models.PolymarketOrderEventNew {
	var resEvent *models.PolymarketOrderEventNew
	for _, orderEvent := range orderEvents {
		isTokenDiffPositive := orderEvent.TokenAmountDiff.Cmp(big.NewInt(0)) > 0
		if isTokenDiffPositive != tokenDiffPositive {
			continue
		}

		if resEvent == nil {
			resEvent = orderEvent.Copy(0)
			if resEvent.Fee == nil {
				resEvent.Fee = big.NewInt(0)
			}
			continue
		}

		resEvent.TokenAmountDiff = new(big.Int).Add(resEvent.TokenAmountDiff, orderEvent.TokenAmountDiff)
		resEvent.UsdcAmountDiff = new(big.Int).Add(resEvent.UsdcAmountDiff, orderEvent.UsdcAmountDiff)
		resEvent.Fee = new(big.Int).Add(resEvent.Fee, orderEvent.Fee)

		if resEvent.SourceAddress == nil || orderEvent.SourceAddress == nil || *resEvent.SourceAddress != *orderEvent.SourceAddress {
			resEvent.SourceAddress = nil
		}
	}

	if resEvent == nil {
		return nil
	}

	if refundAmount != nil {
		resEvent.Fee = resEvent.Fee.Sub(resEvent.Fee, refundAmount)
	}

	return resEvent
}

func (c *Processor) findOrderTaker(allLogs []*types.Log, metadata map[string]any) common.Address {
	if taker, ok := metadata["current_tx_order_taker"].(common.Address); ok {
		return taker
	}

	for _, log := range allLogs {
		if len(log.Topics) != 3 || log.Topics[0] != models.PolymarketOrdersMatchedEventV2SelectorHash {
			continue
		}

		taker := common.BytesToAddress(log.Topics[2].Bytes())
		metadata["current_tx_order_taker"] = taker
		return taker
	}

	users := make(map[common.Address]struct{})
	for _, log := range allLogs {
		if len(log.Topics) == 4 && isPolymarketOrderFilledTopic(log.Topics[0]) {
			users[common.BytesToAddress(log.Topics[2].Bytes())] = struct{}{}
		}
	}

	taker := common.Address{}
	for i := len(allLogs) - 1; i >= 0; i-- {
		log := allLogs[i]
		if len(log.Topics) != 4 || !isPolymarketOrderFilledTopic(log.Topics[0]) {
			continue
		}
		if _, ok := users[common.BytesToAddress(log.Topics[3].Bytes())]; !ok {
			taker = common.BytesToAddress(log.Topics[2].Bytes())
			break
		}
	}

	metadata["current_tx_order_taker"] = taker

	return taker
}

func isPolymarketOrderFilledTopic(topic common.Hash) bool {
	return topic == models.PolymarketOrderFilledEventSelectorHash ||
		topic == models.PolymarketOrderFilledEventV2SelectorHash
}

func (c *Processor) processPositionSplitLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	positionSplit, err := c.polymarketParser.ParsePositionSplitLog(log)
	if err != nil {
		return err
	}

	if positionSplit == nil {
		c.logger.Warn("position split log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	if !positionSplit.IsSplit {
		c.logger.Warn("position split is merge", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	if positionSplit.Amount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	tokens := make([]sharedmodels.PolymarketToken, len(positionSplit.TokenIDs))
	for i, tokenID := range positionSplit.TokenIDs {
		tokens[i] = sharedmodels.PolymarketToken{
			TokenID:            tokenID,
			ConditionID:        positionSplit.ConditionID.Hex(),
			CollateralToken:    positionSplit.CollateralToken.Hex(),
			ParentCollectionID: positionSplit.ParentCollectionId.Hex(),
			Partition:          positionSplit.Partition[i],
			CreatedAt:          time.Unix(int64(log.BlockTimestamp), 0),
		}
	}

	blockTime := time.Unix(int64(log.BlockTimestamp), 0)
	if err := c.insertTokenMapping(positionSplit.ConditionID, tokens, blockTime); err != nil {
		c.logger.Warn("failed to insert token mapping (tokens may already exist)", zap.Error(err))
	}

	userAddress, shouldProcess, err := c.resolveCollateralAdapterStakeholder(
		positionSplit.Stakeholder,
		log,
		allLogs,
		metadata,
		collateralAdapterExpectedTransfers(positionSplit.TokenIDs, positionSplit.Amount),
	)
	if err != nil {
		return err
	}
	if !shouldProcess {
		return nil
	}

	baseEvent := basicEventFromLog(log)
	baseEvent.Source = models.PolymarketOrderEventNewSourcePositionSplit
	baseEvent.SourceAddress = nil
	baseEvent.UserAddress = userAddress

	isStakeholderPolymarketContract := c.isPolymarketContractAddress(userAddress, metadata)

	usdcAmountDiff := new(big.Int).Div(positionSplit.Amount, big.NewInt(int64(len(positionSplit.TokenIDs))))
	usdcAmountDiff.Neg(usdcAmountDiff)

	for i, tokenID := range positionSplit.TokenIDs {
		if !isStakeholderPolymarketContract {
			event := baseEvent.Copy(uint32(i))
			event.TokenID = tokenID
			event.TokenAmountDiff = new(big.Int).Set(positionSplit.Amount)
			event.UsdcAmountDiff = usdcAmountDiff
			event.ConditionID = positionSplit.ConditionID
			if err := c.completeEvent(event); err != nil {
				return fmt.Errorf("error completing event: %w", err)
			}
			c.eventsCh <- event
		}
	}

	return nil
}

func (c *Processor) processNegRiskPositionSplitLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	positionSplit, err := c.polymarketParser.ParseNegRiskPositionSplitLog(log)
	if err != nil {
		c.logger.Warn("failed to parse neg risk position split log", zap.Error(err))
		return nil
	}

	if positionSplit == nil {
		c.logger.Warn("position split log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	if !positionSplit.IsSplit {
		c.logger.Warn("position split is merge", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	if positionSplit.Amount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	userAddress, shouldProcess, err := c.resolveCollateralAdapterStakeholder(
		positionSplit.Stakeholder,
		log,
		allLogs,
		metadata,
		collateralAdapterExpectedTransfers(positionSplit.TokenIDs, positionSplit.Amount),
	)
	if err != nil {
		return err
	}
	if !shouldProcess {
		return nil
	}

	baseEvent := basicEventFromLog(log)
	baseEvent.Source = models.PolymarketOrderEventNewSourcePositionSplit
	baseEvent.SourceAddress = nil
	baseEvent.UserAddress = userAddress

	isStakeholderPolymarketContract := c.isPolymarketContractAddress(userAddress, metadata)

	usdcAmountDiff := new(big.Int).Div(positionSplit.Amount, big.NewInt(int64(len(positionSplit.TokenIDs))))
	usdcAmountDiff.Neg(usdcAmountDiff)

	for i, tokenID := range positionSplit.TokenIDs {
		if !isStakeholderPolymarketContract {
			event := baseEvent.Copy(uint32(i))
			event.TokenID = tokenID
			event.TokenAmountDiff = new(big.Int).Set(positionSplit.Amount)
			event.UsdcAmountDiff = usdcAmountDiff
			event.ConditionID = positionSplit.ConditionID
			if err := c.completeEvent(event); err != nil {
				return fmt.Errorf("error completing event: %w", err)
			}
			c.eventsCh <- event
		}
	}

	return nil
}

func (c *Processor) processPositionMergeLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	positionMerge, err := c.polymarketParser.ParsePositionMergeLog(log)
	if err != nil {
		return err
	}

	if positionMerge == nil {
		c.logger.Warn("position merge log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	if positionMerge.IsSplit {
		c.logger.Warn("position merge is split", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	if positionMerge.Amount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	userAddress, shouldProcess, err := c.resolveCollateralAdapterStakeholder(
		positionMerge.Stakeholder,
		log,
		allLogs,
		metadata,
		collateralAdapterExpectedTransfers(positionMerge.TokenIDs, positionMerge.Amount),
	)
	if err != nil {
		return err
	}
	if !shouldProcess {
		return nil
	}

	baseEvent := basicEventFromLog(log)
	baseEvent.Source = models.PolymarketOrderEventNewSourcePositionMerge
	baseEvent.SourceAddress = nil
	baseEvent.UserAddress = userAddress

	if c.isPolymarketContractAddress(userAddress, metadata) {
		return nil
	}

	usdcAmountDiff := new(big.Int).Div(positionMerge.Amount, big.NewInt(int64(len(positionMerge.TokenIDs))))

	for i, tokenID := range positionMerge.TokenIDs {
		event := baseEvent.Copy(uint32(i))
		event.TokenID = tokenID
		event.TokenAmountDiff = new(big.Int).Neg(positionMerge.Amount)
		event.UsdcAmountDiff = usdcAmountDiff
		event.ConditionID = positionMerge.ConditionID
		if err := c.completeEvent(event); err != nil {
			return fmt.Errorf("error completing event: %w", err)
		}
		c.eventsCh <- event
	}

	return nil
}

func (c *Processor) processNegRiskPositionMergeLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	positionMerge, err := c.polymarketParser.ParseNegRiskPositionMergeLog(log)
	if err != nil {
		c.logger.Warn("failed to parse neg risk position merge log", zap.Error(err))
		return nil
	}

	if positionMerge == nil {
		c.logger.Warn("position merge log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	if positionMerge.IsSplit {
		c.logger.Warn("position merge is split", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	if positionMerge.Amount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	userAddress, shouldProcess, err := c.resolveCollateralAdapterStakeholder(
		positionMerge.Stakeholder,
		log,
		allLogs,
		metadata,
		collateralAdapterExpectedTransfers(positionMerge.TokenIDs, positionMerge.Amount),
	)
	if err != nil {
		return err
	}
	if !shouldProcess {
		return nil
	}

	baseEvent := basicEventFromLog(log)
	baseEvent.Source = models.PolymarketOrderEventNewSourcePositionMerge
	baseEvent.SourceAddress = nil
	baseEvent.UserAddress = userAddress

	if c.isPolymarketContractAddress(userAddress, metadata) {
		return nil
	}

	usdcAmountDiff := new(big.Int).Div(positionMerge.Amount, big.NewInt(int64(len(positionMerge.TokenIDs))))

	for i, tokenID := range positionMerge.TokenIDs {
		event := baseEvent.Copy(uint32(i))
		event.TokenID = tokenID
		event.TokenAmountDiff = new(big.Int).Neg(positionMerge.Amount)
		event.UsdcAmountDiff = usdcAmountDiff
		event.ConditionID = positionMerge.ConditionID
		if err := c.completeEvent(event); err != nil {
			return fmt.Errorf("error completing event: %w", err)
		}
		c.eventsCh <- event
	}

	return nil
}

func (c *Processor) processPositionRedeemLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	payoutRedemption, err := c.polymarketParser.ParsePayoutRedemptionLog(log)
	if err != nil {
		return err
	}

	if payoutRedemption == nil {
		c.logger.Warn("payout redemption log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	if payoutRedemption.Redeemer == models.PolymarketNegRiskAdapterAddress {
		return nil
	}

	amounts, err := c.findPositionRedeemBurnAmounts(log, allLogs, payoutRedemption)
	if err != nil {
		return fmt.Errorf("failed to find position redeem burn amounts: %w", err)
	}
	payoutRedemption.Amounts = amounts

	baseEvent := basicEventFromLog(log)
	return c.processPositionRedeem(payoutRedemption, baseEvent, metadata)
}

func (c *Processor) findPositionRedeemBurnAmounts(currentLog *types.Log, allLogs []*types.Log, payoutRedemption *sharedmodels.PolymarketPayoutRedemption) ([]*big.Int, error) {
	amounts := make([]*big.Int, len(payoutRedemption.TokenIDs))
	tokenIndexes := make(map[string]int, len(payoutRedemption.TokenIDs))
	for i, tokenID := range payoutRedemption.TokenIDs {
		amounts[i] = big.NewInt(0)
		if tokenID != nil {
			tokenIndexes[tokenID.String()] = i
		}
	}

	currentLogIdx := -1
	for i, txLog := range allLogs {
		if txLog == currentLog {
			currentLogIdx = i
			break
		}
	}
	if currentLogIdx == -1 {
		return amounts, nil
	}

	for i := currentLogIdx - 1; i >= 0; i-- {
		txLog := allLogs[i]
		if len(txLog.Topics) == 0 {
			continue
		}
		if txLog.Topics[0] == models.PolymarketPayoutRedemptionEventSelectorHash {
			break
		}
		if txLog.Address != models.PolymarketConditionalTokensAddress {
			continue
		}

		var transfers []*sharedmodels.PolymarketTransfer
		switch txLog.Topics[0] {
		case models.PolymarketTransferSingleEventSelectorHash:
			transfer, err := c.polymarketParser.ParseTransferSingleLog(txLog)
			if err != nil {
				return nil, err
			}
			transfers = []*sharedmodels.PolymarketTransfer{transfer}
		case models.PolymarketTransferBatchEventSelectorHash:
			parsedTransfers, err := c.polymarketParser.ParseTransferBatchLog(txLog)
			if err != nil {
				return nil, err
			}
			transfers = parsedTransfers
		default:
			continue
		}

		matchedBurn := false
		payoutMint := len(transfers) > 0
		for _, transfer := range transfers {
			if transfer == nil {
				payoutMint = false
				continue
			}
			if transfer.From == models.ZeroAddress && transfer.To == payoutRedemption.Redeemer {
				continue
			}
			payoutMint = false
			if transfer.From == payoutRedemption.Redeemer && transfer.To == models.ZeroAddress {
				idx, ok := tokenIndexes[transfer.TokenID.String()]
				if !ok {
					continue
				}
				amounts[idx].Add(amounts[idx], transfer.Amount)
				matchedBurn = true
			}
		}
		if !matchedBurn && !payoutMint {
			break
		}
	}

	return amounts, nil
}

func (c *Processor) processNegRiskPayoutRedemptionLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	payoutRedemption, err := c.polymarketParser.ParseNegRiskPayoutRedemptionLog(log)
	if err != nil {
		return err
	}

	if payoutRedemption == nil {
		c.logger.Warn("payout redemption log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	baseEvent := basicEventFromLog(log)
	return c.processPositionRedeem(payoutRedemption, baseEvent, metadata)
}

func (c *Processor) processPositionRedeem(payoutRedemption *sharedmodels.PolymarketPayoutRedemption, baseEvent models.PolymarketOrderEventNew, metadata map[string]any) error {
	if payoutRedemption.Amounts == nil {
		return fmt.Errorf("payout redemption amounts are nil")
	}
	if c.isPolymarketContractAddress(payoutRedemption.Redeemer, metadata) {
		return nil
	}

	baseEvent.Source = models.PolymarketOrderEventNewSourcePositionRedeem
	baseEvent.SourceAddress = nil
	baseEvent.UserAddress = payoutRedemption.Redeemer

	totalAmount := big.NewInt(0)
	for _, amount := range payoutRedemption.Amounts {
		totalAmount.Add(totalAmount, amount)
	}

	subIdx := uint32(0)
	for i, tokenID := range payoutRedemption.TokenIDs {
		tokenAmount := payoutRedemption.Amounts[i]
		if tokenAmount.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		payout, ok, err := c.getTokenIDToPayout(tokenID)
		if err != nil {
			return fmt.Errorf("failed to get token id payout: %w", err)
		}
		if !ok {
			return fmt.Errorf("token not resolved or found: %s", tokenID.String())
		}
		event := baseEvent.Copy(subIdx)
		event.TokenID = tokenID
		event.TokenAmountDiff = new(big.Int).Neg(payoutRedemption.Amounts[i])
		event.UsdcAmountDiff = new(big.Int).Mul(tokenAmount, payout.Numerator)
		event.UsdcAmountDiff.Div(event.UsdcAmountDiff, payout.Denominator)
		if err := c.completeEvent(event); err != nil {
			return fmt.Errorf("error completing event: %w", err)
		}
		c.eventsCh <- event
		subIdx++
	}

	return nil
}

type collateralAdapterExpectedTransfer struct {
	TokenID *big.Int
	Amount  *big.Int
}

func (c *Processor) resolveCollateralAdapterStakeholder(
	stakeholder common.Address,
	log *types.Log,
	allLogs []*types.Log,
	metadata map[string]any,
	expected []collateralAdapterExpectedTransfer,
) (common.Address, bool, error) {
	if !isPolymarketCollateralAdapterAddress(stakeholder) {
		return stakeholder, true, nil
	}

	userAddress, ok, err := c.findCollateralAdapterTransferUser(stakeholder, allLogs, metadata, expected)
	if err != nil {
		return common.Address{}, false, err
	}
	if !ok {
		// if c.logger != nil {
		// 	c.logger.Warn("failed to resolve collateral adapter user",
		// 		zap.String("adapter", stakeholder.Hex()),
		// 		zap.String("tx_hash", log.TxHash.Hex()),
		// 		zap.Uint64("log_index", uint64(log.Index)),
		// 	)
		// }
		return common.Address{}, false, nil
	}

	return userAddress, true, nil
}

func isPolymarketCollateralAdapterAddress(address common.Address) bool {
	return address == models.PolymarketCtfCollateralAdapterAddress ||
		address == models.PolymarketNegRiskCtfCollateralAdapterAddress ||
		address == models.PolymarketNegRiskCtfCollateralAdapterAddress2
}

func collateralAdapterExpectedTransfers(tokenIDs []*big.Int, amount *big.Int) []collateralAdapterExpectedTransfer {
	expected := make([]collateralAdapterExpectedTransfer, 0, len(tokenIDs))
	for _, tokenID := range tokenIDs {
		expected = append(expected, collateralAdapterExpectedTransfer{
			TokenID: tokenID,
			Amount:  amount,
		})
	}
	return expected
}

func (c *Processor) findCollateralAdapterTransferUser(
	adapter common.Address,
	allLogs []*types.Log,
	metadata map[string]any,
	expected []collateralAdapterExpectedTransfer,
) (common.Address, bool, error) {
	var best common.Address
	var bestDiff *big.Int

	for _, log := range allLogs {
		if log.Address != models.PolymarketConditionalTokensAddress || len(log.Topics) != 4 {
			continue
		}

		var transfers []*sharedmodels.PolymarketTransfer
		switch log.Topics[0] {
		case models.PolymarketTransferSingleEventSelectorHash:
			transfer, err := c.polymarketParser.ParseTransferSingleLog(log)
			if err != nil {
				return common.Address{}, false, err
			}
			transfers = []*sharedmodels.PolymarketTransfer{transfer}
		case models.PolymarketTransferBatchEventSelectorHash:
			parsedTransfers, err := c.polymarketParser.ParseTransferBatchLog(log)
			if err != nil {
				return common.Address{}, false, err
			}
			transfers = parsedTransfers
		default:
			continue
		}

		for _, transfer := range transfers {
			var user common.Address
			switch {
			case transfer.From == adapter:
				user = transfer.To
			case transfer.To == adapter:
				user = transfer.From
			default:
				continue
			}
			if user == models.ZeroAddress || c.isPolymarketContractAddress(user, metadata) {
				continue
			}

			for _, e := range expected {
				if transfer.TokenID.Cmp(e.TokenID) != 0 {
					continue
				}

				diff := absBigInt(new(big.Int).Sub(transfer.Amount, e.Amount))
				if bestDiff == nil || diff.Cmp(bestDiff) < 0 {
					best = user
					bestDiff = diff
				}
			}
		}
	}

	return best, bestDiff != nil, nil
}

func absBigInt(n *big.Int) *big.Int {
	if n.Sign() < 0 {
		return n.Neg(n)
	}
	return n
}

func (c *Processor) processConditionPreparationLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	conditionPreparation, err := c.polymarketParser.ParseConditionPreparationLog(log)
	if err != nil {
		return err
	}

	if conditionPreparation == nil {
		c.logger.Warn("condition preparation log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	payoutNumerators := make([]*big.Int, conditionPreparation.OutcomesCount)
	for i := range payoutNumerators {
		payoutNumerators[i] = big.NewInt(0)
	}

	questionIDLower := strings.ToLower(conditionPreparation.QuestionID.Hex())

	market := sharedmodels.PolymarketMarketNew{
		ConditionID:         conditionPreparation.ConditionID.Hex(),
		QuestionID:          questionIDLower,
		Oracle:              conditionPreparation.Oracle.Hex(),
		PreparedAt:          time.Unix(int64(log.BlockTimestamp), 0),
		PreparedInBlock:     int(log.BlockNumber),
		PreparedInBlockHash: log.BlockHash,
		PreparedInTxHash:    log.TxHash,
		TotalOutcomes:       int(conditionPreparation.OutcomesCount),
		IsResolved:          false,
		ResolvedAt:          nil,
		PayoutNumerators:    payoutNumerators,
	}

	eventData := sharedmodels.PolymarketMarketEventDataCondPrepared{
		ConditionID:   strings.ToLower(conditionPreparation.ConditionID.Hex()),
		Oracle:        strings.ToLower(conditionPreparation.Oracle.Hex()),
		OutcomesCount: conditionPreparation.OutcomesCount,
	}

	event := sharedmodels.PolymarketMarketEvent{
		QuestionID:     strings.ToLower(conditionPreparation.QuestionID.Hex()),
		EventType:      sharedmodels.PolymarketEventCondPrepared,
		BlockNumber:    log.BlockNumber,
		BlockHash:      strings.ToLower(log.BlockHash.Hex()),
		LogIndex:       uint32(log.Index),
		TxIndex:        uint32(log.TxIndex),
		TxHash:         strings.ToLower(log.TxHash.Hex()),
		BlockTimestamp: time.Unix(int64(log.BlockTimestamp), 0),
		EventData:      &eventData,
	}

	if c.notifier != nil {
		if err := notifier.Publish(c.notifier, context.Background(), string(sharedmodels.NotifierEventCondPrepared), event); err != nil {
			c.notifier.Logger.Error("failed to publish condition prepared event", zap.Error(err))
		}
	}

	if err := c.dbClient.TransactionWithRetry(context.Background(), -1, func(ctx context.Context, tx db.DB) error {
		if err := tx.InsertPolymarketMarketNew(ctx, market); err != nil {
			return fmt.Errorf("failed to insert market: %w", err)
		}
		if err := tx.InsertPolymarketMarketEvent(ctx, event); err != nil {
			c.recordPolymarketMarketEventInsertError()
			return fmt.Errorf("failed to insert market event (%s): %w", event.EventType, err)
		}
		return nil
	}); err != nil {
		return err
	}
	c.recordPolymarketMarketEventInserted()

	// Cache condition <-> question mapping
	conditionIDLower := strings.ToLower(conditionPreparation.ConditionID.Hex())
	questionIDLower = strings.ToLower(conditionPreparation.QuestionID.Hex())
	if err := c.cacheConditionQuestionMapping(conditionIDLower, questionIDLower); err != nil {
		c.logger.Warn("failed to cache condition-question mapping", zap.Error(err))
	}

	// TODO: REMOVE
	if event.QuestionID == "0xeeec61bff924dfefb3cf1555f2c7762d99f85a657a8b6dd9c72c2eb2cd296fb3" {
		c.logger.Info("condition preparation log", zap.Any("event", event))
	}

	return nil
}

func (c *Processor) processConditionResolutionLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	start := time.Now()
	conditionResolution, err := c.polymarketParser.ParseConditionResolutionLog(log)
	if err != nil {
		return err
	}

	c.recordProcessDuration("condition_resolution_parse", time.Since(start))
	start = time.Now()

	if conditionResolution == nil {
		c.logger.Warn("condition resolution log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	blockTime := time.Unix(int64(log.BlockTimestamp), 0)
	conditionID := strings.ToLower(conditionResolution.ConditionID.Hex())
	numerators := conditionResolution.Numerators

	questionID, err := c.getQuestionIDByConditionID(conditionID)
	if err != nil {
		return fmt.Errorf("failed to get question id by condition id: %w", err)
	}
	if questionID == "" {
		return fmt.Errorf("question id not found for condition id: %s", conditionID)
	}

	c.recordProcessDuration("condition_resolution_get_question_id", time.Since(start))
	start = time.Now()

	numeratorsStr := make([]string, len(numerators))
	for i, n := range numerators {
		numeratorsStr[i] = n.String()
	}

	eventData := sharedmodels.PolymarketMarketEventDataCondResolved{
		ConditionID: conditionID,
		Oracle:      strings.ToLower(conditionResolution.Oracle.Hex()),
		Numerators:  numeratorsStr,
	}

	blockHash := strings.ToLower(log.BlockHash.Hex())

	event := sharedmodels.PolymarketMarketEvent{
		QuestionID:     strings.ToLower(questionID),
		EventType:      sharedmodels.PolymarketEventCondResolved,
		BlockNumber:    log.BlockNumber,
		BlockHash:      blockHash,
		LogIndex:       uint32(log.Index),
		TxIndex:        uint32(log.TxIndex),
		TxHash:         strings.ToLower(log.TxHash.Hex()),
		BlockTimestamp: blockTime,
		EventData:      &eventData,
	}

	start = time.Now()

	if c.notifier != nil {
		if err := notifier.Publish(c.notifier, context.Background(), string(sharedmodels.NotifierEventCondResolved), event); err != nil {
			c.notifier.Logger.Error("failed to publish condition resolved event", zap.Error(err))
		}
	}

	var resolvedTokens []sharedmodels.PolymarketToken
	c.flushPgSync()
	if err := c.dbClient.TransactionWithRetry(context.Background(), -1, func(ctx context.Context, tx db.DB) error {
		var err error
		resolvedTokens, err = tx.UpdatePolymarketMarketResolutionFull(ctx, conditionID, blockTime, numerators, int64(log.BlockNumber), blockHash)
		if err != nil {
			return fmt.Errorf("failed to update market resolution: %w", err)
		}
		c.recordProcessDuration("condition_resolution_update_market_resolution", time.Since(start))
		start = time.Now()
		if err := tx.InsertPolymarketMarketEvent(ctx, event); err != nil {
			c.recordPolymarketMarketEventInsertError()
			return fmt.Errorf("failed to insert market event (%s): %w", event.EventType, err)
		}
		c.recordProcessDuration("condition_resolution_insert_market_event", time.Since(start))
		return nil
	}); err != nil {
		// c.logger.Warn("failed to update market resolution", zap.Error(err))
		return err
	}
	c.recordPolymarketMarketEventInserted()
	for _, token := range resolvedTokens {
		if token.TokenID == nil {
			continue
		}
		c.setTokenIDToPayout(token.TokenID, tokenIDToPayout{
			Numerator:   token.Numerator,
			Denominator: token.Denominator,
		})
	}

	c.recordProcessDuration("condition_resolution_write", time.Since(start))

	return nil
}

func (c *Processor) processTransferSingleLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	transfer, err := c.polymarketParser.ParseTransferSingleLog(log)
	if err != nil {
		return err
	}

	if transfer == nil {
		c.logger.Warn("transfer single log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	if c.isPolymarketContractAddress(transfer.From, metadata) || c.isPolymarketContractAddress(transfer.To, metadata) {
		return nil
	}

	if transfer.From == models.ZeroAddress || transfer.To == models.ZeroAddress {
		// ignore mints and burns
		return nil
	}

	if transfer.Amount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	baseEvent := basicEventFromLog(log)
	baseEvent.Source = models.PolymarketOrderEventNewSourcePositionTransfer

	event1 := baseEvent.Copy(0)
	event1.UserAddress = transfer.To
	event1.SourceAddress = &transfer.From
	event1.TokenID = transfer.TokenID
	event1.TokenAmountDiff = transfer.Amount
	event1.UsdcAmountDiff = big.NewInt(0)
	if err := c.completeEvent(event1); err != nil {
		return fmt.Errorf("error completing event: %w", err)
	}
	c.eventsCh <- event1

	event2 := baseEvent.Copy(1)
	event2.UserAddress = transfer.From
	event2.SourceAddress = &transfer.To
	event2.TokenID = transfer.TokenID
	event2.TokenAmountDiff = new(big.Int).Neg(transfer.Amount)
	event2.UsdcAmountDiff = big.NewInt(0)
	if err := c.completeEvent(event2); err != nil {
		return fmt.Errorf("error completing event: %w", err)
	}
	c.eventsCh <- event2

	return nil
}

func (c *Processor) processTransferBatchLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	transfers, err := c.polymarketParser.ParseTransferBatchLog(log)
	if err != nil {
		return err
	}

	if len(transfers) == 0 {
		// c.logger.Warn("transfer batch log is empty", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	if c.isPolymarketContractAddress(transfers[0].From, metadata) || c.isPolymarketContractAddress(transfers[0].To, metadata) {
		return nil
	}

	for _, transfer := range transfers {
		if transfer.From == models.ZeroAddress || transfer.To == models.ZeroAddress {
			// ignore mints and burns
			return nil
		}
	}

	baseEvent := basicEventFromLog(log)
	baseEvent.Source = models.PolymarketOrderEventNewSourcePositionTransferBatch

	for i, transfer := range transfers {
		if transfer.Amount.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		event1 := baseEvent.Copy(uint32(i * 2))
		event1.UserAddress = transfer.To
		event1.SourceAddress = &transfer.From
		event1.TokenID = transfer.TokenID
		event1.TokenAmountDiff = transfer.Amount
		event1.UsdcAmountDiff = big.NewInt(0)
		if err := c.completeEvent(event1); err != nil {
			return fmt.Errorf("error completing event: %w", err)
		}
		c.eventsCh <- event1

		event2 := baseEvent.Copy(uint32(i*2 + 1))
		event2.UserAddress = transfer.From
		event2.SourceAddress = &transfer.To
		event2.TokenID = transfer.TokenID
		event2.TokenAmountDiff = new(big.Int).Neg(transfer.Amount)
		event2.UsdcAmountDiff = big.NewInt(0)
		if err := c.completeEvent(event2); err != nil {
			return fmt.Errorf("error completing event: %w", err)
		}
		c.eventsCh <- event2
	}

	return nil
}

func (c *Processor) processFeeRefundedLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	refund, err := c.polymarketParser.ParseFeeRefundedLog(log)
	if err != nil {
		c.logger.Warn("failed to parse fee refunded log", zap.Error(err))
		return nil
	}

	if refund == nil {
		c.logger.Warn("fee refunded log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	if _, ok := metadata["current_tx_fee_refunds"]; !ok {
		metadata["current_tx_fee_refunds"] = make([]*sharedmodels.PolymarketFeeRefunded, 0, 1)
	}
	metadata["current_tx_fee_refunds"] = append(metadata["current_tx_fee_refunds"].([]*sharedmodels.PolymarketFeeRefunded), refund)
	return nil
}

func (c *Processor) processNegRiskPositionConvertedLog(log *types.Log, allLogs []*types.Log, metadata map[string]any) error {
	convertedRaw, err := c.polymarketParser.ParseNegRiskPositionConvertedLogRaw(log)
	if err != nil {
		c.logger.Warn("failed to parse neg risk position converted log", zap.Error(err))
		return nil
	}

	if convertedRaw == nil {
		c.logger.Warn("neg risk position converted log is nil", zap.String("tx_hash", log.TxHash.Hex()), zap.Uint64("log_index", uint64(log.Index)))
		return nil
	}

	var marketData polymarket.MarketData
	marketData, err = c.retryGetMarketData(convertedRaw.MarketId, convertedRaw.LogAddress)
	if err != nil {
		return err
	}

	converted := c.polymarketParser.ParseNegRiskPositionConverted(convertedRaw, marketData)

	expectedTransfers := make([]collateralAdapterExpectedTransfer, 0, len(converted.NoTokensGiven)+len(converted.YesTokensReceived))
	for _, token := range converted.NoTokensGiven {
		expectedTransfers = append(expectedTransfers, collateralAdapterExpectedTransfer{
			TokenID: token.TokenID,
			Amount:  converted.AmountIn,
		})
	}
	for _, token := range converted.YesTokensReceived {
		expectedTransfers = append(expectedTransfers, collateralAdapterExpectedTransfer{
			TokenID: token.TokenID,
			Amount:  converted.AmountOut,
		})
	}

	userAddress, shouldProcess, err := c.resolveCollateralAdapterStakeholder(convertedRaw.Stakeholder, log, allLogs, metadata, expectedTransfers)
	if err != nil {
		return err
	}
	if !shouldProcess {
		return nil
	}

	basicEvent := basicEventFromLog(log)
	basicEvent.Source = models.PolymarketOrderEventNewSourceNegRiskPositionConverted
	basicEvent.UserAddress = userAddress

	totalEvents := 0
	usdcRecvPerEvent := new(big.Int).Div(converted.CollateralOut, big.NewInt(int64(len(converted.NoTokensGiven))))
	for _, noToken := range converted.NoTokensGiven {
		event := basicEvent.Copy(uint32(totalEvents))
		event.TokenID = noToken.TokenID
		event.TokenAmountDiff = new(big.Int).Neg(converted.AmountIn)
		event.UsdcAmountDiff = usdcRecvPerEvent
		event.ConditionID = common.HexToHash(noToken.ConditionID)
		if err := c.completeEvent(event); err != nil {
			return fmt.Errorf("error completing event: %w", err)
		}
		c.eventsCh <- event
		totalEvents++
	}

	for _, yesToken := range converted.YesTokensReceived {
		event := basicEvent.Copy(uint32(totalEvents))
		event.TokenID = yesToken.TokenID
		event.TokenAmountDiff = converted.AmountOut
		event.UsdcAmountDiff = big.NewInt(0)
		event.ConditionID = common.HexToHash(yesToken.ConditionID)
		if err := c.completeEvent(event); err != nil {
			return fmt.Errorf("error completing event: %w", err)
		}
		c.eventsCh <- event
		totalEvents++
	}

	tokenMapping := buildTokenMapping(converted.NoTokensGiven, converted.YesTokensReceived)
	blockTime := time.Unix(int64(log.BlockTimestamp), 0)
	for conditionID, tokens := range tokenMapping {
		if err := c.insertTokenMapping(common.HexToHash(conditionID), tokens, blockTime); err != nil {
			return fmt.Errorf("error inserting token mapping: %w", err)
		}
	}

	return nil
}

func buildTokenMapping(tokens ...[]sharedmodels.PolymarketToken) map[string][]sharedmodels.PolymarketToken {
	tokenMapping := make(map[string][]sharedmodels.PolymarketToken)
	for _, tokenArr := range tokens {
		for _, token := range tokenArr {
			tokenMapping[token.ConditionID] = append(tokenMapping[token.ConditionID], token)
		}
	}
	return tokenMapping
}
