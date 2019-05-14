package file

import (
    "os"
    "io/ioutil"
    "errors"
)

// WriteFile will write a file to disk. An error is returned
// if the file alredy exists
func WriteFile(filePath string, data []byte, perm os.FileMode) error {
    if FileExists(filePath) == true {
        return errors.New(filePath + " already exists.")
    }
    return ioutil.WriteFile(filePath, data, perm)
}

// return true if the file exists
func FileExists(filePath string) bool {
    if _, err := os.Stat(filePath); err != nil {
        if os.IsNotExist(err) == true {
            return false
        }
    }
    return true
}
