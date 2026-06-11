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


type hashResult struct {
	path string
	digest string
}


func worker(c <-chan string, r chan<- hashResult) {
	for path := range c {
		f, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		h := sha256.New()
		io.Copy(h, f)
		f.Close()
		
		digest := hex.EncodeToString(h.Sum(nil))
		
		r <- hashResult{path, digest}
	}
}


func Hash(cfg *config.Config, c <-chan string, r chan<- hashResult) {
	var wg sync.WaitGroup
	
	for range cfg.Workers {
		wg.Go(func() { worker(c, r) })
	}
	
	wg.Wait()
}