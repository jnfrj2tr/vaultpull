package snapshot

import "fmt"

// DriftResult describes a single detected drift between snapshot and live secrets.
type DriftResult struct {
	Key    string
	Status string // "added", "removed", "changed"
}

// DetectDrift compares a saved snapshot against the current live secrets map.
// Returns nil if snap is nil (no baseline to compare against).
func DetectDrift(snap *Snapshot, live map[string]string) []DriftResult {
	if snap == nil {
		return nil
	}

	var results []DriftResult

	// Keys removed or changed
	for k, snapVal := range snap.Secrets {
		liveVal, ok := live[k]
		if !ok {
			results = append(results, DriftResult{Key: k, Status: "removed"})
		} else if liveVal != snapVal {
			results = append(results, DriftResult{Key: k, Status: "changed"})
		}
	}

	// Keys added
	for k := range live {
		if _, ok := snap.Secrets[k]; !ok {
			results = append(results, DriftResult{Key: k, Status: "added"})
		}
	}

	return results
}

// FormatDrift returns a human-readable summary of drift results.
func FormatDrift(results []DriftResult) string {
	if len(results) == 0 {
		return "no drift detected"
	}
	var out string
	for _, r := range results {
		out += fmt.Sprintf("  [%s] %s\n", r.Status, r.Key)
	}
	return out
}
