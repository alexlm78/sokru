// Package cmd
// Description: This file contains utility functions shared across commands.
// (c) 2024 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package cmd

import (
	"os"
	"path/filepath"
	"strings"
)

// expandPath expands ~ to the user's home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

// validateOS checks if the OS is one of the supported values
func validateOS(osName string) bool {
	validOS := map[string]bool{
		"linux":   true,
		"darwin":  true,
		"windows": true,
	}
	return validOS[strings.ToLower(osName)]
}
