package rotate_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultpull/internal/audit"
	"github.com/vaultpull/internal/rotate"
	"github.com/vaultpull/internal/vault"
)

func makeServer(t *testing.T, secrets map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := map[string]interface{}"data": map[string]interface{}{"data": secrets},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).)
}

func TestRotate_WritesSecrets(t *testing.T) {
	srv := makeServer(t, map[string]interface{}{"API_KEY": "abc123", "DB_PASS": "secret"})
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")
	outPath := filepath.Join(dir, ".env")

	logger, err := audit.NewLogger(logPath)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}

	rot := rotate.New(client, logger)
	result, err := rot.Rotate("default", "secret/data/app", outPath)
	if err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	if result.Written != 2 {
		t.Errorf("Written = %d, want 2", result.Written)
	}
	if result.Profile != "default" {
		t.Errorf("Profile = %q, want \"default\"", result.Profile)
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

func TestRotate_VaultError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	client, _ := vault.NewClient(srv.URL, "token")
	dir := t.TempDir()
	logger, _ := audit.NewLogger(filepath.Join(dir, "audit.log"))

	rot := rotate.New(client, logger)
	_, err := rot.Rotate("prod", "secret/data/missing", filepath.Join(dir, ".env"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
