package cert

import (
    "os"
    "fmt"
    "errors"
    "strings"
    "strconv"
    "encoding/json"
    "crypto/x509"
    "math/rand"

    "bmcert/util/auth"
    "bmcert/util/httpclient"
    "bmcert/util/file"
    "bmcert/util/env"

    "github.com/hashicorp/vault/sdk/helper/certutil"
    pkcs12 "github.com/BlueMedoraPublic/go-pkcs12"
)

// Init sets runtime variables such as tls skip verify
func (config *CertConfig) Init() {
    httpclient.ConfigureCertVerification(config.SkipVerify)
}

// CreateCertificate calls the Vault API and returns a signed certifcate
func (config *CertConfig) CreateCertificate() (SignedCertificate, error) {
    var apiResponse apiResponse

	if config.parseArgs() != true {
        return apiResponse.Data, errors.New("parseArgs() returned an error.")
	}

	apiResponse, err := config.requestCertificate()
    if err != nil {
        return apiResponse.Data, err
    }

	if apiResponse.validateCertificate(config) == false {
		return apiResponse.Data, errors.New("Error: failed to validate certifcate")
	}

    return apiResponse.Data, nil
}

func (config *CertConfig) requestCertificate() (apiResponse, error) {
    var r apiResponse

	payload, err := json.Marshal(config.certificateReq)
	if err != nil {
		return r, nil
	}

    token, err := auth.ReadVaultToken()
	if err != nil {
		return r, err
	}

    url, err := env.Read("VAULT_CERT_URL")
    if err != nil {
        return r, err
    }

    body, err := httpclient.Request("POST", url, payload, token)
    if err != nil {
        return r, err
    }

	err = json.Unmarshal(body, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}

// Write the certificate to disk
func (config *CertConfig) WriteCert(c SignedCertificate) error {

	// write a single pem encoded certificate chain
	if config.OutputFormat == "pem" {
		pem := []byte(c.Certificate + "\n" + c.Private_key + "\n" + c.Issuing_ca)
		pem_file := config.getDir() + config.Hostname + ".pem"
        if err := file.WriteFile(pem_file, pem, 0400); err != nil {
            return err
        }

	// write the certificate and private key to seperate files,
	// both pem encoded
	} else if config.OutputFormat == "cert" {
		crt := []byte(c.Certificate + "\n" + c.Issuing_ca)
		crt_file := config.getDir() + config.Hostname + ".crt"
        err := file.WriteFile(crt_file, crt, 0400)
		if err != nil {
			return err
		} else {
			key := []byte(c.Private_key)
			key_file := config.getDir() + config.Hostname + ".key"
            err := file.WriteFile(key_file, key, 0400)
			if err != nil {
				return err
			}
		}

	} else if config.OutputFormat == "pkcs12" || config.OutputFormat == "p12" {
		pem, err := certutil.ParsePEMBundle(c.Certificate + "\n" + c.Private_key + "\n" + c.Issuing_ca)
		if err != nil {
			return err
		}

		// NOTE: This will only use the first certificate in the chain,
		// which will likely cause issues if we use an intermediate certificate
		// during the signing. pkcs12 supports up to 10 certificates in the chain
		ca, err := x509.ParseCertificates(pem.CAChain[0].Bytes)
		if err != nil {
			return err
		}

		rand := strings.NewReader(strconv.Itoa(rand.Int()))

		p12, err := pkcs12.Encode(rand, pem.PrivateKey, pem.Certificate, ca, config.Password)
		if err != nil {
			return err
		}

		p12_file := config.getDir() + config.Hostname + ".p12"
        err = file.WriteFile(p12_file, p12, 0400)
		if err != nil {
			return err
		}
	}
	return nil
}

// set newcert.alt_names if it is valid
func (config *CertConfig) setAltNames() bool {
	if len(config.AltNames) > 0 {
		config.certificateReq.Alt_names = config.AltNames
		return true
	} else {
		return false
	}
}


// set newcert.ip_sans if it is valid
func (config *CertConfig) setIPsans() bool {
	if len(config.IPsans) > 0 {
		config.certificateReq.Ip_sans = config.IPsans
		return true
	} else {
		return false
	}
}


// set newcert.uri_sans if it is valid
func (config *CertConfig) setURISans() bool {
	if len(config.URISans) > 0 {
		config.certificateReq.Uri_sans = config.URISans
		return true
	} else {
		return false
	}
}


// perform basic checks on the certificate before assuming it is valid
func (certresp apiResponse) validateCertificate(config *CertConfig) bool {
	var valid bool = true

	if config.Verbose == true {
		fmt.Println("Validating certificate response. . .")
		fmt.Printf("%#v", certresp)
	}

	// make sure len of cert is not 0
	if len(certresp.Data.Certificate) < 500 {
		valid = false
		fmt.Println("Certificate appears to be shorter than 500 characters, something is wrong.")
	}
	if len(certresp.Data.Issuing_ca) < 500 {
		valid = false
		fmt.Println("Issuing CA appears to be shorter than 500 characters, something is wrong.")
	}
	if len(certresp.Data.Private_key) < 500 {
		valid = false
		fmt.Println("Private key appears to be shorter than 500 characters, something is wrong.")
	}
	return valid
}


// Return the output directory, with a trailing "/"
func (config *CertConfig) getDir() string {
	if len(config.OutputDir) != 0 {
		if config.OutputDir[len(config.OutputDir)-1:] == "/" {
			return config.OutputDir
		} else {
			return config.OutputDir + "/"
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Could not find working directory, trying \"./\"")
		return "./"
	}
	return dir + "/"
}

// parseArgs parses passed arguments, and assigns them to "newcert CertificateReq"
// Returns true if successful, else false
func (config *CertConfig) parseArgs() bool {


	if config.setHostname() != true {
		return false
	}

	// set alt, ip and uri sans but ignore the return
	// value, as they do not matter in this context
	config.setAltNames()
	config.setIPsans()
	config.setURISans()

	return true
}


// set newcert.Hostname if it is valid
func (config *CertConfig) setHostname() bool {


	// split Hostname argument
	stringSlice := strings.Split(config.Hostname, ".")

	// if Hostname is of length zero, return early
	if len(config.Hostname) == 0 {
		fmt.Println("'--Hostname' appears to be empty")
		return false
	}

	// if Hostname appears to be fqdn
	if len(stringSlice) == 3 {
		config.certificateReq.Common_name = config.Hostname
		return true

	// if Hostname appears to be short
	} else if len(stringSlice) == 1 {
		fmt.Println("Hostname appears to be a short Hostname. FQDN is required for --Hostname")
		return false

	// return false if Hostname appears to be invalid
	} else {
		fmt.Println("Hostname appears to be neither a short Hostname nor a FQDN")
		return false
	}
}
