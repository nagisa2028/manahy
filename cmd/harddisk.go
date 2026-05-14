package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newVMHardDiskCmd(configFile *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disk",
		Short: "manage VM hard disk drives",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need valid command")
		},
	}
	cmd.AddCommand(
		newVMHardDiskListCmd(),
		newVMHardDiskAddCmd(configFile),
		newVMHardDiskRemoveCmd(configFile),
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

func newVMHardDiskAddCmd(configFile *string) *cobra.Command {
	var diskAlias string
	c := &cobra.Command{
		Use:   "add",
		Short: "attach a VHD to a VM",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			path := resolveDisk(loadConfig(*configFile), diskAlias)
			return hyperv.AddVMHardDiskDrive(args[0], path)
		},
	}
	c.Flags().StringVarP(&diskAlias, "path", "p", "", "VHD path or disk alias")
	_ = c.MarkFlagRequired("path")
	return c
}

func newVMHardDiskRemoveCmd(configFile *string) *cobra.Command {
	var diskAlias string
	c := &cobra.Command{
		Use:   "remove",
		Short: "detach a VHD from a VM",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			path := resolveDisk(loadConfig(*configFile), diskAlias)
			return hyperv.RemoveVMHardDiskDrive(args[0], path)
		},
	}
	c.Flags().StringVarP(&diskAlias, "path", "p", "", "VHD path or disk alias")
	_ = c.MarkFlagRequired("path")
	return c
}
