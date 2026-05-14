package hyperv

import (
	"strings"
	"testing"
)

// --- Pure validation tests (no PS) ---

func TestCheckDiskType(t *testing.T) {
	tests := []struct {
		name     string
		diskType string
		wantErr  bool
	}{
		{"dynamic", "dynamic", false},
		{"fixed", "fixed", false},
		{"differencing", "differencing", false},
		{"raw", "raw", true},
		{"empty", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkDiskType(tc.diskType)
			if tc.wantErr && err == nil {
				t.Errorf("checkDiskType(%q): expected error, got nil", tc.diskType)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("checkDiskType(%q): expected nil, got %v", tc.diskType, err)
			}
		})
	}
}

func TestCheckDiskSize(t *testing.T) {
	tests := []struct {
		name     string
		diskSize string
		wantErr  bool
	}{
		{"10GB", "10GB", false},
		{"100MB", "100MB", false},
		{"1TB", "1TB", false},
		{"10KB invalid", "10KB", true},
		{"100 no unit", "100", true},
		{"empty", "", true},
		{"10gb lowercase", "10gb", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkDiskSize(tc.diskSize)
			if tc.wantErr && err == nil {
				t.Errorf("checkDiskSize(%q): expected error, got nil", tc.diskSize)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("checkDiskSize(%q): expected nil, got %v", tc.diskSize, err)
			}
		})
	}
}

// --- PS-mocked tests ---

func TestGetVHDInfo(t *testing.T) {
	t.Run("file exists calls outputPS for Get-VHD and returns info", func(t *testing.T) {
		callCount := 0
		const vhdOutput = "Path : C:\\disk.vhd\nVhdType : Dynamic\n"
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				// Test-Path (isFileExist)
				return []byte("True\n"), nil
			}
			// Get-VHD
			return []byte(vhdOutput), nil
		})
		info, err := GetVHDInfo("C:\\disk.vhd")
		if err != nil {
			t.Fatalf("GetVHDInfo: expected nil error, got %v", err)
		}
		if !strings.Contains(info, "Dynamic") {
			t.Errorf("GetVHDInfo: expected output to contain 'Dynamic', got %q", info)
		}
	})

	t.Run("file not found returns error before PS", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil
		})
		_, err := GetVHDInfo("C:\\missing.vhd")
		if err == nil {
			t.Fatal("GetVHDInfo: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("GetVHDInfo: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestResizeVHD(t *testing.T) {
	t.Run("file exists and valid size calls runPS", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return []byte("True\n"), nil
			},
		)
		if err := ResizeVHD("C:\\disk.vhd", "50GB"); err != nil {
			t.Errorf("ResizeVHD: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("ResizeVHD: runPS was not called")
		}
	})

	t.Run("invalid size returns error before PS", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return []byte("True\n"), nil
			},
		)
		err := ResizeVHD("C:\\disk.vhd", "50kb")
		if err == nil {
			t.Fatal("ResizeVHD: expected error for invalid size, got nil")
		}
		if runCalled {
			t.Error("ResizeVHD: runPS should not be called when size is invalid")
		}
	})

	t.Run("file not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil
		})
		err := ResizeVHD("C:\\missing.vhd", "50GB")
		if err == nil {
			t.Fatal("ResizeVHD: expected error for missing file, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("ResizeVHD: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestOptimizeVHD(t *testing.T) {
	t.Run("file exists calls runPS", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return []byte("True\n"), nil
			},
		)
		if err := OptimizeVHD("C:\\disk.vhd"); err != nil {
			t.Errorf("OptimizeVHD: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("OptimizeVHD: runPS was not called")
		}
	})

	t.Run("file not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil
		})
		err := OptimizeVHD("C:\\missing.vhd")
		if err == nil {
			t.Fatal("OptimizeVHD: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("OptimizeVHD: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestMountVHD(t *testing.T) {
	t.Run("file exists calls runPS", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return []byte("True\n"), nil
			},
		)
		if err := MountVHD("C:\\disk.vhd"); err != nil {
			t.Errorf("MountVHD: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("MountVHD: runPS was not called")
		}
	})

	t.Run("file not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil
		})
		err := MountVHD("C:\\missing.vhd")
		if err == nil {
			t.Fatal("MountVHD: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("MountVHD: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestDismountVHD(t *testing.T) {
	t.Run("file exists calls runPS", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return []byte("True\n"), nil
			},
		)
		if err := DismountVHD("C:\\disk.vhd"); err != nil {
			t.Errorf("DismountVHD: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("DismountVHD: runPS was not called")
		}
	})

	t.Run("file not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil
		})
		err := DismountVHD("C:\\missing.vhd")
		if err == nil {
			t.Fatal("DismountVHD: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("DismountVHD: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestMergeVHD(t *testing.T) {
	t.Run("file exists calls runPS", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return []byte("True\n"), nil
			},
		)
		if err := MergeVHD("C:\\child.vhd", ""); err != nil {
			t.Errorf("MergeVHD: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("MergeVHD: runPS was not called")
		}
	})

	t.Run("file not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil
		})
		err := MergeVHD("C:\\missing.vhd", "")
		if err == nil {
			t.Fatal("MergeVHD: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("MergeVHD: error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("with destPath includes DestinationPath in command", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error {
				capturedCmd = c
				return nil
			},
			func(_ string) ([]byte, error) {
				return []byte("True\n"), nil
			},
		)
		if err := MergeVHD("C:\\child.vhd", "C:\\parent.vhd"); err != nil {
			t.Errorf("MergeVHD: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "DestinationPath") {
			t.Errorf("MergeVHD: command %q does not contain 'DestinationPath'", capturedCmd)
		}
	})
}

func TestConvertVHD(t *testing.T) {
	t.Run("file exists calls runPS", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return []byte("True\n"), nil
			},
		)
		if err := ConvertVHD("C:\\disk.vhd", "C:\\disk2.vhdx", ""); err != nil {
			t.Errorf("ConvertVHD: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("ConvertVHD: runPS was not called")
		}
	})

	t.Run("with diskType includes -VHDType in command", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error {
				capturedCmd = c
				return nil
			},
			func(_ string) ([]byte, error) {
				return []byte("True\n"), nil
			},
		)
		if err := ConvertVHD("C:\\disk.vhd", "C:\\disk2.vhdx", "Fixed"); err != nil {
			t.Errorf("ConvertVHD: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "-VHDType") {
			t.Errorf("ConvertVHD: command %q does not contain '-VHDType'", capturedCmd)
		}
	})

	t.Run("file not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil
		})
		err := ConvertVHD("C:\\missing.vhd", "C:\\dest.vhdx", "")
		if err == nil {
			t.Fatal("ConvertVHD: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("ConvertVHD: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

