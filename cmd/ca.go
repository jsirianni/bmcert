package cmd
import (
	"os"
	"fmt"

	"github.com/spf13/cobra"
)

var caCmd = &cobra.Command{
	Use:   "ca",
	Short: "Get the certificate authority",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ca(); err != nil {
            fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
        }
	},
}

func init() {
	rootCmd.AddCommand(caCmd)

    caCmd.Flags().StringVarP(&outputDir, "output-dir", "O", "", "The directory to output to. Defaults to working directory.")
    caCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite if the file already exists")
}

func ca() error {
	return bmcert.WriteCA()
}
