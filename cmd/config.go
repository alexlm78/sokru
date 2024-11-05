// Package: cmd
// Description: This file contains the config command for the cli tool. It is used to set up the configuration of sokru.
// (c) 2023 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
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
	Run:   HelpConfigFunc,
}

var configDotDirCmd = &cobra.Command{
	Use:   "dotdir",
	Short: "Set the directory where the dotfiles are located",
	Long:  `This command will allow you to set the directory where the dotfiles are located.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dotdir command")
	},
}

var configSymlinkFileCmd = &cobra.Command{
	Use:   "symlinkfile",
	Short: "Set the file where the symlinks are located",
	Long:  `This command will allow you to set the file where the symlinks are located.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("symlinkfile command")
	},
}

var configVerboseCmd = &cobra.Command{
	Use:   "verbose",
	Short: "Set the verbosity of the output",
	Long:  `This command will allow you to set the verbosity of the output.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("verbose command")
	},
}

var configDryrunCmd = &cobra.Command{
	Use:   "dryrun",
	Short: "Set the dry run mode",
	Long:  `This command will allow you to set the dry run mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dryrun command")
	},
}

var configOsCmd = &cobra.Command{
	Use:   "os",
	Short: "Set the operating system",
	Long:  `This command will allow you to set the operating system.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("os command")
	},
}

var configHelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Help for the config",
	Long:  `This command will show the help for the config command.`,
	Run:   HelpConfigFunc,
}
