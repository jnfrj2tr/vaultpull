package diff_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/diff"
)

func TestCompare_AddedKeys(t *testing.T) {
	existing := map[string]string{}
	incoming := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}

	result := diff.Compare(existing, incoming)

	if len(result.Added()) != 2 {
		t.Fatalf("expected 2 added, got %d", len(result.Added()))
	}
	if !result.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestCompare_UpdatedKeys(t *testing.T) {
	existing := map[string]string{"DB_HOST": "old-host"}
	incoming := map[string]string{"DB_HOST": "new-host"}

	result := diff.Compare(existing, incoming)

	updated := result.Updated()
	if len(updated) != 1 {
		t.Fatalf("expected 1 updated, got %d", len(updated))
	}
	if updated[0].OldValue != "old-host" {
		t.Errorf("expected OldValue 'old-host', got %q", updated[0].OldValue)
	}
	if updated[0].NewValue != "new-host" {
		t.Errorf("expected NewValue 'new-host', got %q", updated[0].NewValue)
	}
}

func TestCompare_UnchangedKeys(t *testing.T) {
	existing := map[string]string{"API_KEY": "secret"}
	incoming := map[string]string{"API_KEY": "secret"}

	result := diff.Compare(existing, incoming)

	if result.HasChanges() {
		t.Error("expected no changes")
	}
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change entry, got %d", len(result.Changes))
	}
	if result.Changes[0].Kind != diff.Unchanged {
		t.Errorf("expected Unchanged, got %s", result.Changes[0].Kind)
	}
}

func TestCompare_MixedChanges(t *testing.T) {
	existing := map[string]string{
		"KEEP":   "same",
		"CHANGE": "old",
	}
	incoming := map[string]string{
		"KEEP":    "same",
		"CHANGE":  "new",
		"NEW_KEY": "value",
	}

	result := diff.Compare(existing, incoming)

	if len(result.Added()) != 1 {
		t.Errorf("expected 1 added, got %d", len(result.Added()))
	}
	if len(result.Updated()) != 1 {
		t.Errorf("expected 1 updated, got %d", len(result.Updated()))
	}
	if !result.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestCompare_EmptyIncoming(t *testing.T) {
	existing := map[string]string{"FOO": "bar"}
	incoming := map[string]string{}

	result := diff.Compare(existing, incoming)

	if len(result.Changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(result.Changes))
	}
	if result.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestCompare_BothEmpty(t *testing.T) {
	existing := map[string]string{}
	incoming := map[string]string{}

	result := diff.Compare(existing, incoming)

	if len(result.Changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(result.Changes))
	}
	if result.HasChanges() {
		t.Error("expected no changes when both maps are empty")
	}
}
