package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
	"vaultpull/internal/profile"
	"vaultpull/internal/config"
	"vaultpull/internal/vault"
)

var (
	filterInclude string
	filterExclude string
	filterPattern string
	filterOutput  string
)

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Fetch secrets and write only matching keys to an env file",
	RunE:  runFilter,
}

func init() {
	rootCmd.AddCommand(filterCmd)
	filterCmd.Flags().StringVar(&filterInclude, "include-prefix", "", "Only include keys with this prefix")
	filterCmd.Flags().StringVar(&filterExclude, "exclude-prefix", "", "Exclude keys with this prefix")
	filterCmd.Flags().StringVar(&filterPattern, "pattern", "", "Regex pattern keys must match")
	filterCmd.Flags().StringVarP(&filterOutput, "output", "o", "", "Output file (default: profile output_file)")
	filterCmd.Flags().StringVarP(&profileFlag, "profile", "p", "", "Profile to use")
}

func runFilter(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	prof, err := profile.Select(cfg, profileFlag)
	if err != nil {
		return err
	}

	p, err := profile.Resolve(cfg, prof)
	if err != nil {
		return err
	}

	client, err := vault.NewClient(cfg.VaultAddress, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := client.GetSecrets(p.VaultPath)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	filtered, err := env.Filter(secrets, env.FilterOptions{
		IncludePrefix: filterInclude,
		ExcludePrefix: filterExclude,
		Pattern:       filterPattern,
	})
	if err != nil {
		return fmt.Errorf("filter: %w", err)
	}

	out := filterOutput
	if out == "" {
		out = p.OutputFile
	}

	w := env.NewWriter(out)
	if err := w.Write(filtered); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	fmt.Fprintf(os.Stdout, "wrote %d keys to %s\n", len(filtered), out)
	return nil
}
