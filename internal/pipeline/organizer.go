package pipeline

import (
	"fmt"
	"organizer/internal/config"
	"organizer/internal/fsutils"
	"os"
	"path/filepath"
)

func move(f FileEntry, dest string) error {
	return os.Rename(f.path, filepath.Join(dest, f.name))
}

func processDuplicate(f FileEntry, dupPath string) error {
	ex, err := fsutils.DirExists(dupPath)

	if err != nil {
		return err
	}

	if !ex {
		err := fsutils.CreateDir(dupPath)
		if err != nil {
			return err
		}
	}

	return move(f, dupPath)
}

func Organize(cfg *config.Config, c <-chan HashResult) {
	dupMap := make(map[[2]string]bool)
	var duplicatePath = filepath.Join(cfg.Root, "duplicates")
	var folderPath string

	for hr := range c {

		dup := dupMap[[2]string{hr.digest, hr.ext}]

		if dup {
			err := processDuplicate(hr.FileEntry, duplicatePath)
			if err != nil {
				fmt.Printf("Error processing duplicate: \n\n%v", err)
			}
			continue
		}

		folder, ex := config.CategoryMap[hr.ext]

		if !ex {
			continue
		}

		folderPath = filepath.Join(cfg.Root, string(folder))
		ex, err := fsutils.DirExists(folderPath)

		if err != nil {
			fmt.Printf("Error checking if %s exists: \n%v", folderPath, err)
			continue
		}

		if !ex {
			err = fsutils.CreateDir(folderPath)

			if err != nil {
				fmt.Printf("Error creating directory %s: \n%v", folderPath, err)
				continue
			}
		}

		err = move(hr.FileEntry, folderPath)

		if err != nil {
			fmt.Printf("Error processing file %s: \n%v", hr.path, err)
			continue
		}

		dupMap[[2]string{hr.digest, hr.ext}] = true
	}
}
