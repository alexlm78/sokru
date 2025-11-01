// Package test
// Description: Unit tests for configuration management
// (c) 2024 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexlm78/sokru/internal/config"
)

func TestGetDefaultConfig(t *testing.T) {
	cfg := config.GetDefaultConfig()

	if cfg == nil {
		t.Fatal("GetDefaultConfig returned nil")
	}

	if cfg.DotfilesDir == "" {
		t.Error("DotfilesDir should not be empty")
	}

	if cfg.SymlinksFile == "" {
		t.Error("SymlinksFile should not be empty")
	}

	if cfg.OS == "" {
		t.Error("OS should not be empty")
	}

	if cfg.Language != "en" {
		t.Errorf("Expected default language 'en', got '%s'", cfg.Language)
	}

	if cfg.Verbose {
		t.Error("Verbose should be false by default")
	}

	if cfg.DryRun {
		t.Error("DryRun should be false by default")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "sokru-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override the config directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create test config
	testConfig := &config.Config{
		DotfilesDir:  "/test/dotfiles",
		SymlinksFile: "/test/symlinks.yml",
		OS:           "linux",
		Language:     "es",
		Verbose:      true,
		DryRun:       false,
	}

	// Save config
	err = config.SaveConfig(testConfig)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify config file was created
	configPath := filepath.Join(tempDir, ".config", "sokru", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", configPath)
	}

	// Load config
	loadedConfig, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify loaded config matches saved config
	if loadedConfig.DotfilesDir != testConfig.DotfilesDir {
		t.Errorf("DotfilesDir mismatch: expected %s, got %s", testConfig.DotfilesDir, loadedConfig.DotfilesDir)
	}

	if loadedConfig.SymlinksFile != testConfig.SymlinksFile {
		t.Errorf("SymlinksFile mismatch: expected %s, got %s", testConfig.SymlinksFile, loadedConfig.SymlinksFile)
	}

	if loadedConfig.OS != testConfig.OS {
		t.Errorf("OS mismatch: expected %s, got %s", testConfig.OS, loadedConfig.OS)
	}

	if loadedConfig.Language != testConfig.Language {
		t.Errorf("Language mismatch: expected %s, got %s", testConfig.Language, loadedConfig.Language)
	}

	if loadedConfig.Verbose != testConfig.Verbose {
		t.Errorf("Verbose mismatch: expected %v, got %v", testConfig.Verbose, loadedConfig.Verbose)
	}

	if loadedConfig.DryRun != testConfig.DryRun {
		t.Errorf("DryRun mismatch: expected %v, got %v", testConfig.DryRun, loadedConfig.DryRun)
	}
}

func TestLoadConfigNonExistent(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "sokru-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override the config directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Load config when file doesn't exist (should return default config)
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig should not fail when file doesn't exist: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadConfig returned nil")
	}

	// Should return default config
	defaultCfg := config.GetDefaultConfig()
	if cfg.Language != defaultCfg.Language {
		t.Errorf("Expected default language %s, got %s", defaultCfg.Language, cfg.Language)
	}
}

func TestUpdateConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "sokru-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override the config directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create initial config
	initialConfig := config.GetDefaultConfig()
	initialConfig.Verbose = false
	err = config.SaveConfig(initialConfig)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Update config
	err = config.UpdateConfig(func(c *config.Config) {
		c.Verbose = true
		c.Language = "es"
	})
	if err != nil {
		t.Fatalf("UpdateConfig failed: %v", err)
	}

	// Load and verify
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if !cfg.Verbose {
		t.Error("Verbose should be true after update")
	}

	if cfg.Language != "es" {
		t.Errorf("Expected language 'es', got '%s'", cfg.Language)
	}
}

func TestGetAndSetConfig(t *testing.T) {
	// Create test config
	testConfig := &config.Config{
		DotfilesDir:  "/test/dotfiles",
		SymlinksFile: "/test/symlinks.yml",
		OS:           "darwin",
		Language:     "en",
		Verbose:      false,
		DryRun:       true,
	}

	// Set config
	config.SetConfig(testConfig)

	// Get config
	cfg, err := config.GetConfig()
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}

	if cfg != testConfig {
		t.Error("GetConfig should return the same config instance set by SetConfig")
	}

	if cfg.OS != "darwin" {
		t.Errorf("Expected OS 'darwin', got '%s'", cfg.OS)
	}

	if !cfg.DryRun {
		t.Error("DryRun should be true")
	}
}

func TestSaveConfigNil(t *testing.T) {
	err := config.SaveConfig(nil)
	if err == nil {
		t.Error("SaveConfig should fail with nil config")
	}

	expectedError := "config cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestConfigPathGeneration(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "sokru-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override the config directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	configPath, err := config.GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}

	expectedPath := filepath.Join(tempDir, ".config", "sokru", "config.yaml")
	if configPath != expectedPath {
		t.Errorf("Expected config path '%s', got '%s'", expectedPath, configPath)
	}
}
