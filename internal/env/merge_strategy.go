package env

import "fmt"

// Strategy defines how conflicts are resolved when merging env maps.
type Strategy int

const (
	// StrategyOverwrite replaces existing keys with incoming values.
	StrategyOverwrite Strategy = iota
	// StrategyKeepExisting preserves existing keys, only adds new ones.
	StrategyKeepExisting
	// StrategyError returns an error on any conflicting key.
	StrategyError
)

// ParseStrategy converts a string name to a Strategy.
func ParseStrategy(s string) (Strategy, error) {
	switch s {
	case "overwrite":
		return StrategyOverwrite, nil
	case "keep", "keep-existing":
		return StrategyKeepExisting, nil
	case "error":
		return StrategyError, nil
	default:
		return 0, fmt.Errorf("unknown merge strategy %q: must be overwrite, keep, or error", s)
	}
}

// MergeWithStrategy merges incoming into base according to the given strategy.
// base is modified in place and returned.
func MergeWithStrategy(base, incoming map[string]string, strategy Strategy) (map[string]string, error) {
	if base == nil {
		base = make(map[string]string)
	}
	for k, v := range incoming {
		existing, exists := base[k]
		switch strategy {
		case StrategyOverwrite:
			base[k] = v
		case StrategyKeepExisting:
			if !exists {
				base[k] = v
			}
		case StrategyError:
			if exists && existing != v {
				return nil, fmt.Errorf("conflict on key %q: existing=%q incoming=%q", k, existing, v)
			}
			base[k] = v
		}
	}
	return base, nil
}
