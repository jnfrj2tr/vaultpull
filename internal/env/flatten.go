package env

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenOptions controls how nested maps are flattened into env-style keys.
type FlattenOptions struct {
	Separator string // default "_"
	Uppercase bool
	Prefix    string
}

// Flatten converts a nested map[string]any into a flat map[string]string
// using the given separator between key segments.
func Flatten(nested map[string]any, opts FlattenOptions) (map[string]string, error) {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	result := make(map[string]string)
	if err := flattenRecurse(nested, opts.Prefix, opts.Separator, opts.Uppercase, result); err != nil {
		return nil, err
	}
	return result, nil
}

func flattenRecurse(m map[string]any, prefix, sep string, upper bool, out map[string]string) error {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := m[k]
		full := k
		if prefix != "" {
			full = prefix + sep + k
		}
		if upper {
			full = strings.ToUpper(full)
		}
		switch val := v.(type) {
		case map[string]any:
			if err := flattenRecurse(val, full, sep, upper, out); err != nil {
				return err
			}
		case string:
			out[full] = val
		case nil:
			out[full] = ""
		default:
			out[full] = fmt.Sprintf("%v", val)
		}
	}
	return nil
}
