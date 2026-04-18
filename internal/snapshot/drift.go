package snapshot

// DriftResult holds the outcome of comparing a snapshot to current secrets.
type DriftResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string]string // key -> new value
	Clean   bool
}

// DetectDrift compares a saved snapshot against the current secrets map.
// Returns nil if snap is nil (no previous snapshot to compare).
func DetectDrift(snap *Snapshot, current map[string]string) *DriftResult {
	if snap == nil {
		return nil
	}
	result := &DriftResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string]string),
	}
	for k, v := range current {
		old, exists := snap.Secrets[k]
		if !exists {
			result.Added[k] = v
		} else if old != v {
			result.Changed[k] = v
		}
	}
	for k, v := range snap.Secrets {
		if _, exists := current[k]; !exists {
			result.Removed[k] = v
		}
	}
	result.Clean = len(result.Added) == 0 && len(result.Removed) == 0 && len(result.Changed) == 0
	return result
}
