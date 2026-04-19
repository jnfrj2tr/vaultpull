package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlatten_SimpleMap(t *testing.T) {
	input := map[string]any{
		"host": "localhost",
		"port": "5432",
	}
	out, err := Flatten(input, FlattenOptions{})
	require.NoError(t, err)
	assert.Equal(t, "localhost", out["host"])
	assert.Equal(t, "5432", out["port"])
}

func TestFlatten_Nested(t *testing.T) {
	input := map[string]any{
		"db": map[string]any{
			"host": "localhost",
			"port": "5432",
		},
	}
	out, err := Flatten(input, FlattenOptions{Separator: "_"})
	require.NoError(t, err)
	assert.Equal(t, "localhost", out["db_host"])
	assert.Equal(t, "5432", out["db_port"])
}

func TestFlatten_Uppercase(t *testing.T) {
	input := map[string]any{
		"db": map[string]any{
			"host": "localhost",
		},
	}
	out, err := Flatten(input, FlattenOptions{Uppercase: true})
	require.NoError(t, err)
	assert.Equal(t, "localhost", out["DB_HOST"])
}

func TestFlatten_WithPrefix(t *testing.T) {
	input := map[string]any{
		"token": "abc123",
	}
	out, err := Flatten(input, FlattenOptions{Prefix: "APP", Uppercase: true})
	require.NoError(t, err)
	assert.Equal(t, "abc123", out["APP_TOKEN"])
}

func TestFlatten_NilValue(t *testing.T) {
	input := map[string]any{
		"key": nil,
	}
	out, err := Flatten(input, FlattenOptions{})
	require.NoError(t, err)
	assert.Equal(t, "", out["key"])
}

func TestFlatten_NonStringCoerced(t *testing.T) {
	input := map[string]any{
		"count": 42,
	}
	out, err := Flatten(input, FlattenOptions{})
	require.NoError(t, err)
	assert.Equal(t, "42", out["count"])
}

func TestFlatten_DefaultSeparator(t *testing.T) {
	input := map[string]any{
		"a": map[string]any{"b": "val"},
	}
	out, err := Flatten(input, FlattenOptions{Separator: ""})
	require.NoError(t, err)
	assert.Equal(t, "val", out["a_b"])
}
