package env

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"LOG_LEVEL": "info",
		"PORT":     "8080",
	}
}

func TestValidateSchema_RequiredPresent(t *testing.T) {
	rules := []SchemaRule{
		{Key: "APP_ENV", Required: true},
	}
	results, err := ValidateSchema(baseEnv(), rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Passed {
		t.Errorf("expected rule to pass, got: %+v", results)
	}
}

func TestValidateSchema_RequiredMissing(t *testing.T) {
	rules := []SchemaRule{
		{Key: "SECRET_KEY", Required: true},
	}
	results, err := ValidateSchema(baseEnv(), rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Passed {
		t.Errorf("expected rule to fail for missing required key")
	}
}

func TestValidateSchema_PatternMatches(t *testing.T) {
	rules := []SchemaRule{
		{Key: "PORT", Pattern: `^\d+$`},
	}
	results, err := ValidateSchema(baseEnv(), rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Passed {
		t.Errorf("expected pattern to match, got reason: %s", results[0].Reason)
	}
}

func TestValidateSchema_PatternFails(t *testing.T) {
	rules := []SchemaRule{
		{Key: "APP_ENV", Pattern: `^\d+$`},
	}
	results, err := ValidateSchema(baseEnv(), rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Passed {
		t.Errorf("expected pattern to fail for non-numeric value")
	}
}

func TestValidateSchema_AllowedList(t *testing.T) {
	rules := []SchemaRule{
		{Key: "LOG_LEVEL", Allowed: []string{"debug", "info", "warn", "error"}},
	}
	results, err := ValidateSchema(baseEnv(), rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Passed {
		t.Errorf("expected allowed value to pass, got: %s", results[0].Reason)
	}
}

func TestValidateSchema_AllowedListFails(t *testing.T) {
	rules := []SchemaRule{
		{Key: "APP_ENV", Allowed: []string{"staging", "development"}},
	}
	results, err := ValidateSchema(baseEnv(), rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Passed {
		t.Errorf("expected value not in allowed list to fail")
	}
}

func TestValidateSchema_InvalidPattern(t *testing.T) {
	rules := []SchemaRule{
		{Key: "PORT", Pattern: `[invalid`},
	}
	_, err := ValidateSchema(baseEnv(), rules)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestHasSchemaErrors_True(t *testing.T) {
	results := []SchemaResult{{Key: "X", Passed: false, Reason: "missing"}}
	if !HasSchemaErrors(results) {
		t.Error("expected HasSchemaErrors to return true")
	}
}

func TestHasSchemaErrors_False(t *testing.T) {
	results := []SchemaResult{{Key: "X", Passed: true}}
	if HasSchemaErrors(results) {
		t.Error("expected HasSchemaErrors to return false")
	}
}
