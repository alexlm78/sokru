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
