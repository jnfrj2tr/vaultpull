package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/internal/diff"
	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/prompt"
)

var (
	applyDryRun      bool
	applySkipUnchanged bool
	applyNoConfirm   bool
	applyOutputFile  string
)

var applyCmd = &cobra.Command{
	Use:   "apply <base.env> <incoming.env>",
	Short: "Apply diff changes from incoming env into base env file",
	Args:  cobra.ExactArgs(2),
	RunE:  runApply,
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().BoolVar(&applyDryRun, "dry-run", false, "Preview changes without writing")
	applyCmd.Flags().BoolVar(&applySkipUnchanged, "skip-unchanged", false, "Exclude unchanged keys from output")
	applyCmd.Flags().BoolVar(&applyNoConfirm, "yes", false, "Skip confirmation prompt")
	applyCmd.Flags().StringVarP(&applyOutputFile, "output", "o", "", "Write result to file instead of stdout")
}

func runApply(cmd *cobra.Command, args []string) error {
	baseFile := args[0]
	incomingFile := args[1]

	baseEnv, err := env.Parse(baseFile)
	if err != nil {
		return fmt.Errorf("reading base file: %w", err)
	}

	incomingEnv, err := env.Parse(incomingFile)
	if err != nil {
		return fmt.Errorf("reading incoming file: %w", err)
	}

	changes := diff.Compare(baseEnv.ToMap(), incomingEnv.ToMap())

	fmt.Fprintln(cmd.OutOrStdout(), env.SummarizeDiff(changes))

	if applyDryRun {
		fmt.Fprintln(cmd.OutOrStdout(), "[dry-run] no changes written")
		return nil
	}

	if !applyNoConfirm {
		p := prompt.New(os.Stdin, os.Stdout)
		ok, err := p.Ask("Apply these changes?")
		if err != nil || !ok {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
	}

	opts := env.ApplyOptions{
		SkipUnchanged: applySkipUnchanged,
	}

	result, err := env.ApplyDiff(baseEnv.ToMap(), changes, opts)
	if err != nil {
		return fmt.Errorf("applying diff: %w", err)
	}

	dest := applyOutputFile
	if dest == "" {
		dest = baseFile
	}

	w := env.NewWriter(dest)
	if err := w.Write(result); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "wrote %d keys to %s\n", len(result), dest)
	return nil
}
