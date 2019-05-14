package httpclient

import (
    "testing"
)

func TestAPIErrorHelper(t *testing.T) {

    err := APIErrorHelper("https://err.com", 500, []byte("some bad error"))
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
