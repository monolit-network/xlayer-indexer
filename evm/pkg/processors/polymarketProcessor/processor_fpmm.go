package polymarketProcessor

import (
	"math/big"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	sharedModels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

func (s *Processor) processFPMMLogs(logs []*types.Log) error {
	if len(logs) == 0 {
		return nil
	}

	fpmmBuys := make(map[uint]*sharedModels.PolymarketFPMMBuy, 0)
	fpmmSells := make(map[uint]*sharedModels.PolymarketFPMMSell, 0)
	fpmmAddrs := make([]common.Address, 0)
	hasRemoveOrAddFundingLog := false

	for _, log := range logs {
		if len(log.Topics) == 0 {
			continue
		}
		switch log.Topics[0] {
		case models.PolymarketFPMMFundingAddedEventSelectorHash:
			hasRemoveOrAddFundingLog = true
			fpmmAddrs = append(fpmmAddrs, log.Address)
		case models.PolymarketFPMMFundingRemovedEventSelectorHash:
			hasRemoveOrAddFundingLog = true
			fpmmAddrs = append(fpmmAddrs, log.Address)
		case models.PolymarketFPMMBuyEventSelectorHash:
			fpmmBuy, err := s.polymarketParser.ParseFPMMBuyLog(log)
			if err != nil {
				return err
			}
			fpmmBuys[uint(log.Index)] = fpmmBuy
		case models.PolymarketFPMMSellEventSelectorHash:
			fpmmSell, err := s.polymarketParser.ParseFPMMSellLog(log)
			if err != nil {
				return err
			}
			fpmmSells[uint(log.Index)] = fpmmSell
		default:
			continue
		}
	}

	if hasRemoveOrAddFundingLog {
		return nil
	}

	if len(fpmmBuys)+len(fpmmSells) == 0 {
		return nil
	}

	if len(fpmmBuys)+len(fpmmSells) > 1 {
		s.logger.Warn("multiple fpmm buys or sells found in one tx", zap.Int("fpmm_buys", len(fpmmBuys)), zap.Int("fpmm_sells", len(fpmmSells)), zap.String("tx_hash", logs[0].TxHash.Hex()))
		return nil
	}

	baseEvent := basicEventFromLog(logs[0])
	baseEvent.Source = models.PolymarketOrderEventNewSourceOrderFilled

	if len(fpmmBuys) > 0 {
		var buy *sharedModels.PolymarketFPMMBuy
		var logIndex uint
		for logIndex, buy = range fpmmBuys {
			break
		}

		transfers := s.findCTTransfers(logs, buy.FPMMAddress, buy.Buyer)
		if len(transfers) != 1 {
			s.logger.Warn("unexpected number of transfers found for buy", zap.Int("transfers", len(transfers)), zap.String("tx_hash", logs[0].TxHash.Hex()))
			return nil
		}
		transfer := transfers[0]

		event := baseEvent.Copy(0)
		event.LogIndex = uint32(logIndex)
		event.UserAddress = buy.Buyer
		event.SourceAddress = &buy.FPMMAddress
		event.TokenID = transfer.TokenID
		event.TokenAmountDiff = transfer.Amount
		event.UsdcAmountDiff = new(big.Int).Neg(buy.InvestmentAmount)
		if err := s.completeEvent(event); err != nil {
			return err
		}
		s.eventsCh <- event
		return nil
	}

	if len(fpmmSells) > 0 {
		var sell *sharedModels.PolymarketFPMMSell
		var logIndex uint
		for logIndex, sell = range fpmmSells {
			break
		}

		transfers := s.findCTTransfers(logs, sell.Seller, sell.FPMMAddress)
		if len(transfers) != 1 {
			s.logger.Warn("unexpected number of transfers found for buy", zap.Int("transfers", len(transfers)), zap.String("tx_hash", logs[0].TxHash.Hex()))
			return nil
		}
		transfer := transfers[0]

		event := baseEvent.Copy(0)
		event.LogIndex = uint32(logIndex)
		event.UserAddress = sell.Seller
		event.SourceAddress = &sell.FPMMAddress
		event.TokenID = transfer.TokenID
		event.TokenAmountDiff = new(big.Int).Neg(transfer.Amount)
		event.UsdcAmountDiff = sell.ReturnAmount
		if err := s.completeEvent(event); err != nil {
			return err
		}
		s.eventsCh <- event
	}

	return nil
}

func (s *Processor) findCTTransfers(logs []*types.Log, fromAddress common.Address, toAddress common.Address) []*sharedModels.PolymarketTransfer {
	transfers := make([]*sharedModels.PolymarketTransfer, 0)
	for _, log := range logs {
		if len(log.Topics) != 4 {
			continue
		}
		if log.Address != models.PolymarketConditionalTokensAddress {
			continue
		}
		switch log.Topics[0] {
		case models.PolymarketTransferSingleEventSelectorHash:
			transfer, err := s.polymarketParser.ParseTransferSingleLog(log)
			if err != nil {
				continue
			}
			if transfer.From != fromAddress || transfer.To != toAddress {
				continue
			}
			transfers = append(transfers, transfer)
		case models.PolymarketTransferBatchEventSelectorHash:
			transfer, err := s.polymarketParser.ParseTransferBatchLog(log)
			if err != nil {
				continue
			}
			for _, transfer := range transfer {
				if transfer.From != fromAddress || transfer.To != toAddress {
					continue
				}
				transfers = append(transfers, transfer)
			}
		}
	}
	return transfers
}
