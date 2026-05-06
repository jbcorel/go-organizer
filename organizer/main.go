package main

import (
	"fmt"
	"os"
)

func main() {
	err := InitConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	entries, err := readDir(Cfg.root)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, entry := range entries {
		err := processFileEntry(entry)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Finished")
}
