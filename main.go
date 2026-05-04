package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	entries, err := readDir(cfg.root)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	for i := range len(entries) {
		entry := entries[i]

		if entry.isDir {
			continue
		}

		dir, ok := knownFormatsToType[entry.ext]

		if !ok {
			continue
		}

		pathToDir := filepath.Join(cfg.root, dir)

		if exists, _ := dirExists(filepath.Join(cfg.root, dir)); !exists {
			os.Mkdir(pathToDir, os.FileMode(0755))
		}

		renameErr := os.Rename(entry.absPath, filepath.Join(pathToDir, entry.name))

		if renameErr != nil {
			fmt.Printf(
				"Error creating a file %s in directory %s : --> %s \n",
				entry.name,
				filepath.Join(entry.absPath, dir),
				renameErr.Error(),
			)
			continue
		} else {
			fmt.Printf("Created new file %s in directory %s\n", entry.name, pathToDir)
		}
	}

	fmt.Println("Finished")
}
