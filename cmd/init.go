// Package cmd
// Description: This file contains the init command for the cli tool. It is used to initialize sokru for the first time.
// (c) 2023 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package cmd

import (
	"fmt"
	"os"

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
		fmt.Println("Initiliazing sok...")
		homeDir, _ := os.UserHomeDir()
		sokruDir := homeDir + "/.sokru"
		sokruFile := sokruDir + "/config.yaml"

		if _, err := os.Stat(sokruDir); os.IsNotExist(err) {
			os.Mkdir(sokruDir, 0755)
			os.Create(sokruFile)
		} else {
			// check if config.yaml exists if not create it
			if _, err := os.Stat(sokruFile); os.IsNotExist(err) {
				os.Create(sokruFile)
			}
		}
	},
}

func init() {
	initCmd.Flags().BoolP("dry-run", "d", true, "Dry run mode...")

	rootCmd.AddCommand(initCmd)
}
