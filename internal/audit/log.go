package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Profile   string    `json:"profile"`
	SecretPath string   `json:"secret_path"`
	KeysWritten int     `json:"keys_written"`
	OutputFile  string  `json:"output_file"`
	Status      string  `json:"status"`
	Error       string  `json:"error,omitempty"`
}

// Logger writes audit entries to a log file.
type Logger struct {
	path string
}

// NewLogger creates a new Logger that appends entries to the given file path.
func NewLogger(path string) *Logger {
	return &Logger{path: path}
}

// Record appends an audit entry to the log file as a JSON line.
func (l *Logger) Record(e Entry) error {
	e.Timestamp = time.Now().UTC()

	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}

	_, err = fmt.Fprintf(f, "%s\n", data)
	if err != nil {
		return fmt.Errorf("audit: write entry: %w", err)
	}
	return nil
}
