// Command dark-factory-lab prints a greeting.
package main

import (
	"os"

	"github.com/masonbesmer/dark-factory-lab/internal/cli"
)

func main() {
	os.Exit(cli.Main(os.Args[1:], os.Stdout, os.Stderr))
}
