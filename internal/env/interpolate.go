package env

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var interpolateRe = regexp.MustCompile(`\$\{([^}]+)\}`)

// Interpolate replaces ${VAR} references in values using the provided secrets
// map first, then falling back to OS environment variables.
func Interpolate(secrets map[string]string, allowEnv bool) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		resolved, err := interpolateValue(v, secrets, allowEnv)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

func interpolateValue(val string, secrets map[string]string, allowEnv bool) (string, error) {
	var rerr error
	result := interpolateRe.ReplaceAllStringFunc(val, func(match string) string {
		if rerr != nil {
			return match
		}
		inner := match[2 : len(match)-1] // strip ${ and }
		parts := strings.SplitN(inner, ":-", 2)
		name := parts[0]
		defaultVal := ""
		hasDefault := len(parts) == 2
		if v, ok := secrets[name]; ok {
			return v
		}
		if allowEnv {
			if v, ok := os.LookupEnv(name); ok {
				return v
			}
		}
		if hasDefault {
			return defaultVal
		}
		rerr = fmt.Errorf("unresolved variable: %s", name)
		return match
	})
	if rerr != nil {
		return "", rerr
	}
	return result, nil
}
