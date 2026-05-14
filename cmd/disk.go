package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newDiskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disk",
		Short: "management virtual disk",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need valid command")
		},
	}
	cmd.AddCommand(
		newDiskCreateCmd(),
		newDiskRemoveCmd(),
		newDiskInfoCmd(),
		newDiskResizeCmd(),
		newDiskOptimizeCmd(),
		newDiskConvertCmd(),
		newDiskMountCmd(),
		newDiskDismountCmd(),
		newDiskMergeCmd(),
	)
	return cmd
}

func newDiskCreateCmd() *cobra.Command {
	var opt hyperv.Disk
	c := &cobra.Command{
		Use:   "create",
		Short: "Create virtual disk",
		Args:  cobra.RangeArgs(0, 0),
		RunE: func(_ *cobra.Command, _ []string) error {
			return hyperv.CreateDisk(opt, true)
		},
	}
	c.Flags().StringVarP(&opt.Path, "path", "p", "", "set path")
	c.Flags().StringVarP(&opt.Size, "size", "s", "", "set size")
	c.Flags().StringVarP(&opt.Type, "type", "t", "dynamic", "set type")
	c.Flags().StringVarP(&opt.ParentPath, "parent-path", "", "", "set parent path")
	c.Flags().IntVarP(&opt.SourceDisk, "source-disk", "", 0, "set source disk")
	return c
}

func newDiskRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "remove a virtual disk",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.RemoveDisk(args[0], true)
		},
	}
}

func newDiskInfoCmd() *cobra.Command {
	return &cobra.Command{
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
}

func newDiskResizeCmd() *cobra.Command {
	var size string
	c := &cobra.Command{
		Use:   "resize",
		Short: "resize a VHD",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			if size == "" {
				return fmt.Errorf("--size is required")
			}
			return hyperv.ResizeVHD(args[0], size)
		},
	}
	c.Flags().StringVarP(&size, "size", "s", "", "new size (e.g. 20GB)")
	return c
}

func newDiskOptimizeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "optimize",
		Short: "optimize a VHD",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.OptimizeVHD(args[0])
		},
	}
}

func newDiskConvertCmd() *cobra.Command {
	var dest, diskType string
	c := &cobra.Command{
		Use:   "convert",
		Short: "convert a VHD",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			if dest == "" {
				return fmt.Errorf("--dest is required")
			}
			return hyperv.ConvertVHD(args[0], dest, diskType)
		},
	}
	c.Flags().StringVarP(&dest, "dest", "d", "", "destination path")
	c.Flags().StringVarP(&diskType, "type", "t", "", "disk type (Dynamic or Fixed)")
	return c
}

func newDiskMountCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mount",
		Short: "mount a VHD",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.MountVHD(args[0])
		},
	}
}

func newDiskDismountCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dismount",
		Short: "dismount a VHD",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.DismountVHD(args[0])
		},
	}
}

func newDiskMergeCmd() *cobra.Command {
	var dest string
	c := &cobra.Command{
		Use:   "merge",
		Short: "merge a VHD into its parent",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.MergeVHD(args[0], dest)
		},
	}
	c.Flags().StringVarP(&dest, "dest", "d", "", "destination path (default: merge into parent)")
	return c
}
