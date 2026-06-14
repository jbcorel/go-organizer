package pipeline

import (
	"organizer/internal/config"
	"path/filepath"
	"testing"
)

// Pushes results through Organize via a buffered channel we close so it returns.
func feedOrganize(cfg *config.Config, results ...HashResult) {
	c := make(chan HashResult, len(results))
	for _, r := range results {
		c <- r
	}
	close(c)
	Organize(cfg, c)
}

func TestOrganizeCategorizes(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "photo.png")
	writeFile(t, src, "img")

	cfg := &config.Config{Root: root}
	feedOrganize(cfg, HashResult{
		FileEntry: FileEntry{path: src, name: "photo.png", ext: ".png"},
		digest:    "digest-1",
	})

	dest := filepath.Join(root, "pictures", "photo.png")
	if !exists(t, dest) {
		t.Errorf("expected file moved to %s", dest)
	}
	if exists(t, src) {
		t.Errorf("source %s should no longer exist", src)
	}
}

func TestOrganizeDuplicatesGoToDuplicatesFolder(t *testing.T) {
	root := t.TempDir()
	a := filepath.Join(root, "a.png")
	b := filepath.Join(root, "b.png")
	writeFile(t, a, "img")
	writeFile(t, b, "img")

	cfg := &config.Config{Root: root}
	// Same digest + ext => b is a duplicate of a.
	feedOrganize(cfg,
		HashResult{FileEntry{path: a, name: "a.png", ext: ".png"}, "same"},
		HashResult{FileEntry{path: b, name: "b.png", ext: ".png"}, "same"},
	)

	if !exists(t, filepath.Join(root, "pictures", "a.png")) {
		t.Errorf("first file should land in pictures/")
	}
	if !exists(t, filepath.Join(root, "duplicates", "b.png")) {
		t.Errorf("duplicate should land in duplicates/")
	}
}

func TestOrganizeDryRunMovesNothing(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "photo.png")
	writeFile(t, src, "img")

	cfg := &config.Config{Root: root, DryRun: true}
	feedOrganize(cfg, HashResult{
		FileEntry: FileEntry{path: src, name: "photo.png", ext: ".png"},
		digest:    "digest-1",
	})

	if !exists(t, src) {
		t.Errorf("dry run should leave the source file in place")
	}
	if exists(t, filepath.Join(root, "pictures")) {
		t.Errorf("dry run should not create category folders")
	}
}
