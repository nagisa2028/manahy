package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

var remove = &cobra.Command{
	Use:   "remove",
	Short: "remove vm, disk and switch from manahy.yaml",
	Run: func(_ *cobra.Command, _ []string) {
		data, err := hyperv.UnmarshalYaml("manahy.yaml")
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		for _, disk := range data.Disks {
			err = hyperv.RemoveDisk(disk.Path, false)
			if err != nil {
				fmt.Printf("%s\n", err)
			}
		}

		for _, network := range data.Networks {
			err = hyperv.RemoveSwitch(network.Name)
			if err != nil {
				fmt.Printf("%s\n", err)
			}
		}
	},
}
