package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Category string

const (
	Videos    Category = "videos"
	Pictures  Category = "pictures"
	Documents Category = "documents"
)

var CategoryMap = map[string]Category{
	".png":  Pictures,
	".jpg":  Pictures,
	".jpeg": Pictures,
	".heic": Pictures,
	".docx": Documents,
	".doc":  Documents,
	".pdf":  Documents,
	".txt":  Documents,
	".mp4":  Videos,
}

type Config struct {
	Root      string
	Flatten   bool
	Recursive bool
	DryRun    bool
	Verbose   bool
	Workers   int
}

// Logf prints only when verbose output is enabled (always on during a dry run).
func (c *Config) Logf(format string, args ...any) {
	if c.Verbose {
		fmt.Printf(format, args...)
	}
}

func Load() (*Config, error) {
	flagDryRun := flag.Bool("dry-run", false, "Do a dry-run of the program")
	flagFlatten := flag.Bool("flatten", false, "Flatten the directory")
	flagRecursive := flag.Bool("recursive", false, "The program does a recursive run over root's sub-directories")
	flagHelp := flag.Bool("help", false, "Display the help message")
	flagWorkers := flag.Int("workers", 5, "Number of workers to launch with. Default 5")
	flagVerbose := flag.Bool("verbose", false, "Log every action. Enabled automatically by -dry-run")

	flag.Parse()

	var cfg Config

	cfg.DryRun = *flagDryRun
	cfg.Flatten = *flagFlatten
	cfg.Recursive = *flagRecursive
	cfg.Workers = *flagWorkers
	cfg.Verbose = *flagVerbose || *flagDryRun
	root, err := filepath.Abs(".")

	if err != nil {
		return nil, err
	}

	cfg.Root = root

	if *flagHelp {
		flag.Usage()
		os.Exit(0)
	}

	return &cfg, nil
}
