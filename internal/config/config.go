// Package config
// Description: This file contains the configuration management for sokru.
// (c) 2024 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	DotfilesDir  string `yaml:"dotfiles_dir"`
	SymlinksFile string `yaml:"symlinks_file"`
	OS           string `yaml:"os"`
	Verbose      bool   `yaml:"verbose"`
	DryRun       bool   `yaml:"dry_run"`
}

const (
	configDir  = ".sokru"
	configFile = "config.yaml"
)

// Global configuration instance
var globalConfig *Config

// GetDefaultConfig returns a Config with default values
func GetDefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "~"
	}

	return &Config{
		DotfilesDir:  filepath.Join(homeDir, "dotfiles"),
		SymlinksFile: filepath.Join(homeDir, "dotfiles", "symlinks.yaml"),
		OS:           runtime.GOOS,
		Verbose:      false,
		DryRun:       false,
	}
}

// getConfigPath returns the full path to the config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return filepath.Join(homeDir, configDir, configFile), nil
}

// ensureConfigDir creates the config directory if it doesn't exist
func ensureConfigDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDirPath := filepath.Join(homeDir, configDir)
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return nil
}

// LoadConfig reads the configuration from ~/.sokru/config.yaml
// If the file doesn't exist, it returns the default configuration
func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return GetDefaultConfig(), nil
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig writes the configuration to ~/.sokru/config.yaml
func SaveConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Ensure config directory exists
	if err := ensureConfigDir(); err != nil {
		return err
	}

	// Get config file path
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Marshal config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfig returns the global configuration instance
// If not loaded yet, it loads the configuration
func GetConfig() (*Config, error) {
	if globalConfig == nil {
		config, err := LoadConfig()
		if err != nil {
			return nil, err
		}
		globalConfig = config
	}
	return globalConfig, nil
}

// SetConfig sets the global configuration instance
func SetConfig(config *Config) {
	globalConfig = config
}

// UpdateConfig updates specific fields in the configuration and saves it
func UpdateConfig(updateFunc func(*Config)) error {
	config, err := GetConfig()
	if err != nil {
		return err
	}

	updateFunc(config)

	if err := SaveConfig(config); err != nil {
		return err
	}

	return nil
}

// GetConfigPath returns the path to the config file (exported for external use)
func GetConfigPath() (string, error) {
	return getConfigPath()
}
