package hyperv

import (
	"errors"
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
		{"64TB at limit", "64TB", false},
		{"65TB exceeds limit", "65TB", true},
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

// ---------- checkDiskParam ----------

func TestCheckDiskParam(t *testing.T) {
	t.Run("relative disk path returns error", func(t *testing.T) {
		err := checkDiskParam(Disk{Path: `router1.vhd`, Type: "dynamic", Size: "10GB"})
		if err == nil {
			t.Fatal("checkDiskParam: expected error for relative path, got nil")
		}
		if !strings.Contains(err.Error(), "absolute") {
			t.Errorf("checkDiskParam: error %q does not contain 'absolute'", err.Error())
		}
	})

	t.Run("import disk path not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil // Test-Path: import file missing
		})
		err := checkDiskParam(Disk{Path: `C:\missing.vhd`, Import: true})
		if err == nil {
			t.Fatal("checkDiskParam: expected error for missing import disk, got nil")
		}
	})

	t.Run("non-import disk already exists returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("True\n"), nil // Test-Path: file already exists
		})
		err := checkDiskParam(Disk{Path: `C:\existing.vhd`, Type: "dynamic", Size: "10GB"})
		if err == nil {
			t.Fatal("checkDiskParam: expected error for existing non-import disk, got nil")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("checkDiskParam: error %q does not contain 'already exists'", err.Error())
		}
	})

	t.Run("invalid disk type returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil // isNotFileExist: path free
		})
		err := checkDiskParam(Disk{Path: `C:\new.vhd`, Type: "invalid", Size: "10GB"})
		if err == nil {
			t.Fatal("checkDiskParam: expected error for invalid disk type, got nil")
		}
		if !strings.Contains(err.Error(), "invalid disk type") {
			t.Errorf("checkDiskParam: error %q does not contain 'invalid disk type'", err.Error())
		}
	})

	t.Run("differencing disk with relative parent path returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil // isNotFileExist: child path free
		})
		disk := Disk{Path: `C:\VMs\child.vhd`, Type: "differencing", ParentPath: `parent.vhd`}
		err := checkDiskParam(disk)
		if err == nil {
			t.Fatal("checkDiskParam: expected error for relative parent path, got nil")
		}
		if !strings.Contains(err.Error(), "absolute") {
			t.Errorf("checkDiskParam: error %q does not contain 'absolute'", err.Error())
		}
	})

	t.Run("differencing disk with missing parent returns error", func(t *testing.T) {
		callCount := 0
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return []byte("False\n"), nil // isNotFileExist: child path free
			}
			return []byte("False\n"), nil // isFileExist: parent missing
		})
		disk := Disk{Path: `C:\child.vhd`, Type: "differencing", ParentPath: `C:\missing-parent.vhd`}
		err := checkDiskParam(disk)
		if err == nil {
			t.Fatal("checkDiskParam: expected error for missing parent, got nil")
		}
	})

	t.Run("fixed disk with negative SourceDisk returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil // isNotFileExist: path free
		})
		disk := Disk{Path: `C:\new.vhd`, Type: "fixed", Size: "10GB", SourceDisk: -1}
		err := checkDiskParam(disk)
		if err == nil {
			t.Fatal("checkDiskParam: expected error for negative SourceDisk, got nil")
		}
		if !strings.Contains(err.Error(), "non-negative") {
			t.Errorf("checkDiskParam: error %q does not contain 'non-negative'", err.Error())
		}
	})

	t.Run("valid dynamic disk returns nil", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil // isNotFileExist: path free
		})
		disk := Disk{Path: `C:\new.vhd`, Type: "dynamic", Size: "10GB"}
		if err := checkDiskParam(disk); err != nil {
			t.Fatalf("checkDiskParam: expected nil, got %v", err)
		}
	})
}

// ---------- CreateDisk ----------

