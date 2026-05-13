package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "management switch on Hyper-V",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("need valid command")
	},
}

var switchList = &cobra.Command{
	Use:   "list",
	Short: "Print switch list",
	Args:  cobra.RangeArgs(0, 0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if switchListOption.external || switchListOption.internal || switchListOption.private {
			switchListOption.all = false
		}

		switchList, err := hyperv.GetSwitchList()
		if err != nil {
			return err
		}

		if switchListOption.external || switchListOption.all {
			displayList(switchList.External, "External Switch's")
		}
		if switchListOption.internal || switchListOption.all {
			displayList(switchList.Internal, "Internal Switch's")
		}
		if switchListOption.private || switchListOption.all {
			displayList(switchList.Private, "Private Switch's")
		}
		return nil
	},
}

var switchCreate = &cobra.Command{
	Use:   "create",
	Short: "Create switch",
	Args:  cobra.RangeArgs(0, 0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if switchCreateOption.Name == "" || switchCreateOption.Type == "" {
			return fmt.Errorf("error: Please specify switch name and switch type")
		} else if switchCreateOption.Type == "external" && switchCreateOption.ExternalInterface == "" {
			return fmt.Errorf("error: Please specify an external interface")
		}
		return hyperv.CreateSwitch(switchCreateOption, true)
	},
}

var switchRemove = &cobra.Command{
	Use:   "remove",
	Short: "Remove switch",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return hyperv.RemoveSwitch(args[0])
	},
}

var switchRename = &cobra.Command{
	Use:   "rename",
	Short: "Rename switch",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if newSwitchName == "" {
			return fmt.Errorf("error: need new switch name")
		}
		return hyperv.RenameSwitch(args[0], newSwitchName)
	},
}

var switchOptionCfgCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure switch option",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("need valid command")
	},
}

var switchOptionCfgType = &cobra.Command{
	Use:   "type",
	Short: "Configure switch type",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return hyperv.ChangeSwitchType(args[0], switchType)
	},
}

var switchOptionCfgNetAdapter = &cobra.Command{
	Use:   "adapter",
	Short: "Configure network adapter",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return hyperv.ChangeSwitchNetAdapter(args[0], netAdapter)
	},
}
