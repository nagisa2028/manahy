package hyperv

import (
	"fmt"
	"regexp"
	"strconv"
)

// maxDiskSizeGB is the upper bound for disk and memory size validation (64 TB).
const maxDiskSizeGB = 64 * 1024

//nolint:gochecknoglobals
var reDiskSize = regexp.MustCompile(`^([0-9]+)[TGM]B$`)

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

	cmd := cmdNewVHD + " -Path " + ps(newDisk.Path)
	switch newDisk.Type {
	case diskTypeDynamic:
		cmd += " -SizeBytes " + newDisk.Size
	case diskTypeFixed:
		cmd += " -SizeBytes " + newDisk.Size
		if newDisk.SourceDisk > 0 {
			cmd += " -SourceDisk " + strconv.Itoa(newDisk.SourceDisk)
		}
		cmd += " -Fixed"
	case diskTypeDifferencing:
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

	err := runPS(cmdRemoveItem + " " + ps(path))
	printError("Remove Disk", err, output)
	return err
}

func checkDiskParam(newDisk Disk) error {
	if !isWindowsAbsPath(newDisk.Path) {
		return fmt.Errorf("disk path must be absolute: %s", newDisk.Path)
	}
	if newDisk.Import {
		return isFileExist(newDisk.Path)
	}
	if err := isNotFileExist(newDisk.Path); err != nil {
		return err
	}
	if err := checkDiskType(newDisk.Type); err != nil {
		return err
	}
	if newDisk.Type == diskTypeDifferencing {
		if !isWindowsAbsPath(newDisk.ParentPath) {
			return fmt.Errorf("parent path must be absolute: %s", newDisk.ParentPath)
		}
		if err := isFileExist(newDisk.ParentPath); err != nil {
			return err
		}
		return nil // differencing disks inherit size from parent
	}
	if newDisk.Type == diskTypeFixed && newDisk.SourceDisk < 0 {
		return fmt.Errorf("source disk number must be non-negative, got %d", newDisk.SourceDisk)
	}
	return checkDiskSize(newDisk.Size)
}

func checkDiskType(diskType string) error {
	switch diskType {
	case diskTypeDynamic, diskTypeFixed, diskTypeDifferencing:
		return nil
	default:
		return fmt.Errorf("invalid disk type: %s", diskType)
	}
}

func checkDiskSize(diskSize string) error {
	m := reDiskSize.FindStringSubmatch(diskSize)
	if m == nil {
		return fmt.Errorf("invalid disk size format: %s (expected e.g. 10GB)", diskSize)
	}
	n, _ := strconv.Atoi(m[1])
	unit := diskSize[len(m[1]):]
	sizeGB := n
	switch unit {
	case "TB":
		sizeGB = n * 1024
	case "MB":
		sizeGB = 0 // MB values are always well under the cap
	}
	if sizeGB > maxDiskSizeGB {
		return fmt.Errorf("disk size %s exceeds maximum allowed size of %dTB", diskSize, maxDiskSizeGB/1024)
	}
	return nil
}

// GetVHDInfo returns detailed information about a VHD file.
func GetVHDInfo(path string) (string, error) {
	if err := isFileExist(path); err != nil {
		return "", err
	}
	res, err := outputPS(cmdGetVHD + " -Path " + ps(path) + " | Format-List Path, VhdType, FileSize, Size, ParentPath, Attached, DiskNumber")
	if err != nil {
		return "", fmt.Errorf("failed to get VHD info for %s: %w", path, err)
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
	return runPSLong(cmdResizeVHD + " -Path " + ps(path) + " -SizeBytes " + size)
}

// OptimizeVHD compacts a VHD using full mode optimization.
func OptimizeVHD(path string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	return runPSLong(cmdOptimizeVHD + " -Path " + ps(path) + " -Mode Full")
}

// ConvertVHD converts a VHD to a new format at the destination path.
func ConvertVHD(path, destPath, diskType string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	cmd := cmdConvertVHD + " -Path " + ps(path) + " -DestinationPath " + ps(destPath)
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
	return runPS(cmdMountVHD + " -Path " + ps(path))
}

// DismountVHD dismounts a VHD.
func DismountVHD(path string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	return runPS(cmdDismountVHD + " -Path " + ps(path))
}

// MergeVHD merges a differencing VHD into its parent or a specified destination.
func MergeVHD(path, destPath string) error {
	if err := isFileExist(path); err != nil {
		return err
	}
	cmd := cmdMergeVHD + " -Path " + ps(path)
	if destPath != "" {
		cmd += " -DestinationPath " + ps(destPath)
	}
	return runPSLong(cmd)
}

// CloneDisk creates a new VHD by copying an existing one using New-VHD -SourcePath.
// The source must exist; the destination must not.
func CloneDisk(srcPath, destPath string) error {
	if err := isFileExist(srcPath); err != nil {
		return fmt.Errorf("source disk: %w", err)
	}
	if err := isNotFileExist(destPath); err != nil {
		return fmt.Errorf("destination disk: %w", err)
	}
	return runPSLong(cmdNewVHD + " -Path " + ps(destPath) + " -SourcePath " + ps(srcPath))
}
