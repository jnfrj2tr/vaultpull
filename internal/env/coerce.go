package env

import (
	"fmt"
	"strconv"
	"strings"
)

// CoerceType describes how a value should be coerced.
type CoerceType string

const (
	CoerceString CoerceType = "string"
	CoerceBool   CoerceType = "bool"
	CoerceInt    CoerceType = "int"
	CoerceFloat  CoerceType = "float"
	CoerceUpper  CoerceType = "upper"
	CoerceLower  CoerceType = "lower"
)

// CoerceRule maps a key pattern (exact match) to a target type.
type CoerceRule struct {
	Key  string
	Type CoerceType
}

// Coerce applies type coercion rules to the given env map.
// Values are converted to their canonical string representation.
// Returns an error if a value cannot be coerced to the requested type.
func Coerce(env map[string]string, rules []CoerceRule) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	for _, rule := range rules {
		v, ok := out[rule.Key]
		if !ok {
			continue
		}
		coerced, err := coerceValue(v, rule.Type)
		if err != nil {
			return nil, fmt.Errorf("coerce %q to %s: %w", rule.Key, rule.Type, err)
		}
		out[rule.Key] = coerced
	}
	return out, nil
}

// ParseCoerceRules parses rules from strings of the form "KEY:type".
func ParseCoerceRules(raw []string) ([]CoerceRule, error) {
	rules := make([]CoerceRule, 0, len(raw))
	for _, s := range raw {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid coerce rule %q: expected KEY:type", s)
		}
		rules = append(rules, CoerceRule{Key: parts[0], Type: CoerceType(parts[1])})
	}
	return rules, nil
}

func coerceValue(v string, t CoerceType) (string, error) {
	switch t {
	case CoerceString:
		return v, nil
	case CoerceBool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return "", fmt.Errorf("cannot parse %q as bool", v)
		}
		if b {
			return "true", nil
		}
		return "false", nil
	case CoerceInt:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return "", fmt.Errorf("cannot parse %q as int", v)
		}
		return strconv.FormatInt(int64(f), 10), nil
	case CoerceFloat:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return "", fmt.Errorf("cannot parse %q as float", v)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case CoerceUpper:
		return strings.ToUpper(v), nil
	case CoerceLower:
		return strings.ToLower(v), nil
	default:
		return "", fmt.Errorf("unknown coerce type %q", t)
	}
}
