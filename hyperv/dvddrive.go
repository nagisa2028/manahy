package hyperv

import "fmt"

// GetVMDvdDrives returns DVD drives attached to a VM.
func GetVMDvdDrives(vmName string) (string, error) {
	if err := IsVMExist(vmName); err != nil {
		return "", err
	}
	res, err := outputPS(cmdGetVMDvdDrive + " -VMName " + ps(vmName) + " | Format-Table VMName, ControllerType, ControllerNumber, ControllerLocation, Path | Out-String")
	if err != nil {
		return "", fmt.Errorf("failed to get DVD drives for VM %s: %w", vmName, err)
	}
	return string(res), nil
}

// AddVMDvdDrive adds a DVD drive to a VM.
func AddVMDvdDrive(vmName string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return runPS(cmdAddVMDvdDrive + " -VMName " + ps(vmName))
}

// RemoveVMDvdDrive removes the first DVD drive from a VM.
func RemoveVMDvdDrive(vmName string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return runPS(cmdGetVMDvdDrive + " -VMName " + ps(vmName) + " | Select-Object -First 1 | " + cmdRemoveVMDvdDrive)
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
	script := "$dvd = " + cmdGetVMDvdDrive + " -VMName " + ps(vmName) + " | Select-Object -First 1; " +
		"if ($dvd -eq $null) { throw 'no DVD drive found' }; " +
		cmdSetVMDvdDrive + " -VMName " + ps(vmName) + " -ControllerNumber $dvd.ControllerNumber " +
		"-ControllerLocation $dvd.ControllerLocation -Path " + pathParam
	return runPS(script)
}
