// Package prompt provides interactive terminal prompts for user confirmation.
package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Confirmer is the interface for asking yes/no questions.
type Confirmer interface {
	Ask(question string) (bool, error)
}

// Terminal reads confirmation from a terminal reader/writer pair.
type Terminal struct {
	In  io.Reader
	Out io.Writer
}

// New returns a Terminal confirmer backed by stdin/stdout.
func New() *Terminal {
	return &Terminal{
		In:  os.Stdin,
		Out: os.Stdout,
	}
}

// Ask prints question with a [y/N] prompt and returns true only if the user
// types "y" or "yes" (case-insensitive). Any other input — including an empty
// line — is treated as "no".
func (t *Terminal) Ask(question string) (bool, error) {
	fmt.Fprintf(t.Out, "%s [y/N]: ", question)

	scanner := bufio.NewScanner(t.In)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("prompt: read error: %w", err)
		}
		// EOF — treat as no
		return false, nil
	}

	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes", nil
}

// SkipConfirmer always returns true without prompting. Useful for --force flags.
type SkipConfirmer struct{}

func (s SkipConfirmer) Ask(_ string) (bool, error) { return true, nil }
