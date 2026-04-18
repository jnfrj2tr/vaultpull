package lint_test

import (
	"testing"

	"github.com/yourorg/vaultpull/internal/config"
	"github.com/yourorg/vaultpull/internal/lint"
)

func baseConfig() *config.Config {
	return &config.Config{
		Vault: config.VaultConfig{Address: "http://127.0.0.1:8200"},
		Profiles: []config.Profile{
			{Name: "dev", VaultPath: "secret/dev", OutputFile: ".env.dev"},
		},
	}
}

func TestRun_Clean(t *testing.T) {
	issues := lint.Run(baseConfig())
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestRun_MissingVaultAddress(t *testing.T) {
	cfg := baseConfig()
	cfg.Vault.Address = ""
	issues := lint.Run(cfg)
	if !lint.HasErrors(issues) {
		t.Fatal("expected error for missing vault address")
	}
}

func TestRun_NoProfiles(t *testing.T) {
	cfg := baseConfig()
	cfg.Profiles = nil
	issues := lint.Run(cfg)
	if !lint.HasErrors(issues) {
		t.Fatal("expected error for no profiles")
	}
}

func TestRun_MissingVaultPath(t *testing.T) {
	cfg := baseConfig()
	cfg.Profiles[0].VaultPath = ""
	issues := lint.Run(cfg)
	if !lint.HasErrors(issues) {
		t.Fatal("expected error for missing vault_path")
	}
}

func TestRun_MissingOutputFileIsWarn(t *testing.T) {
	cfg := baseConfig()
	cfg.Profiles[0].OutputFile = ""
	issues := lint.Run(cfg)
	if lint.HasErrors(issues) {
		t.Fatal("expected only warnings, not errors")
	}
	if len(issues) == 0 {
		t.Fatal("expected at least one warning")
	}
}

func TestRun_DuplicateProfileName(t *testing.T) {
	cfg := baseConfig()
	cfg.Profiles = append(cfg.Profiles, cfg.Profiles[0])
	issues := lint.Run(cfg)
	if !lint.HasErrors(issues) {
		t.Fatal("expected error for duplicate profile name")
	}
}

func TestIssue_String_WithProfile(t *testing.T) {
	i := lint.Issue{Level: "warn", Profile: "dev", Message: "something"}
	s := i.String()
	if s != `[warn] profile "dev": something` {
		t.Fatalf("unexpected string: %s", s)
	}
}
