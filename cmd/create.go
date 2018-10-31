package cmd
import (
	"os"
	"fmt"
	"strings"
	"strconv"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"crypto/tls"
	"crypto/x509"
	"math/rand"

	"github.com/spf13/cobra"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
	"github.com/hashicorp/vault/helper/certutil"
)


type Request struct {
	Common_name string `json:"common_name"`
	Alt_names   string `json:"alt_names"`
	Ip_sans     string `json:"ip_sans"`
	Uri_sans    string `json:"uri_sans"`
}

type ApiResponse struct {
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


var request      Request  // Certificate struct
var hostname     string
var outputdir    string
var outputformat string
var altnames     string
var ipsans       string
var urisans      string


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
	createCmd.Flags().StringVarP(&hostname, "hostname", "H", "", "The fully qualified hostname.")
	createCmd.Flags().StringVarP(&outputdir, "output-dir", "O", "", "The directory to output to. Defaults to working directory.")
	createCmd.Flags().StringVarP(&outputformat, "format", "F", "pem", "The keyfile formant to output. [pem, p12]")
	createCmd.Flags().StringVarP(&altnames, "alt-names", "", "", "The requested Subject Alternative Names, in a comma-delimited list")
	createCmd.Flags().StringVarP(&ipsans, "ip-sans", "", "", "The requested IP Subject Alternative Names, in a comma-delimited list")
	createCmd.Flags().StringVarP(&urisans, "uri-sans", "", "", "The requested URI Subject Alternative Names, in a comma-delimited list. (ALTHA: Not tested)")

	// require
	createCmd.MarkFlagRequired("hostname")
}


func createCertificate() {
	if parseArgs() != true {
		os.Exit(1)
	}

	var apiresponse ApiResponse = requestCertificate()
	if validateCertificate(apiresponse) == true {
		writeCert(apiresponse.Data)
	} else {
		fmt.Println("Exiting due to certificate validaton failure. . .")
	}
}


func requestCertificate() ApiResponse {
	// NOTE: disable TLS verification for now
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var url string = GetVaultUrl() + pkipath

	// create the json payload
	payload, err := json.Marshal(request)
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
	var r ApiResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return r
}


// Write the certificate to disk
func writeCert(cert SignedCertificate) {

	// write a single pem encoded certificate chain
	if outputformat == "pem" {
		pem := []byte(cert.Certificate + "\n" + cert.Private_key + "\n" + cert.Issuing_ca)
		pem_file := getDir() + hostname + ".pem"
		err := ioutil.WriteFile(pem_file, pem, 0400)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

	// write the certificate and private key to seperate files,
	// both pem encoded
	} else if outputformat == "cert" {
		crt := []byte(cert.Certificate + "\n" + cert.Issuing_ca)
		crt_file := getDir() + hostname + ".crt"
		err := ioutil.WriteFile(crt_file, crt, 0400)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		} else {
			key := []byte(cert.Private_key)
			key_file := getDir() + hostname + ".key"
			err := ioutil.WriteFile(key_file, key, 0400)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

	} else if outputformat == "pkcs12" || outputformat == "p12" {

		pem, err := certutil.ParsePEMBundle(cert.Certificate + "\n" + cert.Private_key + "\n" + cert.Issuing_ca)
		if err != nil {
			fmt.Println("Certutil failed to build PEM")
			os.Exit(1)
		}

		// NOTE: This will only use the first certificate in the chain,
		// which will likely cause issues if we use an intermediate certificate
		// during the signing. pkcs12 supports up to 10 certificates in the chain
		ca, err := x509.ParseCertificates(pem.CAChain[0].Bytes)
		if err != nil {
			fmt.Println("Failed to parse Certificate authority")
		}

		rand := strings.NewReader(strconv.Itoa(rand.Int()))

		p12, err := pkcs12.Encode(rand, pem.PrivateKey, pem.Certificate, ca, "medora")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		p12_file := getDir() + hostname + ".p12"
		err = ioutil.WriteFile(p12_file, p12, 0400)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	return
}



// Parses passed arguments, and assigns them to "newcert CertificateReq"
// Returns true if successful, else false
func parseArgs() bool {
	if setHostname() != true {
		return false
	}

	// set alt, ip and uri sans but ignore the return
	// value, as they do not matter in this context
	setAltNameS()
	setIpSans()
	setUriSans()

	return true
}


// set newcert.hostname if it is valid
func setHostname() bool {
	// split hostname argument
	stringSlice := strings.Split(hostname, ".")

	// if hostname is of length zero, return early
	if len(hostname) == 0 {
		fmt.Println("'--hostname' appears to be empty")
		return false
	}

	// if hostname appears to be fqdn
	if len(stringSlice) == 3 {
		request.Common_name = hostname
		return true

	// if hostname appears to be short
	} else if len(stringSlice) == 1 {
		fmt.Println("Hostname appears to be a short hostname. FQDN is required for --hostname")
		return false

	// return false if hostname appears to be invalid
	} else {
		fmt.Println("Hostname appears to be neither a short hostname nor a FQDN")
		return false
	}
}


// set newcert.alt_names if it is valid
func setAltNameS() bool {
	if len(altnames) > 0 {
		request.Alt_names = altnames
		return true
	} else {
		return false
	}
}


// set newcert.ip_sans if it is valid
func setIpSans() bool {
	if len(ipsans) > 0 {
		request.Ip_sans = ipsans
		return true
	} else {
		return false
	}
}


// set newcert.uri_sans if it is valid
func setUriSans() bool {
	if len(urisans) > 0 {
		request.Uri_sans = urisans
		return true
	} else {
		return false
	}
}


// perform basic checks on the certificate before assuming it is valid
func validateCertificate(certresp ApiResponse) bool {
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
