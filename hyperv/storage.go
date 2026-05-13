// hyperv package is manage Hyper-V
package hyperv

import "os/exec"

// GetStorageList get a list of physical storage
func GetStorageList() (StorageList, error) {
	res, err := exec.Command("powershell", "-NoProfile", "Get-Disk | Format-Table Number,FriendlyName,Size").Output()
	if err != nil {
		return StorageList{}, err
	}
	return storageListingOfExecuteResults(res)
}
