package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/vaultpull/internal/env"
)

var (
	typecheckFile  string
	typecheckRules []string
)

func init() {
	typecheckCmd := &cobra.Command{
		Use:   "typecheck",
		Short: "Validate env variable types against declared rules",
		Long: `Reads a .env file and checks that each key matches its declared type.

Rules are specified as KEY=type pairs where type is one of:
  string, int, float, bool, url, email

Example:
  vaultpull typecheck --file .env --rule PORT=int --rule ENABLED=bool`,
		RunE: runTypecheck,
	}

	typecheckCmd.Flags().StringVarP(&typecheckFile, "file", "f", ".env", "Path to the .env file to validate")
	typecheckCmd.Flags().StringArrayVar(&typecheckRules, "rule", nil, "Type rule in KEY=type format (repeatable)")
	_ = typecheckCmd.MarkFlagRequired("rule")

	rootCmd.AddCommand(typecheckCmd)
}

func runTypecheck(cmd *cobra.Command, _ []string) error {
	rules, err := env.ParseTypeRules(typecheckRules)
	if err != nil {
		return fmt.Errorf("invalid rule: %w", err)
	}

	entries, err := env.Parse(typecheckFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", typecheckFile, err)
	}

	environment := env.ToMap(entries)
	violations := env.CheckTypes(environment, rules)

	if len(violations) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✓ All type checks passed.")
		return nil
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "Type check violations in %s:\n", typecheckFile)
	for _, v := range violations {
		fmt.Fprintf(cmd.ErrOrStderr(), "  %s\n", v.Error())
	}
	os.Exit(1)
	return nil
}
