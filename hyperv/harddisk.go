package hyperv

import (
	"fmt"
	"os/exec"
)

// GetVMHardDiskDrives returns hard disk drives attached to a VM.
func GetVMHardDiskDrives(vmName string) (string, error) {
	if err := IsVMExist(vmName); err != nil {
		return "", err
	}
	res, err := exec.Command("powershell", "-NoProfile", "Get-VMHardDiskDrive -VMName '"+vmName+"' | Format-Table VMName, ControllerType, ControllerNumber, ControllerLocation, Path | Out-String").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get hard disk drives for VM %s", vmName)
	}
	return string(res), nil
}

// AddVMHardDiskDrive attaches a VHD to a VM.
func AddVMHardDiskDrive(vmName, path string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	if err := isFileExist(path); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile", "Add-VMHardDiskDrive -VMName '"+vmName+"' -Path '"+path+"'").Run()
}

// RemoveVMHardDiskDrive detaches a VHD from a VM by path.
func RemoveVMHardDiskDrive(vmName, path string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile", "Get-VMHardDiskDrive -VMName '"+vmName+"' | Where-Object { $_.Path -eq '"+path+"' } | Remove-VMHardDiskDrive").Run()
}
