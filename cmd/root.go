package cmd
import (
	"fmt"
	"os"

	"bmcert/cert"

	"github.com/spf13/cobra"
)

var skipverify bool
var verbose    bool
var hostname     string
var outputdir    string
var outputformat string
var password     string
var altnames     string
var ipsans       string
var urisans      string

var bmcert cert.CertConfig

var rootCmd = &cobra.Command{
	Use:   "bmcert",
	Short: "A CLI for generating certificates with Vault",
	Long: `A CLI for generating certificates with Hashicorp Vault.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		printErrorExit(err, 1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&skipverify, "tls-skip-verify", "", false, "Disable certificate verification when communicating with the Vault API (Defaults to false)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "Enable verbose output --verbose")
}

// assign all command line arguments to the bmcert struct,
// making them accessable by bmcert internal functions
func initConfig() {
	bmcert.AltNames = altnames
	bmcert.Hostname = hostname
	bmcert.IPsans = ipsans
	bmcert.OutputDir = outputdir
	bmcert.OutputFormat = outputformat
	bmcert.Password = password
	bmcert.SkipVerify = skipverify
	bmcert.URISans = urisans
	bmcert.Verbose = verbose
	bmcert.Init()
}

// helper function that prints an error to Stderr and exits
func printErrorExit(err error, code int) {
	fmt.Fprintf(os.Stderr, err.Error())
	os.Exit(code)
}
