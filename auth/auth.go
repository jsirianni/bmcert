/*
bmcert/auth manages retrieving authentication tokens from a
vault server. bmcert/auth relies upon the following environment
variables in order to function:

VAULT_ADDR          // https://vault.mynet.com:8200
VAULT_GITHUB_TOKEN  // required for github auth
*/
package auth
import (
	"fmt"
	"os"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
)


type GithubAuthResp struct {
	Request_id string      `json:"request_id"`
	Lease_id   string      `json:"lease_id"`
	Renewable  bool        `json:"renewable"`
	Lease_duration float32 `json:"lease_duration"`
	Data       string 	   `json:"data"`
	Wrap_info  string      `json:"wrap_info"`
	Warnings   string      `json:"warnings"`
	Auth       GithubAuthData `json:"auth"`
}


type GithubAuthData struct {
	Client_token   string   `json:"client_token"`
	Accessor       string   `json:"accessor"`
	Policies       []string `json:"policies"`
	Token_policies []string `json:"token_policies"`
	Metadata       GithubMetadata `json:"metadata"`
	Lease_duration float32  `json:"lease_duration"`
	Renewable      bool     `json:"renewable"`
	Entity_id      string   `json:"entity_id"`
}


type GithubMetadata struct {
    Org      string `json:"org"`
    Username string `json:"username"`
}


// returns the full URL for github auth endpoint
// relies upon VAULT_ADDR environment variable
func GetVaultAuthUrl() string {
	url := os.Getenv("VAULT_ADDR")
	if len(url) == 0 {
		fmt.Println("Could not read environment VAULT_ADDR")
		os.Exit(1)
	}

    // vault addr includes the protocol and port
    // https://vault.mynetwork.net:8200
	return url + "/v1/auth"
}


// authenticates against vault with a github token
// returns github token as a string
func GithubAuth() string {
    var authurl string = GetVaultAuthUrl() + "/github/login"


    // get the github token from the environment
	token := os.Getenv("VAULT_GITHUB_TOKEN")
    if len(token) == 0 {
        fmt.Println("Could not read environment VAULT_GITHUB_TOKEN")
		os.Exit(1)
    }


	// create the http request
	req, err := http.NewRequest("POST", authurl, bytes.NewBuffer([]byte("{\"token\": \"" + token + "\"}")))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
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


	// handle the http response, be silent if auth success
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else if resp.StatusCode != 200 {
        fmt.Println("Vault response status:", resp.Status)
		fmt.Println("Github auth response Body:", string(body))
		os.Exit(1)
	}


	// unmarshal the response and return the token
	var r GithubAuthResp
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return r.Auth.Client_token
}
