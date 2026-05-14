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
	var configFile string
	cmd := &cobra.Command{
		Use:   "manahy",
		Short: "manahy is a management tool for Hyper-V",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need a valid subcommand")
		},
	}
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file path")
	cmd.AddCommand(
		newVersionCmd(),
		newVMCmd(&configFile),
		newSwitchCmd(),
		newDiskCmd(&configFile),
		newCheckpointCmd(),
		newHostCmd(),
		newBuildCmd(&configFile),
		newRemoveCmd(&configFile),
	)
	return cmd
}
