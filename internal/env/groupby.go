package env

import (
	"fmt"
	"regexp"
	"strings"
)

// GroupByResult holds the grouped keys organized by group name.
type GroupByResult map[string]map[string]string

// GroupByMode controls how keys are grouped.
type GroupByMode string

const (
	GroupByPrefix    GroupByMode = "prefix"
	GroupByDelimiter GroupByMode = "delimiter"
	GroupByPattern   GroupByMode = "pattern"
)

// GroupByOptions configures grouping behaviour.
type GroupByOptions struct {
	Mode      GroupByMode
	Delimiter string
	Pattern   string
	Fallback  string // group name for unmatched keys
}

// GroupBy partitions env vars into named groups according to opts.
func GroupBy(env map[string]string, opts GroupByOptions) (GroupByResult, error) {
	if opts.Fallback == "" {
		opts.Fallback = "other"
	}

	result := make(GroupByResult)

	switch opts.Mode {
	case GroupByPrefix:
		delim := opts.Delimiter
		if delim == "" {
			delim = "_"
		}
		for k, v := range env {
			parts := strings.SplitN(k, delim, 2)
			group := opts.Fallback
			if len(parts) == 2 {
				group = parts[0]
			}
			if result[group] == nil {
				result[group] = make(map[string]string)
			}
			result[group][k] = v
		}

	case GroupByDelimiter:
		delim := opts.Delimiter
		if delim == "" {
			return nil, fmt.Errorf("groupby: delimiter mode requires a delimiter")
		}
		for k, v := range env {
			idx := strings.LastIndex(k, delim)
			group := opts.Fallback
			if idx > 0 {
				group = k[:idx]
			}
			if result[group] == nil {
				result[group] = make(map[string]string)
			}
			result[group][k] = v
		}

	case GroupByPattern:
		if opts.Pattern == "" {
			return nil, fmt.Errorf("groupby: pattern mode requires a pattern")
		}
		re, err := regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, fmt.Errorf("groupby: invalid pattern %q: %w", opts.Pattern, err)
		}
		for k, v := range env {
			group := opts.Fallback
			if m := re.FindStringSubmatch(k); len(m) >= 2 {
				group = m[1]
			}
			if result[group] == nil {
				result[group] = make(map[string]string)
			}
			result[group][k] = v
		}

	default:
		return nil, fmt.Errorf("groupby: unknown mode %q", opts.Mode)
	}

	return result, nil
}
