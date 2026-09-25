package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/monolit-network/xlayer-indexer/shared/db/postgres/queries"
	"github.com/monolit-network/xlayer-indexer/shared/models"
)

func (q *PostgresDB) CreateInvestorsTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, queries.CreateInvestorsTableSQL)
	if err != nil {
		return err
	}

	if _, err := q.getNextInvestorID(ctx); err != nil {
		return fmt.Errorf("error getting next investor id: %v", err)
	}
	return nil
}

func (q *PostgresDB) CreateInvestRoundsTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, queries.CreateInvestRoundsTableSQL)
	return err
}

func (q *PostgresDB) getNextInvestorID(ctx context.Context) (models.InvestorID, error) {
	var id models.InvestorID
	row := q.querier().QueryRow(ctx, queries.SelectNextInvestorIDSQL)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func marshalInvestor(inv models.Investor) (links []byte, roundCat []byte, roundStage []byte, portfolioIDs []int32, err error) {
	links, err = json.Marshal(inv.Links)
	if err != nil {
		return
	}
	roundCat, err = json.Marshal(inv.RoundDistributionByCategory)
	if err != nil {
		return
	}
	roundStage, err = json.Marshal(inv.RoundDistributionByStage)
	if err != nil {
		return
	}
	if len(inv.PortfolioCoins) > 0 {
		portfolioIDs = make([]int32, len(inv.PortfolioCoins))
		for i, v := range inv.PortfolioCoins {
			portfolioIDs[i] = int32(v)
		}
	}
	return
}

func (q *PostgresDB) UpsertInvestor(ctx context.Context, inv models.Investor) error {
	if inv.ID == 0 {
		nid, err := q.getNextInvestorID(ctx)
		if err != nil {
			return err
		}
		inv.ID = nid
	}

	links, roundCat, roundStage, portfolioIDs, err := marshalInvestor(inv)
	if err != nil {
		return err
	}
	_, err = q.querier().Exec(ctx, queries.UpsertInvestorSQL,
		inv.ID,
		inv.Slug,
		inv.Name,
		inv.Image,
		inv.Rank,
		inv.Country,
		inv.Description,
		inv.VentureType,
		inv.Tier,
		inv.LeadInvestments,
		links,
		portfolioIDs,
		inv.PortfolioCoinsSlugs,
		roundCat,
		roundStage,
	)
	return err
}

func (q *PostgresDB) GetExistingInvestorSlugs(ctx context.Context, slugs []string) ([]string, error) {
	rows, err := q.querier().Query(ctx, queries.SelectExistingInvestorSlugsSQL, slugs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		if err == pgx.ErrNoRows {
			return []string{}, nil
		}
		return nil, err
	}

	var existing []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		existing = append(existing, s)
	}
	return existing, nil
}

func (q *PostgresDB) GetInvestorIdsBySlugs(ctx context.Context, slugs []string) (map[string]models.InvestorID, error) {
	if len(slugs) == 0 {
		return map[string]models.InvestorID{}, nil
	}
	rows, err := q.querier().Query(ctx, queries.SelectInvestorIdsBySlugsSQL, slugs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		if err == pgx.ErrNoRows {
			return map[string]models.InvestorID{}, nil
		}
		return nil, err
	}

	res := make(map[string]models.InvestorID, len(slugs))
	for rows.Next() {
		var id models.InvestorID
		var slug string
		if err := rows.Scan(&id, &slug); err != nil {
			return nil, err
		}
		res[slug] = id
	}
	return res, nil
}

func (q *PostgresDB) UpsertInvestRound(ctx context.Context, round models.InvestRound) error {
	investors := make([]int32, 0, len(round.Investors))
	if len(round.Investors) > 0 {
		for _, v := range round.Investors {
			investors = append(investors, int32(v))
		}
	}

	leadInvestors := make([]int32, 0, len(round.LeadInvestors))
	if len(round.LeadInvestors) > 0 {
		for _, v := range round.LeadInvestors {
			leadInvestors = append(leadInvestors, int32(v))
		}
	}

	if round.InvestorsSlugs == nil {
		round.InvestorsSlugs = []string{}
	}
	if round.LeadInvestorsSlugs == nil {
		round.LeadInvestorsSlugs = []string{}
	}

	_, err := q.querier().Exec(ctx, queries.UpsertInvestRoundSQL,
		round.ID,
		round.AssetID,
		round.FundsRaised,
		round.PreValuation,
		round.PreValuationInaccurate,
		round.Stage,
		round.Category,
		round.Date,
		investors,
		leadInvestors,
		round.InvestorsSlugs,
		round.LeadInvestorsSlugs,
	)
	return err
}

func (q *PostgresDB) GetExistingInvestRoundIDs(ctx context.Context, ids []models.InvestRoundID) ([]models.InvestRoundID, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	conv := make([]int64, len(ids))
	for i, v := range ids {
		conv[i] = int64(v)
	}
	rows, err := q.querier().Query(ctx, queries.SelectExistingInvestRoundIDsSQL, conv)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var existing []models.InvestRoundID
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		existing = append(existing, models.InvestRoundID(id))
	}
	return existing, nil
}
