package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/vaultpull/internal/config"
	"github.com/yourorg/vaultpull/internal/lint"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Validate the vaultpull configuration file",
	RunE:  runLint,
}

func init() {
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, _ []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	if cfgPath == "" {
		cfgPath = ".vaultpull.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	issues := lint.Run(cfg)

	if len(issues) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✓ configuration looks good")
		return nil
	}

	for _, i := range issues {
		fmt.Fprintln(cmd.OutOrStdout(), i)
	}

	if lint.HasErrors(issues) {
		return exitWithLintError(cmd)
	}

	return nil
}

// exitWithLintError prints a summary message to stderr and exits with a
// non-zero status code to signal lint failure to the calling process.
func exitWithLintError(cmd *cobra.Command) error {
	fmt.Fprintln(os.Stderr, "lint failed: fix errors before syncing")
	os.Exit(1)
	return nil // unreachable, but satisfies the error return type
}
