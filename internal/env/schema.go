package env

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaRule defines the expected shape of a single env variable.
type SchemaRule struct {
	Key      string
	Required bool
	Pattern  string // optional regex
	Allowed  []string // optional allowlist
}

// SchemaResult holds the outcome of a single rule check.
type SchemaResult struct {
	Key     string
	Passed  bool
	Reason  string
}

// ValidateSchema checks env map against a slice of SchemaRules.
// It returns one result per rule.
func ValidateSchema(env map[string]string, rules []SchemaRule) ([]SchemaResult, error) {
	var results []SchemaResult

	for _, rule := range rules {
		val, exists := env[rule.Key]

		if !exists || val == "" {
			if rule.Required {
				results = append(results, SchemaResult{
					Key:    rule.Key,
					Passed: false,
					Reason: "required key is missing or empty",
				})
			}
			continue
		}

		if rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				return nil, fmt.Errorf("invalid pattern for key %q: %w", rule.Key, err)
			}
			if !re.MatchString(val) {
				results = append(results, SchemaResult{
					Key:    rule.Key,
					Passed: false,
					Reason: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern),
				})
				continue
			}
		}

		if len(rule.Allowed) > 0 {
			matched := false
			for _, a := range rule.Allowed {
				if strings.EqualFold(val, a) {
					matched = true
					break
				}
			}
			if !matched {
				results = append(results, SchemaResult{
					Key:    rule.Key,
					Passed: false,
					Reason: fmt.Sprintf("value %q not in allowed list %v", val, rule.Allowed),
				})
				continue
			}
		}

		results = append(results, SchemaResult{Key: rule.Key, Passed: true})
	}

	return results, nil
}

// HasSchemaErrors returns true if any result failed.
func HasSchemaErrors(results []SchemaResult) bool {
	for _, r := range results {
		if !r.Passed {
			return true
		}
	}
	return false
}
