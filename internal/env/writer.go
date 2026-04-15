package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer handles writing secrets to .env files.
type Writer struct {
	filePath string
}

// NewWriter creates a new Writer for the given file path.
func NewWriter(filePath string) *Writer {
	return &Writer{filePath: filePath}
}

// Write writes the provided secrets map to the .env file.
// Existing file contents are overwritten.
func (w *Writer) Write(secrets map[string]string) error {
	file, err := os.Create(w.filePath)
	if err != nil {
		return fmt.Errorf("failed to create env file %q: %w", w.filePath, err)
	}
	defer file.Close()

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		line := fmt.Sprintf("%s=%s\n", k, quoteValue(secrets[k]))
		if _, err := file.WriteString(line); err != nil {
			return fmt.Errorf("failed to write key %q: %w", k, err)
		}
	}
	return nil
}

// quoteValue wraps a value in double quotes if it contains spaces or special chars.
func quoteValue(v string) string {
	if strings.ContainsAny(v, " \t\n#") {
		v = strings.ReplaceAll(v, `"`, `\"`)
		return `"` + v + `"`
	}
	return v
}
