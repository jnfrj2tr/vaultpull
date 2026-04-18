// Package notify provides hooks for sending notifications after a sync or rotate event.
package notify

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Event describes what happened during a sync or rotate operation.
type Event struct {
	Profile   string
	Operation string // "sync" | "rotate"
	Added     int
	Updated   int
	Removed   int
	Timestamp time.Time
}

// Notifier writes human-readable event summaries to a destination.
type Notifier struct {
	out io.Writer
}

// New returns a Notifier that writes to out. If out is nil, os.Stdout is used.
func New(out io.Writer) *Notifier {
	if out == nil {
		out = os.Stdout
	}
	return &Notifier{out: out}
}

// Notify formats and writes the event summary.
func (n *Notifier) Notify(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	parts := []string{}
	if e.Added > 0 {
		parts = append(parts, fmt.Sprintf("+%d added", e.Added))
	}
	if e.Updated > 0 {
		parts = append(parts, fmt.Sprintf("~%d updated", e.Updated))
	}
	if e.Removed > 0 {
		parts = append(parts, fmt.Sprintf("-%d removed", e.Removed))
	}
	summary := "no changes"
	if len(parts) > 0 {
		summary = strings.Join(parts, ", ")
	}
	_, err := fmt.Fprintf(n.out, "[%s] %s/%s: %s\n",
		e.Timestamp.Format(time.RFC3339),
		e.Operation,
		e.Profile,
		summary,
	)
	return err
}
