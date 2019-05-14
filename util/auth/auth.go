package auth

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// ReadVaultToken returns a vault token from the execution
// environment, first checking for VAULT_TOKEN variable,
// and falling back to ~/.vault-token
func ReadVaultToken() (string, error) {
	if len(os.Getenv("VAULT_TOKEN")) > 0 {
		return os.Getenv("VAULT_TOKEN"), nil
	}

	filePath, err := getVaultTokenFilePath()
	if err != nil {
		return "", err
	}

	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(f), "\n")
	fmt.Println(lines[0])
	return lines[0], nil
}

func getVaultTokenFilePath() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return homeDir + "/.vault-token", nil
}
