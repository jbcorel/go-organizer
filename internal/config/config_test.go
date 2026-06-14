package config_test

import (
	"io"
	"organizer/internal/config"
	"os"
	"testing"
)

// captureStdout swaps os.Stdout for a pipe while fn runs and returns what it
// printed, then restores stdout. The way to test print-based output.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = orig

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read captured output: %v", err)
	}
	return string(out)
}

func TestLogf(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
		want    string
	}{
		{"verbose prints", true, "hello world\n"},
		{"non-verbose is silent", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Verbose: tt.verbose}

			got := captureStdout(t, func() {
				cfg.Logf("hello %s\n", "world")
			})

			if got != tt.want {
				t.Errorf("Logf output = %q, want %q", got, tt.want)
			}
		})
	}
}
