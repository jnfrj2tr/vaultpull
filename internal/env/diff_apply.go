package env

import (
	"fmt"

	"github.com/yourusername/vaultpull/internal/diff"
)

// ApplyOptions controls how a diff is applied to an env map.
type ApplyOptions struct {
	// SkipUnchanged omits keys with no change from the result.
	SkipUnchanged bool
	// DryRun returns the would-be result without modifying the base.
	DryRun bool
	// OnConflict is called when an existing key is about to be overwritten.
	// Return an error to abort the apply.
	OnConflict func(key, oldVal, newVal string) error
}

// ApplyDiff merges diff changes back into base, returning the updated map.
// Added and Updated entries from changes are written into base.
// Removed entries are deleted from base.
func ApplyDiff(base map[string]string, changes []diff.Change, opts ApplyOptions) (map[string]string, error) {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}

	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			result[c.Key] = c.Incoming

		case diff.Updated:
			if opts.OnConflict != nil {
				if err := opts.OnConflict(c.Key, c.Existing, c.Incoming); err != nil {
					return nil, fmt.Errorf("conflict on key %q: %w", c.Key, err)
				}
			}
			result[c.Key] = c.Incoming

		case diff.Removed:
			delete(result, c.Key)

		case diff.Unchanged:
			if opts.SkipUnchanged {
				delete(result, c.Key)
			}
		}
	}

	if opts.DryRun {
		return result, nil
	}

	return result, nil
}

// SummarizeDiff returns a human-readable one-line summary of the changes.
func SummarizeDiff(changes []diff.Change) string {
	var added, updated, removed, unchanged int
	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			added++
		case diff.Updated:
			updated++
		case diff.Removed:
			removed++
		case diff.Unchanged:
			unchanged++
		}
	}
	return fmt.Sprintf("+%d added, ~%d updated, -%d removed, =%d unchanged",
		added, updated, removed, unchanged)
}
