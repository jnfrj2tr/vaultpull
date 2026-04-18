package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/audit"
	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/diff"
	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/profile"
	"github.com/yourusername/vaultpull/internal/prompt"
	"github.com/yourusername/vaultpull/internal/vault"
)

var (
	flagProfile  string
	flagForce    bool
	flagDryRun   bool
	flagAuditLog string
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync secrets from Vault into a local .env file",
	RunE:  runSync,
}

func init() {
	syncCmd.Flags().StringVarP(&flagProfile, "profile", "p", "", "profile name to sync")
	syncCmd.Flags().BoolVar(&flagForce, "force", false, "overwrite without prompting")
	syncCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "preview changes without writing")
	syncCmd.Flags().StringVar(&flagAuditLog, "audit-log", "", "path to append audit entries")
	rootCmd.AddCommand(syncCmd)
}

func runSync(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	prof, err := profile.Select(cfg, flagProfile)
	if err != nil {
		return fmt.Errorf("select profile: %w", err)
	}

	p, err := profile.Resolve(cfg, prof)
	if err != nil {
		return fmt.Errorf("resolve profile: %w", err)
	}

	client, err := vault.NewClient(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	incoming, err := client.GetSecrets(p.VaultPath)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	changes := diff.Compare(p.OutputFile, incoming)
	diff.Print(cmd.OutOrStdout(), changes)

	if flagDryRun {
		fmt.Fprintln(cmd.OutOrStdout(), "dry-run: no changes written")
		return nil
	}

	if len(changes) > 0 && !flagForce {
		var confirmer prompt.Confirmer = prompt.New()
		ok, err := confirmer.Ask(fmt.Sprintf("Write changes to %s?", p.OutputFile))
		if err != nil {
			return fmt.Errorf("prompt: %w", err)
		}
		if !ok {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
	}

	w := env.NewWriter(p.OutputFile)
	if err := w.Write(incoming); err != nil {
		return fmt.Errorf("write env: %w", err)
	}

	if flagAuditLog != "" {
		if err := recordAudit(flagAuditLog, prof, changes); err != nil {
			return fmt.Errorf("audit log: %w", err)
		}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "synced %d secret(s) → %s\n", len(incoming), p.OutputFile)
	os.Exit(0)
	return nil
}

// recordAudit opens the audit log at path and records an entry for each changed key.
func recordAudit(path, prof string, changes []diff.Change) error {
	logger, err := audit.NewLogger(path)
	if err != nil {
		return err
	}
	for _, c := range changes {
		if err := logger.Record(prof, c.Key, c.Type); err != nil {
			return fmt.Errorf("record key %q: %w", c.Key, err)
		}
	}
	return nil
}
