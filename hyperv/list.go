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
// VM, disk, and network queries are issued concurrently; VM and network states are
// fetched in single batch PS calls rather than one call per resource.
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
		stateMap, err := getVMStateMap()
		vms := make([]ResourceStatus, 0, len(cfg.Vms))
		for name, vm := range cfg.Vms {
			for _, instanceName := range vmInstanceNames(name, vm) {
				var state string
				if err != nil {
					state = "error"
				} else {
					var ok bool
					state, ok = stateMap[instanceName]
					if !ok {
						state = vmStateNotFound
					}
				}
				vms = append(vms, ResourceStatus{Name: instanceName, Status: state})
			}
		}
		rl.VMs = vms
	}()

	go func() {
		defer wg.Done()
		multiRefs := multiCountDiskRefs(cfg)
		disks := make([]ResourceStatus, 0, len(cfg.Disks))
		for name, d := range cfg.Disks {
			if multiRefs[name] {
				count := maxCountForDiskRef(cfg, name)
				for i := 1; i <= count; i++ {
					path := numberPath(d.Path, i)
					disks = append(disks, ResourceStatus{
						Name:   name + strconv.Itoa(i),
						Status: fmt.Sprintf("%s (%s)", diskFileStatus(path), path),
					})
				}
			} else {
				disks = append(disks, ResourceStatus{
					Name:   name,
					Status: fmt.Sprintf("%s (%s)", diskFileStatus(d.Path), d.Path),
				})
			}
		}
		rl.Disks = disks
	}()

	go func() {
		defer wg.Done()
		switchMap, err := getSwitchTypeMap()
		networks := make([]ResourceStatus, 0, len(cfg.Networks))
		for name := range cfg.Networks {
			var switchType string
			if err != nil {
				switchType = "error"
			} else {
				var ok bool
				switchType, ok = switchMap[name]
				if !ok {
					switchType = "missing"
				}
			}
			networks = append(networks, ResourceStatus{Name: name, Status: switchType})
		}
		rl.Networks = networks
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
