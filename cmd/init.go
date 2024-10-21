/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize sokru for the first time",
	Long: `Initialize sokru for the first time.

	This command will create the necessary files and directories to start using sokru.

	This will create the following files and directories:
	- .sokru/
	- .sokru/config.yaml
	- .dotfiles/
	- .dotfiles/symlinks.yaml
	- .dotfiles/bash/
	- .dotfiles/zsh/
	- .dotfiles/fish/
	- .dotfiles/pwsh/

	You can change the main start directory with the command sok config set dotdir <path>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initiliazing sok")
	},
}

func init() {
	initCmd.Flags().BoolP("dry-run", "d", true, "Dry run mode...")

	rootCmd.AddCommand(initCmd)
}
