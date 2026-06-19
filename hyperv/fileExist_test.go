package hyperv

import (
	"errors"
	"strings"
	"testing"
)

func TestSearchFilePath(t *testing.T) {
	tests := []struct {
		name       string
		psOutput   []byte
		psErr      error
		wantResult bool
		wantErr    bool
	}{
		{
			name:       "file exists",
			psOutput:   []byte("True\n"),
			wantResult: true,
		},
		{
			name:       "file does not exist",
			psOutput:   []byte("False\n"),
			wantResult: false,
		},
		{
			name:    "PS error",
			psErr:   errors.New("exec error"),
			wantErr: true,
		},
		{
			name:     "garbage output",
			psOutput: []byte("garbage"),
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			withPS(t, nil, func(_ string) ([]byte, error) {
				return tc.psOutput, tc.psErr
			})

			got, err := searchFilePath("C:\\some\\path")
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantResult {
				t.Errorf("searchFilePath: got %v; want %v", got, tc.wantResult)
			}
		})
	}
}

func TestIsFileExist(t *testing.T) {
	t.Run("file exists returns nil", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("True\n"), nil
		})
		if err := isFileExist("C:\\file.vhd"); err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
	})

	t.Run("file not exists returns error containing path", func(t *testing.T) {
		const path = "C:\\missing.vhd"
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil
		})
		err := isFileExist(path)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), path) {
			t.Errorf("error %q does not contain path %q", err.Error(), path)
		}
	})

	t.Run("PS error is propagated", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("exec error")
		})
		if err := isFileExist("C:\\file.vhd"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestIsNotFileExist(t *testing.T) {
	t.Run("file not exists returns nil", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("False\n"), nil
		})
		if err := isNotFileExist("C:\\file.vhd"); err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
	})

	t.Run("file exists returns error containing path", func(t *testing.T) {
		const path = "C:\\existing.vhd"
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte("True\n"), nil
		})
		err := isNotFileExist(path)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), path) {
			t.Errorf("error %q does not contain path %q", err.Error(), path)
		}
	})

	t.Run("PS error is propagated", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("exec error")
		})
		if err := isNotFileExist("C:\\file.vhd"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
