package pipeline

import (
	"organizer/internal/config"
	"path/filepath"
	"slices"
	"testing"
)

func collectScan(cfg *config.Config, root string) []string {
	c := make(chan FileEntry)
	go Scan(cfg, root, c)

	var names []string
	for e := range c {
		names = append(names, e.name)
	}
	slices.Sort(names)
	return names
}

func TestScanNonRecursive(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "photo.png"), "x")
	writeFile(t, filepath.Join(root, "notes.txt"), "x")
	writeFile(t, filepath.Join(root, "skip.xyz"), "x") // unknown extension

	sub := filepath.Join(root, "sub")
	mkDir(t, sub)
	writeFile(t, filepath.Join(sub, "deep.pdf"), "x") // must NOT be seen

	cfg := &config.Config{Recursive: false}

	got := collectScan(cfg, root)
	want := []string{"notes.txt", "photo.png"}

	if !slices.Equal(got, want) {
		t.Errorf("Scan() = %v, want %v", got, want)
	}
}

func TestScanRecursive(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "photo.png"), "x")

	sub := filepath.Join(root, "sub")
	mkDir(t, sub)
	writeFile(t, filepath.Join(sub, "deep.pdf"), "x")

	nested := filepath.Join(sub, "nested")
	mkDir(t, nested)
	writeFile(t, filepath.Join(nested, "more.jpg"), "x")

	cfg := &config.Config{Recursive: true}

	got := collectScan(cfg, root)
	want := []string{"deep.pdf", "more.jpg", "photo.png"}

	if !slices.Equal(got, want) {
		t.Errorf("Scan() = %v, want %v", got, want)
	}
}

func TestScanMissingDirIsNotFatal(t *testing.T) {
	cfg := &config.Config{}
	// Scanning a path that doesn't exist should just yield nothing, not panic.
	got := collectScan(cfg, filepath.Join(t.TempDir(), "does-not-exist"))
	if len(got) != 0 {
		t.Errorf("Scan(missing dir) = %v, want empty", got)
	}
}
