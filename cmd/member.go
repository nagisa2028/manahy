package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newMemberCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
		Short: "manage Hyper-V Administrators group members",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need a valid subcommand")
		},
	}
	cmd.AddCommand(
		newMemberListCmd(),
		newMemberAddCmd(),
		newMemberRemoveCmd(),
	)
	return cmd
}

func newMemberListCmd() *cobra.Command {
	return &cobra.Command{
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
}

func newMemberAddCmd() *cobra.Command {
	var force bool
	c := &cobra.Command{
		Use:   "add",
		Short: "add a user to Hyper-V Administrators",
		Args:  cobra.RangeArgs(1, maxBulkArgs),
		RunE: func(_ *cobra.Command, args []string) error {
			prompt := fmt.Sprintf("Add %d user(s) to Hyper-V Administrators?", len(args))
			if !confirmAction(prompt, force) {
				return nil
			}
			for _, name := range args {
				if err := hyperv.AddGroupMember(name); err != nil {
					return err
				}
			}
			return nil
		},
	}
	c.Flags().BoolVar(&force, "force", false, "skip confirmation prompt")
	return c
}

func newMemberRemoveCmd() *cobra.Command {
	var force bool
	c := &cobra.Command{
		Use:   "remove",
		Short: "remove a user from Hyper-V Administrators",
		Args:  cobra.RangeArgs(1, maxBulkArgs),
		RunE: func(_ *cobra.Command, args []string) error {
			prompt := fmt.Sprintf("Remove %d user(s) from Hyper-V Administrators?", len(args))
			if !confirmAction(prompt, force) {
				return nil
			}
			for _, name := range args {
				if err := hyperv.RemoveGroupMember(name); err != nil {
					return err
				}
			}
			return nil
		},
	}
	c.Flags().BoolVar(&force, "force", false, "skip confirmation prompt")
	return c
}
