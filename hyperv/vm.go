package hyperv

import (
	"fmt"
	"os/exec"
	"strconv"
)

// GetVMList returns a list of all VMs grouped by state.
func GetVMList() (VMList, error) {
	res, err := exec.Command("powershell", "-NoProfile", "Get-VM | Sort-Object State | Format-Table Name, State").Output()
	if err != nil {
		return VMList{}, err
	}
	return vmListingOfExecuteResults(res)
}

// GetVMState returns the current state of a VM, or "NotFound" / "Unknown".
func GetVMState(name string) string {
	res, err := exec.Command("powershell", "-NoProfile", "Get-VM '"+name+"' | Format-Table State").Output()
	if err != nil {
		return "NotFound"
	}
	vmState := listingOfExecuteResults(res, "State")
	if len(vmState) == 1 {
		return vmState[0]
	}
	return "Unknown"
}

// IsVMExist returns an error if the VM does not exist.
func IsVMExist(name string) error {
	switch GetVMState(name) {
	case "Unknown":
		return fmt.Errorf("failed to get state of VM %s", name)
	case "NotFound":
		return fmt.Errorf("VM %s does not exist", name)
	}
	return nil
}

// IsNotVMExist returns an error if the VM already exists.
func IsNotVMExist(name string) error {
	switch GetVMState(name) {
	case "Unknown":
		return fmt.Errorf("failed to get state of VM %s", name)
	case "NotFound":
		return nil
	}
	return fmt.Errorf("VM %s already exists", name)
}

// SetVMProcessor configures the CPU settings of a VM.
func SetVMProcessor(name string, cpu CPU) error {
	if err := IsVMExist(name); err != nil {
		return err
	}
	if err := checkVMProcessorParam(cpu); err != nil {
		return err
	}

	cmd := "Set-VMProcessor '" + name + "'"
	cmd += " -Count " + strconv.Itoa(cpu.Thread)
	cmd += " -ExposeVirtualizationExtensions $" + strconv.FormatBool(cpu.Nested)

	return exec.Command("powershell", "-NoProfile", cmd).Run()
}

// SetVMMemory configures the memory settings of a VM.
func SetVMMemory(name string, memory Memory) error {
	if err := IsVMExist(name); err != nil {
		return err
	}

	cmd := "Set-VMMemory -VMName '" + name + "'"
	cmd += " -StartupBytes " + memory.Size
	cmd += " -DynamicMemoryEnabled $" + strconv.FormatBool(memory.Dynamic)

	return exec.Command("powershell", "-NoProfile", cmd).Run()
}

// SetVMHardDisk attaches hard disk drives to a VM.
func SetVMHardDisk(name string, disks []string) error {
	if err := IsVMExist(name); err != nil {
		return err
	}

	for _, disk := range disks {
		if err := isFileExist(disk); err != nil {
			return err
		}

		cmd := "Add-VMHardDiskDrive -VMName '" + name + "'"
		cmd += " -Path " + disk

		if err := exec.Command("powershell", "-NoProfile", cmd).Run(); err != nil {
			return err
		}
	}
	return nil
}

// SetVMImageFile attaches a DVD/ISO image to a VM.
func SetVMImageFile(name string, image string) error {
	if err := IsVMExist(name); err != nil {
		return err
	}
	if err := isFileExist(image); err != nil {
		return err
	}

	cmd := "Add-VMDvdDrive -VMName '" + name + "'"
	cmd += " -Path " + image

	return exec.Command("powershell", "-NoProfile", cmd).Run()
}

// SetVMSwitch connects network adapters of a VM to virtual switches.
func SetVMSwitch(name string, networks []string) error {
	for _, network := range networks {
		switch GetSwitchType(network) {
		case "NotFound":
			return fmt.Errorf("switch %s does not exist", network)
		case "Unknown":
			return fmt.Errorf("failed to get state of switch %s", network)
		}

		cmd := "Add-VMNetworkAdapter -VMName " + name
		cmd += " -SwitchName " + network

		if err := exec.Command("powershell", "-NoProfile", cmd).Run(); err != nil {
			return err
		}
	}
	return nil
}

