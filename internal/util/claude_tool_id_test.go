package util

import (
	"regexp"
	"testing"
)

func TestSanitizeClaudeToolID_KeepsValidID(t *testing.T) {
	if got := SanitizeClaudeToolID("toolu_01-AbC9"); got != "toolu_01-AbC9" {
		t.Fatalf("valid id changed: %q", got)
	}
}

func TestSanitizeClaudeToolID_ReplacesInvalidCharacters(t *testing.T) {
	if got := SanitizeClaudeToolID("call.1:a b/ü"); got != "call_1_a_b__" {
		t.Fatalf("unexpected sanitized id: %q", got)
	}
}

func TestSanitizeClaudeToolID_EmptyGetsUniqueFallback(t *testing.T) {
	valid := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	a := SanitizeClaudeToolID("")
	b := SanitizeClaudeToolID("")
	if !valid.MatchString(a) || !valid.MatchString(b) {
		t.Fatalf("fallback ids must match Claude regex: %q %q", a, b)
	}
	if a == b {
		t.Fatalf("fallback ids must be unique, got %q twice", a)
	}
}
