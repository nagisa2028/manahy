package hyperv

import (
	"errors"
	"strings"
	"testing"
)

// --- Pure validation tests (no PS) ---

func TestCheckSwitchType(t *testing.T) {
	tests := []struct {
		name       string
		switchType string
		wantErr    bool
	}{
		{"external lowercase", "external", false},
		{"internal lowercase", "internal", false},
		{"private lowercase", "private", false},
		{"External uppercase", "External", true},
		{"empty string", "", true},
		{"other value", "other", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkSwitchType(tc.switchType)
			if tc.wantErr && err == nil {
				t.Errorf("checkSwitchType(%q): expected error, got nil", tc.switchType)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("checkSwitchType(%q): expected nil, got %v", tc.switchType, err)
			}
		})
	}
}

func TestCheckSwitchIntegrity(t *testing.T) {
	tests := []struct {
		name              string
		switchType        string
		externalInterface string
		wantErr           bool
		errContains       string
	}{
		{"external with iface", "external", "eth0", false, ""},
		{"external without iface", "external", "", true, "required"},
		{"internal without iface", "internal", "", false, ""},
		{"internal with iface", "internal", "eth0", true, "only valid"},
		{"private without iface", "private", "", false, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkSwitchIntegrity(tc.switchType, tc.externalInterface)
			if tc.wantErr {
				if err == nil {
					t.Errorf("checkSwitchIntegrity(%q, %q): expected error, got nil", tc.switchType, tc.externalInterface)
					return
				}
				if tc.errContains != "" && !strings.Contains(err.Error(), tc.errContains) {
					t.Errorf("checkSwitchIntegrity(%q, %q): error %q does not contain %q", tc.switchType, tc.externalInterface, err.Error(), tc.errContains)
				}
			} else if err != nil {
				t.Errorf("checkSwitchIntegrity(%q, %q): expected nil, got %v", tc.switchType, tc.externalInterface, err)
			}
		})
	}
}

// --- PS-mocked tests ---

func TestGetSwitchType(t *testing.T) {
	t.Run("returns External on valid output", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("SwitchType\n----------\nExternal\n"), nil
		})
		got := GetSwitchType("mySwitch")
		if got != "External" {
			t.Errorf("GetSwitchType: got %q, want %q", got, "External")
		}
	})

	t.Run("returns NotFound on PS error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("ps error")
		})
		got := GetSwitchType("missing")
		if got != "NotFound" {
			t.Errorf("GetSwitchType: got %q, want %q", got, "NotFound")
		}
	})

	t.Run("returns Unknown on empty output", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("SwitchType\n----------\n"), nil
		})
		got := GetSwitchType("ambiguous")
		if got != "Unknown" {
			t.Errorf("GetSwitchType: got %q, want %q", got, "Unknown")
		}
	})
}

func TestIsSwitchExist(t *testing.T) {
	t.Run("exists returns nil", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("SwitchType\n----------\nExternal\n"), nil
		})
		if err := IsSwitchExist("mySwitch"); err != nil {
			t.Errorf("IsSwitchExist: expected nil, got %v", err)
		}
	})

	t.Run("NotFound returns error does not exist", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		err := IsSwitchExist("missing")
		if err == nil {
			t.Fatal("IsSwitchExist: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("IsSwitchExist: error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("Unknown returns error failed to get", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("SwitchType\n----------\n"), nil
		})
		err := IsSwitchExist("ambiguous")
		if err == nil {
			t.Fatal("IsSwitchExist: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to get") {
			t.Errorf("IsSwitchExist: error %q does not contain 'failed to get'", err.Error())
		}
	})
}

