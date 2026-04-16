package cache_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/cache"
)

func tempCachePath(t *testing.T, name string) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, name+".json")
}

func TestStoreAndLoad_RoundTrip(t *testing.T) {
	path := tempCachePath(t, "dev")
	entry := &cache.Entry{
		Profile:   "dev",
		Secrets:   map[string]string{"KEY": "value"},
		FetchedAt: time.Now().Truncate(time.Second),
		TTL:       5 * time.Minute,
	}
	if err := cache.Store(path, entry); err != nil {
		t.Fatalf("Store: %v", err)
	}
	loaded, err := cache.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Profile != entry.Profile {
		t.Errorf("profile mismatch: got %q want %q", loaded.Profile, entry.Profile)
	}
	if loaded.Secrets["KEY"] != "value" {
		t.Errorf("secret mismatch")
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	entry, err := cache.Load("/tmp/does_not_exist_vaultpull.json")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil entry for missing file")
	}
}

func TestIsExpired_WithTTL(t *testing.T) {
	entry := &cache.Entry{
		FetchedAt: time.Now().Add(-10 * time.Minute),
		TTL:       5 * time.Minute,
	}
	if !entry.IsExpired() {
		t.Error("expected entry to be expired")
	}
}

func TestIsExpired_NotYetExpired(t *testing.T) {
	entry := &cache.Entry{
		FetchedAt: time.Now(),
		TTL:       10 * time.Minute,
	}
	if entry.IsExpired() {
		t.Error("expected entry to not be expired")
	}
}

func TestIsExpired_ZeroTTL(t *testing.T) {
	entry := &cache.Entry{FetchedAt: time.Now(), TTL: 0}
	if !entry.IsExpired() {
		t.Error("zero TTL should always be expired")
	}
}

func TestStore_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "dir", "dev.json")
	entry := &cache.Entry{Profile: "dev", Secrets: map[string]string{}, TTL: time.Minute, FetchedAt: time.Now()}
	if err := cache.Store(path, entry); err != nil {
		t.Fatalf("Store failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