func TestCreateDisk(t *testing.T) {
	t.Run("import disk returns nil immediately without runPS", func(t *testing.T) {
		runCalled := false
		withPS(t, func(_ string) error { runCalled = true; return nil }, nil)
		if err := CreateDisk(Disk{Import: true}, false); err != nil {
			t.Fatalf("CreateDisk: expected nil for import, got %v", err)
		}
		if runCalled {
			t.Error("CreateDisk: runPS called for import disk")
		}
	})

	t.Run("dynamic disk command contains path and size", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error { capturedCmd = c; return nil },
			func(_ string) ([]byte, error) { return []byte("False\n"), nil },
		)
		disk := Disk{Path: `C:\new.vhd`, Type: "dynamic", Size: "20GB"}
		if err := CreateDisk(disk, false); err != nil {
			t.Fatalf("CreateDisk: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, `C:\new.vhd`) {
			t.Errorf("CreateDisk: command %q does not contain path", capturedCmd)
		}
		if !strings.Contains(capturedCmd, "20GB") {
			t.Errorf("CreateDisk: command %q does not contain '20GB'", capturedCmd)
		}
	})

	t.Run("differencing disk command contains -Differencing and ParentPath", func(t *testing.T) {
		var capturedCmd string
		callCount := 0
		withPS(t,
			func(c string) error { capturedCmd = c; return nil },
			func(_ string) ([]byte, error) {
				callCount++
				if callCount == 1 {
					return []byte("False\n"), nil // isNotFileExist: child path free
				}
				return []byte("True\n"), nil // isFileExist: parent exists
			},
		)
		disk := Disk{Path: `C:\child.vhd`, Type: "differencing", ParentPath: `C:\parent.vhd`}
		if err := CreateDisk(disk, false); err != nil {
			t.Fatalf("CreateDisk: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "-Differencing") {
			t.Errorf("CreateDisk: command %q does not contain '-Differencing'", capturedCmd)
		}
		if !strings.Contains(capturedCmd, `C:\parent.vhd`) {
			t.Errorf("CreateDisk: command %q does not contain parent path", capturedCmd)
		}
	})

	t.Run("runPS error is returned", func(t *testing.T) {
		withPS(t,
			func(_ string) error { return errors.New("ps error") },
			func(_ string) ([]byte, error) { return []byte("False\n"), nil },
		)
		disk := Disk{Path: `C:\new.vhd`, Type: "dynamic", Size: "10GB"}
		if err := CreateDisk(disk, false); err == nil {
			t.Fatal("CreateDisk: expected error from runPS, got nil")
		}
	})
}


// ---------- TestRemoveDisk ----------

func TestRemoveDisk(t *testing.T) {
	t.Run("file exists calls runPS and returns nil", func(t *testing.T) {
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
		if err := RemoveDisk(`C:\disk.vhd`, false); err != nil {
			t.Fatalf("RemoveDisk: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("RemoveDisk: runPS was not called")
		}
	})

	t.Run("file not found returns error before runPS", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return []byte("False\n"), nil
			},
		)
		err := RemoveDisk(`C:\missing.vhd`, false)
		if err == nil {
			t.Fatal("RemoveDisk: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("RemoveDisk: error %q does not contain 'does not exist'", err.Error())
		}
		if runCalled {
			t.Error("RemoveDisk: runPS should not be called when file is missing")
		}
	})

	t.Run("runPS failure returns error", func(t *testing.T) {
		withPS(t,
			func(_ string) error {
				return errors.New("access denied")
			},
			func(_ string) ([]byte, error) {
				return []byte("True\n"), nil
			},
		)
		err := RemoveDisk(`C:\disk.vhd`, false)
		if err == nil {
			t.Fatal("RemoveDisk: expected error from runPS, got nil")
		}
	})
}

// ---------- runPS failure paths for VHD operations ----------

func TestResizeVHDRunPSFailure(t *testing.T) {
	withPS(t,
		func(_ string) error {
			return errors.New("resize failed")
		},
		func(_ string) ([]byte, error) {
			return []byte("True\n"), nil
		},
	)
	err := ResizeVHD(`C:\disk.vhd`, "50GB")
	if err == nil {
		t.Fatal("ResizeVHD: expected error from runPS, got nil")
	}
}

func TestOptimizeVHDRunPSFailure(t *testing.T) {
	withPS(t,
		func(_ string) error {
			return errors.New("optimize failed")
		},
		func(_ string) ([]byte, error) {
			return []byte("True\n"), nil
		},
	)
	err := OptimizeVHD(`C:\disk.vhd`)
	if err == nil {
		t.Fatal("OptimizeVHD: expected error from runPS, got nil")
	}
}

