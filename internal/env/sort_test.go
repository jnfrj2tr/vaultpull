package env

import (
	"testing"
)

func TestSort_Alpha(t *testing.T) {
	m := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	keys, err := Sort(m, SortAlpha)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, k, expected[i])
		}
	}
}

func TestSort_AlphaDesc(t *testing.T) {
	m := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	keys, err := Sort(m, SortAlphaDesc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keys[0] != "ZEBRA" {
		t.Errorf("expected ZEBRA first, got %q", keys[0])
	}
}

func TestSort_ByLength(t *testing.T) {
	m := map[string]string{"AB": "1", "ABCDE": "2", "ABC": "3"}
	keys, err := Sort(m, SortByLength)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keys[0] != "AB" {
		t.Errorf("expected shortest key first, got %q", keys[0])
	}
	if keys[2] != "ABCDE" {
		t.Errorf("expected longest key last, got %q", keys[2])
	}
}

func TestSort_Natural(t *testing.T) {
	m := map[string]string{"zebra": "1", "Apple": "2", "mango": "3"}
	keys, err := Sort(m, SortNatural)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keys[0] != "Apple" {
		t.Errorf("expected Apple first (case-insensitive), got %q", keys[0])
	}
}

func TestSort_UnknownOrder(t *testing.T) {
	m := map[string]string{"A": "1"}
	_, err := Sort(m, SortOrder("bogus"))
	if err == nil {
		t.Error("expected error for unknown sort order")
	}
}

func TestSort_EmptyMap(t *testing.T) {
	keys, err := Sort(map[string]string{}, SortAlpha)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}
