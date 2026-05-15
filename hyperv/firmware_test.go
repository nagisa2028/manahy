package hyperv

import (
	"errors"
	"strings"
	"testing"
)

func TestSetVMSecureBoot(t *testing.T) {
	t.Run("enable Secure Boot includes On and no template", func(t *testing.T) {
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
		if err := SetVMSecureBoot("my-vm", true, ""); err != nil {
			t.Fatalf("SetVMSecureBoot: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "-EnableSecureBoot On") {
			t.Errorf("SetVMSecureBoot: command %q does not contain '-EnableSecureBoot On'", capturedCmd)
		}
		if strings.Contains(capturedCmd, "-SecureBootTemplate") {
			t.Errorf("SetVMSecureBoot: command %q should not contain '-SecureBootTemplate' when template is empty", capturedCmd)
		}
	})

	t.Run("disable Secure Boot includes Off", func(t *testing.T) {
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
		if err := SetVMSecureBoot("my-vm", false, ""); err != nil {
			t.Fatalf("SetVMSecureBoot: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "-EnableSecureBoot Off") {
			t.Errorf("SetVMSecureBoot: command %q does not contain '-EnableSecureBoot Off'", capturedCmd)
		}
	})

	t.Run("template name appended to command when provided", func(t *testing.T) {
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
		if err := SetVMSecureBoot("my-vm", true, "MicrosoftWindows"); err != nil {
			t.Fatalf("SetVMSecureBoot: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "-SecureBootTemplate") {
			t.Errorf("SetVMSecureBoot: command %q does not contain '-SecureBootTemplate'", capturedCmd)
		}
		if !strings.Contains(capturedCmd, "MicrosoftWindows") {
			t.Errorf("SetVMSecureBoot: command %q does not contain 'MicrosoftWindows'", capturedCmd)
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
		err := SetVMSecureBoot("my-vm", true, "")
		if err == nil {
			t.Fatal("SetVMSecureBoot: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to set Secure Boot") {
			t.Errorf("SetVMSecureBoot: error %q does not contain 'failed to set Secure Boot'", err.Error())
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmMissingErr()
		})
		err := SetVMSecureBoot("ghost-vm", true, "")
		if err == nil {
			t.Fatal("SetVMSecureBoot: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("SetVMSecureBoot: error %q does not contain 'does not exist'", err.Error())
		}
	})
}

func TestGetVMFirmware(t *testing.T) {
	t.Run("VM exists returns firmware info string", func(t *testing.T) {
		callCount := 0
		const firmwareOutput = "SecureBoot                   : On\nSecureBootTemplate           : MicrosoftWindows\nPreferredNetworkBootProtocol : IPv4\n"
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return vmRunningOutput(), nil
			}
			return []byte(firmwareOutput), nil
		})
		result, err := GetVMFirmware("my-vm")
		if err != nil {
			t.Fatalf("GetVMFirmware: expected nil, got %v", err)
		}
		if !strings.Contains(result, "SecureBoot") {
			t.Errorf("GetVMFirmware: result %q does not contain 'SecureBoot'", result)
		}
	})

	t.Run("VM not found returns error", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, vmMissingErr()
		})
		_, err := GetVMFirmware("ghost-vm")
		if err == nil {
			t.Fatal("GetVMFirmware: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("GetVMFirmware: error %q does not contain 'does not exist'", err.Error())
		}
	})

	t.Run("outputPS error returns wrapped error", func(t *testing.T) {
		callCount := 0
		withPS(t, nil, func(_ string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return vmRunningOutput(), nil
			}
			return nil, errors.New("ps error")
		})
		_, err := GetVMFirmware("my-vm")
		if err == nil {
			t.Fatal("GetVMFirmware: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to get firmware") {
			t.Errorf("GetVMFirmware: error %q does not contain 'failed to get firmware'", err.Error())
		}
	})
}
