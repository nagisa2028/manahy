package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

var diskCmd = &cobra.Command{
	Use:   "disk",
	Short: "management virtual disk",
	RunE: func(_ *cobra.Command, _ []string) error {
		return fmt.Errorf("need valid command")
	},
}

var diskCreate = &cobra.Command{
	Use:   "create",
	Short: "Create virtual disk",
	Args:  cobra.RangeArgs(0, 0),
	RunE: func(_ *cobra.Command, _ []string) error {
		return hyperv.CreateDisk(diskCreateOption, true)
	},
}
