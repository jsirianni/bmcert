package httpclient

import (
    "os"
    "fmt"
    "net/http"
    "io/ioutil"
    "bytes"
    "crypto/tls"
    "encoding/json"
)

// struct is printed as json when debug is passed to
// PerformRequest
type debug struct {
    Host   string
    Path   string
    Method string
    SkipVerify bool
}

// ConfigureCertVerification disables tls verification if true
// is passed
func ConfigureCertVerification(skipVerify bool) {
    if skipVerify == true {
        fmt.Println("Warning: TLS verification disabled due to flag '--tls-skip-verify'")
        http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    }
}

// CreateRequest returns an http request with headers
func CreateRequest(method string, uri string, payload []byte, auth string) (*http.Request, error) {
    req, err := http.NewRequest(method, uri, bytes.NewBuffer(payload))
    if err != nil {
        return req, err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Authorization", "Basic "+auth)
    return req, nil
}

// HttpClientRequest performs an HTTP request
// Pass a http.Request object
// Returns a response body and a status code ([]byte, int)
func PerformRequest(req *http.Request, verbose bool) ([]byte, int, error) {
    if verbose == true {
        if err := printRequest(req); err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
        }
    }

    //client, err := httpClient()
    //if err != nil {
    //    return nil, 0, err
    //}
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

    return body, resp.StatusCode, nil
}

// ConfigureTLS sets the root certificate path to include all
// paths on a Darwin system, which is not done by default
/*func httpClient() (*http.Client, error) {
    tlsConfig := &tls.Config{}
    c := cleanhttp.DefaultClient()
    t := cleanhttp.DefaultTransport()
    t.TLSClientConfig = tlsConfig
    c.Transport = t
    return c, nil
}*/


func printRequest(req *http.Request) error {
    var x debug
    x.Host = req.Host
    x.Path = req.URL.RequestURI()
    x.Method = req.Method
    x.SkipVerify = http.DefaultTransport.(*http.Transport).TLSClientConfig.InsecureSkipVerify

    d, err := json.MarshalIndent(x, "", "\t")
    if err != nil {
        return err
    }
    fmt.Printf("%s\n", d)
    return nil
}
