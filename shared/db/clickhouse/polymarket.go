package clickhouse

import (
	"context"
	"fmt"
	"strings"

	evmmodels "github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/shared/db/clickhouse/queries"
	"github.com/ethereum/go-ethereum/common"
)

func (c *ClickhouseClient) InsertEvmPolymarketOrderEventsBatchNew(ctx context.Context, events []*evmmodels.PolymarketOrderEventNew) error {
	if len(events) == 0 {
		return nil
	}
	batch, err := c.Conn.PrepareBatch(ctx, queries.InsertEvmPolymarketOrderEventsNewSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}
	for _, e := range events {
		sourceAddress := common.Address{}.Hex()
		if e.SourceAddress != nil {
			sourceAddress = e.SourceAddress.Hex()
		}
		if err := batch.Append(
			e.BlockTime,
			e.BlockNumber,
			strings.ToLower(e.BlockHash.Hex()),
			e.TxIdx,
			strings.ToLower(e.TxHash.Hex()),
			strings.ToLower(e.ConditionID.Hex()),
			e.LogIndex,
			e.SubIndex,
			strings.ToLower(e.UserAddress.Hex()),
			string(e.Source),
			strings.ToLower(sourceAddress),
			e.TokenID,
			e.TokenAmountDiff,
			e.UsdcAmountDiff,
			e.IsTaker,
			e.Fee,
		); err != nil {
			return fmt.Errorf("failed to append batch: %w", err)
		}
	}
	return batch.Send()
}
