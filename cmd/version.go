package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print manahy version",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Print("manahy version 0.0.0 (beta)\n")
	},
}
