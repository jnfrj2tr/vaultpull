package env

import (
	"os"
	"testing"
)

func TestInterpolate_ResolvesFromSecrets(t *testing.T) {
	secrets := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api",
	}
	out, err := Interpolate(secrets, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_URL"] != "https://example.com/api" {
		t.Errorf("got %q", out["API_URL"])
	}
}

func TestInterpolate_FallsBackToEnv(t *testing.T) {
	os.Setenv("MY_HOST", "localhost")
	defer os.Unsetenv("MY_HOST")

	secrets := map[string]string{
		"DSN": "postgres://${MY_HOST}/db",
	}
	out, err := Interpolate(secrets, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "postgres://localhost/db" {
		t.Errorf("got %q", out["DSN"])
	}
}

func TestInterpolate_EnvDisabled_ReturnsError(t *testing.T) {
	secrets := map[string]string{
		"VAL": "${MISSING_VAR}",
	}
	_, err := Interpolate(secrets, false)
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestInterpolate_DefaultValue(t *testing.T) {
	secrets := map[string]string{
		"VAL": "${MISSING:-fallback}",
	}
	out, err := Interpolate(secrets, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["VAL"] != "fallback" {
		t.Errorf("got %q", out["VAL"])
	}
}

func TestInterpolate_NoReferences(t *testing.T) {
	secrets := map[string]string{
		"PLAIN": "hello",
	}
	out, err := Interpolate(secrets, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PLAIN"] != "hello" {
		t.Errorf("got %q", out["PLAIN"])
	}
}

func TestInterpolate_EmptyMap(t *testing.T) {
	out, err := Interpolate(map[string]string{}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map")
	}
}
