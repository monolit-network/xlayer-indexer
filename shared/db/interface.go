package db

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/qdrant/go-client/qdrant"

	evmmodels "github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/shared/models"
)

var ErrNoRows = pgx.ErrNoRows

type Client interface {
	Querier() DB
	Transaction(ctx context.Context, fn func(ctx context.Context, tx DB) error) error
	WithRetry(ctx context.Context, maxRetries int, fn func(ctx context.Context, q DB) error) error
	TransactionWithRetry(ctx context.Context, maxRetries int, fn func(ctx context.Context, tx DB) error) error
}

type DB interface {
	// Table creation methods
	CreateUnverifiedTokensTable(ctx context.Context) error
	CreateVerifiedTokensTable(ctx context.Context) error
	CreateDropstabInfoTable(ctx context.Context) error
	CreatePolymarketMarketsNewTable(ctx context.Context) error

	// Unverified tokens
	GetUnverifiedTokenByChainContractAddress(ctx context.Context, chain string, contractAddress string) (*models.Token, error)
	GetUnverifiedTokensBySymbol(ctx context.Context, symbol string) ([]models.Token, error)
	GetUnverifiedTokensByContractAddress(ctx context.Context, contractAddress string) ([]models.Token, error)
	GetUnverifiedTokensByContractAddresses(ctx context.Context, addresses []string) ([]models.Token, error)
	UpsertUnverifiedToken(ctx context.Context, token models.Token) error

	// Verified tokens
	GetVerifiedTokenByAssetIDChain(ctx context.Context, assetID models.VerifiedTokenAssetID, chain string) (*models.VerifiedToken, error)
	GetVerifiedTokenByChainContractAddress(ctx context.Context, chain string, contractAddress string) (*models.VerifiedToken, error)
	GetVerifiedTokensByAssetID(ctx context.Context, assetID models.VerifiedTokenAssetID) ([]models.VerifiedToken, error)
	GetVerifiedTokensByContractAddress(ctx context.Context, contractAddress string) ([]models.VerifiedToken, error)
	GetVerifiedTokensBySymbol(ctx context.Context, symbol string) ([]models.VerifiedToken, error)
	UpsertVerifiedToken(ctx context.Context, token models.VerifiedToken) error
	DeleteVerifiedToken(ctx context.Context, assetID models.VerifiedTokenAssetID, chain string) error
	DeleteVerifiedTokensByAssetID(ctx context.Context, assetID models.VerifiedTokenAssetID) error

	// Dropstab info
	GetDropstabInfoByAssetID(ctx context.Context, assetID models.VerifiedTokenAssetID) (*models.DropstabInfo, error)
	GetDropstabInfosByAssetIDs(ctx context.Context, assetIDs []models.VerifiedTokenAssetID) (map[models.VerifiedTokenAssetID]*models.DropstabInfo, error)
	GetDropstabInfoBySymbol(ctx context.Context, symbol string) (map[models.VerifiedTokenAssetID]*models.DropstabInfo, error)
	GetDropstabInfoBySlug(ctx context.Context, slug string) (*models.DropstabInfo, error)
	GetExistingDropstabInfoSlugs(ctx context.Context, slugs []string) ([]models.DropstabShortCoinDescription, error)
	GetAssetIDsBySlugs(ctx context.Context, slugs []string) (map[string]models.VerifiedTokenAssetID, error)

	UpsertDropstabInfo(ctx context.Context, info models.DropstabInfo) (models.VerifiedTokenAssetID, error)
	DeleteDropstabInfo(ctx context.Context, assetID models.VerifiedTokenAssetID) error

	// Article sources
	UpsertSourceFeed(ctx context.Context, source string, feed string) error
	SelectArticles(ctx context.Context, keywords []string, source string, timeFrom time.Time, timeTo time.Time) ([]models.Article, error)
	GetArticlesByIDs(ctx context.Context, ids []int) ([]models.Article, error)

	// Socials
	CreateSocialPostsTable(ctx context.Context) error
	UpsertSocialPost(ctx context.Context, post models.SocialPost) (bool, error)
	GetSocialPostByLink(ctx context.Context, link string) (*models.SocialPost, error)
	GetSocialPostsByLink(ctx context.Context, link string) ([]models.SocialPost, error)
	GetSocialPostsByLinkBeginning(ctx context.Context, beginning string) ([]models.SocialPost, error)
	GetPostsReferencingLink(ctx context.Context, link string, limit int) ([]models.SocialPost, error)
	GetOldestPostForAuthor(ctx context.Context, source models.SocialSourceTag, author string) (*models.SocialPost, error)
	GetNewestPostForAuthor(ctx context.Context, source models.SocialSourceTag, author string) (*models.SocialPost, error)
	GetSocialPostsByAuthor(ctx context.Context, author string) ([]models.SocialPost, error)

	CreateSocialSourcesTable(ctx context.Context) error
	UpsertSocialSource(ctx context.Context, source models.SocialSource) error
	GetSocialSources(ctx context.Context, tag models.SocialSourceTag) ([]string, error)

	// Backend
	CreateBackendDatabase(ctx context.Context) error
	CreateBackendUsersTable(ctx context.Context) error
	CreateBackendSessionsTable(ctx context.Context) error
	CreateBackendChatHistoryTable(ctx context.Context) error
	CreateBackendChatFoldersTable(ctx context.Context) error
	CreateBackendChatFeedbackTable(ctx context.Context) error
	CreateBackendVizPayloadsTable(ctx context.Context) error
	CreateBackendVizDataRefsTable(ctx context.Context) error
	CreateAPIEndpointBlacklistTable(ctx context.Context) error
	CreateBackendSignalsTable(ctx context.Context) error
	CreateBackendSignalEventsTable(ctx context.Context) error

	GetBackendUserIDByWallet(ctx context.Context, wallet string) (string, error)
	InsertBackendUserWithWallet(ctx context.Context, userID, wallet string) error
	GetBackendUserByID(ctx context.Context, id string) (*models.BackendUser, error)
	GetBackendUserByReferralCode(ctx context.Context, referralCode string) (*models.BackendUser, error)
	UpdateBackendUser(ctx context.Context, user models.BackendUser) error
	SetBackendUserReferralCode(ctx context.Context, userID, referralCode string) (bool, error)
	SetBackendUserReferrerCode(ctx context.Context, userID, referrerCode string) (bool, error)
	SetBackendUserReferralWallet(ctx context.Context, userID, referralWallet string) (bool, error)
	InsertBackendReferralRewardForPaidOrder(ctx context.Context, order models.BillingOrder) (bool, error)
	GetBackendUserReferralInfo(ctx context.Context, userID string) (models.BackendReferralBalance, error)
	GetBackendUserReferrals(ctx context.Context, userID string, limit, offset int) ([]models.BackendReferral, int64, error)
	CreateBackendReferralRewardsTable(ctx context.Context) error
	DeleteBackendUserByID(ctx context.Context, id string) error
	GetProfileExtra(ctx context.Context, userID string) (models.ProfileExtra, error)
	GetProfileExtraForUpdate(ctx context.Context, userID string) (models.ProfileExtra, error)
	UpdateProfileExtra(ctx context.Context, userID string, profile models.ProfileExtra) error
	ListActiveAPIEndpointBlacklistRules(ctx context.Context) ([]models.APIEndpointBlacklistRule, error)

	// Signals
	InsertSignal(ctx context.Context, s models.Signal) error
	// InsertSignalWithCap atomically checks the per-kind active cap and inserts
	// in one transaction; returns false when the cap is reached.
	InsertSignalWithCap(ctx context.Context, s models.Signal, maxActive int) (bool, error)
	GetSignalsByUserID(ctx context.Context, userID string) ([]models.Signal, error)
	GetSignalByID(ctx context.Context, id, userID string) (*models.Signal, error)
	UpdateSignal(ctx context.Context, id, userID string, status *string, cadenceMinutes *int) (bool, error)
	DeleteSignal(ctx context.Context, id, userID string) (bool, error)
	CountActiveSignals(ctx context.Context, userID, kind string) (int64, error)
	SelectDueSignalsForUpdate(ctx context.Context, limit int) ([]models.Signal, error)
	CountDueSignals(ctx context.Context) (int64, error)
	// LeaseSignal pushes next_eval_at forward as a claim lease (does not touch
	// error_count / last_evaluated_at).
	LeaseSignal(ctx context.Context, id string, until time.Time) error
	UpdateSignalAfterEvaluation(ctx context.Context, id string, fired bool, lastValue json.RawMessage, nextEvalAt time.Time, evalErr *string) error
	MarkSignalError(ctx context.Context, id string, lastError string) error
	InsertSignalEvent(ctx context.Context, e models.SignalEvent) error
	GetSignalEventsByUserID(ctx context.Context, userID string, since *time.Time, limit int) ([]models.SignalEvent, error)
	MarkSignalEventSeen(ctx context.Context, id, userID string) (bool, error)
	PruneSignalEvents(ctx context.Context, userID string, keep int) error
	ResetSignalErrorState(ctx context.Context, id, userID string) error

	InsertUserSession(ctx context.Context, session models.BackendUserSession) error
	GetUserSessionByUserIDAndSessionID(ctx context.Context, userID, sessionID string) (*models.BackendUserSession, error)
	GetUserSessionsByUserID(ctx context.Context, userID string) ([]models.BackendUserSession, error)
	GetUserSessionBySessionID(ctx context.Context, sessionID string) (*models.BackendUserSession, error)
	UpdateUserSession(ctx context.Context, session models.BackendUserSession) error
	DeleteUserSession(ctx context.Context, userID, sessionID string) error

	//Billing
	CreateBillingUsersTable(ctx context.Context) error
	CreateBillingTransactionsHistoryTable(ctx context.Context) error
	CreateBillingOrdersTable(ctx context.Context) error
	CreateBillingPromoCodesTable(ctx context.Context) error

	GetBillingUserState(ctx context.Context, userID string) (models.BillingUserState, error)
	GetBillingUserStateForUpdate(ctx context.Context, userID string) (models.BillingUserState, error)
	GetBillingUserStateForUpdateNowait(ctx context.Context, userID string) (models.BillingUserState, error)
	GetBillingUserPlan(ctx context.Context, userID string) (models.PlanCode, error)
	GetBillingUserPlanCredits(ctx context.Context, userID string) (int64, error)
	GetBillingUserAdditionalCredits(ctx context.Context, userID string) (int64, error)
	UpdateBillingUserState(ctx context.Context, state models.BillingUserState) error

	GetBillingApiKeysByUserID(ctx context.Context, userID string) ([]string, error)
	GetBillingUserIDByAPIKey(ctx context.Context, APIKey string) (string, error)

	InsertBillingAPIKeyForUser(ctx context.Context, userID, APIKey string) error
	DeleteBillingAPIKeyFromUser(ctx context.Context, userID, APIKey string) error
	InsertBillingUser(ctx context.Context, userID string) error

	GetBillingPlans(ctx context.Context) ([]models.PlanConfig, error)
	GetBillingUsersForMonthlyRefill(ctx context.Context) ([]models.BillingUserState, error)
	RefillBillingUsersDailyFreeBatch(ctx context.Context, planCredits int64, nextRefillAt, now time.Time, limit int) ([]string, error)
	GetBillingUsersForExpiration(ctx context.Context) ([]models.BillingUserState, error)

	UpdateBillingUserPlanCreditsAndAdditionalCreditsBatch(ctx context.Context, users []models.UserCredits) error

	GetBillingTransactionHistory(ctx context.Context, userID string, eventKind models.BillingEventKind) ([]models.BillingTransactionHistory, error)
	GetBillingTransactionHistoryByOrderID(ctx context.Context, orderID string) (models.BillingTransactionHistory, error)
	InsertBillingTransactionHistory(ctx context.Context, transaction models.BillingTransactionHistory) error
	InsertBillingOrder(ctx context.Context, order models.BillingOrder) error
	GetBillingOrdersByUserID(ctx context.Context, userID string) ([]models.BillingOrder, error)
	GetBillingOrderByID(ctx context.Context, orderID string) (models.BillingOrder, error)
	GetBillingOrderByPaidTxHash(ctx context.Context, txHash string, provider models.BillingPaymentProvider) (models.BillingOrder, error)
	GetBillingOrderByPaymentPayloadHash(ctx context.Context, paymentPayloadHash string, provider models.BillingPaymentProvider) (models.BillingOrder, error)
	GetBillingOrderByIDForUpdate(ctx context.Context, orderID string) (models.BillingOrder, error)
	GetX402PaymentSettledBillingOrders(ctx context.Context) ([]models.BillingOrder, error)
	CountPendingBillingOrdersByUserAndTypes(ctx context.Context, userID string, types []models.BillingOrderType) (int64, error)
	GetExpiredPendingBillingOrders(ctx context.Context) ([]models.BillingOrder, error)
	GetPendingBillingOrderChains(ctx context.Context) ([]evmmodels.Chain, error)
	GetPendingBillingOrdersBatchWithSameScanFromBlock(ctx context.Context, chainName evmmodels.Chain, safeBlock uint64, limit int64) ([]models.BillingOrder, error)
	GetFailedBillingOrders(ctx context.Context) ([]models.BillingOrder, error)
	UpdatePendingBillingOrdersScanFromBlock(ctx context.Context, orderIDs []string, scanFromBlock uint64, updatedAt time.Time) error
	UpdateBillingOrderStatus(ctx context.Context, orderID string, status models.BillingOrderStatus, paidTxHash *string, paidBlockNumber *uint64, updatedAt time.Time) error
	GetBillingPromoCodeByCode(ctx context.Context, code string) (models.PromoCode, error)
	GetBillingPromoCodeByCodeForUpdate(ctx context.Context, code string) (models.PromoCode, error)
	GetBillingPromoCodes(ctx context.Context) ([]models.PromoCode, error)
	InsertBillingPromoCode(ctx context.Context, promo models.PromoCode) error
	UpdateBillingPromoCode(ctx context.Context, promo models.PromoCode) error

	// Payments
	CreatePaymentWalletsTable(ctx context.Context) error
	CreatePaymentChainInfoTable(ctx context.Context) error

	GetAvailablePaymentWalletForUpdate(ctx context.Context) (*models.Wallet, error)
	GetPaymentWallets(ctx context.Context) ([]models.Wallet, error)

	InsertPaymentWallet(ctx context.Context, wallet models.Wallet) error
	ReservePaymentWallet(ctx context.Context, address string, orderID string, reservedAt time.Time) error
	ReleasePaymentWalletReservation(ctx context.Context, address string) error

	GetPaymentChainInfos(ctx context.Context) ([]models.ChainInfo, error)

	// Chat
	InsertChat(ctx context.Context, chat models.BackendChat) error
	GetChatByChatID(ctx context.Context, chatID string) (*models.BackendChat, error)
	GetChatUserByChatID(ctx context.Context, chatID string) (string, error)
	GetMostRecentEmptyChat(ctx context.Context, userID string) (string, error)
	InsertChatFolder(ctx context.Context, folder models.BackendChatFolder) error
	GetChatFoldersByUserID(ctx context.Context, userID string) ([]models.BackendChatFolder, error)
	GetChatFolderByID(ctx context.Context, folderID string) (*models.BackendChatFolder, error)
	UpdateChatFolderName(ctx context.Context, folderID, userID, name string) (bool, error)
	DeleteChatFolder(ctx context.Context, folderID, userID string) (bool, error)
	ClearChatFolderAssignments(ctx context.Context, folderID string) error
	UpsertChatFeedback(ctx context.Context, feedback models.BackendChatFeedback) error
	DeleteChatFeedback(ctx context.Context, chatID, messageID string) error
	GetChatFeedbackByChatID(ctx context.Context, chatID string) (map[string]models.BackendChatFeedbackValue, error)
	InsertVizPayload(ctx context.Context, payload models.BackendVizPayload) error
	UpdateVizPayloadStatus(ctx context.Context, payloadID, userID, status, errorMessage string, updatedAt time.Time) error
	UpdateVizPayloadResult(ctx context.Context, payload models.BackendVizPayload) error
	GetVizPayloadByID(ctx context.Context, payloadID string) (*models.BackendVizPayload, error)
	ListVizPayloadsByUser(ctx context.Context, userID string, before *time.Time, limit int) ([]models.VizPayloadListItem, error)
	InsertVizDataRef(ctx context.Context, ref models.BackendVizDataRef) error
	GetVizDataRefByID(ctx context.Context, refID string) (*models.BackendVizDataRef, error)
	UpdateChatHistory(ctx context.Context, history models.BackendChat) (bool, error)
	UpdateChatShortName(ctx context.Context, chatID, shortName string) (bool, error)
	UpdateInactiveChatShortName(ctx context.Context, chatID, shortName string) (bool, error)
	UpdateChatFolderID(ctx context.Context, chatID string, folderID *string) (bool, error)
	UpdateChatShared(ctx context.Context, chatID string, shared bool) (bool, error)
	GetChatsByUserID(ctx context.Context, userID string, limit, offset int) ([]models.BackendUserChatSummary, error)
	GetChatCountByUserID(ctx context.Context, userID string) (int64, error)
	DeleteChat(ctx context.Context, chatID string) (bool, error)

	LockChat(ctx context.Context, chatID string) (bool, error)
	UnlockChat(ctx context.Context, chatID string) (bool, error)
	GetChatActive(ctx context.Context, chatID string) (bool, error)

	GetChatShareInfo(ctx context.Context, chatID string) (*models.ChatShareInfo, error)

	// CEX methods
	CreateCexSchema(ctx context.Context) error
	UpsertBybitTickerInfo(ctx context.Context, ticker models.TickerInfo) (uint64, error)
	SelectBybitTickerInfoBySymbolAndCategory(ctx context.Context, symbol string, category string) (*models.TickerInfo, error)
	SelectBybitTickerInfoByID(ctx context.Context, id uint64) (*models.TickerInfo, error)
	SelectAllBybitTickerInfo(ctx context.Context) ([]models.TickerInfo, error)

	// EVM parser registry
	CreateEvmParserRegistryTable(ctx context.Context) error
	GetEvmParser(ctx context.Context, chain string, contractAddress string, instructionHash string) (string, string, error)
	GetAllEvmParsersForChain(ctx context.Context, chain string) (map[string]map[string][2]string, error) // [contractAddress][instructionHash][parserName, source]
	UpsertEvmParser(ctx context.Context, chain string, contractAddress string, instructionHash string, parserName string, source string) error

	// Polymarket
	CreatePolymarketTokensTable(ctx context.Context) error
	CreatePolymarketMarketEventsTable(ctx context.Context) error
	InsertPolymarketMarketNew(ctx context.Context, market models.PolymarketMarketNew) error
	InsertPolymarketMarketEvent(ctx context.Context, event models.PolymarketMarketEvent) error
	GetPolymarketQuestionTimestampFromEvents(ctx context.Context, questionID string, maxBlockNumber uint64) (uint64, error)
	UpdatePolymarketMarketAncillaryDataByQuestionID(ctx context.Context, questionID string, ancillaryData []byte, parsedAncillaryData models.ParsedAncillaryData) error
	UpdatePolymarketMarketParsedAncillaryDataByQuestionID(ctx context.Context, questionID string, parsedAncillaryData models.ParsedAncillaryData) error
	UpdatePolymarketMarketResolution(ctx context.Context, conditionID string, resolvedAt time.Time, payoutNumerators []*big.Int, resolvedInBlock int64, resolvedInBlockHash string) error
	UpdatePolymarketMarketResolutionFull(ctx context.Context, conditionID string, resolvedAt time.Time, payoutNumerators []*big.Int, resolvedInBlock int64, resolvedInBlockHash string) ([]models.PolymarketToken, error)
	UpdatePolymarketTokensResolution(ctx context.Context, conditionID string, resolvedAt time.Time, payoutNumerators []*big.Int, resolvedInBlock int64, resolvedInBlockHash string) error
	UpdatePolymarketMarketGammaData(ctx context.Context, conditionID string, gammaInfo *models.PolymarketMarketEssentialGammaInfo) (bool, error)
	UpdatePolymarketMarketGammaDataBatch(ctx context.Context, conditionIDs []string, gammaInfos []*models.PolymarketMarketEssentialGammaInfo) (int64, error)
	InsertPolymarketToken(ctx context.Context, token models.PolymarketToken, blockTime time.Time) error
	InsertPolymarketTokensBatch(ctx context.Context, tokens []models.PolymarketToken) error
	GetPolymarketTokensByConditionID(ctx context.Context, conditionID string) ([]models.PolymarketToken, error)
	GetPolymarketTokenByTokenID(ctx context.Context, tokenID *big.Int) (*models.PolymarketToken, error)
	GetPolymarketMarketMinimalByConditionID(ctx context.Context, conditionID string) (*models.PolymarketMarketNewMinimal, error)
	GetPolymarketMarketByConditionID(ctx context.Context, conditionID string) (*models.PolymarketMarketNew, error)
	GetPolymarketMarketEventsByQuestionID(ctx context.Context, questionID string) ([]models.PolymarketMarketEvent, error)

	// Polymarket reorg handling by block_hash
	DeletePolymarketMarketsByBlockHash(ctx context.Context, blockHash string) (int64, error)
	DeletePolymarketMarketEventsByBlockHash(ctx context.Context, blockHash string) (int64, error)
	ResetPolymarketTokensResolutionByBlockHash(ctx context.Context, blockHash string) (int64, error)
	ResetPolymarketMarketsResolutionByBlockHash(ctx context.Context, blockHash string) (int64, error)

	// Polymarket reorg handling by block_number (>= blockNumber)
	DeletePolymarketMarketsFromBlockNumber(ctx context.Context, blockNumber int64) (int64, error)
	DeletePolymarketMarketEventsFromBlockNumber(ctx context.Context, blockNumber int64) (int64, error)
	ResetPolymarketTokensResolutionFromBlockNumber(ctx context.Context, blockNumber int64) (int64, error)
	ResetPolymarketMarketsResolutionFromBlockNumber(ctx context.Context, blockNumber int64) (int64, error)

	// Investors
	CreateInvestorsTable(ctx context.Context) error
	CreateInvestRoundsTable(ctx context.Context) error
	UpsertInvestor(ctx context.Context, inv models.Investor) error
	UpsertInvestRound(ctx context.Context, round models.InvestRound) error
	GetExistingInvestorSlugs(ctx context.Context, slugs []string) ([]string, error)
	GetExistingInvestRoundIDs(ctx context.Context, ids []models.InvestRoundID) ([]models.InvestRoundID, error)
	GetInvestorIdsBySlugs(ctx context.Context, slugs []string) (map[string]models.InvestorID, error)

	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
}

