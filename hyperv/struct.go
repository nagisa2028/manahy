// Package hyperv manages Hyper-V via PowerShell commands.
package hyperv

// VM defines the configuration for creating a virtual machine.
type VM struct {
	Name       string   `yaml:"-"`
	Count      int      `yaml:"count,omitempty"`
	Generation int      `yaml:"generation"`
	CPU        CPU      `yaml:"cpu"`
	Memory     Memory   `yaml:"memory"`
	Path       string   `yaml:"path"`
	Image      string   `yaml:"image,omitempty"`
	Disks      []string `yaml:"disks"`
	Networks   []string `yaml:"networks"`
}

// CPU defines the processor configuration for a virtual machine.
type CPU struct {
	Thread int  `yaml:"thread"`
	Nested bool `yaml:"nested"`
}

// Memory defines the memory configuration for a virtual machine.
type Memory struct {
	Size    string `yaml:"size"`
	Dynamic bool   `yaml:"dynamic"`
}

// Disk defines the configuration for creating a virtual hard disk.
type Disk struct {
	Path       string `yaml:"path"`
	Size       string `yaml:"size,omitempty"`
	Type       string `yaml:"type,omitempty"`
	ParentPath string `yaml:"parent-path,omitempty"`
	SourceDisk int    `yaml:"source-disk,omitempty"`
	Import     bool   `yaml:"import,omitempty"`
}

// VMSwitch defines the configuration for creating a virtual switch.
type VMSwitch struct {
	Name              string `yaml:"-"`
	Type              string `yaml:"type"`
	ExternalInterface string `yaml:"external-interface,omitempty"`
	AllowManagementOS bool   `yaml:"allow-management-os,omitempty"`
}

// SwitchList groups virtual switches by type.
type SwitchList struct {
	External []string
	Internal []string
	Private  []string
}

// VMList groups virtual machines by state.
type VMList struct {
	Running []string
	Saved   []string
	Paused  []string
	Off     []string
}

// StorageList is physical storage list.
type StorageList struct {
	Number       []string
	FriendlyName []string
	Size         []float64
	SizeUnit     []string
}

// Summarize is top-level structure for manahy.yaml.
type Summarize struct {
	Vms      map[string]VM       `yaml:"vms"`
	Disks    map[string]Disk     `yaml:"disks"`
	Networks map[string]VMSwitch `yaml:"networks"`
}
