package env

import (
	"bufio"
	"os"
	"strings"
)

// Merge reads an existing .env file and merges new secrets into it.
// Keys present in secrets overwrite existing values; other keys are preserved.
func Merge(filePath string, secrets map[string]string) (map[string]string, error) {
	existing, err := readEnvFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	for k, v := range secrets {
		existing[k] = v
	}
	return existing, nil
}

// readEnvFile parses a .env file into a key-value map.
// Lines starting with '#' and empty lines are ignored.
func readEnvFile(filePath string) (map[string]string, error) {
	result := make(map[string]string)

	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		result[key] = val
	}
	return result, scanner.Err()
}
