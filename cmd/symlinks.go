// Package cmd
// Description: This file contains the symlinks command for the cli tool. It is used to manage the symlinks in the system.
// (c) 2023 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/alexlm78/sokru/internal/config"
	"github.com/alexlm78/sokru/internal/i18n"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type SymlinkConfig struct {
	OS      string            `yaml:"os,omitempty"`
	Link    map[string]string `yaml:"link"`
	Common  map[string]string `yaml:"common,omitempty"`
	Linux   map[string]string `yaml:"linux,omitempty"`
	Darwin  map[string]string `yaml:"darwin,omitempty"`
	Windows map[string]string `yaml:"windows,omitempty"`
}

// getLinksForOS returns the appropriate links based on the current OS
func (sc *SymlinkConfig) getLinksForOS(currentOS string) map[string]string {
	return sc.GetLinksForOS(currentOS)
}

// GetLinksForOS is exported for testing purposes
func (sc *SymlinkConfig) GetLinksForOS(currentOS string) map[string]string {
	links := make(map[string]string)

	// If this is a legacy format (only "link" field), return it
	if len(sc.Link) > 0 && sc.OS == "" {
		return sc.Link
	}

	// If OS is specified and doesn't match, skip this entry
	if sc.OS != "" && sc.OS != currentOS {
		return links
	}

	// Add common links first (lowest priority)
	for target, source := range sc.Common {
		links[target] = source
	}

	// Add OS-specific links (higher priority, can override common)
	var osLinks map[string]string
	switch currentOS {
	case "linux":
		osLinks = sc.Linux
	case "darwin":
		osLinks = sc.Darwin
	case "windows":
		osLinks = sc.Windows
	}

	for target, source := range osLinks {
		links[target] = source
	}

	// Legacy "link" field has highest priority
	for target, source := range sc.Link {
		links[target] = source
	}

	return links
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
	Run:   UninstallSymlinksFunc,
}

// listCMd represents the list symlink configured (symlinks.yaml)
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the symlinks",
	Long:  `This command will list the symlinks and their status in the system.`,
	Run:   ListSymlinksFunc,
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
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorLoadingConfig, err))
	}

	// Expand path if it contains ~
	symlinkFile := expandPath(cfg.SymlinksFile)

	// Check if symlinks file exists
	if _, err := os.Stat(symlinkFile); os.IsNotExist(err) {
		log.Fatalf("%s", i18n.Error(i18n.MsgSymlinkFileNotFound, symlinkFile))
	}

	// Verbose output
	if cfg.Verbose {
		fmt.Println(i18n.Info(i18n.MsgReadingSymlinksFrom, symlinkFile))
	}

	// Read the YAML file
	data, err := os.ReadFile(symlinkFile)
	if err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorReadingFile, err))
	}

	// Unmarshal YAML to struct
	var symlinkConfigs []SymlinkConfig
	err = yaml.Unmarshal(data, &symlinkConfigs)
	if err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorParsingYAML, err))
	}

	if cfg.Verbose {
		fmt.Println(i18n.Info(i18n.MsgFoundConfigurations, len(symlinkConfigs)))
	}

	// Iterate over items and create symbolic links
	for _, entry := range symlinkConfigs {
		// Get links for current OS
		links := entry.getLinksForOS(cfg.OS)

		for target, source := range links {
			targetPath := expandPath(target)
			sourcePath := expandPath(source)

			// Check if dry-run mode is enabled
			if cfg.DryRun {
				fmt.Println(i18n.Info(i18n.MsgDryRunWouldCreate, targetPath, sourcePath))
				continue
			}

			existingLink, err := os.Readlink(targetPath)
			if err == nil {
				if existingLink == sourcePath {
					if cfg.Verbose {
						fmt.Println(i18n.Success(i18n.MsgSymlinkAlreadyExists, targetPath, sourcePath))
					}
					continue
				} else {
					// Remove existing link if it points to a different destination
					err = os.Remove(targetPath)
					if err != nil {
						log.Printf("%s", i18n.Error(i18n.MsgErrorRemovingSymlink, targetPath, err))
						continue
					}
					if cfg.Verbose {
						fmt.Println(i18n.Success(i18n.MsgExistingSymlinkRemoved, targetPath))
					}
				}
			} else if !os.IsNotExist(err) {
				log.Printf("%s", i18n.Error(i18n.MsgErrorCheckingFile, targetPath, err))
				continue
			}

			// Create the new symbolic link
			err = os.Symlink(sourcePath, targetPath)
			if err != nil {
				log.Printf("%s", i18n.Error(i18n.MsgErrorCreatingSymlink, targetPath, sourcePath, err))
			} else {
				fmt.Println(i18n.Success(i18n.MsgSymlinkCreated, targetPath, sourcePath))
			}
		}
	}
}

