// Package cmd
// Description: This file contains the root command for the cli tool. It is the entry point for the cli tool.
// (c) 2023 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package cmd

import (
	"fmt"
	"os"

	"github.com/alexlm78/sokru/internal/config"
	"github.com/alexlm78/sokru/internal/i18n"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sok",
	Short: "Sok is a cli tool for managing dotfiles on localhost",
	Long:  `Sok is a cli tool for managing dotfiles on localhost. It allows you to configure and manage your dotfiles easily.`,
	Args:  valArguments,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sok (formerly Sokru) is a cli tool for managing dotfiles on localhost")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load config: %v\n", err)
		cfg = config.GetDefaultConfig()
	}
	config.SetConfig(cfg)

	// Initialize i18n with configured language
	if cfg.Language == "es" {
		i18n.SetLanguage(i18n.Spanish)
	} else {
		i18n.SetLanguage(i18n.English)
	}

	// Set up persistent flags
	rootCmd.PersistentFlags().BoolVarP(&cfg.Verbose, "verbose", "v", cfg.Verbose, "Prints the details of the response such as protocol, status, and headers.")
	rootCmd.PersistentFlags().BoolVar(&cfg.DryRun, "dry-run", cfg.DryRun, "Run in dry-run mode without making actual changes.")
}

func valArguments(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}

	return nil
}
