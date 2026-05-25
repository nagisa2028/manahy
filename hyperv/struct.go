// Package hyperv manages Hyper-V via PowerShell commands.
package hyperv

// VM defines the configuration for creating a virtual machine.
type VM struct {
	Name               string       `yaml:"-"`
	Count              int          `yaml:"count,omitempty"`
	Generation         int          `yaml:"generation"`
	CPU                CPU          `yaml:"cpu"`
	Memory             Memory       `yaml:"memory"`
	Path               string       `yaml:"path"`
	Image              string       `yaml:"image,omitempty"`
	Notes              string       `yaml:"notes,omitempty"`
	Disks              []string     `yaml:"disks"`
	Networks           []string     `yaml:"networks"`
	NetworkRefs        []NetworkRef `yaml:"network-refs,omitempty"`
	SecureBoot         *bool        `yaml:"secure-boot,omitempty"`
	SecureBootTemplate string       `yaml:"secure-boot-template,omitempty"`
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
	// Min, Max, and Buffer are only applied when Dynamic is true.
	// Min and Max follow the same format as Size (e.g. "512MB", "1GB").
	// Buffer is a percentage of demand (5–100); 0 means "use the Hyper-V default".
	Min    string `yaml:"min,omitempty"`
	Max    string `yaml:"max,omitempty"`
	Buffer int    `yaml:"buffer,omitempty"`
}

// NetworkRef describes a virtual switch attachment with an optional VLAN ID
// and adapter name.  It is used in the YAML network-refs field, which
// supersedes the simpler networks: [switchName] shorthand when VLAN or a
// named adapter is needed.
type NetworkRef struct {
	Switch string `yaml:"switch"`
	// VLAN is the IEEE 802.1Q access VLAN ID (1–4094).  0 means untagged.
	VLAN int `yaml:"vlan,omitempty"`
	// Name is the adapter name visible inside Windows.  When empty, Hyper-V
	// assigns a default name ("Network Adapter").
	Name string `yaml:"name,omitempty"`
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

// StorageList holds information about physical storage devices.
type StorageList struct {
	Number       []string
	FriendlyName []string
	Size         []float64
	SizeUnit     []string
}

// Summarize is the top-level structure for manahy.yaml.
type Summarize struct {
	Vms      map[string]VM       `yaml:"vms"`
	Disks    map[string]Disk     `yaml:"disks"`
	Networks map[string]VMSwitch `yaml:"networks"`
}
