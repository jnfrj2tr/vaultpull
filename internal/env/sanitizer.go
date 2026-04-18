package env

import (
	"regexp"
	"strings"
)

// SanitizeResult holds the outcome of a sanitize operation.
type SanitizeResult struct {
	Key     string
	Original string
	Fixed   string
	Changed bool
}

var invalidKeyRe = regexp.MustCompile(`[^A-Z0-9_]`)
var leadingDigitRe = regexp.MustCompile(`^[0-9]`)

// SanitizeKey normalises an env var key to be POSIX-compliant:
// uppercase, only [A-Z0-9_], must not start with a digit.
func SanitizeKey(key string) SanitizeResult {
	original := key
	key = strings.ToUpper(key)
	key = strings.ReplaceAll(key, "-", "_")
	key = strings.ReplaceAll(key, ".", "_")
	key = strings.ReplaceAll(key, " ", "_")
	key = invalidKeyRe.ReplaceAllString(key, "_")
	if leadingDigitRe.MatchString(key) {
		key = "_" + key
	}
	return SanitizeResult{
		Key:      key,
		Original: original,
		Fixed:    key,
		Changed:  key != original,
	}
}

// SanitizeMap applies SanitizeKey to every key in the map.
// Collisions after sanitization are resolved by appending a counter suffix.
func SanitizeMap(in map[string]string) (map[string]string, []SanitizeResult) {
	out := make(map[string]string, len(in))
	var results []SanitizeResult
	counts := map[string]int{}

	for k, v := range in {
		res := SanitizeKey(k)
		final := res.Fixed
		if _, exists := out[final]; exists {
			counts[final]++
			final = final + "_" + strings.Repeat("I", counts[final])
			res.Fixed = final
			res.Changed = true
		}
		out[final] = v
		results = append(results, res)
	}
	return out, results
}
