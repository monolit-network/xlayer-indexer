package main

import (
	"context"
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers/polymarket"
	"github.com/monolit-network/xlayer-indexer/shared/db/postgres"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const BatchSize = 500

func main() {
	godotenv.Load()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbClient, err := postgres.NewPostgresClientFromEnv()
	if err != nil {
		logger.Fatal("error creating database client", zap.Error(err))
	}
	defer dbClient.Close()

	parser := polymarket.NewCTFParser(nil, nil)
	ctx := context.Background()

	rows, err := dbClient.Querier().Query(ctx, `
		SELECT question_id, ancillary_data
		FROM polymarket.markets 
		WHERE ancillary_data IS NOT NULL AND ancillary_data != ''
	`)
	if err != nil {
		logger.Fatal("failed to query markets new", zap.Error(err))
	}
	defer rows.Close()

	var (
		ids   []string
		ques  []*string
		descs []*string
		resD  []*string
		mIDs  []*int64
		inits []*string
		count int
		total int
		start = time.Now()
	)

	logger.Info("Starting re-parsing migration...")

	for rows.Next() {
		var qID string
		var rawData []byte
		if err := rows.Scan(&qID, &rawData); err != nil {
			logger.Error("scan error", zap.Error(err))
			continue
		}

		parsed := parser.ParseAncillaryData(rawData)

		ids = append(ids, qID)
		ques = append(ques, parsed.Question)
		descs = append(descs, parsed.Description)
		resD = append(resD, parsed.ResData)
		mIDs = append(mIDs, parsed.MarketID)
		inits = append(inits, parsed.Initializer)

		count++

		if count >= BatchSize {
			if err := bulkUpdate(ctx, dbClient, ids, ques, descs, resD, mIDs, inits); err != nil {
				logger.Error("bulk update error", zap.Error(err))
			}
			total += count
			logger.Info("Processed...", zap.Int("total", total), zap.Duration("elapsed", time.Since(start)))

			ids, ques, descs, resD, mIDs, inits = nil, nil, nil, nil, nil, nil
			count = 0
		}
	}

	if count > 0 {
		bulkUpdate(ctx, dbClient, ids, ques, descs, resD, mIDs, inits)
		total += count
	}

	logger.Info("Migration finished", zap.Int("total", total), zap.Duration("duration", time.Since(start)))
}

func bulkUpdate(ctx context.Context, db *postgres.PostgresClient, ids []string, ques, descs, resD []*string, mIDs []*int64, inits []*string) error {
	query := `
		UPDATE polymarket.markets AS m SET
			anc_question = t.q,
			anc_description = t.d,
			anc_res_data = t.rd,
			anc_market_id = t.mid,
			anc_initializer = t.init
		FROM (
			SELECT 
				UNNEST($1::text[]) as qid,
				UNNEST($2::text[]) as q,
				UNNEST($3::text[]) as d,
				UNNEST($4::text[]) as rd,
				UNNEST($5::bigint[]) as mid,
				UNNEST($6::text[]) as init
		) AS t
		WHERE m.question_id = t.qid
	`

	_, err := db.Querier().Exec(ctx, query, ids, ques, descs, resD, mIDs, inits)
	return err
}
