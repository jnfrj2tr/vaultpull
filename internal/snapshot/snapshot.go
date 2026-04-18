package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot holds a point-in-time copy of secrets for a profile.
type Snapshot struct {
	Profile   string            `json:"profile"`
	CapturedAt time.Time        `json:"captured_at"`
	Secrets   map[string]string `json:"secrets"`
}

// Save writes a snapshot to disk under dir/<profile>.snap.json.
func Save(dir, profile string, secrets map[string]string) error {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("snapshot: mkdir: %w", err)
	}
	snap := Snapshot{
		Profile:    profile,
		CapturedAt: time.Now().UTC(),
		Secrets:    secrets,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	path := filepath.Join(dir, profile+".snap.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write: %w", err)
	}
	return nil
}

// Load reads a snapshot from disk. Returns nil if the file does not exist.
func Load(dir, profile string) (*Snapshot, error) {
	path := filepath.Join(dir, profile+".snap.json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &snap, nil
}
