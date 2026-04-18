package env

import (
	"testing"
)

func TestApplyDefaults_InsertsDefault(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	defaults := []DefaultEntry{{Key: "BAZ", Default: "qux"}}

	out, err := ApplyDefaults(secrets, defaults)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", out["BAZ"])
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", out["FOO"])
	}
}

func TestApplyDefaults_DoesNotOverwriteExisting(t *testing.T) {
	secrets := map[string]string{"FOO": "original"}
	defaults := []DefaultEntry{{Key: "FOO", Default: "override"}}

	out, err := ApplyDefaults(secrets, defaults)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "original" {
		t.Errorf("expected original, got %q", out["FOO"])
	}
}

func TestApplyDefaults_RequiredMissingReturnsError(t *testing.T) {
	secrets := map[string]string{}
	defaults := []DefaultEntry{{Key: "SECRET", Required: true}}

	_, err := ApplyDefaults(secrets, defaults)
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
}

func TestApplyDefaults_RequiredPresentNoError(t *testing.T) {
	secrets := map[string]string{"SECRET": "val"}
	defaults := []DefaultEntry{{Key: "SECRET", Required: true}}

	_, err := ApplyDefaults(secrets, defaults)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseDefaults_Required(t *testing.T) {
	raw := map[string]string{"TOKEN": "required"}
	entries := ParseDefaults(raw)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if !entries[0].Required {
		t.Error("expected Required=true")
	}
}

func TestParseDefaults_WithDefault(t *testing.T) {
	raw := map[string]string{"REGION": "us-east-1"}
	entries := ParseDefaults(raw)
	if entries[0].Default != "us-east-1" {
		t.Errorf("expected us-east-1, got %q", entries[0].Default)
	}
}
