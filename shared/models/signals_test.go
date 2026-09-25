package models

import (
	"strings"
	"testing"
)

func TestValidateSignalCreate_HappyAlert(t *testing.T) {
	if err := ValidateSignalCreate("alert", 15); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestValidateSignalCreate_BadKind(t *testing.T) {
	err := ValidateSignalCreate("watcher", 15)
	if err == nil || !strings.Contains(err.Error(), "kind") {
		t.Fatalf("expected kind error, got %v", err)
	}
}

func TestValidateSignalCreate_BadCadence(t *testing.T) {
	err := ValidateSignalCreate("alert", 7)
	if err == nil || !strings.Contains(err.Error(), "cadence") {
		t.Fatalf("expected cadence error, got %v", err)
	}
}

func TestValidateSignalCreate_InterestCadenceFloor(t *testing.T) {
	err := ValidateSignalCreate("interest", 60)
	if err == nil || !strings.Contains(err.Error(), "interest") {
		t.Fatalf("expected interest floor error, got %v", err)
	}
	if err := ValidateSignalCreate("interest", 720); err != nil {
		t.Fatalf("720 must be valid for interest: %v", err)
	}
}
