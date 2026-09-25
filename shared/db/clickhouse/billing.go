package clickhouse

import (
	"context"
	"fmt"

	"github.com/monolit-network/xlayer-indexer/shared/db/clickhouse/queries"
	"github.com/monolit-network/xlayer-indexer/shared/models"
)

func (c *ClickhouseClient) CreateBillingLogsTable(ctx context.Context) error {
	if err := c.Conn.Exec(ctx, queries.CreateSchemaBackendClickhouseSQL); err != nil {
		return fmt.Errorf("failed to create backend schema: %w", err)
	}
	if err := c.Conn.Exec(ctx, queries.CreateBillingLogsTableSQL); err != nil {
		return fmt.Errorf("failed to create billing logs table: %w", err)
	}
	if err := c.Conn.Exec(ctx, queries.AddBillingLogsAdditionalDataColumnSQL); err != nil {
		return fmt.Errorf("failed to add billing logs additional data column: %w", err)
	}
	return nil
}

func (c *ClickhouseClient) InsertBillingLogsBatch(ctx context.Context, logs []models.BillingLog) error {
	if len(logs) == 0 {
		return nil
	}

	batch, err := c.Conn.PrepareBatch(ctx, queries.InsertBillingLogsSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare billing logs batch: %w", err)
	}
	for _, log := range logs {
		if err := batch.Append(
			log.UserID,
			log.APIKey,
			log.CreditsAmount,
			string(log.OperationType),
			log.Source,
			log.AdditionalData,
		); err != nil {
			return fmt.Errorf("failed to append billing log: %w", err)
		}
	}
	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send billing logs batch: %w", err)
	}
	return nil
}
