package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
)

var defaultsFile string

var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "Apply default values to an existing .env file",
	Long:  "Reads a .env file and inserts missing keys based on a defaults map defined inline via flags.",
	RunE:  runDefaults,
}

func init() {
	defaultsCmd.Flags().StringVarP(&defaultsFile, "file", "f", ".env", "Path to the .env file to update")
	rootCmd.AddCommand(defaultsCmd)
}

func runDefaults(cmd *cobra.Command, args []string) error {
	data, err := os.ReadFile(defaultsFile)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	entries, err := env.Parse(string(data))
	if err != nil {
		return fmt.Errorf("parsing env file: %w", err)
	}

	secrets := env.ToMap(entries)

	// Example built-in defaults; in production these would come from config.
	builtinDefaults := []env.DefaultEntry{
		{Key: "LOG_LEVEL", Default: "info"},
		{Key: "PORT", Default: "8080"},
	}

	updated, err := env.ApplyDefaults(secrets, builtinDefaults)
	if err != nil {
		return fmt.Errorf("applying defaults: %w", err)
	}

	w := env.NewWriter(defaultsFile)
	if err := w.Write(updated); err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "defaults applied to %s\n", defaultsFile)
	return nil
}
