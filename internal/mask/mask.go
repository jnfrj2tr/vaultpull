// Package mask provides utilities for redacting sensitive secret values
// before they are displayed in logs, diffs, or terminal output.
package mask

import "strings"

const placeholder = "****"

// Redact replaces all values in the provided map with a placeholder string.
func Redact(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k := range secrets {
		out[k] = placeholder
	}
	return out
}

// RedactValue masks a single secret value. If the value is empty it is
// returned as-is so that empty-string distinctions are preserved.
func RedactValue(v string) string {
	if v == "" {
		return v
	}
	return placeholder
}

// PartialRedact shows the first n characters of v followed by asterisks.
// If n >= len(v) the full value is masked.
func PartialRedact(v string, n int) string {
	if v == "" {
		return v
	}
	runes := []rune(v)
	if n <= 0 || n >= len(runes) {
		return placeholder
	}
	return string(runes[:n]) + strings.Repeat("*", len(runes)-n)
}
