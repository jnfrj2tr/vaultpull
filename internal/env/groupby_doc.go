// Package env provides utilities for reading, writing, and transforming
// environment variable maps used by vaultpull profiles.
//
// GroupBy partitions a flat env map into named sub-maps based on key
// structure. Three modes are supported:
//
//   - prefix    – splits on the first occurrence of a delimiter (default "_")
//                 and uses the left-hand side as the group name.
//   - delimiter – splits on the last occurrence of the delimiter so that
//                 compound namespaces like "aws_us_east" become the group.
//   - pattern   – applies a regular expression with a capture group; the
//                 first captured sub-match becomes the group name.
//
// Keys that do not match any group rule are placed in a configurable
// fallback group (default "other").
package env
