package main

import (
	"log/slog"
	"os"
	"path/filepath"
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
