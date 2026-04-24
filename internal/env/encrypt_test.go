package env

import (
	"strings"
	"testing"
)

var testPassphrase = []byte("supersecretpassphrase")

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	original := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
	}

	encrypted, err := EncryptMap(original, testPassphrase, nil)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}

	for k, v := range encrypted {
		if v == original[k] {
			t.Errorf("key %q: value was not encrypted", k)
		}
	}

	decrypted, err := DecryptMap(encrypted, testPassphrase, nil)
	if err != nil {
		t.Fatalf("DecryptMap: %v", err)
	}

	for k, want := range original {
		if got := decrypted[k]; got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestEncryptMap_OnlyTargetedKeys(t *testing.T) {
	env := map[string]string{
		"SECRET": "topsecret",
		"PUBLIC": "visible",
	}

	encrypted, err := EncryptMap(env, testPassphrase, []string{"SECRET"})
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}

	if encrypted["PUBLIC"] != "visible" {
		t.Errorf("PUBLIC should not be encrypted, got %q", encrypted["PUBLIC"])
	}
	if encrypted["SECRET"] == "topsecret" {
		t.Error("SECRET should have been encrypted")
	}
}

func TestDecryptMap_OnlyTargetedKeys(t *testing.T) {
	env := map[string]string{
		"SECRET": "topsecret",
		"PUBLIC": "visible",
	}

	encrypted, err := EncryptMap(env, testPassphrase, []string{"SECRET"})
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}

	decrypted, err := DecryptMap(encrypted, testPassphrase, []string{"SECRET"})
	if err != nil {
		t.Fatalf("DecryptMap: %v", err)
	}

	if decrypted["SECRET"] != "topsecret" {
		t.Errorf("SECRET: got %q, want %q", decrypted["SECRET"], "topsecret")
	}
	if decrypted["PUBLIC"] != "visible" {
		t.Errorf("PUBLIC: got %q, want %q", decrypted["PUBLIC"], "visible")
	}
}

func TestDecryptMap_WrongPassphrase(t *testing.T) {
	env := map[string]string{"KEY": "value"}

	encrypted, err := EncryptMap(env, testPassphrase, nil)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}

	_, err = DecryptMap(encrypted, []byte("wrongpassphrase"), nil)
	if err == nil {
		t.Error("expected error with wrong passphrase, got nil")
	}
}

func TestDecryptMap_InvalidBase64(t *testing.T) {
	env := map[string]string{"KEY": "not-valid-base64!!!"}
	_, err := DecryptMap(env, testPassphrase, nil)
	if err == nil {
		t.Error("expected error for invalid base64, got nil")
	}
	if !strings.Contains(err.Error(), "base64") && !strings.Contains(err.Error(), "decrypt") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestEncryptMap_EmptyMap(t *testing.T) {
	out, err := EncryptMap(map[string]string{}, testPassphrase, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}
