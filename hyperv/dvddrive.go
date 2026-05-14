package hyperv

import "fmt"

// GetVMDvdDrives returns DVD drives attached to a VM.
func GetVMDvdDrives(vmName string) (string, error) {
	if err := IsVMExist(vmName); err != nil {
		return "", err
	}
	res, err := outputPS("Get-VMDvdDrive -VMName " + ps(vmName) + " | Format-Table VMName, ControllerType, ControllerNumber, ControllerLocation, Path | Out-String")
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
	return runPS("Add-VMDvdDrive -VMName " + ps(vmName))
}

// RemoveVMDvdDrive removes the first DVD drive from a VM.
func RemoveVMDvdDrive(vmName string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return runPS("Get-VMDvdDrive -VMName " + ps(vmName) + " | Select-Object -First 1 | Remove-VMDvdDrive")
}

// SetVMDvdDrive sets the ISO image path on the first DVD drive of a VM.
// Pass an empty imagePath to eject the current image.
func SetVMDvdDrive(vmName, imagePath string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	pathParam := "$null"
	if imagePath != "" {
		pathParam = ps(imagePath)
	}
	script := "$dvd = Get-VMDvdDrive -VMName " + ps(vmName) + " | Select-Object -First 1; " +
		"Set-VMDvdDrive -VMName " + ps(vmName) + " -ControllerNumber $dvd.ControllerNumber " +
		"-ControllerLocation $dvd.ControllerLocation -Path " + pathParam
	return runPS(script)
}
