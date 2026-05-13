package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

var checkpointCmd = &cobra.Command{
	Use:   "checkpoint",
	Short: "manage VM checkpoints",
	RunE: func(_ *cobra.Command, _ []string) error {
		return fmt.Errorf("need valid command")
	},
}

var checkpointCreate = &cobra.Command{
	Use:   "create",
	Short: "create a checkpoint for a VM",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		return hyperv.CreateCheckpoint(args[0], checkpointName)
	},
}

var checkpointList = &cobra.Command{
	Use:   "list",
	Short: "list checkpoints for a VM",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		out, err := hyperv.GetCheckpoints(args[0])
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	},
}

var checkpointRestore = &cobra.Command{
	Use:   "restore",
	Short: "restore a VM to a checkpoint",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		if checkpointName == "" {
			return fmt.Errorf("--name is required")
		}
		return hyperv.RestoreCheckpoint(args[0], checkpointName)
	},
}

var checkpointRemove = &cobra.Command{
	Use:   "remove",
	Short: "delete a checkpoint from a VM",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		if checkpointName == "" {
			return fmt.Errorf("--name is required")
		}
		return hyperv.RemoveCheckpoint(args[0], checkpointName)
	},
}

var checkpointRename = &cobra.Command{
	Use:   "rename",
	Short: "rename a checkpoint",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		if checkpointName == "" {
			return fmt.Errorf("--name is required")
		}
		if newCheckpointName == "" {
			return fmt.Errorf("--new-name is required")
		}
		return hyperv.RenameCheckpoint(args[0], checkpointName, newCheckpointName)
	},
}

var checkpointExport = &cobra.Command{
	Use:   "export",
	Short: "export a checkpoint to a directory",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		if checkpointName == "" {
			return fmt.Errorf("--name is required")
		}
		if exportPath == "" {
			return fmt.Errorf("--path is required")
		}
		return hyperv.ExportCheckpoint(args[0], checkpointName, exportPath)
	},
}
