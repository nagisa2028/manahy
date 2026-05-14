package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

var vmDvdCmd = &cobra.Command{
	Use:   "dvd",
	Short: "manage VM DVD drives",
	RunE: func(_ *cobra.Command, _ []string) error {
		return fmt.Errorf("need valid command")
	},
}

var vmDvdList = &cobra.Command{
	Use:   "list",
	Short: "list DVD drives for a VM",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		out, err := hyperv.GetVMDvdDrives(args[0])
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	},
}

var vmDvdAdd = &cobra.Command{
	Use:   "add",
	Short: "add a DVD drive to a VM",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		return hyperv.AddVMDvdDrive(args[0])
	},
}

var vmDvdRemove = &cobra.Command{
	Use:   "remove",
	Short: "remove a DVD drive from a VM",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		return hyperv.RemoveVMDvdDrive(args[0])
	},
}

var vmDvdSet = &cobra.Command{
	Use:   "set",
	Short: "set an ISO image on the DVD drive",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		return hyperv.SetVMDvdDrive(args[0], dvdImagePath)
	},
}

var vmDvdEject = &cobra.Command{
	Use:   "eject",
	Short: "eject the ISO image from the DVD drive",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		return hyperv.SetVMDvdDrive(args[0], "")
	},
}
