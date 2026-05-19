package hyperv

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Hyper-V supports only generation 1 and 2.
const (
	minVMGeneration     = 1
	maxVMGeneration     = 2
	minVMProcessorCount = 1
	maxVMCount          = 256
	maxMemorySizeGB     = 16 * 1024 // 16 TB upper bound for memory size validation
)

//nolint:gochecknoglobals
var reMemorySize = regexp.MustCompile(`^[0-9]+[TGM]B$`)

// GetVMList returns a list of all VMs grouped by state.
func GetVMList() (VMList, error) {
	res, err := outputPS(cmdGetVM + " | Sort-Object State | Format-Table Name, State")
	if err != nil {
		return VMList{}, err
	}
	return vmListingOfExecuteResults(res)
}

// GetVMState returns the current state of a VM, or "NotFound" / "Unknown".
func GetVMState(name string) string {
	res, err := outputPS(cmdGetVM + " " + ps(name) + " | Format-Table State")
	if err != nil {
		return vmStateNotFound
	}
	vmState := listingOfExecuteResults(res, "State")
	if len(vmState) == 1 {
		return vmState[0]
	}
	return vmStateUnknown
}

// getVMStateMap returns name→state for all VMs on the host in one PS call.
// VMs whose state cannot be parsed are omitted; callers treat missing entries as NotFound.
func getVMStateMap() map[string]string {
	res, err := outputPS(cmdGetVM + " | Format-Table Name, State")
	if err != nil {
		return map[string]string{}
	}
	m := make(map[string]string)
	for _, line := range reSplit.Split(string(res), -1) {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "Name") || reDashOnly.MatchString(line) {
			continue
		}
		state := reVMState.FindString(line)
		if state == "" {
			continue
		}
		name := strings.TrimSpace(reVMState.ReplaceAllString(line, ""))
		if name != "" {
			m[name] = state
		}
	}
	return m
}

// IsVMExist returns an error if the VM does not exist.
func IsVMExist(name string) error {
	switch GetVMState(name) {
	case vmStateUnknown:
		return fmt.Errorf("failed to get state of VM %s", name)
	case vmStateNotFound:
		return &notFoundError{msg: fmt.Sprintf("VM %s does not exist", name)}
	}
	return nil
}

// IsNotVMExist returns an error if the VM already exists.
func IsNotVMExist(name string) error {
	switch GetVMState(name) {
	case vmStateUnknown:
		return fmt.Errorf("failed to get state of VM %s", name)
	case vmStateNotFound:
		return nil
	}
	return fmt.Errorf("VM %s already exists", name)
}

// SetVMProcessor configures the CPU settings of a VM.
func SetVMProcessor(name string, cpu CPU) error {
	if err := IsVMExist(name); err != nil {
		return err
	}
	if err := checkVMProcessor(cpu); err != nil {
		return err
	}

	cmd := cmdSetVMProcessor + " " + ps(name)
	cmd += " -Count " + strconv.Itoa(cpu.Thread)
	cmd += " -ExposeVirtualizationExtensions $" + strconv.FormatBool(cpu.Nested)

	return runPS(cmd)
}

// SetVMMemory configures the memory settings of a VM.
func SetVMMemory(name string, memory Memory) error {
	if err := IsVMExist(name); err != nil {
		return err
	}

	cmd := cmdSetVMMemory + " -VMName " + ps(name)
	cmd += " -StartupBytes " + memory.Size
	cmd += " -DynamicMemoryEnabled $" + strconv.FormatBool(memory.Dynamic)

	return runPS(cmd)
}

// SetVMHardDisk attaches hard disk drives to a VM.
func SetVMHardDisk(name string, disks []string) error {
	if err := IsVMExist(name); err != nil {
		return err
	}
	if len(disks) == 0 {
		return nil
	}
	for _, disk := range disks {
		if err := isFileExist(disk); err != nil {
			return err
		}
	}
	var sb strings.Builder
	for i, disk := range disks {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(cmdAddVMHardDiskDrive + " -VMName " + ps(name) + " -Path " + ps(disk))
	}
	return runPS(sb.String())
}