// CreateVM creates a new VM with the specified configuration.
func CreateVM(newVM VM, output bool) error {
	err := checkVMParam(newVM)
	printError("Check VM Param", err, output)
	if err != nil {
		return err
	}

	cmd := "New-VM -Name " + newVM.Name
	cmd += " -Generation " + strconv.Itoa(newVM.Generation)
	cmd += " -Path " + newVM.Path

	err = exec.Command("powershell", "-NoProfile", cmd).Run()
	printError("Create VM", err, output)
	if err != nil {
		return fmt.Errorf("failed to create VM %s", newVM.Name)
	}

	err = SetVMProcessor(newVM.Name, newVM.CPU)
	printError("Set Processor", err, output)
	if err != nil {
		return err
	}

	err = SetVMMemory(newVM.Name, newVM.Memory)
	printError("Set Memory", err, output)
	if err != nil {
		return err
	}

	err = SetVMHardDisk(newVM.Name, newVM.Disks)
	printError("Set HardDisk", err, output)
	if err != nil {
		return err
	}

	err = SetVMImageFile(newVM.Name, newVM.Image)
	printError("Set Image File", err, output)
	if err != nil {
		return err
	}

	err = SetVMSwitch(newVM.Name, newVM.Networks)
	printError("Set VMSwitch", err, output)
	return err
}

// RemoveVM deletes a VM forcefully.
func RemoveVM(name string, output bool) error {
	if err := IsVMExist(name); err != nil {
		return err
	}

	err := exec.Command("powershell", "-NoProfile", "Remove-VM -Name '"+name+"' -Force").Run()
	printError("Remove VM", err, output)
	return err
}

// RenameVM renames a VM.
func RenameVM(name string, newName string) error {
	if err := IsVMExist(name); err != nil {
		return err
	}
	if err := IsNotVMExist(newName); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile", "Rename-VM -Name '"+name+"' -NewName '"+newName+"'").Run()
}

// ConnectVM opens a VM console connection.
func ConnectVM(name string) error {
	if GetVMState(name) != "Running" {
		return fmt.Errorf("VM %s is not running", name)
	}
	return exec.Command("powershell", "-NoProfile", "vmconnect localhost '"+name+"'").Run()
}

// StartVM starts a VM.
func StartVM(name string) error {
	if GetVMState(name) == "Running" {
		return fmt.Errorf("VM %s is already running", name)
	}
	return exec.Command("powershell", "-NoProfile", "Start-VM '"+name+"'").Run()
}

// StopVM shuts down a VM gracefully.
func StopVM(name string) error {
	if GetVMState(name) != "Running" {
		return fmt.Errorf("VM %s is not running", name)
	}
	return exec.Command("powershell", "-NoProfile", "Stop-VM -Name '"+name+"'").Run()
}

// DestroyVM force-stops a VM.
func DestroyVM(name string) error {
	if GetVMState(name) != "Running" {
		return fmt.Errorf("VM %s is not running", name)
	}
	return exec.Command("powershell", "-NoProfile", "Stop-VM -Force -Name '"+name+"'").Run()
}

// SaveVM saves the state of a VM.
func SaveVM(name string) error {
	if GetVMState(name) != "Running" {
		return fmt.Errorf("VM %s is not running", name)
	}
	return exec.Command("powershell", "-NoProfile", "Save-VM -Name '"+name+"'").Run()
}

// SuspendVM pauses a VM.
func SuspendVM(name string) error {
	if GetVMState(name) != "Running" {
		return fmt.Errorf("VM %s is not running", name)
	}
	return exec.Command("powershell", "-NoProfile", "Suspend-VM -Name '"+name+"'").Run()
}

// RestartVM restarts a VM.
func RestartVM(name string) error {
	if GetVMState(name) != "Running" {
		return fmt.Errorf("VM %s is not running", name)
	}
	return exec.Command("powershell", "-NoProfile", "Restart-VM -Name '"+name+"' -Force").Run()
}

func checkVMParam(newVM VM) error {
	if err := IsNotVMExist(newVM.Name); err != nil {
		return err
	}
	if err := checkVMGeneration(newVM.Generation); err != nil {
		return err
	}
	if err := checkVMPath(newVM.Name, newVM.Path); err != nil {
		return err
	}
	if newVM.Image != "" {
		return isFileExist(newVM.Image)
	}
	return nil
}

func checkVMGeneration(generation int) error {
	if generation < 1 || generation > 2 {
		return fmt.Errorf("generation must be 1 or 2")
	}
	return nil
}

func checkVMPath(name string, path string) error {
	return isNotFileExist(path + "\\" + name)
}

func checkVMProcessorParam(cpu CPU) error {
	if cpu.Thread < 1 {
		return fmt.Errorf("vcpu count must be at least 1")
	}
	return nil
}
