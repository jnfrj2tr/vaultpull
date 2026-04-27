package env

import (
	"testing"
)

func TestAlias_BasicRename(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	rules := []AliasRule{{Source: "DB_HOST", Aliases: []string{"DATABASE_HOST"}, Keep: false}}
	out, err := Alias(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected alias to carry value, got %q", out["DATABASE_HOST"])
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("expected source key to be removed")
	}
}

func TestAlias_KeepSource(t *testing.T) {
	env := map[string]string{"API_KEY": "secret"}
	rules := []AliasRule{{Source: "API_KEY", Aliases: []string{"SERVICE_API_KEY"}, Keep: true}}
	out, err := Alias(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_KEY"] != "secret" {
		t.Error("expected source key to be retained")
	}
	if out["SERVICE_API_KEY"] != "secret" {
		t.Error("expected alias to be set")
	}
}

func TestAlias_MultipleAliases(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	rules := []AliasRule{{Source: "PORT", Aliases: []string{"APP_PORT", "HTTP_PORT"}, Keep: false}}
	out, err := Alias(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_PORT"] != "8080" || out["HTTP_PORT"] != "8080" {
		t.Error("expected both aliases to be set")
	}
}

func TestAlias_MissingSourceSkipped(t *testing.T) {
	env := map[string]string{"OTHER": "val"}
	rules := []AliasRule{{Source: "MISSING", Aliases: []string{"ALIAS"}, Keep: false}}
	out, err := Alias(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["ALIAS"]; ok {
		t.Error("expected alias not to be set when source is missing")
	}
}

func TestAlias_EmptyAliasReturnsError(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	rules := []AliasRule{{Source: "KEY", Aliases: []string{""}}, }
	_, err := Alias(env, rules)
	if err == nil {
		t.Error("expected error for empty alias name")
	}
}

func TestParseAliasRules_Valid(t *testing.T) {
	rules, err := ParseAliasRules([]string{"DB_HOST:DATABASE_HOST,PG_HOST+keep"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if !rules[0].Keep {
		t.Error("expected Keep=true")
	}
	if len(rules[0].Aliases) != 2 {
		t.Errorf("expected 2 aliases, got %d", len(rules[0].Aliases))
	}
}

func TestParseAliasRules_Invalid(t *testing.T) {
	_, err := ParseAliasRules([]string{"NOCORON"})
	if err == nil {
		t.Error("expected error for missing colon")
	}
}
