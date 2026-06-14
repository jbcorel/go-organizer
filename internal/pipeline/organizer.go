package pipeline

import (
	"fmt"
	"organizer/internal/config"
	"organizer/internal/fsutils"
	"os"
	"path/filepath"
)

func move(cfg *config.Config, f FileEntry, dest string) error {
	cfg.Logf("Moving %s -> %s\n", f.path, filepath.Join(dest, f.name))

	if cfg.DryRun {
		return nil
	}

	return os.Rename(f.path, filepath.Join(dest, f.name))
}

func ensureDir(cfg *config.Config, path string) error {
	ex, err := fsutils.DirExists(path)

	if err != nil {
		return err
	}

	if ex {
		return nil
	}

	cfg.Logf("Creating directory %s\n", path)

	if cfg.DryRun {
		return nil
	}

	return fsutils.CreateDir(path)
}

func processDuplicate(cfg *config.Config, f FileEntry, dupPath string) error {
	if err := ensureDir(cfg, dupPath); err != nil {
		return err
	}

	return move(cfg, f, dupPath)
}

func Organize(cfg *config.Config, c <-chan HashResult) {
	dupMap := make(map[[2]string]bool)
	var duplicatePath = filepath.Join(cfg.Root, "duplicates")
	var folderPath string

	for hr := range c {

		dup := dupMap[[2]string{hr.digest, hr.ext}]

		if dup {
			err := processDuplicate(cfg, hr.FileEntry, duplicatePath)
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

		if err := ensureDir(cfg, folderPath); err != nil {
			fmt.Printf("Error creating directory %s: \n%v", folderPath, err)
			continue
		}

		if err := move(cfg, hr.FileEntry, folderPath); err != nil {
			fmt.Printf("Error processing file %s: \n%v", hr.path, err)
			continue
		}

		dupMap[[2]string{hr.digest, hr.ext}] = true
	}
}
