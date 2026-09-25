package clickhouse

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	evmmodels "github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/shared/db/clickhouse/queries"
	"github.com/monolit-network/xlayer-indexer/util/waitgroup"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (c *ClickhouseClient) CreateSchemaEvm(ctx context.Context) error {
	return c.Conn.Exec(ctx, queries.CreateSchemaEvmClickhouseSQL)
}

func (c *ClickhouseClient) CreateEvmPolymarketOrderEventsTable(ctx context.Context) error {
	return c.Conn.Exec(ctx, queries.CreateTableEvmPolymarketOrderEventsSQL)
}

func (c *ClickhouseClient) CreateEvmSwapEventsTable(ctx context.Context) error {
	return c.Conn.Exec(ctx, queries.CreateTableEvmSwapEventsSQL)
}

func (c *ClickhouseClient) CreateEvmTransferEventsTable(ctx context.Context) error {
	return c.Conn.Exec(ctx, queries.CreateTableEvmTransferEventsSQL)
}

func (c *ClickhouseClient) CreateEvmDefiEventsTable(ctx context.Context) error {
	return c.Conn.Exec(ctx, queries.CreateTableEvmDefiEventsSQL)
}

func (c *ClickhouseClient) CreateEvmErrorEventsTable(ctx context.Context) error {
	return c.Conn.Exec(ctx, queries.CreateTableEvmErrorEventsSQL)
}

func (c *ClickhouseClient) InsertEvmSwapEventsBatch(ctx context.Context, rows []evmSwapEventRow) error {
	if len(rows) == 0 {
		return nil
	}
	batch, err := c.Conn.PrepareBatch(ctx, queries.InsertEvmSwapEventsSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}
	for _, r := range rows {
		if r.Value.BitLen() > 256 {
			c.logger.Warn("value is too large", zap.String("chain", r.Chain), zap.String("tx_hash", r.TxHash), zap.String("amount", r.Value.String()), zap.Int("bit_len", r.Value.BitLen()))
			continue
		}
		if r.BaseCoinAmount.BitLen() > 256 {
			c.logger.Warn("base coin amount is too large", zap.String("chain", r.Chain), zap.String("tx_hash", r.TxHash), zap.String("amount", r.BaseCoinAmount.String()), zap.Int("bit_len", r.BaseCoinAmount.BitLen()))
			continue
		}
		if r.QuoteCoinAmount.BitLen() > 256 {
			c.logger.Warn("quote coin amount is too large", zap.String("chain", r.Chain), zap.String("tx_hash", r.TxHash), zap.String("amount", r.QuoteCoinAmount.String()), zap.Int("bit_len", r.QuoteCoinAmount.BitLen()))
			continue
		}
		if err := batch.Append(
			r.Chain,
			r.BlockTime,
			r.BlockNumber,
			r.BlockHash,
			r.TxIdx,
			r.TxHash,
			r.GasUsed,
			r.Value,
			r.TxFromAddress,
			r.TxToAddress,
			r.Source,
			r.InstructionHash,
			r.BaseCoin,
			r.QuoteCoin,
			r.BaseCoinAmount,
			r.BaseCoinDecimals,
			r.QuoteCoinAmount,
			r.QuoteCoinDecimals,
			r.Sender,
			r.Receiver,
		); err != nil {
			return fmt.Errorf("failed to append batch: %w", err)
		}
	}
	return batch.Send()
}

func (c *ClickhouseClient) InsertEvmTransferEventsBatch(ctx context.Context, rows []evmTransferEventRow) error {
	if len(rows) == 0 {
		return nil
	}
	batch, err := c.Conn.PrepareBatch(ctx, queries.InsertEvmTransferEventsSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}
	for _, r := range rows {
		if r.Amount.BitLen() > 256 {
			c.logger.Warn("amount is too large", zap.String("chain", r.Chain), zap.String("tx_hash", r.TxHash), zap.String("amount", r.Amount.String()), zap.Int("bit_len", r.Amount.BitLen()))
			continue
		}
		if err := batch.Append(
			r.Chain,
			r.BlockTime,
			r.BlockNumber,
			r.BlockHash,
			r.TxIdx,
			r.TxHash,
			r.GasUsed,
			r.Sender,
			r.Receiver,
			r.TokenAddress,
			r.TokenDecimals,
			r.Amount,
			r.Source,
			// TODO mb add AmountFloat
		); err != nil {
			return fmt.Errorf("failed to append batch: %w", err)
		}
	}
	return batch.Send()
}

