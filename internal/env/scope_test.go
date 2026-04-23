package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var baseEnv = map[string]string{
	"APP_HOST":  "localhost",
	"APP_PORT":  "8080",
	"DB_HOST":   "db.local",
	"DB_PASS":   "secret",
	"LOG_LEVEL": "info",
}

func TestScope_IncludePrefix(t *testing.T) {
	result, err := Scope(baseEnv, ScopeOptions{
		Mode:   ScopeModeInclude,
		Scopes: []string{"APP_"},
	})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}, result)
}

func TestScope_ExcludePrefix(t *testing.T) {
	result, err := Scope(baseEnv, ScopeOptions{
		Mode:   ScopeModeExclude,
		Scopes: []string{"DB_"},
	})
	require.NoError(t, err)
	assert.NotContains(t, result, "DB_HOST")
	assert.NotContains(t, result, "DB_PASS")
	assert.Contains(t, result, "APP_HOST")
	assert.Contains(t, result, "LOG_LEVEL")
}

func TestScope_IncludeExact(t *testing.T) {
	result, err := Scope(baseEnv, ScopeOptions{
		Mode:   ScopeModeInclude,
		Scopes: []string{"DB_PASS", "LOG_LEVEL"},
		Exact:  true,
	})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{
		"DB_PASS":   "secret",
		"LOG_LEVEL": "info",
	}, result)
}

func TestScope_ExcludeExact(t *testing.T) {
	result, err := Scope(baseEnv, ScopeOptions{
		Mode:   ScopeModeExclude,
		Scopes: []string{"APP_HOST"},
		Exact:  true,
	})
	require.NoError(t, err)
	assert.NotContains(t, result, "APP_HOST")
	assert.Contains(t, result, "APP_PORT")
}

func TestScope_UnknownMode_ReturnsError(t *testing.T) {
	_, err := Scope(baseEnv, ScopeOptions{
		Mode:   "upsert",
		Scopes: []string{"APP_"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown mode")
}

func TestScope_EmptyScopes_ReturnsError(t *testing.T) {
	_, err := Scope(baseEnv, ScopeOptions{
		Mode:   ScopeModeInclude,
		Scopes: []string{},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least one scope")
}

func TestParseScopes_CommaSeparated(t *testing.T) {
	scopes := ParseScopes("APP_, DB_ , LOG_")
	assert.Equal(t, []string{"APP_", "DB_", "LOG_"}, scopes)
}

func TestParseScopes_EmptyTokensSkipped(t *testing.T) {
	scopes := ParseScopes("APP_,,DB_")
	assert.Equal(t, []string{"APP_", "DB_"}, scopes)
}
