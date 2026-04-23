package env

import (
	"testing"
)

func TestCoerce_Bool(t *testing.T) {
	env := map[string]string{"ENABLED": "1", "DISABLED": "0"}
	rules := []CoerceRule{
		{Key: "ENABLED", Type: CoerceBool},
		{Key: "DISABLED", Type: CoerceBool},
	}
	out, err := Coerce(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ENABLED"] != "true" {
		t.Errorf("expected true, got %q", out["ENABLED"])
	}
	if out["DISABLED"] != "false" {
		t.Errorf("expected false, got %q", out["DISABLED"])
	}
}

func TestCoerce_Int(t *testing.T) {
	env := map[string]string{"PORT": "8080.9"}
	rules := []CoerceRule{{Key: "PORT", Type: CoerceInt}}
	out, err := Coerce(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected 8080, got %q", out["PORT"])
	}
}

func TestCoerce_Float(t *testing.T) {
	env := map[string]string{"RATE": "3.14"}
	rules := []CoerceRule{{Key: "RATE", Type: CoerceFloat}}
	out, err := Coerce(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["RATE"] != "3.14" {
		t.Errorf("expected 3.14, got %q", out["RATE"])
	}
}

func TestCoerce_UpperLower(t *testing.T) {
	env := map[string]string{"ENV": "Production", "MODE": "DEBUG"}
	rules := []CoerceRule{
		{Key: "ENV", Type: CoerceLower},
		{Key: "MODE", Type: CoerceLower},
	}
	out, err := Coerce(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ENV"] != "production" {
		t.Errorf("expected production, got %q", out["ENV"])
	}
	if out["MODE"] != "debug" {
		t.Errorf("expected debug, got %q", out["MODE"])
	}
}

func TestCoerce_MissingKeySkipped(t *testing.T) {
	env := map[string]string{"A": "hello"}
	rules := []CoerceRule{{Key: "MISSING", Type: CoerceBool}}
	out, err := Coerce(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "hello" {
		t.Errorf("expected A to be unchanged")
	}
}

func TestCoerce_InvalidBoolReturnsError(t *testing.T) {
	env := map[string]string{"FLAG": "maybe"}
	rules := []CoerceRule{{Key: "FLAG", Type: CoerceBool}}
	_, err := Coerce(env, rules)
	if err == nil {
		t.Fatal("expected error for invalid bool")
	}
}

func TestParseCoerceRules_Valid(t *testing.T) {
	rules, err := ParseCoerceRules([]string{"PORT:int", "DEBUG:bool", "NAME:upper"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
	if rules[0].Key != "PORT" || rules[0].Type != CoerceInt {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
}

func TestParseCoerceRules_Invalid(t *testing.T) {
	_, err := ParseCoerceRules([]string{"NOCOLON"})
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestCoerce_UnknownTypeReturnsError(t *testing.T) {
	env := map[string]string{"X": "val"}
	rules := []CoerceRule{{Key: "X", Type: CoerceType("base64")}}
	_, err := Coerce(env, rules)
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}
