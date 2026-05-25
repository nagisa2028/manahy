// Package main is the entry point for the manahy CLI.
package main

import (
	"fmt"
	"os"

	"github.com/DevelopNaoki/manahy/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
