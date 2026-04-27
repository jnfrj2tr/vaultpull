package env

import (
	"fmt"
	"strings"
)

// AliasRule maps a source key to one or more alias keys.
type AliasRule struct {
	Source  string
	Aliases []string
	Keep    bool // retain the original key alongside aliases
}

// Alias duplicates values from source keys into alias keys.
// If Keep is false the source key is removed after aliasing.
func Alias(env map[string]string, rules []AliasRule) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	for _, rule := range rules {
		val, ok := out[rule.Source]
		if !ok {
			continue
		}
		for _, alias := range rule.Aliases {
			if alias == "" {
				return nil, fmt.Errorf("alias: empty alias name for source %q", rule.Source)
			}
			out[alias] = val
		}
		if !rule.Keep {
			delete(out, rule.Source)
		}
	}
	return out, nil
}

// ParseAliasRules parses rules from strings of the form "SRC:ALIAS1,ALIAS2[+keep]".
// Append "+keep" to retain the source key.
func ParseAliasRules(raw []string) ([]AliasRule, error) {
	rules := make([]AliasRule, 0, len(raw))
	for _, s := range raw {
		keep := false
		if strings.HasSuffix(s, "+keep") {
			keep = true
			s = strings.TrimSuffix(s, "+keep")
		}
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("alias: invalid rule %q, expected SRC:ALIAS1,ALIAS2", s)
		}
		aliases := strings.Split(parts[1], ",")
		for i, a := range aliases {
			aliases[i] = strings.TrimSpace(a)
		}
		rules = append(rules, AliasRule{
			Source:  strings.TrimSpace(parts[0]),
			Aliases: aliases,
			Keep:    keep,
		})
	}
	return rules, nil
}