func (c *ClickhouseClient) InsertEvmDefiEventsBatch(ctx context.Context, rows []evmDefiEventRow) error {
	if len(rows) == 0 {
		return nil
	}
	batch, err := c.Conn.PrepareBatch(ctx, queries.InsertEvmDefiEventsSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}
	for _, r := range rows {
		var foundTooBigAmount bool
		for _, amount := range r.InputAmounts {
			if amount.BitLen() > 256 {
				c.logger.Warn("input amount is too large", zap.String("chain", r.Chain), zap.String("tx_hash", r.TxHash), zap.String("amount", amount.String()), zap.Int("bit_len", amount.BitLen()))
				foundTooBigAmount = true
			}
		}
		for _, amount := range r.OutputAmounts {
			if amount.BitLen() > 256 {
				c.logger.Warn("output amount is too large", zap.String("chain", r.Chain), zap.String("tx_hash", r.TxHash), zap.String("amount", amount.String()), zap.Int("bit_len", amount.BitLen()))
				foundTooBigAmount = true
			}
		}
		if foundTooBigAmount {
			continue
		}
		if err := batch.Append(
			r.Chain,
			r.BlockTime,
			r.BlockNumber,
			r.BlockHash,
			r.TxIdx,
			r.TxHash,
			r.GasUsed,
			r.Value,
			r.TxFromAddress,
			r.TxToAddress,
			r.Source,
			r.InstructionHash,
			r.InputAddresses,
			r.InputAmounts,
			r.InputDecimals,
			r.OutputAddresses,
			r.OutputAmounts,
			r.OutputDecimals,
		); err != nil {
			return fmt.Errorf("failed to append batch: %w", err)
		}
	}
	return batch.Send()
}

func (c *ClickhouseClient) InsertEvmErrorEventsBatch(ctx context.Context, rows []evmErrorEventRow) error {
	if len(rows) == 0 {
		return nil
	}
	batch, err := c.Conn.PrepareBatch(ctx, queries.InsertEvmErrorEventsSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}
	for _, r := range rows {
		if err := batch.Append(
			r.Chain,
			r.BlockNumber,
			r.TxHash,
			r.Error,
		); err != nil {
			return fmt.Errorf("failed to append batch: %w", err)
		}
	}
	return batch.Send()
}

func (c *ClickhouseClient) InsertEvmPolymarketOrderEventsBatch(ctx context.Context, rows []evmPolymarketOrderEventRow) error {
	if len(rows) == 0 {
		return nil
	}
	batch, err := c.Conn.PrepareBatch(ctx, queries.InsertEvmPolymarketOrderEventsSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}
	for _, r := range rows {
		fee := decimal.NewFromInt(0)
		if r.Fee != nil {
			fee = decimal.NewFromBigInt(r.Fee, 0)
		}
		if err := batch.Append(
			r.BlockTime,
			r.BlockNumber,
			r.BlockHash,
			r.TxIdx,
			r.TxHash,
			r.TxFromAddress,
			r.TxToAddress,
			r.InstructionHash,
			r.ID,
			r.TransactionHash,
			r.OrderHash,
			r.Maker,
			r.Taker,
			r.MakerAssetID,
			r.TakerAssetID,
			r.MakerAmountFilled,
			r.TakerAmountFilled,
			fee,
			r.IsDeleted,
		); err != nil {
			return fmt.Errorf("failed to append batch: %w", err)
		}
	}
	return batch.Send()
}

