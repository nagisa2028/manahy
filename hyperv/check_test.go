package hyperv

import (
	"errors"
	"strings"
	"testing"
)

// cmdletsFromBatchCmd extracts single-quoted cmdlet names from a batch Get-Command
// call and returns them newline-separated, simulating PowerShell finding all of them.
func cmdletsFromBatchCmd(cmd string) []byte {
	parts := strings.Split(cmd, "'")
	var sb strings.Builder
	for i := 1; i < len(parts); i += 2 {
		if name := parts[i]; name != "" {
			sb.WriteString(name)
			sb.WriteByte('\n')
		}
	}
	return []byte(sb.String())
}

// cmdletsFromBatchCmdExcept is like cmdletsFromBatchCmd but omits one name,
// simulating PowerShell not finding that particular cmdlet.
func cmdletsFromBatchCmdExcept(cmd, exclude string) []byte {
	parts := strings.Split(cmd, "'")
	var sb strings.Builder
	for i := 1; i < len(parts); i += 2 {
		if name := parts[i]; name != "" && name != exclude {
			sb.WriteString(name)
			sb.WriteByte('\n')
		}
	}
	return []byte(sb.String())
}

// --- checkPSAvailable ---

func TestCheckPSAvailable(t *testing.T) {
	t.Run("PS available returns OK with version", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("5.1.19041.0\n"), nil
		})
		r := checkPSAvailable()
		if r.Status != CheckOK {
			t.Errorf("checkPSAvailable: expected CheckOK, got %v", r.Status)
		}
		if !strings.Contains(r.Message, "5.1") {
			t.Errorf("checkPSAvailable: expected version in message, got %q", r.Message)
		}
	})

	t.Run("PS unavailable returns Fail", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("executable not found")
		})
		r := checkPSAvailable()
		if r.Status != CheckFail {
			t.Errorf("checkPSAvailable: expected CheckFail, got %v", r.Status)
		}
	})
}

// --- checkHyperVEnabled ---

func TestCheckHyperVEnabled(t *testing.T) {
	t.Run("Hyper-V enabled returns OK", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte(""), nil
		})
		r := checkHyperVEnabled()
		if r.Status != CheckOK {
			t.Errorf("checkHyperVEnabled: expected CheckOK, got %v", r.Status)
		}
	})

	t.Run("Hyper-V not enabled returns Fail", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		r := checkHyperVEnabled()
		if r.Status != CheckFail {
			t.Errorf("checkHyperVEnabled: expected CheckFail, got %v", r.Status)
		}
	})
}

// --- checkHyperVPermission ---

func TestCheckHyperVPermission(t *testing.T) {
	t.Run("member returns OK", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("True\n"), nil
		})
		r := checkHyperVPermission()
		if r.Status != CheckOK {
			t.Errorf("checkHyperVPermission: expected CheckOK, got %v", r.Status)
		}
	})

	t.Run("not member returns Fail", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil
		})
		r := checkHyperVPermission()
		if r.Status != CheckFail {
			t.Errorf("checkHyperVPermission: expected CheckFail, got %v", r.Status)
		}
	})

	t.Run("PS error returns Warn", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("access denied")
		})
		r := checkHyperVPermission()
		if r.Status != CheckWarn {
			t.Errorf("checkHyperVPermission: expected CheckWarn, got %v", r.Status)
		}
	})
}

// --- checkCmdletExists ---

func TestCheckCmdletExists(t *testing.T) {
	t.Run("cmdlet found returns OK", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte(""), nil
		})
		r := checkCmdletExists("Get-VM")
		if r.Status != CheckOK {
			t.Errorf("checkCmdletExists: expected CheckOK, got %v", r.Status)
		}
		if r.Name != "Get-VM" {
			t.Errorf("checkCmdletExists: expected Name=Get-VM, got %q", r.Name)
		}
	})

	t.Run("cmdlet not found returns Fail", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		r := checkCmdletExists("Get-VM")
		if r.Status != CheckFail {
			t.Errorf("checkCmdletExists: expected CheckFail, got %v", r.Status)
		}
	})
}

// --- CheckSystem ---

func TestCheckSystem(t *testing.T) {
	withPS(t, nil, func(_ string) ([]byte, error) {
		return []byte("True\n"), nil
	})
	results := CheckSystem()
	if len(results) != 3 {
		t.Fatalf("CheckSystem: expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Name == "" {
			t.Error("CheckSystem: result has empty Name")
		}
	}
}

// --- CheckVMCommands / CheckDiskCommands / CheckNetworkCommands ---

