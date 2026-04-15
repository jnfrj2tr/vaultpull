package audit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Summary holds aggregated statistics from an audit log file.
type Summary struct {
	TotalRuns    int
	SuccessCount int
	FailureCount int
	Profiles     map[string]int
}

// Summarize reads the audit log at path and returns aggregated stats.
func Summarize(path string) (*Summary, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("audit: open log for summary: %w", err)
	}
	defer f.Close()

	s := &Summary{
		Profiles: make(map[string]int),
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			continue // skip malformed lines
		}
		s.TotalRuns++
		if e.Status == "success" {
			s.SuccessCount++
		} else {
			s.FailureCount++
		}
		if e.Profile != "" {
			s.Profiles[e.Profile]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("audit: scan log: %w", err)
	}
	return s, nil
}
