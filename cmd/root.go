// Package cmd implements the manahy CLI commands.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RootCmd is the top-level manahy command.
var RootCmd = &cobra.Command{
	Use:   "manahy",
	Short: "manahy is management tool on Hyper-V",
	RunE: func(_ *cobra.Command, _ []string) error {
		return fmt.Errorf("need valid command")
	},
}
