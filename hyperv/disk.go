package hyperv

import (
	"fmt"
	"regexp"
	"strconv"
)

//nolint:gochecknoglobals
var reDiskSize = regexp.MustCompile("^[0-9]*[TGM]B$")

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

	cmd := "New-VHD -Path " + ps(newDisk.Path)
	switch newDisk.Type {
	case "dynamic":
		cmd += " -SizeBytes " + newDisk.Size
	case "fixed":
		cmd += " -SizeBytes " + newDisk.Size
		cmd += " -SourceDisk " + strconv.Itoa(newDisk.SourceDisk)
		cmd += " -Fixed"
	case "differencing":
		cmd += " -ParentPath " + ps(newDisk.ParentPath)
		cmd += " -Differencing"
	}

	err = runPS(cmd)
	printError("Create Disk", err, output)
	return err
}

// RemoveDisk deletes a virtual hard disk file.
func RemoveDisk(path string, output bool) error {
	if err := isFileExist(path); err != nil {
		return err
	}

	err := runPS("Remove-Item " + ps(path))
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
	if err := checkDiskType(newDisk.Type); err != nil {
		return err
	}
	if newDisk.Type == "differencing" {
		if err := isFileExist(newDisk.ParentPath); err != nil {
			return err
		}
	}
	return checkDiskSize(newDisk.Size)
}

func checkDiskType(diskType string) error {
	switch diskType {
	case "dynamic", "fixed", "differencing":
		return nil
	default:
		return fmt.Errorf("undefined disk type: %s", diskType)
	}
}

func checkDiskSize(diskSize string) error {
	if reDiskSize.FindString(diskSize) == "" {
		return fmt.Errorf("invalid disk size format: %s (expected e.g. 10GB)", diskSize)
	}
	return nil
}

// GetVHDInfo returns detailed information about a VHD file.
func GetVHDInfo(path string) (string, error) {
	if err := isFileExist(path); err != nil {
		return "", err
	}
	res, err := outputPS("Get-VHD -Path " + ps(path) + " | Format-List Path, VhdType, FileSize, Size, ParentPath, Attached, DiskNumber")
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
	if err := checkDiskSize(size); err != nil {
		return err
	}
	return runPSLong("Resize-VHD -Path " + ps(path) + " -SizeBytes " + size)
}

// OptimizeVHD compacts a VHD using full mode optimization.
func OptimizeVHD(path string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	return runPSLong("Optimize-VHD -Path " + ps(path) + " -Mode Full")
}

// ConvertVHD converts a VHD to a new format at the destination path.
func ConvertVHD(path, destPath, diskType string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	cmd := "Convert-VHD -Path " + ps(path) + " -DestinationPath " + ps(destPath)
	if diskType != "" {
		cmd += " -VHDType " + ps(diskType)
	}
	return runPSLong(cmd)
}

// MountVHD mounts a VHD.
func MountVHD(path string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	return runPS("Mount-VHD -Path " + ps(path))
}

// DismountVHD dismounts a VHD.
func DismountVHD(path string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	return runPS("Dismount-VHD -Path " + ps(path))
}

// MergeVHD merges a differencing VHD into its parent or a specified destination.
func MergeVHD(path, destPath string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	cmd := "Merge-VHD -Path " + ps(path)
	if destPath != "" {
		cmd += " -DestinationPath " + ps(destPath)
	}
	return runPSLong(cmd)
}
