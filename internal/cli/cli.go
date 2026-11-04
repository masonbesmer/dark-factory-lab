// Package cli holds the command-line logic behind cmd/dark-factory-lab.
package cli

import (
	"fmt"
	"io"

	"github.com/masonbesmer/dark-factory-lab/internal/greet"
)

// Run greets the first argument (or the world) on w.
func Run(args []string, w io.Writer) error {
	var name string
	if len(args) > 0 {
		name = args[0]
	}
	_, err := fmt.Fprintln(w, greet.Greet(name))
	return err
}

// Main runs the CLI and returns the process exit code.
func Main(args []string, stdout, stderr io.Writer) int {
	if err := Run(args, stdout); err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}
