package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newVMNicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nic",
		Short: "manage VM network adapters",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need a valid subcommand")
		},
	}
	cmd.AddCommand(
		newVMNicListCmd(),
		newVMNicAddCmd(),
		newVMNicRemoveCmd(),
		newVMNicConnectCmd(),
		newVMNicVlanCmd(),
	)
	return cmd
}

func newVMNicListCmd() *cobra.Command {
	return &cobra.Command{
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
}

func newVMNicAddCmd() *cobra.Command {
	var name, switchName string
	c := &cobra.Command{
		Use:   "add",
		Short: "add a network adapter to a VM",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.AddVMNetworkAdapter(args[0], name, switchName)
		},
	}
	c.Flags().StringVarP(&name, "name", "n", "", "adapter name")
	c.Flags().StringVarP(&switchName, "switch", "s", "", "switch name")
	return c
}

func newVMNicRemoveCmd() *cobra.Command {
	var name string
	c := &cobra.Command{
		Use:   "remove",
		Short: "remove a network adapter from a VM",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.RemoveVMNetworkAdapter(args[0], name)
		},
	}
	c.Flags().StringVarP(&name, "name", "n", "", "adapter name")
	_ = c.MarkFlagRequired("name")
	return c
}

func newVMNicVlanCmd() *cobra.Command {
	var name string
	var vlanID int
	c := &cobra.Command{
		Use:   "vlan",
		Short: "set VLAN ID on a network adapter (0 = untagged)",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.SetVMNetworkAdapterVlan(args[0], name, vlanID)
		},
	}
	c.Flags().StringVarP(&name, "name", "n", "", "adapter name")
	_ = c.MarkFlagRequired("name")
	c.Flags().IntVar(&vlanID, "vlan", 0, "VLAN ID (0 = untagged/remove VLAN tag)")
	return c
}

func newVMNicConnectCmd() *cobra.Command {
	var name, switchName string
	c := &cobra.Command{
		Use:   "connect",
		Short: "connect a network adapter to a switch",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.ConnectVMNetworkAdapter(args[0], name, switchName)
		},
	}
	c.Flags().StringVarP(&name, "name", "n", "", "adapter name")
	c.Flags().StringVarP(&switchName, "switch", "s", "", "switch name")
	_ = c.MarkFlagRequired("name")
	_ = c.MarkFlagRequired("switch")
	return c
}
