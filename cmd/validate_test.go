package cmd

import (
	"testing"

	"github.com/DevelopNaoki/manahy/hyperv"
)

func TestValidateConfig(t *testing.T) {
	t.Run("empty config is valid", func(t *testing.T) {
		errs := validateConfig(hyperv.Summarize{})
		if len(errs) != 0 {
			t.Errorf("empty config: expected no errors, got %v", errs)
		}
	})

	t.Run("valid VM config produces no errors", func(t *testing.T) {
		data := hyperv.Summarize{
			Vms: map[string]hyperv.VM{
				"myvm": {
					Generation: 2,
					Path:       `C:\VMs`,
					Memory:     hyperv.Memory{Size: "2GB", Dynamic: true},
					CPU:        hyperv.CPU{Thread: 2},
				},
			},
		}
		errs := validateConfig(data)
		if len(errs) != 0 {
			t.Errorf("valid config: expected no errors, got %v", errs)
		}
	})

	t.Run("invalid generation returns error", func(t *testing.T) {
		data := hyperv.Summarize{
			Vms: map[string]hyperv.VM{
				"badvm": {
					Generation: 3,
					Path:       `C:\VMs`,
					Memory:     hyperv.Memory{Size: "1GB"},
					CPU:        hyperv.CPU{Thread: 1},
				},
			},
		}
		errs := validateConfig(data)
		if len(errs) == 0 {
			t.Error("invalid generation: expected at least one error")
		}
		found := false
		for _, e := range errs {
			if contains(e, "generation") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("invalid generation: expected error mentioning 'generation', got %v", errs)
		}
	})

	t.Run("missing memory size returns error", func(t *testing.T) {
		data := hyperv.Summarize{
			Vms: map[string]hyperv.VM{
				"vm1": {Generation: 1, Path: `C:\VMs`, CPU: hyperv.CPU{Thread: 1}},
			},
		}
		errs := validateConfig(data)
		found := false
		for _, e := range errs {
			if contains(e, "memory.size") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing memory size: expected error mentioning 'memory.size', got %v", errs)
		}
	})

	t.Run("missing path returns error", func(t *testing.T) {
		data := hyperv.Summarize{
			Vms: map[string]hyperv.VM{
				"vm1": {Generation: 1, Memory: hyperv.Memory{Size: "1GB"}, CPU: hyperv.CPU{Thread: 1}},
			},
		}
		errs := validateConfig(data)
		found := false
		for _, e := range errs {
			if contains(e, "path is required") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing path: expected error mentioning 'path is required', got %v", errs)
		}
	})

	t.Run("memory buffer out of range returns error", func(t *testing.T) {
		data := hyperv.Summarize{
			Vms: map[string]hyperv.VM{
				"vm1": {
					Generation: 1,
					Path:       `C:\VMs`,
					Memory:     hyperv.Memory{Size: "1GB", Dynamic: true, Buffer: 150},
					CPU:        hyperv.CPU{Thread: 1},
				},
			},
		}
		errs := validateConfig(data)
		found := false
		for _, e := range errs {
			if contains(e, "buffer") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("buffer out of range: expected error mentioning 'buffer', got %v", errs)
		}
	})

	t.Run("network-refs with invalid VLAN returns error", func(t *testing.T) {
		data := hyperv.Summarize{
			Vms: map[string]hyperv.VM{
				"vm1": {
					Generation: 1,
					Path:       `C:\VMs`,
					Memory:     hyperv.Memory{Size: "1GB"},
					CPU:        hyperv.CPU{Thread: 1},
					NetworkRefs: []hyperv.NetworkRef{
						{Switch: "mySwitch", VLAN: 9999},
					},
				},
			},
		}
		errs := validateConfig(data)
		found := false
		for _, e := range errs {
			if contains(e, "vlan") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("invalid VLAN: expected error mentioning 'vlan', got %v", errs)
		}
	})

	t.Run("network-refs with empty switch name returns error", func(t *testing.T) {
		data := hyperv.Summarize{
			Vms: map[string]hyperv.VM{
				"vm1": {
					Generation:  1,
					Path:        `C:\VMs`,
					Memory:      hyperv.Memory{Size: "1GB"},
					CPU:         hyperv.CPU{Thread: 1},
					NetworkRefs: []hyperv.NetworkRef{{Switch: ""}},
				},
			},
		}
		errs := validateConfig(data)
		found := false
		for _, e := range errs {
			if contains(e, "empty switch") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("empty switch: expected error mentioning 'empty switch', got %v", errs)
		}
	})

	t.Run("disk missing path returns error", func(t *testing.T) {
		data := hyperv.Summarize{
			Disks: map[string]hyperv.Disk{
				"mydisk": {Size: "10GB", Type: "dynamic"},
			},
		}
		errs := validateConfig(data)
		found := false
		for _, e := range errs {
			if contains(e, "path is required") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing disk path: expected error mentioning 'path is required', got %v", errs)
		}
	})

	t.Run("disk missing type (non-import) returns error", func(t *testing.T) {
		data := hyperv.Summarize{
			Disks: map[string]hyperv.Disk{
				"mydisk": {Path: `C:\disk.vhdx`},
			},
		}
		errs := validateConfig(data)
		found := false
		for _, e := range errs {
			if contains(e, "type is required") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing disk type: expected error mentioning 'type is required', got %v", errs)
		}
	})

	t.Run("network missing type returns error", func(t *testing.T) {
		data := hyperv.Summarize{
			Networks: map[string]hyperv.VMSwitch{
				"mynet": {},
			},
		}
		errs := validateConfig(data)
		found := false
		for _, e := range errs {
			if contains(e, "type is required") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing network type: expected error mentioning 'type is required', got %v", errs)
		}
	})
}

// contains is a helper for test assertions.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
