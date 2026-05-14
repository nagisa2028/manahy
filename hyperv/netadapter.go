package hyperv

import "fmt"

// GetVMNetworkAdapters returns network adapters attached to a VM.
func GetVMNetworkAdapters(vmName string) (string, error) {
	if err := IsVMExist(vmName); err != nil {
		return "", err
	}
	res, err := outputPS(cmdGetVMNetworkAdapter + " -VMName " + ps(vmName) + " | Format-Table Name, SwitchName, MacAddress, Status | Out-String")
	if err != nil {
		return "", fmt.Errorf("failed to get network adapters for VM %s", vmName)
	}
	return string(res), nil
}

// AddVMNetworkAdapter adds a network adapter to a VM.
func AddVMNetworkAdapter(vmName, name, switchName string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	cmd := cmdAddVMNetworkAdapter + " -VMName " + ps(vmName)
	if name != "" {
		cmd += " -Name " + ps(name)
	}
	if switchName != "" {
		cmd += " -SwitchName " + ps(switchName)
	}
	return runPS(cmd)
}

// RemoveVMNetworkAdapter removes a named network adapter from a VM.
func RemoveVMNetworkAdapter(vmName, name string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return runPS(cmdRemoveVMNetworkAdapter + " -VMName " + ps(vmName) + " -Name " + ps(name))
}

// ConnectVMNetworkAdapter connects a VM network adapter to a virtual switch.
func ConnectVMNetworkAdapter(vmName, name, switchName string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return runPS(cmdConnectVMNetworkAdapter + " -VMName " + ps(vmName) + " -Name " + ps(name) + " -SwitchName " + ps(switchName))
}

// SetVMNetworkAdapterVlan assigns an access VLAN ID to a VM network adapter.
// Pass vlanID=0 to remove VLAN tagging (untagged mode).
func SetVMNetworkAdapterVlan(vmName, name string, vlanID int) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	var cmd string
	if vlanID == 0 {
		cmd = cmdSetVMNetworkAdapterVlan + " -VMName " + ps(vmName) + " -VMNetworkAdapterName " + ps(name) + " -Untagged"
	} else {
		cmd = cmdSetVMNetworkAdapterVlan + " -VMName " + ps(vmName) + " -VMNetworkAdapterName " + ps(name) +
			" -Access -VlanId " + fmt.Sprintf("%d", vlanID)
	}
	if err := runPS(cmd); err != nil {
		return fmt.Errorf("failed to set VLAN for adapter %s on VM %s", name, vmName)
	}
	return nil
}
