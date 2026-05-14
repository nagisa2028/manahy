package hyperv

import "fmt"

// CreateCheckpoint creates a checkpoint (snapshot) for a VM.
func CreateCheckpoint(vmName, name string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	cmd := "Checkpoint-VM -Name " + ps(vmName)
	if name != "" {
		cmd += " -SnapshotName " + ps(name)
	}
	return runPS(cmd)
}

// GetCheckpoints returns a formatted table of checkpoints for a VM.
func GetCheckpoints(vmName string) (string, error) {
	if err := IsVMExist(vmName); err != nil {
		return "", err
	}
	out, err := outputPS("Get-VMCheckpoint -VMName " + ps(vmName) + " | Sort-Object CreationTime | Format-Table Name, CreationTime | Out-String")
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
	return runPS("Restore-VMCheckpoint -VMName " + ps(vmName) + " -Name " + ps(name) + " -Confirm:$false")
}

// RemoveCheckpoint deletes a checkpoint from a VM.
func RemoveCheckpoint(vmName, name string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return runPS("Remove-VMCheckpoint -VMName " + ps(vmName) + " -Name " + ps(name) + " -Confirm:$false")
}

// RenameCheckpoint renames an existing checkpoint.
func RenameCheckpoint(vmName, name, newName string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return runPS("Rename-VMCheckpoint -VMName " + ps(vmName) + " -Name " + ps(name) + " -NewName " + ps(newName))
}

// ExportCheckpoint exports a checkpoint to the specified path.
func ExportCheckpoint(vmName, name, path string) error {
	if err := IsVMExist(vmName); err != nil {
		return err
	}
	return runPS("Export-VMCheckpoint -VMName " + ps(vmName) + " -Name " + ps(name) + " -Path " + ps(path))
}
