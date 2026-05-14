package hyperv

import "fmt"

// GetVMNetworkAdapters returns network adapters attached to a VM.
func GetVMNetworkAdapters(vmName string) (string, error) {
	if err := IsVMExist(vmName); err != nil {
		return "", err
	}
	res, err := outputPS("Get-VMNetworkAdapter -VMName " + ps(vmName) + " | Format-Table Name, SwitchName, MacAddress, Status | Out-String")
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
	cmd := "Add-VMNetworkAdapter -VMName " + ps(vmName)
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
	return runPS("Remove-VMNetworkAdapter -VMName " + ps(vmName) + " -Name " + ps(name))
}

// ConnectVMNetworkAdapter connects a VM network adapter to a virtual switch.
func ConnectVMNetworkAdapter(vmName, name, switchName string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return runPS("Connect-VMNetworkAdapter -VMName " + ps(vmName) + " -Name " + ps(name) + " -SwitchName " + ps(switchName))
}
