package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newCheckpointCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoint",
		Short: "manage VM checkpoints",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need valid command")
		},
	}
	cmd.AddCommand(
		newCheckpointCreateCmd(),
		newCheckpointListCmd(),
		newCheckpointRestoreCmd(),
		newCheckpointRemoveCmd(),
		newCheckpointRenameCmd(),
		newCheckpointExportCmd(),
	)
	return cmd
}

func newCheckpointCreateCmd() *cobra.Command {
	var name string
	c := &cobra.Command{
		Use:   "create",
		Short: "create a checkpoint for a VM",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.CreateCheckpoint(args[0], name)
		},
	}
	c.Flags().StringVarP(&name, "name", "n", "", "checkpoint name")
	return c
}

func newCheckpointListCmd() *cobra.Command {
	return &cobra.Command{
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
}

func newCheckpointRestoreCmd() *cobra.Command {
	var name string
	c := &cobra.Command{
		Use:   "restore",
		Short: "restore a VM to a checkpoint",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			return hyperv.RestoreCheckpoint(args[0], name)
		},
	}
	c.Flags().StringVarP(&name, "name", "n", "", "checkpoint name")
	return c
}

func newCheckpointRemoveCmd() *cobra.Command {
	var name string
	c := &cobra.Command{
		Use:   "remove",
		Short: "delete a checkpoint from a VM",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			return hyperv.RemoveCheckpoint(args[0], name)
		},
	}
	c.Flags().StringVarP(&name, "name", "n", "", "checkpoint name")
	return c
}

func newCheckpointRenameCmd() *cobra.Command {
	var name, newName string
	c := &cobra.Command{
		Use:   "rename",
		Short: "rename a checkpoint",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if newName == "" {
				return fmt.Errorf("--new-name is required")
			}
			return hyperv.RenameCheckpoint(args[0], name, newName)
		},
	}
	c.Flags().StringVarP(&name, "name", "n", "", "checkpoint name")
	c.Flags().StringVarP(&newName, "new-name", "", "", "new checkpoint name")
	return c
}

func newCheckpointExportCmd() *cobra.Command {
	var name, path string
	c := &cobra.Command{
		Use:   "export",
		Short: "export a checkpoint to a directory",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if path == "" {
				return fmt.Errorf("--path is required")
			}
			return hyperv.ExportCheckpoint(args[0], name, path)
		},
	}
	c.Flags().StringVarP(&name, "name", "n", "", "checkpoint name")
	c.Flags().StringVarP(&path, "path", "p", "", "export destination directory")
	return c
}
