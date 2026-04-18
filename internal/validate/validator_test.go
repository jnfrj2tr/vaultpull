package validate_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/validate"
)

func TestValidate_RequiredPresent(t *testing.T) {
	secrets := map[string]string{"DB_URL": "postgres://localhost/db"}
	rules := []validate.Rule{{Key: "DB_URL", Required: true}}
	violations := validate.Validate(secrets, rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestValidate_RequiredMissing(t *testing.T) {
	secrets := map[string]string{}
	rules := []validate.Rule{{Key: "DB_URL", Required: true}}
	violations := validate.Validate(secrets, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DB_URL" {
		t.Errorf("unexpected key %q", violations[0].Key)
	}
}

func TestValidate_RequiredEmpty(t *testing.T) {
	secrets := map[string]string{"API_KEY": "   "}
	rules := []validate.Rule{{Key: "API_KEY", Required: true}}
	violations := validate.Validate(secrets, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_PatternMatches(t *testing.T) {
	secrets := map[string]string{"PORT": "8080"}
	rules := []validate.Rule{{Key: "PORT", Pattern: `^\d+$`}}
	violations := validate.Validate(secrets, rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestValidate_PatternFails(t *testing.T) {
	secrets := map[string]string{"PORT": "not-a-number"}
	rules := []validate.Rule{{Key: "PORT", Pattern: `^\d+$`}}
	violations := validate.Validate(secrets, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_InvalidPattern(t *testing.T) {
	secrets := map[string]string{"KEY": "value"}
	rules := []validate.Rule{{Key: "KEY", Pattern: `[invalid`}}
	violations := validate.Validate(secrets, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for bad regex, got %d", len(violations))
	}
}

func TestValidate_NoRules(t *testing.T) {
	secrets := map[string]string{"X": "y"}
	violations := validate.Validate(secrets, nil)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}
