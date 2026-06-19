package cmd

import (
	"fmt"
	"os"

	"github.com/DevelopNaoki/manahy/hyperv"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate [config-file]",
		Short: "validate a manahy config file",
		Long: `Validate parses the given config file (default: manahy.yaml) and reports
any structural or semantic errors without contacting Hyper-V.

Exit code 0 means the config is valid; non-zero means it has errors.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			path := "manahy.yaml"
			if len(args) == 1 {
				path = args[0]
			}
			data, err := hyperv.UnmarshalYaml(path)
			if err != nil {
				return fmt.Errorf("config parse error: %w", err)
			}
			errs := validateConfig(data)
			if len(errs) == 0 {
				fmt.Fprintf(os.Stdout, "config %s is valid\n", path)
				return nil
			}
			for _, e := range errs {
				fmt.Fprintf(os.Stderr, "error: %s\n", e)
			}
			return fmt.Errorf("%d validation error(s) in %s", len(errs), path)
		},
	}
}

// validateConfig performs semantic validation on a parsed Summarize config.
// It returns a list of error strings; an empty slice means the config is valid.
func validateConfig(data hyperv.Summarize) []string {
	var errs []string

	for name, vm := range data.Vms {
		if vm.Generation < 1 || vm.Generation > 2 {
			errs = append(errs, fmt.Sprintf("vm %q: generation must be 1 or 2, got %d", name, vm.Generation))
		}
		if vm.Memory.Size == "" {
			errs = append(errs, fmt.Sprintf("vm %q: memory.size is required", name))
		} else if err := hyperv.CheckMemorySize(vm.Memory.Size); err != nil {
			errs = append(errs, fmt.Sprintf("vm %q: memory.size: %s", name, err))
		}
		if vm.Path == "" {
			errs = append(errs, fmt.Sprintf("vm %q: path is required", name))
		}
		if vm.Memory.Dynamic && vm.Memory.Buffer != 0 && (vm.Memory.Buffer < 5 || vm.Memory.Buffer > 100) {
			errs = append(errs, fmt.Sprintf("vm %q: memory.buffer must be between 5 and 100 (or 0 to use default), got %d", name, vm.Memory.Buffer))
		}
		// Verify disk aliases exist in the config.
		for _, diskRef := range vm.Disks {
			if _, ok := data.Disks[diskRef]; !ok {
				// Could be a raw absolute path — only warn if it doesn't look like one.
				if !hyperv.IsWindowsAbsPath(diskRef) {
					errs = append(errs, fmt.Sprintf("vm %q: disk ref %q not found in disks section", name, diskRef))
				}
			}
		}
		// Verify network switch aliases exist.
		for _, netRef := range vm.Networks {
			if _, ok := data.Networks[netRef]; !ok {
				// May be an existing switch defined outside the config — warn but don't error.
				_ = netRef
			}
		}
		for _, ref := range vm.NetworkRefs {
			if ref.Switch == "" {
				errs = append(errs, fmt.Sprintf("vm %q: network-refs entry has empty switch name", name))
			}
			if ref.VLAN < 0 || ref.VLAN > 4094 {
				errs = append(errs, fmt.Sprintf("vm %q: network-refs switch %q: vlan %d out of range (0–4094)", name, ref.Switch, ref.VLAN))
			}
		}
	}

	for name, disk := range data.Disks {
		if disk.Path == "" {
			errs = append(errs, fmt.Sprintf("disk %q: path is required", name))
		}
		if !disk.Import {
			if disk.Type == "" {
				errs = append(errs, fmt.Sprintf("disk %q: type is required when import is false", name))
			} else if err := hyperv.CheckDiskType(disk.Type); err != nil {
				errs = append(errs, fmt.Sprintf("disk %q: %s", name, err))
			} else if disk.Type != "differencing" {
				if disk.Size == "" {
					errs = append(errs, fmt.Sprintf("disk %q: size is required for type %s", name, disk.Type))
				} else if err := hyperv.CheckDiskSize(disk.Size); err != nil {
					errs = append(errs, fmt.Sprintf("disk %q: size: %s", name, err))
				}
			}
		}
	}

	for name, sw := range data.Networks {
		if sw.Type == "" {
			errs = append(errs, fmt.Sprintf("network %q: type is required", name))
		}
	}

	return errs
}
