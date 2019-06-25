package cmd
import (
	"fmt"
	"os"

	"bmcert/cert"

	"github.com/spf13/cobra"
)

// boolean command line flags
var skipVerify, verbose, force bool

// string command line flags
var (
	hostname     string
	outputDir    string
	outputFormat string
	password     string
	altNames     string
	ipSans       string
	uriSans      string
)

// bmcert is the certificate configuration
var bmcert cert.Cert

var rootCmd = &cobra.Command{
	Use:   "bmcert",
	Short: "A CLI for generating certificates with Vault",
	Long: `A CLI for generating certificates with Hashicorp Vault.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		printErrorExit(err, 1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&skipVerify, "tls-skip-verify", "", false, "Disable certificate verification when communicating with the Vault API (Defaults to false)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "Enable verbose output --verbose")
}

// assign all command line arguments to the bmcert struct,
// making them accessable by bmcert internal functions
func initConfig() {
	bmcert.AltNames = altNames
	bmcert.Hostname = hostname
	bmcert.IPsans = ipSans
	bmcert.OutputDir = outputDir
	bmcert.OutputFormat = outputFormat
	bmcert.Password = password
	bmcert.SkipVerify = skipVerify
	bmcert.URISans = uriSans
	bmcert.Verbose = verbose
	bmcert.OverWrite = force
	bmcert.Init()
}

// helper function that prints an error to Stderr and exits
func printErrorExit(err error, code int) {
	fmt.Fprintf(os.Stderr, err.Error())
	os.Exit(code)
}
