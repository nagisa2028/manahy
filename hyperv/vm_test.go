package hyperv

import (
	"errors"
	"strings"
	"testing"
)

// stateOutput returns the canonical PS table output for a given VM state string.
func stateOutput(state string) []byte {
	return []byte("State\n-----\n" + state + "\n")
}

// ---------- pure-validation helpers (no PS needed) ----------

func TestCheckVMGeneration(t *testing.T) {
	tests := []struct {
		generation int
		wantErr    bool
	}{
		{1, false},
		{2, false},
		{0, true},
		{3, true},
		{-1, true},
	}
	for _, tc := range tests {
		err := checkVMGeneration(tc.generation)
		if tc.wantErr && err == nil {
			t.Errorf("checkVMGeneration(%d): expected error, got nil", tc.generation)
		}
		if !tc.wantErr && err != nil {
			t.Errorf("checkVMGeneration(%d): unexpected error: %v", tc.generation, err)
		}
	}
}

func TestCheckVMProcessor(t *testing.T) {
	tests := []struct {
		name    string
		cpu     CPU
		wantErr bool
	}{
		{"single thread", CPU{Thread: 1}, false},
		{"4 threads nested", CPU{Thread: 4, Nested: true}, false},
		{"zero threads", CPU{Thread: 0}, true},
		{"negative threads", CPU{Thread: -1}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkVMProcessor(tc.cpu)
			if tc.wantErr && err == nil {
				t.Errorf("checkVMProcessor(%+v): expected error, got nil", tc.cpu)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("checkVMProcessor(%+v): unexpected error: %v", tc.cpu, err)
			}
		})
	}
}

func TestCheckMemorySize(t *testing.T) {
	tests := []struct {
		name    string
		size    string
		wantErr bool
	}{
		{"GB valid", "1GB", false},
		{"MB valid", "512MB", false},
		{"TB valid", "2TB", false},
		{"multi digit", "1024MB", false},
		{"bare number no unit", "1024", true},
		{"KB not supported", "512KB", true},
		{"empty string", "", true},
		{"lowercase", "1gb", true},
		{"with space", "1 GB", true},
		{"16TB at limit", "16TB", false},
		{"17TB exceeds limit", "17TB", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkMemorySize(tc.size)
			if tc.wantErr && err == nil {
				t.Errorf("checkMemorySize(%q): expected error, got nil", tc.size)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("checkMemorySize(%q): unexpected error: %v", tc.size, err)
			}
		})
	}
}

// ---------- SetVMProcessor ----------

