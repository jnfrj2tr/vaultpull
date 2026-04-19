package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
	"vaultpull/internal/env/parser"
)

var (
	interpolateInput  string
	interpolateOutput string
	interpolateAllowEnv bool
)

func init() {
	interpolateCmd := &cobra.Command{
		Use:   "interpolate",
		Short: "Resolve \${VAR} references in an env file using secrets or OS env",
		RunE:  runInterpolate,
	}
	interpolateCmd.Flags().StringVarP(&interpolateInput, "input", "i", "", "Input .env file (required)")
	interpolateCmd.Flags().StringVarP(&interpolateOutput, "output", "o", "", "Output .env file (defaults to input)")
	interpolateCmd.Flags().BoolVar(&interpolateAllowEnv, "allow-env", false, "Fall back to OS environment variables")
	_ = interpolateCmd.MarkFlagRequired("input")
	rootCmd.AddCommand(interpolateCmd)
}

func runInterpolate(cmd *cobra.Command, args []string) error {
	entries, err := parser.ParseFile(interpolateInput)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	secrets := parser.ToMap(entries)

	resolved, err := env.Interpolate(secrets, interpolateAllowEnv)
	if err != nil {
		return fmt.Errorf("interpolate: %w", err)
	}

	dest := interpolateOutput
	if dest == "" {
		dest = interpolateInput
	}

	w := env.NewWriter(dest)
	if err := w.Write(resolved); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	fmt.Fprintf(os.Stdout, "interpolated %d keys -> %s\n", len(resolved), dest)
	return nil
}
