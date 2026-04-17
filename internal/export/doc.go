// Package export provides multi-format secret export for vaultpull.
//
// Supported formats:
//
//   - dotenv  — plain KEY=VALUE lines suitable for .env files
//   - json    — pretty-printed JSON object
//   - export  — shell-sourceable lines prefixed with "export"
//
// Usage:
//
//	e, err := export.New(export.FormatJSON)
//	if err != nil { ... }
//	err = e.Write(secrets, "/path/to/output.json")
//
// Pass "-" as the path to write to stdout.
package export