// SetVMImageFile attaches a DVD/ISO image to a VM.
func SetVMImageFile(name string, image string) error {
	if err := IsVMExist(name); err != nil {
		return err
	}
	if err := isFileExist(image); err != nil {
		return err
	}

	cmd := cmdAddVMDvdDrive + " -VMName " + ps(name)
	cmd += " -Path " + ps(image)

	return runPS(cmd)
}

// SetVMSwitch connects network adapters of a VM to virtual switches.
func SetVMSwitch(name string, networks []string) error {
	if len(networks) == 0 {
		return nil
	}
	for _, network := range networks {
		switch GetSwitchType(network) {
		case vmStateNotFound:
			return fmt.Errorf("switch %s does not exist", network)
		case vmStateUnknown:
			return fmt.Errorf("failed to get state of switch %s", network)
		}
	}
	var sb strings.Builder
	for i, network := range networks {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(cmdAddVMNetworkAdapter + " -VMName " + ps(name) + " -SwitchName " + ps(network))
	}
	return runPS(sb.String())
}

// CreateVM creates a new VM with the specified configuration.
func CreateVM(newVM VM, output bool) error {
	err := checkVMParam(newVM)
	printError("Check VM Param", err, output)
	if err != nil {
		return err
	}

	cmd := cmdNewVM + " -Name " + ps(newVM.Name)
	cmd += " -Generation " + strconv.Itoa(newVM.Generation)
	cmd += " -Path " + ps(newVM.Path)

	err = runPS(cmd)
	printError("Create VM", err, output)
	if err != nil {
		return fmt.Errorf("failed to create VM %s: %w", newVM.Name, err)
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

	if newVM.Image != "" {
		err = SetVMImageFile(newVM.Name, newVM.Image)
		printError("Set Image File", err, output)
		if err != nil {
			return err
		}
	}

	err = SetVMSwitch(newVM.Name, newVM.Networks)
	printError("Set VMSwitch", err, output)
	if err != nil {
		return err
	}

	if newVM.Generation == maxVMGeneration && newVM.SecureBoot != nil {
		err = SetVMSecureBoot(newVM.Name, *newVM.SecureBoot, newVM.SecureBootTemplate)
		printError("Set Secure Boot", err, output)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveVM deletes a VM forcefully.
func RemoveVM(name string, output bool) error {
	if err := IsVMExist(name); err != nil {
		return err
	}

	err := runPS(cmdRemoveVM + " -Name " + ps(name) + " -Force")
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
	return runPS(cmdRenameVM + " -Name " + ps(name) + " -NewName " + ps(newName))
}

// ConnectVM opens a VM console connection.
func ConnectVM(name string) error {
	if GetVMState(name) != vmStateRunning {
		return fmt.Errorf("VM %s is not running", name)
	}
	return runPS(cmdVMConnect + " localhost " + ps(name))
}

// StartVM starts a VM.
func StartVM(name string) error {
	if GetVMState(name) == vmStateRunning {
		return fmt.Errorf("VM %s is already running", name)
	}
	return runPS(cmdStartVM + " " + ps(name))
}

// StopVM shuts down a VM gracefully.
func StopVM(name string) error {
	if GetVMState(name) != vmStateRunning {
		return fmt.Errorf("VM %s is not running", name)
	}
	return runPS(cmdStopVM + " -Name " + ps(name))
}

// DestroyVM force-stops a VM.
func DestroyVM(name string) error {
	if GetVMState(name) != vmStateRunning {
		return fmt.Errorf("VM %s is not running", name)
	}
	return runPS(cmdStopVM + " -Force -Name " + ps(name))
}

// SaveVM saves the state of a VM.
func SaveVM(name string) error {
	if GetVMState(name) != vmStateRunning {
		return fmt.Errorf("VM %s is not running", name)
	}
	return runPS(cmdSaveVM + " -Name " + ps(name))
}

// SuspendVM pauses a VM.
func SuspendVM(name string) error {
	if GetVMState(name) != vmStateRunning {
		return fmt.Errorf("VM %s is not running", name)
	}
	return runPS(cmdSuspendVM + " -Name " + ps(name))
}

// RestartVM restarts a VM.
func RestartVM(name string) error {
	if GetVMState(name) != vmStateRunning {
		return fmt.Errorf("VM %s is not running", name)
	}
	return runPS(cmdRestartVM + " -Name " + ps(name) + " -Force")
}

// ResumeVM resumes a paused VM.
func ResumeVM(name string) error {
	if GetVMState(name) != vmStatePaused {
		return fmt.Errorf("VM %s is not paused", name)
	}
	return runPS(cmdResumeVM + " -Name " + ps(name))
}

// ExportVM exports a VM to the specified directory.
func ExportVM(name, path string) error {
	if err := IsVMExist(name); err != nil {
		return err
	}
	return runPSLong(cmdExportVM + " -Name " + ps(name) + " -Path " + ps(path))
}

// ImportVM registers a VM from a .vmcx file path.
func ImportVM(path string) error {
	return runPSLong(cmdImportVM + " -Path " + ps(path))
}

// MoveVMStorage moves all VM storage files to a new directory on the same host.
func MoveVMStorage(name, destPath string) error {
	if err := IsVMExist(name); err != nil {
		return err
	}
	return runPSLong(cmdMoveVMStorage + " -VMName " + ps(name) + " -DestinationStoragePath " + ps(destPath))
}

// CopyVM exports the source VM then imports it as a new VM with a different name.
// The exported files are placed under destPath and remain after the copy completes.
func CopyVM(name, newName, destPath string) error {
	if err := IsVMExist(name); err != nil {
		return err
	}
	if err := IsNotVMExist(newName); err != nil {
		return err
	}
	// Guard against backslashes in VM names breaking the intermediate path construction.
	if strings.ContainsAny(name, `\/`) {
		return fmt.Errorf("VM name %q must not contain path separators", name)
	}
	script := cmdExportVM + " -Name " + ps(name) + " -Path " + ps(destPath) + "; " +
		"$vmcx = (Get-ChildItem -Recurse -Path " + ps(destPath+"\\"+name) + " -Filter '*.vmcx' | Select-Object -First 1).FullName; " +
		"$newVM = " + cmdImportVM + " -Path $vmcx -Copy -GenerateNewId; " +
		cmdRenameVM + " -VM $newVM -NewName " + ps(newName)
	return runPSLong(script)
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
	if err := checkMemorySize(newVM.Memory.Size); err != nil {
		return err
	}
	if newVM.Image != "" {
		if !isWindowsAbsPath(newVM.Image) {
			return fmt.Errorf("image path must be absolute: %s", newVM.Image)
		}
		return isFileExist(newVM.Image)
	}
	return nil
}

func checkVMGeneration(generation int) error {
	if generation < minVMGeneration || generation > maxVMGeneration {
		return fmt.Errorf("generation must be %d or %d", minVMGeneration, maxVMGeneration)
	}
	return nil
}

func checkVMPath(name string, path string) error {
	if !isWindowsAbsPath(path) {
		return fmt.Errorf("VM path must be absolute: %s", path)
	}
	return isNotFileExist(path + "\\" + name)
}

func checkVMProcessor(cpu CPU) error {
	if cpu.Thread < minVMProcessorCount {
		return fmt.Errorf("vcpu count must be at least %d", minVMProcessorCount)
	}
	return nil
}

func checkMemorySize(size string) error {
	if reMemorySize.FindString(size) == "" {
		return fmt.Errorf("invalid memory size format: %s (expected e.g. 512MB, 1GB)", size)
	}
	n, _ := strconv.Atoi(size[:len(size)-2])
	unit := size[len(size)-2:]
	sizeGB := n
	if unit == "TB" {
		sizeGB = n * 1024
	}
	if sizeGB > maxMemorySizeGB {
		return fmt.Errorf("memory size %s exceeds maximum allowed size of %dTB", size, maxMemorySizeGB/1024)
	}
	return nil
}
