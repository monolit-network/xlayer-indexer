package polymarketProcessor

import (
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers/polymarket"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

const (
	retryableCallMaxRetries = 5
	retryableCallBaseDelay  = 200 * time.Millisecond
	retryableCallMaxDelay   = 5 * time.Second
)

func (c *Processor) retryGetMarketData(marketID common.Hash, adapterAddress common.Address) (polymarket.MarketData, error) {
	return retryGetMarketDataCall(c.logger, func() (polymarket.MarketData, error) {
		return c.polymarketParser.GetMarketData(marketID, adapterAddress)
	})
}

func retryGetMarketDataCall(logger *zap.Logger, fn func() (polymarket.MarketData, error)) (polymarket.MarketData, error) {
	delay := retryableCallBaseDelay
	for attempt := 0; ; attempt++ {
		value, err := fn()
		if err == nil {
			return value, nil
		}

		if !evmclient.IsRetryableError(err) || attempt >= retryableCallMaxRetries {
			return polymarket.MarketData{}, err
		}

		logger.Warn(
			"retryable error in polymarket processor",
			zap.String("op", "getMarketData"),
			zap.Int("attempt", attempt+1),
			zap.Duration("backoff", delay),
			zap.Error(err),
		)

		time.Sleep(delay)
		if delay < retryableCallMaxDelay {
			delay *= 2
			if delay > retryableCallMaxDelay {
				delay = retryableCallMaxDelay
			}
		}
	}
}
