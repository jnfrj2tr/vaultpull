package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
)

var (
	coerceRulesFlag []string
	coerceInputFlag string
	coerceOutputFlag string
)

var coerceCmd = &cobra.Command{
	Use:   "coerce",
	Short: "Apply type coercion rules to an env file",
	Long: `Reads a .env file and rewrites values according to coercion rules.

Rules are specified as KEY:type pairs, where type is one of:
  string, bool, int, float, upper, lower

Example:
  vaultpull coerce --input .env --rule PORT:int --rule DEBUG:bool --output .env.coerced`,
	RunE: runCoerce,
}

func init() {
	rootCmd.AddCommand(coerceCmd)
	coerceCmd.Flags().StringVarP(&coerceInputFlag, "input", "i", ".env", "Input .env file")
	coerceCmd.Flags().StringVarP(&coerceOutputFlag, "output", "o", "", "Output file (defaults to input file)")
	coerceCmd.Flags().StringArrayVarP(&coerceRulesFlag, "rule", "r", nil, "Coercion rule as KEY:type (repeatable)")
	_ = coerceCmd.MarkFlagRequired("rule")
}

func runCoerce(cmd *cobra.Command, _ []string) error {
	rules, err := env.ParseCoerceRules(coerceRulesFlag)
	if err != nil {
		return fmt.Errorf("parsing rules: %w", err)
	}

	entries, err := env.Parse(coerceInputFlag)
	if err != nil {
		return fmt.Errorf("reading %s: %w", coerceInputFlag, err)
	}

	envMap := env.ToMap(entries)

	coerced, err := env.Coerce(envMap, rules)
	if err != nil {
		return fmt.Errorf("coercing values: %w", err)
	}

	dest := coerceOutputFlag
	if dest == "" {
		dest = coerceInputFlag
	}

	w := env.NewWriter(dest)
	if err := w.Write(coerced); err != nil {
		return fmt.Errorf("writing %s: %w", dest, err)
	}

	fmt.Fprintf(os.Stdout, "coerced %d rule(s) → %s\n", len(rules), dest)
	return nil
}
