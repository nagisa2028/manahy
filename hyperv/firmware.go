package hyperv

import "fmt"

// SetVMSecureBoot enables or disables Secure Boot on a Gen2 VM.
// template specifies the Secure Boot template (e.g. "MicrosoftWindows",
// "MicrosoftUEFICertificateAuthority"); leave empty to keep the current template.
func SetVMSecureBoot(name string, enabled bool, template string) error {
	if err := IsVMExist(name); err != nil {
		return err
	}
	onOff := "On"
	if !enabled {
		onOff = "Off"
	}
	cmd := cmdSetVMFirmware + " -VMName " + ps(name) + " -EnableSecureBoot " + onOff
	if template != "" {
		cmd += " -SecureBootTemplate " + ps(template)
	}
	if err := runPS(cmd); err != nil {
		return fmt.Errorf("failed to set Secure Boot for VM %s", name)
	}
	return nil
}

// GetVMFirmware returns the firmware configuration of a Gen2 VM.
func GetVMFirmware(name string) (string, error) {
	if err := IsVMExist(name); err != nil {
		return "", err
	}
	out, err := outputPS(cmdGetVMFirmware + " -VMName " + ps(name) +
		" | Format-List SecureBoot, SecureBootTemplate, PreferredNetworkBootProtocol, BootOrder | Out-String")
	if err != nil {
		return "", fmt.Errorf("failed to get firmware info for VM %s", name)
	}
	return string(out), nil
}
