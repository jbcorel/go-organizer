// package fsutils_test (not fsutils) = black-box: only the exported API is
// reachable, like a real caller.
package fsutils_test

import (
	"organizer/internal/fsutils"
	"os"
	"path/filepath"
	"testing"
)

func TestDirExists(t *testing.T) {
	root := t.TempDir()

	file := filepath.Join(root, "file.txt")
	if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	// Table-driven: one row per scenario, each run as a named subtest.
	tests := []struct {
		name      string
		path      string
		wantOK    bool
		wantError bool
	}{
		{"existing directory", root, true, false},
		{"missing path", filepath.Join(root, "nope"), false, false},
		{"path is a file", file, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := fsutils.DirExists(tt.path)

			if ok != tt.wantOK {
				t.Errorf("DirExists(%q) ok = %v, want %v", tt.path, ok, tt.wantOK)
			}
			if (err != nil) != tt.wantError {
				t.Errorf("DirExists(%q) err = %v, wantError = %v", tt.path, err, tt.wantError)
			}
		})
	}
}

func TestCreateDir(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "created")

	if err := fsutils.CreateDir(target); err != nil {
		t.Fatalf("CreateDir: %v", err)
	}

	ok, err := fsutils.DirExists(target)
	if err != nil {
		t.Fatalf("DirExists after create: %v", err)
	}
	if !ok {
		t.Errorf("CreateDir did not create %s", target)
	}
}
