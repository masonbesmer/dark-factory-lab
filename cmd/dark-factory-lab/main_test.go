package main_test

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"testing"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "dark-factory-lab")
	if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	return bin
}

func TestBinary(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)

	tests := []struct {
		name string
		args []string
		want string
	}{
		{"no args greets world", nil, "Hello, world!\n"},
		{"first arg is greeted", []string{"Ada"}, "Hello, Ada!\n"},
		{"blank arg greets world", []string{"  "}, "Hello, world!\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var stdout, stderr bytes.Buffer
			cmd := exec.Command(bin, tc.args...)
			cmd.Stdout, cmd.Stderr = &stdout, &stderr
			if err := cmd.Run(); err != nil {
				t.Fatalf("run %v: %v (stderr %q)", tc.args, err, stderr.String())
			}
			if got := stdout.String(); got != tc.want {
				t.Errorf("stdout = %q, want %q", got, tc.want)
			}
		})
	}
}
