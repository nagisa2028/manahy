package hyperv

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

// ResourceStatus holds the name and current status of a single resource.
type ResourceStatus struct {
	Name   string
	Status string
}

// ResourceList groups the status of all resources defined in a YAML config.
type ResourceList struct {
	VMs      []ResourceStatus
	Disks    []ResourceStatus
	Networks []ResourceStatus
}

// GetResourceList returns the live status of all resources defined in the config file.
// VM, disk, and network queries are issued concurrently.
func GetResourceList(configFile string) (ResourceList, error) {
	cfg, err := UnmarshalYaml(configFile)
	if err != nil {
		return ResourceList{}, err
	}

	var rl ResourceList
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		for name, vm := range cfg.Vms {
			for _, instanceName := range vmInstanceNames(name, vm) {
				rl.VMs = append(rl.VMs, ResourceStatus{Name: instanceName, Status: GetVMState(instanceName)})
			}
		}
	}()

	go func() {
		defer wg.Done()
		multiRefs := multiCountDiskRefs(cfg)
		for name, d := range cfg.Disks {
			if multiRefs[name] {
				count := maxCountForDiskRef(cfg, name)
				for i := 1; i <= count; i++ {
					path := numberPath(d.Path, i)
					rl.Disks = append(rl.Disks, ResourceStatus{
						Name:   name + strconv.Itoa(i),
						Status: fmt.Sprintf("%s (%s)", diskFileStatus(path), path),
					})
				}
			} else {
				rl.Disks = append(rl.Disks, ResourceStatus{
					Name:   name,
					Status: fmt.Sprintf("%s (%s)", diskFileStatus(d.Path), d.Path),
				})
			}
		}
	}()

	go func() {
		defer wg.Done()
		for name := range cfg.Networks {
			switchStatus := GetSwitchType(name)
			if switchStatus == vmStateNotFound {
				switchStatus = "missing"
			}
			rl.Networks = append(rl.Networks, ResourceStatus{Name: name, Status: switchStatus})
		}
	}()

	wg.Wait()

	sort.Slice(rl.VMs, func(i, j int) bool { return rl.VMs[i].Name < rl.VMs[j].Name })
	sort.Slice(rl.Disks, func(i, j int) bool { return rl.Disks[i].Name < rl.Disks[j].Name })
	sort.Slice(rl.Networks, func(i, j int) bool { return rl.Networks[i].Name < rl.Networks[j].Name })

	return rl, nil
}

func diskFileStatus(path string) string {
	if err := isFileExist(path); err != nil {
		return "missing"
	}
	return "present"
}
