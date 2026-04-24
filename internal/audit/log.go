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

// ReadAll reads and returns all audit log entries from the log file.
// It returns an error if the file cannot be opened or if any line
// contains malformed JSON.
func (l *Logger) ReadAll() ([]Entry, error) {
	data, err := os.ReadFile(l.path)
	if err != nil {
		return nil, fmt.Errorf("audit: read log file: %w", err)
	}

	var entries []Entry
	decoder := json.NewDecoder(bytesReader(data))
	for decoder.More() {
		var e Entry
		if err := decoder.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
