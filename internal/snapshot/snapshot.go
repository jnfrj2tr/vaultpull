// Package snapshot provides functionality to capture and compare
// secret state over time, enabling drift detection.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of secrets for a profile.
type Snapshot struct {
	Profile   string            `json:"profile"`
	CapturedAt time.Time        `json:"captured_at"`
	Secrets   map[string]string `json:"secrets"`
}

// Save writes a snapshot to the given directory as a JSON file.
func Save(dir, profile string, secrets map[string]string) error {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("snapshot: create dir: %w", err)
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
	path := filepath.Join(dir, profile+".json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write: %w", err)
	}
	return nil
}

// Load reads a previously saved snapshot for the given profile.
func Load(dir, profile string) (*Snapshot, error) {
	path := filepath.Join(dir, profile+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &snap, nil
}
