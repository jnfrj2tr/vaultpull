package env

import (
	"strings"
)

// TrimOptions controls which trimming operations are applied.
type TrimOptions struct {
	// TrimKeys removes leading/trailing whitespace from all keys.
	TrimKeys bool
	// TrimValues removes leading/trailing whitespace from all values.
	TrimValues bool
	// TrimPrefix removes a specific prefix from all keys (case-sensitive).
	TrimPrefix string
	// TrimSuffix removes a specific suffix from all keys (case-sensitive).
	TrimSuffix string
}

// Trim applies whitespace and affix trimming to a map of environment variables.
// Keys that become empty after trimming are dropped. If two keys collide after
// trimming, the last value encountered (in iteration order) wins.
func Trim(env map[string]string, opts TrimOptions) (map[string]string, error) {
	out := make(map[string]string, len(env))

	for k, v := range env {
		newKey := k
		newVal := v

		if opts.TrimKeys {
			newKey = strings.TrimSpace(newKey)
		}
		if opts.TrimPrefix != "" {
			newKey = strings.TrimPrefix(newKey, opts.TrimPrefix)
		}
		if opts.TrimSuffix != "" {
			newKey = strings.TrimSuffix(newKey, opts.TrimSuffix)
		}
		if opts.TrimValues {
			newVal = strings.TrimSpace(newVal)
		}

		if newKey == "" {
			continue
		}

		out[newKey] = newVal
	}

	return out, nil
}
