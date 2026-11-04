// Package greet builds greetings.
package greet

import (
	"fmt"
	"strings"
)

// Greet returns "Hello, <name>!", or "Hello, world!" when name is blank.
func Greet(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "world"
	}
	return fmt.Sprintf("Hello, %s!", name)
}
