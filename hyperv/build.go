// hyperv package is manage Hyper-V
package hyperv

import "strconv"

// BuildByStruct creates VMs, disks, and switches from a Summarize config
func BuildByStruct(summarize Summarize) error {
	for _, disk := range summarize.Disks {
		err := CreateDisk(disk, true)
		if err != nil {
			return err
		}
	}

	for _, network := range summarize.Networks {
		if GetSwitchType(network.Name) == "NotFound" {
			err := CreateSwitch(network, true)
			if err != nil {
				return err
			}
		}
	}

	for _, vm := range summarize.Vms {
		if vm.Count == 0 {
			vm.Count = 1
		}
		for i := 1; i <= vm.Count; i++ {
			if vm.Count != 1 {
				vm.Name += strconv.Itoa(i)
			}
			err := CreateVm(vm, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
