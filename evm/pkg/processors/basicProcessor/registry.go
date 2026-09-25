package basicProcessor

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	evmclient "github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers/defi"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers/swaps"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/parsers/transfers"
	"github.com/monolit-network/xlayer-indexer/shared/db"
)

type Registry struct {
	logger         *zap.Logger
	parsersMapping map[string]map[string]parserWithSource
	mu             sync.RWMutex
	chain          string

	db db.Client

	ParsersMap map[string]parsers.Parser

	IdentifierMap map[parsers.Identifier]parsers.Parser

	FallbackSource string
	FallbackParser parsers.Parser

	// priorityIdentifiers are checked deterministically BEFORE the
	// IdentifierMap loop. Go map iteration order is random, so a tx that
	// pool Swap log in the same aggregator tx) was previously attributed
	// to a random one of them. Identifiers listed here always get their
	// parser attached; the map loop then proceeds as before (minus these).
	priorityIdentifiers []parsers.Identifier
	// exclusiveIdentifiers: a matching priority identifier from this set
	// suppresses the whole IdentifierMap pass (bundles must not be
	// re-parsed as one net swap by the uniswap identifiers).
	exclusiveIdentifiers map[parsers.Identifier]bool

	preprocessFuncs []requestPreprocessor
}

type parserWithSource struct {
	parser parsers.Parser
	Source string
}

func newParserRegistry(logger *zap.Logger, db db.Client, chain string, evmClient *evmclient.Client) *Registry {
	r := &Registry{
		logger:          logger.Named("evm-parser-registry"),
		parsersMapping:  make(map[string]map[string]parserWithSource),
		db:              db,
		chain:           chain,
		preprocessFuncs: allPreprocessFuncs[models.Chain(chain)],
	}

	r.initParser(evmClient, logger)

	if err := r.updateFromDB(); err != nil {
		logger.Error("error updating from db", zap.Error(err))
		panic(err)
	}

	go func() {
		for {
			time.Sleep(time.Minute)
			if err := r.updateFromDB(); err != nil {
				logger.Error("error updating from db", zap.Error(err))
			}
		}
	}()

	return r
}

func (r *Registry) initParser(evmClient *evmclient.Client, logger *zap.Logger) {
	allParsers := []parsers.Parser{
		transfers.NewERC20OrNativeTransfersParser(evmClient, logger),
		swaps.NewByTransfersParser(evmClient, logger),
		swaps.NewByTransfersParserV2(evmClient, logger),
		swaps.NewSimpleByTransfersParser(evmClient, logger),
		swaps.NewCurveParser(evmClient, logger),
		swaps.NewUniswapPoolV2Parser(evmClient, logger),
		swaps.NewUniswapPoolV3Parser(evmClient, logger),
		swaps.NewUniswapUniversalRouterParser(evmClient, logger),
		defi.NewDefiParser(evmClient, logger),
	}
	r.ParsersMap = make(map[string]parsers.Parser)
	for _, parser := range allParsers {
		r.ParsersMap[parser.Name()] = parser
	}

	aaIdentifier := swaps.NewERC4337Identifier(evmClient, logger)
	aaParser := swaps.NewERC4337BundleParser(evmClient, logger)
	r.ParsersMap[aaParser.Name()] = aaParser
	r.IdentifierMap = map[parsers.Identifier]parsers.Parser{
		transfers.NewERC20OrNativeTransfersIdentifier(evmClient, logger): nil,
		swaps.NewUniswapV4Identifier(evmClient, logger):                  nil,
		swaps.NewUniswapV3LogIdentifier(evmClient, logger):               nil,
		swaps.NewUniswapV2LogIdentifier(evmClient, logger):               nil,
		swaps.NewIziSwapLogIdentifier(evmClient, logger):                 nil,
		swaps.NewDodoRouteLogIdentifier(evmClient, logger):               nil,
		swaps.NewCurvePoolLogIdentifier(evmClient, logger):               nil,
		aaIdentifier:   nil,
	}
	r.priorityIdentifiers = []parsers.Identifier{aaIdentifier}
	r.exclusiveIdentifiers = map[parsers.Identifier]bool{aaIdentifier: true}
	for identifier := range r.IdentifierMap {
		r.IdentifierMap[identifier] = r.ParsersMap[identifier.GetParserName()]
	}

	r.FallbackParser = r.ParsersMap[defi.DefiParser{}.Name()]
	r.FallbackSource = "defi"
}

