// Package cmd
// Description: This file contains the init command for the cli tool. It is used to initialize sokru for the first time.
// (c) 2023 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexlm78/sokru/internal/config"
	"github.com/spf13/cobra"
)

const symlinkTemplate = `# Sokru Symlinks Configuration
# Format:
# - link:
#     ~/target/path: ~/.dotfiles/source/file

# Examples (uncomment to use):
# - link:
#     ~/.bashrc: ~/.dotfiles/bash/bashrc
#     ~/.zshrc: ~/.dotfiles/zsh/zshrc
`

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize sokru for the first time",
	Long: `Initialize sokru for the first time.

	This command will create the necessary files and directories to start using sokru.

	This will create the following files and directories:
	- .config/sokru/
	- .config/sokru/config.yaml
	- .dotfiles/
	- .dotfiles/symlinks.yaml
	- .dotfiles/bash/
	- .dotfiles/zsh/
	- .dotfiles/fish/
	- .dotfiles/pwsh/

	You can change the main start directory with the command sok config set dotdir <path>
	`,
	Run: InitFunc,
}

func InitFunc(cmd *cobra.Command, args []string) {
	fmt.Println("Initializing sokru...")
	fmt.Println()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get home directory: %v\n", err)
		os.Exit(1)
	}

	// Paths
	sokruDir := filepath.Join(homeDir, ".config", "sokru")
	configFile := filepath.Join(sokruDir, "config.yaml")
	dotfilesDir := filepath.Join(homeDir, "dotfiles")
	symlinkFile := filepath.Join(dotfilesDir, "symlinks.yml")

	// Track what was created
	var created []string
	var skipped []string

	// 1. Create ~/.config/sokru/ directory
	if _, err := os.Stat(sokruDir); os.IsNotExist(err) {
		if err := os.MkdirAll(sokruDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to create %s: %v\n", sokruDir, err)
			os.Exit(1)
		}
		created = append(created, sokruDir)
		fmt.Printf("✓ Created directory: %s\n", sokruDir)
	} else {
		skipped = append(skipped, sokruDir)
		fmt.Printf("⚠ Directory already exists: %s\n", sokruDir)
	}

	// 2. Create ~/.config/sokru/config.yaml with defaults
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		cfg := config.GetDefaultConfig()
		// Update paths to use the dotfiles directory we're creating
		cfg.DotfilesDir = dotfilesDir
		cfg.SymlinksFile = symlinkFile

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to create config file: %v\n", err)
			os.Exit(1)
		}
		created = append(created, configFile)
		fmt.Printf("✓ Created config file: %s\n", configFile)
	} else {
		skipped = append(skipped, configFile)
		fmt.Printf("⚠ Config file already exists: %s\n", configFile)
	}

	// 3. Create ~/.dotfiles/ directory
	if _, err := os.Stat(dotfilesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dotfilesDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to create %s: %v\n", dotfilesDir, err)
			os.Exit(1)
		}
		created = append(created, dotfilesDir)
		fmt.Printf("✓ Created directory: %s\n", dotfilesDir)
	} else {
		skipped = append(skipped, dotfilesDir)
		fmt.Printf("⚠ Directory already exists: %s\n", dotfilesDir)
	}

	// 4. Create shell subdirectories
	shellDirs := []string{"bash", "zsh", "fish", "pwsh"}
	for _, shell := range shellDirs {
		shellDir := filepath.Join(dotfilesDir, shell)
		if _, err := os.Stat(shellDir); os.IsNotExist(err) {
			if err := os.MkdirAll(shellDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to create %s: %v\n", shellDir, err)
				os.Exit(1)
			}
			created = append(created, shellDir)
			fmt.Printf("✓ Created directory: %s\n", shellDir)
		} else {
			skipped = append(skipped, shellDir)
			fmt.Printf("⚠ Directory already exists: %s\n", shellDir)
		}
	}

	// 5. Create symlinks.yml template
	if _, err := os.Stat(symlinkFile); os.IsNotExist(err) {
		if err := os.WriteFile(symlinkFile, []byte(symlinkTemplate), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to create %s: %v\n", symlinkFile, err)
			os.Exit(1)
		}
		created = append(created, symlinkFile)
		fmt.Printf("✓ Created symlinks file: %s\n", symlinkFile)
	} else {
		skipped = append(skipped, symlinkFile)
		fmt.Printf("⚠ Symlinks file already exists: %s\n", symlinkFile)
	}

	// Print summary
	fmt.Println()
	fmt.Println("=== Initialization Summary ===")
	fmt.Printf("Created: %d file(s)/directory(ies)\n", len(created))
	if len(skipped) > 0 {
		fmt.Printf("Skipped: %d file(s)/directory(ies) (already exist)\n", len(skipped))
	}

	// Print next steps
	fmt.Println()
	fmt.Println("=== Next Steps ===")
	fmt.Println("1. Add your dotfiles to the shell directories:")
	fmt.Printf("   %s\n", filepath.Join(dotfilesDir, "bash/"))
	fmt.Printf("   %s\n", filepath.Join(dotfilesDir, "zsh/"))
	fmt.Printf("   %s\n", filepath.Join(dotfilesDir, "fish/"))
	fmt.Printf("   %s\n", filepath.Join(dotfilesDir, "pwsh/"))
	fmt.Println()
	fmt.Println("2. Edit the symlinks configuration:")
	fmt.Printf("   %s\n", symlinkFile)
	fmt.Println()
	fmt.Println("3. Install your symlinks:")
	fmt.Println("   sok symlinks install")
	fmt.Println()
	fmt.Println("4. View your configuration:")
	fmt.Println("   sok config")
	fmt.Println()
	fmt.Println("✓ Sokru initialized successfully!")
}

func init() {
	rootCmd.AddCommand(initCmd)
}
