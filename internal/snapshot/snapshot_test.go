package snapshot

import (
	"os"
	"testing"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "snap-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := tempDir(t)
	secrets := map[string]string{"DB_PASS": "secret", "API_KEY": "abc123"}

	if err := Save(dir, "prod", secrets); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := Load(dir, "prod")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if snap == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if snap.Profile != "prod" {
		t.Errorf("profile = %q, want prod", snap.Profile)
	}
	if snap.Secrets["DB_PASS"] != "secret" {
		t.Errorf("DB_PASS = %q, want secret", snap.Secrets["DB_PASS"])
	}
	if snap.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestLoad_NonExistent_ReturnsNil(t *testing.T) {
	dir := tempDir(t)
	snap, err := Load(dir, "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap != nil {
		t.Error("expected nil snapshot for missing file")
	}
}

func TestSave_InvalidDir_ReturnsError(t *testing.T) {
	err := Save("/dev/null/no-such-dir", "prod", map[string]string{})
	if err == nil {
		t.Error("expected error for invalid dir")
	}
}
