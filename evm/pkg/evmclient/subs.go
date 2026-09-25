package evmclient

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"

	"go.uber.org/zap"
)

func (c *Client) SubscribeHeads(ctx context.Context) (<-chan *types.Header, error) {
	if c.wsEthClient == nil {
		return nil, fmt.Errorf("ws eth client is not set")
	}

	subFn := func() (ethereum.Subscription, chan *types.Header, error) {
		headChan := make(chan *types.Header, 1024)
		sub, err := c.wsEthClient.SubscribeNewHead(ctx, headChan)
		if err != nil {
			close(headChan)
			return nil, nil, err
		}
		return sub, headChan, nil
	}

	return subscribeWithReconnect(ctx, c.logger, subFn)
}

// subscribeFn should close result channel if error subscribing
func subscribeWithReconnect[T any](ctx context.Context, logger *zap.Logger, subscribeFn func() (ethereum.Subscription, chan T, error)) (<-chan T, error) {
	resChan := make(chan T)

	sub, subChan, err := subscribeFn()
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			sub.Unsubscribe()
			close(resChan)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case err := <-sub.Err():
				logger.Warn("subscription error, reconnecting", zap.Error(err))

				sub.Unsubscribe()

			RcnLoop:
				for {
					select {
					case <-ctx.Done():
						return
					case <-time.After(5 * time.Second):
						logger.Info("attempting to reconnect subscription")
						sub, subChan, err = subscribeFn()
						if err == nil {
							break RcnLoop
						}
						logger.Warn("reconnect failed, retrying", zap.Error(err))
					}
				}

			case data := <-subChan:
				logger.Debug("data received", zap.Any("data", data))
				select {
				case <-ctx.Done():
					return
				case resChan <- data:
				}
			}
		}
	}()

	return resChan, nil
}
