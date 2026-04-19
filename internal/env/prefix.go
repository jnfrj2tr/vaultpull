package env

import "strings"

// PrefixOptions controls how prefix addition or stripping is applied.
type PrefixOptions struct {
	Prefix    string
	Strip     bool
	Overwrite bool
}

// AddPrefix returns a new map with the given prefix added to every key.
// If Overwrite is false and a prefixed key already exists, it is skipped.
func AddPrefix(secrets map[string]string, opts PrefixOptions) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		newKey := opts.Prefix + k
		if _, exists := out[newKey]; exists && !opts.Overwrite {
			continue
		}
		out[newKey] = v
	}
	return out
}

// StripPrefix returns a new map with the given prefix removed from matching keys.
// Keys that do not carry the prefix are passed through unchanged unless they
// would collide with a stripped key, in which case the stripped key wins.
func StripPrefix(secrets map[string]string, prefix string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if strings.HasPrefix(k, prefix) {
			out[strings.TrimPrefix(k, prefix)] = v
		} else {
			if _, exists := out[k]; !exists {
				out[k] = v
			}
		}
	}
	return out
}
