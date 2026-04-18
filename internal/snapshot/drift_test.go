package snapshot_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/snapshot"
)

func baseSnap(secrets map[string]string) *snapshot.Snapshot {
	return &snapshot.Snapshot{Profile: "test", Secrets: secrets}
}

func TestDetectDrift_NilSnapshot_ReturnsNil(t *testing.T) {
	result := snapshot.DetectDrift(nil, map[string]string{"A": "1"})
	if result != nil {
		t.Error("expected nil result for nil snapshot")
	}
}

func TestDetectDrift_Clean(t *testing.T) {
	snap := baseSnap(map[string]string{"A": "1", "B": "2"})
	result := snapshot.DetectDrift(snap, map[string]string{"A": "1", "B": "2"})
	if !result.Clean {
		t.Error("expected clean drift")
	}
}

func TestDetectDrift_Added(t *testing.T) {
	snap := baseSnap(map[string]string{"A": "1"})
	result := snapshot.DetectDrift(snap, map[string]string{"A": "1", "B": "new"})
	if result.Clean {
		t.Error("expected dirty drift")
	}
	if result.Added["B"] != "new" {
		t.Errorf("Added[B] = %q, want new", result.Added["B"])
	}
}

func TestDetectDrift_Removed(t *testing.T) {
	snap := baseSnap(map[string]string{"A": "1", "B": "2"})
	result := snapshot.DetectDrift(snap, map[string]string{"A": "1"})
	if _, ok := result.Removed["B"]; !ok {
		t.Error("expected B in Removed")
	}
}

func TestDetectDrift_Changed(t *testing.T) {
	snap := baseSnap(map[string]string{"A": "old"})
	result := snapshot.DetectDrift(snap, map[string]string{"A": "new"})
	if result.Changed["A"] != "new" {
		t.Errorf("Changed[A] = %q, want new", result.Changed["A"])
	}
}
