package snapshot_test

import (
	"os"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/snapshot"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "snapshot-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := tempDir(t)
	secrets := map[string]string{"KEY": "value", "TOKEN": "abc123"}

	if err := snapshot.Save(dir, "prod", secrets); err != nil {
		t.Fatalf("Save: %v", err)
	}
	snap, err := snapshot.Load(dir, "prod")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if snap == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if snap.Profile != "prod" {
		t.Errorf("profile = %q, want prod", snap.Profile)
	}
	if snap.Secrets["KEY"] != "value" {
		t.Errorf("KEY = %q, want value", snap.Secrets["KEY"])
	}
	if snap.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
	_ = time.Now()
}

func TestLoad_NonExistent_ReturnsNil(t *testing.T) {
	dir := tempDir(t)
	snap, err := snapshot.Load(dir, "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap != nil {
		t.Error("expected nil snapshot for missing profile")
	}
}

func TestSave_InvalidDir_ReturnsError(t *testing.T) {
	err := snapshot.Save("/proc/invalid-vaultpull-test", "prod", map[string]string{})
	if err == nil {
		t.Error("expected error for invalid dir")
	}
}
