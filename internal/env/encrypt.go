package env

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// EncryptMap encrypts the values of the specified keys in the map using AES-GCM.
// If keys is empty, all values are encrypted.
// The key must be 16, 24, or 32 bytes long (AES-128, AES-192, AES-256).
func EncryptMap(env map[string]string, passphrase []byte, keys []string) (map[string]string, error) {
	block, err := aes.NewCipher(normalizeKey(passphrase))
	if err != nil {
		return nil, fmt.Errorf("encrypt: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("encrypt: create gcm: %w", err)
	}

	target := toSet(keys)
	out := make(map[string]string, len(env))
	for k, v := range env {
		if len(target) == 0 || target[k] {
			enc, err := encryptValue(gcm, v)
			if err != nil {
				return nil, fmt.Errorf("encrypt: key %q: %w", k, err)
			}
			out[k] = enc
		} else {
			out[k] = v
		}
	}
	return out, nil
}

// DecryptMap decrypts AES-GCM encrypted values for the specified keys.
// If keys is empty, all values are decrypted.
func DecryptMap(env map[string]string, passphrase []byte, keys []string) (map[string]string, error) {
	block, err := aes.NewCipher(normalizeKey(passphrase))
	if err != nil {
		return nil, fmt.Errorf("decrypt: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("decrypt: create gcm: %w", err)
	}

	target := toSet(keys)
	out := make(map[string]string, len(env))
	for k, v := range env {
		if len(target) == 0 || target[k] {
			dec, err := decryptValue(gcm, v)
			if err != nil {
				return nil, fmt.Errorf("decrypt: key %q: %w", k, err)
			}
			out[k] = dec
		} else {
			out[k] = v
		}
	}
	return out, nil
}

func encryptValue(gcm cipher.AEAD, plaintext string) (string, error) {
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptValue(gcm cipher.AEAD, encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", errors.New("ciphertext too short")
	}
	plaintext, err := gcm.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", fmt.Errorf("aes-gcm open: %w", err)
	}
	return string(plaintext), nil
}

// normalizeKey pads or truncates passphrase to 32 bytes for AES-256.
func normalizeKey(p []byte) []byte {
	key := make([]byte, 32)
	copy(key, p)
	return key
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
