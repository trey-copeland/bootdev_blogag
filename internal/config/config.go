package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("User home directory not found: %w", err)
	}
	return filepath.Join(home, ".gatorconfig.json"), nil
}

func Read() (Config, error) {
	cfgPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	return readAt(cfgPath)
}

func readAt(cfgPath string) (Config, error) {
	bytes, err := os.ReadFile(cfgPath)
	if err != nil {
		return Config{}, fmt.Errorf("Error reading .gatorconfig: %v", err)
	}

	cfg := Config{}
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return Config{}, fmt.Errorf("Error unmarshalling .gatorconfig: %v", err)
	}

	return cfg, nil
}

func writeConfig(cfg Config) error {
	cfgPath, err := getConfigPath()
	if err != nil {
		return err
	}

	return writeConfigAt(cfgPath, cfg)
}

func writeConfigAt(cfgPath string, cfg Config) error {
	cfgJson, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Error marshalling config: %w", err)
	}
	const configFilePerm os.FileMode = 0o600
	if err := os.WriteFile(cfgPath, cfgJson, configFilePerm); err != nil {
		return fmt.Errorf("Config write error: %w", err)
	}

	return nil
}

func (c Config) SetUser(current_user_name string) error {
	c.CurrentUserName = current_user_name

	if err := writeConfig(c); err != nil {
		return err
	}

	return nil
}
