package models

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	SignalKindAlert    = "alert"
	SignalKindInterest = "interest"

	SignalStatusActive = "active"
	SignalStatusPaused = "paused"
	SignalStatusError  = "error"
)

var AllowedCadences = map[int]struct{}{5: {}, 15: {}, 60: {}, 360: {}, 720: {}, 1440: {}}

type Signal struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id"`
	Kind            string          `json:"kind"`
	Status          string          `json:"status"`
	NLText          string          `json:"nl_text"`
	Spec            json.RawMessage `json:"spec"`
	HumanReadable   string          `json:"human_readable"`
	CadenceMinutes  int             `json:"cadence_minutes"`
	CooldownMinutes int             `json:"cooldown_minutes"`
	LastValue       json.RawMessage `json:"last_value,omitempty"`
	LastEvaluatedAt *time.Time      `json:"last_evaluated_at,omitempty"`
	LastFiredAt     *time.Time      `json:"last_fired_at,omitempty"`
	NextEvalAt      time.Time       `json:"next_eval_at"`
	ErrorCount      int             `json:"error_count"`
	LastError       *string         `json:"last_error,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type SignalEvent struct {
	ID            string          `json:"id"`
	SignalID      string          `json:"signal_id"`
	UserID        string          `json:"user_id"`
	FiredAt       time.Time       `json:"fired_at"`
	Title         string          `json:"title"`
	Body          string          `json:"body"`
	ValueSnapshot json.RawMessage `json:"value_snapshot,omitempty"`
	SeedQuery     string          `json:"seed_query,omitempty"`
	SeenAt        *time.Time      `json:"seen_at,omitempty"`
	Channel       string          `json:"channel"`
}

// ValidateSignalCreate checks kind + cadence invariants shared by the create
// handler and the migration-era seeding paths. Returns *ValidationError
// (profile_validation.go) so handlers map it to 400 via the existing
// errors.As dispatch.
func ValidateSignalCreate(kind string, cadenceMinutes int) error {
	if kind != SignalKindAlert && kind != SignalKindInterest {
		return &ValidationError{Field: "kind", Message: fmt.Sprintf("unknown kind %q", kind)}
	}
	if _, ok := AllowedCadences[cadenceMinutes]; !ok {
		return &ValidationError{Field: "cadence_minutes", Message: fmt.Sprintf("cadence %d not in allowed set", cadenceMinutes)}
	}
	if kind == SignalKindInterest && cadenceMinutes < 720 {
		return &ValidationError{Field: "cadence_minutes", Message: "interest signals require cadence ≥720 minutes"}
	}
	return nil
}
