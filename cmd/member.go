package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

var memberCmd = &cobra.Command{
	Use:   "member",
	Short: "management Hyper-V Administrators group members",
	RunE: func(_ *cobra.Command, _ []string) error {
		return fmt.Errorf("need valid command")
	},
}

var memberListCmd = &cobra.Command{
	Use:   "list",
	Short: "show Hyper-V Administrators group members",
	RunE: func(_ *cobra.Command, _ []string) error {
		members, err := hyperv.GetGroupMember()
		if err != nil {
			return err
		}

		fmt.Println("Hyper-V Administrators")
		for _, m := range members {
			fmt.Printf("- %s\n", m)
		}
		fmt.Println()
		return nil
	},
}

var memberAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add a user to Hyper-V Administrators",
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(_ *cobra.Command, args []string) error {
		for _, name := range args {
			if err := hyperv.AddGroupMember(name); err != nil {
				return err
			}
		}
		return nil
	},
}

var memberRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove a user from Hyper-V Administrators",
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(_ *cobra.Command, args []string) error {
		for _, name := range args {
			if err := hyperv.RemoveGroupMember(name); err != nil {
				return err
			}
		}
		return nil
	},
}
