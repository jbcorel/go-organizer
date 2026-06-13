package main

import (
	"fmt"
	"organizer/internal/config"
	"organizer/internal/pipeline"
	"sync"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		panic(err)
	}

	cEntry := make(chan pipeline.FileEntry)
	cResult := make(chan pipeline.HashResult)

	var wg sync.WaitGroup
	wg.Go(func() { pipeline.Scan(cfg, cfg.Root, cEntry) })
	wg.Go(func() { pipeline.Hash(cfg, cEntry, cResult) })
	wg.Go(func() { pipeline.Organize(cfg, cResult) })

	wg.Wait()

	fmt.Println("Finished")

	// add flatten too

}
