package env

import (
	"fmt"
	"strconv"
	"strings"
)

// TypeRule defines an expected type for an environment variable key.
type TypeRule struct {
	Key      string
	Expected string // "string", "int", "float", "bool", "url", "email"
}

// TypeViolation describes a key whose value does not match its expected type.
type TypeViolation struct {
	Key      string
	Expected string
	Actual   string
}

func (v TypeViolation) Error() string {
	return fmt.Sprintf("key %q: expected %s, got %q", v.Key, v.Expected, v.Actual)
}

// CheckTypes validates that each key in env matches its declared type rule.
// Keys not listed in rules are silently skipped.
func CheckTypes(env map[string]string, rules []TypeRule) []TypeViolation {
	var violations []TypeViolation
	for _, r := range rules {
		val, ok := env[r.Key]
		if !ok {
			continue
		}
		if !matchesType(val, r.Expected) {
			violations = append(violations, TypeViolation{
				Key:      r.Key,
				Expected: r.Expected,
				Actual:   val,
			})
		}
	}
	return violations
}

// ParseTypeRules parses rules from a slice of "KEY=type" strings.
func ParseTypeRules(raw []string) ([]TypeRule, error) {
	rules := make([]TypeRule, 0, len(raw))
	for _, s := range raw {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid type rule %q: expected KEY=type", s)
		}
		rules = append(rules, TypeRule{Key: parts[0], Expected: parts[1]})
	}
	return rules, nil
}

func matchesType(val, typ string) bool {
	switch strings.ToLower(typ) {
	case "string":
		return true
	case "int":
		_, err := strconv.Atoi(val)
		return err == nil
	case "float":
		_, err := strconv.ParseFloat(val, 64)
		return err == nil
	case "bool":
		_, err := strconv.ParseBool(val)
		return err == nil
	case "url":
		return strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://")
	case "email":
		return strings.Contains(val, "@") && strings.Contains(val, ".")
	default:
		return false
	}
}
