package hyperv

import (
	"fmt"
	"os/exec"
)

// GetVMInfo returns processor and memory configuration for a VM.
func GetVMInfo(name string) (string, error) {
	if err := IsVMExist(name); err != nil {
		return "", err
	}
	procOut, err := exec.Command("powershell", "-NoProfile",
		"Get-VMProcessor -VMName '"+name+"' | Format-List VMName, Count, ExposeVirtualizationExtensions").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get processor info for VM %s", name)
	}
	memOut, err := exec.Command("powershell", "-NoProfile",
		"Get-VMMemory -VMName '"+name+"' | Format-List VMName, DynamicMemoryEnabled, Startup, Minimum, Maximum, Buffer, Priority").Output()
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
	script := "Enable-VMResourceMetering -VMName '" + name + "'; " +
		"Measure-VM -VMName '" + name + "' | Format-List AvgCPUUsage, AvgRAMUsage, TotalDisk, NetworkMeteredTrafficReport"
	out, err := exec.Command("powershell", "-NoProfile", script).Output()
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
	out, err := exec.Command("powershell", "-NoProfile",
		"Get-VMIntegrationService -VMName '"+name+"' | Format-Table Name, Enabled, PrimaryStatusDescription | Out-String").Output()
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
	return exec.Command("powershell", "-NoProfile",
		"Enable-VMIntegrationService -VMName '"+name+"' -Name '"+service+"'").Run()
}

// DisableVMIntegrationService disables a named integration service on a VM.
func DisableVMIntegrationService(name, service string) error {
	if err := IsVMExist(name); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile",
		"Disable-VMIntegrationService -VMName '"+name+"' -Name '"+service+"'").Run()
}

// GetVMHost returns Hyper-V host configuration.
func GetVMHost() (string, error) {
	out, err := exec.Command("powershell", "-NoProfile",
		"Get-VMHost | Format-List VirtualHardDiskPath, VirtualMachinePath, MacAddressMinimum, MacAddressMaximum, NumaSpanningEnabled, EnableEnhancedSessionMode, MaximumVirtualMachineMigrations, MaximumStorageMigrations").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get Hyper-V host information")
	}
	return string(out), nil
}
