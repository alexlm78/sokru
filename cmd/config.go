// Package: cmd
// Description: This file contains the config command for the cli tool. It is used to set up the configuration of sokru.
// (c) 2023 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alexlm78/sokru/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configDotDirCmd)
	configCmd.AddCommand(configSymlinkFileCmd)
	configCmd.AddCommand(configVerboseCmd)
	configCmd.AddCommand(configDryrunCmd)
	configCmd.AddCommand(configOsCmd)
	configCmd.AddCommand(configHelpCmd)
	rootCmd.AddCommand(configCmd)
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Tool to set up the configuration of sokru",
	Long:  `This command will allow you to set up the configuration of sokru.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show all config values by default
		cfg, err := config.GetConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get config: %v\n", err)
			os.Exit(1)
		}

		configPath, _ := config.GetConfigPath()
		fmt.Printf("Configuration file: %s\n\n", configPath)
		fmt.Printf("Current configuration:\n")
		fmt.Printf("  Dotfiles Directory: %s\n", cfg.DotfilesDir)
		fmt.Printf("  Symlinks File:      %s\n", cfg.SymlinksFile)
		fmt.Printf("  Operating System:   %s\n", cfg.OS)
		fmt.Printf("  Verbose:            %v\n", cfg.Verbose)
		fmt.Printf("  Dry Run:            %v\n", cfg.DryRun)
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show all configuration values",
	Long:  `This command will display all current configuration values.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.GetConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get config: %v\n", err)
			os.Exit(1)
		}

		configPath, _ := config.GetConfigPath()
		fmt.Printf("Configuration file: %s\n\n", configPath)
		fmt.Printf("Current configuration:\n")
		fmt.Printf("  Dotfiles Directory: %s\n", cfg.DotfilesDir)
		fmt.Printf("  Symlinks File:      %s\n", cfg.SymlinksFile)
		fmt.Printf("  Operating System:   %s\n", cfg.OS)
		fmt.Printf("  Verbose:            %v\n", cfg.Verbose)
		fmt.Printf("  Dry Run:            %v\n", cfg.DryRun)
	},
}

var configDotDirCmd = &cobra.Command{
	Use:   "dotdir [path]",
	Short: "Set the directory where the dotfiles are located",
	Long:  `This command will allow you to set the directory where the dotfiles are located.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.GetConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get config: %v\n", err)
			os.Exit(1)
		}

		if len(args) == 0 {
			// Display current value
			fmt.Printf("Current dotfiles directory: %s\n", cfg.DotfilesDir)
			return
		}

		// Expand path and update value
		expandedPath := expandPath(args[0])
		err = config.UpdateConfig(func(c *config.Config) {
			c.DotfilesDir = expandedPath
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to update config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Dotfiles directory set to: %s\n", expandedPath)
	},
}

var configSymlinkFileCmd = &cobra.Command{
	Use:   "symlinkfile [path]",
	Short: "Set the file where the symlinks are located",
	Long:  `This command will allow you to set the file where the symlinks are located.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.GetConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get config: %v\n", err)
			os.Exit(1)
		}

		if len(args) == 0 {
			// Display current value
			fmt.Printf("Current symlinks file: %s\n", cfg.SymlinksFile)
			return
		}

		// Expand path and update value
		expandedPath := expandPath(args[0])
		err = config.UpdateConfig(func(c *config.Config) {
			c.SymlinksFile = expandedPath
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to update config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Symlinks file set to: %s\n", expandedPath)
	},
}

var configVerboseCmd = &cobra.Command{
	Use:   "verbose [true|false]",
	Short: "Set the verbosity of the output",
	Long:  `This command will allow you to set the verbosity of the output.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.GetConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get config: %v\n", err)
			os.Exit(1)
		}

		if len(args) == 0 {
			// Display current value
			fmt.Printf("Current verbose setting: %v\n", cfg.Verbose)
			return
		}

		// Parse and update value
		verbose, err := strconv.ParseBool(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid boolean value. Use 'true' or 'false'\n")
			os.Exit(1)
		}

		err = config.UpdateConfig(func(c *config.Config) {
			c.Verbose = verbose
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to update config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Verbose mode set to: %v\n", verbose)
	},
}

var configDryrunCmd = &cobra.Command{
	Use:   "dryrun [true|false]",
	Short: "Set the dry run mode",
	Long:  `This command will allow you to set the dry run mode.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.GetConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get config: %v\n", err)
			os.Exit(1)
		}

		if len(args) == 0 {
			// Display current value
			fmt.Printf("Current dry-run setting: %v\n", cfg.DryRun)
			return
		}

		// Parse and update value
		dryRun, err := strconv.ParseBool(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid boolean value. Use 'true' or 'false'\n")
			os.Exit(1)
		}

		err = config.UpdateConfig(func(c *config.Config) {
			c.DryRun = dryRun
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to update config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Dry-run mode set to: %v\n", dryRun)
	},
}

var configOsCmd = &cobra.Command{
	Use:   "os [operating-system]",
	Short: "Set the operating system",
	Long:  `This command will allow you to set the operating system (e.g., linux, darwin, windows).`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.GetConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get config: %v\n", err)
			os.Exit(1)
		}

		if len(args) == 0 {
			// Display current value
			fmt.Printf("Current OS setting: %s\n", cfg.OS)
			return
		}

		// Validate OS
		osName := strings.ToLower(args[0])
		if !validateOS(osName) {
			fmt.Fprintf(os.Stderr, "Error: Invalid OS '%s'. Valid options are: linux, darwin, windows\n", args[0])
			os.Exit(1)
		}

		// Update value
		err = config.UpdateConfig(func(c *config.Config) {
			c.OS = osName
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to update config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("OS set to: %s\n", osName)
	},
}

var configHelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Help for the config",
	Long:  `This command will show the help for the config command.`,
	Run:   HelpConfigFunc,
}
