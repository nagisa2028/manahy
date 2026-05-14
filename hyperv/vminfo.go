package hyperv

import "fmt"

// GetVMInfo returns processor and memory configuration for a VM.
func GetVMInfo(name string) (string, error) {
	if err := IsVMExist(name); err != nil {
		return "", err
	}
	script := cmdGetVMProcessor + " -VMName " + ps(name) + " | Format-List VMName, Count, ExposeVirtualizationExtensions; " +
		cmdGetVMMemory + " -VMName " + ps(name) + " | Format-List VMName, DynamicMemoryEnabled, Startup, Minimum, Maximum, Buffer, Priority"
	out, err := outputPS(script)
	if err != nil {
		return "", fmt.Errorf("failed to get info for VM %s", name)
	}
	return string(out), nil
}

// MeasureVM returns resource usage metrics for a VM.
// Resource metering is enabled automatically if not already active.
func MeasureVM(name string) (string, error) {
	if err := IsVMExist(name); err != nil {
		return "", err
	}
	script := cmdEnableVMResourceMetering + " -VMName " + ps(name) + "; " +
		cmdMeasureVM + " -VMName " + ps(name) + " | Format-List AvgCPUUsage, AvgRAMUsage, TotalDisk, NetworkMeteredTrafficReport"
	out, err := outputPS(script)
	if err != nil {
		return "", fmt.Errorf("failed to measure VM %s", name)
	}
	return string(out), nil
}

// GetVMIntegrationServices returns a formatted table of integration services for a VM.
func GetVMIntegrationServices(name string) (string, error) {
	if err := IsVMExist(name); err != nil {
		return "", err
	}
	out, err := outputPS(cmdGetVMIntegrationService + " -VMName " + ps(name) + " | Format-Table Name, Enabled, PrimaryStatusDescription | Out-String")
	if err != nil {
		return "", fmt.Errorf("failed to get integration services for VM %s", name)
	}
	return string(out), nil
}

// EnableVMIntegrationService enables a named integration service on a VM.
func EnableVMIntegrationService(name, service string) error {
	if err := IsVMExist(name); err != nil {
		return err
	}
	return runPS(cmdEnableVMIntegrationService + " -VMName " + ps(name) + " -Name " + ps(service))
}

// DisableVMIntegrationService disables a named integration service on a VM.
func DisableVMIntegrationService(name, service string) error {
	if err := IsVMExist(name); err != nil {
		return err
	}
	return runPS(cmdDisableVMIntegrationService + " -VMName " + ps(name) + " -Name " + ps(service))
}

// GetVMHost returns Hyper-V host configuration.
func GetVMHost() (string, error) {
	out, err := outputPS(cmdGetVMHost + " | Format-List VirtualHardDiskPath, VirtualMachinePath, MacAddressMinimum, MacAddressMaximum, NumaSpanningEnabled, EnableEnhancedSessionMode, MaximumVirtualMachineMigrations, MaximumStorageMigrations")
	if err != nil {
		return "", fmt.Errorf("failed to get Hyper-V host information")
	}
	return string(out), nil
}
