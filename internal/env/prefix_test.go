package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPrefix_Basic(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	out := AddPrefix(secrets, PrefixOptions{Prefix: "APP_"})
	assert.Equal(t, "localhost", out["APP_DB_HOST"])
	assert.Equal(t, "5432", out["APP_DB_PORT"])
	assert.Len(t, out, 2)
}

func TestAddPrefix_EmptyPrefix(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	out := AddPrefix(secrets, PrefixOptions{Prefix: ""})
	assert.Equal(t, "val", out["KEY"])
}

func TestAddPrefix_EmptyMap(t *testing.T) {
	out := AddPrefix(map[string]string{}, PrefixOptions{Prefix: "X_"})
	assert.Empty(t, out)
}

func TestStripPrefix_Basic(t *testing.T) {
	secrets := map[string]string{"APP_DB_HOST": "localhost", "APP_DB_PORT": "5432"}
	out := StripPrefix(secrets, "APP_")
	assert.Equal(t, "localhost", out["DB_HOST"])
	assert.Equal(t, "5432", out["DB_PORT"])
	assert.Len(t, out, 2)
}

func TestStripPrefix_NoMatch_PassThrough(t *testing.T) {
	secrets := map[string]string{"OTHER_KEY": "value"}
	out := StripPrefix(secrets, "APP_")
	assert.Equal(t, "value", out["OTHER_KEY"])
}

func TestStripPrefix_CollisionStrippedWins(t *testing.T) {
	// APP_KEY strips to KEY; KEY also present — stripped wins
	secrets := map[string]string{"APP_KEY": "from_vault", "KEY": "original"}
	out := StripPrefix(secrets, "APP_")
	assert.Equal(t, "from_vault", out["KEY"])
}

func TestStripPrefix_EmptyPrefix(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	out := StripPrefix(secrets, "")
	assert.Equal(t, "val", out["KEY"])
}
