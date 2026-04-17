// Package export provides functionality to export secrets to various output formats.
package export

import (
	"encoding/json"
	"fmt"
	"os"\n	"sort"
	"strings"
)

// Format represents an export output format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatExport Format = "export"
)

// Exporter writes secrets to stdout or a file in a given format.
type Exporter struct {
	format Format
}

// New creates a new Exporter for the given format.
func New(format Format) (*Exporter, error) {
	switch format {
	case FormatDotenv, FormatJSON, FormatExport:
		return &Exporter{format: format}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %q", format)
	}
}

// Write serializes secrets to the given file path, or stdout if path is "-".
func (e *Exporter) Write(secrets map[string]string, path string) error {
	var out []byte
	var err error

	switch e.format {
	case FormatJSON:
		out, err = json.MarshalIndent(secrets, "", "  ")
		if err != nil {
			return fmt.Errorf("json marshal: %w", err)
		}
		out = append(out, '\n')
	case FormatDotenv:
		out = []byte(renderDotenv(secrets, false))
	case FormatExport:
		out = []byte(renderDotenv(secrets, true))
	}

	if path == "-" {
		_, err = os.Stdout.Write(out)
		return err
	}
	return os.WriteFile(path, out, 0600)
}

func renderDotenv(secrets map[string]string, withExport bool) string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := secrets[k]
		if strings.ContainsAny(v, " \t\n#") {
			v = fmt.Sprintf("%q", v)
		}
		if withExport {
			fmt.Fprintf(&sb, "export %s=%s\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}
	return sb.String()
}
