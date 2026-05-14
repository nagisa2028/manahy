package hyperv

import "fmt"

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

	for name := range cfg.Vms {
		state := GetVMState(name)
		rl.VMs = append(rl.VMs, ResourceStatus{Name: name, Status: state})
	}

	for name, d := range cfg.Disks {
		path := d.Path
		status := "present"
		if err2 := isFileExist(path); err2 != nil {
			status = "missing"
		}
		rl.Disks = append(rl.Disks, ResourceStatus{Name: name, Status: fmt.Sprintf("%s (%s)", status, path)})
	}

	for name := range cfg.Networks {
		switchStatus := GetSwitchType(name)
		if switchStatus == vmStateNotFound {
			switchStatus = "missing"
		}
		rl.Networks = append(rl.Networks, ResourceStatus{Name: name, Status: switchStatus})
	}

	return rl, nil
}
