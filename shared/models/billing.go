package models

import (
	"encoding/json"
	"time"

	evmmodels "github.com/monolit-network/xlayer-indexer/evm/pkg/models"
)

type PlanCode string

const (
	PlanCodeFree PlanCode = "free"
)

const (
	PlanDuration     = time.Hour * 24 * 30
	MonthDuration    = time.Hour * 24 * 30
	FreeRefillPeriod = 24 * time.Hour
)

// NextFreePlanRefillAt returns the next shared daily refill boundary. Free plan
// credits reset for every user at 00:00 UTC instead of 24 hours after signup.
func NextFreePlanRefillAt(now time.Time) time.Time {
	now = now.UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Add(FreeRefillPeriod)
}

const (
	AdditionalCreditsMinimumAmount  = 1_000_000
	USDToAdditionalCreditsRate      = 1_000_000
	USDCentsToAdditionalCreditsRate = 10_000
	OrderTTL                        = time.Hour
	RefbackBps                      = 500
	ReferralDiscountBps             = 500
)

type FeatureCode string

const (
	FeatureFreeAccess  FeatureCode = "free-handlers"
	FeatureBasicAccess FeatureCode = "basic-handlers"
	FeatureProAccess   FeatureCode = "pro-handlers"
	FeatureMaxAccess   FeatureCode = "max-handlers"
)

type PlanConfig struct {
	RequestsPerMinute uint64        `json:"requests_per_minute"`
	Credits           int64         `json:"credits"`
	Name              string        `json:"name"`
	Code              PlanCode      `json:"code"`
	Features          []FeatureCode `json:"features"`
	PriceUSD          int64         `json:"price_usd"`
}

type UserCredits struct {
	UserID            string
	PlanCredits       int64
	AdditionalCredits int64
}

type BillingOperationType string

const (
	BillingOperationSpend  BillingOperationType = "spend"
	BillingOperationRefund BillingOperationType = "refund"
)

const (
	BillingSourceAPICall      = "api_call"
	BillingSourceChatQuery    = "chat_query"
	BillingSourceMCPQuery     = "mcp_query"
	BillingSourceWSConnection = "ws_connection"
)

type BillingLog struct {
	UserID         string
	APIKey         string
	CreditsAmount  int64
	OperationType  BillingOperationType
	Source         string
	AdditionalData json.RawMessage
}

type BillingUserState struct {
	UserID                 string
	APIKeys                []string
	Plan                   PlanCode
	PlanCredits            int64
	AdditionalCredits      int64
	BillingPeriodPaidCents int64
	PlanStartedAt          *time.Time
	NextPlanRefillAt       *time.Time
	PlanEndingAt           *time.Time
}

type BillingOrderType string

const (
	BillingOrderTypePlanPurchase      BillingOrderType = "plan_purchase"
	BillingOrderTypePlanUpgrade       BillingOrderType = "plan_upgrade"
	BillingOrderTypeAdditionalCredits BillingOrderType = "additional_credits"
)

type PromoScope string

const (
	PromoScopeBuyPlan    PromoScope = "buy_plan"
	PromoScopeUpgrade    PromoScope = "upgrade"
	PromoScopeBuyCredits PromoScope = "buy_credits"
)

type PromoDiscountType string

const (
	PromoDiscountTypePercent PromoDiscountType = "percent"
	PromoDiscountTypeFixed   PromoDiscountType = "fixed"
)

type BillingOrderStatus string

const (
	BillingOrderStatusPending           BillingOrderStatus = "pending"
	BillingOrderStatusPaid              BillingOrderStatus = "paid"
	BillingOrderStatusPaymentSettled    BillingOrderStatus = "payment_settled"
	BillingOrderStatusPaymentConflict   BillingOrderStatus = "payment_conflict"
	BillingOrderStatusFulfillmentFailed BillingOrderStatus = "fulfillment_failed"
	BillingOrderStatusCanceled          BillingOrderStatus = "canceled"
	BillingOrderStatusExpired           BillingOrderStatus = "expired"
)

type BillingPaymentProvider string

const (
	BillingPaymentProviderWalletTransfer BillingPaymentProvider = "wallet_transfer"
	BillingPaymentProviderX402           BillingPaymentProvider = "x402"
)

type BillingEventKind string

const (
	BillingEventOrderCreated           BillingEventKind = "order_created"
	BillingEventOrderPaid              BillingEventKind = "order_paid"
	BillingEventOrderFulfilled         BillingEventKind = "order_fulfilled"
	BillingEventOrderFulfillmentFailed BillingEventKind = "order_fulfillment_failed"
	BillingEventOrderCanceled          BillingEventKind = "order_canceled"
	BillingEventOrderExpired           BillingEventKind = "order_expired"
	BillingEventLatePaymentIgnored     BillingEventKind = "late_payment_ignored"
	BillingEventPlanMonthlyRefill      BillingEventKind = "plan_monthly_refill"
	BillingEventFreeDailyRefill        BillingEventKind = "free_daily_refill"
	BillingEventPlanExpiredToFree      BillingEventKind = "plan_expired_to_free"
	BillingEventManualPlanGranted      BillingEventKind = "manual_plan_granted"
	BillingEventDemoGranted            BillingEventKind = "demo_granted"
)

type BillingOrder struct {
	ID                 string
	UserID             string
	Status             BillingOrderStatus
	Type               BillingOrderType
	Chain              evmmodels.Chain
	TokenAddress       string
	AmountUSDCents     int64
	WalletAddress      string
	PromoCode          *string
	ProductData        json.RawMessage
	ExpiresAt          time.Time
	ScanFromBlock      uint64
	PaidTxHash         *string
	PaidBlockNumber    *uint64
	PaymentProvider    BillingPaymentProvider
	PaymentPayloadHash *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type BillingTransactionHistory struct {
	OrderID        string
	UserID         string
	EventKind      BillingEventKind
	AmountUSDCents *int64
	CreatedAt      time.Time
	AdditionalInfo json.RawMessage
}

type PlanPurchaseProductData struct {
	Plan   PlanCode `json:"plan"`
	Months int64    `json:"months"`
}

type PlanUpgradeProductData struct {
	FromPlan           PlanCode `json:"from_plan"`
	ToPlan             PlanCode `json:"to_plan"`
	Months             int64    `json:"months"`
	NewPeriodPaidCents int64    `json:"new_period_paid_cents"`
}

type AdditionalCreditsProductData struct {
	Credits int64 `json:"credits"`
}

type PromoCode struct {
	Code          string            `json:"code"`
	Active        bool              `json:"active"`
	Scopes        []PromoScope      `json:"scopes"`
	DiscountType  PromoDiscountType `json:"discount_type"`
	DiscountValue int64             `json:"discount_value"`
	MaxUses       *int64            `json:"max_uses,omitempty"`
	UsesCount     int64             `json:"uses_count"`
	TargetPlans   []PlanCode        `json:"target_plans"`
	MinMonths     *int64            `json:"min_months,omitempty"`
	MaxMonths     *int64            `json:"max_months,omitempty"`
	MinCredits    *int64            `json:"min_credits,omitempty"`
	MaxCredits    *int64            `json:"max_credits,omitempty"`
	ExpiresAt     *time.Time        `json:"expires_at,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type HandlerCost struct {
	BaseCost     int64
	DiscountCost int64
}
