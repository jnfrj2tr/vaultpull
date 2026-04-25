package env

import (
	"fmt"
	"strings"
)

// CloneOptions controls how a map is cloned and optionally transformed.
type CloneOptions struct {
	// KeyPrefix adds a prefix to all keys in the output.
	KeyPrefix string
	// StripKeyPrefix removes a prefix from all keys before cloning.
	StripKeyPrefix string
	// UppercaseKeys normalizes all keys to uppercase.
	UppercaseKeys bool
	// Exclude is a list of exact keys to omit from the clone.
	Exclude []string
}

// Clone returns a deep copy of src, applying any transformations defined in
// opts. Keys listed in Exclude are dropped. StripKeyPrefix is applied before
// KeyPrefix so the final key is: KeyPrefix + (original - StripKeyPrefix).
func Clone(src map[string]string, opts CloneOptions) (map[string]string, error) {
	excludeSet := make(map[string]struct{}, len(opts.Exclude))
	for _, k := range opts.Exclude {
		excludeSet[k] = struct{}{}
	}

	out := make(map[string]string, len(src))
	for k, v := range src {
		if _, skip := excludeSet[k]; skip {
			continue
		}

		newKey := k

		if opts.StripKeyPrefix != "" {
			if !strings.HasPrefix(newKey, opts.StripKeyPrefix) {
				continue
			}
			newKey = strings.TrimPrefix(newKey, opts.StripKeyPrefix)
		}

		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		}

		newKey = opts.KeyPrefix + newKey

		if newKey == "" {
			return nil, fmt.Errorf("clone: key transformation produced an empty key for original key %q", k)
		}

		out[newKey] = v
	}

	return out, nil
}