type DBClick interface {
	// Schema and tables
	CreateSchemaCex(ctx context.Context) error
	CreateSchemaEvm(ctx context.Context) error
	CreateBybitKlineTable(ctx context.Context) error
	CreateBybitFundingRateTable(ctx context.Context) error
	CreateBybitTickerInfoTable(ctx context.Context) error
	CreateBybitLiquidationTable(ctx context.Context) error
	CreateBybitTradeTable(ctx context.Context) error
	CreateEvmSwapEventsTable(ctx context.Context) error
	CreateEvmTransferEventsTable(ctx context.Context) error
	CreateEvmDefiEventsTable(ctx context.Context) error
	CreateEvmPolymarketOrderEventsTable(ctx context.Context) error
	CreateBillingLogsTable(ctx context.Context) error

	Query(ctx context.Context, query string, args ...interface{}) (driver.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) error
	InsertBillingLogsBatch(ctx context.Context, logs []models.BillingLog) error

	// Kline
	InsertBybitKline(ctx context.Context, kline *models.KlineInfo) error
	InsertBybitKlineBatch(ctx context.Context, klines []*models.KlineInfo) error
	GetBybitKlineLatestOpenTime(ctx context.Context, tickerID uint64) (time.Time, error)
	GetBybitKlineData(ctx context.Context, tickerID uint64, startTime, endTime time.Time) ([]*models.KlineInfo, error)
	GetBybitKlineDataAggregated(ctx context.Context, tickerID uint64, intervalMinutes int, startTime, endTime time.Time) ([]*models.KlineInfo, error)
	GetBybitKlineLatestClosePrice(ctx context.Context, tickerID uint64) (float64, error)

	// Funding rate
	InsertBybitFundingRate(ctx context.Context, fundingRate *models.FundingRateInfo) error
	InsertBybitFundingRateBatch(ctx context.Context, fundingRates []*models.FundingRateInfo) error
	GetBybitFundingRateLatestFundingRateTimestamp(ctx context.Context, tickerID uint64) (uint64, error)
	GetBybitFundingRateData(ctx context.Context, tickerID uint64, startTimestamp, endTimestamp uint64) ([]*models.FundingRateInfo, error)

	// Ticker info (ClickHouse replica)
	InsertBybitTickerInfo(ctx context.Context, ticker *models.TickerInfo) error
	InsertBybitTickerInfoBatch(ctx context.Context, tickers []models.TickerInfo) error
	GetAllBybitTickerInfo(ctx context.Context) ([]models.TickerInfo, error)

	// Liquidations
	InsertBybitLiquidation(ctx context.Context, liq *models.LiquidationInfo) error
	InsertBybitLiquidationBatch(ctx context.Context, liqs []*models.LiquidationInfo) error
	GetBybitLiquidationLatestTimestamp(ctx context.Context, tickerID uint64) (uint64, error)
	GetBybitLiquidationData(ctx context.Context, tickerID uint64, startTimestamp, endTimestamp uint64) ([]*models.LiquidationInfo, error)

	// Trades
	InsertBybitTrade(ctx context.Context, trade *models.TradeInfo) error
	InsertBybitTradeBatch(ctx context.Context, trades []*models.TradeInfo) error
	GetBybitTradeLatestTimestamp(ctx context.Context, tickerID uint64) (uint64, error)
	GetBybitTradeData(ctx context.Context, tickerID uint64, startTimestamp, endTimestamp uint64) ([]*models.TradeInfo, error)

	// EVM events
	InsertEvmEventsBatch(ctx context.Context, chain evmmodels.Chain, events []evmmodels.Event) error
	DeleteEvmDefiEventsBatch(ctx context.Context, chain evmmodels.Chain, blockNumbers []*big.Int, txHashes []common.Hash) error
	DeleteEvmDefiEventsByContractAddress(ctx context.Context, chain evmmodels.Chain, contractAddress common.Address, instructionHash string) error

	// EVM select queries
	GetEvmSwapEvents(ctx context.Context, chain evmmodels.Chain, startBlock, endBlock *big.Int) ([]evmmodels.Event, error)
	GetEvmSwapEventsChan(ctx context.Context, chain evmmodels.Chain, startBlock, endBlock *big.Int, buffer int) (<-chan evmmodels.Event, <-chan error, error)
	GetEvmTransferEvents(ctx context.Context, chain evmmodels.Chain, startBlock, endBlock *big.Int) ([]evmmodels.Event, error)
	GetEvmTransferEventsChan(ctx context.Context, chain evmmodels.Chain, startBlock, endBlock *big.Int, buffer int) (<-chan evmmodels.Event, <-chan error, error)
	GetEvmDefiEvents(ctx context.Context, chain evmmodels.Chain, startBlock, endBlock *big.Int) ([]evmmodels.Event, error)
	GetEvmDefiEventsChan(ctx context.Context, chain evmmodels.Chain, contractAddress common.Address, instructionHash string, startBlock, endBlock *big.Int, buffer int) (<-chan evmmodels.Event, <-chan error, error)
	GetEvmErrorEvents(ctx context.Context, chain evmmodels.Chain, startBlock, endBlock *big.Int) ([]evmmodels.Event, error)
	GetEvmErrorEventsChan(ctx context.Context, chain evmmodels.Chain, startBlock, endBlock *big.Int, buffer int) (<-chan evmmodels.Event, <-chan error, error)

	// Listings
	InsertListing(ctx context.Context, cexName string, ticker string, dt time.Time) error

	// Polymarket
	InsertEvmPolymarketOrderEventsBatchNew(ctx context.Context, events []*evmmodels.PolymarketOrderEventNew) error

	// Delete by block_hash for reorg handling
	DeleteEvmSwapEventsByBlockHash(ctx context.Context, chain evmmodels.Chain, blockHash string) error
	DeleteEvmTransferEventsByBlockHash(ctx context.Context, chain evmmodels.Chain, blockHash string) error
	DeleteEvmDefiEventsByBlockHash(ctx context.Context, chain evmmodels.Chain, blockHash string) error
	DeleteEvmPolymarketOrderEventsByBlockHash(ctx context.Context, blockHash string) error
	DeleteEvmEventsByBlockHash(ctx context.Context, chain evmmodels.Chain, blockHash string) error

	// Delete by block_number for reorg handling (>= blockNumber)
	DeleteEvmSwapEventsFromBlockNumber(ctx context.Context, chain evmmodels.Chain, blockNumber uint64) error
	DeleteEvmTransferEventsFromBlockNumber(ctx context.Context, chain evmmodels.Chain, blockNumber uint64) error
	DeleteEvmDefiEventsFromBlockNumber(ctx context.Context, chain evmmodels.Chain, blockNumber uint64) error
	DeleteEvmPolymarketOrderEventsFromBlockNumber(ctx context.Context, blockNumber uint64) error
	DeleteEvmEventsFromBlockNumber(ctx context.Context, chain evmmodels.Chain, blockNumber uint64) error
}

const DropstabDescriptionCollection = "dropstab_description"
const ArticlesCollection = "articles"
const SocialPostsCollection = "social_posts"

const DropstabDescriptionVectorSize = 384
const SocialPostsVectorSize = 384

type VecDB interface {
	Search(ctx context.Context, query string, collection string, limit int, threshold float32, filter *qdrant.Filter) (map[int]float32, error)
	SearchStringID(ctx context.Context, query string, collection string, limit int, threshold float32, filter *qdrant.Filter) (map[string]float32, error)
	SearchAggregated(ctx context.Context, query string, aggregationField string, collection string, limit int, threshold float32, filter *qdrant.Filter) (map[any]float32, error)

	Store(ctx context.Context, text string, collection string, id any, payload map[string]any) error

	CreateCollection(ctx context.Context, collection string, vectorSize int, indexes *qdrant.CreateFieldIndexCollection) error

	StoreDescription(ctx context.Context, description string, assetID models.VerifiedTokenAssetID) error
	SearchDescriptions(ctx context.Context, query string, limit int, threshold float32, filter *qdrant.Filter) (map[models.VerifiedTokenAssetID]float32, error)
}

const (
	UniqueViolation = "23505"
)

type KeyValueDB interface {
}
