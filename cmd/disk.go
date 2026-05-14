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

var diskRemove = &cobra.Command{
	Use:   "remove",
	Short: "remove a virtual disk",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		return hyperv.RemoveDisk(args[0], true)
	},
}

var diskInfo = &cobra.Command{
	Use:   "info",
	Short: "show VHD info",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		out, err := hyperv.GetVHDInfo(args[0])
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	},
}

var diskResize = &cobra.Command{
	Use:   "resize",
	Short: "resize a VHD",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		if resizeSize == "" {
			return fmt.Errorf("--size is required")
		}
		return hyperv.ResizeVHD(args[0], resizeSize)
	},
}

var diskOptimize = &cobra.Command{
	Use:   "optimize",
	Short: "optimize a VHD",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		return hyperv.OptimizeVHD(args[0])
	},
}

var diskConvert = &cobra.Command{
	Use:   "convert",
	Short: "convert a VHD",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		if convertDestPath == "" {
			return fmt.Errorf("--dest is required")
		}
		return hyperv.ConvertVHD(args[0], convertDestPath, convertDiskType)
	},
}

var diskMount = &cobra.Command{
	Use:   "mount",
	Short: "mount a VHD",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		return hyperv.MountVHD(args[0])
	},
}

var diskDismount = &cobra.Command{
	Use:   "dismount",
	Short: "dismount a VHD",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		return hyperv.DismountVHD(args[0])
	},
}

var diskMerge = &cobra.Command{
	Use:   "merge",
	Short: "merge a VHD into its parent",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		return hyperv.MergeVHD(args[0], mergeDest)
	},
}