func (c *ClickhouseClient) InsertEvmEventsBatch(ctx context.Context, chain models.Chain, events []evmmodels.Event) error {
	if len(events) == 0 {
		return nil
	}
	defiRows := make([]evmDefiEventRow, 0, len(events))
	transferRows := make([]evmTransferEventRow, 0, len(events))
	swapRows := make([]evmSwapEventRow, 0, len(events))
	errorRows := make([]evmErrorEventRow, 0, len(events))
	polymarketRows := make([]evmPolymarketOrderEventRow, 0, len(events))

	for _, ev := range events {
		switch e := ev.(type) {
		case *evmmodels.DefiEvent:
			row := (evmDefiEventRow{}).FromDefiEvent(chain, e)
			defiRows = append(defiRows, row)
		case *evmmodels.TransferEvent:
			row := (evmTransferEventRow{}).FromTransferEvent(chain, e)
			transferRows = append(transferRows, row)
		case *evmmodels.SwapEvent:
			row := (evmSwapEventRow{}).FromSwapEvent(chain, e)
			swapRows = append(swapRows, row)
		case *evmmodels.ErrorEvent:
			row := (evmErrorEventRow{}).FromErrorEvent(chain, e)
			errorRows = append(errorRows, row)
		case *evmmodels.PolymarketOrderEvent:
			row := (evmPolymarketOrderEventRow{}).FromPolymarketOrderEvent(e)
			polymarketRows = append(polymarketRows, row)
		}
	}

	var wg waitgroup.WaitGroup
	wg.Go(func() error {
		return c.InsertEvmDefiEventsBatch(ctx, defiRows)
	})
	wg.Go(func() error {
		return c.InsertEvmTransferEventsBatch(ctx, transferRows)
	})
	wg.Go(func() error {
		return c.InsertEvmSwapEventsBatch(ctx, swapRows)
	})
	wg.Go(func() error {
		return c.InsertEvmErrorEventsBatch(ctx, errorRows)
	})
	wg.Go(func() error {
		return c.InsertEvmPolymarketOrderEventsBatch(ctx, polymarketRows)
	})
	return wg.Wait()
}

func (c *ClickhouseClient) GetEvmSwapEvents(ctx context.Context, chain models.Chain, startBlock, endBlock *big.Int) ([]evmmodels.Event, error) {
	rows, err := c.Conn.Query(ctx, queries.SelectEvmSwapEventsSQL, chain, startBlock, endBlock)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []evmmodels.Event
	for rows.Next() {
		var (
			chainOut          string
			blockTime         time.Time
			blockNumber       big.Int
			blockHash         string
			txIdx             uint32
			txHash            string
			gasUsed           uint64
			value             big.Int
			txFrom            string
			txTo              string
			source            string
			instructionHash   string
			baseCoin          string
			quoteCoin         string
			baseCoinAmount    big.Int
			baseCoinDecimals  *uint8
			quoteCoinAmount   big.Int
			quoteCoinDecimals *uint8
			sender            string
			receiver          string
		)

		if err := rows.Scan(
			&chainOut,
			&blockTime,
			&blockNumber,
			&blockHash,
			&txIdx,
			&txHash,
			&gasUsed,
			&value,
			&txFrom,
			&txTo,
			&source,
			&instructionHash,
			&baseCoin,
			&quoteCoin,
			&baseCoinAmount,
			&baseCoinDecimals,
			&quoteCoinAmount,
			&quoteCoinDecimals,
			&sender,
			&receiver,
		); err != nil {
			return nil, err
		}

		r := evmSwapEventRow{
			Chain:             chainOut,
			BlockTime:         blockTime,
			BlockNumber:       new(big.Int).Set(&blockNumber),
			BlockHash:         blockHash,
			TxIdx:             txIdx,
			TxHash:            txHash,
			GasUsed:           gasUsed,
			Value:             new(big.Int).Set(&value),
			TxFromAddress:     txFrom,
			TxToAddress:       txTo,
			Source:            source,
			InstructionHash:   instructionHash,
			BaseCoin:          baseCoin,
			QuoteCoin:         quoteCoin,
			BaseCoinAmount:    new(big.Int).Set(&baseCoinAmount),
			BaseCoinDecimals:  baseCoinDecimals,
			QuoteCoinAmount:   new(big.Int).Set(&quoteCoinAmount),
			QuoteCoinDecimals: quoteCoinDecimals,
			Sender:            sender,
			Receiver:          receiver,
		}
		events = append(events, r.ToEvent())
	}

	return events, rows.Err()
}

