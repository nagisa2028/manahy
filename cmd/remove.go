package cmd

import (
	"os"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newRemoveCmd(configFile *string) *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "remove vm, disk and switch from config file",
		RunE: func(_ *cobra.Command, _ []string) error {
			path := *configFile
			if path == "" {
				path = "manahy.yaml"
			}
			data, err := hyperv.UnmarshalYaml(path)
			if err != nil {
				return err
			}
			return hyperv.RemoveByStruct(data, os.Stderr)
		},
	}
}
