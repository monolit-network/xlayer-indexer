package mockProcessor

import (
	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"go.uber.org/zap"
)

type MockProcessor struct {
	logger            *zap.Logger
	errorsCh          chan error
	requiredDataTypes evmclient.RequiredDataTypes
}

func NewMockProcessor(logger *zap.Logger, requiredDataTypes evmclient.RequiredDataTypes) *MockProcessor {
	return &MockProcessor{
		logger:            logger,
		errorsCh:          make(chan error, 1024),
		requiredDataTypes: requiredDataTypes,
	}
}

func (p *MockProcessor) QueueRequest(request parsers.ParseTxRequest) {
}

func (p *MockProcessor) Stop() {
}

func (p *MockProcessor) Errors() chan error {
	return p.errorsCh
}

func (p *MockProcessor) RequiredDataTypes() evmclient.RequiredDataTypes {
	return p.requiredDataTypes
}

func (p *MockProcessor) FlushSync() {
}

func (p *MockProcessor) ProcessRequest(request parsers.ParseTxRequest) error {
	return nil
}
