package file

import (
    "testing"
)

func TestExists(t *testing.T) {
    if Exists("/etc/hosts") == false {
        t.Errorf("Expected FileExists(\"/etc/hosts\") to return true, got false")
    }
}
