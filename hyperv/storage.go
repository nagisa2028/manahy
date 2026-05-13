package hyperv

import "os/exec"

// GetStorageList returns a list of physical storage devices.
func GetStorageList() (StorageList, error) {
	res, err := exec.Command("powershell", "-NoProfile", "Get-Disk | Format-Table Number,FriendlyName,Size").Output()
	if err != nil {
		return StorageList{}, err
	}
	return storageListingOfExecuteResults(res)
}
