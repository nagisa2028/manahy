package hyperv

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// dryrunOutputMock handles all outputPS calls used by DryRunBuild / DryRunRemove.
// Get-VMSwitch must be checked before Get-VM (substring order).
func dryrunOutputMock(existsFile, existsSwitch, existsVM bool) func(string) ([]byte, error) {
	return func(cmd string) ([]byte, error) {
		switch {
		case strings.Contains(cmd, "Test-Path"):
			if existsFile {
				return []byte("True\n"), nil
			}
			return []byte("False\n"), nil
		case strings.Contains(cmd, "Get-VMSwitch"):
			if existsSwitch {
				return []byte("SwitchType\n----------\nInternal\n"), nil
			}
			return nil, errors.New("not found")
		case strings.Contains(cmd, "Get-VM"):
			if existsVM {
				return stateOutput("Off"), nil
			}
			return nil, errors.New("not found")
		default:
			return nil, errors.New("unexpected call")
		}
	}
}

// --- DryRunBuild ---

func TestDryRunBuild(t *testing.T) {
	t.Run("new resources show create", func(t *testing.T) {
		withPS(t, nil, dryrunOutputMock(false, false, false))
		config := Summarize{
			Disks: map[string]Disk{
				"boot": {Path: `C:\boot.vhd`, Size: "10GB", Type: "dynamic"},
			},
			Networks: map[string]VMSwitch{
				"internal": {Type: "internal"},
			},
			Vms: map[string]VM{
				"test-vm": {Generation: 1, Path: `C:\VMs`, Memory: Memory{Size: "512MB"}, CPU: CPU{Thread: 1}},
			},
		}
		var buf bytes.Buffer
		DryRunBuild(config, &buf)
		out := buf.String()
		if !strings.Contains(out, "create") {
			t.Errorf("DryRunBuild: expected 'create' in output, got:\n%s", out)
		}
	})

	t.Run("existing resources show skip", func(t *testing.T) {
		withPS(t, nil, dryrunOutputMock(true, true, true))
		config := Summarize{
			Disks: map[string]Disk{
				"boot": {Path: `C:\boot.vhd`, Size: "10GB", Type: "dynamic"},
			},
			Networks: map[string]VMSwitch{
				"internal": {Type: "internal"},
			},
			Vms: map[string]VM{
				"test-vm": {Generation: 1, Path: `C:\VMs`, Memory: Memory{Size: "512MB"}, CPU: CPU{Thread: 1}},
			},
		}
		var buf bytes.Buffer
		DryRunBuild(config, &buf)
		out := buf.String()
		if strings.Contains(out, "create") {
			t.Errorf("DryRunBuild: expected no 'create' for existing resources, got:\n%s", out)
		}
		if !strings.Contains(out, "skip") {
			t.Errorf("DryRunBuild: expected 'skip' in output, got:\n%s", out)
		}
	})

	t.Run("import disk is skipped", func(t *testing.T) {
		withPS(t, nil, dryrunOutputMock(false, false, false))
		config := Summarize{
			Disks: map[string]Disk{
				"imported": {Path: `C:\existing.vhd`, Import: true},
			},
		}
		var buf bytes.Buffer
		DryRunBuild(config, &buf)
		out := buf.String()
		if !strings.Contains(out, "skip") || !strings.Contains(out, "import: true") {
			t.Errorf("DryRunBuild: expected import skip, got:\n%s", out)
		}
	})

	t.Run("invalid memory size shows error", func(t *testing.T) {
		withPS(t, nil, dryrunOutputMock(false, false, false))
		config := Summarize{
			Vms: map[string]VM{
				"bad-vm": {Generation: 1, Path: `C:\VMs`, Memory: Memory{Size: "1024"}, CPU: CPU{Thread: 1}},
			},
		}
		var buf bytes.Buffer
		DryRunBuild(config, &buf)
		if !strings.Contains(buf.String(), "error") {
			t.Errorf("DryRunBuild: expected 'error' for invalid memory size, got:\n%s", buf.String())
		}
	})

	t.Run("count > 1 generates suffixed names", func(t *testing.T) {
		withPS(t, nil, dryrunOutputMock(false, false, false))
		config := Summarize{
			Vms: map[string]VM{
				"web": {Generation: 1, Path: `C:\VMs`, Memory: Memory{Size: "512MB"}, CPU: CPU{Thread: 1}, Count: 3},
			},
		}
		var buf bytes.Buffer
		DryRunBuild(config, &buf)
		out := buf.String()
		for _, name := range []string{"web1", "web2", "web3"} {
			if !strings.Contains(out, name) {
				t.Errorf("DryRunBuild: expected %q in output, got:\n%s", name, out)
			}
		}
	})
}

// --- DryRunRemove ---

func TestDryRunRemove(t *testing.T) {
	t.Run("existing resources show remove", func(t *testing.T) {
		withPS(t, nil, dryrunOutputMock(true, true, true))
		config := Summarize{
			Vms:      map[string]VM{"my-vm": {}},
			Networks: map[string]VMSwitch{"internal": {Type: "internal"}},
			Disks:    map[string]Disk{"data": {Path: `C:\data.vhd`}},
		}
		var buf bytes.Buffer
		DryRunRemove(config, &buf)
		out := buf.String()
		if !strings.Contains(out, "remove") {
			t.Errorf("DryRunRemove: expected 'remove' in output, got:\n%s", out)
		}
	})

	t.Run("missing resources show skip", func(t *testing.T) {
		withPS(t, nil, dryrunOutputMock(false, false, false))
		config := Summarize{
			Vms:      map[string]VM{"ghost-vm": {}},
			Networks: map[string]VMSwitch{"missing-net": {Type: "private"}},
			Disks:    map[string]Disk{"gone": {Path: `C:\gone.vhd`}},
		}
		var buf bytes.Buffer
		DryRunRemove(config, &buf)
		out := buf.String()
		if strings.Contains(out, "  remove") {
			t.Errorf("DryRunRemove: expected no 'remove' action for missing resources, got:\n%s", out)
		}
		if !strings.Contains(out, "skip") {
			t.Errorf("DryRunRemove: expected 'skip' in output, got:\n%s", out)
		}
	})

	t.Run("import disk is skipped", func(t *testing.T) {
		withPS(t, nil, dryrunOutputMock(true, false, false))
		config := Summarize{
			Disks: map[string]Disk{
				"imported": {Path: `C:\existing.vhd`, Import: true},
			},
		}
		var buf bytes.Buffer
		DryRunRemove(config, &buf)
		out := buf.String()
		if !strings.Contains(out, "import: true") {
			t.Errorf("DryRunRemove: expected import skip, got:\n%s", out)
		}
	})

	t.Run("count > 1 generates suffixed names", func(t *testing.T) {
		withPS(t, nil, dryrunOutputMock(false, false, true))
		config := Summarize{
			Vms: map[string]VM{"web": {Count: 2}},
		}
		var buf bytes.Buffer
		DryRunRemove(config, &buf)
		out := buf.String()
		for _, name := range []string{"web1", "web2"} {
			if !strings.Contains(out, name) {
				t.Errorf("DryRunRemove: expected %q in output, got:\n%s", name, out)
			}
		}
	})
}
