package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/monolit-network/xlayer-indexer/shared/db/postgres/queries"
	"github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func bigIntToNumeric(val *big.Int) pgtype.Numeric {
	if val == nil {
		return pgtype.Numeric{Valid: false}
	}
	return pgtype.Numeric{
		Int:   new(big.Int).Set(val),
		Exp:   0,
		Valid: true,
	}
}

func bigIntSliceToNumericSlice(vals []*big.Int) []pgtype.Numeric {
	result := make([]pgtype.Numeric, len(vals))
	for i, val := range vals {
		result[i] = bigIntToNumeric(val)
	}
	return result
}

func numericToBigInt(num pgtype.Numeric) *big.Int {
	if !num.Valid || num.Int == nil {
		return nil
	}
	return new(big.Int).Set(num.Int)
}

func numericSliceToBigIntSlice(nums []pgtype.Numeric) []*big.Int {
	result := make([]*big.Int, len(nums))
	for i, num := range nums {
		result[i] = numericToBigInt(num)
	}
	return result
}

func (c *PostgresDB) CreatePolymarketMarketsNewTable(ctx context.Context) error {
	_, err := c.querier().Exec(ctx, queries.CreatePolymarketMarketsNewTableSQL)
	return err
}

func (c *PostgresDB) CreatePolymarketTokensTable(ctx context.Context) error {
	_, err := c.querier().Exec(ctx, queries.CreatePolymarketTokensTableSQL)
	return err
}

func (c *PostgresDB) InsertPolymarketMarketNew(ctx context.Context, market models.PolymarketMarketNew) error {
	_, err := c.querier().Exec(ctx, queries.InsertPolymarketMarketNewSQL,
		strings.ToLower(market.ConditionID),
		market.QuestionID,
		strings.ToLower(market.Oracle),
		market.PreparedAt,
		market.PreparedInBlock,
		strings.ToLower(market.PreparedInBlockHash.Hex()),
		strings.ToLower(market.PreparedInTxHash.Hex()),
		market.TotalOutcomes,
		bigIntSliceToNumericSlice(market.PayoutNumerators),
	)
	return err
}

func (c *PostgresDB) UpdatePolymarketMarketResolution(ctx context.Context, conditionID string, resolvedAt time.Time, payoutNumerators []*big.Int, resolvedInBlock int64, resolvedInBlockHash string) error {
	_, err := c.querier().Exec(ctx, queries.UpdatePolymarketMarketResolutionSQL,
		conditionID,
		resolvedAt,
		bigIntSliceToNumericSlice(payoutNumerators),
		resolvedInBlock,
		resolvedInBlockHash,
	)
	return err
}

