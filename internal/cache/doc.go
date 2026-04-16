// Package cache implements a simple disk-based cache for Vault secret snapshots.
//
// Cache entries are stored as JSON files under .vaultpull_cache/ and are keyed
// by profile name. Each entry records the fetched secrets, the time of fetch,
// and a TTL. Callers should check IsExpired before using a cached entry to
// decide whether a fresh Vault fetch is required.
//
// Files are written with mode 0600 to prevent other users from reading secrets.
package cache
