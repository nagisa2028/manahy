package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newSwitchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch",
		Short: "manage virtual switches",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need a valid subcommand")
		},
	}
	cmd.AddCommand(
		newSwitchListCmd(),
		newSwitchCreateCmd(),
		newSwitchRemoveCmd(),
		newSwitchRenameCmd(),
		newSwitchConfigureCmd(),
	)
	return cmd
}

func newSwitchListCmd() *cobra.Command {
	opts := struct {
		external bool
		internal bool
		private  bool
		all      bool
	}{all: true}
	c := &cobra.Command{
		Use:   "list",
		Short: "Print switch list",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if opts.external || opts.internal || opts.private {
				opts.all = false
			}
			switchList, err := hyperv.GetSwitchList()
			if err != nil {
				return err
			}
			if opts.external || opts.all {
				displayList(switchList.External, "External Switch's")
			}
			if opts.internal || opts.all {
				displayList(switchList.Internal, "Internal Switch's")
			}
			if opts.private || opts.all {
				displayList(switchList.Private, "Private Switch's")
			}
			return nil
		},
	}
	c.Flags().BoolVarP(&opts.external, "external", "e", false, "list external switches")
	c.Flags().BoolVarP(&opts.internal, "internal", "i", false, "list internal switches")
	c.Flags().BoolVarP(&opts.private, "private", "p", false, "list private switches")
	c.Flags().BoolVarP(&opts.all, "all", "a", true, "list all switches")
	return c
}

func newSwitchCreateCmd() *cobra.Command {
	var opt hyperv.VMSwitch
	c := &cobra.Command{
		Use:   "create",
		Short: "Create switch",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if opt.Name == "" || opt.Type == "" {
				return fmt.Errorf("--name and --type are required")
			}
			if opt.Type == "external" && opt.ExternalInterface == "" {
				return fmt.Errorf("--external-interface is required for external switches")
			}
			return hyperv.CreateSwitch(opt, true)
		},
	}
	c.Flags().StringVarP(&opt.Name, "name", "n", "", "set name")
	c.Flags().StringVarP(&opt.Type, "type", "t", "", "set type")
	c.Flags().StringVarP(&opt.ExternalInterface, "external-interface", "", "", "set external interface")
	c.Flags().BoolVarP(&opt.AllowManagementOS, "allow-management-os", "", false, "set allow management os")
	return c
}

func newSwitchRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove switch",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.RemoveSwitch(args[0])
		},
	}
}

func newSwitchRenameCmd() *cobra.Command {
	var newName string
	c := &cobra.Command{
		Use:   "rename",
		Short: "Rename switch",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.RenameSwitch(args[0], newName)
		},
	}
	c.Flags().StringVarP(&newName, "new-name", "n", "", "rename switch")
	_ = c.MarkFlagRequired("new-name")
	return c
}

func newSwitchConfigureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure switch option",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need a valid subcommand")
		},
	}
	cmd.AddCommand(
		newSwitchConfigureTypeCmd(),
		newSwitchConfigureAdapterCmd(),
	)
	return cmd
}

func newSwitchConfigureTypeCmd() *cobra.Command {
	var switchType string
	c := &cobra.Command{
		Use:   "type",
		Short: "Configure switch type",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.ChangeSwitchType(args[0], switchType)
		},
	}
	c.Flags().StringVarP(&switchType, "type", "t", "", "change switch type")
	return c
}

func newSwitchConfigureAdapterCmd() *cobra.Command {
	var netAdapter string
	c := &cobra.Command{
		Use:   "adapter",
		Short: "Configure network adapter",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.ChangeSwitchNetAdapter(args[0], netAdapter)
		},
	}
	c.Flags().StringVarP(&netAdapter, "net-adapter", "n", "", "change network adapter")
	return c
}