func TestSetVMProcessor(t *testing.T) {
	t.Run("VM exists calls runPS with count and nested flag", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error { capturedCmd = c; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		if err := SetVMProcessor("my-vm", CPU{Thread: 4, Nested: true}); err != nil {
			t.Fatalf("SetVMProcessor: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "-Count 4") {
			t.Errorf("SetVMProcessor: command %q does not contain '-Count 4'", capturedCmd)
		}
		if !strings.Contains(capturedCmd, "$true") {
			t.Errorf("SetVMProcessor: command %q does not contain '$true'", capturedCmd)
		}
	})

	t.Run("invalid CPU count returns error without runPS call", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		err := SetVMProcessor("my-vm", CPU{Thread: 0})
		if err == nil {
			t.Fatal("SetVMProcessor: expected error for Thread=0, got nil")
		}
		if runCalled {
			t.Error("SetVMProcessor: runPS should not be called for invalid CPU")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		err := SetVMProcessor("ghost-vm", CPU{Thread: 2})
		if err == nil {
			t.Fatal("SetVMProcessor: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("SetVMProcessor: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

// ---------- SetVMMemory ----------

func TestSetVMMemory(t *testing.T) {
	t.Run("VM exists calls runPS with size and dynamic flag", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error { capturedCmd = c; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Off"), nil },
		)
		if err := SetVMMemory("my-vm", Memory{Size: "1GB", Dynamic: false}); err != nil {
			t.Fatalf("SetVMMemory: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "1GB") {
			t.Errorf("SetVMMemory: command %q does not contain '1GB'", capturedCmd)
		}
		if !strings.Contains(capturedCmd, "$false") {
			t.Errorf("SetVMMemory: command %q does not contain '$false'", capturedCmd)
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		err := SetVMMemory("ghost-vm", Memory{Size: "512MB"})
		if err == nil {
			t.Fatal("SetVMMemory: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("SetVMMemory: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

// ---------- SetVMHardDisk ----------

func TestSetVMHardDisk(t *testing.T) {
	t.Run("two disks calls runPS twice", func(t *testing.T) {
		runCount := 0
		callCount := 0
		withPS(t,
			func(_ string) error { runCount++; return nil },
			func(_ string) ([]byte, error) {
				callCount++
				if callCount == 1 {
					return stateOutput("Off"), nil // IsVMExist
				}
				return []byte("True\n"), nil // Test-Path for each disk
			},
		)
		if err := SetVMHardDisk("my-vm", []string{`C:\disk1.vhd`, `C:\disk2.vhd`}); err != nil {
			t.Fatalf("SetVMHardDisk: expected nil, got %v", err)
		}
		if runCount != 2 {
			t.Errorf("SetVMHardDisk: expected runPS called 2 times, got %d", runCount)
		}
	})

	t.Run("empty disk list returns nil without runPS", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Off"), nil },
		)
		if err := SetVMHardDisk("my-vm", []string{}); err != nil {
			t.Fatalf("SetVMHardDisk: expected nil for empty list, got %v", err)
		}
		if runCalled {
			t.Error("SetVMHardDisk: runPS should not be called for empty disk list")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		err := SetVMHardDisk("ghost-vm", []string{`C:\disk.vhd`})
		if err == nil {
			t.Fatal("SetVMHardDisk: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("SetVMHardDisk: error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("disk file not found returns error", func(t *testing.T) {
		callCount := 0
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return stateOutput("Off"), nil // IsVMExist
			}
			return []byte("False\n"), nil // Test-Path: disk not found
		})
		err := SetVMHardDisk("my-vm", []string{`C:\missing.vhd`})
		if err == nil {
			t.Fatal("SetVMHardDisk: expected error for missing disk, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("SetVMHardDisk: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

// ---------- SetVMImageFile ----------

func TestSetVMImageFile(t *testing.T) {
	t.Run("VM and image file exist calls runPS", func(t *testing.T) {
		runCalled := false
		callCount := 0
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) {
				callCount++
				if callCount == 1 {
					return stateOutput("Off"), nil // IsVMExist
				}
				return []byte("True\n"), nil // isFileExist for image
			},
		)
		if err := SetVMImageFile("my-vm", `C:\image.iso`); err != nil {
			t.Fatalf("SetVMImageFile: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("SetVMImageFile: runPS was not called")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		err := SetVMImageFile("ghost-vm", `C:\image.iso`)
		if err == nil {
			t.Fatal("SetVMImageFile: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("SetVMImageFile: error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("image file not found returns error", func(t *testing.T) {
		callCount := 0
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return stateOutput("Off"), nil // IsVMExist
			}
			return []byte("False\n"), nil // isFileExist: image not found
		})
		err := SetVMImageFile("my-vm", `C:\missing.iso`)
		if err == nil {
			t.Fatal("SetVMImageFile: expected error for missing image, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("SetVMImageFile: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

// ---------- SetVMSwitch ----------

func TestSetVMSwitch(t *testing.T) {
	t.Run("switch exists calls runPS", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) {
				return []byte("SwitchType\n----------\nInternal\n"), nil
			},
		)
		if err := SetVMSwitch("my-vm", []string{"mySwitch"}); err != nil {
			t.Fatalf("SetVMSwitch: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("SetVMSwitch: runPS was not called")
		}
	})

	t.Run("empty switch list returns nil without runPS", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) {
				return []byte("SwitchType\n----------\nInternal\n"), nil
			},
		)
		if err := SetVMSwitch("my-vm", []string{}); err != nil {
			t.Fatalf("SetVMSwitch: expected nil for empty list, got %v", err)
		}
		if runCalled {
			t.Error("SetVMSwitch: runPS should not be called for empty list")
		}
	})

	t.Run("switch not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		err := SetVMSwitch("my-vm", []string{"missing-switch"})
		if err == nil {
			t.Fatal("SetVMSwitch: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("SetVMSwitch: error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("switch unknown state returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("SwitchType\n----------\n"), nil // Unknown
		})
		err := SetVMSwitch("my-vm", []string{"ambiguous-switch"})
		if err == nil {
			t.Fatal("SetVMSwitch: expected error for Unknown switch state, got nil")
		}
		if !strings.Contains(err.Error(), "failed to get") {
			t.Errorf("SetVMSwitch: error %q does not contain 'failed to get'", err.Error())
		}
	})
}

// ---------- CreateVM ----------

func TestCheckVMPath(t *testing.T) {
	t.Run("relative path returns error without PS call", func(t *testing.T) {
		outputCalled := false
		withPS(t, nil, func(_ string) ([]byte, error) {
			outputCalled = true
			return []byte("False\n"), nil
		})
		err := checkVMPath("my-vm", `VMs`)
		if err == nil {
			t.Fatal("checkVMPath: expected error for relative path, got nil")
		}
		if !strings.Contains(err.Error(), "absolute") {
			t.Errorf("checkVMPath: error %q does not contain 'absolute'", err.Error())
		}
		if outputCalled {
			t.Error("checkVMPath: PS should not be called for relative path")
		}
	})
}

func TestCreateVM(t *testing.T) {
	// makeCreateMock returns an outputPS mock for the minimal VM creation sequence:
	// call 1: IsNotVMExist (Get-VM) → not found
	// call 2: isNotFileExist (Test-Path for vm path) → False (path free)
	// subsequent Get-VM calls: Running (IsVMExist in Set* functions)
	makeCreateMock := func() func(string) ([]byte, error) {
		callCount := 0
		return func(cmd string) ([]byte, error) {
			if strings.Contains(cmd, "Test-Path") {
				return []byte("False\n"), nil
			}
			callCount++
			if callCount == 1 {
				return nil, errors.New("not found") // IsNotVMExist → VM free
			}
			return stateOutput("Running"), nil // subsequent IsVMExist calls
		}
	}

	minimalVM := VM{
		Name:       "test-vm",
		Generation: 1,
		Path:       `C:\VMs`,
		Memory:     Memory{Size: "512MB"},
		CPU:        CPU{Thread: 1},
	}

	t.Run("minimal VM no disks no image no networks succeeds", func(t *testing.T) {
		runCount := 0
		withPS(t,
			func(_ string) error { runCount++; return nil },
			makeCreateMock(),
		)
		if err := CreateVM(minimalVM, false); err != nil {
			t.Fatalf("CreateVM: expected nil, got %v", err)
		}
		// New-VM + Set-VMProcessor + Set-VMMemory = 3 runPS calls
		if runCount != 3 {
			t.Errorf("CreateVM: expected 3 runPS calls, got %d", runCount)
		}
	})

	t.Run("VM already exists returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return stateOutput("Running"), nil // IsNotVMExist sees existing VM
		})
		err := CreateVM(minimalVM, false)
		if err == nil {
			t.Fatal("CreateVM: expected error for existing VM, got nil")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("CreateVM: error %q does not contain 'already exists'", err.Error())
		}
	})

	t.Run("invalid generation returns error without runPS", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return nil, errors.New("not found") },
		)
		vm := minimalVM
		vm.Generation = 3
		err := CreateVM(vm, false)
		if err == nil {
			t.Fatal("CreateVM: expected error for generation=3, got nil")
		}
		if !strings.Contains(err.Error(), "generation") {
			t.Errorf("CreateVM: error %q does not contain 'generation'", err.Error())
		}
		if runCalled {
			t.Error("CreateVM: runPS should not be called for invalid params")
		}
	})

	t.Run("New-VM runPS failure returns error", func(t *testing.T) {
		withPS(t,
			func(_ string) error { return errors.New("access denied") },
			makeCreateMock(),
		)
		err := CreateVM(minimalVM, false)
		if err == nil {
			t.Fatal("CreateVM: expected error when New-VM fails, got nil")
		}
		if !strings.Contains(err.Error(), "failed to create VM") {
			t.Errorf("CreateVM: error %q does not contain 'failed to create VM'", err.Error())
		}
	})

	t.Run("image skipped when Image is empty", func(t *testing.T) {
		var capturedCmds []string
		withPS(t,
			func(c string) error { capturedCmds = append(capturedCmds, c); return nil },
			makeCreateMock(),
		)
		vm := minimalVM
		vm.Image = ""
		if err := CreateVM(vm, false); err != nil {
			t.Fatalf("CreateVM: expected nil, got %v", err)
		}
		for _, cmd := range capturedCmds {
			if strings.Contains(cmd, "Add-VMDvdDrive") {
				t.Errorf("CreateVM: Add-VMDvdDrive called when Image is empty: %q", cmd)
			}
		}
	})
}

// ---------- GetVMState ----------

func TestGetVMState(t *testing.T) {
	tests := []struct {
		name      string
		psOutput  []byte
		psErr     error
		wantState string
	}{
		{"running", stateOutput("Running"), nil, "Running"},
		{"off", stateOutput("Off"), nil, "Off"},
		{"saved", stateOutput("Saved"), nil, "Saved"},
		{"paused", stateOutput("Paused"), nil, "Paused"},
		{"ps error → NotFound", nil, errors.New("not found"), "NotFound"},
		{"empty output → Unknown", []byte("State\n-----\n"), nil, "Unknown"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			withPS(t, nil, func(_ string) ([]byte, error) {
				return tc.psOutput, tc.psErr
			})
			got := GetVMState("test-vm")
			if got != tc.wantState {
				t.Errorf("GetVMState: got %q; want %q", got, tc.wantState)
			}
		})
	}
}

// ---------- IsVMExist ----------

func TestIsVMExist(t *testing.T) {
	t.Run("running VM → nil", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return stateOutput("Running"), nil
		})
		if err := IsVMExist("my-vm"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("NotFound → error containing 'does not exist'", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		err := IsVMExist("my-vm")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("Unknown → error containing 'failed to get'", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("State\n-----\n"), nil
		})
		err := IsVMExist("my-vm")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to get") {
			t.Errorf("error %q does not contain 'failed to get'", err.Error())
		}
	})
}

// ---------- IsNotVMExist ----------

func TestIsNotVMExist(t *testing.T) {
	t.Run("NotFound → nil", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		if err := IsNotVMExist("my-vm"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("running VM → error containing 'already exists'", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return stateOutput("Running"), nil
		})
		err := IsNotVMExist("my-vm")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("error %q does not contain 'already exists'", err.Error())
		}
	})

	t.Run("Unknown → error containing 'failed to get'", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("State\n-----\n"), nil
		})
		err := IsNotVMExist("my-vm")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to get") {
			t.Errorf("error %q does not contain 'failed to get'", err.Error())
		}
	})
}

// ---------- StartVM ----------

func TestStartVM(t *testing.T) {
	t.Run("not running → runPS called → nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Off"), nil },
		)
		if err := StartVM("my-vm"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !runCalled {
			t.Error("expected runPS to be called")
		}
	})

	t.Run("already Running → error without PS run call", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		err := StartVM("my-vm")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if runCalled {
			t.Error("runPS should not be called when VM is already Running")
		}
	})
}

