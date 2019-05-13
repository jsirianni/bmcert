package httpclient

import (
    "strconv"
    "net/http"
	"testing"
)

func TestConfigureCertVerification(t *testing.T) {
    // test disabling verification
    ConfigureCertVerification(true)
    if http.DefaultTransport.(*http.Transport).TLSClientConfig.InsecureSkipVerify != true {
        t.Errorf("InsecureSkipVerify was not enabled when it should have been.")
    }
}

func TestCreateRequest(t *testing.T) {
    x, err := CreateRequest("GET", "https://google.com", nil, "")
    if err != nil {
        t.Errorf(err.Error())

    } else if x.Header.Get("Content-Type") != "application/json" {
        t.Errorf("Content-Type header was not set to application/json")

    } else if x.Header.Get("Content-Type") != "application/json" {
        t.Errorf("Accept header was not set to application/json")

    } else if len(x.Header.Get("Authorization")) == 0 {
        t.Errorf("Authorization header was not set to application/json")
    }
}

func TestPerformRequest(t *testing.T) {
    x, err := CreateRequest("GET", "https://google.com", nil, "")
    if err != nil {
        t.Errorf(err.Error())
    }

    resp, status, err := PerformRequest(x, true)
    if err != nil {
        t.Errorf(err.Error())

    } else if status != 200 {
        t.Errorf("Expected status 200, got" + strconv.Itoa(status))

    } else if len(resp) == 0 {
        t.Errorf("Response length was 0")
    }

}


// not a great test, but shows that the function should
// not return an error
/*func TesthttpClient(t *testing.T) {
    _, err := httpClient()
    if err != nil {
        t.Errorf(err.Error())
    }
}*/
