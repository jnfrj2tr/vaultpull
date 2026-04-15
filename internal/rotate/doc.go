// Package rotate provides functionality for fetching secrets from HashiCorp
// Vault and writing them to local .env files, with optional schedule-based
// rotation support.
//
// Basic usage:
//
//	client, _ := vault.NewClient(addr, token)
//	logger, _ := audit.NewLogger(logPath)
//	rot := rotate.New(client, logger)
//	result, err := rot.Rotate("production", "secret/data/app", ".env")
//
// Scheduled rotation:
//
//	sched := &rotate.Schedule{
//		Interval:    24 * time.Hour,
//		LastRotated: lastRun,
//	}
//	if sched.IsDue() {
//		rot.Rotate(...)
//	}
package rotate
