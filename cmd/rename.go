package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
)

var (
	renameRules    []string
	renameInput    string
	renameOutput   string
	renameOverwrite bool
)

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename keys in a .env file using FROM=TO rules",
	RunE:  runRename,
}

func init() {
	renameCmd.Flags().StringArrayVarP(&renameRules, "rule", "r", nil, "Rename rule in FROM=TO format (repeatable)")
	renameCmd.Flags().StringVarP(&renameInput, "input", "i", ".env", "Input .env file")
	renameCmd.Flags().StringVarP(&renameOutput, "output", "o", "", "Output file (defaults to input file)")
	renameCmd.Flags().BoolVar(&renameOverwrite, "overwrite", false, "Overwrite target key if it already exists")
	_ = renameCmd.MarkFlagRequired("rule")
	RootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) error {
	rules, err := env.ParseRenameRules(renameRules)
	if err != nil {
		return fmt.Errorf("invalid rename rules: %w", err)
	}

	entries, err := env.Parse(renameInput)
	if err != nil {
		return fmt.Errorf("failed to parse input file: %w", err)
	}

	secrets := env.ToMap(entries)

	renamed, err := env.Rename(secrets, rules, renameOverwrite)
	if err != nil {
		return fmt.Errorf("rename failed: %w", err)
	}

	dest := renameOutput
	if dest == "" {
		dest = renameInput
	}

	w := env.NewWriter(dest)
	if err := w.Write(renamed); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Renamed %d rule(s) → wrote %s\n", len(rules), dest)
	return nil
}
