package file

import (
    "os"
    "io/ioutil"
    "errors"
)

// WriteFile will write a file to disk. An error is returned
// if the file alredy exists
func WriteFile(filePath string, data []byte, perm os.FileMode, OverWrite bool) error {
    if Exists(filePath) == true && OverWrite == false {
        return errors.New(filePath + " already exists.")
    }
    return ioutil.WriteFile(filePath, data, perm)
}

// Exists returns true if the file exists
func Exists(filePath string) bool {
    if _, err := os.Stat(filePath); err != nil {
        if os.IsNotExist(err) == true {
            return false
        }
    }
    return true
}
