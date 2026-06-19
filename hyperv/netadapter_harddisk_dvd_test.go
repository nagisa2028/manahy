package hyperv

import (
	"errors"
	"strings"
	"testing"
)

// vmRunningOutput returns PS table output for a running VM (used by IsVMExist).
func vmRunningOutput() []byte {
	return []byte("State\n-----\nRunning\n")
}

// vmMissingErr returns an error representing a missing VM.
func vmMissingErr() error {
	return errors.New("not found")
}

// =============================================================================
// netadapter
// =============================================================================

func TestGetVMNetworkAdapters(t *testing.T) {
	t.Run("VM exists outputPS called and returns result", func(t *testing.T) {
		callCount := 0
		const adapterOutput = "Name  SwitchName  MacAddress  Status\n----  ----------  ----------  ------\neth0  mySwitch    AABBCC      OK\n"
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return vmRunningOutput(), nil
			}
			return []byte(adapterOutput), nil
		})
		result, err := GetVMNetworkAdapters("my-vm")
		if err != nil {
			t.Fatalf("GetVMNetworkAdapters: expected nil, got %v", err)
		}
		if !strings.Contains(result, "eth0") {
			t.Errorf("GetVMNetworkAdapters: result %q does not contain 'eth0'", result)
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmMissingErr()
		})
		_, err := GetVMNetworkAdapters("ghost-vm")
		if err == nil {
			t.Fatal("GetVMNetworkAdapters: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("GetVMNetworkAdapters: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestAddVMNetworkAdapter(t *testing.T) {
	t.Run("VM exists calls runPS and returns nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmRunningOutput(), nil
			},
		)
		if err := AddVMNetworkAdapter("my-vm", "", ""); err != nil {
			t.Fatalf("AddVMNetworkAdapter: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("AddVMNetworkAdapter: runPS was not called")
		}
	})

	t.Run("with name and switch included in command", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error {
				capturedCmd = c
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmRunningOutput(), nil
			},
		)
		if err := AddVMNetworkAdapter("my-vm", "my-adapter", "my-switch"); err != nil {
			t.Fatalf("AddVMNetworkAdapter: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "my-adapter") {
			t.Errorf("AddVMNetworkAdapter: command %q does not contain 'my-adapter'", capturedCmd)
		}
		if !strings.Contains(capturedCmd, "my-switch") {
			t.Errorf("AddVMNetworkAdapter: command %q does not contain 'my-switch'", capturedCmd)
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmMissingErr()
		})
		err := AddVMNetworkAdapter("ghost-vm", "", "")
		if err == nil {
			t.Fatal("AddVMNetworkAdapter: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("AddVMNetworkAdapter: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestRemoveVMNetworkAdapter(t *testing.T) {
	t.Run("VM exists calls runPS and returns nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmRunningOutput(), nil
			},
		)
		if err := RemoveVMNetworkAdapter("my-vm", "eth0"); err != nil {
			t.Fatalf("RemoveVMNetworkAdapter: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("RemoveVMNetworkAdapter: runPS was not called")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmMissingErr()
		})
		err := RemoveVMNetworkAdapter("ghost-vm", "eth0")
		if err == nil {
			t.Fatal("RemoveVMNetworkAdapter: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("RemoveVMNetworkAdapter: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestConnectVMNetworkAdapter(t *testing.T) {
	t.Run("VM exists calls runPS and returns nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmRunningOutput(), nil
			},
		)
		if err := ConnectVMNetworkAdapter("my-vm", "eth0", "mySwitch"); err != nil {
			t.Fatalf("ConnectVMNetworkAdapter: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("ConnectVMNetworkAdapter: runPS was not called")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmMissingErr()
		})
		err := ConnectVMNetworkAdapter("ghost-vm", "eth0", "mySwitch")
		if err == nil {
			t.Fatal("ConnectVMNetworkAdapter: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("ConnectVMNetworkAdapter: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

// =============================================================================
// harddisk
// =============================================================================

func TestGetVMHardDiskDrives(t *testing.T) {
	t.Run("VM exists outputPS called and returns result", func(t *testing.T) {
		callCount := 0
		const diskOutput = "VMName  ControllerType  ControllerNumber  ControllerLocation  Path\n------  --------------  ----------------  ------------------  ----\nmy-vm   SCSI            0                 0                   C:\\disk.vhd\n"
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return vmRunningOutput(), nil
			}
			return []byte(diskOutput), nil
		})
		result, err := GetVMHardDiskDrives("my-vm")
		if err != nil {
			t.Fatalf("GetVMHardDiskDrives: expected nil, got %v", err)
		}
		if !strings.Contains(result, "disk.vhd") {
			t.Errorf("GetVMHardDiskDrives: result %q does not contain 'disk.vhd'", result)
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmMissingErr()
		})
		_, err := GetVMHardDiskDrives("ghost-vm")
		if err == nil {
			t.Fatal("GetVMHardDiskDrives: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("GetVMHardDiskDrives: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestAddVMHardDiskDrive(t *testing.T) {
	t.Run("VM exists and file exists calls runPS", func(t *testing.T) {
		// outputPS call 1: IsVMExist (GetVMState)
		// outputPS call 2: isFileExist (Test-Path)
		callCount := 0
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				callCount++
				if callCount == 1 {
					// IsVMExist
					return vmRunningOutput(), nil
				}
				// isFileExist (Test-Path)
				return []byte("True\n"), nil
			},
		)
		if err := AddVMHardDiskDrive("my-vm", "C:\\disk.vhd"); err != nil {
			t.Fatalf("AddVMHardDiskDrive: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("AddVMHardDiskDrive: runPS was not called")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmMissingErr()
		})
		err := AddVMHardDiskDrive("ghost-vm", "C:\\disk.vhd")
		if err == nil {
			t.Fatal("AddVMHardDiskDrive: expected error for missing VM, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("AddVMHardDiskDrive: error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("VM exists but file not found returns error", func(t *testing.T) {
		callCount := 0
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return vmRunningOutput(), nil
			}
			// Test-Path returns False
			return []byte("False\n"), nil
		})
		err := AddVMHardDiskDrive("my-vm", "C:\\missing.vhd")
		if err == nil {
			t.Fatal("AddVMHardDiskDrive: expected error for missing file, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("AddVMHardDiskDrive: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestRemoveVMHardDiskDrive(t *testing.T) {
	t.Run("VM exists calls runPS and returns nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmRunningOutput(), nil
			},
		)
		if err := RemoveVMHardDiskDrive("my-vm", "C:\\disk.vhd"); err != nil {
			t.Fatalf("RemoveVMHardDiskDrive: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("RemoveVMHardDiskDrive: runPS was not called")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmMissingErr()
		})
		err := RemoveVMHardDiskDrive("ghost-vm", "C:\\disk.vhd")
		if err == nil {
			t.Fatal("RemoveVMHardDiskDrive: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("RemoveVMHardDiskDrive: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

// =============================================================================
// dvddrive
// =============================================================================

func TestGetVMDvdDrives(t *testing.T) {
	t.Run("VM exists outputPS called and returns result", func(t *testing.T) {
		callCount := 0
		const dvdOutput = "VMName  ControllerType  ControllerNumber  ControllerLocation  Path\n------  --------------  ----------------  ------------------  ----\nmy-vm   IDE             1                 0                   C:\\image.iso\n"
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return vmRunningOutput(), nil
			}
			return []byte(dvdOutput), nil
		})
		result, err := GetVMDvdDrives("my-vm")
		if err != nil {
			t.Fatalf("GetVMDvdDrives: expected nil, got %v", err)
		}
		if !strings.Contains(result, "image.iso") {
			t.Errorf("GetVMDvdDrives: result %q does not contain 'image.iso'", result)
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmMissingErr()
		})
		_, err := GetVMDvdDrives("ghost-vm")
		if err == nil {
			t.Fatal("GetVMDvdDrives: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("GetVMDvdDrives: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestAddVMDvdDrive(t *testing.T) {
	t.Run("VM exists calls runPS and returns nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmRunningOutput(), nil
			},
		)
		if err := AddVMDvdDrive("my-vm"); err != nil {
			t.Fatalf("AddVMDvdDrive: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("AddVMDvdDrive: runPS was not called")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmMissingErr()
		})
		err := AddVMDvdDrive("ghost-vm")
		if err == nil {
			t.Fatal("AddVMDvdDrive: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("AddVMDvdDrive: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestRemoveVMDvdDrive(t *testing.T) {
	t.Run("VM exists calls runPS and returns nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmRunningOutput(), nil
			},
		)
		if err := RemoveVMDvdDrive("my-vm"); err != nil {
			t.Fatalf("RemoveVMDvdDrive: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("RemoveVMDvdDrive: runPS was not called")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmMissingErr()
		})
		err := RemoveVMDvdDrive("ghost-vm")
		if err == nil {
			t.Fatal("RemoveVMDvdDrive: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("RemoveVMDvdDrive: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestSetVMNetworkAdapterVlan(t *testing.T) {
	t.Run("VM exists sets VLAN access mode with ID in command", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error {
				capturedCmd = c
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmRunningOutput(), nil
			},
		)
		if err := SetVMNetworkAdapterVlan("my-vm", "eth0", 100); err != nil {
			t.Fatalf("SetVMNetworkAdapterVlan: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "-Access") {
			t.Errorf("SetVMNetworkAdapterVlan: command %q does not contain '-Access'", capturedCmd)
		}
		if !strings.Contains(capturedCmd, "100") {
			t.Errorf("SetVMNetworkAdapterVlan: command %q does not contain '100'", capturedCmd)
		}
		if strings.Contains(capturedCmd, "-Untagged") {
			t.Errorf("SetVMNetworkAdapterVlan: command %q should not contain '-Untagged'", capturedCmd)
		}
	})

	t.Run("vlanID 0 sets Untagged mode", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error {
				capturedCmd = c
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmRunningOutput(), nil
			},
		)
		if err := SetVMNetworkAdapterVlan("my-vm", "eth0", 0); err != nil {
			t.Fatalf("SetVMNetworkAdapterVlan: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "-Untagged") {
			t.Errorf("SetVMNetworkAdapterVlan: command %q does not contain '-Untagged'", capturedCmd)
		}
		if strings.Contains(capturedCmd, "-Access") {
			t.Errorf("SetVMNetworkAdapterVlan: command %q should not contain '-Access'", capturedCmd)
		}
	})

	t.Run("adapter name included in command", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error {
				capturedCmd = c
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmRunningOutput(), nil
			},
		)
		if err := SetVMNetworkAdapterVlan("my-vm", "my-adapter", 10); err != nil {
			t.Fatalf("SetVMNetworkAdapterVlan: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "my-adapter") {
			t.Errorf("SetVMNetworkAdapterVlan: command %q does not contain 'my-adapter'", capturedCmd)
		}
	})

	t.Run("PS error returns wrapped error message", func(t *testing.T) {
		withPS(t,
			func(_ string) error {
				return errors.New("ps error")
			},
			func(_ string) ([]byte, error) {
				return vmRunningOutput(), nil
			},
		)
		err := SetVMNetworkAdapterVlan("my-vm", "eth0", 100)
		if err == nil {
			t.Fatal("SetVMNetworkAdapterVlan: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to set VLAN") {
			t.Errorf("SetVMNetworkAdapterVlan: error %q does not contain 'failed to set VLAN'", err.Error())
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmMissingErr()
		})
		err := SetVMNetworkAdapterVlan("ghost-vm", "eth0", 10)
		if err == nil {
			t.Fatal("SetVMNetworkAdapterVlan: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("SetVMNetworkAdapterVlan: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestSetVMDvdDrive(t *testing.T) {
	t.Run("VM exists and empty imagePath uses $null in command", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error {
				capturedCmd = c
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmRunningOutput(), nil
			},
		)
		if err := SetVMDvdDrive("my-vm", ""); err != nil {
			t.Fatalf("SetVMDvdDrive: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "$null") {
			t.Errorf("SetVMDvdDrive: command %q does not contain '$null' for empty imagePath", capturedCmd)
		}
	})

	t.Run("VM exists and imagePath set uses quoted path in command", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error {
				capturedCmd = c
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmRunningOutput(), nil
			},
		)
		if err := SetVMDvdDrive("my-vm", "C:\\image.iso"); err != nil {
			t.Fatalf("SetVMDvdDrive: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "C:\\image.iso") {
			t.Errorf("SetVMDvdDrive: command %q does not contain 'C:\\image.iso'", capturedCmd)
		}
		if strings.Contains(capturedCmd, "$null") {
			t.Errorf("SetVMDvdDrive: command %q should not contain '$null' when imagePath is set", capturedCmd)
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmMissingErr()
		})
		err := SetVMDvdDrive("ghost-vm", "C:\\image.iso")
		if err == nil {
			t.Fatal("SetVMDvdDrive: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("SetVMDvdDrive: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

// =============================================================================
// runPS failure paths
// =============================================================================

func TestAddVMNetworkAdapterRunPSFailure(t *testing.T) {
	withPS(t,
		func(_ string) error { return errors.New("runPS failed") },
		func(_ string) ([]byte, error) { return vmRunningOutput(), nil },
	)
	err := AddVMNetworkAdapter("my-vm", "", "")
	if err == nil {
		t.Fatal("AddVMNetworkAdapter: expected error from runPS, got nil")
	}
}

func TestRemoveVMNetworkAdapterRunPSFailure(t *testing.T) {
	withPS(t,
		func(_ string) error { return errors.New("runPS failed") },
		func(_ string) ([]byte, error) { return vmRunningOutput(), nil },
	)
	err := RemoveVMNetworkAdapter("my-vm", "eth0")
	if err == nil {
		t.Fatal("RemoveVMNetworkAdapter: expected error from runPS, got nil")
	}
}

func TestConnectVMNetworkAdapterRunPSFailure(t *testing.T) {
	withPS(t,
		func(_ string) error { return errors.New("runPS failed") },
		func(_ string) ([]byte, error) { return vmRunningOutput(), nil },
	)
	err := ConnectVMNetworkAdapter("my-vm", "eth0", "mySwitch")
	if err == nil {
		t.Fatal("ConnectVMNetworkAdapter: expected error from runPS, got nil")
	}
}

func TestAddVMHardDiskDriveRunPSFailure(t *testing.T) {
	callCount := 0
	withPS(t,
		func(_ string) error { return errors.New("runPS failed") },
		func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return vmRunningOutput(), nil
			}
			return []byte("True\n"), nil
		},
	)
	err := AddVMHardDiskDrive("my-vm", `C:\disk.vhd`)
	if err == nil {
		t.Fatal("AddVMHardDiskDrive: expected error from runPS, got nil")
	}
}

func TestRemoveVMHardDiskDriveRunPSFailure(t *testing.T) {
	withPS(t,
		func(_ string) error { return errors.New("runPS failed") },
		func(_ string) ([]byte, error) { return vmRunningOutput(), nil },
	)
	err := RemoveVMHardDiskDrive("my-vm", `C:\disk.vhd`)
	if err == nil {
		t.Fatal("RemoveVMHardDiskDrive: expected error from runPS, got nil")
	}
}

func TestAddVMDvdDriveRunPSFailure(t *testing.T) {
	withPS(t,
		func(_ string) error { return errors.New("runPS failed") },
		func(_ string) ([]byte, error) { return vmRunningOutput(), nil },
	)
	err := AddVMDvdDrive("my-vm")
	if err == nil {
		t.Fatal("AddVMDvdDrive: expected error from runPS, got nil")
	}
}

func TestRemoveVMDvdDriveRunPSFailure(t *testing.T) {
	withPS(t,
		func(_ string) error { return errors.New("runPS failed") },
		func(_ string) ([]byte, error) { return vmRunningOutput(), nil },
	)
	err := RemoveVMDvdDrive("my-vm")
	if err == nil {
		t.Fatal("RemoveVMDvdDrive: expected error from runPS, got nil")
	}
}

func TestSetVMDvdDriveRunPSFailure(t *testing.T) {
	withPS(t,
		func(_ string) error { return errors.New("runPS failed") },
		func(_ string) ([]byte, error) { return vmRunningOutput(), nil },
	)
	err := SetVMDvdDrive("my-vm", "")
	if err == nil {
		t.Fatal("SetVMDvdDrive: expected error from runPS, got nil")
	}
}
