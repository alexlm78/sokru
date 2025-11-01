// Package test
// Description: Unit tests for utility functions
// (c) 2024 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexlm78/sokru/cmd"
)

func TestExpandPath(t *testing.T) {
	// Get the actual home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Path with tilde",
			input:    "~/dotfiles",
			expected: filepath.Join(homeDir, "dotfiles"),
		},
		{
			name:     "Path with tilde and subdirectory",
			input:    "~/.config/nvim",
			expected: filepath.Join(homeDir, ".config/nvim"),
		},
		{
			name:     "Absolute path without tilde",
			input:    "/usr/local/bin",
			expected: "/usr/local/bin",
		},
		{
			name:     "Relative path without tilde",
			input:    "./local/file",
			expected: "./local/file",
		},
		{
			name:     "Just tilde",
			input:    "~",
			expected: "~",
		},
		{
			name:     "Tilde not at start",
			input:    "/home/user~/file",
			expected: "/home/user~/file",
		},
		{
			name:     "Empty path",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cmd.ExpandPathForTesting(tt.input)
			if result != tt.expected {
				t.Errorf("expandPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateOS(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid OS: linux",
			input:    "linux",
			expected: true,
		},
		{
			name:     "Valid OS: darwin",
			input:    "darwin",
			expected: true,
		},
		{
			name:     "Valid OS: windows",
			input:    "windows",
			expected: true,
		},
		{
			name:     "Valid OS: Linux (uppercase)",
			input:    "Linux",
			expected: true,
		},
		{
			name:     "Valid OS: DARWIN (uppercase)",
			input:    "DARWIN",
			expected: true,
		},
		{
			name:     "Valid OS: Windows (mixed case)",
			input:    "Windows",
			expected: true,
		},
		{
			name:     "Invalid OS: freebsd",
			input:    "freebsd",
			expected: false,
		},
		{
			name:     "Invalid OS: macos",
			input:    "macos",
			expected: false,
		},
		{
			name:     "Invalid OS: empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Invalid OS: random string",
			input:    "invalid",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cmd.ValidateOSForTesting(tt.input)
			if result != tt.expected {
				t.Errorf("validateOS(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExpandPathWithCustomHome(t *testing.T) {
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set custom HOME
	customHome := "/custom/home"
	os.Setenv("HOME", customHome)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Tilde expansion with custom home",
			input:    "~/dotfiles",
			expected: filepath.Join(customHome, "dotfiles"),
		},
		{
			name:     "Tilde with nested path",
			input:    "~/.config/app/settings.json",
			expected: filepath.Join(customHome, ".config/app/settings.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cmd.ExpandPathForTesting(tt.input)
			if result != tt.expected {
				t.Errorf("expandPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
