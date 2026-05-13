// Package hyperv manages Hyper-V via PowerShell commands.
package hyperv

// VM is create vm option.
type VM struct {
	Name       string   `yaml:"name" json:"name"`
	Count      int      `yaml:"count,omitempty"`
	Generation int      `yaml:"generation" json:"generation"`
	CPU        CPU      `yaml:"cpu" json:"cpu"`
	Memory     Memory   `yaml:"memory" json:"memory"`
	Path       string   `yaml:"path" json:"path"`
	Image      string   `yaml:"image,omitempty" json:"image"`
	Disks      []string `yaml:"disk"`
	Networks   []string `yaml:"network"`
}

// CPU is set-processor option.
type CPU struct {
	Thread int  `yaml:"thread"`
	Nested bool `yaml:"nested"`
}

// Memory is set memory option.
type Memory struct {
	Size    string `yaml:"size"`
	Dynamic bool   `yaml:"dynamic"`
}

// Disk is create disk option.
type Disk struct {
	Path       string `yaml:"path"`
	Size       string `yaml:"size,omitempty"`
	Type       string `yaml:"type,omitempty"`
	ParentPath string `yaml:"parent-path,omitempty"`
	SourceDisk int    `yaml:"source-disk,omitempty"`
	Import     bool   `yaml:"import,omitempty"`
}

// VMSwitch is create switch option.
type VMSwitch struct {
	Name              string `yaml:"name"`
	Type              string `yaml:"type"`
	ExternalInterface string `yaml:"external-interface,omitempty"`
	AllowManagementOs bool   `yaml:"allow-management-os,omitempty"`
}

// SwitchList is all type switch list.
type SwitchList struct {
	External []string
	Internal []string
	Private  []string
}

// VMList is all status vm list.
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
