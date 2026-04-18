// Package validate checks secrets against simple rules before writing.
package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for a secret key.
type Rule struct {
	Key      string
	Required bool
	Pattern  string // optional regex
}

// Violation describes a failed rule.
type Violation struct {
	Key     string
	Reason  string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Reason)
}

// Validate checks secrets against the provided rules.
// It returns a slice of violations (empty means all passed).
func Validate(secrets map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for _, r := range rules {
		val, exists := secrets[r.Key]

		if r.Required && (!exists || strings.TrimSpace(val) == "") {
			violations = append(violations, Violation{
				Key:    r.Key,
				Reason: "required but missing or empty",
			})
			continue
		}

		if r.Pattern != "" && exists && val != "" {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				violations = append(violations, Violation{
					Key:    r.Key,
					Reason: fmt.Sprintf("invalid pattern %q: %v", r.Pattern, err),
				})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, Violation{
					Key:    r.Key,
					Reason: fmt.Sprintf("value does not match pattern %q", r.Pattern),
				})
			}
		}
	}

	return violations
}
