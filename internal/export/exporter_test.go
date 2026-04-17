package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []Format{FormatDotenv, FormatJSON, FormatExport} {
		_, err := New(f)
		if err != nil {
			t.Errorf("expected no error for format %q, got %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := New("xml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_JSON(t *testing.T) {
	e, _ := New(FormatJSON)
	tmp := filepath.Join(t.TempDir(), "out.json")
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}

	if err := e.Write(secrets, tmp); err != nil {
		t.Fatalf("Write: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	var got map[string]string
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got["FOO"] != "bar" || got["BAZ"] != "qux" {
		t.Errorf("unexpected json output: %v", got)
	}
}

func TestWrite_Dotenv(t *testing.T) {
	e, _ := New(FormatDotenv)
	tmp := filepath.Join(t.TempDir(), ".env")
	secrets := map[string]string{"KEY": "value", "OTHER": "hello world"}

	if err := e.Write(secrets, tmp); err != nil {
		t.Fatalf("Write: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	content := string(data)
	if !strings.Contains(content, "KEY=value") {
		t.Errorf("expected KEY=value in output")
	}
	if !strings.Contains(content, "OTHER=") {
		t.Errorf("expected OTHER= in output")
	}
}

func TestWrite_Export(t *testing.T) {
	e, _ := New(FormatExport)
	tmp := filepath.Join(t.TempDir(), "env.sh")
	secrets := map[string]string{"TOKEN": "abc123"}

	if err := e.Write(secrets, tmp); err != nil {
		t.Fatalf("Write: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), "export TOKEN=abc123") {
		t.Errorf("expected 'export TOKEN=abc123' in output, got: %s", data)
	}
}

func TestWrite_FilePermissions(t *testing.T) {
	e, _ := New(FormatDotenv)
	tmp := filepath.Join(t.TempDir(), ".env")
	_ = e.Write(map[string]string{"A": "b"}, tmp)

	info, err := os.Stat(tmp)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
