package hyperv

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// BuildByStruct creates VMs, disks, and switches from a Summarize config.
func BuildByStruct(summarize Summarize) error {
	// Disk aliases referenced by count>1 VMs are created per VM instance below,
	// not here, so skip them in this standalone pass.
	multiRefs := multiCountDiskRefs(summarize)

	for alias, disk := range summarize.Disks {
		if disk.Import || multiRefs[alias] {
			continue
		}
		if err := CreateDisk(disk, true); err != nil {
			return err
		}
	}

	for key, network := range summarize.Networks {
		network.Name = key
		if GetSwitchType(network.Name) == vmStateNotFound {
			if err := CreateSwitch(network, true); err != nil {
				return err
			}
		}
	}

	for key, vm := range summarize.Vms {
		vm.Name = key
		if vm.Count == 0 {
			vm.Count = 1
		}
		if vm.Count > maxVMCount {
			return fmt.Errorf("VM %s: count %d exceeds maximum of %d", key, vm.Count, maxVMCount)
		}
		for i := 1; i <= vm.Count; i++ {
			named := vm
			if vm.Count != 1 {
				named.Name = vm.Name + strconv.Itoa(i)
			}
			diskPaths, err := resolveAndCreateDisks(summarize, vm.Disks, i, vm.Count)
			if err != nil {
				return err
			}
			named.Disks = diskPaths
			if err := CreateVM(named, true); err != nil {
				return err
			}
		}
	}
	return nil
}

// RemoveByStruct removes VMs, switches, and disks defined in a Summarize config.
// Resources that do not exist are silently skipped.
// Partial failures are written to w; the last non-not-found error encountered is returned.
func RemoveByStruct(summarize Summarize, w io.Writer) error {
	var (
		lastErr error
		nfe     *notFoundError
	)
	skipIfNotFound := func(err error) {
		if err == nil || errors.As(err, &nfe) {
			return
		}
		_, _ = fmt.Fprintf(w, "%s\n", err)
		lastErr = err
	}

	for key, vm := range summarize.Vms {
		count := vm.Count
		if count == 0 {
			count = 1
		}
		if count == 1 {
			skipIfNotFound(RemoveVM(key, false))
			continue
		}
		for i := 1; i <= count; i++ {
			skipIfNotFound(RemoveVM(key+strconv.Itoa(i), false))
		}
	}
	for key, network := range summarize.Networks {
		network.Name = key
		skipIfNotFound(RemoveSwitch(network.Name))
	}

	multiRefs := multiCountDiskRefs(summarize)
	for alias, disk := range summarize.Disks {
		if disk.Import {
			continue
		}
		if multiRefs[alias] {
			// Remove per-VM numbered copies created by BuildByStruct.
			count := maxCountForDiskRef(summarize, alias)
			for i := 1; i <= count; i++ {
				skipIfNotFound(RemoveDisk(numberPath(disk.Path, i), false))
			}
			continue
		}
		skipIfNotFound(RemoveDisk(disk.Path, false))
	}
	return lastErr
}

// resolveDisks replaces disk aliases with actual paths from the config.
// Used by the cmd layer for single-disk operations.
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

// resolveAndCreateDisks resolves disk aliases to paths. For count>1 VMs,
// non-import disks are duplicated with a per-instance numeric suffix and created.
func resolveAndCreateDisks(summarize Summarize, refs []string, index, count int) ([]string, error) {
	paths := make([]string, len(refs))
	for i, ref := range refs {
		disk, ok := summarize.Disks[ref]
		if !ok {
			// Raw path: use as-is for all instances.
			paths[i] = ref
			continue
		}
		if disk.Import || count == 1 {
			paths[i] = disk.Path
			continue
		}
		// count>1 and non-import: create a numbered disk copy for this VM instance.
		numbered := disk
		numbered.Path = numberPath(disk.Path, index)
		if err := CreateDisk(numbered, true); err != nil {
			return nil, err
		}
		paths[i] = numbered.Path
	}
	return paths, nil
}

// multiCountDiskRefs returns the set of disk aliases referenced by any VM with count>1.
func multiCountDiskRefs(summarize Summarize) map[string]bool {
	refs := make(map[string]bool)
	for _, vm := range summarize.Vms {
		if vm.Count <= 1 {
			continue
		}
		for _, ref := range vm.Disks {
			if _, ok := summarize.Disks[ref]; ok {
				refs[ref] = true
			}
		}
	}
	return refs
}

// maxCountForDiskRef returns the highest count among VMs that reference the alias.
func maxCountForDiskRef(summarize Summarize, alias string) int {
	highest := 0
	for _, vm := range summarize.Vms {
		for _, ref := range vm.Disks {
			if ref == alias && vm.Count > highest {
				highest = vm.Count
			}
		}
	}
	return highest
}

// numberPath inserts n before the last file extension.
// e.g. "C:\disks\test.vhdx", 2 → "C:\disks\test2.vhdx".

func numberPath(path string, n int) string {
	dot := strings.LastIndex(path, ".")
	if dot < 0 {
		return path + strconv.Itoa(n)
	}
	return path[:dot] + strconv.Itoa(n) + path[dot:]
}
