package hyperv

import (
	"fmt"
	"io"
	"strconv"
)

// DryRunBuild prints what BuildByStruct would do without making any changes.
func DryRunBuild(summarize Summarize, w io.Writer) {
	_, _ = fmt.Fprintln(w, "[dry-run] build")

	for name, disk := range summarize.Disks {
		if disk.Import {
			dryRunLine(w, "disk", name, "skip", "import: true")
			continue
		}
		if err := checkDiskType(disk.Type); err != nil {
			dryRunLine(w, "disk", name, "error", err.Error())
			continue
		}
		if err := checkDiskSize(disk.Size); err != nil {
			dryRunLine(w, "disk", name, "error", err.Error())
			continue
		}
		if isNotFileExist(disk.Path) != nil {
			dryRunLine(w, "disk", name, "skip", "already exists: "+disk.Path)
			continue
		}
		dryRunLine(w, "disk", name, "create", disk.Path+" ("+disk.Size+", "+disk.Type+")")
	}

	for name, network := range summarize.Networks {
		network.Name = name
		if err := checkSwitchType(network.Type); err != nil {
			dryRunLine(w, "network", name, "error", err.Error())
			continue
		}
		if GetSwitchType(name) != vmStateNotFound {
			dryRunLine(w, "network", name, "skip", "already exists")
			continue
		}
		detail := network.Type
		if network.ExternalInterface != "" {
			detail += " via " + network.ExternalInterface
		}
		dryRunLine(w, "network", name, "create", detail)
	}

	for name, vm := range summarize.Vms {
		vm.Name = name
		if err := checkVMGeneration(vm.Generation); err != nil {
			dryRunLine(w, "vm", name, "error", err.Error())
			continue
		}
		if err := checkMemorySize(vm.Memory.Size); err != nil {
			dryRunLine(w, "vm", name, "error", err.Error())
			continue
		}
		if err := checkVMProcessor(vm.CPU); err != nil {
			dryRunLine(w, "vm", name, "error", err.Error())
			continue
		}

		count := vm.Count
		if count == 0 {
			count = 1
		}
		for i := 1; i <= count; i++ {
			vmName := name
			if count != 1 {
				vmName = name + strconv.Itoa(i)
			}
			if GetVMState(vmName) != vmStateNotFound {
				dryRunLine(w, "vm", vmName, "skip", "already exists")
				continue
			}
			detail := fmt.Sprintf("gen%d, %s, %d vCPU", vm.Generation, vm.Memory.Size, vm.CPU.Thread)
			dryRunLine(w, "vm", vmName, "create", detail)
		}
	}
}

// DryRunRemove prints what RemoveByStruct would do without making any changes.
func DryRunRemove(summarize Summarize, w io.Writer) {
	_, _ = fmt.Fprintln(w, "[dry-run] remove")

	for name, vm := range summarize.Vms {
		count := vm.Count
		if count == 0 {
			count = 1
		}
		for i := 1; i <= count; i++ {
			vmName := name
			if count != 1 {
				vmName = name + strconv.Itoa(i)
			}
			state := GetVMState(vmName)
			if state == vmStateNotFound {
				dryRunLine(w, "vm", vmName, "skip", "not found")
			} else {
				dryRunLine(w, "vm", vmName, "remove", "state: "+state)
			}
		}
	}

	for name := range summarize.Networks {
		if GetSwitchType(name) == vmStateNotFound {
			dryRunLine(w, "network", name, "skip", "not found")
		} else {
			dryRunLine(w, "network", name, "remove", "")
		}
	}

	for name, disk := range summarize.Disks {
		if disk.Import {
			dryRunLine(w, "disk", name, "skip", "import: true")
			continue
		}
		if isFileExist(disk.Path) != nil {
			dryRunLine(w, "disk", name, "skip", "not found: "+disk.Path)
		} else {
			dryRunLine(w, "disk", name, "remove", disk.Path)
		}
	}
}

func dryRunLine(w io.Writer, category, name, action, detail string) {
	if detail != "" {
		_, _ = fmt.Fprintf(w, "  %-8s  %-24s  %-8s  %s\n", category, name, action, detail)
	} else {
		_, _ = fmt.Fprintf(w, "  %-8s  %-24s  %s\n", category, name, action)
	}
}
