package env

import (
	"testing"
)

func TestCheckTypes_AllValid(t *testing.T) {
	env := map[string]string{
		"PORT":    "8080",
		"RATIO":   "0.75",
		"ENABLED": "true",
		"HOST":    "localhost",
	}
	rules := []TypeRule{
		{Key: "PORT", Expected: "int"},
		{Key: "RATIO", Expected: "float"},
		{Key: "ENABLED", Expected: "bool"},
		{Key: "HOST", Expected: "string"},
	}
	violations := CheckTypes(env, rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestCheckTypes_IntViolation(t *testing.T) {
	env := map[string]string{"PORT": "not-a-number"}
	rules := []TypeRule{{Key: "PORT", Expected: "int"}}
	violations := CheckTypes(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "PORT" {
		t.Errorf("unexpected key %q", violations[0].Key)
	}
}

func TestCheckTypes_MissingKeySkipped(t *testing.T) {
	env := map[string]string{}
	rules := []TypeRule{{Key: "MISSING", Expected: "int"}}
	violations := CheckTypes(env, rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations for missing key, got %v", violations)
	}
}

func TestCheckTypes_URLAndEmail(t *testing.T) {
	env := map[string]string{
		"WEBHOOK": "https://example.com/hook",
		"EMAIL":   "user@example.com",
		"BAD_URL": "ftp://old.example.com",
	}
	rules := []TypeRule{
		{Key: "WEBHOOK", Expected: "url"},
		{Key: "EMAIL", Expected: "email"},
		{Key: "BAD_URL", Expected: "url"},
	}
	violations := CheckTypes(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(violations), violations)
	}
	if violations[0].Key != "BAD_URL" {
		t.Errorf("expected BAD_URL violation, got %q", violations[0].Key)
	}
}

func TestParseTypeRules_Valid(t *testing.T) {
	rules, err := ParseTypeRules([]string{"PORT=int", "ENABLED=bool"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Key != "PORT" || rules[0].Expected != "int" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
}

func TestParseTypeRules_Invalid(t *testing.T) {
	_, err := ParseTypeRules([]string{"NOTYPE"})
	if err == nil {
		t.Fatal("expected error for malformed rule")
	}
}

func TestCheckTypes_BoolVariants(t *testing.T) {
	cases := []string{"true", "false", "1", "0", "TRUE", "FALSE"}
	for _, v := range cases {
		env := map[string]string{"FLAG": v}
		violations := CheckTypes(env, []TypeRule{{Key: "FLAG", Expected: "bool"}})
		if len(violations) != 0 {
			t.Errorf("value %q should be valid bool, got violation", v)
		}
	}
}
