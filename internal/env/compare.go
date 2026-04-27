package env

import (
	"fmt"
	"sort"
	"strings"
)

// CompareResult holds the outcome of comparing two env maps.
type CompareResult struct {
	OnlyInLeft  map[string]string
	OnlyInRight map[string]string
	Different   map[string][2]string // key -> [left, right]
	Identical   map[string]string
}

// HasDifferences returns true if there are any keys that differ between the two maps.
func (r *CompareResult) HasDifferences() bool {
	return len(r.OnlyInLeft) > 0 || len(r.OnlyInRight) > 0 || len(r.Different) > 0
}

// Compare performs a key-by-key comparison of two env maps.
// left is considered the "base" (e.g. current .env), right is "incoming" (e.g. from Vault).
func Compare(left, right map[string]string) *CompareResult {
	res := &CompareResult{
		OnlyInLeft:  make(map[string]string),
		OnlyInRight: make(map[string]string),
		Different:   make(map[string][2]string),
		Identical:   make(map[string]string),
	}

	for k, lv := range left {
		if rv, ok := right[k]; ok {
			if lv == rv {
				res.Identical[k] = lv
			} else {
				res.Different[k] = [2]string{lv, rv}
			}
		} else {
			res.OnlyInLeft[k] = lv
		}
	}

	for k, rv := range right {
		if _, ok := left[k]; !ok {
			res.OnlyInRight[k] = rv
		}
	}

	return res
}

// Format returns a human-readable summary of the comparison result.
func (r *CompareResult) Format() string {
	var sb strings.Builder

	keys := func(m map[string]string) []string {
		out := make([]string, 0, len(m))
		for k := range m {
			out = append(out, k)
		}
		sort.Strings(out)
		return out
	}

	for _, k := range keys(r.OnlyInLeft) {
		sb.WriteString(fmt.Sprintf("- %s (only in left)\n", k))
	}
	for _, k := range keys(r.OnlyInRight) {
		sb.WriteString(fmt.Sprintf("+ %s (only in right)\n", k))
	}

	diffKeys := make([]string, 0, len(r.Different))
	for k := range r.Different {
		diffKeys = append(diffKeys, k)
	}
	sort.Strings(diffKeys)
	for _, k := range diffKeys {
		pair := r.Different[k]
		sb.WriteString(fmt.Sprintf("~ %s: %q -> %q\n", k, pair[0], pair[1]))
	}

	return sb.String()
}
