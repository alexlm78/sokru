/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// helpCmd represents the help command
var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Help for sokru",
	Long:  `This command will show the help for sokru.`,
	Run:   HelpFunc,
}

var helpSymlinksCmd = &cobra.Command{
	Use:   "symlinks",
	Short: "Help for the symlinks",
	Long:  `This command will show the help for the symlinks command.`,
	Run:   HelpSymlinksFunc,
}

var helpConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Help for the config",
	Long:  `This command will show the help for the config command.`,
	Run:   HelpConfigFunc,
}

func HelpSymlinksFunc(cmd *cobra.Command, args []string) {
	fmt.Println("SymLinks command help::")
	fmt.Println("sok symlinks install        # Install the symlinks")
	fmt.Println("sok symlinks uninstall      # Uninstall the symlinks")
	fmt.Println("sok symlinks list           # List the symlinks")
	fmt.Println("sok symlinks help           # Show this help")
}

func HelpConfigFunc(cmd *cobra.Command, args []string) {
	fmt.Println("Config command help::")
	fmt.Println("sok config dotDir <dir>     # Set the directory where the dotfiles are stored (default: ~/.dotfiles)")
	fmt.Println("sok config symlinks <file>  # Set the file that contains the symlinks (default: ~/.dotfiles/symlinks.yml)")
	fmt.Println("sok config os <os>          # Set the OS to use (default: linux)")
	fmt.Println("sok config verbose <bool>   # Set the verbose mode (default: false)")
	fmt.Println("sok config dryRun <bool>    # Set the dry run mode (default: false)")
	fmt.Println("sok config help             # Show this help")
}

func HelpFunc(cmd *cobra.Command, args []string) {
	fmt.Println("Sokru help::")
	fmt.Println("sok init                    # Initialize the configuration")
	fmt.Println("sok apply                   # Apply the changes in memory and reload the symlinks and dotfiles")
	fmt.Println("sok version                 # Show the version")
	fmt.Println("sok config                  # Show the configuration options")
	fmt.Println("sok symlinks                # Show the symlinks options")
	fmt.Println("sok help                    # Show this help")
}

func init() {
	rootCmd.AddCommand(helpCmd)
	helpCmd.AddCommand(helpSymlinksCmd)
	helpCmd.AddCommand(helpConfigCmd)
}
