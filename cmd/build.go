package cmd

import (
	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

var build = &cobra.Command{
	Use:   "build",
	Short: "create vm, disk and switch from manahy.yaml",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := hyperv.UnmarshalYaml("manahy.yaml")
		if err != nil {
			return err
		}
		return hyperv.BuildByStruct(data)
	},
}
