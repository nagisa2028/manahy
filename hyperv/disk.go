package hyperv

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

// CreateDisk creates a new virtual hard disk.
func CreateDisk(newDisk Disk, output bool) error {
	if newDisk.Import {
		return nil
	}

	err := checkDiskParam(newDisk)
	printError("Check Disk Param", err, output)
	if err != nil {
		return err
	}

	cmd := "New-VHD -Path " + newDisk.Path
	switch newDisk.Type {
	case "dynamic":
		cmd += " -SizeBytes " + newDisk.Size
	case "fixed":
		cmd += " -SizeBytes " + newDisk.Size
		cmd += " -SourceDisk " + strconv.Itoa(newDisk.SourceDisk)
		cmd += " -Fixed"
	case "differencing":
		cmd += " -ParentPath " + newDisk.ParentPath
		cmd += " -Differencing"
	}

	err = exec.Command("powershell", "-NoProfile", cmd).Run()
	printError("Create Disk", err, output)
	return err
}

// RemoveDisk deletes a virtual hard disk file.
func RemoveDisk(path string, output bool) error {
	if err := isFileExist(path); err != nil {
		return err
	}

	err := exec.Command("powershell", "-NoProfile", "rm '"+path+"'").Run()
	printError("Remove Disk", err, output)
	return err
}

func checkDiskParam(newDisk Disk) error {
	if newDisk.Import {
		return isFileExist(newDisk.Path)
	}
	if err := isNotFileExist(newDisk.Path); err != nil {
		return err
	}
	if err := checkDiskTypeParam(newDisk.Type); err != nil {
		return err
	}
	if newDisk.Type == "differencing" {
		if err := isFileExist(newDisk.ParentPath); err != nil {
			return err
		}
	}
	return checkDiskSizeParam(newDisk.Size)
}

func checkDiskTypeParam(diskType string) error {
	switch diskType {
	case "dynamic", "fixed", "differencing":
		return nil
	default:
		return fmt.Errorf("undefined disk type: %s", diskType)
	}
}

func checkDiskSizeParam(diskSize string) error {
	if regexp.MustCompile("^[0-9]*[TGM]B$").FindString(diskSize) == "" {
		return fmt.Errorf("invalid disk size format: %s (expected e.g. 10GB)", diskSize)
	}
	return nil
}

// GetVHDInfo returns detailed information about a VHD file.
func GetVHDInfo(path string) (string, error) {
	if err := isFileExist(path); err != nil {
		return "", err
	}
	res, err := exec.Command("powershell", "-NoProfile", "Get-VHD -Path '"+path+"' | Format-List Path, VhdType, FileSize, Size, ParentPath, Attached, DiskNumber").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get VHD info for %s", path)
	}
	return string(res), nil
}

// ResizeVHD resizes a VHD to the specified size.
func ResizeVHD(path, size string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	if err := checkDiskSizeParam(size); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile", "Resize-VHD -Path '"+path+"' -SizeBytes "+size).Run()
}

// OptimizeVHD compacts a VHD using full mode optimization.
func OptimizeVHD(path string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile", "Optimize-VHD -Path '"+path+"' -Mode Full").Run()
}

// ConvertVHD converts a VHD to a new format at the destination path.
func ConvertVHD(path, destPath, diskType string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	cmd := "Convert-VHD -Path '" + path + "' -DestinationPath '" + destPath + "'"
	if diskType != "" {
		cmd += " -VHDType " + diskType
	}
	return exec.Command("powershell", "-NoProfile", cmd).Run()
}

// MountVHD mounts a VHD.
func MountVHD(path string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile", "Mount-VHD -Path '"+path+"'").Run()
}

// DismountVHD dismounts a VHD.
func DismountVHD(path string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile", "Dismount-VHD -Path '"+path+"'").Run()
}

// MergeVHD merges a differencing VHD into its parent or a specified destination.
func MergeVHD(path, destPath string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	cmd := "Merge-VHD -Path '" + path + "'"
	if destPath != "" {
		cmd += " -DestinationPath '" + destPath + "'"
	}
	return exec.Command("powershell", "-NoProfile", cmd).Run()
}