func (c *ClickhouseClient) GetEvmSwapEventsChan(ctx context.Context, chain models.Chain, startBlock, endBlock *big.Int, buffer int) (<-chan evmmodels.Event, <-chan error, error) {
	if buffer <= 0 {
		buffer = 1
	}

	rows, err := c.Conn.Query(ctx, queries.SelectEvmSwapEventsSQL, chain, startBlock, endBlock)
	if err != nil {
		return nil, nil, err
	}

	eventCh := make(chan evmmodels.Event, buffer)
	errCh := make(chan error, 1)

	go func() {
		defer rows.Close()
		defer close(eventCh)
		defer close(errCh)

		for rows.Next() {
			var (
				chainOut          string
				blockTime         time.Time
				blockNumber       big.Int
				blockHash         string
				txIdx             uint32
				txHash            string
				gasUsed           uint64
				value             big.Int
				txFrom            string
				txTo              string
				source            string
				instructionHash   string
				baseCoin          string
				quoteCoin         string
				baseCoinAmount    big.Int
				baseCoinDecimals  *uint8
				quoteCoinAmount   big.Int
				quoteCoinDecimals *uint8
				sender            string
				receiver          string
			)

			if err := rows.Scan(
				&chainOut,
				&blockTime,
				&blockNumber,
				&blockHash,
				&txIdx,
				&txHash,
				&gasUsed,
				&value,
				&txFrom,
				&txTo,
				&source,
				&instructionHash,
				&baseCoin,
				&quoteCoin,
				&baseCoinAmount,
				&baseCoinDecimals,
				&quoteCoinAmount,
				&quoteCoinDecimals,
				&sender,
				&receiver,
			); err != nil {
				errCh <- err
				continue
			}

			r := evmSwapEventRow{
				Chain:             chainOut,
				BlockTime:         blockTime,
				BlockNumber:       new(big.Int).Set(&blockNumber),
				BlockHash:         blockHash,
				TxIdx:             txIdx,
				TxHash:            txHash,
				GasUsed:           gasUsed,
				Value:             new(big.Int).Set(&value),
				TxFromAddress:     txFrom,
				TxToAddress:       txTo,
				Source:            source,
				InstructionHash:   instructionHash,
				BaseCoin:          baseCoin,
				QuoteCoin:         quoteCoin,
				BaseCoinAmount:    new(big.Int).Set(&baseCoinAmount),
				BaseCoinDecimals:  baseCoinDecimals,
				QuoteCoinAmount:   new(big.Int).Set(&quoteCoinAmount),
				QuoteCoinDecimals: quoteCoinDecimals,
				Sender:            sender,
				Receiver:          receiver,
			}

			eventCh <- r.ToEvent()
		}

		if err := rows.Err(); err != nil {
			errCh <- err
		}
	}()

	return eventCh, errCh, nil
}

func (c *ClickhouseClient) GetEvmTransferEvents(ctx context.Context, chain models.Chain, startBlock, endBlock *big.Int) ([]evmmodels.Event, error) {
	rows, err := c.Conn.Query(ctx, queries.SelectEvmTransferEventsSQL, chain, startBlock, endBlock)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []evmmodels.Event
	for rows.Next() {
		var (
			chainOut      string
			blockTime     time.Time
			blockNumber   big.Int
			blockHash     string
			txIdx         uint32
			txHash        string
			gasUsed       uint64
			sender        string
			receiver      string
			tokenAddress  string
			tokenDecimals *uint8
			amount        big.Int
		)

		if err := rows.Scan(
			&chainOut,
			&blockTime,
			&blockNumber,
			&blockHash,
			&txIdx,
			&txHash,
			&gasUsed,
			&sender,
			&receiver,
			&tokenAddress,
			&tokenDecimals,
			&amount,
		); err != nil {
			return nil, err
		}

		r := evmTransferEventRow{
			Chain:         chainOut,
			BlockTime:     blockTime,
			BlockNumber:   new(big.Int).Set(&blockNumber),
			BlockHash:     blockHash,
			TxIdx:         txIdx,
			TxHash:        txHash,
			GasUsed:       gasUsed,
			Sender:        sender,
			Receiver:      receiver,
			TokenAddress:  tokenAddress,
			TokenDecimals: tokenDecimals,
			Amount:        new(big.Int).Set(&amount),
		}
		events = append(events, r.ToEvent())
	}

	return events, rows.Err()
}

