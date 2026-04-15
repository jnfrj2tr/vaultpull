package profile

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/yourusername/vaultpull/internal/config"
)

const envVarProfileKey = "VAULTPULL_PROFILE"

// Select determines which profile name to use based on the following
// priority order:
//  1. Explicit flag value (non-empty flagValue)
//  2. VAULTPULL_PROFILE environment variable
//  3. The single profile defined in config (if exactly one exists)
//
// Returns an error if no profile can be determined.
func Select(cfg *config.Config, flagValue string) (string, error) {
	if strings.TrimSpace(flagValue) != "" {
		return flagValue, nil
	}

	if env := os.Getenv(envVarProfileKey); env != "" {
		return env, nil
	}

	names := ListNames(cfg)
	if len(names) == 0 {
		return "", errors.New("no profiles defined in config")
	}
	if len(names) == 1 {
		return names[0], nil
	}

	sort.Strings(names)
	return "", fmt.Errorf(
		"multiple profiles defined, please specify one with --profile or $%s. Available: %s",
		envVarProfileKey,
		strings.Join(names, ", "),
	)
}
