package cmd
import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


var cfgFile   string
var domain    string
var vaulthost string
var vaultport string
var tls       bool


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
	cobra.OnInitialize(initConfig)

	// if no config is passed, initConfig() will set it
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.bmcert.yaml)")

	// global arguments
	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "bluemedora.localnet", "The domain name of the host, used if FQDN not present")
	rootCmd.PersistentFlags().StringVarP(&vaulthost, "vault-host", "v", "vault.bluemedora.localnet", "The vault server" )
	rootCmd.PersistentFlags().StringVarP(&vaultport, "vault-port", "p", "8200", "The vault http port")
	rootCmd.PersistentFlags().BoolVarP(&tls, "tls", "", true, "Enable or disable TLS encryption \"--tls=true\"")
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bmcert" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bmcert")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
