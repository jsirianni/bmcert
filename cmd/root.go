package cmd
import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bmcert",
	Short: "A CLI for generating certificates with Vault",
	Long: `A CLI for generating certificates with Hashicorp Vault.`,
}


// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


func init() {
	// global arguments
	rootCmd.PersistentFlags().StringVarP(&vaulthost, "vault-host", "", "vault.bluemedora.localnet", "The vault server" )
	rootCmd.PersistentFlags().StringVarP(&vaultport, "vault-port", "", "8200", "The vault http port")
	rootCmd.PersistentFlags().StringVarP(&pkipath, "pkipath", "", "/v1/bm-pki-int/issue/bluemedora-dot-localnet", "The vault certificate authority mount point")
	rootCmd.PersistentFlags().BoolVarP(&tlsenable, "tls", "", true, "Enable or disable TLS encryption \"--tls=true\"")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "Enable verbose output --verbose")
}



func GetVaultUrl() string {
	if tlsenable == true {
		return "https://" + vaulthost + ":" + vaultport
	} else {
		return "http://" + vaulthost + ":" + vaultport
	}
}
