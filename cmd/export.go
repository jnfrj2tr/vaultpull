package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/export"
	"vaultpull/internal/profile"
	"vaultpull/internal/vault"
)

var (
	exportFormat  string
	exportOutput  string
	exportProfile string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export secrets from Vault to a file in a specified format",
	RunE:  runExport,
}

func init() {
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "dotenv", "Output format: dotenv, json, export")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file path (required)")
	exportCmd.Flags().StringVarP(&exportProfile, "profile", "p", "", "Profile to use")
	_ = exportCmd.MarkFlagRequired("output")
	RootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	prof, err := profile.Select(cfg, exportProfile, os.Getenv("VAULTPULL_PROFILE"))
	if err != nil {
		return fmt.Errorf("select profile: %w", err)
	}

	resolvedProf, err := profile.Resolve(cfg, prof)
	if err != nil {
		return fmt.Errorf("resolve profile: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddress, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := client.GetSecrets(resolvedProf.VaultPath)
	if err != nil {
		return fmt.Errorf("get secrets: %w", err)
	}

	exporter, err := export.New(exportFormat)
	if err != nil {
		return fmt.Errorf("create exporter: %w", err)
	}

	f, err := os.Create(exportOutput)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer f.Close()

	if err := exporter.Write(f, secrets); err != nil {
		return fmt.Errorf("write export: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Exported %d secrets to %s (format: %s)\n", len(secrets), exportOutput, exportFormat)
	return nil
}
