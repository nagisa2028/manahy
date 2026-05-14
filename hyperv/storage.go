package hyperv

// GetStorageList returns a list of physical storage devices.
func GetStorageList() (StorageList, error) {
	res, err := outputPS(cmdGetDisk + " | Format-Table Number,FriendlyName,Size")
	if err != nil {
		return StorageList{}, err
	}
	return storageListingOfExecuteResults(res)
}