func UninstallSymlinksFunc(*cobra.Command, []string) {
	// Get configuration
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorLoadingConfig, err))
	}

	// Expand path if it contains ~
	symlinkFile := expandPath(cfg.SymlinksFile)

	// Check if symlinks file exists
	if _, err := os.Stat(symlinkFile); os.IsNotExist(err) {
		log.Fatalf("%s", i18n.Error(i18n.MsgSymlinkFileNotFound, symlinkFile))
	}

	// Verbose output
	if cfg.Verbose {
		fmt.Println(i18n.Info(i18n.MsgReadingSymlinksFrom, symlinkFile))
	}

	// Read the YAML file
	data, err := os.ReadFile(symlinkFile)
	if err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorReadingFile, err))
	}

	// Unmarshal YAML to struct
	var symlinkConfigs []SymlinkConfig
	err = yaml.Unmarshal(data, &symlinkConfigs)
	if err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorParsingYAML, err))
	}

	if cfg.Verbose {
		fmt.Println(i18n.Info(i18n.MsgFoundConfigurations, len(symlinkConfigs)))
	}

	// Counters for summary
	var removed, skipped, notFound, notSymlink int

	// Iterate over items and remove symbolic links
	for _, entry := range symlinkConfigs {
		// Get links for current OS
		links := entry.getLinksForOS(cfg.OS)

		for target, source := range links {
			targetPath := expandPath(target)
			sourcePath := expandPath(source)

			// Check if target exists
			fileInfo, err := os.Lstat(targetPath)
			if os.IsNotExist(err) {
				if cfg.Verbose {
					fmt.Println(i18n.Info(i18n.MsgSymlinkNotFound, targetPath))
				}
				notFound++
				continue
			}
			if err != nil {
				log.Printf("%s", i18n.Error(i18n.MsgErrorCheckingFile, targetPath, err))
				skipped++
				continue
			}

			// Check if it's a symlink
			if fileInfo.Mode()&os.ModeSymlink == 0 {
				log.Printf("%s", i18n.Warning(i18n.MsgNotSymlink, targetPath))
				notSymlink++
				continue
			}

			// Read the symlink to verify it points to the expected source
			existingLink, err := os.Readlink(targetPath)
			if err != nil {
				log.Printf("%s", i18n.Error(i18n.MsgErrorReadingSymlink, targetPath, err))
				skipped++
				continue
			}

			// Verify the symlink points to the expected source
			if existingLink != sourcePath {
				if cfg.Verbose {
					fmt.Println(i18n.Warning(i18n.MsgSymlinkWrongTarget, targetPath, existingLink, sourcePath))
				}
				skipped++
				continue
			}

			// Check if dry-run mode is enabled
			if cfg.DryRun {
				fmt.Println(i18n.Info(i18n.MsgDryRunWouldRemove, targetPath, sourcePath))
				removed++
				continue
			}

			// Remove the symlink
			err = os.Remove(targetPath)
			if err != nil {
				log.Printf("%s", i18n.Error(i18n.MsgErrorRemovingSymlink, targetPath, err))
				skipped++
			} else {
				fmt.Println(i18n.Success(i18n.MsgSymlinkRemoved, targetPath, sourcePath))
				removed++
			}
		}
	}

	// Print summary
	fmt.Printf("\n=== %s ===\n", i18n.T(i18n.MsgUninstallSummary))
	if cfg.DryRun {
		fmt.Println(i18n.Info(i18n.MsgWouldRemove, removed))
	} else {
		fmt.Println(i18n.Success(i18n.MsgRemoved, removed))
	}
	if notFound > 0 {
		fmt.Println(i18n.Info(i18n.MsgNotFound, notFound))
	}
	if notSymlink > 0 {
		fmt.Println(i18n.Warning(i18n.MsgNotSymlinks, notSymlink))
	}
	if skipped > 0 {
		fmt.Println(i18n.Warning(i18n.MsgSkipped, skipped))
	}
}

