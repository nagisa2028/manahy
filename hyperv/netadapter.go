package hyperv

import (
	"fmt"
	"os/exec"
)

// GetVMNetworkAdapters returns network adapters attached to a VM.
func GetVMNetworkAdapters(vmName string) (string, error) {
	if err := IsVMExist(vmName); err != nil {
		return "", err
	}
	res, err := exec.Command("powershell", "-NoProfile", "Get-VMNetworkAdapter -VMName '"+vmName+"' | Format-Table Name, SwitchName, MacAddress, Status | Out-String").Output()
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
	cmd := "Add-VMNetworkAdapter -VMName '" + vmName + "'"
	if name != "" {
		cmd += " -Name '" + name + "'"
	}
	if switchName != "" {
		cmd += " -SwitchName '" + switchName + "'"
	}
	return exec.Command("powershell", "-NoProfile", cmd).Run()
}

// RemoveVMNetworkAdapter removes a named network adapter from a VM.
func RemoveVMNetworkAdapter(vmName, name string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile", "Remove-VMNetworkAdapter -VMName '"+vmName+"' -Name '"+name+"'").Run()
}

// ConnectVMNetworkAdapter connects a VM network adapter to a virtual switch.
func ConnectVMNetworkAdapter(vmName, name, switchName string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile", "Connect-VMNetworkAdapter -VMName '"+vmName+"' -Name '"+name+"' -SwitchName '"+switchName+"'").Run()
}
