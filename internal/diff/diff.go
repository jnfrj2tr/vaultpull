// Package diff provides utilities for comparing existing .env file
// contents against secrets fetched from Vault, producing a structured
// change set that can be reported or acted upon.
package diff

// ChangeKind describes the type of change for a single key.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Updated  ChangeKind = "updated"
	Unchanged ChangeKind = "unchanged"
)

// Change represents a single key-level diff entry.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Result holds the full diff between existing and incoming secrets.
type Result struct {
	Changes []Change
}

// Added returns only the added changes.
func (r *Result) Added() []Change {
	return r.filter(Added)
}

// Updated returns only the updated changes.
func (r *Result) Updated() []Change {
	return r.filter(Updated)
}

// Unchanged returns only the unchanged changes.
func (r *Result) Unchanged() []Change {
	return r.filter(Unchanged)
}

// HasChanges reports whether any keys were added or updated.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Kind != Unchanged {
			return true
		}
	}
	return false
}

// Summary returns a count of added, updated, and unchanged keys.
func (r *Result) Summary() (added, updated, unchanged int) {
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			added++
		case Updated:
			updated++
		case Unchanged:
			unchanged++
		}
	}
	return
}

func (r *Result) filter(kind ChangeKind) []Change {
	var out []Change
	for _, c := range r.Changes {
		if c.Kind == kind {
			out = append(out, c)
		}
	}
	return out
}

// Compare produces a Result by comparing existing key/value pairs with
// the incoming secrets map. Keys present only in incoming are marked
// Added; keys present in both are marked Updated or Unchanged.
func Compare(existing map[string]string, incoming map[string]string) *Result {
	result := &Result{}

	for key, newVal := range incoming {
		oldVal, exists := existing[key]
		switch {
		case !exists:
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Kind:     Added,
				NewValue: newVal,
			})
		case oldVal != newVal:
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Kind:     Updated,
				OldValue: oldVal,
				NewValue: newVal,
			})
		default:
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Kind:     Unchanged,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	return result
}
