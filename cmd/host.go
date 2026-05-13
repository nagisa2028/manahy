package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "show Hyper-V host information",
	RunE: func(_ *cobra.Command, _ []string) error {
		return fmt.Errorf("need valid command")
	},
}

var hostShow = &cobra.Command{
	Use:   "show",
	Short: "show Hyper-V host configuration",
	Args:  cobra.RangeArgs(0, 0),
	RunE: func(_ *cobra.Command, _ []string) error {
		out, err := hyperv.GetVMHost()
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	},
}
