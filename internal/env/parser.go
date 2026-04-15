package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair parsed from a .env file.
type Entry struct {
	Key   string
	Value string
}

// Parse reads a .env file and returns a slice of Entries.
// Lines starting with '#' are treated as comments and skipped.
// Empty lines are skipped. Lines without '=' are treated as errors.
func Parse(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: open %q: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("parser: %q line %d: missing '=' separator", path, lineNum)
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = stripQuotes(val)

		if key == "" {
			return nil, fmt.Errorf("parser: %q line %d: empty key", path, lineNum)
		}

		entries = append(entries, Entry{Key: key, Value: val})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: scanning %q: %w", path, err)
	}

	return entries, nil
}

// ToMap converts a slice of Entries into a map of key-value pairs.
// If duplicate keys are present, the last value wins.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// stripQuotes removes surrounding single or double quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
