package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

var vmNicCmd = &cobra.Command{
	Use:   "nic",
	Short: "manage VM network adapters",
	RunE: func(_ *cobra.Command, _ []string) error {
		return fmt.Errorf("need valid command")
	},
}

var vmNicList = &cobra.Command{
	Use:   "list",
	Short: "list network adapters for a VM",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		out, err := hyperv.GetVMNetworkAdapters(args[0])
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	},
}

var vmNicAdd = &cobra.Command{
	Use:   "add",
	Short: "add a network adapter to a VM",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		return hyperv.AddVMNetworkAdapter(args[0], vmNicName, vmNicSwitch)
	},
}

var vmNicRemove = &cobra.Command{
	Use:   "remove",
	Short: "remove a network adapter from a VM",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		if vmNicName == "" {
			return fmt.Errorf("--name is required")
		}
		return hyperv.RemoveVMNetworkAdapter(args[0], vmNicName)
	},
}

var vmNicConnect = &cobra.Command{
	Use:   "connect",
	Short: "connect a network adapter to a switch",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		if vmNicName == "" {
			return fmt.Errorf("--name is required")
		}
		if vmNicSwitch == "" {
			return fmt.Errorf("--switch is required")
		}
		return hyperv.ConnectVMNetworkAdapter(args[0], vmNicName, vmNicSwitch)
	},
}
