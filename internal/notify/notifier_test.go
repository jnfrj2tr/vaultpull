package notify

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func baseEvent() Event {
	return Event{
		Profile:   "staging",
		Operation: "sync",
		Added:     2,
		Updated:   1,
		Removed:   0,
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
}

func TestNotify_ContainsProfile(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	_ = n.Notify(baseEvent())
	if !strings.Contains(buf.String(), "staging") {
		t.Errorf("expected profile in output, got: %s", buf.String())
	}
}

func TestNotify_ContainsCounts(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	_ = n.Notify(baseEvent())
	out := buf.String()
	if !strings.Contains(out, "+2 added") {
		t.Errorf("expected added count, got: %s", out)
	}
	if !strings.Contains(out, "~1 updated") {
		t.Errorf("expected updated count, got: %s", out)
	}
}

func TestNotify_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	e := baseEvent()
	e.Added, e.Updated, e.Removed = 0, 0, 0
	_ = n.Notify(e)
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes', got: %s", buf.String())
	}
}

func TestNotify_DefaultsTimestamp(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	e := baseEvent()
	e.Timestamp = time.Time{}
	err := n.Notify(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	n := New(nil)
	if n.out == nil {
		t.Error("expected non-nil writer")
	}
}
