package models

import (
	"testing"
	"time"
)

func TestProfileExtraIsEmpty(t *testing.T) {
	if !EmptyProfileExtra().IsEmpty() {
		t.Fatalf("EmptyProfileExtra must be empty")
	}
	p := EmptyProfileExtra()
	p.Wallets = []ProfileWallet{{Address: "0x", Chain: "evm", Label: "x", AddedAt: time.Now()}}
	if p.IsEmpty() {
		t.Fatalf("non-empty wallets must not be empty")
	}
}

func TestProfileExtraIsEmpty_PreferencesNonEmpty(t *testing.T) {
	p := EmptyProfileExtra()
	chain := "evm"
	p.Preferences.DefaultChain = &chain
	if p.IsEmpty() {
		t.Fatalf("non-empty preferences must not be empty")
	}
}

func TestProfileExtraIsEmpty_NotesNonEmpty(t *testing.T) {
	p := EmptyProfileExtra()
	p.Notes = []ProfileNote{{ID: "n_x", Text: "x", Kind: "trait", Source: "user", CapturedAt: time.Now()}}
	if p.IsEmpty() {
		t.Fatalf("non-empty notes must not be empty")
	}
}
