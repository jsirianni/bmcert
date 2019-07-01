package cert

import (
    "testing"
)

func TestValidOutputFormats(t *testing.T) {
    x := ValidOutputFormats()

    if len(x) != 4 {
        t.Errorf("Expected ValidOutputFormats() to return only 4 items")
    }

    found := false
    for _, f := range []string{"pem", "cert", "p12", "pkcs12"} {
        for _, format := range x {
            if f == format {
                found = true
                break
            }
        }

        // if format not found, set error
        if found != true {
            t.Errorf("expected to find " + f + " in ValidOutputFormats()")
        // if format found, reset to false and try next format
        } else {
            found = false
        }
    }
}

func TestIsValidOutputFormat(t *testing.T) {
    for _, f := range []string{"pem", "cert", "p12", "pkcs12"} {

        if err := IsValidOutputFormat(f); err != nil {
            t.Errorf("Expected format '" + f + "' to be valid, got: " + err.Error() )
        }
    }
}