func TestMountVHDRunPSFailure(t *testing.T) {
	withPS(t,
		func(_ string) error {
			return errors.New("mount failed")
		},
		func(_ string) ([]byte, error) {
			return []byte("True\n"), nil
		},
	)
	err := MountVHD(`C:\disk.vhd`)
	if err == nil {
		t.Fatal("MountVHD: expected error from runPS, got nil")
	}
}

func TestDismountVHDRunPSFailure(t *testing.T) {
	withPS(t,
		func(_ string) error {
			return errors.New("dismount failed")
		},
		func(_ string) ([]byte, error) {
			return []byte("True\n"), nil
		},
	)
	err := DismountVHD(`C:\disk.vhd`)
	if err == nil {
		t.Fatal("DismountVHD: expected error from runPS, got nil")
	}
}

func TestMergeVHDRunPSFailure(t *testing.T) {
	withPS(t,
		func(_ string) error {
			return errors.New("merge failed")
		},
		func(_ string) ([]byte, error) {
			return []byte("True\n"), nil
		},
	)
	err := MergeVHD(`C:\child.vhd`, "")
	if err == nil {
		t.Fatal("MergeVHD: expected error from runPS, got nil")
	}
}

func TestConvertVHDRunPSFailure(t *testing.T) {
	withPS(t,
		func(_ string) error {
			return errors.New("convert failed")
		},
		func(_ string) ([]byte, error) {
			return []byte("True\n"), nil
		},
	)
	err := ConvertVHD(`C:\disk.vhd`, `C:\disk2.vhdx`, "")
	if err == nil {
		t.Fatal("ConvertVHD: expected error from runPS, got nil")
	}
}

// ---------- checkDiskSize GB boundary ----------

func TestCheckDiskSizeGBBoundary(t *testing.T) {
	// maxDiskSizeGB = 64 * 1024 = 65536 GB
	t.Run("65536GB at limit is valid", func(t *testing.T) {
		if err := checkDiskSize("65536GB"); err != nil {
			t.Errorf("checkDiskSize(65536GB): expected nil, got %v", err)
		}
	})
	t.Run("65537GB exceeds limit is invalid", func(t *testing.T) {
		if err := checkDiskSize("65537GB"); err == nil {
			t.Error("checkDiskSize(65537GB): expected error, got nil")
		}
	})
}

// ---------- CloneDisk ----------

func TestCloneDisk(t *testing.T) {
	t.Run("source exists and dest absent calls runPS with SourcePath", func(t *testing.T) {
		var capturedCmd string
		callCount := 0
		withPS(t,
			func(c string) error {
				capturedCmd = c
				return nil
			},
			func(_ string) ([]byte, error) {
				callCount++
				if callCount == 1 {
					return []byte("True\n"), nil // source exists
				}
				return []byte("False\n"), nil // dest does not exist
			},
		)
		if err := CloneDisk(`C:\src.vhdx`, `C:\clone.vhdx`); err != nil {
			t.Fatalf("CloneDisk: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "SourcePath") {
			t.Errorf("CloneDisk: command %q does not contain 'SourcePath'", capturedCmd)
		}
		if !strings.Contains(capturedCmd, `C:\src.vhdx`) {
			t.Errorf("CloneDisk: command %q does not contain source path", capturedCmd)
		}
		if !strings.Contains(capturedCmd, `C:\clone.vhdx`) {
			t.Errorf("CloneDisk: command %q does not contain dest path", capturedCmd)
		}
	})

	t.Run("source missing returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil // source missing
		})
		err := CloneDisk(`C:\missing.vhdx`, `C:\clone.vhdx`)
		if err == nil {
			t.Fatal("CloneDisk: expected error for missing source, got nil")
		}
		if !strings.Contains(err.Error(), "source disk") {
			t.Errorf("CloneDisk: error %q does not mention 'source disk'", err.Error())
		}
	})

	t.Run("dest already exists returns error", func(t *testing.T) {
		callCount := 0
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return []byte("True\n"), nil // source exists
			}
			return []byte("True\n"), nil // dest already exists
		})
		err := CloneDisk(`C:\src.vhdx`, `C:\existing.vhdx`)
		if err == nil {
			t.Fatal("CloneDisk: expected error for existing dest, got nil")
		}
		if !strings.Contains(err.Error(), "destination disk") {
			t.Errorf("CloneDisk: error %q does not mention 'destination disk'", err.Error())
		}
	})
}
