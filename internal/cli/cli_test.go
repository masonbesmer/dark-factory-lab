package cli_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/masonbesmer/dark-factory-lab/internal/cli"
)

type failWriter struct{}

func (failWriter) Write([]byte) (int, error) { return 0, errors.New("boom") }

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{"no args greets world", nil, "Hello, world!\n", false},
		{"first arg is greeted", []string{"Ada"}, "Hello, Ada!\n", false},
		{"extra args ignored", []string{"Ada", "Grace"}, "Hello, Ada!\n", false},
		{"blank arg greets world", []string{"  "}, "Hello, world!\n", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			err := cli.Run(tc.args, &buf)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Run err = %v, wantErr %v", err, tc.wantErr)
			}
			if got := buf.String(); got != tc.want {
				t.Errorf("Run output = %q, want %q", got, tc.want)
			}
		})
	}

	t.Run("write failure is returned", func(t *testing.T) {
		t.Parallel()
		if err := cli.Run([]string{"Ada"}, failWriter{}); err == nil {
			t.Fatal("Run with failing writer returned nil error")
		}
	})
}

func TestMainExit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout string
		wantStderr bool
	}{
		{name: "success", args: []string{"Ada"}, wantCode: 0, wantStdout: "Hello, Ada!\n"},
		{name: "no args", wantCode: 0, wantStdout: "Hello, world!\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var stdout, stderr bytes.Buffer
			code := cli.Main(tc.args, &stdout, &stderr)
			if code != tc.wantCode {
				t.Errorf("Main code = %d, want %d", code, tc.wantCode)
			}
			if stdout.String() != tc.wantStdout {
				t.Errorf("Main stdout = %q, want %q", stdout.String(), tc.wantStdout)
			}
			if (stderr.Len() > 0) != tc.wantStderr {
				t.Errorf("Main stderr = %q, wantStderr %v", stderr.String(), tc.wantStderr)
			}
		})
	}

	t.Run("write failure exits 1 and reports on stderr", func(t *testing.T) {
		t.Parallel()
		var stderr bytes.Buffer
		if code := cli.Main([]string{"Ada"}, failWriter{}, &stderr); code != 1 {
			t.Errorf("Main code = %d, want 1", code)
		}
		if stderr.Len() == 0 {
			t.Error("Main wrote nothing to stderr on failure")
		}
	})
}
