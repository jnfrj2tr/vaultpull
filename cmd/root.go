package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	profile    string
	configFile string
)

var rootCmd = &cobra.Command{
	Use:   "vaultpull",
	Short: "Sync secrets from HashiCorp Vault into local .env files",
	Long: `vaultpull is a CLI tool that pulls secrets from HashiCorp Vault
and writes them into local .env files with support for multiple profiles.`,
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "profile name to use from config")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", ".vaultpull.yaml", "path to config file")
}
