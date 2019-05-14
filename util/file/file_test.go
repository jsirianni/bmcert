package file

import (
    "testing"
)

func TestFileExists(t *testing.T) {
    if FileExists("/etc/hosts") == false {
        t.Errorf("Expected FileExists(\"/etc/hosts\") to return true, got false")
    }
}
