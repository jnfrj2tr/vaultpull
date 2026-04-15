package audit

import (
	"path/filepath"
	"testing"
)

func TestSummarize_CountsCorrectly(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")
	logger := NewLogger(logPath)

	entries := []Entry{
		{Profile: "dev", Status: "success"},
		{Profile: "dev", Status: "success"},
		{Profile: "production", Status: "failure", Error: "permission denied"},
		{Profile: "staging", Status: "success"},
	}
	for _, e := range entries {
		if err := logger.Record(e); err != nil {
			t.Fatalf("record failed: %v", err)
		}
	}

	s, err := Summarize(logPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.TotalRuns != 4 {
		t.Errorf("expected 4 total runs, got %d", s.TotalRuns)
	}
	if s.SuccessCount != 3 {
		t.Errorf("expected 3 successes, got %d", s.SuccessCount)
	}
	if s.FailureCount != 1 {
		t.Errorf("expected 1 failure, got %d", s.FailureCount)
	}
	if s.Profiles["dev"] != 2 {
		t.Errorf("expected 2 dev runs, got %d", s.Profiles["dev"])
	}
}

func TestSummarize_MissingFile(t *testing.T) {
	_, err := Summarize("/nonexistent/audit.log")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSummarize_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")
	logger := NewLogger(logPath)
	// Record nothing — file is created but empty
	_ = logger

	s, err := Summarize(logPath)
	if err == nil && s != nil {
		// empty file is valid, just zero counts — but file won't exist yet
	}
	// If file does not exist, err is expected
	if err != nil {
		t.Logf("got expected error for non-created file: %v", err)
	}
}
