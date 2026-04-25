package env

import (
	"testing"
)

func TestClone_BasicCopy(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Clone(src, CloneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("unexpected values: %v", out)
	}
	// Ensure it is a copy, not a reference.
	out["FOO"] = "mutated"
	if src["FOO"] != "bar" {
		t.Error("Clone modified the source map")
	}
}

func TestClone_ExcludesKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	out, err := Clone(src, CloneOptions{Exclude: []string{"B"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["B"]; ok {
		t.Error("excluded key B should not be present")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestClone_AddsKeyPrefix(t *testing.T) {
	src := map[string]string{"NAME": "alice"}
	out, err := Clone(src, CloneOptions{KeyPrefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "alice" {
		t.Errorf("expected APP_NAME=alice, got %v", out)
	}
}

func TestClone_StripsKeyPrefix(t *testing.T) {
	src := map[string]string{"PROD_HOST": "localhost", "OTHER": "val"}
	out, err := Clone(src, CloneOptions{StripKeyPrefix: "PROD_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// PROD_HOST -> HOST; OTHER has no prefix so it is dropped.
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %v", out)
	}
	if _, ok := out["OTHER"]; ok {
		t.Error("OTHER should be dropped when StripKeyPrefix is set and key lacks the prefix")
	}
}

func TestClone_UppercaseKeys(t *testing.T) {
	src := map[string]string{"db_host": "127.0.0.1"}
	out, err := Clone(src, CloneOptions{UppercaseKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "127.0.0.1" {
		t.Errorf("expected DB_HOST, got %v", out)
	}
}

func TestClone_StripAndPrefix_Combined(t *testing.T) {
	src := map[string]string{"DEV_PORT": "5432"}
	out, err := Clone(src, CloneOptions{StripKeyPrefix: "DEV_", KeyPrefix: "PROD_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PROD_PORT"] != "5432" {
		t.Errorf("expected PROD_PORT=5432, got %v", out)
	}
}

func TestClone_EmptyKeyAfterTransform_ReturnsError(t *testing.T) {
	// Stripping the full key name leaves an empty string.
	src := map[string]string{"PREFIX_": "val"}
	_, err := Clone(src, CloneOptions{StripKeyPrefix: "PREFIX_"})
	if err == nil {
		t.Error("expected error for empty resulting key, got nil")
	}
}
