package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderString_BasicSecret(t *testing.T) {
	r := New(map[string]string{"DB_PASS": "hunter2"})
	out, err := r.RenderString(`DB_PASS={{ .Secrets.DB_PASS }}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "DB_PASS=hunter2" {
		t.Errorf("got %q", out)
	}
}

func TestRenderString_MissingKeyErrors(t *testing.T) {
	r := New(map[string]string{})
	_, err := r.RenderString(`{{ .Secrets.MISSING }}`)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRenderString_EnvAccess(t *testing.T) {
	t.Setenv("HOME_TEST_VAR", "testvalue")
	r := New(map[string]string{})
	out, err := r.RenderString(`VAL={{ .Env.HOME_TEST_VAR }}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "testvalue") {
		t.Errorf("expected env value in output, got %q", out)
	}
}

func TestRenderFile_WritesOutput(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "tmpl.env")
	dst := filepath.Join(dir, "out.env")

	content := "API_KEY={{ .Secrets.API_KEY }}\nAPP=myapp\n"
	if err := os.WriteFile(src, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	r := New(map[string]string{"API_KEY": "abc123"})
	if err := r.RenderFile(src, dst); err != nil {
		t.Fatalf("RenderFile error: %v", err)
	}

	b, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "API_KEY=abc123") {
		t.Errorf("rendered output missing expected value: %s", b)
	}
}

func TestRenderFile_MissingSource(t *testing.T) {
	r := New(map[string]string{})
	err := r.RenderFile("/nonexistent/tmpl.env", "/tmp/out.env")
	if err == nil {
		t.Fatal("expected error for missing source file")
	}
}

func TestRenderFile_BadDestination(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "tmpl.env")
	_ = os.WriteFile(src, []byte("KEY=val"), 0644)

	r := New(map[string]string{})
	err := r.RenderFile(src, "/nonexistent/dir/out.env")
	if err == nil {
		t.Fatal("expected error for bad destination")
	}
}
