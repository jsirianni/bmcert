package env

import (
    "os"
    "testing"
)

func TestRead(t *testing.T) {
    _, err := Read("SOME_FAKE_VAR")
    if err == nil {
        t.Errorf("Expected Read(\"SOME_FAKE_VAR\") to return an error, got nil")
    }

    err = os.Setenv("SOME_REAL_VAR", "value")
    if err != nil {
        t.Errorf("Got an error while setting an env variable, this should never happen")
        return
    }

    x, err := Read("SOME_REAL_VAR")
    if err != nil {
        t.Errorf("Expected Read(\"SOME_REAL_VAR\") to return an error, got nil")
    }

    if x != "value" {
        t.Errorf("Expected Read(\"SOME_REAL_VAR\") to return an 'value', got: " + x)
    }

    if len(x) < 1 {
        t.Errorf("Expected Read(\"SOME_REAL_VAR\") to return an 'value', got: " + x)
    }
}
