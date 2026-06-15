package pipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"organizer/internal/config"
	"path/filepath"
	"testing"
)

func sha256Hex(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

func TestHashComputesDigests(t *testing.T) {
	root := t.TempDir()

	a := filepath.Join(root, "a.txt")
	b := filepath.Join(root, "b.txt") // identical content to a -> identical digest
	c := filepath.Join(root, "c.txt") // different content -> different digest
	writeFile(t, a, "hello")
	writeFile(t, b, "hello")
	writeFile(t, c, "world")

	in := make(chan FileEntry, 3)
	out := make(chan HashResult, 3)
	in <- FileEntry{path: a, name: "a.txt", ext: ".txt"}
	in <- FileEntry{path: b, name: "b.txt", ext: ".txt"}
	in <- FileEntry{path: c, name: "c.txt", ext: ".txt"}
	close(in)

	cfg := &config.Config{Workers: 2}
	Hash(cfg, in, out)

	got := make(map[string]string)
	for r := range out {
		got[r.name] = r.digest
	}

	if len(got) != 3 {
		t.Fatalf("got %d results, want 3", len(got))
	}
	if want := sha256Hex("hello"); got["a.txt"] != want {
		t.Errorf("a.txt digest = %s, want %s", got["a.txt"], want)
	}
	if got["a.txt"] != got["b.txt"] {
		t.Errorf("identical files produced different digests: %s vs %s", got["a.txt"], got["b.txt"])
	}
	if got["a.txt"] == got["c.txt"] {
		t.Errorf("different files produced the same digest: %s", got["a.txt"])
	}
}
