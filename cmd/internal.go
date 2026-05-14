package cmd

import (
	"fmt"
	"os"

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

// loadConfig loads a Summarize config from the given path.
// Returns an empty Summarize if path is empty.
// On read/parse error, prints a warning to stderr and returns an empty Summarize
// so the caller can still fall back to treating the argument as a literal value.
func loadConfig(path string) hyperv.Summarize {
	if path == "" {
		return hyperv.Summarize{}
	}
	config, err := hyperv.UnmarshalYaml(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not load config %s: %s\n", path, err)
		return hyperv.Summarize{}
	}
	return config
}

// resolveDisk resolves a disk alias to a path using the loaded config.
// If the arg matches a key in config.Disks, returns its path; otherwise returns arg as-is.
func resolveDisk(config hyperv.Summarize, arg string) string {
	if d, ok := config.Disks[arg]; ok {
		return d.Path
	}
	return arg
}
