package env

import (
	"sort"
	"strings"
)

// CompactOptions controls how compaction is performed.
type CompactOptions struct {
	// RemoveEmpty drops keys whose values are empty strings.
	RemoveEmpty bool
	// RemoveDuplicateValues drops all but the first key that shares a value.
	RemoveDuplicateValues bool
	// TrimSpace trims leading/trailing whitespace from every value.
	TrimSpace bool
}

// Compact cleans up an env map according to the supplied options and returns
// a new map.  The original is never modified.
func Compact(env map[string]string, opts CompactOptions) (map[string]string, error) {
	out := make(map[string]string, len(env))

	// Collect keys in deterministic order so duplicate-value removal is stable.
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	seen := make(map[string]struct{})

	for _, k := range keys {
		v := env[k]

		if opts.TrimSpace {
			v = strings.TrimSpace(v)
		}

		if opts.RemoveEmpty && v == "" {
			continue
		}

		if opts.RemoveDuplicateValues {
			if _, exists := seen[v]; exists {
				continue
			}
			seen[v] = struct{}{}
		}

		out[k] = v
	}

	return out, nil
}
