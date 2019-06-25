package httpclient

import (
    "net/http"
    "testing"
)

func TestAPIErrorHelper(t *testing.T) {

    req, err := http.NewRequest("POST", "https://err.com", nil)
    if err != nil {
        t.Errorf("Failed to create new HTTP Request, this error should not happen.\n" + err.Error())
    }

    err = APIErrorHelper(req, 500)
    if err == nil {
        t.Errorf("Expected an error when calling APIErrorHelper(), got 'nil'")
        return
    }

    if len(err.Error()) < 1 {
        t.Errorf("Expected APIErrorHelper() to return an error greater than length zero")
    }
}

func TestStatusValid(t *testing.T) {
    // should return false
    if StatusValid(500) == true {
        t.Errorf("Expected StatusValid(500) to return false, got true")
    }

    if StatusValid(400) == true {
        t.Errorf("Expected StatusValid(500) to return false, got true")
    }

    if StatusValid(300) == true {
        t.Errorf("Expected StatusValid(500) to return false, got true")
    }

    if StatusValid(202) == true {
        t.Errorf("Expected StatusValid(202) to return false, got true")
    }

    // should return true
    if StatusValid(201) == false {
        t.Errorf("Expected StatusValid(201) to return true, got false")
    }

    if StatusValid(200) == false {
        t.Errorf("Expected StatusValid(200) to return true, got false")
    }
}

func TestCreateRequest(t *testing.T) {
    req, err := CreateRequest("POST", "https://test.com", []byte("payload"), "token")
    if err != nil {
        t.Errorf("Expected CreateRequest() to NOT return an error, got " + err.Error())
        return
    }

    if req.Header.Get("X-Vault-Token") != "token" {
        t.Errorf("Expected CreateRequest() to return a http request with header X-Vault-Token='token'")
    }

    if req.Header.Get("Content-Type") != "application/json" {
        t.Errorf("Expected CreateRequest() to return a http request with header Content-Type='application/json'")
    }

}
