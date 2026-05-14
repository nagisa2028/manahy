package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

var vmHardDiskCmd = &cobra.Command{
	Use:   "disk",
	Short: "manage VM hard disk drives",
	RunE: func(_ *cobra.Command, _ []string) error {
		return fmt.Errorf("need valid command")
	},
}

var vmHardDiskList = &cobra.Command{
	Use:   "list",
	Short: "list hard disk drives for a VM",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		out, err := hyperv.GetVMHardDiskDrives(args[0])
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	},
}

var vmHardDiskAdd = &cobra.Command{
	Use:   "add",
	Short: "add a hard disk drive to a VM",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		if vmHardDiskPath == "" {
			return fmt.Errorf("--path is required")
		}
		return hyperv.AddVMHardDiskDrive(args[0], vmHardDiskPath)
	},
}

var vmHardDiskRemove = &cobra.Command{
	Use:   "remove",
	Short: "remove a hard disk drive from a VM",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		if vmHardDiskPath == "" {
			return fmt.Errorf("--path is required")
		}
		return hyperv.RemoveVMHardDiskDrive(args[0], vmHardDiskPath)
	},
}
