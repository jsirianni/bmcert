package cmd
import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/certutil"
)


type Certificate struct {
	Common_name string `json:"common_name"`
}


var newcert Certificate  // Certificate struct
var hostname string
var domain   string

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	createCmd.PersistentFlags().StringVar(&hostname, "hostname", "", "The short hostname or FQDN")
	createCmd.PersistentFlags().StringVar(&domain, "domain", "bluemedora.localnet", "The domain name of the host, used if FQDN not present")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


func createCertificate() {
	if parseHostname() == true {
		fmt.Println("Using hostname:", newcert.Common_name)
	}

	return
}


// Sets the fqdn if hostname argument appears to be valid
// Returns true if successful, else false
func parseHostname() bool {
	// split hostname argument
	stringSlice := strings.Split(hostname, ".")

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
