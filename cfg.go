package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var knownFormatsToType = map[string]string{
	".png":  "pics",
	".jpg":  "pics",
	".jpeg": "pics",
	".heic": "pics",
	".docx": "docs",
	".doc":  "docs",
	".pdf":  "docs",
	".txt":  "docs",
	".mp4":  "vids",
}

type UnknownFlagError struct {
	flags []string
}

func (e *UnknownFlagError) Error() string {
	return fmt.Sprintf("Unknown flags: %s", strings.Join(e.flags, ", "))
}

type Config struct {
	root      string
	flatten   bool
	recursive bool
	dryRun    bool
}

func (cfg *Config) setRoot(root string) error {
	cfg.root = root
	return nil
}

func (cfg *Config) setFlatten(flatten bool) error {
	cfg.flatten = flatten
	return nil
}

func (cfg *Config) setRecursive(recursive bool) error {
	cfg.recursive = recursive
	return nil
}

func (cfg *Config) setDryRun(dryRun bool) error {
	cfg.dryRun = dryRun
	return nil
}

func NewConfig() (*Config, error) {
	flagDryRun := flag.Bool("dry-run", false, "Do a dry-run of the program")
	flagFlatten := flag.Bool("flatten", false, "Flatten the directory")
	flagRecursive := flag.Bool("recursive", false, "The program does a recursive run over root's sub-directories")
	flagHelp := flag.Bool("help", false, "Display the help message")

	flag.Parse()

	root, _ := rootFromArgs()

	cfg := &Config{}

	cfg.setDryRun(*flagDryRun)
	cfg.setFlatten(*flagFlatten)
	cfg.setRecursive(*flagRecursive)
	cfg.setRoot(root)

	if *flagHelp {
		flag.Usage()
		os.Exit(0)
	}

	return cfg, nil
}
