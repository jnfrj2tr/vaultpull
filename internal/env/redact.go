package env

import (
	"regexp"
	"strings"
)

// RedactOptions controls how values are redacted in an env map.
type RedactOptions struct {
	// Keys is an explicit list of key names to redact.
	Keys []string
	// Patterns is a list of regex patterns; matching keys are redacted.
	Patterns []string
	// PartialReveal shows the first N characters before masking the rest.
	PartialReveal int
}

// RedactMap returns a copy of env with sensitive values masked.
// Keys listed in opts.Keys or matching any pattern in opts.Patterns are
// replaced with "***" (or a partial reveal if opts.PartialReveal > 0).
func RedactMap(env map[string]string, opts RedactOptions) (map[string]string, error) {
	compiled := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[strings.ToUpper(k)] = struct{}{}
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if shouldRedact(k, keySet, compiled) {
			out[k] = maskValue(v, opts.PartialReveal)
		} else {
			out[k] = v
		}
	}
	return out, nil
}

func shouldRedact(key string, keySet map[string]struct{}, patterns []*regexp.Regexp) bool {
	if _, ok := keySet[strings.ToUpper(key)]; ok {
		return true
	}
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

func maskValue(v string, reveal int) string {
	if reveal <= 0 || reveal >= len(v) {
		return "***"
	}
	return v[:reveal] + "***"
}
