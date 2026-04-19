package env

import (
	"sort"
	"strings"
)

// SortOrder defines how keys should be sorted.
type SortOrder string

const (
	SortAlpha      SortOrder = "alpha"
	SortAlphaDesc  SortOrder = "alpha-desc"
	SortByLength   SortOrder = "length"
	SortNatural    SortOrder = "natural"
)

// Sort returns a new map with keys ordered according to the given order.
// The returned slice of key-value pairs preserves the requested ordering.
func Sort(m map[string]string, order SortOrder) ([]string, error) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	switch order {
	case SortAlpha, "":
		sort.Strings(keys)
	case SortAlphaDesc:
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	case SortByLength:
		sort.Slice(keys, func(i, j int) bool {
			if len(keys[i]) == len(keys[j]) {
				return keys[i] < keys[j]
			}
			return len(keys[i]) < len(keys[j])
		})
	case SortNatural:
		sort.Slice(keys, func(i, j int) bool {
			return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
		})
	default:
		return nil, fmt.Errorf("unknown sort order: %q", order)
	}

	return keys, nil
}
