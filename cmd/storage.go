package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "management storage",
	RunE: func(_ *cobra.Command, _ []string) error {
		return fmt.Errorf("need valid command")
	},
}

var storageList = &cobra.Command{
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
