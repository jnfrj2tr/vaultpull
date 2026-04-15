package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/vaultpull/internal/config"
)

const sampleYAML = `
profiles:
  default:
    vault_addr: "http://127.0.0.1:8200"
    vault_token: "root"
    mount_path: "secret"
    secret_path: "myapp/prod"
    output_file: ".env"
    mapping:
      DB_PASSWORD: db_password
      API_KEY: api_key
  staging:
    vault_addr: "http://staging-vault:8200"
    vault_token: "staging-token"
    mount_path: "secret"
    secret_path: "myapp/staging"
    output_file: ".env.staging"
`

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".vaultpull.yaml")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return p
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, sampleYAML)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Profiles) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(cfg.Profiles))
	}
}

func TestGetProfile_Exists(t *testing.T) {
	path := writeTempConfig(t, sampleYAML)
	cfg, _ := config.Load(path)
	p, err := cfg.GetProfile("default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("unexpected vault_addr: %s", p.VaultAddr)
	}
	if p.Mapping["DB_PASSWORD"] != "db_password" {
		t.Errorf("unexpected mapping value for DB_PASSWORD")
	}
}

func TestGetProfile_Missing(t *testing.T) {
	path := writeTempConfig(t, sampleYAML)
	cfg, _ := config.Load(path)
	_, err := cfg.GetProfile("nonexistent")
	if err == nil {
		t.Error("expected error for missing profile, got nil")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/.vaultpull.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
