package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "management vm on Hyper-V",
	RunE: func(_ *cobra.Command, _ []string) error {
		return fmt.Errorf("need valid command")
	},
}

var vmList = &cobra.Command{
	Use:   "list",
	Short: "Print VM list",
	Args:  cobra.RangeArgs(0, 0),
	RunE: func(_ *cobra.Command, _ []string) error {
		if vmListOption.saved || vmListOption.inactive || vmListOption.paused || vmListOption.all {
			vmListOption.active = false
		}

		vmList, err := hyperv.GetVMList()
		if err != nil {
			return err
		}

		if vmListOption.active || vmListOption.all {
			displayList(vmList.Running, "Running VM's")
		}
		if vmListOption.saved || vmListOption.all {
			displayList(vmList.Saved, "Saved VM's")
		}
		if vmListOption.paused || vmListOption.all {
			displayList(vmList.Paused, "Paused VM's")
		}
		if vmListOption.inactive || vmListOption.all {
			displayList(vmList.Off, "Inactive VM's")
		}

		return nil
	},
}

var vmState = &cobra.Command{
	Use:   "state",
	Short: "Print VM state",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(_ *cobra.Command, args []string) {
		fmt.Println(hyperv.GetVMState(args[0]))
	},
}

var vmCreate = &cobra.Command{
	Use:   "create",
	Short: "create VM",
	Args:  cobra.RangeArgs(0, 0),
	RunE: func(_ *cobra.Command, _ []string) error {
		vm.Memory.Dynamic = !vm.Memory.Dynamic
		if vmDisk != "" {
			vm.Disks = append(vm.Disks, vmDisk)
		}
		if vmSwitch != "" {
			vm.Networks = append(vm.Networks, vmSwitch)
		}
		return hyperv.CreateVM(vm, true)
	},
}

var vmRemove = &cobra.Command{
	Use:   "remove",
	Short: "remove VM",
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(_ *cobra.Command, args []string) error {
		for _, name := range args {
			if err := hyperv.RemoveVM(name, false); err != nil {
				return err
			}
		}
		return nil
	},
}

var vmRename = &cobra.Command{
	Use:   "rename",
	Short: "rename VM",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(_ *cobra.Command, args []string) error {
		if newVMName == "" {
			return fmt.Errorf("--new-name is required")
		}
		return hyperv.RenameVM(args[0], newVMName)
	},
}

var vmStart = &cobra.Command{
	Use:   "start",
	Short: "start VM",
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(_ *cobra.Command, args []string) error {
		for _, vmName := range args {
			if err := hyperv.StartVM(vmName); err != nil {
				return err
			}
		}
		return nil
	},
}

var vmSave = &cobra.Command{
	Use:   "save",
	Short: "save VM",
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(_ *cobra.Command, args []string) error {
		for _, vmName := range args {
			if err := hyperv.SaveVM(vmName); err != nil {
				return err
			}
		}
		return nil
	},
}

var vmShutdown = &cobra.Command{
	Use:   "shutdown",
	Short: "shutdown VM",
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(_ *cobra.Command, args []string) error {
		for _, vmName := range args {
			if err := hyperv.StopVM(vmName); err != nil {
				return err
			}
		}
		return nil
	},
}

var vmDestroy = &cobra.Command{
	Use:   "destroy",
	Short: "destroy VM",
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(_ *cobra.Command, args []string) error {
		for _, vmName := range args {
			if err := hyperv.DestroyVM(vmName); err != nil {
				return err
			}
		}
		return nil
	},
}

var vmSuspend = &cobra.Command{
	Use:   "suspend",
	Short: "suspend VM",
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(_ *cobra.Command, args []string) error {
		for _, vmName := range args {
			if err := hyperv.SuspendVM(vmName); err != nil {
				return err
			}
		}
		return nil
	},
}

var vmRestart = &cobra.Command{
	Use:   "restart",
	Short: "restart VM",
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(_ *cobra.Command, args []string) error {
		for _, vmName := range args {
			if err := hyperv.RestartVM(vmName); err != nil {
				return err
			}
		}
		return nil
	},
}

var vmConnect = &cobra.Command{
	Use:   "connect",
	Short: "connect VM",
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(_ *cobra.Command, args []string) error {
		for _, vmName := range args {
			if err := hyperv.ConnectVM(vmName); err != nil {
				return err
			}
		}
		return nil
	},
}
