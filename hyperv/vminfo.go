package hyperv

import "fmt"

// GetVMInfo returns processor and memory configuration for a VM.
func GetVMInfo(name string) (string, error) {
	if err := IsVMExist(name); err != nil {
		return "", err
	}
	procOut, err := outputPS("Get-VMProcessor -VMName " + ps(name) + " | Format-List VMName, Count, ExposeVirtualizationExtensions")
	if err != nil {
		return "", fmt.Errorf("failed to get processor info for VM %s", name)
	}
	memOut, err := outputPS("Get-VMMemory -VMName " + ps(name) + " | Format-List VMName, DynamicMemoryEnabled, Startup, Minimum, Maximum, Buffer, Priority")
	if err != nil {
		return "", fmt.Errorf("failed to get memory info for VM %s", name)
	}
	return string(procOut) + string(memOut), nil
}

// MeasureVM returns resource usage metrics for a VM.
// Resource metering is enabled automatically if not already active.
func MeasureVM(name string) (string, error) {
	if err := IsVMExist(name); err != nil {
		return "", err
	}
	script := "Enable-VMResourceMetering -VMName " + ps(name) + "; " +
		"Measure-VM -VMName " + ps(name) + " | Format-List AvgCPUUsage, AvgRAMUsage, TotalDisk, NetworkMeteredTrafficReport"
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
	out, err := outputPS("Get-VMIntegrationService -VMName " + ps(name) + " | Format-Table Name, Enabled, PrimaryStatusDescription | Out-String")
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
	return runPS("Enable-VMIntegrationService -VMName " + ps(name) + " -Name " + ps(service))
}

// DisableVMIntegrationService disables a named integration service on a VM.
func DisableVMIntegrationService(name, service string) error {
	if err := IsVMExist(name); err != nil {
		return err
	}
	return runPS("Disable-VMIntegrationService -VMName " + ps(name) + " -Name " + ps(service))
}

// GetVMHost returns Hyper-V host configuration.
func GetVMHost() (string, error) {
	out, err := outputPS("Get-VMHost | Format-List VirtualHardDiskPath, VirtualMachinePath, MacAddressMinimum, MacAddressMaximum, NumaSpanningEnabled, EnableEnhancedSessionMode, MaximumVirtualMachineMigrations, MaximumStorageMigrations")
	if err != nil {
		return "", fmt.Errorf("failed to get Hyper-V host information")
	}
	return string(out), nil
}
