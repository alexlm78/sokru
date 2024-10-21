package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
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
	//var Verbose bool
	//var Request string
	//var Headers []string
	//var Data string

	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//rootCmd.Flags().StringArrayVarP(&Headers, "header", "H", []string{}, "Pass custom headers to the server.")
	//rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Prints the details of the response such as protocol, status, and headers.")
	//rootCmd.PersistentFlags().StringVarP(&Request, "request", "X", "GET", "Specifies the request command to use.")
	//rootCmd.PersistentFlags().StringVarP(&Data, "data", "d", "", "Sends the specified data in a POST request to the HTTP server.")
}

func valArguments(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}

	return nil
}
