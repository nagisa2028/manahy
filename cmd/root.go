// Package cmd implements the manahy CLI commands.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Execute runs the root command.
func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manahy",
		Short: "manahy is management tool on Hyper-V",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need valid command")
		},
	}
	cmd.AddCommand(
		newVersionCmd(),
		newVMCmd(),
		newSwitchCmd(),
		newDiskCmd(),
		newStorageCmd(),
		newMemberCmd(),
		newCheckpointCmd(),
		newHostCmd(),
		newBuildCmd(),
		newRemoveCmd(),
	)
	return cmd
}
