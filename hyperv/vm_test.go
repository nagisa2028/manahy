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
