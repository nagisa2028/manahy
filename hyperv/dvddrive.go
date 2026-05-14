package hyperv

import (
	"fmt"
	"os/exec"
)

// GetVMDvdDrives returns DVD drives attached to a VM.
func GetVMDvdDrives(vmName string) (string, error) {
	if err := IsVMExist(vmName); err != nil {
		return "", err
	}
	res, err := exec.Command("powershell", "-NoProfile", "Get-VMDvdDrive -VMName '"+vmName+"' | Format-Table VMName, ControllerType, ControllerNumber, ControllerLocation, Path | Out-String").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get DVD drives for VM %s", vmName)
	}
	return string(res), nil
}

// AddVMDvdDrive adds a DVD drive to a VM.
func AddVMDvdDrive(vmName string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile", "Add-VMDvdDrive -VMName '"+vmName+"'").Run()
}

// RemoveVMDvdDrive removes the first DVD drive from a VM.
func RemoveVMDvdDrive(vmName string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile", "Get-VMDvdDrive -VMName '"+vmName+"' | Select-Object -First 1 | Remove-VMDvdDrive").Run()
}

// SetVMDvdDrive sets the ISO image path on the first DVD drive of a VM.
func SetVMDvdDrive(vmName, imagePath string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	pathParam := "$null"
	if imagePath != "" {
		pathParam = "'" + imagePath + "'"
	}
	script := "$dvd = Get-VMDvdDrive -VMName '" + vmName + "' | Select-Object -First 1; " +
		"Set-VMDvdDrive -VMName '" + vmName + "' -ControllerNumber $dvd.ControllerNumber " +
		"-ControllerLocation $dvd.ControllerLocation -Path " + pathParam
	return exec.Command("powershell", "-NoProfile", script).Run()
}
