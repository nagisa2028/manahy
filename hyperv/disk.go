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
