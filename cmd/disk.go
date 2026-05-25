package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newDiskCmd(configFile *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disk",
		Short: "manage virtual disks",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need a valid subcommand")
		},
	}
	cmd.AddCommand(
		newDiskCreateCmd(),
		newDiskCloneCmd(configFile),
		newDiskRemoveCmd(configFile),
		newDiskInfoCmd(configFile),
		newDiskResizeCmd(configFile),
		newDiskOptimizeCmd(configFile),
		newDiskConvertCmd(configFile),
		newDiskMountCmd(configFile),
		newDiskDismountCmd(configFile),
		newDiskMergeCmd(configFile),
	)
	return cmd
}

func newDiskCreateCmd() *cobra.Command {
	var opt hyperv.Disk
	c := &cobra.Command{
		Use:   "create",
		Short: "Create virtual disk",
		Args:  cobra.NoArgs,
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

func newDiskRemoveCmd(configFile *string) *cobra.Command {
	var force bool
	c := &cobra.Command{
		Use:   "remove",
		Short: "remove a virtual disk",
		Args:  cobra.RangeArgs(1, maxBulkArgs),
		RunE: func(_ *cobra.Command, args []string) error {
			prompt := fmt.Sprintf("Remove %d disk(s)?", len(args))
			if !confirmAction(prompt, force) {
				return nil
			}
			for _, arg := range args {
				path := resolveDisk(loadConfig(*configFile), arg)
				if err := hyperv.RemoveDisk(path, true); err != nil {
					return err
				}
			}
			return nil
		},
	}
	c.Flags().BoolVar(&force, "force", false, "skip confirmation prompt")
	return c
}

func newDiskInfoCmd(configFile *string) *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "show VHD info",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			path := resolveDisk(loadConfig(*configFile), args[0])
			out, err := hyperv.GetVHDInfo(path)
			if err != nil {
				return err
			}
			fmt.Print(out)
			return nil
		},
	}
}

func newDiskResizeCmd(configFile *string) *cobra.Command {
	var size string
	c := &cobra.Command{
		Use:   "resize",
		Short: "resize a VHD",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			path := resolveDisk(loadConfig(*configFile), args[0])
			return hyperv.ResizeVHD(path, size)
		},
	}
	c.Flags().StringVarP(&size, "size", "s", "", "new size (e.g. 20GB)")
	_ = c.MarkFlagRequired("size")
	return c
}

func newDiskOptimizeCmd(configFile *string) *cobra.Command {
	return &cobra.Command{
		Use:   "optimize",
		Short: "optimize a VHD",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			path := resolveDisk(loadConfig(*configFile), args[0])
			return hyperv.OptimizeVHD(path)
		},
	}
}

func newDiskConvertCmd(configFile *string) *cobra.Command {
	var dest, diskType string
	c := &cobra.Command{
		Use:   "convert",
		Short: "convert a VHD",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			path := resolveDisk(loadConfig(*configFile), args[0])
			return hyperv.ConvertVHD(path, dest, diskType)
		},
	}
	c.Flags().StringVarP(&dest, "dest", "d", "", "destination path")
	c.Flags().StringVarP(&diskType, "type", "t", "", "disk type (Dynamic or Fixed)")
	_ = c.MarkFlagRequired("dest")
	return c
}

func newDiskMountCmd(configFile *string) *cobra.Command {
	return &cobra.Command{
		Use:   "mount",
		Short: "mount a VHD",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			path := resolveDisk(loadConfig(*configFile), args[0])
			return hyperv.MountVHD(path)
		},
	}
}

func newDiskDismountCmd(configFile *string) *cobra.Command {
	return &cobra.Command{
		Use:   "dismount",
		Short: "dismount a VHD",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			path := resolveDisk(loadConfig(*configFile), args[0])
			return hyperv.DismountVHD(path)
		},
	}
}

func newDiskMergeCmd(configFile *string) *cobra.Command {
	var dest string
	c := &cobra.Command{
		Use:   "merge",
		Short: "merge a VHD into its parent",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			path := resolveDisk(loadConfig(*configFile), args[0])
			return hyperv.MergeVHD(path, dest)
		},
	}
	c.Flags().StringVarP(&dest, "dest", "d", "", "destination path (default: merge into parent)")
	return c
}

func newDiskCloneCmd(configFile *string) *cobra.Command {
	var dest string
	c := &cobra.Command{
		Use:   "clone <source>",
		Short: "clone a VHD to a new path",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			src := resolveDisk(loadConfig(*configFile), args[0])
			return hyperv.CloneDisk(src, dest)
		},
	}
	c.Flags().StringVarP(&dest, "dest", "d", "", "destination path for the cloned VHD")
	_ = c.MarkFlagRequired("dest")
	return c
}
