// Package mask provides helpers for redacting or partially obscuring secret
// values before they are surfaced in user-visible output such as diff previews,
// audit log entries, or terminal prompts.
//
// Use Redact to blank an entire map of secrets, RedactValue for a single value,
// or PartialRedact when a short prefix hint is acceptable (e.g. showing the
// first character of a token to aid debugging without exposing the full value).
package mask
