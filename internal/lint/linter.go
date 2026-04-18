// Package lint provides validation of .vaultpull.yaml configuration files,
// warning about common misconfigurations before a sync is attempted.
package lint

import (
	"fmt"
	"strings"

	"github.com/yourorg/vaultpull/internal/config"
)

// Issue represents a single lint warning or error.
type Issue struct {
	Level   string // "warn" or "error"
	Profile string
	Message string
}

func (i Issue) String() string {
	if i.Profile != "" {
		return fmt.Sprintf("[%s] profile %q: %s", i.Level, i.Profile, i.Message)
	}
	return fmt.Sprintf("[%s] %s", i.Level, i.Message)
}

// Run lints cfg and returns any issues found.
func Run(cfg *config.Config) []Issue {
	var issues []Issue

	if cfg.Vault.Address == "" {
		issues = append(issues, Issue{Level: "error", Message: "vault.address is required"})
	}

	if len(cfg.Profiles) == 0 {
		issues = append(issues, Issue{Level: "error", Message: "no profiles defined"})
		return issues
	}

	seen := map[string]bool{}
	for _, p := range cfg.Profiles {
		name := p.Name
		if seen[name] {
			issues = append(issues, Issue{Level: "error", Profile: name, Message: "duplicate profile name"})
		}
		seen[name] = true

		if p.VaultPath == "" {
			issues = append(issues, Issue{Level: "error", Profile: name, Message: "vault_path is required"})
		}

		if p.OutputFile == "" {
			issues = append(issues, Issue{Level: "warn", Profile: name, Message: "output_file not set; will default to .env"})
		}

		if strings.Contains(p.VaultPath, " ") {
			issues = append(issues, Issue{Level: "warn", Profile: name, Message: "vault_path contains spaces"})
		}
	}

	return issues
}

// HasErrors returns true if any issue has level "error".
func HasErrors(issues []Issue) bool {
	for _, i := range issues {
		if i.Level == "error" {
			return true
		}
	}
	return false
}
