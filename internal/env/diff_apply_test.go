package env

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultpull/internal/diff"
)

func baseMap() map[string]string {
	return map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
		"NAME": "mydb",
	}
}

func TestApplyDiff_AddedKey(t *testing.T) {
	changes := []diff.Change{
		{Key: "USER", Type: diff.Added, Incoming: "admin"},
	}
	result, err := ApplyDiff(baseMap(), changes, ApplyOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["USER"] != "admin" {
		t.Errorf("expected USER=admin, got %q", result["USER"])
	}
	if result["HOST"] != "localhost" {
		t.Errorf("expected HOST to be preserved")
	}
}

func TestApplyDiff_UpdatedKey(t *testing.T) {
	changes := []diff.Change{
		{Key: "PORT", Type: diff.Updated, Existing: "5432", Incoming: "5433"},
	}
	result, err := ApplyDiff(baseMap(), changes, ApplyOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["PORT"] != "5433" {
		t.Errorf("expected PORT=5433, got %q", result["PORT"])
	}
}

func TestApplyDiff_RemovedKey(t *testing.T) {
	changes := []diff.Change{
		{Key: "NAME", Type: diff.Removed, Existing: "mydb"},
	}
	result, err := ApplyDiff(baseMap(), changes, ApplyOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["NAME"]; ok {
		t.Errorf("expected NAME to be removed")
	}
}

func TestApplyDiff_OnConflictAborts(t *testing.T) {
	changes := []diff.Change{
		{Key: "PORT", Type: diff.Updated, Existing: "5432", Incoming: "9999"},
	}
	opts := ApplyOptions{
		OnConflict: func(key, old, new string) error {
			return errors.New("conflict not allowed")
		},
	}
	_, err := ApplyDiff(baseMap(), changes, opts)
	if err == nil {
		t.Fatal("expected error from OnConflict, got nil")
	}
}

func TestApplyDiff_SkipUnchanged(t *testing.T) {
	changes := []diff.Change{
		{Key: "HOST", Type: diff.Unchanged, Existing: "localhost", Incoming: "localhost"},
	}
	result, err := ApplyDiff(baseMap(), changes, ApplyOptions{SkipUnchanged: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["HOST"]; ok {
		t.Errorf("expected HOST to be removed when SkipUnchanged=true")
	}
}

func TestSummarizeDiff(t *testing.T) {
	changes := []diff.Change{
		{Key: "A", Type: diff.Added},
		{Key: "B", Type: diff.Added},
		{Key: "C", Type: diff.Updated},
		{Key: "D", Type: diff.Removed},
		{Key: "E", Type: diff.Unchanged},
	}
	summary := SummarizeDiff(changes)
	expected := "+2 added, ~1 updated, -1 removed, =1 unchanged"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}
