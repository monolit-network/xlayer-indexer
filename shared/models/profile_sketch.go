package models

// IsEmpty reports whether the profile has no addressable content.
// When empty, the full document is not sent to Python (saves WS payload bytes).
func (p ProfileExtra) IsEmpty() bool {
	if len(p.Wallets) > 0 || len(p.Notes) > 0 {
		return false
	}
	prefs := p.Preferences
	if prefs.DefaultChain != nil || prefs.RiskProfile != nil || len(prefs.FavoriteTokens) > 0 {
		return false
	}
	return true
}