func (c *ClickhouseClient) GetEvmTransferEventsChan(ctx context.Context, chain models.Chain, startBlock, endBlock *big.Int, buffer int) (<-chan evmmodels.Event, <-chan error, error) {
	if buffer <= 0 {
		buffer = 1
	}

	rows, err := c.Conn.Query(ctx, queries.SelectEvmTransferEventsSQL, chain, startBlock, endBlock)
	if err != nil {
		return nil, nil, err
	}

	eventCh := make(chan evmmodels.Event, buffer)
	errCh := make(chan error, 1)

	go func() {
		defer rows.Close()
		defer close(eventCh)
		defer close(errCh)

		for rows.Next() {
			var (
				chainOut      string
				blockTime     time.Time
				blockNumber   big.Int
				blockHash     string
				txIdx         uint32
				txHash        string
				gasUsed       uint64
				sender        string
				receiver      string
				tokenAddress  string
				tokenDecimals *uint8
				amount        big.Int
			)

			if err := rows.Scan(
				&chainOut,
				&blockTime,
				&blockNumber,
				&blockHash,
				&txIdx,
				&txHash,
				&gasUsed,
				&sender,
				&receiver,
				&tokenAddress,
				&tokenDecimals,
				&amount,
			); err != nil {
				errCh <- err
				continue
			}

			r := evmTransferEventRow{
				Chain:         chainOut,
				BlockTime:     blockTime,
				BlockNumber:   new(big.Int).Set(&blockNumber),
				BlockHash:     blockHash,
				TxIdx:         txIdx,
				TxHash:        txHash,
				GasUsed:       gasUsed,
				Sender:        sender,
				Receiver:      receiver,
				TokenAddress:  tokenAddress,
				TokenDecimals: tokenDecimals,
				Amount:        new(big.Int).Set(&amount),
			}

			eventCh <- r.ToEvent()
		}

		if err := rows.Err(); err != nil {
			errCh <- err
		}
	}()

	return eventCh, errCh, nil
}

