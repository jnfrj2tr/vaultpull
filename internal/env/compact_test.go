package env

import (
	"testing"
)

func TestCompact_RemoveEmpty(t *testing.T) {
	input := map[string]string{
		"KEY_A": "hello",
		"KEY_B": "",
		"KEY_C": "world",
	}
	out, err := Compact(input, CompactOptions{RemoveEmpty: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["KEY_B"]; ok {
		t.Error("expected KEY_B to be removed")
	}
	if out["KEY_A"] != "hello" || out["KEY_C"] != "world" {
		t.Error("non-empty keys should be preserved")
	}
}

func TestCompact_TrimSpace(t *testing.T) {
	input := map[string]string{
		"KEY_A": "  trimmed  ",
		"KEY_B": "clean",
	}
	out, err := Compact(input, CompactOptions{TrimSpace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY_A"] != "trimmed" {
		t.Errorf("expected 'trimmed', got %q", out["KEY_A"])
	}
	if out["KEY_B"] != "clean" {
		t.Errorf("expected 'clean', got %q", out["KEY_B"])
	}
}

func TestCompact_RemoveDuplicateValues(t *testing.T) {
	input := map[string]string{
		"KEY_A": "same",
		"KEY_B": "same",
		"KEY_C": "different",
	}
	out, err := Compact(input, CompactOptions{RemoveDuplicateValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Stable sort means KEY_A comes before KEY_B; KEY_A should survive.
	if _, ok := out["KEY_A"]; !ok {
		t.Error("expected KEY_A to be kept (first occurrence)")
	}
	if _, ok := out["KEY_B"]; ok {
		t.Error("expected KEY_B to be removed (duplicate value)")
	}
	if out["KEY_C"] != "different" {
		t.Error("expected KEY_C to be preserved")
	}
}

func TestCompact_NoOptions_PassThrough(t *testing.T) {
	input := map[string]string{
		"A": "",
		"B": "  spaced  ",
		"C": "val",
	}
	out, err := Compact(input, CompactOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(input) {
		t.Errorf("expected %d keys, got %d", len(input), len(out))
	}
}

func TestCompact_CombinedOptions(t *testing.T) {
	input := map[string]string{
		"A": "  ",
		"B": "val",
		"C": "val",
	}
	out, err := Compact(input, CompactOptions{TrimSpace: true, RemoveEmpty: true, RemoveDuplicateValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "A" trims to "" and is removed; "B" survives; "C" is duplicate of "B" and removed.
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d: %v", len(out), out)
	}
	if out["B"] != "val" {
		t.Errorf("expected B=val, got %q", out["B"])
	}
}
