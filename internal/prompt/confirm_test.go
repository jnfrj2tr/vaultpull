package prompt_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/prompt"
)

func newTerminal(input string) *prompt.Terminal {
	return &prompt.Terminal{
		In:  strings.NewReader(input),
		Out: &bytes.Buffer{},
	}
}

func TestAsk_YesAnswers(t *testing.T) {
	for _, ans := range []string{"y", "Y", "yes", "YES", "Yes"} {
		t.Run(ans, func(t *testing.T) {
			term := newTerminal(ans + "\n")
			ok, err := term.Ask("Continue?")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Errorf("expected true for input %q", ans)
			}
		})
	}
}

func TestAsk_NoAnswers(t *testing.T) {
	for _, ans := range []string{"n", "N", "no", "", "maybe"} {
		t.Run("input="+ans, func(t *testing.T) {
			term := newTerminal(ans + "\n")
			ok, err := term.Ask("Continue?")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ok {
				t.Errorf("expected false for input %q", ans)
			}
		})
	}
}

func TestAsk_EOF(t *testing.T) {
	term := newTerminal("") // empty reader → immediate EOF
	ok, err := term.Ask("Continue?")
	if err != nil {
		t.Fatalf("unexpected error on EOF: %v", err)
	}
	if ok {
		t.Error("expected false on EOF")
	}
}

func TestAsk_PrintsQuestion(t *testing.T) {
	out := &bytes.Buffer{}
	term := &prompt.Terminal{In: strings.NewReader("y\n"), Out: out}
	_, _ = term.Ask("Overwrite secrets?")
	if !strings.Contains(out.String(), "Overwrite secrets?") {
		t.Errorf("output %q does not contain question", out.String())
	}
}

func TestSkipConfirmer_AlwaysTrue(t *testing.T) {
	s := prompt.SkipConfirmer{}
	ok, err := s.Ask("anything")
	if err != nil || !ok {
		t.Errorf("SkipConfirmer.Ask() = (%v, %v), want (true, nil)", ok, err)
	}
}