func (c *ClickhouseClient) GetEvmDefiEvents(ctx context.Context, chain models.Chain, startBlock, endBlock *big.Int) ([]evmmodels.Event, error) {
	rows, err := c.Conn.Query(ctx, queries.SelectEvmDefiEventsSQL, chain, startBlock, endBlock)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []evmmodels.Event
	for rows.Next() {
		var (
			chainOut        string
			blockTime       time.Time
			blockNumber     big.Int
			blockHash       string
			txIdx           uint32
			txHash          string
			gasUsed         uint64
			value           big.Int
			txFrom          string
			txTo            string
			source          string
			instructionHash string
			inputAddresses  []string
			inputAmounts    []*big.Int
			inputDecimals   []uint8
			outputAddresses []string
			outputAmounts   []*big.Int
			outputDecimals  []uint8
		)

		if err := rows.Scan(
			&chainOut,
			&blockTime,
			&blockNumber,
			&blockHash,
			&txIdx,
			&txHash,
			&gasUsed,
			&value,
			&txFrom,
			&txTo,
			&source,
			&instructionHash,
			&inputAddresses,
			&inputAmounts,
			&inputDecimals,
			&outputAddresses,
			&outputAmounts,
			&outputDecimals,
		); err != nil {
			return nil, err
		}

		cloneAddresses := func(src []string) []string {
			if len(src) == 0 {
				return nil
			}
			dst := make([]string, len(src))
			copy(dst, src)
			return dst
		}
		cloneAmounts := func(src []*big.Int) []*big.Int {
			if len(src) == 0 {
				return nil
			}
			dst := make([]*big.Int, len(src))
			for i, amt := range src {
				if amt != nil {
					dst[i] = new(big.Int).Set(amt)
				}
			}
			return dst
		}
		cloneDecimals := func(src []uint8) []*uint8 {
			if len(src) == 0 {
				return nil
			}
			dst := make([]*uint8, len(src))
			for i, dec := range src {
				dst[i] = &dec
			}
			return dst
		}

		r := evmDefiEventRow{
			Chain:           chainOut,
			BlockTime:       blockTime,
			BlockNumber:     new(big.Int).Set(&blockNumber),
			BlockHash:       blockHash,
			TxIdx:           txIdx,
			TxHash:          txHash,
			GasUsed:         gasUsed,
			Value:           new(big.Int).Set(&value),
			TxFromAddress:   txFrom,
			TxToAddress:     txTo,
			Source:          source,
			InstructionHash: instructionHash,
			InputAddresses:  cloneAddresses(inputAddresses),
			InputAmounts:    cloneAmounts(inputAmounts),
			InputDecimals:   cloneDecimals(inputDecimals),
			OutputAddresses: cloneAddresses(outputAddresses),
			OutputAmounts:   cloneAmounts(outputAmounts),
			OutputDecimals:  cloneDecimals(outputDecimals),
		}
		events = append(events, r.ToEvent())
	}

	return events, rows.Err()
}

func (c *ClickhouseClient) GetEvmDefiEventsChan(ctx context.Context, chain models.Chain, contractAddress common.Address, instructionHash string, startBlock, endBlock *big.Int, buffer int) (<-chan evmmodels.Event, <-chan error, error) {
	if buffer <= 0 {
		buffer = 1
	}

	rows, err := c.Conn.Query(ctx, queries.SelectEvmDefiEventsSQL, chain, startBlock, endBlock, normalizeAddress(contractAddress), instructionHash)
	if err != nil {
		return nil, nil, err
	}

	eventCh := make(chan evmmodels.Event, buffer)
	errCh := make(chan error, 1)

	go func() {
		defer rows.Close()
		defer close(eventCh)
		defer close(errCh)

		for rows.Next() {
			var (
				chainOut        string
				blockTime       time.Time
				blockNumber     big.Int
				blockHash       string
				txIdx           uint32
				txHash          string
				gasUsed         uint64
				value           big.Int
				txFrom          string
				txTo            string
				source          string
				instructionHash string
				inputAddresses  []string
				inputAmounts    []*big.Int
				inputDecimals   []uint8
				outputAddresses []string
				outputAmounts   []*big.Int
				outputDecimals  []uint8
			)

			if err := rows.Scan(
				&chainOut,
				&blockTime,
				&blockNumber,
				&blockHash,
				&txIdx,
				&txHash,
				&gasUsed,
				&value,
				&txFrom,
				&txTo,
				&source,
				&instructionHash,
				&inputAddresses,
				&inputAmounts,
				&inputDecimals,
				&outputAddresses,
				&outputAmounts,
				&outputDecimals,
			); err != nil {
				errCh <- err
				continue
			}

			cloneAddresses := func(src []string) []string {
				if len(src) == 0 {
					return nil
				}
				dst := make([]string, len(src))
				copy(dst, src)
				return dst
			}
			cloneAmounts := func(src []*big.Int) []*big.Int {
				if len(src) == 0 {
					return nil
				}
				dst := make([]*big.Int, len(src))
				for i, amt := range src {
					if amt != nil {
						dst[i] = new(big.Int).Set(amt)
					}
				}
				return dst
			}
			cloneDecimals := func(src []uint8) []*uint8 {
				if len(src) == 0 {
					return nil
				}
				dst := make([]*uint8, len(src))
				for i, dec := range src {
					dst[i] = &dec
				}
				return dst
			}

			r := evmDefiEventRow{
				Chain:           chainOut,
				BlockTime:       blockTime,
				BlockNumber:     new(big.Int).Set(&blockNumber),
				BlockHash:       blockHash,
				TxIdx:           txIdx,
				TxHash:          txHash,
				GasUsed:         gasUsed,
				Value:           new(big.Int).Set(&value),
				TxFromAddress:   txFrom,
				TxToAddress:     txTo,
				Source:          source,
				InstructionHash: instructionHash,
				InputAddresses:  cloneAddresses(inputAddresses),
				InputAmounts:    cloneAmounts(inputAmounts),
				InputDecimals:   cloneDecimals(inputDecimals),
				OutputAddresses: cloneAddresses(outputAddresses),
				OutputAmounts:   cloneAmounts(outputAmounts),
				OutputDecimals:  cloneDecimals(outputDecimals),
			}

			eventCh <- r.ToEvent()
		}

		if err := rows.Err(); err != nil {
			errCh <- err
		}
	}()

	return eventCh, errCh, nil
}

