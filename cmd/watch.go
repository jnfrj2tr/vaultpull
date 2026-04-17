package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/internal/watch"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Re-sync secrets whenever the config file changes",
	Long: `watch monitors .vaultpull.yaml for changes and automatically
runs the sync command whenever the file is saved.`,
	RunE: runWatch,
}

var watchDebounce time.Duration

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().DurationVar(&watchDebounce, "debounce", 300*time.Millisecond, "debounce delay between file and sync")
}

func runWatch(cmd *cobra.Command, args []string) error {
	cfg, _ := cmd.Root().PersistentFlags().GetString("config")
	if cfgFile == "" {
		cfgFile = ".vaultpull.yaml"
	}

	if _, err := os.Stat(cfgFile); err != nil {
		return fmt.Errorf("config file not found: %w", err)
	}

	fmt.Printf("Watching %s for changes (debounce: %s)…\n", cfgFile, watchDebounce)

	handler := func() error {
		fmt.Println("Config changed — running sync…")
		return runSync(cmd, args)
	}

	w := watch.New(cfgFile, watchDebounce, handler)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := w.Run(ctx)
	if errors.Is(err, context.Canceled) {
		fmt.Println("\nWatch stopped.")
		return nil
	}
	return err
}
