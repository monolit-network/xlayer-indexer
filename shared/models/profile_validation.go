package models

import (
	"fmt"
	"strings"
)

// ValidationError is returned when a profile_extra field fails an invariant check.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return e.Field + ": " + e.Message
}

// CapReachedError is returned when a profile_extra collection hits its size limit.
type CapReachedError struct {
	What string // "wallet" | "note" | "favorite_tokens"
	Max  int
}

func (e *CapReachedError) Error() string {
	return fmt.Sprintf("%s cap reached (max %d)", e.What, e.Max)
}

const (
	maxWallets        = 20
	maxNotes          = 50
	maxLabelLen       = 32
	maxNoteTextLen    = 500
	maxFavoriteTokens = 10
)

var (
	validChains      = map[string]struct{}{"evm": {}, "solana": {}}
	validNoteKinds   = map[string]struct{}{"trait": {}, "preference": {}, "goal": {}, "context": {}}
	validNoteSources = map[string]struct{}{"agent": {}, "user": {}}
	validRisks       = map[string]struct{}{"conservative": {}, "moderate": {}, "degen": {}}
)

func ValidateProfileExtra(p ProfileExtra) error {
	if p.SchemaVersion != 1 {
		return &ValidationError{Message: fmt.Sprintf("unsupported schema_version: %d", p.SchemaVersion)}
	}
	if len(p.Wallets) > maxWallets {
		// Error string contains "wallet cap" to satisfy existing tests.
		return &CapReachedError{What: "wallet", Max: maxWallets}
	}
	primaryCount := 0
	for i, w := range p.Wallets {
		if w.IsPrimary {
			primaryCount++
		}
		if _, ok := validChains[w.Chain]; !ok {
			return &ValidationError{Field: fmt.Sprintf("wallets[%d]", i), Message: fmt.Sprintf("unknown chain %q", w.Chain)}
		}
		if len(w.Label) == 0 {
			return &ValidationError{Field: fmt.Sprintf("wallets[%d]", i), Message: "label required"}
		}
		if len(w.Label) > maxLabelLen {
			return &ValidationError{Field: fmt.Sprintf("wallets[%d]", i), Message: fmt.Sprintf("label too long (%d > %d)", len(w.Label), maxLabelLen)}
		}
		if strings.TrimSpace(w.Address) == "" {
			return &ValidationError{Field: fmt.Sprintf("wallets[%d]", i), Message: "address required"}
		}
	}
	if primaryCount > 1 {
		// Error string contains "primary" to satisfy existing tests.
		return &ValidationError{Message: "more than one primary wallet"}
	}
	if len(p.Notes) > maxNotes {
		// Error string contains "note cap" to satisfy existing tests.
		return &CapReachedError{What: "note", Max: maxNotes}
	}
	for i, n := range p.Notes {
		if _, ok := validNoteKinds[n.Kind]; !ok {
			// Error string contains "kind" to satisfy existing tests.
			return &ValidationError{Field: fmt.Sprintf("notes[%d]", i), Message: fmt.Sprintf("unknown kind %q", n.Kind)}
		}
		if _, ok := validNoteSources[n.Source]; !ok {
			return &ValidationError{Field: fmt.Sprintf("notes[%d]", i), Message: fmt.Sprintf("unknown source %q", n.Source)}
		}
		if len(n.Text) == 0 {
			return &ValidationError{Field: fmt.Sprintf("notes[%d]", i), Message: "text required"}
		}
		if len(n.Text) > maxNoteTextLen {
			// Error string contains "note text" to satisfy existing tests.
			return &ValidationError{Field: fmt.Sprintf("notes[%d]", i), Message: fmt.Sprintf("note text too long (%d > %d)", len(n.Text), maxNoteTextLen)}
		}
	}
	if p.Preferences.RiskProfile != nil {
		if _, ok := validRisks[*p.Preferences.RiskProfile]; !ok {
			return &ValidationError{Field: "preferences.risk_profile", Message: fmt.Sprintf("unknown value %q", *p.Preferences.RiskProfile)}
		}
	}
	if p.Preferences.DefaultChain != nil {
		if _, ok := validChains[*p.Preferences.DefaultChain]; !ok {
			return &ValidationError{Field: "preferences.default_chain", Message: fmt.Sprintf("unknown value %q", *p.Preferences.DefaultChain)}
		}
	}
	if len(p.Preferences.FavoriteTokens) > maxFavoriteTokens {
		// Error string contains "favorite" to satisfy existing tests.
		return &CapReachedError{What: "favorite_tokens", Max: maxFavoriteTokens}
	}
	return nil
}