// ---------- StopVM ----------

func TestStopVM(t *testing.T) {
	t.Run("running → runPS called → nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		if err := StopVM("my-vm"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !runCalled {
			t.Error("expected runPS to be called")
		}
	})

	t.Run("not running → error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return stateOutput("Off"), nil
		})
		if err := StopVM("my-vm"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ---------- DestroyVM ----------

func TestDestroyVM(t *testing.T) {
	t.Run("running → runPS called → nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		if err := DestroyVM("my-vm"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !runCalled {
			t.Error("expected runPS to be called")
		}
	})

	t.Run("not running → error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return stateOutput("Off"), nil
		})
		if err := DestroyVM("my-vm"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ---------- SaveVM ----------

func TestSaveVM(t *testing.T) {
	t.Run("running → runPS called → nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		if err := SaveVM("my-vm"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !runCalled {
			t.Error("expected runPS to be called")
		}
	})

	t.Run("not running → error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return stateOutput("Saved"), nil
		})
		if err := SaveVM("my-vm"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ---------- SuspendVM ----------

func TestSuspendVM(t *testing.T) {
	t.Run("running → runPS called → nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		if err := SuspendVM("my-vm"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !runCalled {
			t.Error("expected runPS to be called")
		}
	})

	t.Run("not running → error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return stateOutput("Off"), nil
		})
		if err := SuspendVM("my-vm"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ---------- RestartVM ----------

func TestRestartVM(t *testing.T) {
	t.Run("running → runPS called → nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		if err := RestartVM("my-vm"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !runCalled {
			t.Error("expected runPS to be called")
		}
	})

	t.Run("not running → error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return stateOutput("Off"), nil
		})
		if err := RestartVM("my-vm"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ---------- ResumeVM ----------

func TestResumeVM(t *testing.T) {
	t.Run("paused → runPS called → nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Paused"), nil },
		)
		if err := ResumeVM("my-vm"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !runCalled {
			t.Error("expected runPS to be called")
		}
	})

	t.Run("not paused → error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return stateOutput("Running"), nil
		})
		if err := ResumeVM("my-vm"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ---------- RemoveVM ----------

func TestRemoveVM(t *testing.T) {
	t.Run("exists → runPS called → nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		if err := RemoveVM("my-vm", false); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !runCalled {
			t.Error("expected runPS to be called")
		}
	})

	t.Run("not found → error from IsVMExist", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		if err := RemoveVM("ghost-vm", false); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ---------- RenameVM ----------

func TestRenameVM(t *testing.T) {
	// RenameVM flow:
	//   call 1 (outputPS): IsVMExist("src")  → checks source exists
	//   call 2 (outputPS): IsNotVMExist("dst") → checks destination is free
	//   call 3 (runPS):   Rename-VM
	t.Run("exists + new name free → rename succeeds", func(t *testing.T) {
		callCount := 0
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) {
				callCount++
				switch callCount {
				case 1:
					// IsVMExist("src") → source exists
					return stateOutput("Running"), nil
				default:
					// IsNotVMExist("dst") → destination not found
					return nil, errors.New("not found")
				}
			},
		)
		if err := RenameVM("src", "dst"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !runCalled {
			t.Error("expected runPS to be called for the rename")
		}
	})

	t.Run("source not found → error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		if err := RenameVM("ghost", "dst"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("destination already exists → error", func(t *testing.T) {
		callCount := 0
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			switch callCount {
			case 1:
				// IsVMExist("src") → source exists
				return stateOutput("Running"), nil
			default:
				// IsNotVMExist("dst") → destination already exists
				return stateOutput("Off"), nil
			}
		})
		err := RenameVM("src", "existing-dst")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("error %q does not contain 'already exists'", err.Error())
		}
	})
}
