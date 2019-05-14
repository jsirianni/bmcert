package env

import (
    "os"
    "errors"
)

// Read returns an env variable as a string, will return an
// error if the length of the variable is < 1
func Read(variable string) (string, error) {
    url := os.Getenv(variable)
    if len(url) < 1 {
        return "", errors.New("Could not read environment'" + variable + "'")
    }
    return url, nil
}
