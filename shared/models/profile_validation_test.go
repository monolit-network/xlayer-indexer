package models

import (
	"strings"
	"testing"
	"time"
)

func validWallet(addr string) ProfileWallet {
	return ProfileWallet{
		Address: addr, Chain: "evm", Label: "x", IsPrimary: false,
		AddedAt: time.Now().UTC(),
	}
}

func TestValidateProfileExtra_HappyPath(t *testing.T) {
	p := EmptyProfileExtra()
	p.Wallets = []ProfileWallet{validWallet("0xabc")}
	p.Wallets[0].IsPrimary = true
	if err := ValidateProfileExtra(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateProfileExtra_WalletsCap(t *testing.T) {
	p := EmptyProfileExtra()
	for i := 0; i < 21; i++ {
		p.Wallets = append(p.Wallets, validWallet("0x"))
	}
	if err := ValidateProfileExtra(p); err == nil || !strings.Contains(err.Error(), "wallet cap") {
		t.Fatalf("expected wallet cap error, got %v", err)
	}
}

func TestValidateProfileExtra_MultiplePrimaryRejected(t *testing.T) {
	p := EmptyProfileExtra()
	a, b := validWallet("0xa"), validWallet("0xb")
	a.IsPrimary = true
	b.IsPrimary = true
	p.Wallets = []ProfileWallet{a, b}
	if err := ValidateProfileExtra(p); err == nil || !strings.Contains(err.Error(), "primary") {
		t.Fatalf("expected single-primary error, got %v", err)
	}
}

func TestValidateProfileExtra_LabelTooLong(t *testing.T) {
	p := EmptyProfileExtra()
	w := validWallet("0xa")
	w.Label = strings.Repeat("a", 33)
	p.Wallets = []ProfileWallet{w}
	if err := ValidateProfileExtra(p); err == nil || !strings.Contains(err.Error(), "label") {
		t.Fatalf("expected label length error, got %v", err)
	}
}

func TestValidateProfileExtra_BadChainRejected(t *testing.T) {
	p := EmptyProfileExtra()
	w := validWallet("0xa")
	w.Chain = "bitcoin"
	p.Wallets = []ProfileWallet{w}
	if err := ValidateProfileExtra(p); err == nil || !strings.Contains(err.Error(), "chain") {
		t.Fatalf("expected chain error, got %v", err)
	}
}

func TestValidateProfileExtra_NoteTextTooLong(t *testing.T) {
	p := EmptyProfileExtra()
	p.Notes = []ProfileNote{{
		ID: "n_x", Text: strings.Repeat("x", 501), Kind: "trait", Source: "user",
		CapturedAt: time.Now().UTC(),
	}}
	if err := ValidateProfileExtra(p); err == nil || !strings.Contains(err.Error(), "note text") {
		t.Fatalf("expected note text length error, got %v", err)
	}
}

func TestValidateProfileExtra_NotesCap(t *testing.T) {
	p := EmptyProfileExtra()
	for i := 0; i < 51; i++ {
		p.Notes = append(p.Notes, ProfileNote{
			ID: "n_x", Text: "x", Kind: "trait", Source: "user", CapturedAt: time.Now().UTC(),
		})
	}
	if err := ValidateProfileExtra(p); err == nil || !strings.Contains(err.Error(), "note cap") {
		t.Fatalf("expected note cap error, got %v", err)
	}
}

func TestValidateProfileExtra_BadKindRejected(t *testing.T) {
	p := EmptyProfileExtra()
	p.Notes = []ProfileNote{{
		ID: "n_x", Text: "x", Kind: "medical", Source: "user", CapturedAt: time.Now().UTC(),
	}}
	if err := ValidateProfileExtra(p); err == nil || !strings.Contains(err.Error(), "kind") {
		t.Fatalf("expected kind error, got %v", err)
	}
}

func TestValidateProfileExtra_FavoriteTokensCap(t *testing.T) {
	p := EmptyProfileExtra()
	for i := 0; i < 11; i++ {
		p.Preferences.FavoriteTokens = append(p.Preferences.FavoriteTokens, "TKN")
	}
	if err := ValidateProfileExtra(p); err == nil || !strings.Contains(err.Error(), "favorite") {
		t.Fatalf("expected favorite_tokens cap error, got %v", err)
	}
}
