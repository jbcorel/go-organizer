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

	return true, nil
}

func createDirIfNotExists(path string) error {
	exists, err := dirExists(path)

	if err != nil {
		return err
	}

	if exists {
		return nil
	} else {
		return os.Mkdir(path, os.FileMode(0755))
	}

}