func (c *ClickhouseClient) GetEvmErrorEvents(ctx context.Context, chain models.Chain, startBlock, endBlock *big.Int) ([]evmmodels.Event, error) {
	rows, err := c.Conn.Query(ctx, queries.SelectEvmErrorEventsSQL, chain, startBlock, endBlock)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []evmmodels.Event
	for rows.Next() {
		var (
			chainOut    string
			blockNumber big.Int
			txHash      string
			errStr      string
		)

		if err := rows.Scan(
			&chainOut,
			&blockNumber,
			&txHash,
			&errStr,
		); err != nil {
			return nil, err
		}

		r := evmErrorEventRow{
			Chain:       chainOut,
			BlockNumber: new(big.Int).Set(&blockNumber),
			TxHash:      txHash,
			Error:       errStr,
		}
		events = append(events, r.ToEvent())
	}

	return events, rows.Err()
}

func (c *ClickhouseClient) GetEvmErrorEventsChan(ctx context.Context, chain models.Chain, startBlock, endBlock *big.Int, buffer int) (<-chan evmmodels.Event, <-chan error, error) {
	if buffer <= 0 {
		buffer = 1
	}

	rows, err := c.Conn.Query(ctx, queries.SelectEvmErrorEventsSQL, chain, startBlock, endBlock)
	if err != nil {
		return nil, nil, err
	}

	eventCh := make(chan evmmodels.Event, buffer)
	errCh := make(chan error, 1)

	go func() {
		defer rows.Close()
		defer close(eventCh)
		defer close(errCh)

		for rows.Next() {
			var (
				chainOut    string
				blockNumber big.Int
				txHash      string
				errStr      string
			)

			if err := rows.Scan(
				&chainOut,
				&blockNumber,
				&txHash,
				&errStr,
			); err != nil {
				errCh <- err
				continue
			}

			r := evmErrorEventRow{
				Chain:       chainOut,
				BlockNumber: new(big.Int).Set(&blockNumber),
				TxHash:      txHash,
				Error:       errStr,
			}

			eventCh <- r.ToEvent()
		}

		if err := rows.Err(); err != nil {
			errCh <- err
		}
	}()

	return eventCh, errCh, nil
}

func (c *ClickhouseClient) DeleteEvmDefiEventsBatch(ctx context.Context, chain models.Chain, blockNumbers []*big.Int, txHashes []common.Hash) error {
	if len(txHashes) == 0 {
		return nil
	}
	bns := make([]uint64, len(blockNumbers))
	hashes := make([]string, len(txHashes))

	for i, h := range txHashes {
		bns[i] = blockNumbers[i].Uint64()
		hashes[i] = normalizeHash(h)
	}

	return c.Conn.Exec(ctx, queries.DeleteEvmDefiEventsSQL, chain, bns, hashes)
}

func (c *ClickhouseClient) DeleteEvmDefiEventsByContractAddress(ctx context.Context, chain models.Chain, contractAddress common.Address, instructionHash string) error {
	return c.Conn.Exec(ctx, queries.DeleteEvmDefiEventsByContractAddressSQL, chain, normalizeAddress(contractAddress), instructionHash)
}

