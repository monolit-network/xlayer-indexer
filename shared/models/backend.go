package models

import (
	"encoding/json"
	"time"

	thirdweb "github.com/monolit-network/xlayer-indexer/shared/apis/thirdweb"
)

type BackendUser struct {
	ID             string                     `json:"id"`
	ThirdwebUserID string                     `json:"thirdweb_user_id"`
	Source         string                     `json:"source"`
	Profiles       []thirdweb.UserInfoProfile `json:"profiles"`
	ReferralCode   *string                    `json:"referral_code,omitempty"`
	ReferrerCode   *string                    `json:"referrer_code,omitempty"`
	ReferralWallet *string                    `json:"referral_wallet,omitempty"`
	ProfileExtra   ProfileExtra               `json:"profile_extra"`
}

type BackendReferral struct {
	WalletPrefix             string `json:"wallet_prefix"`
	ReferralEarningsUSDCents int64  `json:"referral_earnings_usd_cents"`
}

type BackendReferralRewardStatus string

const (
	BackendReferralRewardStatusAvailable BackendReferralRewardStatus = "available"
	BackendReferralRewardStatusPaidOut   BackendReferralRewardStatus = "paid_out"
)

type BackendReferralBalance struct {
	ReferralCode      *string `json:"referral_code,omitempty"`
	ReferralCount     int64   `json:"referral_count"`
	TotalUSDCents     int64   `json:"total_usd_cents"`
	PaidOutUSDCents   int64   `json:"paid_out_usd_cents"`
	AvailableUSDCents int64   `json:"available_usd_cents"`
	ReferralWallet    string  `json:"referral_wallet"`
}

type BackendUserSession struct {
	UserID    string    `json:"user_id"`
	SessionID string    `json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type APIEndpointBlacklistRule struct {
	ID            int64      `json:"id"`
	Method        string     `json:"method"`
	MatchType     string     `json:"match_type"`
	Pattern       string     `json:"pattern"`
	Reason        string     `json:"reason"`
	DisabledUntil *time.Time `json:"disabled_until,omitempty"`
}

type BackendChat struct {
	ChatID    string          `json:"chat_id"`
	UserID    string          `json:"user_id"`
	ShortName string          `json:"short_name"`
	FolderID  *string         `json:"folder_id,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	History   json.RawMessage `json:"history"`
	Active    bool            `json:"active"`
}

type BackendUserChatSummary struct {
	ChatID    string    `json:"chat_id"`
	ShortName string    `json:"short_name"`
	FolderID  *string   `json:"folder_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BackendChatFolder struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BackendChatFeedback struct {
	ChatID    string                      `json:"chat_id"`
	MessageID string                      `json:"message_id"`
	Mark      *int                        `json:"mark"`
	Comment   *BackendChatFeedbackComment `json:"comment"`
}

type BackendChatFeedbackValue struct {
	Mark    *int                        `json:"mark"`
	Comment *BackendChatFeedbackComment `json:"comment"`
}

type BackendChatFeedbackComment struct {
	Tags []string `json:"tags"`
	Text string   `json:"text"`
}

const (
	BackendVizPayloadStatusGenerating = "generating"
	BackendVizPayloadStatusCompleted  = "completed"
	BackendVizPayloadStatusFailed     = "failed"
)

type BackendVizPayload struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	ChatID    *string         `json:"chat_id,omitempty"`
	MessageID string          `json:"message_id,omitempty"`
	RefID     string          `json:"ref_id,omitempty"`
	Payload   json.RawMessage `json:"payload"`
	Status    string          `json:"status"`
	Error     string          `json:"error,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type VizPayloadListItem struct {
	ID        string    `json:"id"`
	ChatID    *string   `json:"chat_id,omitempty"`
	MessageID string    `json:"message_id,omitempty"`
	RefID     string    `json:"ref_id,omitempty"`
	Title     string    `json:"title"`
	VizType   string    `json:"viz_type"`
	CreatedAt time.Time `json:"created_at"`
}

type BackendVizDataRef struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	ChatID    *string         `json:"chat_id,omitempty"`
	MessageID string          `json:"message_id,omitempty"`
	RefID     string          `json:"ref_id,omitempty"`
	Path      string          `json:"path,omitempty"`
	Spec      json.RawMessage `json:"spec"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ChatShareInfo struct {
	OwnerID string `json:"owner_id"`
	Shared  bool   `json:"shared"`
	Active  bool   `json:"active"`
}

// ProfileExtra is the shape stored in backend.users.profile_extra (JSONB).
// schema_version=1 is the only version currently defined.
type ProfileExtra struct {
	SchemaVersion int                `json:"schema_version"`
	Wallets       []ProfileWallet    `json:"wallets"`
	Notes         []ProfileNote      `json:"notes"`
	Preferences   ProfilePreferences `json:"preferences"`
}

type ProfileWallet struct {
	Address   string    `json:"address"`
	Chain     string    `json:"chain"` // "evm" | "solana"
	Label     string    `json:"label"` // ≤32 chars
	IsPrimary bool      `json:"is_primary"`
	AddedAt   time.Time `json:"added_at"`
}

type ProfileNote struct {
	ID           string    `json:"id"`     // "n_" + UUIDv7
	Text         string    `json:"text"`   // ≤500 chars
	Kind         string    `json:"kind"`   // "trait" | "preference" | "goal" | "context"
	Source       string    `json:"source"` // "agent" | "user"
	SourceChatID *string   `json:"source_chat_id,omitempty"`
	CapturedAt   time.Time `json:"captured_at"`
}

type ProfilePreferences struct {
	DefaultChain   *string  `json:"default_chain,omitempty"`   // "evm" | "solana"
	RiskProfile    *string  `json:"risk_profile,omitempty"`    // "conservative" | "moderate" | "degen"
	FavoriteTokens []string `json:"favorite_tokens,omitempty"` // uppercase, ≤10
}

// EmptyProfileExtra returns the default zero-value profile used as DB column default
// and as initial value for new users.
func EmptyProfileExtra() ProfileExtra {
	return ProfileExtra{
		SchemaVersion: 1,
		Wallets:       []ProfileWallet{},
		Notes:         []ProfileNote{},
		Preferences:   ProfilePreferences{},
	}
}
