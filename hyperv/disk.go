// hyperv package is manage Hyper-V
package hyperv

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

// ------------------ //
// Create/Remove Disk
// ------------------ //

func CreateDisk(newDisk Disk, output bool) error {
	if newDisk.Import {
		return nil
	}

	err := checkDiskParam(newDisk)
	if output {
		PrintError("Check Disk Param", err)
	}
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
	if output {
		PrintError("Create Disk", err)
	}
	if err != nil {
		return err
	}
	return nil
}

func RemoveDisk(path string, output bool) error {
	err := isFileExist(path)
	if err != nil {
		return err
	}

	err = exec.Command("powershell", "-NoProfile", "rm '"+path+"'").Run()
	if output {
		PrintError("Remove Disk", err)
	}
	if err != nil {
		return err
	}
	return nil
}

// ----------------- //
// Check Disk Option
// ----------------- //

func checkDiskParam(newDisk Disk) error {
	if newDisk.Import {
		err := isFileExist(newDisk.Path)
		if err != nil {
			return err
		}
		return nil
	}
	err := isNotFileExist(newDisk.Path)
	if err != nil {
		return err
	}

	err = checkDiskTypeParam(newDisk.Type)
	if err != nil {
		return err
	}

	switch newDisk.Type {
	case "differencing":
		err := isFileExist(newDisk.ParentPath)
		if err != nil {
			return err
		}
	case "fixed":
	}

	err = checkDiskSizeParam(newDisk.Size)
	if err != nil {
		return err
	}

	return nil
}

func checkDiskTypeParam(diskType string) error {
	switch diskType {
	case "dynamic":
	case "fixed":
	case "differencing":
	default:
		return fmt.Errorf("error: Undefined disk type")
	}

	return nil
}

func checkDiskSizeParam(diskSize string) error {
	diskSize = regexp.MustCompile("^[0-9]*[TGM]B$").FindString(diskSize)
	if diskSize == "" {
		return fmt.Errorf("error: unknown size")
	}

	return nil
}
