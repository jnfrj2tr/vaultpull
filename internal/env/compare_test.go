package env

import (
	"strings"
	"testing"
)

func TestCompare_OnlyInLeft(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1"}

	res := Compare(left, right)

	if _, ok := res.OnlyInLeft["B"]; !ok {
		t.Error("expected B to be only in left")
	}
	if len(res.OnlyInRight) != 0 {
		t.Errorf("expected no right-only keys, got %v", res.OnlyInRight)
	}
}

func TestCompare_OnlyInRight(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"A": "1", "C": "3"}

	res := Compare(left, right)

	if _, ok := res.OnlyInRight["C"]; !ok {
		t.Error("expected C to be only in right")
	}
	if len(res.OnlyInLeft) != 0 {
		t.Errorf("expected no left-only keys, got %v", res.OnlyInLeft)
	}
}

func TestCompare_Different(t *testing.T) {
	left := map[string]string{"KEY": "old"}
	right := map[string]string{"KEY": "new"}

	res := Compare(left, right)

	pair, ok := res.Different["KEY"]
	if !ok {
		t.Fatal("expected KEY to appear in Different")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("unexpected pair: %v", pair)
	}
}

func TestCompare_Identical(t *testing.T) {
	left := map[string]string{"X": "same"}
	right := map[string]string{"X": "same"}

	res := Compare(left, right)

	if v, ok := res.Identical["X"]; !ok || v != "same" {
		t.Errorf("expected X to be identical, got %v", res.Identical)
	}
	if res.HasDifferences() {
		t.Error("expected no differences")
	}
}

func TestCompare_HasDifferences(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"A": "2"}

	res := Compare(left, right)
	if !res.HasDifferences() {
		t.Error("expected HasDifferences to return true")
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	res := Compare(map[string]string{}, map[string]string{})
	if res.HasDifferences() {
		t.Error("expected no differences for empty maps")
	}
}

func TestFormat_ContainsExpectedLines(t *testing.T) {
	left := map[string]string{"OLD": "x", "SAME": "y"}
	right := map[string]string{"NEW": "z", "SAME": "y", "OLD": "changed"}

	res := Compare(left, right)
	out := res.Format()

	if !strings.Contains(out, "+ NEW") {
		t.Errorf("expected '+ NEW' in output:\n%s", out)
	}
	if !strings.Contains(out, "~ OLD") {
		t.Errorf("expected '~ OLD' in output:\n%s", out)
	}
	if strings.Contains(out, "SAME") {
		t.Errorf("identical key SAME should not appear in format output:\n%s", out)
	}
}
