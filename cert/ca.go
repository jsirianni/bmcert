package cert

import (
    "fmt"

    "github.com/BlueMedoraPublic/bmcert/util/env"
    "github.com/BlueMedoraPublic/bmcert/util/vaultauth"
    "github.com/BlueMedoraPublic/bmcert/util/httpclient"
    "github.com/BlueMedoraPublic/bmcert/util/file"
)

const envPKIURL = "VAULT_PKI_URL"

// CA returns the configured certificate authority
// certificate in PEM format
func (config *Cert) CA() ([]byte, error) {
    pkiURL, err := env.Read(envPKIURL)
    if err != nil {
        return nil, err
    }

    token, err := vaultauth.ReadVaultToken()
    if err != nil {
        return nil, err
    }

    url := pkiURL + "/ca/pem"
    body, err := httpclient.Request("GET", url, nil, token)
    if err != nil {
        return nil, checkAPIError(body, err)
    }

    return body, nil
}

// WriteCA calls CA() and then writes ca.crt to disk
func (config *Cert) WriteCA() error {
    ca, err := config.CA()
    if err != nil {
        return err
    }

    // get the output
    filePath := config.getDir() + "ca.crt"

    // WriteFile takes a file path, []byte payload, os permissions
    // and overwrite boolean. It returns an error if the write
    // failse
    if err := file.WriteFile(filePath, ca, 0600, config.OverWrite); err != nil {
        return err
    }

    fmt.Println("ca written to " + filePath)
    return nil
}
