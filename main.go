// Package main is the entry point for the manahy CLI.
package main

import (
	"os"

	"github.com/DevelopNaoki/manahy/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
