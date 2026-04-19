package env

import (
	"testing"
)

func TestDedupe_KeepFirst(t *testing.T) {
	pairs := [][2]string{
		{"FOO", "first"},
		{"BAR", "bar"},
		{"FOO", "second"},
	}
	res := Dedupe(pairs, DedupeKeepFirst)
	if res.Removed != 1 {
		t.Fatalf("expected 1 removed, got %d", res.Removed)
	}
	if res.Kept["FOO"] != "first" {
		t.Errorf("expected FOO=first, got %s", res.Kept["FOO"])
	}
}

func TestDedupe_KeepLast(t *testing.T) {
	pairs := [][2]string{
		{"FOO", "first"},
		{"BAR", "bar"},
		{"FOO", "second"},
	}
	res := Dedupe(pairs, DedupeKeepLast)
	if res.Removed != 1 {
		t.Fatalf("expected 1 removed, got %d", res.Removed)
	}
	if res.Kept["FOO"] != "second" {
		t.Errorf("expected FOO=second, got %s", res.Kept["FOO"])
	}
}

func TestDedupe_NoDuplicates(t *testing.T) {
	pairs := [][2]string{
		{"A", "1"},
		{"B", "2"},
	}
	res := Dedupe(pairs, DedupeKeepFirst)
	if res.Removed != 0 {
		t.Errorf("expected 0 removed, got %d", res.Removed)
	}
	if len(res.Kept) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Kept))
	}
}

func TestDedupeMap_KeepFirst(t *testing.T) {
	base := map[string]string{"FOO": "base", "BAR": "bar"}
	override := map[string]string{"FOO": "override", "BAZ": "baz"}
	result := DedupeMap(base, override, DedupeKeepFirst)
	if result["FOO"] != "base" {
		t.Errorf("expected FOO=base, got %s", result["FOO"])
	}
	if result["BAZ"] != "baz" {
		t.Errorf("expected BAZ=baz, got %s", result["BAZ"])
	}
}

func TestDedupeMap_KeepLast(t *testing.T) {
	base := map[string]string{"FOO": "base"}
	override := map[string]string{"FOO": "override"}
	result := DedupeMap(base, override, DedupeKeepLast)
	if result["FOO"] != "override" {
		t.Errorf("expected FOO=override, got %s", result["FOO"])
	}
}

func TestDedupeMap_EmptyOverride(t *testing.T) {
	base := map[string]string{"X": "1"}
	result := DedupeMap(base, nil, DedupeKeepLast)
	if result["X"] != "1" {
		t.Errorf("expected X=1, got %s", result["X"])
	}
}
