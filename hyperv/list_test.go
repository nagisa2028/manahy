package hyperv

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestGetResourceList(t *testing.T) {
	const diskPath = `C:\disks\test.vhdx`
	const yamlContent = `
vms:
  test-vm:
    generation: 1
    cpu:
      thread: 1
    memory:
      size: 512MB
    path: C:\VMs
    disks: []
    networks: []
networks:
  test-switch:
    type: internal
disks:
  test-disk:
    path: C:\disks\test.vhdx
`
	f, err := os.CreateTemp(t.TempDir(), "manahy-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err = f.WriteString(yamlContent); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	if err = f.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	// makeHook returns a hook that answers True for all Test-Path calls except
	// for the disk path, which returns the value controlled by diskPresent.
	makeHook := func(vmState string, switchType string, diskPresent bool) func(string) ([]byte, error) {
		return func(cmd string) ([]byte, error) {
			// Check Get-VMSwitch before Get-VM because the latter is a substring of the former.
			if strings.Contains(cmd, "Get-VMSwitch") {
				if switchType == "" {
					// Returning an error causes GetSwitchType to return vmStateNotFound.
					return nil, errors.New("not found")
				}
				return []byte("SwitchType\n----------\n" + switchType + "\n"), nil
			}
			if strings.Contains(cmd, "Get-VM") {
				return []byte("State\n-----\n" + vmState + "\n"), nil
			}
			// Test-Path: return True for config yaml, controlled for disk path.
			if strings.Contains(cmd, diskPath) {
				if diskPresent {
					return []byte("True\n"), nil
				}
				return []byte("False\n"), nil
			}
			// Config yaml or any other path: file exists.
			return []byte("True\n"), nil
		}
	}

	t.Run("returns VM state from GetVMState", func(t *testing.T) {
		withPS(t, nil, makeHook("Running", "Internal", false))

		rl, err := GetResourceList(f.Name())
		if err != nil {
			t.Fatalf("GetResourceList: expected nil, got %v", err)
		}
		if len(rl.VMs) != 1 {
			t.Fatalf("GetResourceList: expected 1 VM, got %d", len(rl.VMs))
		}
		if rl.VMs[0].Name != "test-vm" {
			t.Errorf("GetResourceList: VM name %q, want 'test-vm'", rl.VMs[0].Name)
		}
		if rl.VMs[0].Status != "Running" {
			t.Errorf("GetResourceList: VM status %q, want 'Running'", rl.VMs[0].Status)
		}
	})

	t.Run("disk present when file exists", func(t *testing.T) {
		withPS(t, nil, makeHook("Off", "Internal", true))

		rl, err := GetResourceList(f.Name())
		if err != nil {
			t.Fatalf("GetResourceList: expected nil, got %v", err)
		}
		if len(rl.Disks) != 1 {
			t.Fatalf("GetResourceList: expected 1 disk, got %d", len(rl.Disks))
		}
		if !strings.Contains(rl.Disks[0].Status, "present") {
			t.Errorf("GetResourceList: disk status %q does not contain 'present'", rl.Disks[0].Status)
		}
	})

	t.Run("disk missing when file does not exist", func(t *testing.T) {
		withPS(t, nil, makeHook("Off", "", false))

		rl, err := GetResourceList(f.Name())
		if err != nil {
			t.Fatalf("GetResourceList: expected nil, got %v", err)
		}
		if len(rl.Disks) != 1 {
			t.Fatalf("GetResourceList: expected 1 disk, got %d", len(rl.Disks))
		}
		if !strings.Contains(rl.Disks[0].Status, "missing") {
			t.Errorf("GetResourceList: disk status %q does not contain 'missing'", rl.Disks[0].Status)
		}
	})

	t.Run("network switch type shown in status", func(t *testing.T) {
		withPS(t, nil, makeHook("Off", "Internal", false))

		rl, err := GetResourceList(f.Name())
		if err != nil {
			t.Fatalf("GetResourceList: expected nil, got %v", err)
		}
		if len(rl.Networks) != 1 {
			t.Fatalf("GetResourceList: expected 1 network, got %d", len(rl.Networks))
		}
		if rl.Networks[0].Name != "test-switch" {
			t.Errorf("GetResourceList: network name %q, want 'test-switch'", rl.Networks[0].Name)
		}
		if !strings.Contains(rl.Networks[0].Status, "Internal") {
			t.Errorf("GetResourceList: network status %q does not contain 'Internal'", rl.Networks[0].Status)
		}
	})

	t.Run("missing switch shown as missing in status", func(t *testing.T) {
		withPS(t, nil, makeHook("Off", "", false))

		rl, err := GetResourceList(f.Name())
		if err != nil {
			t.Fatalf("GetResourceList: expected nil, got %v", err)
		}
		if len(rl.Networks) != 1 {
			t.Fatalf("GetResourceList: expected 1 network, got %d", len(rl.Networks))
		}
		if !strings.Contains(rl.Networks[0].Status, "missing") {
			t.Errorf("GetResourceList: network status %q does not contain 'missing'", rl.Networks[0].Status)
		}
	})

	t.Run("invalid config file returns error", func(t *testing.T) {
		_, err := GetResourceList("/nonexistent/path.yaml")
		if err == nil {
			t.Fatal("GetResourceList: expected error for missing file, got nil")
		}
	})
}
