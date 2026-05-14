package cmd

import (
	"fmt"
	"os"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newBuildCmd(configFile *string) *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "build",
		Short: "create vm, disk and switch from config file",
		RunE: func(_ *cobra.Command, _ []string) error {
			path := *configFile
			if path == "" {
				path = "manahy.yaml"
			}
			data, err := hyperv.UnmarshalYaml(path)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "using config: %s\n", path)
			if dryRun {
				hyperv.DryRunBuild(data, os.Stdout)
				return nil
			}
			return hyperv.BuildByStruct(data)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be created without making changes")
	return cmd
}
