package processors

import (
	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
)

type IndexerProcessor interface {
	QueueRequest(request parsers.ParseTxRequest)
	ProcessRequest(request parsers.ParseTxRequest) error
	Stop()
	Errors() chan error
	RequiredDataTypes() evmclient.RequiredDataTypes
	FlushSync()
}
