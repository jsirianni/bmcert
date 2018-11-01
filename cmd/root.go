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
	rootCmd.PersistentFlags().BoolVarP(&tlsenable, "tls", "", true, "Enable or disable TLS encryption \"--tls=true\" (Defaults to true)")
	rootCmd.PersistentFlags().BoolVarP(&skipverify, "tls-skip-verify", "", false, "Disable certificate verifiction when communicating with the Vault API (Defaults to false)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "Enable verbose output --verbose")
}



// returns the full URL for the VAULT PKI endpoint
// example: https://vault.localnet:8200/v1/pki/issue/myrole
func GetVaultUrl() string {
	url := os.Getenv("VAULT_CERT_URL")
	if len(url) == 0 {
		fmt.Println("Could not read environment VAULT_CERT_URL")
		os.Exit(1)
	}
	return url
}



// returns the vault token
func GetVaultToken() string {
	token := os.Getenv("VAULT_TOKEN")
	if len(token) == 0 {
		fmt.Println("Could not read environment VAULT_TOKEN")
		os.Exit(1)
	}
	return token
}
