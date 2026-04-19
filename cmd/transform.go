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
	transformName  string
	transformKeys  []string
	transformWrite bool
)

func init() {
	transformCmd := &cobra.Command{
		Use:   "transform",
		Short: "Fetch secrets and apply a value transformation before writing",
		RunE:  runTransform,
	}
	transformCmd.Flags().StringVarP(&transformName, "transform", "t", "trim", "Transform to apply (upper|lower|trim|base64)")
	transformCmd.Flags().StringSliceVar(&transformKeys, "keys", nil, "Limit transform to specific keys (default: all)")
	transformCmd.Flags().BoolVar(&transformWrite, "write", false, "Write transformed secrets to output file")
	transformCmd.Flags().StringVar(&profileFlag, "profile", "", "Profile name")
	rootCmd.AddCommand(transformCmd)
}

func runTransform(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load(configFile)
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
	tr := env.NewTransformer()
	transformed, err := tr.Apply(secrets, transformName, transformKeys)
	if err != nil {
		return err
	}
	if !transformWrite {
		for k, v := range transformed {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
		}
		return nil
	}
	w := env.NewWriter(p.OutputFile)
	if err := w.Write(transformed); err != nil {
		fmt.Fprintf(os.Stderr, "write error: %v\n", err)
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "wrote %d keys to %s\n", len(transformed), p.OutputFile)
	return nil
}
