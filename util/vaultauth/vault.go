package vaultauth

import (
	"io/ioutil"
	"strings"

	"bmcert/util/env"

	"github.com/mitchellh/go-homedir"
)

// ReadVaultToken returns a vault token from the execution
// environment, first checking for VAULT_TOKEN variable,
// and falling back to ~/.vault-token
func ReadVaultToken() (string, error) {
	// env.Read() returns an error if length is < 1,
	// therefor, if err is nil, return x
	x, err := env.Read("VAULT_TOKEN")
	if err == nil {
		return x, nil
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
	return lines[0], nil
}

func getVaultTokenFilePath() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return homeDir + "/.vault-token", nil
}
