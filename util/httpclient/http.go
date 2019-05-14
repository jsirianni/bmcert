package httpclient

import (
    "fmt"
    "crypto/tls"
    "net/http"
    "io/ioutil"
    "bytes"
    "errors"
    "strconv"
)

// ConfigureCertVerification allows tls verification to be disabled
func ConfigureCertVerification(skipVerify bool) {
    if skipVerify == true {
        fmt.Println("Warning: TLS verification disabled due to flag '--tls-skip-verify'")
        http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    }
}

// Request returns a response body and status code
func Request(method string, uri string, payload []byte, token string) ([]byte, error) {
    req, err := createRequest(method, uri, payload, token)
    if err != nil {
        return nil, err
    }

    body, status, err := performRequest(req)

    if StatusValid(status) == false {
        return body, APIErrorHelper(uri, status, body)
    }
    return body, err
}

// APIErrorHelper formats an error message
func APIErrorHelper(uri string, status int, respBody []byte) error {
    return errors.New(uri + " returned " + strconv.Itoa(status) + "\n" + string(respBody))
}

// StatusValid takes a status code, returns true if status
// is 200 or 201
func StatusValid(status int) bool {
    switch status {
    case 200:
        return true
    case 201:
        return true
    default:
        return false
    }
}

// CreateRequest returns an http request with headers
func createRequest(method string, uri string, payload []byte, token string) (*http.Request, error) {
    fmt.Println(method, uri, string(payload), token)
    req, err := http.NewRequest(method, uri, bytes.NewBuffer(payload))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    req.Header.Set("X-Vault-Token", token)
    return req, err
}

// PerformRequest performs an HTTP request and returns a
// response body and status code
func performRequest(req *http.Request) ([]byte, int, error) {
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, 0, err
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, 0, err
    }
    defer resp.Body.Close()
    fmt.Println(resp.StatusCode)
    return body, resp.StatusCode, err
}
