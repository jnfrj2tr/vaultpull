package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/vaultpull/vaultpull/internal/env"
)

var (
	groupByMode      string
	groupByDelimiter string
	groupByPattern   string
	groupByFallback  string
	groupByInput     string
	groupByJSON      bool
)

func init() {
	groupByCmd := &cobra.Command{
		Use:   "groupby",
		Short: "Partition env vars into named groups by prefix, delimiter, or pattern",
		RunE:  runGroupBy,
	}

	groupByCmd.Flags().StringVarP(&groupByMode, "mode", "m", "prefix", "Grouping mode: prefix | delimiter | pattern")
	groupByCmd.Flags().StringVarP(&groupByDelimiter, "delimiter", "d", "_", "Delimiter used in prefix/delimiter modes")
	groupByCmd.Flags().StringVarP(&groupByPattern, "pattern", "p", "", "Regex with one capture group (pattern mode)")
	groupByCmd.Flags().StringVar(&groupByFallback, "fallback", "other", "Group name for unmatched keys")
	groupByCmd.Flags().StringVarP(&groupByInput, "file", "f", "", "Input .env file (default: stdin)")
	groupByCmd.Flags().BoolVar(&groupByJSON, "json", false, "Output result as JSON")

	rootCmd.AddCommand(groupByCmd)
}

func runGroupBy(cmd *cobra.Command, args []string) error {
	var entries []env.Entry
	var err error

	if groupByInput != "" {
		entries, err = env.Parse(groupByInput)
	} else {
		entries, err = env.Parse(os.Stdin.Name())
	}
	if err != nil {
		return fmt.Errorf("groupby: failed to parse input: %w", err)
	}

	envMap := env.ToMap(entries)

	opts := env.GroupByOptions{
		Mode:      env.GroupByMode(groupByMode),
		Delimiter: groupByDelimiter,
		Pattern:   groupByPattern,
		Fallback:  groupByFallback,
	}

	result, err := env.GroupBy(envMap, opts)
	if err != nil {
		return err
	}

	if groupByJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	groups := make([]string, 0, len(result))
	for g := range result {
		groups = append(groups, g)
	}
	sort.Strings(groups)

	for _, g := range groups {
		fmt.Printf("[%s]\n", g)
		keys := make([]string, 0, len(result[g]))
		for k := range result[g] {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("  %s=%s\n", k, result[g][k])
		}
	}
	return nil
}
