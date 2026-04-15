package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SecretData holds the key-value pairs returned from a Vault secret path.
type SecretData map[string]string

// GetSecrets retrieves secrets from the given Vault KV path.
// It supports KV v2 by reading from the "data" sub-key when present.
func (c *Client) GetSecrets(path string) (SecretData, error) {
	url := fmt.Sprintf("%s/v1/%s", c.address, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("secret path not found: %s", path)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for path: %s", resp.StatusCode, path)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var raw struct {
		Data struct {
			Data map[string]interface{} `json:"data"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	result := make(SecretData, len(raw.Data.Data))
	for k, v := range raw.Data.Data {
		switch val := v.(type) {
		case string:
			result[k] = val
		default:
			result[k] = fmt.Sprintf("%v", val)
		}
	}

	return result, nil
}
