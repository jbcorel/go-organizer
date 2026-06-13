package fsutils

import (
	"errors"
	"os"
)

func DirExists(path string) (bool, error) {
	fi, err := os.Stat(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	if !fi.IsDir() {
		return false, errors.New("Not a directory")
	}

	return true, nil
}

func CreateDir(path string) error {
	return os.Mkdir(path, 0755)
}