func ListSymlinksFunc(*cobra.Command, []string) {
	// Get configuration
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorLoadingConfig, err))
	}

	// Expand path if it contains ~
	symlinkFile := expandPath(cfg.SymlinksFile)

	// Check if symlinks file exists
	if _, err := os.Stat(symlinkFile); os.IsNotExist(err) {
		log.Fatalf("%s", i18n.Error(i18n.MsgSymlinkFileNotFound, symlinkFile))
	}

	// Verbose output
	if cfg.Verbose {
		fmt.Printf("%s\n\n", i18n.Info(i18n.MsgReadingSymlinksFrom, symlinkFile))
	}

	// Read the YAML file
	data, err := os.ReadFile(symlinkFile)
	if err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorReadingFile, err))
	}

	// Unmarshal YAML to struct
	var symlinkConfigs []SymlinkConfig
	err = yaml.Unmarshal(data, &symlinkConfigs)
	if err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorParsingYAML, err))
	}

	// Print header
	fmt.Println(i18n.T(i18n.MsgSymlinksStatus))
	fmt.Println("================")
	fmt.Println()
	fmt.Printf("%-8s %-40s -> %s\n", i18n.T(i18n.MsgStatus), i18n.T(i18n.MsgTarget), i18n.T(i18n.MsgSource))
	fmt.Println("─────────────────────────────────────────────────────────────────────────────────────────────")

	// Counters for summary
	var installed, wrongTarget, notInstalled, regularFile int

	// Iterate over items and check status
	for _, entry := range symlinkConfigs {
		// Get links for current OS
		links := entry.getLinksForOS(cfg.OS)

		for target, source := range links {
			targetPath := expandPath(target)
			sourcePath := expandPath(source)

			// Check if target exists
			fileInfo, err := os.Lstat(targetPath)

			if os.IsNotExist(err) {
				// Not installed
				fmt.Printf("%-8s %-40s -> %s\n", "❌", targetPath, sourcePath)
				notInstalled++
				continue
			}

			if err != nil {
				log.Printf("%s", i18n.Error(i18n.MsgErrorCheckingFile, targetPath, err))
				continue
			}

			// Check if it's a symlink
			if fileInfo.Mode()&os.ModeSymlink == 0 {
				// Regular file exists at target location
				fmt.Printf("%-8s %-40s -> %s\n", "⛔", targetPath, sourcePath)
				regularFile++
				continue
			}

			// Read the symlink
			existingLink, err := os.Readlink(targetPath)
			if err != nil {
				log.Printf("%s", i18n.Error(i18n.MsgErrorReadingSymlink, targetPath, err))
				continue
			}

			// Check if it points to the correct source
			if existingLink == sourcePath {
				// Installed correctly
				fmt.Printf("%-8s %-40s -> %s\n", "✅", targetPath, sourcePath)
				installed++
			} else {
				// Installed but wrong target
				fmt.Printf("%-8s %-40s -> %s\n", "⚠️", targetPath, sourcePath)
				if cfg.Verbose {
					fmt.Printf("         (currently points to: %s)\n", existingLink)
				}
				wrongTarget++
			}
		}
	}

	// Print summary
	fmt.Println()
	fmt.Println(i18n.T(i18n.MsgSummary))
	fmt.Println("--------")
	fmt.Printf("✅ %s\n", i18n.T(i18n.MsgInstalledCorrectly, installed))
	if wrongTarget > 0 {
		fmt.Printf("⚠️  %s\n", i18n.T(i18n.MsgWrongTarget, wrongTarget))
	}
	if notInstalled > 0 {
		fmt.Printf("❌ %s\n", i18n.T(i18n.MsgNotInstalled, notInstalled))
	}
	if regularFile > 0 {
		fmt.Printf("⛔ %s\n", i18n.T(i18n.MsgRegularFileExists, regularFile))
	}
	fmt.Printf("\n%s\n", i18n.T(i18n.MsgTotalSymlinks, installed+wrongTarget+notInstalled+regularFile))

	// Show legend
	fmt.Printf("\n%s\n", i18n.T(i18n.MsgLegend))
	fmt.Printf("  %s\n", i18n.T(i18n.MsgLegendInstalled))
	fmt.Printf("  %s\n", i18n.T(i18n.MsgLegendWrongTarget))
	fmt.Printf("  %s\n", i18n.T(i18n.MsgLegendNotInstalled))
	fmt.Printf("  %s\n", i18n.T(i18n.MsgLegendRegularFile))
}

func init() {
	symlinksCmd.AddCommand(installCmd)
	symlinksCmd.AddCommand(uninstallCmd)
	symlinksCmd.AddCommand(listCmd)
	symlinksCmd.AddCommand(symhelpCmd)
	rootCmd.AddCommand(symlinksCmd)
}
