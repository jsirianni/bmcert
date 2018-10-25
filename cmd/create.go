package cmd
import (
	"os"
	"fmt"
	"strings"
	"encoding/json"

	"github.com/spf13/cobra"
)


type Certificate struct {
	Common_name string `json:"common_name"`
}


var newcert      Certificate  // Certificate struct
var hostname     string
var outputdir    string
var outputformat string


// NOTE : forces bluemedora.localnet, for now
const fixedDomain string = "bluemedora.localnet"


// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a certificate",
	Run: func(cmd *cobra.Command, args []string) {
		createCertificate()
	},
}


func init() {
	rootCmd.AddCommand(createCmd)

	// set flags
	createCmd.Flags().StringVarP(&hostname, "hostname", "H", "", "The short hostname or FQDN")
	createCmd.Flags().StringVarP(&outputdir, "output-dir", "O", "", "The directory to output to")
	createCmd.Flags().StringVarP(&outputformat, "format", "F", "pem", "The keyfile formant to output. [pem, p12]")

	// require
	createCmd.MarkFlagRequired("hostname")
}


func createCertificate() {
	if parseHostname() != true {
		fmt.Println("Failed to parse hostname: \"" + hostname + "\"" )
		os.Exit(1)
	}
	fmt.Println("Using hostname:", newcert.Common_name)

	payload, err := json.Marshal(newcert)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(string(payload))
	fmt.Println(tls)
}


/*
  Sets the fqdn if hostname argument appears to be valid
   --hostname vault                       // valid
   --hostname vault.bluemedora.localnet   // valid
   --hostname vault.blue                  // !valid

  Returns true if successful, else false
*/
func parseHostname() bool {
	// split hostname argument
	stringSlice := strings.Split(hostname, ".")

	// if hostname is of length zero, return early
	if len(hostname) == 0 {
		fmt.Println("'--hostname' appears to be empty")
		return false
	}

	// if hostname appears to be fqdn
	if len(stringSlice) == 3 {
		// compare domain to fixed domain constant
		d := stringSlice[1] + "." + stringSlice[2]
		if d == fixedDomain {
			newcert.Common_name = hostname
			return true

		// return false if the domain appears to not be bluemedora.localnet
		} else {
			fmt.Println("Domain appears to be malformed, or not equal to", fixedDomain)
			return false
		}

	// if hostname appears to be short
	} else if len(stringSlice) == 1 {
		newcert.Common_name = hostname + "." + fixedDomain
		return true

	// return false if hostname appears to be invalid
	} else {
		fmt.Println("Hostname appears to be neither a short hostname nor a FQDN")
		return false
	}
}
