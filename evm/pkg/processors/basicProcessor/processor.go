package basicProcessor

import (
	"context"
	"sync"
	"time"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	evmmetrics "github.com/monolit-network/xlayer-indexer/evm/pkg/metrics"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"github.com/monolit-network/xlayer-indexer/shared/db"
	sharedmodels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/monolit-network/xlayer-indexer/util/notifier"
	"go.uber.org/zap"
)

type Processor struct {
	logger    *zap.Logger
	registry  *Registry
	dbClick   db.DBClick
	eventsCh  chan models.Event
	evmClient *evmclient.Client
	notifier  *notifier.Notifier

	requests  chan parsers.ParseTxRequest
	errorsCh  chan error
	workersWg *sync.WaitGroup
	writerWg  *sync.WaitGroup

	flushCh chan chan struct{}

	requiredDataTypes evmclient.RequiredDataTypes

	chain models.Chain
}

func NewProcessor(logger *zap.Logger, dbClient db.Client, dbClick db.DBClick, evmClient *evmclient.Client, notif *notifier.Notifier, chain string) *Processor {
	registry := newParserRegistry(logger, dbClient, chain, evmClient)
	c := &Processor{
		logger:            logger,
		registry:          registry,
		dbClick:           dbClick,
		evmClient:         evmClient,
		notifier:          notif,
		chain:             models.Chain(chain),
		errorsCh:          make(chan error, 1024),
		eventsCh:          make(chan models.Event, 16384),
		requests:          make(chan parsers.ParseTxRequest, 16384),
		flushCh:           make(chan chan struct{}),
		requiredDataTypes: evmclient.RequiredDataTypesTx | evmclient.RequiredDataTypesReceipt | evmclient.RequiredDataTypesTraces | evmclient.RequiredDataTypesBlockHeader,
	}

	c.workersWg = &sync.WaitGroup{}
	c.writerWg = &sync.WaitGroup{}

	c.writerWg.Go(c.dbWriter)

	for i := 0; i < 500; i++ {
		c.workersWg.Go(func() {
			for request := range c.requests {
				events, err := c.registry.ParseTx(request, c.evmClient)
				if err != nil {
					c.errorsCh <- err
					continue
				}
				c.enqueueEvents(events)
			}
		})
	}
	return c
}

func (c *Processor) Errors() chan error {
	return c.errorsCh
}

func (c *Processor) QueueRequest(request parsers.ParseTxRequest) {
	c.requests <- request
}

func (c *Processor) ProcessRequest(request parsers.ParseTxRequest) error {
	events, err := c.registry.ParseTx(request, c.evmClient)
	if err != nil {
		return err
	}
	c.enqueueEvents(events)
	return nil
}

func (c *Processor) Stop() {
	close(c.requests)
	c.workersWg.Wait()
	close(c.eventsCh)
	c.writerWg.Wait()
	close(c.errorsCh)
}

func (c *Processor) dbWriter() {
	flush := func(batch []models.Event) []models.Event {
		if len(batch) == 0 {
			return batch
		}
		counts := eventTableCounts(batch)
		if err := c.dbClick.InsertEvmEventsBatch(context.Background(), c.chain, batch); err != nil {
			c.logger.Error("error inserting events batch", zap.Error(err))
			for table := range counts {
				evmmetrics.InsertErrors.WithLabelValues(string(c.chain), "basic", "clickhouse", table).Inc()
			}
			time.Sleep(time.Second * 5)
			return batch
		}
		for table, count := range counts {
			evmmetrics.InsertedRows.WithLabelValues(string(c.chain), "basic", "clickhouse", table).Add(float64(count))
		}
		return batch[:0]
	}

	batch := make([]models.Event, 0, 4096)
	for {
		select {
		case event, ok := <-c.eventsCh:
			if !ok {
				flush(batch)
				return
			}
			batch = append(batch, event)
			if len(batch) >= 4096 {
				batch = flush(batch)
			}
		case done := <-c.flushCh:
			batch = flush(batch)
			close(done)
		}
	}
}

func eventTableCounts(events []models.Event) map[string]int {
	counts := make(map[string]int, 4)
	for _, event := range events {
		switch event.(type) {
		case *models.SwapEvent:
			counts["evm.swap_events"]++
		case *models.TransferEvent:
			counts["evm.transfer_events"]++
		case *models.DefiEvent:
			counts["evm.defi_events"]++
		case *models.ErrorEvent:
			counts["evm.error_events"]++
		case *models.PolymarketOrderEvent:
			counts["evm.polymarket_order_events"]++
		}
	}
	return counts
}

func (c *Processor) FlushSync() {
	done := make(chan struct{})
	c.flushCh <- done
	<-done
}

func (c *Processor) RequiredDataTypes() evmclient.RequiredDataTypes {
	return c.requiredDataTypes
}

func (c *Processor) enqueueEvents(events []models.Event) {
	for _, event := range events {
		c.eventsCh <- event
		c.publishEvent(event)
	}
}

func (c *Processor) publishEvent(event models.Event) {
	if c.notifier == nil {
		return
	}

	ctx := context.Background()
	switch typed := event.(type) {
	case *models.SwapEvent:
		if err := notifier.Publish(c.notifier, ctx, string(sharedmodels.NotifierEventEVMSwap), *typed); err != nil {
			c.logger.Warn("failed to publish evm swap event", zap.String("topic", string(sharedmodels.NotifierEventEVMSwap)), zap.Error(err))
		}
	case *models.TransferEvent:
		if err := notifier.Publish(c.notifier, ctx, string(sharedmodels.NotifierEventEVMTransfer), *typed); err != nil {
			c.logger.Warn("failed to publish evm transfer event", zap.String("topic", string(sharedmodels.NotifierEventEVMTransfer)), zap.Error(err))
		}
	case *models.DefiEvent:
		if err := notifier.Publish(c.notifier, ctx, string(sharedmodels.NotifierEventEVMDefi), *typed); err != nil {
			c.logger.Warn("failed to publish evm defi event", zap.String("topic", string(sharedmodels.NotifierEventEVMDefi)), zap.Error(err))
		}
	}
}