func (c *PostgresDB) UpdatePolymarketMarketResolutionFull(ctx context.Context, conditionID string, resolvedAt time.Time, payoutNumerators []*big.Int, resolvedInBlock int64, resolvedInBlockHash string) ([]models.PolymarketToken, error) {
	rows, err := c.querier().Query(ctx, queries.UpdatePolymarketMarketResolutionFullSQL,
		conditionID,
		resolvedAt,
		bigIntSliceToNumericSlice(payoutNumerators),
		resolvedInBlock,
		resolvedInBlockHash,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := make([]models.PolymarketToken, 0)
	for rows.Next() {
		token, err := scanPolymarketResolvedTokenRow(rows)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

func (c *PostgresDB) UpdatePolymarketTokensResolution(ctx context.Context, conditionID string, resolvedAt time.Time, payoutNumerators []*big.Int, resolvedInBlock int64, resolvedInBlockHash string) error {
	_, err := c.querier().Exec(ctx, queries.UpdatePolymarketTokensResolutionSQL,
		conditionID,
		resolvedAt,
		bigIntSliceToNumericSlice(payoutNumerators),
		resolvedInBlock,
		resolvedInBlockHash,
	)
	return err
}

func (c *PostgresDB) InsertPolymarketToken(ctx context.Context, token models.PolymarketToken, blockTime time.Time) error {
	_, err := c.querier().Exec(ctx, queries.InsertPolymarketTokenSQL,
		bigIntToNumeric(token.TokenID),
		token.ConditionID,
		token.CollateralToken,
		token.ParentCollectionID,
		bigIntToNumeric(token.Partition),
		blockTime,
	)
	return err
}

func (c *PostgresDB) InsertPolymarketTokensBatch(ctx context.Context, tokens []models.PolymarketToken) error {
	if len(tokens) == 0 {
		return nil
	}

	deduped := make([]models.PolymarketToken, 0, len(tokens))
	tokenIndexes := make(map[string]int, len(tokens))
	for _, token := range tokens {
		if token.TokenID == nil {
			deduped = append(deduped, token)
			continue
		}
		tokenID := token.TokenID.String()
		idx, ok := tokenIndexes[tokenID]
		if !ok {
			tokenIndexes[tokenID] = len(deduped)
			deduped = append(deduped, token)
			continue
		}
		if token.CreatedAt.Before(deduped[idx].CreatedAt) {
			deduped[idx] = token
		}
	}
	tokens = deduped

	valueStrings := make([]string, 0, len(tokens))
	valueArgs := make([]interface{}, 0, len(tokens)*6)

	for i, token := range tokens {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d::NUMERIC(78,0), $%d::TEXT, $%d::TEXT, $%d::TEXT, $%d::NUMERIC(78,0), $%d::TIMESTAMP)",
			i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))
		valueArgs = append(valueArgs, bigIntToNumeric(token.TokenID), token.ConditionID, token.CollateralToken, token.ParentCollectionID, bigIntToNumeric(token.Partition), token.CreatedAt)
	}

	query := fmt.Sprintf(queries.InsertPolymarketTokensBatchSQL, strings.Join(valueStrings, ","))
	_, err := c.querier().Exec(ctx, query, valueArgs...)
	return err
}

