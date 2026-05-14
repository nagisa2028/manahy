package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newHostCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "check Hyper-V environment and command availability",
		RunE: func(_ *cobra.Command, _ []string) error {
			printCheckSection("System", hyperv.CheckSystem())
			printCheckSection("VM Commands", hyperv.CheckVMCommands())
			printCheckSection("Disk Commands", hyperv.CheckDiskCommands())
			printCheckSection("Network Commands", hyperv.CheckNetworkCommands())
			return nil
		},
	}
	cmd.AddCommand(
		newHostCheckVMCmd(),
		newHostCheckDiskCmd(),
		newHostCheckNetworkCmd(),
	)
	return cmd
}

func newHostCheckVMCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "vm",
		Short: "check VM-related command availability",
		RunE: func(_ *cobra.Command, _ []string) error {
			printCheckSection("System", hyperv.CheckSystem())
			printCheckSection("VM Commands", hyperv.CheckVMCommands())
			return nil
		},
	}
}

func newHostCheckDiskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disk",
		Short: "check disk-related command availability",
		RunE: func(_ *cobra.Command, _ []string) error {
			printCheckSection("System", hyperv.CheckSystem())
			printCheckSection("Disk Commands", hyperv.CheckDiskCommands())
			return nil
		},
	}
}

func newHostCheckNetworkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "network",
		Short: "check network-related command availability",
		RunE: func(_ *cobra.Command, _ []string) error {
			printCheckSection("System", hyperv.CheckSystem())
			printCheckSection("Network Commands", hyperv.CheckNetworkCommands())
			return nil
		},
	}
}

func printCheckSection(title string, results []hyperv.CheckResult) {
	fmt.Printf("=== %s ===\n", title)
	for _, r := range results {
		var icon string
		switch r.Status {
		case hyperv.CheckOK:
			icon = "[OK]  "
		case hyperv.CheckWarn:
			icon = "[WARN]"
		default:
			icon = "[FAIL]"
		}
		fmt.Printf("  %s  %-44s  %s\n", icon, r.Name, r.Message)
	}
	fmt.Println()
}
