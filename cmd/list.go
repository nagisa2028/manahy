package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newListCmd(configFile *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "show status of all resources defined in the config file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if *configFile == "" {
				return fmt.Errorf("--config is required")
			}
			rl, err := hyperv.GetResourceList(*configFile)
			if err != nil {
				return err
			}
			if len(rl.VMs) > 0 {
				fmt.Println("VMs:")
				for _, r := range rl.VMs {
					fmt.Printf("  %-30s %s\n", r.Name, r.Status)
				}
			}
			if len(rl.Disks) > 0 {
				fmt.Println("Disks:")
				for _, r := range rl.Disks {
					fmt.Printf("  %-30s %s\n", r.Name, r.Status)
				}
			}
			if len(rl.Networks) > 0 {
				fmt.Println("Networks:")
				for _, r := range rl.Networks {
					fmt.Printf("  %-30s %s\n", r.Name, r.Status)
				}
			}
			return nil
		},
	}
}
