package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newVMHardDiskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disk",
		Short: "manage VM hard disk drives",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need valid command")
		},
	}
	cmd.AddCommand(
		newVMHardDiskListCmd(),
		newVMHardDiskAddCmd(),
		newVMHardDiskRemoveCmd(),
	)
	return cmd
}

func newVMHardDiskListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list hard disk drives attached to a VM",
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
}

func newVMHardDiskAddCmd() *cobra.Command {
	var path string
	c := &cobra.Command{
		Use:   "add",
		Short: "attach a VHD to a VM",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			if path == "" {
				return fmt.Errorf("--path is required")
			}
			return hyperv.AddVMHardDiskDrive(args[0], path)
		},
	}
	c.Flags().StringVarP(&path, "path", "p", "", "VHD path")
	return c
}

func newVMHardDiskRemoveCmd() *cobra.Command {
	var path string
	c := &cobra.Command{
		Use:   "remove",
		Short: "detach a VHD from a VM",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			if path == "" {
				return fmt.Errorf("--path is required")
			}
			return hyperv.RemoveVMHardDiskDrive(args[0], path)
		},
	}
	c.Flags().StringVarP(&path, "path", "p", "", "VHD path")
	return c
}
