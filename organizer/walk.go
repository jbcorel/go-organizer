package main

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"
)

type fileEntry struct {
	name    string
	ext     string
	relPath string
	absPath string
	isDir   bool
}

// walks a directory `root`, with or without including other directories`
func readDir(root string) ([]fileEntry, error) {
	entries, err := os.ReadDir(root)
	var fileEntries []fileEntry

	if err != nil {
		slog.Error(err.Error())
		return []fileEntry{}, err
	}

	for _, entry := range entries {
		name := entry.Name()
		ext := filepath.Ext(name)
		relPath := filepath.Join(root, name)
		absPath, _ := filepath.Abs(relPath)
		fileEntries = append(
			fileEntries,
			fileEntry{
				name:    name,
				ext:     ext,
				relPath: relPath,
				absPath: absPath,
				isDir:   entry.IsDir(),
			},
		)
	}

	return fileEntries, nil
}

// TODO change to dest string and move file to dest, not to root
func processFileEntry(entry fileEntry) error {
	if entry.isDir {
		errString := fmt.Sprint(entry.absPath, " is a directory")
		return errors.New(errString)
	}

	dir, ok := knownFormatsToType[entry.ext]

	if !ok {
		errString := fmt.Sprint(entry.ext, " not found in allowed extenstion list")
		return errors.New(errString)
	}

	pathToDir := filepath.Join(Cfg.root, dir.String())

	if err := createDirIfNotExists(pathToDir); err != nil {
		return err
	}

	if err := os.Rename(entry.absPath, filepath.Join(pathToDir, entry.name)); err != nil {
		errString := fmt.Sprint("Cannot move file", entry.name, " : ", err)
		return errors.New(errString)
	}

	fmt.Printf("Created new file %s in directory %s\n", entry.name, pathToDir)

	return nil
}

// destDir - destination of files from path, path - dir to flatten
func flattenDirectory(destDir string, path string) error {
	entries, err := readDir(path)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.isDir && slices.Contains(slices.Collect(maps.Values(folderMap)), entry.name) {
			flattenDirectory(destDir, entry.absPath)
		} else if destDir != path { //its okay since we start with root and wont move to non-media folders anyway. Not the cleanest tho
			if err := os.Rename(entry.absPath, filepath.Join(destDir, entry.name)); err != nil {
				fmt.Println(err)
			}
		}
	}
	os.Remove(path)

	return nil
}
