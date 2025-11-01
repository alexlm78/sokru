// Package test
// Description: Unit tests for symlinks functionality
// (c) 2024 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package test

import (
	"reflect"
	"testing"

	"github.com/alexlm78/sokru/cmd"
)

func TestSymlinkConfig_GetLinksForOS(t *testing.T) {
	tests := []struct {
		name     string
		config   cmd.SymlinkConfig
		os       string
		expected map[string]string
	}{
		{
			name: "Legacy format - link field only",
			config: cmd.SymlinkConfig{
				Link: map[string]string{
					"~/.bashrc": "~/.dotfiles/bash/bashrc",
					"~/.vimrc":  "~/.dotfiles/vim/vimrc",
				},
			},
			os: "linux",
			expected: map[string]string{
				"~/.bashrc": "~/.dotfiles/bash/bashrc",
				"~/.vimrc":  "~/.dotfiles/vim/vimrc",
			},
		},
		{
			name: "Common only",
			config: cmd.SymlinkConfig{
				Common: map[string]string{
					"~/.gitconfig": "~/.dotfiles/git/gitconfig",
					"~/.vimrc":     "~/.dotfiles/vim/vimrc",
				},
			},
			os: "linux",
			expected: map[string]string{
				"~/.gitconfig": "~/.dotfiles/git/gitconfig",
				"~/.vimrc":     "~/.dotfiles/vim/vimrc",
			},
		},
		{
			name: "Linux specific only",
			config: cmd.SymlinkConfig{
				Linux: map[string]string{
					"~/.config/i3/config": "~/.dotfiles/i3/config",
				},
			},
			os: "linux",
			expected: map[string]string{
				"~/.config/i3/config": "~/.dotfiles/i3/config",
			},
		},
		{
			name: "Darwin specific only",
			config: cmd.SymlinkConfig{
				Darwin: map[string]string{
					"~/Library/Preferences/app.plist": "~/.dotfiles/macos/app.plist",
				},
			},
			os: "darwin",
			expected: map[string]string{
				"~/Library/Preferences/app.plist": "~/.dotfiles/macos/app.plist",
			},
		},
		{
			name: "Windows specific only",
			config: cmd.SymlinkConfig{
				Windows: map[string]string{
					"~/AppData/config.ini": "~/.dotfiles/windows/config.ini",
				},
			},
			os: "windows",
			expected: map[string]string{
				"~/AppData/config.ini": "~/.dotfiles/windows/config.ini",
			},
		},
		{
			name: "Common + Linux specific (Linux OS)",
			config: cmd.SymlinkConfig{
				Common: map[string]string{
					"~/.gitconfig": "~/.dotfiles/git/gitconfig",
					"~/.vimrc":     "~/.dotfiles/vim/vimrc",
				},
				Linux: map[string]string{
					"~/.config/i3/config": "~/.dotfiles/i3/config",
				},
			},
			os: "linux",
			expected: map[string]string{
				"~/.gitconfig":        "~/.dotfiles/git/gitconfig",
				"~/.vimrc":            "~/.dotfiles/vim/vimrc",
				"~/.config/i3/config": "~/.dotfiles/i3/config",
			},
		},
		{
			name: "Common + Darwin specific (Darwin OS)",
			config: cmd.SymlinkConfig{
				Common: map[string]string{
					"~/.gitconfig": "~/.dotfiles/git/gitconfig",
				},
				Darwin: map[string]string{
					"~/Library/Preferences/app.plist": "~/.dotfiles/macos/app.plist",
				},
			},
			os: "darwin",
			expected: map[string]string{
				"~/.gitconfig":                    "~/.dotfiles/git/gitconfig",
				"~/Library/Preferences/app.plist": "~/.dotfiles/macos/app.plist",
			},
		},
		{
			name: "OS-specific override common",
			config: cmd.SymlinkConfig{
				Common: map[string]string{
					"~/.bashrc": "~/.dotfiles/bash/bashrc.common",
				},
				Linux: map[string]string{
					"~/.bashrc": "~/.dotfiles/bash/bashrc.linux",
				},
			},
			os: "linux",
			expected: map[string]string{
				"~/.bashrc": "~/.dotfiles/bash/bashrc.linux",
			},
		},
		{
			name: "Link field overrides OS-specific",
			config: cmd.SymlinkConfig{
				Link: map[string]string{
					"~/.bashrc": "~/.dotfiles/bash/bashrc.link",
				},
				Common: map[string]string{
					"~/.bashrc": "~/.dotfiles/bash/bashrc.common",
				},
				Linux: map[string]string{
					"~/.bashrc": "~/.dotfiles/bash/bashrc.linux",
				},
			},
			os: "linux",
			expected: map[string]string{
				"~/.bashrc": "~/.dotfiles/bash/bashrc.link",
			},
		},
		{
			name: "Wrong OS - should return empty",
			config: cmd.SymlinkConfig{
				Linux: map[string]string{
					"~/.config/i3/config": "~/.dotfiles/i3/config",
				},
			},
			os:       "darwin",
			expected: map[string]string{},
		},
		{
			name: "OS filter - matching OS",
			config: cmd.SymlinkConfig{
				OS: "linux",
				Link: map[string]string{
					"~/.config/i3/config": "~/.dotfiles/i3/config",
				},
			},
			os: "linux",
			expected: map[string]string{
				"~/.config/i3/config": "~/.dotfiles/i3/config",
			},
		},
		{
			name: "OS filter - non-matching OS",
			config: cmd.SymlinkConfig{
				OS: "linux",
				Link: map[string]string{
					"~/.config/i3/config": "~/.dotfiles/i3/config",
				},
			},
			os:       "darwin",
			expected: map[string]string{},
		},
		{
			name: "Complex: Common + OS-specific with different files",
			config: cmd.SymlinkConfig{
				Common: map[string]string{
					"~/file1": "~/dotfiles/common/file1",
					"~/file2": "~/dotfiles/common/file2",
				},
				Linux: map[string]string{
					"~/file2": "~/dotfiles/linux/file2",
					"~/file3": "~/dotfiles/linux/file3",
				},
			},
			os: "linux",
			expected: map[string]string{
				"~/file1": "~/dotfiles/common/file1",
				"~/file2": "~/dotfiles/linux/file2",
				"~/file3": "~/dotfiles/linux/file3",
			},
		},
		{
			name:     "Empty config",
			config:   cmd.SymlinkConfig{},
			os:       "linux",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetLinksForOS(tt.os)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetLinksForOS() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSymlinkConfig_GetLinksForOS_AllOperatingSystems(t *testing.T) {
	config := cmd.SymlinkConfig{
		Common: map[string]string{
			"~/.gitconfig": "~/.dotfiles/git/gitconfig",
		},
		Linux: map[string]string{
			"~/.config/i3/config": "~/.dotfiles/i3/config",
		},
		Darwin: map[string]string{
			"~/Library/Preferences/app.plist": "~/.dotfiles/macos/app.plist",
		},
		Windows: map[string]string{
			"~/AppData/config.ini": "~/.dotfiles/windows/config.ini",
		},
	}

	// Test Linux
	linuxResult := config.GetLinksForOS("linux")
	if len(linuxResult) != 2 {
		t.Errorf("Linux: expected 2 links, got %d", len(linuxResult))
	}
	if linuxResult["~/.gitconfig"] != "~/.dotfiles/git/gitconfig" {
		t.Error("Linux: common link not found")
	}
	if linuxResult["~/.config/i3/config"] != "~/.dotfiles/i3/config" {
		t.Error("Linux: linux-specific link not found")
	}

	// Test Darwin
	darwinResult := config.GetLinksForOS("darwin")
	if len(darwinResult) != 2 {
		t.Errorf("Darwin: expected 2 links, got %d", len(darwinResult))
	}
	if darwinResult["~/.gitconfig"] != "~/.dotfiles/git/gitconfig" {
		t.Error("Darwin: common link not found")
	}
	if darwinResult["~/Library/Preferences/app.plist"] != "~/.dotfiles/macos/app.plist" {
		t.Error("Darwin: darwin-specific link not found")
	}

	// Test Windows
	windowsResult := config.GetLinksForOS("windows")
	if len(windowsResult) != 2 {
		t.Errorf("Windows: expected 2 links, got %d", len(windowsResult))
	}
	if windowsResult["~/.gitconfig"] != "~/.dotfiles/git/gitconfig" {
		t.Error("Windows: common link not found")
	}
	if windowsResult["~/AppData/config.ini"] != "~/.dotfiles/windows/config.ini" {
		t.Error("Windows: windows-specific link not found")
	}
}
