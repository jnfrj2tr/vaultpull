package audit

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRecord_WritesEntry(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	logger := NewLogger(logPath)
	err := logger.Record(Entry{
		Profile:     "production",
		SecretPath:  "secret/myapp",
		KeysWritten: 5,
		OutputFile:  ".env",
		Status:      "success",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f, err := os.Open(logPath)
	if err != nil {
		t.Fatalf("could not open log file: %v", err)
	}
	defer f.Close()

	var entry Entry
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("expected at least one log line")
	}
	if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
		t.Fatalf("invalid JSON in log: %v", err)
	}

	if entry.Profile != "production" {
		t.Errorf("expected profile 'production', got %q", entry.Profile)
	}
	if entry.KeysWritten != 5 {
		t.Errorf("expected 5 keys written, got %d", entry.KeysWritten)
	}
	if entry.Status != "success" {
		t.Errorf("expected status 'success', got %q", entry.Status)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestRecord_AppendsMultipleEntries(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")
	logger := NewLogger(logPath)

	for i := 0; i < 3; i++ {
		if err := logger.Record(Entry{Status: "success", Profile: "dev"}); err != nil {
			t.Fatalf("record %d failed: %v", i, err)
		}
	}

	f, _ := os.Open(logPath)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 log lines, got %d", count)
	}
}

func TestRecord_InvalidPath(t *testing.T) {
	logger := NewLogger("/nonexistent/dir/audit.log")
	err := logger.Record(Entry{Status: "success"})
	if err == nil {
		t.Error("expected error for invalid log path")
	}
}
