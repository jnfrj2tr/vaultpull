package profile_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/profile"
)

func singleProfileConfig() *config.Config {
	return &config.Config{
		Profiles: map[string]config.Profile{
			"dev": {VaultPath: "secret/data/dev"},
		},
	}
}

func TestSelect_FlagTakesPriority(t *testing.T) {
	cfg := makeConfig()
	t.Setenv("VAULTPULL_PROFILE", "prod")
	name, err := profile.Select(cfg, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "dev" {
		t.Errorf("expected 'dev', got %q", name)
	}
}

func TestSelect_EnvVarFallback(t *testing.T) {
	cfg := makeConfig()
	t.Setenv("VAULTPULL_PROFILE", "prod")
	name, err := profile.Select(cfg, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "prod" {
		t.Errorf("expected 'prod', got %q", name)
	}
}

func TestSelect_SingleProfileAutoSelected(t *testing.T) {
	cfg := singleProfileConfig()
	t.Setenv("VAULTPULL_PROFILE", "")
	name, err := profile.Select(cfg, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "dev" {
		t.Errorf("expected 'dev', got %q", name)
	}
}

func TestSelect_MultipleProfilesNoFlag_ReturnsError(t *testing.T) {
	cfg := makeConfig()
	t.Setenv("VAULTPULL_PROFILE", "")
	_, err := profile.Select(cfg, "")
	if err == nil {
		t.Fatal("expected error for ambiguous profile, got nil")
	}
}

func TestSelect_NoProfiles_ReturnsError(t *testing.T) {
	cfg := &config.Config{Profiles: map[string]config.Profile{}}
	t.Setenv("VAULTPULL_PROFILE", "")
	_, err := profile.Select(cfg, "")
	if err == nil {
		t.Fatal("expected error for empty profiles, got nil")
	}
}
