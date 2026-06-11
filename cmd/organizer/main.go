package main

import (
	"organizer/internal/pipeline"
	"organizer/internal/config"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		panic(err)
	}

	// add flatten too
	
}