func (r *Registry) updateFromDB() error {
	mapping, err := r.db.Querier().GetAllEvmParsersForChain(context.Background(), r.chain)
	if err != nil {
		return fmt.Errorf("error getting all evm parsers for chain: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	oldLen := len(r.parsersMapping)
	r.parsersMapping = make(map[string]map[string]parserWithSource)
	for contractAddress, instructionHash := range mapping {
		for instructionHash, parserData := range instructionHash {
			parserName, source := parserData[0], parserData[1]
			parser, ok := r.ParsersMap[strings.ToLower(parserName)]
			if !ok {
				return fmt.Errorf("parser not found for name: %s", parserName)
			}
			if _, ok := r.parsersMapping[strings.ToLower(contractAddress)]; !ok {
				r.parsersMapping[strings.ToLower(contractAddress)] = make(map[string]parserWithSource)
			}
			r.parsersMapping[strings.ToLower(contractAddress)][strings.ToLower(instructionHash)] = parserWithSource{parser, source}
		}
	}
	if oldLen != len(r.parsersMapping) {
		r.logger.Info("updated parsers mapping", zap.Int("oldLen", oldLen), zap.Int("newLen", len(r.parsersMapping)))
	}
	return nil
}

func (r *Registry) addrKey(addr common.Address) string {
	return strings.ToLower(addr.Hex())
}

type LookupResult int

const (
	LookupResultCANotFound LookupResult = iota
	LookupResultInstructionHashNotFound
	LookupResultFound
	LookupResultIdentified
	LookupResultContractCreation
	LookupResultOnlyMinor
)

func (r *Registry) GetSource(addr common.Address) (string, bool) {
	key := r.addrKey(addr)
	r.mu.RLock()
	instrMapping, ok := r.parsersMapping[key]
	r.mu.RUnlock()
	if !ok {
		return "", false
	}

	for _, instr := range instrMapping {
		return instr.Source, true
	}
	return "", false
}

func (r *Registry) GetParserSimple(addr common.Address) (parsers.Parser, bool) {
	key := r.addrKey(addr)
	r.mu.RLock()
	instrMapping, ok := r.parsersMapping[key]
	r.mu.RUnlock()
	if !ok {
		return nil, false
	}
	var parser parsers.Parser
	for _, instr := range instrMapping {
		parser = instr.parser
		break
	}
	return parser, true
}

func (r *Registry) LookupParsersByTx(req parsers.ParseTxRequest) ([]parsers.Parser, []string, LookupResult) {
	if req.Tx.To() == nil {
		return nil, nil, LookupResultContractCreation
	}
	key := r.addrKey(*req.Tx.To())
	var res LookupResult
	resParsers := make([]parsers.Parser, 0)
	resSources := make([]string, 0)
	r.mu.RLock()
	instrMapping, ok := r.parsersMapping[key]
	r.mu.RUnlock()
	if !ok {
		res = LookupResultCANotFound
	}
	if ok {
		if len(req.Tx.Data()) == 0 {
			return nil, nil, LookupResultInstructionHashNotFound
		}

		instructionHash := models.SelectorFromTx(req.Tx)
		data, ok := instrMapping[instructionHash]
		if !ok {
			res = LookupResultInstructionHashNotFound
		} else {
			res = LookupResultFound
			resParsers = append(resParsers, data.parser)
			resSources = append(resSources, data.Source)
		}
	}

	// Deterministic pass: priority identifiers always contribute their
	// parser if they match, regardless of what else matches below.
	priorityMatched := false
	exclusiveMatched := false
	isPriority := make(map[parsers.Identifier]bool, len(r.priorityIdentifiers))
	for _, ident := range r.priorityIdentifiers {
		isPriority[ident] = true
		if ok, err := ident.CheckTx(req); err == nil && ok {
			res = LookupResultFound
			resParsers = append(resParsers, r.IdentifierMap[ident])
			resSources = append(resSources, ident.Source())
			priorityMatched = true
			if r.exclusiveIdentifiers[ident] {
				exclusiveMatched = true
			}
		} else if err != nil {
			r.logger.Error("error checking tx", zap.Error(err))
		}
	}
	if exclusiveMatched {
		return resParsers, resSources, LookupResultFound
	}

	for ident, parser := range r.IdentifierMap {
		if isPriority[ident] {
			continue
		}
		if ok, err := ident.CheckTx(req); err == nil && ok {
			res = LookupResultFound
			resParsers = append(resParsers, parser)
			resSources = append(resSources, ident.Source())
			if ident.ShouldStop() {
				return resParsers, resSources, res
			}
		} else if err != nil {
			r.logger.Error("error checking tx", zap.Error(err))
		}
	}

	if priorityMatched {
		// A priority identifier is a "major" match: no defi fallback,
		// same as the pre-existing ShouldStop early-return behavior.
		return resParsers, resSources, LookupResultFound
	}

	if res == LookupResultFound && len(resParsers) > 0 {
		return resParsers, resSources, LookupResultOnlyMinor
	}

	return resParsers, resSources, res
}

func (r *Registry) ParseTx(req parsers.ParseTxRequest, evmClient *evmclient.Client) ([]models.Event, error) {
	var err error

	for _, preprocessFunc := range r.preprocessFuncs {
		req, err = preprocessFunc(req)
		if err != nil {
			return nil, fmt.Errorf("error preprocessing tx request: %w", err)
		}
	}

	parser, sources, result := r.LookupParsersByTx(req)
	if result == LookupResultCANotFound || result == LookupResultInstructionHashNotFound || result == LookupResultOnlyMinor {
		parser = append(parser, r.FallbackParser)
		sources = append(sources, r.FallbackSource)
	}
	if result == LookupResultContractCreation {
		return nil, nil
	}

	var allRequiredData evmclient.RequiredDataTypes
	for _, parser := range parser {
		allRequiredData = allRequiredData.Add(parser.RequiredData())
	}

	errorEvent := models.ErrorEvent{Chain: r.chain, BlockNumber: req.BlockHeader.Number, TxHash: req.Tx.Hash()}

	allEvents := make([]models.Event, 0)
	for i, parser := range parser {
		req.Source = sources[i]
		events, err := parser.ParseTx(req)
		if err != nil {
			if parser.OptionalData() != 0 {
				if events, err = parser.ParseTx(req); err != nil {
					errorEvent.Error = err.Error()
					allEvents = append(allEvents, &errorEvent)
					continue
				}

				allEvents = append(allEvents, events...)
				continue
			}
			errorEvent.Error = err.Error()
			allEvents = append(allEvents, &errorEvent)
		}
		if len(events) > 0 {
			allEvents = append(allEvents, events...)
		}
	}

	return allEvents, nil
}
