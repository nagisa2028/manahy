package hyperv

import "fmt"

// GetVMHardDiskDrives returns hard disk drives attached to a VM.
func GetVMHardDiskDrives(vmName string) (string, error) {
	if err := IsVMExist(vmName); err != nil {
		return "", err
	}
	res, err := outputPS(cmdGetVMHardDiskDrive + " -VMName " + ps(vmName) + " | Format-Table VMName, ControllerType, ControllerNumber, ControllerLocation, Path | Out-String")
	if err != nil {
		return "", fmt.Errorf("failed to get hard disk drives for VM %s: %w", vmName, err)
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
	return runPS(cmdAddVMHardDiskDrive + " -VMName " + ps(vmName) + " -Path " + ps(path))
}

// RemoveVMHardDiskDrive detaches a VHD from a VM by path.
func RemoveVMHardDiskDrive(vmName, path string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return runPS(cmdGetVMHardDiskDrive + " -VMName " + ps(vmName) + " | Where-Object { $_.Path -eq " + ps(path) + " } | " + cmdRemoveVMHardDiskDrive)
}
