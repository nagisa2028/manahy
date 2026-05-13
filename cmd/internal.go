package cmd

import (
	"fmt"

	"github.com/DevelopNaoki/manahy/hyperv"
)

func displayList(list []string, message string) {
	fmt.Println(message)
	for _, item := range list {
		fmt.Println("-", item)
	}
	fmt.Println()
}

func displayStorageList(storageList hyperv.StorageList) {
	fmt.Print("Storage\n")
	for i := range storageList.Number {
		fmt.Printf("- %s: %s: %.2f %s\n", storageList.Number[i], storageList.FriendlyName[i], storageList.Size[i], storageList.SizeUnit[i])
	}
	fmt.Print("\n")
	fmt.Print("More information, execute 'Get-Disk'\n")
}
