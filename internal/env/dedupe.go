package env

// DedupeStrategy controls which value wins when duplicate keys are found.
type DedupeStrategy string

const (
	DedupeKeepFirst DedupeStrategy = "first"
	DedupeKeepLast  DedupeStrategy = "last"
)

// DedupeResult holds statistics about a deduplication pass.
type DedupeResult struct {
	Removed  int
	Kept     map[string]string
}

// Dedupe removes duplicate keys from an ordered list of key-value pairs.
// The strategy controls whether the first or last occurrence is retained.
// Input is a slice of [2]string{key, value} to preserve order information.
func Dedupe(pairs [][2]string, strategy DedupeStrategy) DedupeResult {
	seen := make(map[string]int) // key -> index in out
	out := make([][2]string, 0, len(pairs))
	removed := 0

	for _, pair := range pairs {
		k := pair[0]
		if idx, exists := seen[k]; exists {
			if strategy == DedupeKeepLast {
				out[idx] = pair
			}
			removed++
		} else {
			seen[k] = len(out)
			out = append(out, pair)
		}
	}

	kept := make(map[string]string, len(out))
	for _, pair := range out {
		kept[pair[0]] = pair[1]
	}

	return DedupeResult{Removed: removed, Kept: kept}
}

// DedupeMap removes duplicate keys from a flat map — no-op since maps
// cannot have duplicate keys, but merges two maps respecting strategy.
func DedupeMap(base, override map[string]string, strategy DedupeStrategy) map[string]string {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}
	if strategy == DedupeKeepLast {
		for k, v := range override {
			result[k] = v
		}
	} else {
		for k, v := range override {
			if _, exists := result[k]; !exists {
				result[k] = v
			}
		}
	}
	return result
}
