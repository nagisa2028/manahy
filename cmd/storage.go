package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newStorageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage",
		Short: "manage physical storage",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need a valid subcommand")
		},
	}
	cmd.AddCommand(newStorageListCmd())
	return cmd
}

func newStorageListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "listing all storage",
		Args:  cobra.RangeArgs(0, 0),
		RunE: func(_ *cobra.Command, _ []string) error {
			storageList, err := hyperv.GetStorageList()
			if err != nil {
				return err
			}
			displayStorageList(storageList)
			return nil
		},
	}
}
