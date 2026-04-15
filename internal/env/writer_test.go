package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWrite_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := NewWriter(path)
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", content)
	}
	if !strings.Contains(content, "DB_PORT=5432") {
		t.Errorf("expected DB_PORT=5432 in output, got:\n%s", content)
	}
}

func TestWrite_QuotesSpecialValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := NewWriter(path)
	secrets := map[string]string{
		"MSG": "hello world",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", string(data))
	}
}

func TestWrite_InvalidPath(t *testing.T) {
	w := NewWriter("/nonexistent/dir/.env")
	err := w.Write(map[string]string{"KEY": "val"})
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestQuoteValue(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"has space", `"has space"`},
		{"has#hash", `"has#hash"`},
		{"normal123", "normal123"},
	}
	for _, c := range cases {
		got := quoteValue(c.input)
		if got != c.expected {
			t.Errorf("quoteValue(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}
