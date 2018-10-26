package cmd
import (
	"os"
	"fmt"
	"strings"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"crypto/tls"

	"github.com/spf13/cobra"
)


type CertificateReq struct {
	Common_name string `json:"common_name"`
}

type CertificateResp struct {
	Request_id string      `json:request_id`
	Lease_id   string      `json:lease_id`    // usually null
	Renewable  bool        `json:renewable`
	Lease_duration float32 `json:lease_duration`
	Data SignedCertificate `json:data`
	Wrap_info  string      `json:wrap_info`   // usually null
	Warnings   string      `json:warnings`    // usually null
	Auth       string      `json:auth`        // usually null
}

type SignedCertificate struct {
	Certificate string      `json:certificate`
	Issuing_ca string       `json:issuing_ca`
	Private_key string 	    `json:private_key`
	Private_key_type string `json:private_key_type`
	Serial_number string    `json:serial_number`
}


var newcert      CertificateReq  // Certificate struct
var hostname     string
var outputdir    string
var outputformat string


// NOTE : forces bluemedora.localnet, for now
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

	// set flags
	createCmd.Flags().StringVarP(&hostname, "hostname", "H", "", "The short hostname or FQDN.")
	createCmd.Flags().StringVarP(&outputdir, "output-dir", "O", "", "The directory to output to. Defaults to working directory.")
	createCmd.Flags().StringVarP(&outputformat, "format", "F", "pem", "The keyfile formant to output. [pem, p12]")

	// require
	createCmd.MarkFlagRequired("hostname")
}


func createCertificate() {
	if parseHostname() != true {
		fmt.Println("Failed to parse hostname: \"" + hostname + "\"" )
		os.Exit(1)
	}


	var certresp CertificateResp = requestCertificate()
	var newcert SignedCertificate = certresp.Data
	if validateCertificate(certresp) == true {
		writeCert(newcert)
	} else {
		fmt.Println("Exiting due to certificate validaton failure. . .")
	}
}


func requestCertificate() CertificateResp {
	// NOTE: disable TLS verification for now
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var url string = GetVaultUrl() + pkipath

	// create the json payload
	payload, err := json.Marshal(newcert)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// create the http request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	req.Header.Set("X-Vault-Token", os.Getenv("VAULT_TOKEN"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")


	// create a http client, and perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	// handle the http response
	fmt.Println("Vault response status:", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else if resp.StatusCode != 200 {
		fmt.Println("Response Body:", string(body))
		os.Exit(1)
	}


	// return the certificate response
	var certresp CertificateResp
	err = json.Unmarshal(body, &certresp)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return certresp
}


// Write the certificate to disk
func writeCert(cert SignedCertificate) {
	pem := []byte(cert.Certificate + "\n" + cert.Private_key + "\n" + cert.Issuing_ca)
	file := getDir() + hostname + ".crt"
	err := ioutil.WriteFile(file, pem, 0400)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return
}


/*
  Sets the fqdn if hostname argument appears to be valid
   --hostname vault                       // valid
   --hostname vault.bluemedora.localnet   // valid
   --hostname vault.blue                  // !valid

  Returns true if successful, else false
*/
func parseHostname() bool {
	// split hostname argument
	stringSlice := strings.Split(hostname, ".")

	// if hostname is of length zero, return early
	if len(hostname) == 0 {
		fmt.Println("'--hostname' appears to be empty")
		return false
	}

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
}


// perform basic checks on the certificate before assuming it is valid
func validateCertificate(certresp CertificateResp) bool {
	var valid bool = true

	if verbose == true {
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
func getDir() string {
	if len(outputdir) != 0 {
		if outputdir[len(outputdir)-1:] == "/" {
			return outputdir
		} else {
			return outputdir + "/"
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Could not find working directory, trying \"./\"")
		return "./"
	}
	return dir + "/"
}
