package pipeline

import (
	"organizer/internal/config"
	"os"
	"path/filepath"
)


func scan(cfg *config.Config, dir string, c chan<- string, eCh chan<- error) {
	defer close(c)
	
	entries, err := os.ReadDir(dir)

	if err != nil { 
		eCh <- err
		return
	}

	for _, e := range entries {
		path := filepath.Join(dir, e.Name())
		if e.IsDir() {
			if cfg.Recursive {
				scan(cfg, path, c, eCh)
			}
			continue
		}
		c <- path
	}
}