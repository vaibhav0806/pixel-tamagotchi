package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ScanDirs []string `json:"scan_dirs"`
}

func DefaultDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".terminal-pet")
}

func DefaultStatePath() string {
	return filepath.Join(DefaultDir(), "state.json")
}

func DefaultConfigPath() string {
	return filepath.Join(DefaultDir(), "config.json")
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func Save(c *Config, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
