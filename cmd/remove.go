package cmd

import (
	"fmt"
	"os"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newRemoveCmd(configFile *string) *cobra.Command {
	var dryRun, force bool
	cmd := &cobra.Command{
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
			fmt.Fprintf(os.Stderr, "using config: %s\n", path)
			if dryRun {
				hyperv.DryRunRemove(data, os.Stdout)
				return nil
			}
			if !confirmAction("Remove all resources defined in "+path+"?", force) {
				return nil
			}
			return hyperv.RemoveByStruct(data, os.Stderr)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be removed without making changes")
	cmd.Flags().BoolVar(&force, "force", false, "skip confirmation prompt")
	return cmd
}
