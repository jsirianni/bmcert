package auth

import (
    "os"
    "testing"
)

func TestReadVaultToken(t *testing.T) {
    // Set VAULT_TOKEN
    if os.Setenv("VAULT_TOKEN", "fake_token") != nil {
        t.Errorf("Failed to set VAULT_TOKEN env while testing, this should never happen..")
        return // return early as this should not happen
    }

    token, err := ReadVaultToken()
    if err != nil {
        t.Errorf("Expected ReadVaultToken() to NOT return an error, got: " + err.Error())
    }

    if token != "fake_token" {
        t.Errorf("Expected ReadVaultToken() to return 'fake_token', got: " + token)
    }
}
