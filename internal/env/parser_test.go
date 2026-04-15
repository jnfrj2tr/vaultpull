package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnvFile: %v", err)
	}
	return p
}

func TestParse_BasicKeyValue(t *testing.T) {
	p := writeTempEnvFile(t, "FOO=bar\nBAZ=qux\n")
	entries, err := Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "FOO" || entries[0].Value != "bar" {
		t.Errorf("unexpected entry[0]: %+v", entries[0])
	}
}

func TestParse_SkipsComments(t *testing.T) {
	p := writeTempEnvFile(t, "# comment\nKEY=value\n")
	entries, err := Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestParse_StripsDoubleQuotes(t *testing.T) {
	p := writeTempEnvFile(t, `TOKEN="my secret token"`+"\n")
	entries, err := Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "my secret token" {
		t.Errorf("expected unquoted value, got %q", entries[0].Value)
	}
}

func TestParse_StripsSingleQuotes(t *testing.T) {
	p := writeTempEnvFile(t, "KEY='hello world'\n")
	entries, err := Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "hello world" {
		t.Errorf("expected unquoted value, got %q", entries[0].Value)
	}
}

func TestParse_MissingSeparatorReturnsError(t *testing.T) {
	p := writeTempEnvFile(t, "INVALIDLINE\n")
	_, err := Parse(p)
	if err == nil {
		t.Fatal("expected error for missing '=', got nil")
	}
}

func TestParse_MissingFileReturnsError(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParse_EmptyFileReturnsNoEntries(t *testing.T) {
	p := writeTempEnvFile(t, "")
	entries, err := Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
