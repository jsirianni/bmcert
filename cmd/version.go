package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "0.1.4"

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
