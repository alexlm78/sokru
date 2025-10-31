// Package cmd
// Description: This file contains the symlinks command for the cli tool. It is used to manage the symlinks in the system.
// (c) 2023 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alexlm78/sokru/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type SymlinkConfig struct {
	Link map[string]string `yaml:"link"`
}

// symlinksCmd represents the symlinks command
var symlinksCmd = &cobra.Command{
	Use:   "symlinks",
	Short: "Manage the symlinks",
	Long:  `This command will allow you to manage the symlinks in the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Main sumlinks function called")
	},
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the symlinks",
	Long:  `This command will install the symlinks in the system.`,
	Run:   InstallSymlinksFunc,
}

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the symlinks",
	Long:  `This command will uninstall the symlinks in the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Uninstalling symlinks")
	},
}

// listCMd represents the list symlink configured (symlinks.yaml)
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the symlinks",
	Long:  `This command will list the symlinks in the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing symlinks")
	},
}

// helpCmd represents the help command
var symhelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Help for the symlinks",
	Long:  `This command will show the help for the symlinks command.`,
	Run:   HelpSymlinksFunc,
}

func InstallSymlinksFunc(*cobra.Command, []string) {
	// Get configuration
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Expand path if it contains ~
	symlinkFile := expandPath(cfg.SymlinksFile)

	// Check if symlinks file exists
	if _, err := os.Stat(symlinkFile); os.IsNotExist(err) {
		log.Fatalf("Symlinks file not found: %s\nPlease create the file or update the configuration with: sok config symlinkfile <path>", symlinkFile)
	}

	// Verbose output
	if cfg.Verbose {
		fmt.Printf("Reading symlinks configuration from: %s\n", symlinkFile)
	}

	// Read the YAML file
	data, err := os.ReadFile(symlinkFile)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Unmarshal YAML to struct
	var symlinkConfigs []SymlinkConfig
	err = yaml.Unmarshal(data, &symlinkConfigs)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}

	if cfg.Verbose {
		fmt.Printf("Found %d symlink configuration(s)\n", len(symlinkConfigs))
	}

	// Iterate over items and create symbolic links
	for _, entry := range symlinkConfigs {
		for target, source := range entry.Link {
			targetPath := expandPath(target)
			sourcePath := expandPath(source)

			// Check if dry-run mode is enabled
			if cfg.DryRun {
				fmt.Printf("[DRY-RUN] Would create symlink: %s -> %s\n", targetPath, sourcePath)
				continue
			}

			existingLink, err := os.Readlink(targetPath)
			if err == nil {
				if existingLink == sourcePath {
					if cfg.Verbose {
						fmt.Printf("Symlink already exists and is correct: %s -> %s\n", targetPath, sourcePath)
					}
					continue
				} else {
					// Remove existing link if it points to a different destination
					err = os.Remove(targetPath)
					if err != nil {
						log.Printf("Error removing existing symlink at %s: %v", targetPath, err)
						continue
					}
					if cfg.Verbose {
						fmt.Printf("Existing symlink removed: %s\n", targetPath)
					}
				}
			} else if !os.IsNotExist(err) {
				log.Printf("Error checking symlink at %s: %v", targetPath, err)
				continue
			}

			// Create the new symbolic link
			err = os.Symlink(sourcePath, targetPath)
			if err != nil {
				log.Printf("Error creating symlink from %s to %s: %v", targetPath, sourcePath, err)
			} else {
				fmt.Printf("Symlink created: %s -> %s\n", targetPath, sourcePath)
			}
		}
	}
}

func init() {
	symlinksCmd.AddCommand(installCmd)
	symlinksCmd.AddCommand(uninstallCmd)
	symlinksCmd.AddCommand(listCmd)
	symlinksCmd.AddCommand(symhelpCmd)
	rootCmd.AddCommand(symlinksCmd)
}

// expandPath reemplaza ~ con el directorio HOME
func expandPath(path string) string {
	if path[:2] == "~/" {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, path[2:])
	}
	return path
}
