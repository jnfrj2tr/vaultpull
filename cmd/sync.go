package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/env"
	"vaultpull/internal/vault"
)

var (
	profileFlag string
	mergeFlag   bool
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync secrets from Vault into a local .env file",
	RunE:  runSync,
}

func init() {
	syncCmd.Flags().StringVarP(&profileFlag, "profile", "p", "default", "Config profile to use")
	syncCmd.Flags().BoolVarP(&mergeFlag, "merge", "m", false, "Merge with existing .env instead of overwriting")
	RootCmd.AddCommand(syncCmd)
}

func runSync(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	profile, err := cfg.GetProfile(profileFlag)
	if err != nil {
		return fmt.Errorf("profile %q not found: %w", profileFlag, err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.Token)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	secrets, err := client.GetSecrets(profile.SecretPath)
	if err != nil {
		return fmt.Errorf("fetching secrets from %q: %w", profile.SecretPath, err)
	}

	output := secrets
	if mergeFlag {
		output, err = env.Merge(profile.EnvFile, secrets)
		if err != nil {
			return fmt.Errorf("merging env file: %w", err)
		}
	}

	writer := env.NewWriter(profile.EnvFile)
	if err := writer.Write(output); err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}

	log.Printf("Synced %d secret(s) to %s (profile: %s)\n", len(secrets), profile.EnvFile, profileFlag)
	return nil
}
