package hyperv

import (
	"errors"
	"strings"
	"testing"
)

// vmStateOutput returns a PS table output for a given VM state (used by IsVMExist).
func vmStateOutput(state string) []byte {
	return []byte("State\n-----\n" + state + "\n")
}

// --- GetVMInfo ---

func TestGetVMInfo(t *testing.T) {
	t.Run("VM exists outputPS called with batched script and returns result", func(t *testing.T) {
		callCount := 0
		const infoOutput = "VMName : my-vm\nCount  : 2\n"
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				// IsVMExist
				return vmStateOutput("Running"), nil
			}
			// batched Get-VMProcessor / Get-VMMemory script
			return []byte(infoOutput), nil
		})
		result, err := GetVMInfo("my-vm")
		if err != nil {
			t.Fatalf("GetVMInfo: expected nil, got %v", err)
		}
		if !strings.Contains(result, "my-vm") {
			t.Errorf("GetVMInfo: result %q does not contain 'my-vm'", result)
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		_, err := GetVMInfo("ghost-vm")
		if err == nil {
			t.Fatal("GetVMInfo: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("GetVMInfo: error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("PS error after VM exists returns error failed to get info", func(t *testing.T) {
		callCount := 0
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return vmStateOutput("Running"), nil
			}
			return nil, errors.New("ps error")
		})
		_, err := GetVMInfo("my-vm")
		if err == nil {
			t.Fatal("GetVMInfo: expected error on PS failure, got nil")
		}
		if !strings.Contains(err.Error(), "failed to get info") {
			t.Errorf("GetVMInfo: error %q does not contain 'failed to get info'", err.Error())
		}
	})
}

// --- MeasureVM ---

func TestMeasureVM(t *testing.T) {
	t.Run("VM exists outputPS called and returns result", func(t *testing.T) {
		callCount := 0
		const measureOutput = "AvgCPUUsage : 5\nAvgRAMUsage : 512\n"
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return vmStateOutput("Running"), nil
			}
			return []byte(measureOutput), nil
		})
		result, err := MeasureVM("my-vm")
		if err != nil {
			t.Fatalf("MeasureVM: expected nil, got %v", err)
		}
		if !strings.Contains(result, "AvgCPUUsage") {
			t.Errorf("MeasureVM: result %q does not contain 'AvgCPUUsage'", result)
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		_, err := MeasureVM("ghost-vm")
		if err == nil {
			t.Fatal("MeasureVM: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("MeasureVM: error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("PS error after VM exists returns error failed to measure", func(t *testing.T) {
		callCount := 0
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return vmStateOutput("Running"), nil
			}
			return nil, errors.New("ps error")
		})
		_, err := MeasureVM("my-vm")
		if err == nil {
			t.Fatal("MeasureVM: expected error on PS failure, got nil")
		}
		if !strings.Contains(err.Error(), "failed to measure") {
			t.Errorf("MeasureVM: error %q does not contain 'failed to measure'", err.Error())
		}
	})
}

// --- GetVMIntegrationServices ---

func TestGetVMIntegrationServices(t *testing.T) {
	t.Run("VM exists outputPS called and returns result", func(t *testing.T) {
		callCount := 0
		const svcOutput = "Name                    Enabled  PrimaryStatusDescription\n----                    -------  ------------------------\nTime Synchronization    True     OK\n"
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return vmStateOutput("Running"), nil
			}
			return []byte(svcOutput), nil
		})
		result, err := GetVMIntegrationServices("my-vm")
		if err != nil {
			t.Fatalf("GetVMIntegrationServices: expected nil, got %v", err)
		}
		if !strings.Contains(result, "Time Synchronization") {
			t.Errorf("GetVMIntegrationServices: result %q does not contain 'Time Synchronization'", result)
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		_, err := GetVMIntegrationServices("ghost-vm")
		if err == nil {
			t.Fatal("GetVMIntegrationServices: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("GetVMIntegrationServices: error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("PS error after VM exists returns error failed to get integration", func(t *testing.T) {
		callCount := 0
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return vmStateOutput("Running"), nil
			}
			return nil, errors.New("ps error")
		})
		_, err := GetVMIntegrationServices("my-vm")
		if err == nil {
			t.Fatal("GetVMIntegrationServices: expected error on PS failure, got nil")
		}
		if !strings.Contains(err.Error(), "failed to get integration") {
			t.Errorf("GetVMIntegrationServices: error %q does not contain 'failed to get integration'", err.Error())
		}
	})
}

// --- EnableVMIntegrationService ---

func TestEnableVMIntegrationService(t *testing.T) {
	t.Run("VM exists calls runPS and returns nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmStateOutput("Running"), nil
			},
		)
		if err := EnableVMIntegrationService("my-vm", "Time Synchronization"); err != nil {
			t.Fatalf("EnableVMIntegrationService: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("EnableVMIntegrationService: runPS was not called")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		err := EnableVMIntegrationService("ghost-vm", "Time Synchronization")
		if err == nil {
			t.Fatal("EnableVMIntegrationService: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("EnableVMIntegrationService: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

// --- DisableVMIntegrationService ---

func TestDisableVMIntegrationService(t *testing.T) {
	t.Run("VM exists calls runPS and returns nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmStateOutput("Running"), nil
			},
		)
		if err := DisableVMIntegrationService("my-vm", "Time Synchronization"); err != nil {
			t.Fatalf("DisableVMIntegrationService: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("DisableVMIntegrationService: runPS was not called")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		err := DisableVMIntegrationService("ghost-vm", "Time Synchronization")
		if err == nil {
			t.Fatal("DisableVMIntegrationService: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("DisableVMIntegrationService: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

// --- GetVMHost ---

func TestGetVMHost(t *testing.T) {
	t.Run("success outputPS called and returns result", func(t *testing.T) {
		const hostOutput = "VirtualHardDiskPath : C:\\VMs\nVirtualMachinePath  : C:\\VMs\n"
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte(hostOutput), nil
		})
		result, err := GetVMHost()
		if err != nil {
			t.Fatalf("GetVMHost: expected nil, got %v", err)
		}
		if !strings.Contains(result, "VirtualHardDiskPath") {
			t.Errorf("GetVMHost: result %q does not contain 'VirtualHardDiskPath'", result)
		}
	})

	t.Run("PS error returns error failed to get Hyper-V host", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("ps error")
		})
		_, err := GetVMHost()
		if err == nil {
			t.Fatal("GetVMHost: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to get Hyper-V host") {
			t.Errorf("GetVMHost: error %q does not contain 'failed to get Hyper-V host'", err.Error())
		}
	})
}
