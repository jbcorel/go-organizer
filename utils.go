package main

import (
	"errors"
	"os"
)

func dirExists(path string) (bool, error) {
	d, err := os.Stat(path)

	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if !d.IsDir() {
		return false, errors.New("Not a directory")
	}

	return true, err
}

// Not yet supporting non-current-dir path
func rootFromArgs() (string, error) {
	return ".", nil
}
