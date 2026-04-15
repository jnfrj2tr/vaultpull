package profile_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/profile"
)

func makeConfig() *config.Config {
	return &config.Config{
		Profiles: map[string]config.Profile{
			"dev": {
				VaultPath:  "secret/data/dev",
				OutputFile: ".env.dev",
				Merge:      true,
			},
			"prod": {
				VaultPath: "secret/data/prod",
			},
		},
	}
}

func TestResolve_KnownProfile(t *testing.T) {
	cfg := makeConfig()
	rp, err := profile.Resolve(cfg, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rp.VaultPath != "secret/data/dev" {
		t.Errorf("expected vault path 'secret/data/dev', got %q", rp.VaultPath)
	}
	if rp.OutputFile != ".env.dev" {
		t.Errorf("expected output file '.env.dev', got %q", rp.OutputFile)
	}
	if !rp.Merge {
		t.Error("expected merge to be true")
	}
}

func TestResolve_DefaultsOutputFile(t *testing.T) {
	cfg := makeConfig()
	rp, err := profile.Resolve(cfg, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rp.OutputFile != ".env" {
		t.Errorf("expected default output file '.env', got %q", rp.OutputFile)
	}
}

func TestResolve_UnknownProfile(t *testing.T) {
	cfg := makeConfig()
	_, err := profile.Resolve(cfg, "staging")
	if err == nil {
		t.Fatal("expected error for unknown profile, got nil")
	}
}

func TestResolve_MissingVaultPath(t *testing.T) {
	cfg := &config.Config{
		Profiles: map[string]config.Profile{
			"empty": {},
		},
	}
	_, err := profile.Resolve(cfg, "empty")
	if err == nil {
		t.Fatal("expected error for missing vault_path, got nil")
	}
}

func TestListNames_ReturnsAllProfiles(t *testing.T) {
	cfg := makeConfig()
	names := profile.ListNames(cfg)
	if len(names) != 2 {
		t.Errorf("expected 2 profile names, got %d", len(names))
	}
}
