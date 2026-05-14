package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newVMCmd(configFile *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vm",
		Short: "management vm on Hyper-V",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need valid command")
		},
	}
	cmd.AddCommand(
		newVMListCmd(),
		newVMStateCmd(),
		newVMCreateCmd(),
		newVMRemoveCmd(),
		newVMRenameCmd(),
		newVMStartCmd(),
		newVMSaveCmd(),
		newVMShutdownCmd(),
		newVMDestroyCmd(),
		newVMSuspendCmd(),
		newVMResumeCmd(),
		newVMRestartCmd(),
		newVMConnectCmd(),
		newVMExportCmd(),
		newVMImportCmd(),
		newVMMoveCmd(),
		newVMCopyCmd(),
		newVMInfoCmd(),
		newVMMeasureCmd(),
		newVMIntegrationCmd(),
		newVMNicCmd(),
		newVMHardDiskCmd(configFile),
		newVMDvdCmd(configFile),
	)
	return cmd
}

func newVMListCmd() *cobra.Command {
	opts := struct {
		active   bool
		saved    bool
		inactive bool
		paused   bool
		all      bool
	}{active: true}
	c := &cobra.Command{
		Use:   "list",
		Short: "Print VM list",
		Args:  cobra.RangeArgs(0, 0),
		RunE: func(_ *cobra.Command, _ []string) error {
			if opts.saved || opts.inactive || opts.paused || opts.all {
				opts.active = false
			}
			vmList, err := hyperv.GetVMList()
			if err != nil {
				return err
			}
			if opts.active || opts.all {
				displayList(vmList.Running, "Running VM's")
			}
			if opts.saved || opts.all {
				displayList(vmList.Saved, "Saved VM's")
			}
			if opts.paused || opts.all {
				displayList(vmList.Paused, "Paused VM's")
			}
			if opts.inactive || opts.all {
				displayList(vmList.Off, "Inactive VM's")
			}
			return nil
		},
	}
	c.Flags().BoolVarP(&opts.active, "active", "", true, "list active vm's")
	c.Flags().BoolVarP(&opts.inactive, "inactive", "i", false, "list inactive vm's")
	c.Flags().BoolVarP(&opts.saved, "saved", "s", false, "list saved vm's")
	c.Flags().BoolVarP(&opts.paused, "paused", "p", false, "list paused vm's")
	c.Flags().BoolVarP(&opts.all, "all", "a", false, "list all vm's")
	return c
}

func newVMStateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "state",
		Short: "Print VM state",
		Args:  cobra.RangeArgs(1, 1),
		Run: func(_ *cobra.Command, args []string) {
			fmt.Println(hyperv.GetVMState(args[0]))
		},
	}
}