func TestIsNotSwitchExist(t *testing.T) {
	t.Run("NotFound returns nil", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		if err := IsNotSwitchExist("missing"); err != nil {
			t.Errorf("IsNotSwitchExist: expected nil, got %v", err)
		}
	})

	t.Run("exists returns error already exists", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("SwitchType\n----------\nInternal\n"), nil
		})
		err := IsNotSwitchExist("existing")
		if err == nil {
			t.Fatal("IsNotSwitchExist: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("IsNotSwitchExist: error %q does not contain 'already exists'", err.Error())
		}
	})

	t.Run("Unknown returns error failed to get", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("SwitchType\n----------\n"), nil
		})
		err := IsNotSwitchExist("ambiguous")
		if err == nil {
			t.Fatal("IsNotSwitchExist: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to get") {
			t.Errorf("IsNotSwitchExist: error %q does not contain 'failed to get'", err.Error())
		}
	})
}

func TestRemoveSwitch(t *testing.T) {
	t.Run("exists calls runPS and returns nil", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return []byte("SwitchType\n----------\nExternal\n"), nil
			},
		)
		if err := RemoveSwitch("mySwitch"); err != nil {
			t.Errorf("RemoveSwitch: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("RemoveSwitch: runPS was not called")
		}
	})

	t.Run("not found returns error from IsSwitchExist", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		err := RemoveSwitch("missing")
		if err == nil {
			t.Fatal("RemoveSwitch: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("RemoveSwitch: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestRenameSwitch(t *testing.T) {
	t.Run("source exists and new name free calls runPS", func(t *testing.T) {
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
					// IsSwitchExist for source: exists
					return []byte("SwitchType\n----------\nExternal\n"), nil
				}
				// IsNotSwitchExist for newName: not found
				return nil, errors.New("not found")
			},
		)
		if err := RenameSwitch("oldName", "newName"); err != nil {
			t.Errorf("RenameSwitch: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("RenameSwitch: runPS was not called")
		}
	})

	t.Run("source not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		err := RenameSwitch("missing", "newName")
		if err == nil {
			t.Fatal("RenameSwitch: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("RenameSwitch: error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("dest already exists returns error", func(t *testing.T) {
		callCount := 0
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				// IsSwitchExist for source: exists
				return []byte("SwitchType\n----------\nExternal\n"), nil
			}
			// IsNotSwitchExist for newName: already exists
			return []byte("SwitchType\n----------\nInternal\n"), nil
		})
		err := RenameSwitch("source", "existing")
		if err == nil {
			t.Fatal("RenameSwitch: expected error for existing dest, got nil")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("RenameSwitch: error %q does not contain 'already exists'", err.Error())
		}
	})
}

func TestChangeSwitchType(t *testing.T) {
	t.Run("exists and valid new type different calls runPS", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return []byte("SwitchType\n----------\nInternal\n"), nil
			},
		)
		if err := ChangeSwitchType("mySwitch", "private"); err != nil {
			t.Errorf("ChangeSwitchType: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("ChangeSwitchType: runPS was not called")
		}
	})

	t.Run("Unknown state returns error failed to get", func(t *testing.T) {
		// Empty output causes GetSwitchType to return "Unknown".
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("SwitchType\n----------\n"), nil
		})
		err := ChangeSwitchType("mySwitch", "internal")
		if err == nil {
			t.Fatal("ChangeSwitchType: expected error for Unknown state, got nil")
		}
		if !strings.Contains(err.Error(), "failed to get") {
			t.Errorf("ChangeSwitchType: error %q does not contain 'failed to get'", err.Error())
		}
	})

	t.Run("not found returns error does not exist", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("not found")
		})
		err := ChangeSwitchType("missing", "internal")
		if err == nil {
			t.Fatal("ChangeSwitchType: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("ChangeSwitchType: error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("invalid type returns error before PS", func(t *testing.T) {
		outputCalled := 0
		withPS(t,
			func(_ string) error {
				return nil
			},
			func(_ string) ([]byte, error) {
				outputCalled++
				// First call is GetSwitchType to check if switch exists
				return []byte("SwitchType\n----------\nInternal\n"), nil
			},
		)
		err := ChangeSwitchType("mySwitch", "InvalidType")
		if err == nil {
			t.Fatal("ChangeSwitchType: expected error for invalid type, got nil")
		}
		if !strings.Contains(err.Error(), "undefined switch type") {
			t.Errorf("ChangeSwitchType: error %q does not contain 'undefined switch type'", err.Error())
		}
	})
}
