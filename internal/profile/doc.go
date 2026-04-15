// Package profile provides utilities for resolving and selecting
// vaultpull sync profiles from a loaded configuration.
//
// Profile resolution (Resolve) validates a named profile and fills in
// sensible defaults (e.g. output file path). Profile selection (Select)
// determines which profile name to use by checking, in order:
//
//  1. An explicit --profile flag value
//  2. The VAULTPULL_PROFILE environment variable
//  3. Auto-selection when exactly one profile is configured
//
// Example usage:
//
//	cfg, _ := config.Load(".vaultpull.yaml")
//	name, err := profile.Select(cfg, flagProfile)
//	rp, err   := profile.Resolve(cfg, name)
package profile
