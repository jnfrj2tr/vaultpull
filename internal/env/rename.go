package env

import (
	"fmt"
	"strings"
)

// RenameRule describes a single key rename operation.
type RenameRule struct {
	From string
	To   string
}

// Rename applies a set of rename rules to a map of env vars.
// Keys not matching any rule are passed through unchanged.
// Returns an error if a target key already exists and overwrite is false.
func Rename(secrets map[string]string, rules []RenameRule, overwrite bool) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, r := range rules {
		if r.From == "" || r.To == "" {
			return nil, fmt.Errorf("rename rule has empty from or to field")
		}
		val, exists := out[r.From]
		if !exists {
			continue
		}
		if _, conflict := out[r.To]; conflict && !overwrite {
			return nil, fmt.Errorf("rename target %q already exists", r.To)
		}
		out[r.To] = val
		delete(out, r.From)
	}
	return out, nil
}

// ParseRenameRules parses a slice of "FROM=TO" strings into RenameRule values.
func ParseRenameRules(raw []string) ([]RenameRule, error) {
	rules := make([]RenameRule, 0, len(raw))
	for _, s := range raw {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid rename rule %q: expected FROM=TO", s)
		}
		rules = append(rules, RenameRule{From: parts[0], To: parts[1]})
	}
	return rules, nil
}
