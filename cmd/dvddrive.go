package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newVMDvdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dvd",
		Short: "manage VM DVD drives",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need valid command")
		},
	}
	cmd.AddCommand(
		newVMDvdListCmd(),
		newVMDvdAddCmd(),
		newVMDvdRemoveCmd(),
		newVMDvdSetCmd(),
		newVMDvdEjectCmd(),
	)
	return cmd
}

func newVMDvdListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list DVD drives attached to a VM",
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
}

func newVMDvdAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "add a DVD drive to a VM",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.AddVMDvdDrive(args[0])
		},
	}
}

func newVMDvdRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "remove the first DVD drive from a VM",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.RemoveVMDvdDrive(args[0])
		},
	}
}

func newVMDvdSetCmd() *cobra.Command {
	var image string
	c := &cobra.Command{
		Use:   "set",
		Short: "set the ISO image for the DVD drive",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			if image == "" {
				return fmt.Errorf("--image is required")
			}
			return hyperv.SetVMDvdDrive(args[0], image)
		},
	}
	c.Flags().StringVarP(&image, "image", "i", "", "ISO image path")
	return c
}

func newVMDvdEjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "eject",
		Short: "eject the ISO image from the DVD drive",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.SetVMDvdDrive(args[0], "")
		},
	}
}
