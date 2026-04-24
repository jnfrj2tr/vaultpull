package env

import (
	"testing"
)

func TestRedactMap_ExplicitKey(t *testing.T) {
	env := map[string]string{
		"API_KEY": "secret123",
		"HOST":    "localhost",
	}
	out, err := RedactMap(env, RedactOptions{Keys: []string{"API_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_KEY"] != "***" {
		t.Errorf("expected API_KEY to be redacted, got %q", out["API_KEY"])
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST to be unchanged, got %q", out["HOST"])
	}
}

func TestRedactMap_PatternMatch(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "hunter2",
		"DB_HOST":     "db.local",
		"APP_SECRET":  "topsecret",
	}
	out, err := RedactMap(env, RedactOptions{Patterns: []string{"(?i)password|secret"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASSWORD"] != "***" {
		t.Errorf("expected DB_PASSWORD redacted, got %q", out["DB_PASSWORD"])
	}
	if out["APP_SECRET"] != "***" {
		t.Errorf("expected APP_SECRET redacted, got %q", out["APP_SECRET"])
	}
	if out["DB_HOST"] != "db.local" {
		t.Errorf("expected DB_HOST unchanged, got %q", out["DB_HOST"])
	}
}

func TestRedactMap_PartialReveal(t *testing.T) {
	env := map[string]string{"TOKEN": "abcdef9999"}
	out, err := RedactMap(env, RedactOptions{Keys: []string{"TOKEN"}, PartialReveal: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN"] != "abc***" {
		t.Errorf("expected partial redact, got %q", out["TOKEN"])
	}
}

func TestRedactMap_PartialReveal_ExceedsLength(t *testing.T) {
	env := map[string]string{"K": "hi"}
	out, err := RedactMap(env, RedactOptions{Keys: []string{"K"}, PartialReveal: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["K"] != "***" {
		t.Errorf("expected full redact when reveal >= len, got %q", out["K"])
	}
}

func TestRedactMap_InvalidPattern(t *testing.T) {
	env := map[string]string{"X": "y"}
	_, err := RedactMap(env, RedactOptions{Patterns: []string{"[invalid"}})
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestRedactMap_EmptyOptions(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	out, err := RedactMap(env, RedactOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Error("expected map unchanged when no redact rules given")
	}
}

func TestRedactMap_CaseInsensitiveKey(t *testing.T) {
	env := map[string]string{"api_key": "secret"}
	out, err := RedactMap(env, RedactOptions{Keys: []string{"API_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["api_key"] != "***" {
		t.Errorf("expected case-insensitive key match, got %q", out["api_key"])
	}
}