// Delete by block_hash for reorg handling

func (c *ClickhouseClient) DeleteEvmSwapEventsByBlockHash(ctx context.Context, chain models.Chain, blockHash string) error {
	return c.Conn.Exec(ctx, queries.DeleteEvmSwapEventsByBlockHashSQL, chain, blockHash)
}

func (c *ClickhouseClient) DeleteEvmTransferEventsByBlockHash(ctx context.Context, chain models.Chain, blockHash string) error {
	return c.Conn.Exec(ctx, queries.DeleteEvmTransferEventsByBlockHashSQL, chain, blockHash)
}

func (c *ClickhouseClient) DeleteEvmDefiEventsByBlockHash(ctx context.Context, chain models.Chain, blockHash string) error {
	return c.Conn.Exec(ctx, queries.DeleteEvmDefiEventsByBlockHashSQL, chain, blockHash)
}

func (c *ClickhouseClient) DeleteEvmPolymarketOrderEventsByBlockHash(ctx context.Context, blockHash string) error {
	return c.Conn.Exec(ctx, queries.DeleteEvmPolymarketOrderEventsByBlockHashSQL, blockHash)
}

func (c *ClickhouseClient) DeleteEvmEventsByBlockHash(ctx context.Context, chain models.Chain, blockHash string) error {
	var wg waitgroup.WaitGroup
	wg.Go(func() error {
		return c.DeleteEvmSwapEventsByBlockHash(ctx, chain, blockHash)
	})
	wg.Go(func() error {
		return c.DeleteEvmTransferEventsByBlockHash(ctx, chain, blockHash)
	})
	wg.Go(func() error {
		return c.DeleteEvmDefiEventsByBlockHash(ctx, chain, blockHash)
	})
	wg.Go(func() error {
		return c.DeleteEvmPolymarketOrderEventsByBlockHash(ctx, blockHash)
	})
	return wg.Wait()
}

// Delete by block_number for reorg handling (>= blockNumber)

func (c *ClickhouseClient) DeleteEvmSwapEventsFromBlockNumber(ctx context.Context, chain models.Chain, blockNumber uint64) error {
	return c.Conn.Exec(ctx, queries.DeleteEvmSwapEventsFromBlockNumberSQL, chain, blockNumber)
}

func (c *ClickhouseClient) DeleteEvmTransferEventsFromBlockNumber(ctx context.Context, chain models.Chain, blockNumber uint64) error {
	return c.Conn.Exec(ctx, queries.DeleteEvmTransferEventsFromBlockNumberSQL, chain, blockNumber)
}

func (c *ClickhouseClient) DeleteEvmDefiEventsFromBlockNumber(ctx context.Context, chain models.Chain, blockNumber uint64) error {
	return c.Conn.Exec(ctx, queries.DeleteEvmDefiEventsFromBlockNumberSQL, chain, blockNumber)
}

func (c *ClickhouseClient) DeleteEvmPolymarketOrderEventsFromBlockNumber(ctx context.Context, blockNumber uint64) error {
	return c.Conn.Exec(ctx, queries.DeleteEvmPolymarketOrderEventsFromBlockNumberSQL, blockNumber)
}

func (c *ClickhouseClient) DeleteEvmEventsFromBlockNumber(ctx context.Context, chain models.Chain, blockNumber uint64) error {
	var wg waitgroup.WaitGroup
	wg.Go(func() error {
		return c.DeleteEvmSwapEventsFromBlockNumber(ctx, chain, blockNumber)
	})
	wg.Go(func() error {
		return c.DeleteEvmTransferEventsFromBlockNumber(ctx, chain, blockNumber)
	})
	wg.Go(func() error {
		return c.DeleteEvmDefiEventsFromBlockNumber(ctx, chain, blockNumber)
	})
	if chain == models.ChainPolygon {
		wg.Go(func() error {
			return c.DeleteEvmPolymarketOrderEventsFromBlockNumber(ctx, blockNumber)
		})
	}
	return wg.Wait()
}