func (c *PostgresDB) GetPolymarketTokensByConditionID(ctx context.Context, conditionID string) ([]models.PolymarketToken, error) {
	rows, err := c.querier().Query(ctx, queries.GetPolymarketTokensByConditionIDSQL, conditionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := make([]models.PolymarketToken, 0)
	for rows.Next() {
		token, err := scanPolymarketTokenRow(rows)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

func (c *PostgresDB) GetPolymarketTokenByTokenID(ctx context.Context, tokenID *big.Int) (*models.PolymarketToken, error) {
	row := c.querier().QueryRow(ctx, queries.GetPolymarketTokenByTokenIDSQL, bigIntToNumeric(tokenID))

	token, err := scanPolymarketResolvedTokenRow(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &token, nil
}

func (c *PostgresDB) GetPolymarketMarketByConditionID(ctx context.Context, conditionID string) (*models.PolymarketMarketNew, error) {
	row := c.querier().QueryRow(ctx, queries.GetPolymarketMarketByConditionIDSQL, conditionID)

	var market models.PolymarketMarketNew
	var payoutNumeratorsNum []pgtype.Numeric

	if err := row.Scan(
		&market.ConditionID,
		&market.QuestionID,
		&market.Question,
		&market.Description,
		&market.AncillaryData,
		&market.Oracle,
		&market.TotalOutcomes,
		&market.IsResolved,
		&market.ResolvedAt,
		&payoutNumeratorsNum,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	market.PayoutNumerators = numericSliceToBigIntSlice(payoutNumeratorsNum)

	return &market, nil
}

func (c *PostgresDB) GetPolymarketMarketMinimalByConditionID(ctx context.Context, conditionID string) (*models.PolymarketMarketNewMinimal, error) {
	row := c.querier().QueryRow(ctx, queries.GetPolymarketMarketMinimalByConditionIDSQL, conditionID)

	var market models.PolymarketMarketNewMinimal
	if err := row.Scan(&market.ConditionID, &market.QuestionID); err != nil {
		return nil, err
	}
	return &market, nil
}

func (c *PostgresDB) GetPolymarketQuestionTimestampFromEvents(ctx context.Context, questionID string, maxBlockNumber uint64) (uint64, error) {
	row := c.querier().QueryRow(ctx, queries.GetPolymarketQuestionTimestampFromEventsSQL,
		strings.ToLower(questionID),
		maxBlockNumber,
	)

	var eventDataRaw []byte
	var eventType models.PolymarketMarketEventType
	if err := row.Scan(&eventDataRaw, &eventType); err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("no ADAPTER_INITIALIZED or ADAPTER_RESET event found for question %s up to block %d", questionID, maxBlockNumber)
		}
		return 0, fmt.Errorf("failed to query event: %w", err)
	}

	switch eventType {
	case models.PolymarketEventAdapterInitialized:
		var eventData models.PolymarketMarketEventDataAdapterInitialized
		if err := json.Unmarshal(eventDataRaw, &eventData); err != nil {
			return 0, fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		if eventData.RequestTimestamp == 0 {
			return 0, fmt.Errorf("request_timestamp is 0 in event data for question %s", questionID)
		}
		return eventData.RequestTimestamp, nil
	case models.PolymarketEventAdapterReset:
		var eventData models.PolymarketMarketEventDataAdapterReset
		if err := json.Unmarshal(eventDataRaw, &eventData); err != nil {
			return 0, fmt.Errorf("failed to unmarshal event data: %w", err)
		}
		if eventData.NewTimestamp == 0 {
			return 0, fmt.Errorf("new_timestamp is 0 in event data for question %s", questionID)
		}
		return eventData.NewTimestamp, nil
	default:
		return 0, fmt.Errorf("unsupported event type: %s", eventType)
	}
}

func (c *PostgresDB) UpdatePolymarketMarketAncillaryDataByQuestionID(ctx context.Context, questionID string, ancillaryData []byte, parsedAncillaryData models.ParsedAncillaryData) error {
	_, err := c.querier().Exec(ctx, queries.UpdatePolymarketMarketAncillaryDataByQuestionIDSQL,
		strings.ToLower(questionID),
		ancillaryData,
		parsedAncillaryData.Question,
		parsedAncillaryData.Description,
		parsedAncillaryData.ResData,
		parsedAncillaryData.MarketID,
		parsedAncillaryData.Initializer,
	)
	return err
}

func (c *PostgresDB) UpdatePolymarketMarketParsedAncillaryDataByQuestionID(ctx context.Context, questionID string, parsedAncillaryData models.ParsedAncillaryData) error {
	_, err := c.querier().Exec(ctx, queries.UpdatePolymarketMarketParsedAncillaryDataByQuestionIDSQL,
		strings.ToLower(questionID),
		parsedAncillaryData.Question,
		parsedAncillaryData.Description,
		parsedAncillaryData.ResData,
		parsedAncillaryData.MarketID,
		parsedAncillaryData.Initializer,
	)
	return err
}

func scanPolymarketTokenRow(rows pgx.Rows) (models.PolymarketToken, error) {
	var token models.PolymarketToken
	var tokenIDNum, partitionNum pgtype.Numeric

	if err := rows.Scan(
		&tokenIDNum,
		&token.ConditionID,
		&token.CollateralToken,
		&token.ParentCollectionID,
		&partitionNum,
	); err != nil {
		return models.PolymarketToken{}, err
	}

	token.TokenID = numericToBigInt(tokenIDNum)
	token.Partition = numericToBigInt(partitionNum)

	return token, nil
}

type polymarketTokenScanner interface {
	Scan(dest ...any) error
}

func scanPolymarketResolvedTokenRow(row polymarketTokenScanner) (models.PolymarketToken, error) {
	var token models.PolymarketToken
	var resolvedInBlock pgtype.Int8
	var resolvedInBlockHash pgtype.Text
	var tokenIDNum, partitionNum, numeratorNum, denominatorNum pgtype.Numeric

	if err := row.Scan(
		&tokenIDNum,
		&token.ConditionID,
		&token.CollateralToken,
		&token.ParentCollectionID,
		&partitionNum,
		&numeratorNum,
		&denominatorNum,
		&token.IsResolved,
		&resolvedInBlock,
		&resolvedInBlockHash,
	); err != nil {
		return models.PolymarketToken{}, err
	}

	token.TokenID = numericToBigInt(tokenIDNum)
	token.Partition = numericToBigInt(partitionNum)
	token.Numerator = numericToBigInt(numeratorNum)
	token.Denominator = numericToBigInt(denominatorNum)
	if resolvedInBlock.Valid {
		token.ResolvedInBlock = &resolvedInBlock.Int64
	}
	if resolvedInBlockHash.Valid {
		token.ResolvedInBlockHash = &resolvedInBlockHash.String
	}

	return token, nil
}

func (c *PostgresDB) CreatePolymarketMarketEventsTable(ctx context.Context) error {
	_, err := c.querier().Exec(ctx, queries.CreatePolymarketMarketEventsTableSQL)
	return err
}

func (c *PostgresDB) InsertPolymarketMarketEvent(ctx context.Context, event models.PolymarketMarketEvent) error {
	eventData, err := json.Marshal(event.EventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}
	_, err = c.querier().Exec(ctx, queries.InsertPolymarketMarketEventSQL,
		strings.ToLower(event.QuestionID),
		string(event.EventType),
		event.BlockNumber,
		strings.ToLower(event.BlockHash),
		int64(event.LogIndex),
		int64(event.TxIndex),
		strings.ToLower(event.TxHash),
		event.BlockTimestamp,
		eventData,
	)
	return err
}

func (c *PostgresDB) GetPolymarketMarketEventsByQuestionID(ctx context.Context, questionID string) ([]models.PolymarketMarketEvent, error) {
	rows, err := c.querier().Query(ctx, queries.GetPolymarketMarketEventsByQuestionIDSQL, strings.ToLower(questionID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]models.PolymarketMarketEvent, 0)
	for rows.Next() {
		var event models.PolymarketMarketEvent
		var eventType string
		dataRaw := []byte{}
		if err := rows.Scan(
			&event.QuestionID,
			&eventType,
			&event.BlockNumber,
			&event.BlockHash,
			&event.LogIndex,
			&event.TxIndex,
			&event.TxHash,
			&event.BlockTimestamp,
			&dataRaw,
		); err != nil {
			return nil, err
		}
		event.EventType = models.PolymarketMarketEventType(eventType)
		event.EventData, err = models.UnmarshalPolymarketMarketEventData(event.EventType, dataRaw)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

func (c *PostgresDB) UpdatePolymarketMarketGammaData(ctx context.Context, conditionID string, gammaInfo *models.PolymarketMarketEssentialGammaInfo) (bool, error) {
	result, err := c.querier().Exec(ctx, queries.UpdatePolymarketMarketGammaDataSQL,
		strings.ToLower(conditionID),
		gammaInfo.IsPresentedInGamma,
		gammaInfo.ID,
		gammaInfo.Question,
		gammaInfo.Description,
		gammaInfo.Slug,
		gammaInfo.EventSlug,
		gammaInfo.EventLink,
		gammaInfo.ResolutionSource,
		gammaInfo.StartDate,
		gammaInfo.EndDate,
		gammaInfo.CreatedAt,
		gammaInfo.UpdatedAt,
		gammaInfo.ClosedAt,
		gammaInfo.EventImageURL,
		gammaInfo.EventIconURL,
		gammaInfo.MarketImageURL,
		gammaInfo.MarketIconURL,
		gammaInfo.Outcomes,
		gammaInfo.GroupItemTitle,
		gammaInfo.AcceptingOrders,
		gammaInfo.OrderPriceMinTickSize,
		gammaInfo.OrderMinSize,
		gammaInfo.TagSlugs,
		gammaInfo.NegRisk,
		gammaInfo.NegRiskRequestID,
		gammaInfo.NegRiskOther,
		gammaInfo.UmaResolutionStatus,
		gammaInfo.UmaResolutionStatuses,
		gammaInfo.UmaBond,
		gammaInfo.UmaReward,
		gammaInfo.FeesEnabled,
		gammaInfo.Fee,
		gammaInfo.RawResponse,
	)
	if err != nil {
		return false, err
	}

	rowsAffected := result.RowsAffected()
	return rowsAffected > 0, nil
}

func (c *PostgresDB) UpdatePolymarketMarketGammaDataBatch(
	ctx context.Context,
	conditionIDs []string,
	gammaInfos []*models.PolymarketMarketEssentialGammaInfo,
) (int64, error) {
	if len(conditionIDs) != len(gammaInfos) {
		return 0, fmt.Errorf("conditionIDs and gammaInfos slices must have the same length")
	}

	if len(conditionIDs) == 0 {
		return 0, nil
	}

	batch := &pgx.Batch{}

	for i := range conditionIDs {
		batch.Queue(queries.UpdatePolymarketMarketGammaDataSQL,
			strings.ToLower(conditionIDs[i]),
			gammaInfos[i].IsPresentedInGamma,
			gammaInfos[i].ID,
			gammaInfos[i].Question,
			gammaInfos[i].Description,
			gammaInfos[i].Slug,
			gammaInfos[i].EventSlug,
			gammaInfos[i].EventLink,
			gammaInfos[i].ResolutionSource,
			gammaInfos[i].StartDate,
			gammaInfos[i].EndDate,
			gammaInfos[i].CreatedAt,
			gammaInfos[i].UpdatedAt,
			gammaInfos[i].ClosedAt,
			gammaInfos[i].EventImageURL,
			gammaInfos[i].EventIconURL,
			gammaInfos[i].MarketImageURL,
			gammaInfos[i].MarketIconURL,
			gammaInfos[i].Outcomes,
			gammaInfos[i].GroupItemTitle,
			gammaInfos[i].AcceptingOrders,
			gammaInfos[i].OrderPriceMinTickSize,
			gammaInfos[i].OrderMinSize,
			gammaInfos[i].TagSlugs,
			gammaInfos[i].NegRisk,
			gammaInfos[i].NegRiskRequestID,
			gammaInfos[i].NegRiskOther,
			gammaInfos[i].UmaResolutionStatus,
			gammaInfos[i].UmaResolutionStatuses,
			gammaInfos[i].UmaBond,
			gammaInfos[i].UmaReward,
			gammaInfos[i].FeesEnabled,
			gammaInfos[i].Fee,
			// gammaInfos[i].RawResponse,
		)
	}

	br := c.querier().SendBatch(ctx, batch)
	defer br.Close()

	var totalAffected int64
	for i := 0; i < len(conditionIDs); i++ {
		ct, err := br.Exec()
		if err != nil {
			return totalAffected, fmt.Errorf("error executing batch at index %d: %w", i, err)
		}
		totalAffected += ct.RowsAffected()
	}

	return totalAffected, nil
}

// Reorg handling methods

func (c *PostgresDB) DeletePolymarketMarketsByBlockHash(ctx context.Context, blockHash string) (int64, error) {
	result, err := c.querier().Exec(ctx, queries.DeletePolymarketMarketsByBlockHashSQL, blockHash)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (c *PostgresDB) DeletePolymarketMarketEventsByBlockHash(ctx context.Context, blockHash string) (int64, error) {
	result, err := c.querier().Exec(ctx, queries.DeletePolymarketMarketEventsByBlockHashSQL, blockHash)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (c *PostgresDB) ResetPolymarketTokensResolutionByBlockHash(ctx context.Context, blockHash string) (int64, error) {
	result, err := c.querier().Exec(ctx, queries.ResetPolymarketTokensResolutionByBlockHashSQL, blockHash)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// Reorg handling by block_number

func (c *PostgresDB) DeletePolymarketMarketsFromBlockNumber(ctx context.Context, blockNumber int64) (int64, error) {
	result, err := c.querier().Exec(ctx, queries.DeletePolymarketMarketsFromBlockNumberSQL, blockNumber)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (c *PostgresDB) DeletePolymarketMarketEventsFromBlockNumber(ctx context.Context, blockNumber int64) (int64, error) {
	result, err := c.querier().Exec(ctx, queries.DeletePolymarketMarketEventsFromBlockNumberSQL, blockNumber)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (c *PostgresDB) ResetPolymarketTokensResolutionFromBlockNumber(ctx context.Context, blockNumber int64) (int64, error) {
	result, err := c.querier().Exec(ctx, queries.ResetPolymarketTokensResolutionFromBlockNumberSQL, blockNumber)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (c *PostgresDB) ResetPolymarketMarketsResolutionByBlockHash(ctx context.Context, blockHash string) (int64, error) {
	result, err := c.querier().Exec(ctx, queries.ResetPolymarketMarketsResolutionByBlockHashSQL, blockHash)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (c *PostgresDB) ResetPolymarketMarketsResolutionFromBlockNumber(ctx context.Context, blockNumber int64) (int64, error) {
	result, err := c.querier().Exec(ctx, queries.ResetPolymarketMarketsResolutionFromBlockNumberSQL, blockNumber)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
