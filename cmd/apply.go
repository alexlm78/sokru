/*
Copyright © 2024 Alejanjdro Lopez Monzon <alejandro@kreaker.dev>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the changes in memory and reload the symlinks and dotfiles",
	Long:  `This command will apply the changes in memory and reload the symlinks and dotfiles.`,
	Run:   ApplyFunc,
}

func ApplyFunc(cmd *cobra.Command, args []string) {
	fmt.Println("apply called")
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
