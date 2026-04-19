package env

import (
	"testing"
)

func TestParseStrategy_Valid(t *testing.T) {
	cases := []struct {
		input    string
		want     Strategy
	}{
		{"overwrite", StrategyOverwrite},
		{"keep", StrategyKeepExisting},
		{"keep-existing", StrategyKeepExisting},
		{"error", StrategyError},
	}
	for _, c := range cases {
		got, err := ParseStrategy(c.input)
		if err != nil {
			t.Fatalf("ParseStrategy(%q) unexpected error: %v", c.input, err)
		}
		if got != c.want {
			t.Errorf("ParseStrategy(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestParseStrategy_Invalid(t *testing.T) {
	_, err := ParseStrategy("unknown")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestMergeWithStrategy_Overwrite(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	incoming := map[string]string{"B": "99", "C": "3"}
	result, err := MergeWithStrategy(base, incoming, StrategyOverwrite)
	if err != nil {
		t.Fatal(err)
	}
	if result["B"] != "99" {
		t.Errorf("expected B=99, got %s", result["B"])
	}
	if result["C"] != "3" {
		t.Errorf("expected C=3, got %s", result["C"])
	}
}

func TestMergeWithStrategy_KeepExisting(t *testing.T) {
	base := map[string]string{"A": "original"}
	incoming := map[string]string{"A": "new", "B": "added"}
	result, err := MergeWithStrategy(base, incoming, StrategyKeepExisting)
	if err != nil {
		t.Fatal(err)
	}
	if result["A"] != "original" {
		t.Errorf("expected A=original, got %s", result["A"])
	}
	if result["B"] != "added" {
		t.Errorf("expected B=added, got %s", result["B"])
	}
}

func TestMergeWithStrategy_ErrorOnConflict(t *testing.T) {
	base := map[string]string{"A": "1"}
	incoming := map[string]string{"A": "2"}
	_, err := MergeWithStrategy(base, incoming, StrategyError)
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

func TestMergeWithStrategy_ErrorNoConflict(t *testing.T) {
	base := map[string]string{"A": "1"}
	incoming := map[string]string{"A": "1", "B": "2"}
	result, err := MergeWithStrategy(base, incoming, StrategyError)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["B"] != "2" {
		t.Errorf("expected B=2, got %s", result["B"])
	}
}

func TestMergeWithStrategy_NilBase(t *testing.T) {
	incoming := map[string]string{"X": "10"}
	result, err := MergeWithStrategy(nil, incoming, StrategyOverwrite)
	if err != nil {
		t.Fatal(err)
	}
	if result["X"] != "10" {
		t.Errorf("expected X=10, got %s", result["X"])
	}
}
