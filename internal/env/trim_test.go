package env

import (
	"testing"
)

func TestTrim_TrimKeys(t *testing.T) {
	input := map[string]string{
		"  APP_HOST  ": "localhost",
		"APP_PORT": "8080",
	}
	out, err := Trim(input, TrimOptions{TrimKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST key after trimming")
	}
	if _, ok := out["APP_PORT"]; !ok {
		t.Error("expected APP_PORT key to remain")
	}
}

func TestTrim_TrimValues(t *testing.T) {
	input := map[string]string{
		"DB_HOST": "  postgres  ",
		"DB_PORT": "5432",
	}
	out, err := Trim(input, TrimOptions{TrimValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "postgres" {
		t.Errorf("expected 'postgres', got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected '5432', got %q", out["DB_PORT"])
	}
}

func TestTrim_TrimPrefix(t *testing.T) {
	input := map[string]string{
		"PROD_APP_HOST": "example.com",
		"PROD_APP_PORT": "443",
		"DEV_APP_HOST":  "localhost",
	}
	out, err := Trim(input, TrimOptions{TrimPrefix: "PROD_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST after stripping PROD_ prefix")
	}
	if _, ok := out["APP_PORT"]; !ok {
		t.Error("expected APP_PORT after stripping PROD_ prefix")
	}
	if _, ok := out["DEV_APP_HOST"]; !ok {
		t.Error("expected DEV_APP_HOST to remain unchanged")
	}
}

func TestTrim_TrimSuffix(t *testing.T) {
	input := map[string]string{
		"APP_HOST_V1": "host-a",
		"APP_PORT_V1": "9000",
		"APP_NAME":    "myapp",
	}
	out, err := Trim(input, TrimOptions{TrimSuffix: "_V1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST after stripping _V1 suffix")
	}
	if _, ok := out["APP_NAME"]; !ok {
		t.Error("expected APP_NAME to remain unchanged")
	}
}

func TestTrim_EmptyKeyDropped(t *testing.T) {
	input := map[string]string{
		"PREFIX_": "value",
	}
	out, err := Trim(input, TrimOptions{TrimPrefix: "PREFIX_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestTrim_NoOptions_PassThrough(t *testing.T) {
	input := map[string]string{
		"KEY": "  value  ",
	}
	out, err := Trim(input, TrimOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "  value  " {
		t.Errorf("expected value unchanged, got %q", out["KEY"])
	}
}
