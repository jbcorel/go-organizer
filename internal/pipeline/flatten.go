package pipeline

import (
	"fmt"
	"organizer/internal/config"
	"os"
	"path/filepath"
)

// "duplicates" is intentionally excluded: flattening it would collide with the
// originals it was deduplicated against.
func isCategoryDir(name string) bool {
	for _, cat := range config.CategoryMap {
		if string(cat) == name {
			return true
		}
	}
	return false
}

// Flatten moves files out of the category folders back into the root and
// removes the emptied folders.
func Flatten(cfg *config.Config) error {
	return flatten(cfg, cfg.Root, cfg.Root)
}

func flatten(cfg *config.Config, dest, dir string) error {
	entries, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, e := range entries {
		path := filepath.Join(dir, e.Name())

		if e.IsDir() {
			if isCategoryDir(e.Name()) {
				if err := flatten(cfg, dest, path); err != nil {
					fmt.Println(err)
				}
			}
			continue
		}

		if dir == dest {
			continue
		}

		if err := move(cfg, FileEntry{path: path, name: e.Name()}, dest); err != nil {
			fmt.Println(err)
		}
	}

	if dir != dest {
		cfg.Logf("Removing %s\n", dir)
		if !cfg.DryRun {
			if err := os.Remove(dir); err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}
