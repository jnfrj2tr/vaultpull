package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpull/internal/env"
)

var (
	aliasRulesFlag []string
	aliasInputFlag string
	aliasOutputFlag string
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Duplicate env keys under one or more alias names",
	Long: `Read a .env file, apply alias rules, and write the result.

Rules have the form SRC:ALIAS1,ALIAS2 (append +keep to retain the source key).

Example:
  vaultpull alias --rule DB_HOST:DATABASE_HOST,PG_HOST+keep \
                  --input .env --output .env.aliased`,
	RunE: runAlias,
}

func init() {
	aliasCmd.Flags().StringArrayVar(&aliasRulesFlag, "rule", nil, "alias rule: SRC:ALIAS1,ALIAS2[+keep] (repeatable)")
	aliasCmd.Flags().StringVar(&aliasInputFlag, "input", ".env", "source .env file")
	aliasCmd.Flags().StringVar(&aliasOutputFlag, "output", "", "destination file (default: overwrite input)")
	_ = aliasCmd.MarkFlagRequired("rule")
	rootCmd.AddCommand(aliasCmd)
}

func runAlias(cmd *cobra.Command, _ []string) error {
	parsed, err := env.Parse(aliasInputFlag)
	if err != nil {
		return fmt.Errorf("alias: reading %s: %w", aliasInputFlag, err)
	}
	src := env.ToMap(parsed)

	rules, err := env.ParseAliasRules(aliasRulesFlag)
	if err != nil {
		return fmt.Errorf("alias: %w", err)
	}

	out, err := env.Alias(src, rules)
	if err != nil {
		return fmt.Errorf("alias: %w", err)
	}

	dest := aliasOutputFlag
	if dest == "" {
		dest = aliasInputFlag
	}

	w, err := env.NewWriter(dest)
	if err != nil {
		return fmt.Errorf("alias: opening output: %w", err)
	}
	if err := w.Write(out); err != nil {
		fmt.Fprintf(os.Stderr, "alias: write error: %v\n", err)
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "wrote %d keys to %s\n", len(out), dest)
	return nil
}
