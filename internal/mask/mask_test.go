package mask_test

import (
	"testing"

	"github.com/yourorg/vaultpull/internal/mask"
)

func TestRedact_AllValuesReplaced(t *testing.T) {
	input := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
	}
	out := mask.Redact(input)
	for k, v := range out {
		if v != "****" {
			t.Errorf("key %s: expected ****, got %s", k, v)
		}
	}
	if len(out) != len(input) {
		t.Errorf("expected %d keys, got %d", len(input), len(out))
	}
}

func TestRedact_EmptyMap(t *testing.T) {
	out := mask.Redact(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map")
	}
}

func TestRedactValue_NonEmpty(t *testing.T) {
	if got := mask.RedactValue("secret"); got != "****" {
		t.Errorf("expected ****, got %s", got)
	}
}

func TestRedactValue_Empty(t *testing.T) {
	if got := mask.RedactValue(""); got != "" {
		t.Errorf("expected empty string, got %s", got)
	}
}

func TestPartialRedact_ShowsPrefix(t *testing.T) {
	got := mask.PartialRedact("supersecret", 3)
	if got != "sup********" {
		t.Errorf("unexpected result: %s", got)
	}
}

func TestPartialRedact_NZero(t *testing.T) {
	got := mask.PartialRedact("secret", 0)
	if got != "****" {
		t.Errorf("expected ****, got %s", got)
	}
}

func TestPartialRedact_NExceedsLength(t *testing.T) {
	got := mask.PartialRedact("hi", 100)
	if got != "****" {
		t.Errorf("expected ****, got %s", got)
	}
}

func TestPartialRedact_EmptyValue(t *testing.T) {
	if got := mask.PartialRedact("", 2); got != "" {
		t.Errorf("expected empty string, got %s", got)
	}
}
