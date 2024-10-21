// Package cmd
// Description: This file contains the version command for the cli tool. It is used to print the version of the tool.
// (c) 2023 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Sok",
	Long:  `All software has versions. This is Sok's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sok v0.1 -- HEAD")
	},
}
