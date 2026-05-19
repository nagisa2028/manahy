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

	// makeHook returns a hook for the batch queries used by GetResourceList.
	// vmState is returned for "test-vm"; switchType (empty string → missing) for "test-switch".
	makeHook := func(vmState string, switchType string, diskPresent bool) func(string) ([]byte, error) {
		return func(cmd string) ([]byte, error) {
			// Check Get-VMSwitch before Get-VM because the latter is a substring of the former.
			if strings.Contains(cmd, "Get-VMSwitch") {
				if switchType == "" {
					return nil, errors.New("not found")
				}
				return []byte("Name        SwitchType\n----        ----------\ntest-switch " + switchType + "\n"), nil
			}
			if strings.Contains(cmd, "Get-VM") {
				return []byte("Name     State\n----     -----\ntest-vm  " + vmState + "\n"), nil
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

	t.Run("VM not on host shown as NotFound", func(t *testing.T) {
		withPS(t, nil, makeHook("Off", "Internal", false))
		// makeHook returns "test-vm Off"; to simulate a missing VM we use an empty batch response.
		withPS(t, nil, func(cmd string) ([]byte, error) {
			if strings.Contains(cmd, "Get-VMSwitch") {
				return []byte("Name        SwitchType\n----        ----------\ntest-switch Internal\n"), nil
			}
			if strings.Contains(cmd, "Get-VM") {
				// Empty batch: no VMs on host.
				return []byte("Name  State\n----  -----\n"), nil
			}
			return []byte("True\n"), nil
		})

		rl, err := GetResourceList(f.Name())
		if err != nil {
			t.Fatalf("GetResourceList: expected nil, got %v", err)
		}
		if len(rl.VMs) != 1 {
			t.Fatalf("GetResourceList: expected 1 VM, got %d", len(rl.VMs))
		}
		if rl.VMs[0].Status != vmStateNotFound {
			t.Errorf("GetResourceList: VM status %q, want %q", rl.VMs[0].Status, vmStateNotFound)
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

	t.Run("count>1 VM expands to numbered instances", func(t *testing.T) {
		const multiYAML = `
vms:
  router:
    count: 3
    generation: 1
    cpu:
      thread: 1
    memory:
      size: 512MB
    path: C:\VMs
disks: {}
networks: {}
`
		mf, err := os.CreateTemp(t.TempDir(), "manahy-multi-*.yaml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		if _, err = mf.WriteString(multiYAML); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}
		_ = mf.Close()

		withPS(t, nil, func(cmd string) ([]byte, error) {
			if strings.Contains(cmd, "Get-VM") {
				return []byte("Name     State\n----     -----\nrouter1  Running\nrouter2  Running\nrouter3  Running\n"), nil
			}
			return []byte("True\n"), nil
		})

		rl, err := GetResourceList(mf.Name())
		if err != nil {
			t.Fatalf("GetResourceList count>1: expected nil, got %v", err)
		}
		if len(rl.VMs) != 3 {
			t.Fatalf("GetResourceList count>1: expected 3 VMs, got %d: %v", len(rl.VMs), rl.VMs)
		}
		for i, want := range []string{"router1", "router2", "router3"} {
			if rl.VMs[i].Name != want {
				t.Errorf("GetResourceList count>1: VMs[%d].Name = %q, want %q", i, rl.VMs[i].Name, want)
			}
		}
	})

	t.Run("multiRef disk expands to numbered paths", func(t *testing.T) {
		const multiDiskYAML = `
vms:
  router:
    count: 2
    generation: 1
    cpu:
      thread: 1
    memory:
      size: 512MB
    path: C:\VMs
    disks:
      - rtr-disk
networks: {}
disks:
  rtr-disk:
    path: C:\Hyper-V\Disks\router.vhd
    size: 10GB
    type: dynamic
`
		mf, err := os.CreateTemp(t.TempDir(), "manahy-multidisk-*.yaml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		if _, err = mf.WriteString(multiDiskYAML); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}
		_ = mf.Close()

		withPS(t, nil, func(cmd string) ([]byte, error) {
			if strings.Contains(cmd, "Get-VM") {
				return []byte("Name     State\n----     -----\nrouter1  Running\nrouter2  Running\n"), nil
			}
			return []byte("True\n"), nil // all Test-Path return present
		})

		rl, err := GetResourceList(mf.Name())
		if err != nil {
			t.Fatalf("GetResourceList multiRef disk: expected nil, got %v", err)
		}
		if len(rl.Disks) != 2 {
			t.Fatalf("GetResourceList multiRef disk: expected 2 disks, got %d: %v", len(rl.Disks), rl.Disks)
		}
		for i, want := range []string{"rtr-disk1", "rtr-disk2"} {
			if rl.Disks[i].Name != want {
				t.Errorf("GetResourceList multiRef disk: Disks[%d].Name = %q, want %q", i, rl.Disks[i].Name, want)
			}
		}
	})
}

// TestGetResourceListConcurrent verifies correctness of the concurrent implementation.
// Run with -race to detect data races.
func TestGetResourceListConcurrent(t *testing.T) {
	// richYAML has multiple resources in every group so all 3 goroutines
	// do real work simultaneously and append to their respective slices.
	const richYAML = `
vms:
  alpha:
    count: 1
    generation: 1
    cpu:
      thread: 1
    memory:
      size: 512MB
    path: C:\VMs
  beta:
    count: 3
    generation: 1
    cpu:
      thread: 1
    memory:
      size: 512MB
    path: C:\VMs
networks:
  net-a:
    type: internal
  net-b:
    type: private
  net-c:
    type: internal
disks:
  disk-x:
    path: C:\disks\x.vhd
  disk-y:
    path: C:\disks\y.vhd
`
	const vmBatchRunning = "Name   State\n----   -----\nalpha  Running\nbeta1  Running\nbeta2  Running\nbeta3  Running\n"
	const vmBatchOff = "Name   State\n----   -----\nalpha  Off\nbeta1  Off\nbeta2  Off\nbeta3  Off\n"
	const switchBatch = "Name   SwitchType\n----   ----------\nnet-a  Internal\nnet-b  Private\nnet-c  Internal\n"

	writeYAML := func(t *testing.T, content string) string {
		t.Helper()
		mf, err := os.CreateTemp(t.TempDir(), "manahy-concurrent-*.yaml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		if _, err = mf.WriteString(content); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}
		_ = mf.Close()
		return mf.Name()
	}

	t.Run("all groups complete and counts are exact", func(t *testing.T) {
		withPS(t, nil, func(cmd string) ([]byte, error) {
			if strings.Contains(cmd, "Get-VMSwitch") {
				return []byte(switchBatch), nil
			}
			if strings.Contains(cmd, "Get-VM") {
				return []byte(vmBatchRunning), nil
			}
			return []byte("True\n"), nil
		})

		rl, err := GetResourceList(writeYAML(t, richYAML))
		if err != nil {
			t.Fatalf("GetResourceList: expected nil, got %v", err)
		}
		// alpha(count=1) + beta1 + beta2 + beta3 = 4 VMs
		if len(rl.VMs) != 4 {
			t.Errorf("VMs: expected 4, got %d: %v", len(rl.VMs), rl.VMs)
		}
		// disk-x + disk-y = 2 disks
		if len(rl.Disks) != 2 {
			t.Errorf("Disks: expected 2, got %d: %v", len(rl.Disks), rl.Disks)
		}
		// net-a + net-b + net-c = 3 networks
		if len(rl.Networks) != 3 {
			t.Errorf("Networks: expected 3, got %d: %v", len(rl.Networks), rl.Networks)
		}
	})

	t.Run("output is sorted despite concurrent execution", func(t *testing.T) {
		withPS(t, nil, func(cmd string) ([]byte, error) {
			if strings.Contains(cmd, "Get-VMSwitch") {
				return []byte(switchBatch), nil
			}
			if strings.Contains(cmd, "Get-VM") {
				return []byte(vmBatchOff), nil
			}
			return []byte("True\n"), nil
		})

		rl, err := GetResourceList(writeYAML(t, richYAML))
		if err != nil {
			t.Fatalf("GetResourceList: expected nil, got %v", err)
		}
		assertSorted(t, "VMs", rl.VMs)
		assertSorted(t, "Disks", rl.Disks)
		assertSorted(t, "Networks", rl.Networks)
	})

	t.Run("no data race under repeated concurrent calls", func(t *testing.T) {
		// Call GetResourceList many times to stress the goroutine implementation.
		// Run this test with -race to catch any data races.
		withPS(t, nil, func(cmd string) ([]byte, error) {
			if strings.Contains(cmd, "Get-VMSwitch") {
				return []byte(switchBatch), nil
			}
			if strings.Contains(cmd, "Get-VM") {
				return []byte(vmBatchRunning), nil
			}
			return []byte("True\n"), nil
		})
		path := writeYAML(t, richYAML)
		for i := 0; i < 20; i++ {
			rl, err := GetResourceList(path)
			if err != nil {
				t.Fatalf("iteration %d: unexpected error: %v", i, err)
			}
			if len(rl.VMs) != 4 || len(rl.Disks) != 2 || len(rl.Networks) != 3 {
				t.Fatalf("iteration %d: unexpected counts VMs=%d Disks=%d Networks=%d",
					i, len(rl.VMs), len(rl.Disks), len(rl.Networks))
			}
		}
	})
}

// assertSorted verifies that the ResourceStatus slice is in ascending name order.
func assertSorted(t *testing.T, label string, rs []ResourceStatus) {
	t.Helper()
	for i := 1; i < len(rs); i++ {
		if rs[i].Name < rs[i-1].Name {
			t.Errorf("%s: not sorted at index %d: %q < %q", label, i, rs[i].Name, rs[i-1].Name)
		}
	}
}
