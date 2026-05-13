package hyperv

import (
	"fmt"
	"os/exec"
)

// CreateCheckpoint creates a checkpoint (snapshot) for a VM.
func CreateCheckpoint(vmName, name string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	cmd := "Checkpoint-VM -Name '" + vmName + "'"
	if name != "" {
		cmd += " -SnapshotName '" + name + "'"
	}
	return exec.Command("powershell", "-NoProfile", cmd).Run()
}

// GetCheckpoints returns a formatted table of checkpoints for a VM.
func GetCheckpoints(vmName string) (string, error) {
	if err := IsVMExist(vmName); err != nil {
		return "", err
	}
	out, err := exec.Command("powershell", "-NoProfile",
		"Get-VMCheckpoint -VMName '"+vmName+"' | Sort-Object CreationTime | Format-Table Name, CreationTime | Out-String").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get checkpoints for VM %s", vmName)
	}
	return string(out), nil
}

// RestoreCheckpoint restores a VM to the named checkpoint.
func RestoreCheckpoint(vmName, name string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile",
		"Restore-VMCheckpoint -VMName '"+vmName+"' -Name '"+name+"' -Confirm:$false").Run()
}

// RemoveCheckpoint deletes a checkpoint from a VM.
func RemoveCheckpoint(vmName, name string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile",
		"Remove-VMCheckpoint -VMName '"+vmName+"' -Name '"+name+"' -Confirm:$false").Run()
}

// RenameCheckpoint renames an existing checkpoint.
func RenameCheckpoint(vmName, name, newName string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile",
		"Rename-VMCheckpoint -VMName '"+vmName+"' -Name '"+name+"' -NewName '"+newName+"'").Run()
}

// ExportCheckpoint exports a checkpoint to the specified path.
func ExportCheckpoint(vmName, name, path string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile",
		"Export-VMCheckpoint -VMName '"+vmName+"' -Name '"+name+"' -Path '"+path+"'").Run()
}
