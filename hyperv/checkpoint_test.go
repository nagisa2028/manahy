package hyperv

import (
	"errors"
	"strings"
	"testing"
)

// vmExistsOutput is the canonical PS table output for a running VM (used by IsVMExist).
func vmExistsOutput() []byte {
	return []byte("State\n-----\nRunning\n")
}

// vmNotFoundErr is the error returned when a VM does not exist.
func vmNotFoundErr() error {
	return errors.New("not found")
}

// --- CreateCheckpoint ---

func TestCreateCheckpoint(t *testing.T) {
	t.Run("VM exists calls runPS and returns nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmExistsOutput(), nil
			},
		)
		if err := CreateCheckpoint("my-vm", "snap1"); err != nil {
			t.Fatalf("CreateCheckpoint: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("CreateCheckpoint: runPS was not called")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmNotFoundErr()
		})
		err := CreateCheckpoint("ghost-vm", "snap1")
		if err == nil {
			t.Fatal("CreateCheckpoint: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("CreateCheckpoint: error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("with name includes SnapshotName in command", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error {
				capturedCmd = c
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmExistsOutput(), nil
			},
		)
		if err := CreateCheckpoint("my-vm", "my-snapshot"); err != nil {
			t.Fatalf("CreateCheckpoint: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "SnapshotName") {
			t.Errorf("CreateCheckpoint: command %q does not contain 'SnapshotName'", capturedCmd)
		}
		if !strings.Contains(capturedCmd, "my-snapshot") {
			t.Errorf("CreateCheckpoint: command %q does not contain 'my-snapshot'", capturedCmd)
		}
	})

	t.Run("without name command has no SnapshotName", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error {
				capturedCmd = c
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmExistsOutput(), nil
			},
		)
		if err := CreateCheckpoint("my-vm", ""); err != nil {
			t.Fatalf("CreateCheckpoint: expected nil, got %v", err)
		}
		if strings.Contains(capturedCmd, "SnapshotName") {
			t.Errorf("CreateCheckpoint: command %q should not contain 'SnapshotName' when name is empty", capturedCmd)
		}
	})
}

// --- GetCheckpoints ---

func TestGetCheckpoints(t *testing.T) {
	t.Run("VM exists calls outputPS for Get-VMCheckpoint and returns result", func(t *testing.T) {
		callCount := 0
		const checkpointOutput = "Name       CreationTime\n----       ------------\nsnap1      2024-01-01\n"
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				// IsVMExist
				return vmExistsOutput(), nil
			}
			// Get-VMCheckpoint
			return []byte(checkpointOutput), nil
		})
		result, err := GetCheckpoints("my-vm")
		if err != nil {
			t.Fatalf("GetCheckpoints: expected nil, got %v", err)
		}
		if !strings.Contains(result, "snap1") {
			t.Errorf("GetCheckpoints: result %q does not contain 'snap1'", result)
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmNotFoundErr()
		})
		_, err := GetCheckpoints("ghost-vm")
		if err == nil {
			t.Fatal("GetCheckpoints: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("GetCheckpoints: error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("PS error returns error", func(t *testing.T) {
		callCount := 0
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return vmExistsOutput(), nil
			}
			return nil, errors.New("ps error")
		})
		_, err := GetCheckpoints("my-vm")
		if err == nil {
			t.Fatal("GetCheckpoints: expected error on PS failure, got nil")
		}
		if !strings.Contains(err.Error(), "failed to get checkpoints") {
			t.Errorf("GetCheckpoints: error %q does not contain 'failed to get checkpoints'", err.Error())
		}
	})
}

// --- RestoreCheckpoint ---

func TestRestoreCheckpoint(t *testing.T) {
	t.Run("VM exists calls runPS and returns nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmExistsOutput(), nil
			},
		)
		if err := RestoreCheckpoint("my-vm", "snap1"); err != nil {
			t.Fatalf("RestoreCheckpoint: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("RestoreCheckpoint: runPS was not called")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmNotFoundErr()
		})
		err := RestoreCheckpoint("ghost-vm", "snap1")
		if err == nil {
			t.Fatal("RestoreCheckpoint: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("RestoreCheckpoint: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

// --- RemoveCheckpoint ---

func TestRemoveCheckpoint(t *testing.T) {
	t.Run("VM exists calls runPS and returns nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmExistsOutput(), nil
			},
		)
		if err := RemoveCheckpoint("my-vm", "snap1"); err != nil {
			t.Fatalf("RemoveCheckpoint: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("RemoveCheckpoint: runPS was not called")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmNotFoundErr()
		})
		err := RemoveCheckpoint("ghost-vm", "snap1")
		if err == nil {
			t.Fatal("RemoveCheckpoint: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("RemoveCheckpoint: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

// --- RenameCheckpoint ---

func TestRenameCheckpoint(t *testing.T) {
	t.Run("VM exists calls runPS and returns nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmExistsOutput(), nil
			},
		)
		if err := RenameCheckpoint("my-vm", "old-snap", "new-snap"); err != nil {
			t.Fatalf("RenameCheckpoint: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("RenameCheckpoint: runPS was not called")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmNotFoundErr()
		})
		err := RenameCheckpoint("ghost-vm", "old-snap", "new-snap")
		if err == nil {
			t.Fatal("RenameCheckpoint: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("RenameCheckpoint: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

// --- ExportCheckpoint ---

func TestExportCheckpoint(t *testing.T) {
	t.Run("VM exists calls runPS and returns nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return vmExistsOutput(), nil
			},
		)
		if err := ExportCheckpoint("my-vm", "snap1", "C:\\exports"); err != nil {
			t.Fatalf("ExportCheckpoint: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("ExportCheckpoint: runPS was not called")
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmNotFoundErr()
		})
		err := ExportCheckpoint("ghost-vm", "snap1", "C:\\exports")
		if err == nil {
			t.Fatal("ExportCheckpoint: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("ExportCheckpoint: error %q does not contain 'does not exist'", err.Error())
		}
	})
}
