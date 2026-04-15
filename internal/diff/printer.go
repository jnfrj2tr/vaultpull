package diff

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
)

// PrintOptions controls how the diff is rendered.
type PrintOptions struct {
	// NoColor disables ANSI colour codes.
	NoColor bool
	// ShowUnchanged includes unchanged keys in the output.
	ShowUnchanged bool
}

// Print writes a human-readable diff summary to w.
func Print(w io.Writer, result *Result, opts PrintOptions) {
	if len(result.Changes) == 0 {
		fmt.Fprintln(w, "No secrets to sync.")
		return
	}

	for _, c := range result.Changes {
		switch c.Kind {
		case Added:
			prefix := colorize("+ ", colorGreen, opts.NoColor)
			fmt.Fprintf(w, "%s%s\n", prefix, c.Key)
		case Updated:
			prefix := colorize("~ ", colorYellow, opts.NoColor)
			fmt.Fprintf(w, "%s%s\n", prefix, c.Key)
		case Unchanged:
			if opts.ShowUnchanged {
				prefix := colorize("  ", colorGray, opts.NoColor)
				fmt.Fprintf(w, "%s%s\n", prefix, c.Key)
			}
		}
	}

	added := len(result.Added())
	updated := len(result.Updated())

	parts := []string{}
	if added > 0 {
		parts = append(parts, colorize(fmt.Sprintf("%d added", added), colorGreen, opts.NoColor))
	}
	if updated > 0 {
		parts = append(parts, colorize(fmt.Sprintf("%d updated", updated), colorYellow, opts.NoColor))
	}
	if len(parts) > 0 {
		fmt.Fprintf(w, "\nSummary: %s\n", strings.Join(parts, ", "))
	}
}

func colorize(s, color string, noColor bool) string {
	if noColor {
		return s
	}
	return color + s + colorReset
}
