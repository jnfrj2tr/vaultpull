package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Profile defines a named set of Vault settings and secret mappings.
type Profile struct {
	VaultAddr  string            `yaml:"vault_addr"`
	VaultToken string            `yaml:"vault_token"`
	MountPath  string            `yaml:"mount_path"`
	SecretPath string            `yaml:"secret_path"`
	OutputFile string            `yaml:"output_file"`
	Mapping    map[string]string `yaml:"mapping"`
}

// Config holds all named profiles.
type Config struct {
	Profiles map[string]Profile `yaml:"profiles"`
}

// Load reads and parses the YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}

	return &cfg, nil
}

// GetProfile retrieves a named profile, returning an error if not found.
func (c *Config) GetProfile(name string) (Profile, error) {
	p, ok := c.Profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("profile %q not found in config", name)
	}
	return p, nil
}
