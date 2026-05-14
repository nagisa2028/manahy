package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "remove vm, disk and switch from manahy.yaml",
		Run: func(_ *cobra.Command, _ []string) {
			data, err := hyperv.UnmarshalYaml("manahy.yaml")
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			for _, disk := range data.Disks {
				if err = hyperv.RemoveDisk(disk.Path, false); err != nil {
					fmt.Printf("%s\n", err)
				}
			}
			for _, network := range data.Networks {
				if err = hyperv.RemoveSwitch(network.Name); err != nil {
					fmt.Printf("%s\n", err)
				}
			}
		},
	}
}
