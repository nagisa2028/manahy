package cmd

import (
	"fmt"
	"os"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newStackCmd(configFile *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stack",
		Short: "manage a stack of Hyper-V resources defined in a config file",
		RunE: func(_ *cobra.Command, _ []string) error {
			return fmt.Errorf("need a valid subcommand")
		},
	}
	cmd.AddCommand(
		newStackBuildCmd(configFile),
		newStackRemoveCmd(configFile),
		newStackListCmd(configFile),
		newStackStartCmd(configFile),
		newStackStopCmd(configFile),
		newStackRestartCmd(configFile),
		newStackSaveCmd(configFile),
		newStackResumeCmd(configFile),
	)
	return cmd
}

func loadStack(configFile *string) (string, hyperv.Summarize, error) {
	path := *configFile
	if path == "" {
		path = "manahy.yaml"
	}
	data, err := hyperv.UnmarshalYaml(path)
	if err != nil {
		return "", hyperv.Summarize{}, err
	}
	fmt.Fprintf(os.Stderr, "using config: %s\n", path)
	return path, data, nil
}

func newStackBuildCmd(configFile *string) *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "build",
		Short: "create VMs, disks, and switches from config file",
		RunE: func(_ *cobra.Command, _ []string) error {
			_, data, err := loadStack(configFile)
			if err != nil {
				return err
			}
			if dryRun {
				hyperv.DryRunBuild(data, os.Stdout)
				return nil
			}
			return hyperv.BuildByStruct(data)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be created without making changes")
	return cmd
}

func newStackRemoveCmd(configFile *string) *cobra.Command {
	var dryRun, force bool
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "remove VMs, disks, and switches from config file",
		RunE: func(_ *cobra.Command, _ []string) error {
			path, data, err := loadStack(configFile)
			if err != nil {
				return err
			}
			if dryRun {
				hyperv.DryRunRemove(data, os.Stdout)
				return nil
			}
			if !confirmAction("Remove all resources defined in "+path+"?", force) {
				return nil
			}
			if err := hyperv.RemoveByStruct(data, os.Stderr); err != nil {
				// Per-resource failures already written to stderr by RemoveByStruct.
				return fmt.Errorf("one or more resources could not be removed")
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be removed without making changes")
	cmd.Flags().BoolVar(&force, "force", false, "skip confirmation prompt")
	return cmd
}

func newStackListCmd(configFile *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "show status of all resources defined in config file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			path := *configFile
			if path == "" {
				path = "manahy.yaml"
			}
			rl, err := hyperv.GetResourceList(path)
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

func newStackStartCmd(configFile *string) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start all VMs defined in config file",
		RunE: func(_ *cobra.Command, _ []string) error {
			_, data, err := loadStack(configFile)
			if err != nil {
				return err
			}
			if err := hyperv.StartByStruct(data, os.Stderr); err != nil {
				// Per-VM failures already written to stderr by StartByStruct.
				return fmt.Errorf("one or more VMs failed to start")
			}
			return nil
		},
	}
}

func newStackStopCmd(configFile *string) *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop all VMs defined in config file",
		RunE: func(_ *cobra.Command, _ []string) error {
			path, data, err := loadStack(configFile)
			if err != nil {
				return err
			}
			if !confirmAction("Stop all VMs defined in "+path+"?", force) {
				return nil
			}
			if err := hyperv.StopByStruct(data, os.Stderr); err != nil {
				// Per-VM failures already written to stderr by StopByStruct.
				return fmt.Errorf("one or more VMs failed to stop")
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "skip confirmation prompt")
	return cmd
}

func newStackRestartCmd(configFile *string) *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "restart all running VMs defined in config file",
		RunE: func(_ *cobra.Command, _ []string) error {
			path, data, err := loadStack(configFile)
			if err != nil {
				return err
			}
			if !confirmAction("Restart all running VMs defined in "+path+"?", force) {
				return nil
			}
			if err := hyperv.RestartByStruct(data, os.Stderr); err != nil {
				// Per-VM failures already written to stderr by RestartByStruct.
				return fmt.Errorf("one or more VMs failed to restart")
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "skip confirmation prompt")
	return cmd
}

func newStackSaveCmd(configFile *string) *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "save",
		Short: "save state of all running VMs defined in config file",
		RunE: func(_ *cobra.Command, _ []string) error {
			path, data, err := loadStack(configFile)
			if err != nil {
				return err
			}
			if !confirmAction("Save state of all running VMs defined in "+path+"?", force) {
				return nil
			}
			if err := hyperv.SaveByStruct(data, os.Stderr); err != nil {
				// Per-VM failures already written to stderr by SaveByStruct.
				return fmt.Errorf("one or more VMs failed to save")
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "skip confirmation prompt")
	return cmd
}

func newStackResumeCmd(configFile *string) *cobra.Command {
	return &cobra.Command{
		Use:   "resume",
		Short: "resume all saved VMs defined in config file",
		RunE: func(_ *cobra.Command, _ []string) error {
			_, data, err := loadStack(configFile)
			if err != nil {
				return err
			}
			if err := hyperv.ResumeByStruct(data, os.Stderr); err != nil {
				// Per-VM failures already written to stderr by ResumeByStruct.
				return fmt.Errorf("one or more VMs failed to resume")
			}
			return nil
		},
	}
}
