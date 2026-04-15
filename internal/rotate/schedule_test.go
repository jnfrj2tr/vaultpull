package rotate_test

import (
	"testing"
	"time"

	"github.com/vaultpull/internal/rotate"
)

func TestIsDue_ZeroLastRotated(t *testing.T) {
	s := &rotate.Schedule{Interval: time.Hour, LastRotated: time.Time{}}
	if !s.IsDue() {
		t.Error("expected IsDue=true when LastRotated is zero")
	}
}

func TestIsDue_NotYetDue(t *testing.T) {
	s := &rotate.Schedule{
		Interval:    24 * time.Hour,
		LastRotated: time.Now().Add(-1 * time.Hour),
	}
	if s.IsDue() {
		t.Error("expected IsDue=false when interval not elapsed")
	}
}

func TestIsDue_Overdue(t *testing.T) {
	s := &rotate.Schedule{
		Interval:    time.Hour,
		LastRotated: time.Now().Add(-2 * time.Hour),
	}
	if !s.IsDue() {
		t.Error("expected IsDue=true when interval has elapsed")
	}
}

func TestIsDue_ZeroInterval(t *testing.T) {
	s := &rotate.Schedule{Interval: 0, LastRotated: time.Now().Add(-48 * time.Hour)}
	if s.IsDue() {
		t.Error("expected IsDue=false when interval is zero")
	}
}

func TestParseInterval_Valid(t *testing.T) {
	d, err := rotate.ParseInterval("12h")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 12*time.Hour {
		t.Errorf("got %v, want 12h", d)
	}
}

func TestParseInterval_Empty(t *testing.T) {
	d, err := rotate.ParseInterval("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 0 {
		t.Errorf("expected 0 duration, got %v", d)
	}
}

func TestParseInterval_Invalid(t *testing.T) {
	_, err := rotate.ParseInterval("notaduration")
	if err == nil {
		t.Error("expected error for invalid interval")
	}
}

func TestParseInterval_Negative(t *testing.T) {
	_, err := rotate.ParseInterval("-5m")
	if err == nil {
		t.Error("expected error for negative interval")
	}
}
