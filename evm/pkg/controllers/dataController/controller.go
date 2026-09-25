package dataController

import (
	"context"
	"math/big"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/shared/db"
	"go.uber.org/zap"
)

type DataController struct {
	logger    *zap.Logger
	clickDB   db.DBClick
	db        db.Client
	evmClient *evmclient.Client
	chain     string
}

func NewDataController(
	logger *zap.Logger,
	clickDB db.DBClick,
	db db.Client,
	evmClient *evmclient.Client,
	chain string,
) *DataController {
	return &DataController{
		logger:    logger,
		clickDB:   clickDB,
		db:        db,
		evmClient: evmClient,
		chain:     chain,
	}
}

func (c *DataController) SubscribeBlocks(ctx context.Context, requiredDataTypes evmclient.RequiredDataTypes) (chan *evmclient.BlocksInfo, error) {
	sub, err := c.evmClient.SubscribeHeads(ctx)
	if err != nil {
		c.logger.Error("error subscribing to heads", zap.Error(err))
		return nil, err
	}

	blockInfosCh := make(chan *evmclient.BlocksInfo, 64)

	go func() {
		defer close(blockInfosCh)

		for {
			select {
			case <-ctx.Done():
				return
			case head := <-sub:
				if head == nil {
					c.logger.Error("head is nil")
					continue
				}

				blockInfos := c.evmClient.BatchGetBlocksInfo(ctx, big.NewInt(head.Number.Int64()), big.NewInt(head.Number.Int64()), requiredDataTypes)
				blockInfosCh <- blockInfos
			}
		}
	}()

	return blockInfosCh, nil
}