func TestCheckCommandGroups(t *testing.T) {
	t.Run("all VM cmdlets available", func(t *testing.T) {
		// Return all cmdlet names from the batch command as found.
		withPS(t, nil, func(cmd string) ([]byte, error) {
			return cmdletsFromBatchCmd(cmd), nil
		})
		for _, r := range CheckVMCommands() {
			if r.Status != CheckOK {
				t.Errorf("CheckVMCommands: %q expected OK, got Fail", r.Name)
			}
		}
	})

	t.Run("one disk cmdlet missing", func(t *testing.T) {
		// Return all disk cmdlets except New-VHD from the batch command.
		withPS(t, nil, func(cmd string) ([]byte, error) {
			return cmdletsFromBatchCmdExcept(cmd, "New-VHD"), nil
		})
		results := CheckDiskCommands()
		found := false
		for _, r := range results {
			if r.Name == "New-VHD" && r.Status == CheckFail {
				found = true
			}
		}
		if !found {
			t.Error("CheckDiskCommands: expected New-VHD to be CheckFail")
		}
		// All other disk cmdlets should be OK.
		for _, r := range results {
			if r.Name != "New-VHD" && r.Status != CheckOK {
				t.Errorf("CheckDiskCommands: %q expected OK, got Fail", r.Name)
			}
		}
	})

	t.Run("all network cmdlets available", func(t *testing.T) {
		withPS(t, nil, func(cmd string) ([]byte, error) {
			return cmdletsFromBatchCmd(cmd), nil
		})
		for _, r := range CheckNetworkCommands() {
			if r.Status != CheckOK {
				t.Errorf("CheckNetworkCommands: %q expected OK, got Fail", r.Name)
			}
		}
	})
}

// TestCheckCmdletsBatch verifies the single-call batch implementation of checkCmdlets.
func TestCheckCmdletsBatch(t *testing.T) {
	t.Run("result slice length matches input", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte(""), nil
		})
		cmdlets := []string{"Get-VM", "New-VM", "Remove-VM"}
		results := checkCmdlets(cmdlets)
		if len(results) != len(cmdlets) {
			t.Errorf("checkCmdlets: expected %d results, got %d", len(cmdlets), len(results))
		}
		for i, r := range results {
			if r.Name != cmdlets[i] {
				t.Errorf("checkCmdlets: results[%d].Name = %q, want %q", i, r.Name, cmdlets[i])
			}
		}
	})

	t.Run("empty input returns nil", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte(""), nil
		})
		if results := checkCmdlets(nil); results != nil {
			t.Errorf("checkCmdlets(nil): expected nil, got %v", results)
		}
	})

	t.Run("order preserved in results", func(t *testing.T) {
		// New-VM is absent; Get-VM and Remove-VM are present.
		withPS(t, nil, func(cmd string) ([]byte, error) {
			return cmdletsFromBatchCmdExcept(cmd, "New-VM"), nil
		})
		cmdlets := []string{"Get-VM", "New-VM", "Remove-VM"}
		results := checkCmdlets(cmdlets)
		if results[0].Status != CheckOK {
			t.Errorf("checkCmdlets: results[0] (Get-VM) expected OK, got %v", results[0].Status)
		}
		if results[1].Status != CheckFail {
			t.Errorf("checkCmdlets: results[1] (New-VM) expected Fail, got %v", results[1].Status)
		}
		if results[2].Status != CheckOK {
			t.Errorf("checkCmdlets: results[2] (Remove-VM) expected OK, got %v", results[2].Status)
		}
	})

	t.Run("all cmdlets missing when PS call fails", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("PS unavailable")
		})
		cmdlets := []string{"Get-VM", "New-VM"}
		results := checkCmdlets(cmdlets)
		for _, r := range results {
			if r.Status != CheckFail {
				t.Errorf("checkCmdlets PS error: %q expected CheckFail, got %v", r.Name, r.Status)
			}
		}
	})

	t.Run("repeated calls produce consistent results", func(t *testing.T) {
		withPS(t, nil, func(cmd string) ([]byte, error) {
			return cmdletsFromBatchCmd(cmd), nil
		})
		for i := 0; i < 20; i++ {
			results := CheckVMCommands()
			if len(results) == 0 {
				t.Fatalf("iteration %d: expected results, got none", i)
			}
			for _, r := range results {
				if r.Status != CheckOK {
					t.Errorf("iteration %d: %q expected CheckOK, got %v", i, r.Name, r.Status)
				}
			}
		}
	})
}
