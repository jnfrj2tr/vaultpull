// Package cache provides local caching of Vault secrets to reduce API calls.
package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a cached secret snapshot.
type Entry struct {
	Profile   string            `json:"profile"`
	Secrets   map[string]string `json:"secrets"`
	FetchedAt time.Time         `json:"fetched_at"`
	TTL       time.Duration     `json:"ttl"`
}

// IsExpired returns true if the cache entry is older than its TTL.
func (e *Entry) IsExpired() bool {
	if e.TTL == 0 {
		return true
	}
	return time.Since(e.FetchedAt) > e.TTL
}

// Store writes a cache entry to disk at the given path.
func Store(path string, entry *Entry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(entry)
}

// Load reads a cache entry from disk. Returns nil, nil if file does not exist.
func Load(path string) (*Entry, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var entry Entry
	if err := json.NewDecoder(f).Decode(&entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// CachePath returns the default cache file path for a profile.
func CachePath(profile string) string {
	return filepath.Join(".vaultpull_cache", profile+".json")
}
