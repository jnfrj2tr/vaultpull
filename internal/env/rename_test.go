package env

import (
	"testing"
)

func TestRename_BasicRule(t *testing.T) {
	in := map[string]string{"OLD_KEY": "value"}
	rules := []RenameRule{{From: "OLD_KEY", To: "NEW_KEY"}}
	out, err := Rename(in, rules, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %q", out["NEW_KEY"])
	}
	if _, exists := out["OLD_KEY"]; exists {
		t.Error("OLD_KEY should have been removed")
	}
}

func TestRename_MissingSourceIsSkipped(t *testing.T) {
	in := map[string]string{"KEEP": "v"}
	rules := []RenameRule{{From: "MISSING", To: "OTHER"}}
	out, err := Rename(in, rules, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := out["OTHER"]; exists {
		t.Error("OTHER should not exist")
	}
}

func TestRename_ConflictWithoutOverwrite(t *testing.T) {
	in := map[string]string{"A": "1", "B": "2"}
	rules := []RenameRule{{From: "A", To: "B"}}
	_, err := Rename(in, rules, false)
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

func TestRename_ConflictWithOverwrite(t *testing.T) {
	in := map[string]string{"A": "new", "B": "old"}
	rules := []RenameRule{{From: "A", To: "B"}}
	out, err := Rename(in, rules, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["B"] != "new" {
		t.Errorf("expected B=new, got %q", out["B"])
	}
}

func TestRename_EmptyRuleReturnsError(t *testing.T) {
	in := map[string]string{"X": "1"}
	rules := []RenameRule{{From: "", To: "Y"}}
	_, err := Rename(in, rules, false)
	if err == nil {
		t.Fatal("expected error for empty From")
	}
}

func TestParseRenameRules_Valid(t *testing.T) {
	rules, err := ParseRenameRules([]string{"OLD=NEW", "FOO=BAR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].From != "OLD" || rules[0].To != "NEW" {
		t.Errorf("unexpected rule: %+v", rules[0])
	}
}

func TestParseRenameRules_Invalid(t *testing.T) {
	_, err := ParseRenameRules([]string{"NODIVIDER"})
	if err == nil {
		t.Fatal("expected parse error")
	}
}