func newVMCreateCmd() *cobra.Command {
	var vm hyperv.VM
	var vmDisk, vmSwitch string
	c := &cobra.Command{
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
	c.Flags().StringVarP(&vm.Name, "name", "n", "", "new vm name")
	c.Flags().IntVarP(&vm.Generation, "generation", "g", 1, "set vm generation")
	c.Flags().IntVarP(&vm.CPU.Thread, "vcpus", "v", 1, "set vm vcpus")
	c.Flags().BoolVarP(&vm.CPU.Nested, "nested", "", false, "enable nested virtualization")
	c.Flags().StringVarP(&vm.Memory.Size, "memory", "m", "", "set vm memory")
	c.Flags().BoolVarP(&vm.Memory.Dynamic, "nodynamic", "", false, "disable dynamic memory")
	c.Flags().StringVarP(&vm.Path, "path", "p", "", "new vm path")
	c.Flags().StringVarP(&vm.Image, "image", "i", "", "image path")
	c.Flags().StringVarP(&vmDisk, "disk", "d", "", "disk path")
	c.Flags().StringVarP(&vmSwitch, "network", "s", "", "switch name")
	return c
}

func newVMRemoveCmd() *cobra.Command {
	return &cobra.Command{
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
}

func newVMRenameCmd() *cobra.Command {
	var newName string
	c := &cobra.Command{
		Use:   "rename",
		Short: "rename VM",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.RenameVM(args[0], newName)
		},
	}
	c.Flags().StringVarP(&newName, "new-name", "n", "", "new vm name")
	_ = c.MarkFlagRequired("new-name")
	return c
}

func newVMStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start VM",
		Args:  cobra.RangeArgs(1, 100),
		RunE: func(_ *cobra.Command, args []string) error {
			for _, name := range args {
				if err := hyperv.StartVM(name); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func newVMSaveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "save",
		Short: "save VM",
		Args:  cobra.RangeArgs(1, 100),
		RunE: func(_ *cobra.Command, args []string) error {
			for _, name := range args {
				if err := hyperv.SaveVM(name); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func newVMShutdownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "shutdown",
		Short: "shutdown VM",
		Args:  cobra.RangeArgs(1, 100),
		RunE: func(_ *cobra.Command, args []string) error {
			for _, name := range args {
				if err := hyperv.StopVM(name); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func newVMDestroyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "destroy",
		Short: "destroy VM",
		Args:  cobra.RangeArgs(1, 100),
		RunE: func(_ *cobra.Command, args []string) error {
			for _, name := range args {
				if err := hyperv.DestroyVM(name); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func newVMSuspendCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "suspend",
		Short: "suspend VM",
		Args:  cobra.RangeArgs(1, 100),
		RunE: func(_ *cobra.Command, args []string) error {
			for _, name := range args {
				if err := hyperv.SuspendVM(name); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func newVMResumeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resume",
		Short: "resume paused VM",
		Args:  cobra.RangeArgs(1, 100),
		RunE: func(_ *cobra.Command, args []string) error {
			for _, name := range args {
				if err := hyperv.ResumeVM(name); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func newVMRestartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "restart VM",
		Args:  cobra.RangeArgs(1, 100),
		RunE: func(_ *cobra.Command, args []string) error {
			for _, name := range args {
				if err := hyperv.RestartVM(name); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func newVMConnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connect",
		Short: "connect VM",
		Args:  cobra.RangeArgs(1, 100),
		RunE: func(_ *cobra.Command, args []string) error {
			for _, name := range args {
				if err := hyperv.ConnectVM(name); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func newVMExportCmd() *cobra.Command {
	var path string
	c := &cobra.Command{
		Use:   "export",
		Short: "export VM to a directory",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.ExportVM(args[0], path)
		},
	}
	c.Flags().StringVarP(&path, "path", "p", "", "export destination directory")
	_ = c.MarkFlagRequired("path")
	return c
}

func newVMImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import",
		Short: "import VM from a .vmcx file",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.ImportVM(args[0])
		},
	}
}

func newVMMoveCmd() *cobra.Command {
	var path string
	c := &cobra.Command{
		Use:   "move",
		Short: "move VM storage to another directory",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.MoveVMStorage(args[0], path)
		},
	}
	c.Flags().StringVarP(&path, "path", "p", "", "destination storage directory")
	_ = c.MarkFlagRequired("path")
	return c
}

func newVMCopyCmd() *cobra.Command {
	var newName, path string
	c := &cobra.Command{
		Use:   "copy",
		Short: "copy VM with a new name",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return hyperv.CopyVM(args[0], newName, path)
		},
	}
	c.Flags().StringVarP(&newName, "new-name", "n", "", "new vm name")
	c.Flags().StringVarP(&path, "path", "p", "", "directory for exported source files")
	_ = c.MarkFlagRequired("new-name")
	_ = c.MarkFlagRequired("path")
	return c
}

func newVMInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "show VM processor and memory configuration",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			info, err := hyperv.GetVMInfo(args[0])
			if err != nil {
				return err
			}
			fmt.Print(info)
			return nil
		},
	}
}

func newVMMeasureCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "measure",
		Short: "show VM resource usage metrics",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			metrics, err := hyperv.MeasureVM(args[0])
			if err != nil {
				return err
			}
			fmt.Print(metrics)
			return nil
		},
	}
}

func newVMIntegrationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "integration",
		Short: "manage VM integration services",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need valid command")
		},
	}
	cmd.AddCommand(
		newVMIntegrationListCmd(),
		newVMIntegrationEnableCmd(),
		newVMIntegrationDisableCmd(),
	)
	return cmd
}

func newVMIntegrationListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list integration services for a VM",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			out, err := hyperv.GetVMIntegrationServices(args[0])
			if err != nil {
				return err
			}
			fmt.Print(out)
			return nil
		},
	}
}

func newVMIntegrationEnableCmd() *cobra.Command {
	var name string
	c := &cobra.Command{
		Use:   "enable",
		Short: "enable an integration service on a VM",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			return hyperv.EnableVMIntegrationService(args[0], name)
		},
	}
	c.Flags().StringVarP(&name, "name", "n", "", "integration service name")
	return c
}

func newVMIntegrationDisableCmd() *cobra.Command {
	var name string
	c := &cobra.Command{
		Use:   "disable",
		Short: "disable an integration service on a VM",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			return hyperv.DisableVMIntegrationService(args[0], name)
		},
	}
	c.Flags().StringVarP(&name, "name", "n", "", "integration service name")
	return c
}
