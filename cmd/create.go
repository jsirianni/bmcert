package cmd
import (
	"os"
	"fmt"
	"strings"
	"encoding/json"

	"github.com/spf13/cobra"
	//"github.com/hashicorp/vault/api"
	//"github.com/hashicorp/vault/helper/certutil"
)


type Certificate struct {
	Common_name string `json:"common_name"`
}


var newcert Certificate  // Certificate struct
var hostname string

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
	createCmd.Flags().StringVarP(&hostname, "hostname", "H", "", "The short hostname or FQDN")
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


// Sets the fqdn if hostname argument appears to be valid
// Returns true if successful, else false
func parseHostname() bool {
	// split hostname argument
	stringSlice := strings.Split(hostname, ".")


	if len(hostname) == 0 {
		fmt.Println("bruh")

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

	// TODO: Figure out why if / else if / else was not good enough
	fmt.Println("This should never happen, but the compiler made me put a return here..")
	return false
}
