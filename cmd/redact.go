package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/env"
)

var (
	redactKeys     []string
	redactPatterns []string
	redactReveal   int
	redactInput    string
	redactOutput   string
)

func init() {
	redactCmd := &cobra.Command{
		Use:   "redact",
		Short: "Redact sensitive values in an env file",
		Long: `Reads an env file and writes a copy with sensitive values masked.

Keys can be targeted explicitly with --key or by regex with --pattern.
Use --reveal N to show the first N characters before masking.`,
		RunE: runRedact,
	}

	redactCmd.Flags().StringVarP(&redactInput, "input", "i", "", "Input .env file (required)")
	redactCmd.Flags().StringVarP(&redactOutput, "output", "o", "", "Output file (defaults to stdout)")
	redactCmd.Flags().StringArrayVar(&redactKeys, "key", nil, "Key name to redact (repeatable, case-insensitive)")
	redactCmd.Flags().StringArrayVar(&redactPatterns, "pattern", nil, "Regex pattern for keys to redact (repeatable)")
	redactCmd.Flags().IntVar(&redactReveal, "reveal", 0, "Show first N chars of redacted value")
	_ = redactCmd.MarkFlagRequired("input")

	rootCmd.AddCommand(redactCmd)
}

func runRedact(cmd *cobra.Command, _ []string) error {
	entries, err := env.Parse(redactInput)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	envm := env.ToMap(entries)

	opts := env.RedactOptions{
		Keys:          redactKeys,
		Patterns:      redactPatterns,
		PartialReveal: redactReveal,
	}

	redacted, err := env.RedactMap(envm, opts)
	if err != nil {
		return fmt.Errorf("redact: %w", err)
	}

	var sb strings.Builder
	for k, v := range redacted {
		fmt.Fprintf(&sb, "%s=%s\n", k, v)
	}

	if redactOutput == "" {
		fmt.Fprint(cmd.OutOrStdout(), sb.String())
		return nil
	}

	if err := os.WriteFile(redactOutput, []byte(sb.String()), 0o600); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "redacted env written to %s\n", redactOutput)
	return nil
}
