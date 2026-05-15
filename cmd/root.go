// Package cmd implements the manahy CLI commands.
package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

// Execute runs the root command.
func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	var configFile string
	var verbose bool
	cmd := &cobra.Command{
		Use:   "manahy",
		Short: "manahy is a management tool for Hyper-V",
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			hyperv.SetVerbose(verbose)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need a valid subcommand")
		},
	}
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file path")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "print each PowerShell command before execution")
	cmd.AddCommand(
		newVersionCmd(),
		newVMCmd(&configFile),
		newSwitchCmd(),
		newDiskCmd(&configFile),
		newCheckpointCmd(),
		newHostCmd(),
		newBuildCmd(&configFile),
		newRemoveCmd(&configFile),
		newListCmd(&configFile),
	)
	return cmd
}
