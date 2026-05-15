package hyperv

import (
	"fmt"
	"sort"
	"strconv"
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
func GetResourceList(configFile string) (ResourceList, error) {
	cfg, err := UnmarshalYaml(configFile)
	if err != nil {
		return ResourceList{}, err
	}

	var rl ResourceList

	for name, vm := range cfg.Vms {
		count := vm.Count
		if count == 0 {
			count = 1
		}
		if count == 1 {
			rl.VMs = append(rl.VMs, ResourceStatus{Name: name, Status: GetVMState(name)})
		} else {
			for i := 1; i <= count; i++ {
				instanceName := name + strconv.Itoa(i)
				rl.VMs = append(rl.VMs, ResourceStatus{Name: instanceName, Status: GetVMState(instanceName)})
			}
		}
	}

	multiRefs := multiCountDiskRefs(cfg)
	for name, d := range cfg.Disks {
		if multiRefs[name] {
			count := maxCountForDiskRef(cfg, name)
			for i := 1; i <= count; i++ {
				path := numberPath(d.Path, i)
				status := diskFileStatus(path)
				rl.Disks = append(rl.Disks, ResourceStatus{
					Name:   name + strconv.Itoa(i),
					Status: fmt.Sprintf("%s (%s)", status, path),
				})
			}
		} else {
			rl.Disks = append(rl.Disks, ResourceStatus{
				Name:   name,
				Status: fmt.Sprintf("%s (%s)", diskFileStatus(d.Path), d.Path),
			})
		}
	}

	for name := range cfg.Networks {
		switchStatus := GetSwitchType(name)
		if switchStatus == vmStateNotFound {
			switchStatus = "missing"
		}
		rl.Networks = append(rl.Networks, ResourceStatus{Name: name, Status: switchStatus})
	}

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
