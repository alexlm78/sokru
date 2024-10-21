/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Tool to set up the configuration of sokru",
	Long:  `This command will allow you to set up the configuration of sokru.`,
	Run:   HelpConfigFunc,
}

var configHelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Help for the config",
	Long:  `This command will show the help for the config command.`,
	Run:   HelpConfigFunc,
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
