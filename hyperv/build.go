package hyperv

import (
	"fmt"
	"io"
	"strconv"
)

// BuildByStruct creates VMs, disks, and switches from a Summarize config.
func BuildByStruct(summarize Summarize) error {
	for _, disk := range summarize.Disks {
		if disk.Import {
			continue
		}
		if err := CreateDisk(disk, true); err != nil {
			return err
		}
	}

	for key, network := range summarize.Networks {
		network.Name = key
		if GetSwitchType(network.Name) == "NotFound" {
			if err := CreateSwitch(network, true); err != nil {
				return err
			}
		}
	}

	for key, vm := range summarize.Vms {
		vm.Name = key
		vm.Disks = resolveDisks(summarize, vm.Disks)
		if vm.Count == 0 {
			vm.Count = 1
		}
		for i := 1; i <= vm.Count; i++ {
			named := vm
			if vm.Count != 1 {
				named.Name = vm.Name + strconv.Itoa(i)
			}
			if err := CreateVM(named, true); err != nil {
				return err
			}
		}
	}
	return nil
}

// RemoveByStruct removes VMs, switches, and disks defined in a Summarize config.
// Partial failures are written to w; the last error encountered is returned.
func RemoveByStruct(summarize Summarize, w io.Writer) error {
	var lastErr error
	for key, vm := range summarize.Vms {
		count := vm.Count
		if count == 0 {
			count = 1
		}
		if count == 1 {
			if err := RemoveVM(key, false); err != nil {
				_, _ = fmt.Fprintf(w, "%s\n", err)
				lastErr = err
			}
			continue
		}
		for i := 1; i <= count; i++ {
			if err := RemoveVM(key+strconv.Itoa(i), false); err != nil {
				_, _ = fmt.Fprintf(w, "%s\n", err)
				lastErr = err
			}
		}
	}
	for key, network := range summarize.Networks {
		network.Name = key
		if err := RemoveSwitch(network.Name); err != nil {
			_, _ = fmt.Fprintf(w, "%s\n", err)
			lastErr = err
		}
	}
	for _, disk := range summarize.Disks {
		if disk.Import {
			continue
		}
		if err := RemoveDisk(disk.Path, false); err != nil {
			_, _ = fmt.Fprintf(w, "%s\n", err)
			lastErr = err
		}
	}
	return lastErr
}

// resolveDisks replaces disk aliases with actual paths from the config.
func resolveDisks(summarize Summarize, refs []string) []string {
	paths := make([]string, len(refs))
	for i, ref := range refs {
		if d, ok := summarize.Disks[ref]; ok {
			paths[i] = d.Path
		} else {
			paths[i] = ref
		}
	}
	return paths
}
