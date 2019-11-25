package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VERSION const is used for the version command and
// by the Makefile for determining file names
const VERSION = "1.0.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the bmcert version and then exits",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
