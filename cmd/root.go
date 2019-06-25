package cmd
import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"errors"

	"bmcert/cert"
	"bmcert/util/timecalc"

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
	ttl          string
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
	var err error

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

	if len(ttl) > 0 {
		bmcert.TTL, err = setTTL()
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	bmcert.Init()
}

// helper function that prints an error to Stderr and exits
func printErrorExit(err error, code int) {
	fmt.Fprintf(os.Stderr, err.Error())
	os.Exit(code)
}

// setTTL takes the ttl argument and returns
// number of seconds as a string
func setTTL() (string, error) {
	ttlNumber, err := stripTTLSuffix()
	if err != nil {
		return "", err
	}

	unit := strings.ToLower(ttl[len(ttl)-1:])
	switch unit {
	case "s":
		return strconv.FormatInt(ttlNumber, timecalc.BASE), nil
	case "d":
		return timecalc.SecondsDayString(ttlNumber), nil
	case "m":
		return timecalc.SecondsMonthString(ttlNumber), nil
	}
	return "", errors.New("TTL unit must be seconds, days, or months (s, d, m), got: " + ttl)
}

// stripTTLSuffix returns the time to live value without the
// unit type suffix
// exampple: 600d is returned as 600
func stripTTLSuffix() (int64, error) {
	t := ttl[:len(ttl)-1]

	// if t (ttl without the unit suffix) cannot be converted
	// to an int, return an error
	i, err := strconv.ParseInt(t, timecalc.BASE, 64)
	if err != nil {
		return 0, errors.New("Failed to convert --ttl value to an int64.\n" + err.Error())
	}

	// make sure the user did not pass '0d' or something similar
	if i < 1 {
		return 0, errors.New("TTL value passed is less than 1: " + ttl)
	}
	return i, nil
}
