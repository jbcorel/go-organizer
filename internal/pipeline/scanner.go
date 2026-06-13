package pipeline

import (
	"fmt"
	"organizer/internal/config"
	"os"
	"path/filepath"
)

type FileEntry struct {
	path string
	name string
	ext  string
}

func Scan(cfg *config.Config, dir string, c chan<- FileEntry) {
	defer close(c)
	scan(cfg, dir, c)
}

func scan(cfg *config.Config, dir string, c chan<- FileEntry) {
	entries, err := os.ReadDir(dir)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, e := range entries {
		path := filepath.Join(dir, e.Name())
		if e.IsDir() {
			if cfg.Recursive {
				scan(cfg, path, c)
			}
			continue
		}

		ext := filepath.Ext(path)
		_, ex := config.CategoryMap[ext]

		if !ex {
			continue
		}

		c <- FileEntry{
			path: path,
			ext:  ext,
			name: e.Name(),
		}
	}
}
