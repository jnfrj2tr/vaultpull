package env

import "fmt"

// DefaultEntry represents a key with an optional default value and required flag.
type DefaultEntry struct {
	Key      string
	Default  string
	Required bool
}

// ApplyDefaults merges defaults into the given secrets map.
// If a key is missing and has a default, it is inserted.
// If a key is missing, has no default, and is required, an error is returned.
func ApplyDefaults(secrets map[string]string, defaults []DefaultEntry) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	for _, d := range defaults {
		if _, ok := out[d.Key]; !ok {
			if d.Default != "" {
				out[d.Key] = d.Default
			} else if d.Required {
				return nil, fmt.Errorf("required key %q is missing and has no default", d.Key)
			}
		}
	}

	return out, nil
}

// ParseDefaults converts a map[string]string (from config) into []DefaultEntry.
// Format of each value: "default=<val>" or "required" or "required,default=<val>".
func ParseDefaults(raw map[string]string) []DefaultEntry {
	entries := make([]DefaultEntry, 0, len(raw))
	for k, v := range raw {
		entry := DefaultEntry{Key: k}
		if v == "required" {
			entry.Required = true
		} else {
			entry.Default = v
		}
		entries = append(entries, entry)
	}
	return entries
}
