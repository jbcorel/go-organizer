package pipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"organizer/internal/config"
	"os"
	"sync"
)

type HashResult struct {
	FileEntry
	digest string
}

func worker(cfg *config.Config, c <-chan FileEntry, r chan<- HashResult) {
	for entry := range c {
		f, err := os.Open(entry.path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		h := sha256.New()
		io.Copy(h, f)
		f.Close()

		digest := hex.EncodeToString(h.Sum(nil))

		cfg.Logf("Digest of %s is %s\n", entry.path, digest)

		r <- HashResult{entry, digest}
	}
}

func Hash(cfg *config.Config, c <-chan FileEntry, r chan<- HashResult) {
	defer close(r)
	var wg sync.WaitGroup

	for range cfg.Workers {
		wg.Go(func() { worker(cfg, c, r) })
	}

	wg.Wait()
}
