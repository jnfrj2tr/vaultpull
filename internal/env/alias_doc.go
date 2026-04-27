// Package env provides utilities for manipulating environment variable maps.
//
// # Alias
//
// The Alias function duplicates values from one key to one or more alias keys.
// This is useful when a secret stored under a canonical name must also be
// available under legacy or framework-specific names without running a full
// sync twice.
//
// Rules are expressed as structs or parsed from compact string notation:
//
//	"DB_HOST:DATABASE_HOST,PG_HOST"        — rename to two aliases, drop source
//	"DB_HOST:DATABASE_HOST,PG_HOST+keep"   — same but retain the original key
//
// ParseAliasRules converts the string form into []AliasRule values suitable
// for passing to Alias.
package env
