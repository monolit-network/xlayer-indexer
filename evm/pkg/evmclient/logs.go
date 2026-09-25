package evmclient

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

func (c *Client) getLogs(ctx context.Context, addresses []common.Address, topics [][]common.Hash, startBlock *big.Int, endBlock *big.Int) ([]types.Log, error) {
	c.recordRPCCall()
	logs, err := c.ethClient.FilterLogs(ctx, ethereum.FilterQuery{
		Addresses: addresses,
		Topics:    topics,
		FromBlock: startBlock,
		ToBlock:   endBlock,
	})
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (c *Client) GetLogs(ctx context.Context, addresses [][]common.Address, topics [][][]common.Hash, startBlock *big.Int, endBlock *big.Int, infiniteRetries bool) ([]types.Log, error) {
	if len(addresses) != len(topics) {
		return nil, fmt.Errorf("addresses and topics must have the same length")
	}

	allLogs := make([]types.Log, 0)
	i := 0
	for i < len(addresses) {
		logs, err := c.getLogs(ctx, addresses[i], topics[i], startBlock, endBlock)
		if err != nil {
			if isBatchSplitError(err) {
				c.recordBatchSplit()
				batchSize := new(big.Int).Sub(endBlock, startBlock).Int64()
				if batchSize <= 1 {
					return nil, err
				}

				c.logger.Warn("logs response too large, splitting batch", zap.Error(err), zap.Int("batch_size", int(batchSize)), zap.String("chain", string(c.Chain)))
				midBlock := new(big.Int).Add(startBlock, endBlock)
				midBlock.Div(midBlock, big.NewInt(2))
				logs1, err1 := c.GetLogs(ctx, [][]common.Address{addresses[i]}, [][][]common.Hash{topics[i]}, startBlock, midBlock, infiniteRetries)
				if err1 != nil {
					return nil, err1
				}
				logs2, err2 := c.GetLogs(ctx, [][]common.Address{addresses[i]}, [][][]common.Hash{topics[i]}, new(big.Int).Add(midBlock, big.NewInt(1)), endBlock, infiniteRetries)
				if err2 != nil {
					return nil, err2
				}
				allLogs = append(allLogs, logs1...)
				allLogs = append(allLogs, logs2...)
				i++
				continue
			}
			if infiniteRetries {
				if errors.Is(err, context.Canceled) {
					return nil, err
				}
				c.recordRetry()
				c.logger.Warn("error getting logs, retrying", zap.Error(err))
				time.Sleep(10 * time.Second)
				continue
			}
			return nil, err
		}

		allLogs = append(allLogs, logs...)
		i++
	}

	allLogs = sortLogs(allLogs)
	return allLogs, nil
}

func sortLogs(logs []types.Log) []types.Log {
	sort.SliceStable(logs, func(i, j int) bool {
		left := logs[i]
		right := logs[j]
		if left.BlockNumber != right.BlockNumber {
			return left.BlockNumber < right.BlockNumber
		}
		if left.TxIndex != right.TxIndex {
			return left.TxIndex < right.TxIndex
		}
		return int32(left.Index) < int32(right.Index)
	})

	if len(logs) < 2 {
		return logs
	}

	write := 1
	for read := 1; read < len(logs); read++ {
		prev := logs[write-1]
		curr := logs[read]
		isDuplicate := prev.BlockNumber == curr.BlockNumber &&
			prev.TxIndex == curr.TxIndex &&
			prev.Index == curr.Index
		if isDuplicate {
			continue
		}
		logs[write] = curr
		write++
	}

	return logs[:write]
}

type LogResult struct {
	Log types.Log
	Err error
}

func (c *Client) StreamGetLogs(
	ctx context.Context,
	addresses [][]common.Address,
	topics [][][]common.Hash,
	startBlock *big.Int,
	endBlock *big.Int,
) <-chan LogResult {
	if len(addresses) != len(topics) {
		ch := make(chan LogResult, 1)
		ch <- LogResult{Err: fmt.Errorf("addresses and topics must have the same length")}
		close(ch)
		return ch
	}

	workersNum := 40
	promisesQueueSize := workersNum * 10
	batchSize := c.filterLogsBatchSize

	out := make(chan LogResult, workersNum*batchSize)

	innerCtx, cancel := context.WithCancel(ctx)

	type batchPromise struct {
		dataCh chan []types.Log
		errCh  chan error
	}

	workersSem := make(chan struct{}, workersNum)
	promisesQueue := make(chan batchPromise, promisesQueueSize)

	go func() {
		defer close(promisesQueue)

		curr := new(big.Int).Set(startBlock)
		for curr.Cmp(endBlock) <= 0 {
			lStart := new(big.Int).Set(curr)
			lEnd := new(big.Int).Add(lStart, big.NewInt(int64(batchSize-1)))
			if lEnd.Cmp(endBlock) > 0 {
				lEnd.Set(endBlock)
			}

			promise := batchPromise{
				dataCh: make(chan []types.Log, 1),
				errCh:  make(chan error, 1),
			}

			select {
			case <-innerCtx.Done():
				return
			case workersSem <- struct{}{}:
			}

			select {
			case <-innerCtx.Done():
				<-workersSem
				return
			case promisesQueue <- promise:
			}

			go func(s, e *big.Int, p batchPromise) {
				defer func() {
					<-workersSem
				}()

				c.logger.Info("getting logs", zap.Int64("start_block", s.Int64()), zap.Int64("end_block", e.Int64()))
				logs, err := c.GetLogs(innerCtx, addresses, topics, s, e, true)
				c.logger.Info("got logs", zap.Int("logs", len(logs)))
				if err != nil {
					select {
					case <-innerCtx.Done():
						return
					case p.errCh <- err:
					}
					return
				}
				select {
				case <-innerCtx.Done():
					return
				case p.dataCh <- logs:
				}
			}(lStart, lEnd, promise)

			curr.Add(lEnd, big.NewInt(1))
		}
	}()

	go func() {
		defer cancel()
		defer close(out)

		i := 0
		for {
			i++
			select {
			case <-innerCtx.Done():
				c.logger.Info("inner context done")
				return
			case promise, ok := <-promisesQueue:
				if !ok {
					c.logger.Info("promises queue closed")
					return
				}

				select {
				case <-innerCtx.Done():
					return
				case err := <-promise.errCh:
					out <- LogResult{Err: err}
					return
				case logs := <-promise.dataCh:
					for _, log := range logs {
						select {
						case <-innerCtx.Done():
							return
						case out <- LogResult{Log: log}:
						}
					}
				}
			}
		}
	}()

	return out
}
