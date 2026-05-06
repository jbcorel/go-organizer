package main

import (
	"flag"
	"os"
	"path/filepath"
)

type Folder int

const (
	Videos Folder = iota
	Pictures
	Documents
)

var folderMap = map[Folder]string{
	Videos:    "vids",
	Pictures:  "pics",
	Documents: "docs",
}

func (f Folder) String() string {
	return folderMap[f]
}

var knownFormatsToType = map[string]Folder{
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
	root      string
	flatten   bool
	recursive bool
	dryRun    bool
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

// Not yet supporting non-current-dir path
func (cfg *Config) rootFromArgs() error {
	fp, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	cfg.root = fp
	return nil
}

var Cfg Config

func InitConfig() error {
	flagDryRun := flag.Bool("dry-run", false, "Do a dry-run of the program")
	flagFlatten := flag.Bool("flatten", false, "Flatten the directory")
	flagRecursive := flag.Bool("recursive", false, "The program does a recursive run over root's sub-directories")
	flagHelp := flag.Bool("help", false, "Display the help message")

	flag.Parse()

	Cfg.setDryRun(*flagDryRun)
	Cfg.setFlatten(*flagFlatten)
	Cfg.setRecursive(*flagRecursive)
	Cfg.rootFromArgs()

	if *flagHelp {
		flag.Usage()
		os.Exit(0)
	}

	return nil
}
