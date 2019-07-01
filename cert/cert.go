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

    "bmcert/util/vaultauth"
    "bmcert/util/httpclient"
    "bmcert/util/file"
    "bmcert/util/env"

    "github.com/hashicorp/vault/sdk/helper/certutil"
    pkcs12 "github.com/BlueMedoraPublic/go-pkcs12"
)

// Init sets runtime variables such as tls skip verify
func (config *Cert) Init() {
    httpclient.ConfigureCertVerification(config.SkipVerify)
}

// CreateCertificate calls the Vault API and returns a signed certifcate
func (config *Cert) CreateCertificate() (SignedCertificate, error) {
    var apiResponse apiResponse

	if err := config.parseArgs(); err != nil {
        return apiResponse.Data, err
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

func (config *Cert) requestCertificate() (apiResponse, error) {
    var r apiResponse

	payload, err := json.Marshal(config.certificateReq)
	if err != nil {
		return r, nil
	}

    token, err := vaultauth.ReadVaultToken()
	if err != nil {
		return r, err
	}

    url, err := env.Read("VAULT_CERT_URL")
    if err != nil {
        return r, err
    }

    body, err := httpclient.Request("POST", url, payload, token)
    if err != nil {
        return r, checkAPIError(body, err)
    }

	err = json.Unmarshal(body, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}

// WriteCert the certificate to disk
func (config *Cert) WriteCert(c SignedCertificate) error {

	// write a single pem encoded certificate chain
	if config.OutputFormat == "pem" {
		pem := []byte(c.Certificate + "\n" + c.PrivateKey + "\n" + c.IssuingCa)
		pemFile := config.getDir() + config.Hostname + ".pem"
        if err := file.WriteFile(pemFile, pem, 0600, config.OverWrite); err != nil {
            return err
        }

	// write the certificate and private key to seperate files,
	// both pem encoded
	} else if config.OutputFormat == "cert" {
		crt := []byte(c.Certificate + "\n" + c.IssuingCa)
		crtFile := config.getDir() + config.Hostname + ".crt"
        err := file.WriteFile(crtFile, crt, 0600, config.OverWrite)
		if err != nil {
			return err
		}

		key := []byte(c.PrivateKey)
		keyFile := config.getDir() + config.Hostname + ".key"
        err = file.WriteFile(keyFile, key, 0600, config.OverWrite)
		if err != nil {
			return err
		}

	} else if config.OutputFormat == "pkcs12" || config.OutputFormat == "p12" {
		pem, err := certutil.ParsePEMBundle(c.Certificate + "\n" + c.PrivateKey + "\n" + c.IssuingCa)
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

		p12File := config.getDir() + config.Hostname + ".p12"
        err = file.WriteFile(p12File, p12, 0600, config.OverWrite)
		if err != nil {
			return err
		}
	}
	return nil
}

// perform basic checks on the certificate before assuming it is valid
func (certresp apiResponse) validateCertificate(config *Cert) bool {
	valid := true

	if config.Verbose == true {
		fmt.Println("Validating certificate response. . .")
		fmt.Printf("%#v", certresp)
	}

	// make sure len of cert is not 0
	if len(certresp.Data.Certificate) < 500 {
		valid = false
		fmt.Println("Certificate appears to be shorter than 500 characters, something is wrong.")
	}
	if len(certresp.Data.IssuingCa) < 500 {
		valid = false
		fmt.Println("Issuing CA appears to be shorter than 500 characters, something is wrong.")
	}
	if len(certresp.Data.PrivateKey) < 500 {
		valid = false
		fmt.Println("Private key appears to be shorter than 500 characters, something is wrong.")
	}
	return valid
}

// Return the output directory, with a trailing "/"
func (config *Cert) getDir() string {
	if len(config.OutputDir) != 0 {
		if config.OutputDir[len(config.OutputDir)-1:] == "/" {
			return config.OutputDir
		}
		return config.OutputDir + "/"
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
func (config *Cert) parseArgs() error {
	if err := config.setHostname(); err != nil {
		return err
	}

	// set alt, ip and uri sans
    config.certificateReq.AltNames = config.AltNames
    config.certificateReq.IPSans = config.IPsans
    config.certificateReq.URISans = config.URISans
    config.certificateReq.TTL = config.TTL

	return nil
}


// set newcert.Hostname if it is valid
func (config *Cert) setHostname() error {
	// split Hostname argument
	stringSlice := strings.Split(config.Hostname, ".")

	// if Hostname is of length zero, return early
	if len(config.Hostname) == 0 {
		return errors.New("'--Hostname' appears to be empty")
	}

	// if Hostname appears to be fqdn
	if len(stringSlice) == 3 {
		config.certificateReq.CommonName = config.Hostname
		return nil

	// if Hostname appears to be short
	} else if len(stringSlice) == 1 {
		return errors.New("Hostname appears to be a short Hostname. FQDN is required for --Hostname")

	// return false if Hostname appears to be invalid
	} else {
		return errors.New("Hostname appears to be neither a short Hostname nor a FQDN")
	}
}

// checkAPIError returns a nice error if it is known, otherwise
// returns the origonal error with the response body
func checkAPIError(body []byte, err error) error {
    type apiErrorResponse struct {
        Errors []string `json:"errors"`
    }

    var apiError apiErrorResponse

    // if unmarshaling the api response fails, just return
    // the origonal error
    if json.Unmarshal(body, &apiError) != nil {
        return errors.New(err.Error() + "\n" + string(body))
    }

    for _, e := range apiError.Errors {
        if strings.Contains(e, "unknown role") == true {
            return errors.New(err.Error() + "\nMake sure VAULT_CERT_URL is correct.")
        }
    }
    return errors.New(err.Error() + "\n" + string(body))
}

// ValidOutputFormats returns all valid format options
func ValidOutputFormats() []string {
    return []string{"pem", "cert", "p12", "pkcs12"}
}

// IsValidOutputFormat returns nil if a format is valid
func IsValidOutputFormat(format string) error {
    for _, f := range ValidOutputFormats() {
        if f == format {
            return nil
        }
    }
    return errors.New("error: " + format + " is not valid")
}
