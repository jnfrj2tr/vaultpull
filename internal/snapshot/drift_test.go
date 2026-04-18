package snapshot

import "testing"

func baseSnap() *Snapshot {
	return &Snapshot{
		Secrets: map[string]string{
			"FOO": "bar",
			"BAZ": "qux",
		},
	}
}

func TestDetectDrift_NilSnapshot_ReturnsNil(t *testing.T) {
	results := DetectDrift(nil, map[string]string{"FOO": "bar"})
	if results != nil {
		t.Errorf("expected nil, got %v", results)
	}
}

func TestDetectDrift_Clean(t *testing.T) {
	results := DetectDrift(baseSnap(), map[string]string{"FOO": "bar", "BAZ": "qux"})
	if len(results) != 0 {
		t.Errorf("expected no drift, got %v", results)
	}
}

func TestDetectDrift_Added(t *testing.T) {
	live := map[string]string{"FOO": "bar", "BAZ": "qux", "NEW": "val"}
	results := DetectDrift(baseSnap(), live)
	if len(results) != 1 || results[0].Key != "NEW" || results[0].Status != "added" {
		t.Errorf("expected added NEW, got %v", results)
	}
}

func TestDetectDrift_Removed(t *testing.T) {
	live := map[string]string{"FOO": "bar"}
	results := DetectDrift(baseSnap(), live)
	if len(results) != 1 || results[0].Key != "BAZ" || results[0].Status != "removed" {
		t.Errorf("expected removed BAZ, got %v", results)
	}
}

func TestDetectDrift_Changed(t *testing.T) {
	live := map[string]string{"FOO": "CHANGED", "BAZ": "qux"}
	results := DetectDrift(baseSnap(), live)
	if len(results) != 1 || results[0].Key != "FOO" || results[0].Status != "changed" {
		t.Errorf("expected changed FOO, got %v", results)
	}
}

func TestFormatDrift_NoResults(t *testing.T) {
	out := FormatDrift(nil)
	if out != "no drift detected" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatDrift_WithResults(t *testing.T) {
	results := []DriftResult{{Key: "FOO", Status: "changed"}}
	out := FormatDrift(results)
	if out == "" {
		t.Error("expected non-empty output")
	}
}
