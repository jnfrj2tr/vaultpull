package env

import (
	"fmt"
	"strings"
)

// ScopeMode controls how scoping is applied to a map of env vars.
type ScopeMode string

const (
	// ScopeModeInclude keeps only keys matching the given scopes.
	ScopeModeInclude ScopeMode = "include"
	// ScopeModeExclude removes keys matching the given scopes.
	ScopeModeExclude ScopeMode = "exclude"
)

// ScopeOptions configures the Scope operation.
type ScopeOptions struct {
	Mode   ScopeMode
	Scopes []string // prefixes or exact key names
	Exact  bool     // if true, match full key name; otherwise match prefix
}

// Scope filters env vars by scope rules, returning a new map.
// In include mode only matching keys are kept; in exclude mode matching keys are removed.
func Scope(env map[string]string, opts ScopeOptions) (map[string]string, error) {
	if opts.Mode != ScopeModeInclude && opts.Mode != ScopeModeExclude {
		return nil, fmt.Errorf("scope: unknown mode %q, must be \"include\" or \"exclude\"", opts.Mode)
	}
	if len(opts.Scopes) == 0 {
		return nil, fmt.Errorf("scope: at least one scope must be provided")
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		matched := matchesAnyScope(k, opts.Scopes, opts.Exact)
		keep := (opts.Mode == ScopeModeInclude) == matched
		if keep {
			result[k] = v
		}
	}
	return result, nil
}

// ParseScopes parses a comma-separated list of scope tokens into a slice.
func ParseScopes(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func matchesAnyScope(key string, scopes []string, exact bool) bool {
	for _, s := range scopes {
		if exact {
			if key == s {
				return true
			}
		} else {
			if strings.HasPrefix(key, s) {
				return true
			}
		}
	}
	return false
}
