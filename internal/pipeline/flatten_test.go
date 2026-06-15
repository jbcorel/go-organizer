package pipeline

import (
	"organizer/internal/config"
	"path/filepath"
	"testing"
)

func TestFlatten(t *testing.T) {
	root := t.TempDir()

	pics := filepath.Join(root, "pictures")
	docs := filepath.Join(root, "documents")
	dups := filepath.Join(root, "duplicates")
	mkDir(t, pics)
	mkDir(t, docs)
	mkDir(t, dups)

	writeFile(t, filepath.Join(pics, "a.png"), "img")
	writeFile(t, filepath.Join(docs, "b.pdf"), "doc")
	writeFile(t, filepath.Join(dups, "a.png"), "img-copy") 

	cfg := &config.Config{Root: root}
	if err := Flatten(cfg); err != nil {
		t.Fatalf("Flatten: %v", err)
	}

	if !exists(t, filepath.Join(root, "a.png")) {
		t.Errorf("pictures/a.png should be moved to root")
	}
	if !exists(t, filepath.Join(root, "b.pdf")) {
		t.Errorf("documents/b.pdf should be moved to root")
	}

	if exists(t, pics) {
		t.Errorf("emptied pictures/ should be removed")
	}
	if exists(t, docs) {
		t.Errorf("emptied documents/ should be removed")
	}

	// duplicates/ is intentionally skipped 
	if !exists(t, dups) || !exists(t, filepath.Join(dups, "a.png")) {
		t.Errorf("duplicates/ should be left untouched")
	}
}

func TestFlattenDryRunChangesNothing(t *testing.T) {
	root := t.TempDir()
	pics := filepath.Join(root, "pictures")
	mkDir(t, pics)
	writeFile(t, filepath.Join(pics, "a.png"), "img")

	cfg := &config.Config{Root: root, DryRun: true}
	if err := Flatten(cfg); err != nil {
		t.Fatalf("Flatten: %v", err)
	}

	if exists(t, filepath.Join(root, "a.png")) {
		t.Errorf("dry run should not move files")
	}
	if !exists(t, filepath.Join(pics, "a.png")) {
		t.Errorf("dry run should leave files in place")
	}
	if !exists(t, pics) {
		t.Errorf("dry run should not remove folders")
	}
}
