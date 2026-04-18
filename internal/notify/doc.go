// Package notify provides lightweight notification support for vaultpull events.
//
// Two notifier types are available:
//
//   - Notifier: writes a human-readable summary line to any io.Writer.
//   - WebhookNotifier: POSTs a JSON payload to an HTTP endpoint.
//
// Both accept an Event struct that describes what changed during a sync or
// rotate operation.
package notify
