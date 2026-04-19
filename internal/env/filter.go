package env

import (
	"regexp"
	"strings"
)

// FilterOptions controls which keys are included or excluded.
type FilterOptions struct {
	IncludePrefix string
	ExcludePrefix string
	Pattern       string
}

// Filter returns a new map containing only keys matching the given options.
func Filter(secrets map[string]string, opts FilterOptions) (map[string]string, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
	}

	result := make(map[string]string)
	for k, v := range secrets {
		if opts.IncludePrefix != "" && !strings.HasPrefix(k, opts.IncludePrefix) {
			continue
		}
		if opts.ExcludePrefix != "" && strings.HasPrefix(k, opts.ExcludePrefix) {
			continue
		}
		if re != nil && !re.MatchString(k) {
			continue
		}
		result[k] = v
	}
	return result, nil
}
