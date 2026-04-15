package profile

import (
	"fmt"

	"github.com/yourusername/vaultpull/internal/config"
)

// ResolvedProfile holds the fully resolved settings for a given profile.
type ResolvedProfile struct {
	Name       string
	VaultPath  string
	OutputFile string
	Merge      bool
	Environment map[string]string
}

// Resolve looks up the named profile in cfg and returns a ResolvedProfile.
// It applies any top-level defaults before returning.
func Resolve(cfg *config.Config, profileName string) (*ResolvedProfile, error) {
	p, err := cfg.GetProfile(profileName)
	if err != nil {
		return nil, fmt.Errorf("profile %q not found: %w", profileName, err)
	}

	if p.VaultPath == "" {
		return nil, fmt.Errorf("profile %q has no vault_path configured", profileName)
	}

	outputFile := p.OutputFile
	if outputFile == "" {
		outputFile = ".env"
	}

	return &ResolvedProfile{
		Name:        profileName,
		VaultPath:   p.VaultPath,
		OutputFile:  outputFile,
		Merge:       p.Merge,
		Environment: p.Environment,
	}, nil
}

// ListNames returns all profile names defined in the config.
func ListNames(cfg *config.Config) []string {
	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	return names
}
