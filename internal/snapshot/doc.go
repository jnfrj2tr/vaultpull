// Package snapshot captures point-in-time secret state for a profile
// and provides drift detection by comparing a saved snapshot against
// the current secrets fetched from Vault.
//
// Typical usage:
//
//	snap, _ := snapshot.Load(dir, profile)
//	drift := snapshot.DetectDrift(snap, currentSecrets)
//	if drift != nil && !drift.Clean { /* handle drift */ }
//	snapshot.Save(dir, profile, currentSecrets)
package snapshot
