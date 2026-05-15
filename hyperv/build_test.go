package hyperv

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// --- numberPath ---

func TestNumberPath(t *testing.T) {
	tests := []struct {
		path string
		n    int
		want string
	}{
		{`C:\disks\test.vhdx`, 1, `C:\disks\test1.vhdx`},
		{`C:\disks\test.vhdx`, 3, `C:\disks\test3.vhdx`},
		{`C:\vms\boot.vhd`, 2, `C:\vms\boot2.vhd`},
		{`nodot`, 1, `nodot1`},
		{`C:\multi.part.vhdx`, 2, `C:\multi.part2.vhdx`},
	}
	for _, tc := range tests {
		got := numberPath(tc.path, tc.n)
		if got != tc.want {
			t.Errorf("numberPath(%q, %d) = %q, want %q", tc.path, tc.n, got, tc.want)
		}
	}
}

// --- resolveDisks ---

func TestResolveDisks(t *testing.T) {
	config := Summarize{
		Disks: map[string]Disk{
			"data-disk": {Path: `C:\VMs\data.vhd`},
			"boot-disk": {Path: `C:\VMs\boot.vhd`},
		},
	}

	t.Run("alias is resolved to path", func(t *testing.T) {
		got := resolveDisks(config, []string{"data-disk"})
		if len(got) != 1 || got[0] != `C:\VMs\data.vhd` {
			t.Errorf("resolveDisks: got %v, want [C:\\VMs\\data.vhd]", got)
		}
	})

	t.Run("unknown key is passed through as-is", func(t *testing.T) {
		got := resolveDisks(config, []string{`C:\literal\path.vhd`})
		if len(got) != 1 || got[0] != `C:\literal\path.vhd` {
			t.Errorf("resolveDisks: got %v, want literal path", got)
		}
	})

	t.Run("multiple refs including mix of alias and literal", func(t *testing.T) {
		got := resolveDisks(config, []string{"boot-disk", `C:\other.vhd`, "data-disk"})
		want := []string{`C:\VMs\boot.vhd`, `C:\other.vhd`, `C:\VMs\data.vhd`}
		for i, w := range want {
			if got[i] != w {
				t.Errorf("resolveDisks[%d]: got %q, want %q", i, got[i], w)
			}
		}
	})

	t.Run("empty refs returns empty slice", func(t *testing.T) {
		got := resolveDisks(config, []string{})
		if len(got) != 0 {
			t.Errorf("resolveDisks: expected empty, got %v", got)
		}
	})

	t.Run("empty config returns literal paths", func(t *testing.T) {
		got := resolveDisks(Summarize{}, []string{"some-disk"})
		if len(got) != 1 || got[0] != "some-disk" {
			t.Errorf("resolveDisks with empty config: got %v", got)
		}
	})
}

// --- BuildByStruct ---

// buildOutputMock routes outputPS responses by command content.
func buildOutputMock(cmd string) ([]byte, error) {
	switch {
	case strings.Contains(cmd, "Test-Path"):
		return []byte("False\n"), nil // file does not exist yet → OK to create
	case strings.Contains(cmd, "Get-VMSwitch"):
		return nil, errors.New("not found") // switch does not exist → create it
	case strings.Contains(cmd, "Get-VM"):
		return nil, errors.New("not found") // VM does not exist → create it
	default:
		return nil, errors.New("not found")
	}
}

func TestBuildByStruct(t *testing.T) {
	t.Run("import disk is skipped", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return nil, errors.New("not found")
			},
		)
		config := Summarize{
			Disks: map[string]Disk{
				"import-disk": {Path: `C:\existing.vhd`, Import: true},
			},
		}
		_ = BuildByStruct(config)
		if runCalled {
			t.Error("BuildByStruct: runPS was called for import disk")
		}
	})

	t.Run("non-import disk calls runPS", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			buildOutputMock,
		)
		config := Summarize{
			Disks: map[string]Disk{
				"new-disk": {Path: `C:\new.vhd`, Size: "10GB", Type: "dynamic"},
			},
		}
		if err := BuildByStruct(config); err != nil {
			t.Fatalf("BuildByStruct: expected nil, got %v", err)
		}
		if !runCalled {
			t.Error("BuildByStruct: runPS was not called for non-import disk")
		}
	})

	t.Run("VM count exceeds maximum returns error", func(t *testing.T) {
		withPS(t,
			func(_ string) error { return nil },
			func(_ string) ([]byte, error) { return nil, errors.New("not found") },
		)
		config := Summarize{
			Vms: map[string]VM{
				"cluster": {Count: maxVMCount + 1, Generation: 1, Path: `C:\VMs`, Memory: Memory{Size: "512MB"}, CPU: CPU{Thread: 1}},
			},
		}
		err := BuildByStruct(config)
		if err == nil {
			t.Fatal("BuildByStruct: expected error for count exceeding maximum, got nil")
		}
		if !strings.Contains(err.Error(), "exceeds maximum") {
			t.Errorf("BuildByStruct: error %q does not contain 'exceeds maximum'", err.Error())
		}
	})

	t.Run("disk alias is resolved before VM creation", func(t *testing.T) {
		var capturedCmds []string
		withPS(t,
			func(c string) error {
				capturedCmds = append(capturedCmds, c)
				return nil
			},
			buildOutputMock,
		)
		config := Summarize{
			Disks: map[string]Disk{
				"boot": {Path: `C:\VMs\boot.vhd`, Size: "10GB", Type: "dynamic"},
			},
			Vms: map[string]VM{
				"test-vm": {
					Generation: 1,
					Path:       `C:\VMs`,
					Memory:     Memory{Size: "512MB"},
					CPU:        CPU{Thread: 1},
					Disks:      []string{"boot"},
				},
			},
		}
		_ = BuildByStruct(config)
		found := false
		for _, cmd := range capturedCmds {
			if strings.Contains(cmd, `C:\VMs\boot.vhd`) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("BuildByStruct: resolved disk path not found in PS commands: %v", capturedCmds)
		}
	})

	t.Run("existing VM is silently skipped", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(cmd string) ([]byte, error) {
				if strings.Contains(cmd, "Get-VM") {
					return stateOutput("Running"), nil // VM already exists
				}
				return []byte("False\n"), nil
			},
		)
		config := Summarize{
			Vms: map[string]VM{
				"existing-vm": {Generation: 1, Path: `C:\VMs`, Memory: Memory{Size: "512MB"}, CPU: CPU{Thread: 1}},
			},
		}
		if err := BuildByStruct(config); err != nil {
			t.Fatalf("BuildByStruct: expected nil for existing VM, got %v", err)
		}
		if runCalled {
			t.Error("BuildByStruct: runPS was called for already-existing VM")
		}
	})

	t.Run("existing disk is silently skipped", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(cmd string) ([]byte, error) {
				if strings.Contains(cmd, "Test-Path") {
					return []byte("True\n"), nil // disk already exists
				}
				return nil, errors.New("not found")
			},
		)
		config := Summarize{
			Disks: map[string]Disk{
				"existing-disk": {Path: `C:\VMs\data.vhd`, Size: "10GB", Type: "dynamic"},
			},
		}
		if err := BuildByStruct(config); err != nil {
			t.Fatalf("BuildByStruct: expected nil for existing disk, got %v", err)
		}
		if runCalled {
			t.Error("BuildByStruct: runPS was called for already-existing disk")
		}
	})

}

// --- resolveAndCreateDisks ---

func TestResolveAndCreateDisks(t *testing.T) {
	config := Summarize{
		Disks: map[string]Disk{
			"boot":       {Path: `C:\VMs\boot.vhd`, Size: "10GB", Type: "dynamic"},
			"shared-iso": {Path: `C:\ISOs\install.iso`, Import: true},
		},
	}

	t.Run("count=1 returns base path without creating a new disk", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			buildOutputMock,
		)
		paths, err := resolveAndCreateDisks(config, []string{"boot"}, 1, 1)
		if err != nil {
			t.Fatalf("resolveAndCreateDisks count=1: expected nil, got %v", err)
		}
		if len(paths) != 1 || paths[0] != `C:\VMs\boot.vhd` {
			t.Errorf("resolveAndCreateDisks count=1: got %v, want [C:\\VMs\\boot.vhd]", paths)
		}
		if runCalled {
			t.Error("resolveAndCreateDisks count=1: runPS should not be called")
		}
	})

	t.Run("count=3 index=2 creates boot2.vhd and returns numbered path", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error {
				capturedCmd = c
				return nil
			},
			buildOutputMock,
		)
		paths, err := resolveAndCreateDisks(config, []string{"boot"}, 2, 3)
		if err != nil {
			t.Fatalf("resolveAndCreateDisks count=3 idx=2: expected nil, got %v", err)
		}
		if len(paths) != 1 || paths[0] != `C:\VMs\boot2.vhd` {
			t.Errorf("resolveAndCreateDisks count=3 idx=2: got %v, want [C:\\VMs\\boot2.vhd]", paths)
		}
		if !strings.Contains(capturedCmd, "boot2.vhd") {
			t.Errorf("resolveAndCreateDisks count=3 idx=2: New-VHD command %q does not contain 'boot2.vhd'", capturedCmd)
		}
	})

	t.Run("count=3 creates all three numbered disks", func(t *testing.T) {
		for i := 1; i <= 3; i++ {
			index := i
			withPS(t,
				func(_ string) error { return nil },
				buildOutputMock,
			)
			paths, err := resolveAndCreateDisks(config, []string{"boot"}, index, 3)
			if err != nil {
				t.Fatalf("resolveAndCreateDisks count=3 idx=%d: expected nil, got %v", index, err)
			}
			want := `C:\VMs\boot` + strconv.Itoa(index) + `.vhd`
			if len(paths) != 1 || paths[0] != want {
				t.Errorf("resolveAndCreateDisks count=3 idx=%d: got %v, want [%s]", index, paths, want)
			}
		}
	})

	t.Run("import disk uses same path regardless of count", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			buildOutputMock,
		)
		paths, err := resolveAndCreateDisks(config, []string{"shared-iso"}, 2, 3)
		if err != nil {
			t.Fatalf("resolveAndCreateDisks import count=3: expected nil, got %v", err)
		}
		if len(paths) != 1 || paths[0] != `C:\ISOs\install.iso` {
			t.Errorf("resolveAndCreateDisks import count=3: got %v, want base path", paths)
		}
		if runCalled {
			t.Error("resolveAndCreateDisks import: runPS should not be called for import disk")
		}
	})

	t.Run("raw path reference is passed through as-is", func(t *testing.T) {
		withPS(t, func(_ string) error { return nil }, buildOutputMock)
		paths, err := resolveAndCreateDisks(config, []string{`C:\raw\path.vhd`}, 2, 3)
		if err != nil {
			t.Fatalf("resolveAndCreateDisks raw path: expected nil, got %v", err)
		}
		if len(paths) != 1 || paths[0] != `C:\raw\path.vhd` {
			t.Errorf("resolveAndCreateDisks raw path: got %v, want literal path", paths)
		}
	})
}

// removeOutputMock routes outputPS responses for RemoveByStruct tests.
// Get-VMSwitch must be checked before Get-VM because the former contains the latter.
func removeOutputMock(cmd string) ([]byte, error) {
	switch {
	case strings.Contains(cmd, "Test-Path"):
		return []byte("True\n"), nil // file exists → can remove
	case strings.Contains(cmd, "Get-VMSwitch"):
		return []byte("SwitchType\n----------\nInternal\n"), nil // switch exists
	case strings.Contains(cmd, "Get-VM"):
		return stateOutput("Running"), nil // VM exists
	default:
		return nil, errors.New("unexpected call")
	}
}

// --- RemoveByStruct ---

func TestRemoveByStruct(t *testing.T) {
	t.Run("VMs networks and disks are removed", func(t *testing.T) {
		withPS(t,
			func(_ string) error { return nil },
			removeOutputMock,
		)
		var buf bytes.Buffer
		config := Summarize{
			Vms: map[string]VM{
				"my-vm": {Count: 1},
			},
			Networks: map[string]VMSwitch{
				"internal": {Type: "internal"},
			},
			Disks: map[string]Disk{
				"data": {Path: `C:\data.vhd`},
			},
		}
		err := RemoveByStruct(config, &buf)
		if err != nil {
			t.Errorf("RemoveByStruct: expected nil, got %v", err)
		}
		if buf.Len() > 0 {
			t.Errorf("RemoveByStruct: unexpected output: %s", buf.String())
		}
	})

	t.Run("import disk is skipped", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error {
				runCalled = true
				return nil
			},
			func(_ string) ([]byte, error) {
				return stateOutput("Running"), nil
			},
		)
		var buf bytes.Buffer
		config := Summarize{
			Disks: map[string]Disk{
				"imported": {Path: `C:\existing.vhd`, Import: true},
			},
		}
		_ = RemoveByStruct(config, &buf)
		if runCalled {
			t.Error("RemoveByStruct: runPS was called for import disk")
		}
	})

	t.Run("count 0 treated as 1 — no suffix", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error {
				capturedCmd = c
				return nil
			},
			func(_ string) ([]byte, error) {
				return stateOutput("Running"), nil
			},
		)
		var buf bytes.Buffer
		config := Summarize{
			Vms: map[string]VM{"solo-vm": {Count: 0}},
		}
		_ = RemoveByStruct(config, &buf)
		if !strings.Contains(capturedCmd, "solo-vm") {
			t.Errorf("RemoveByStruct: expected 'solo-vm' in command, got %q", capturedCmd)
		}
		if strings.Contains(capturedCmd, "solo-vm1") {
			t.Errorf("RemoveByStruct: count=0 should not add numeric suffix, got %q", capturedCmd)
		}
	})

	t.Run("count > 1 generates suffixed names", func(t *testing.T) {
		var capturedCmds []string
		withPS(t,
			func(c string) error {
				capturedCmds = append(capturedCmds, c)
				return nil
			},
			func(_ string) ([]byte, error) {
				return stateOutput("Running"), nil
			},
		)
		var buf bytes.Buffer
		config := Summarize{
			Vms: map[string]VM{"web": {Count: 3}},
		}
		_ = RemoveByStruct(config, &buf)
		for _, suffix := range []string{"web1", "web2", "web3"} {
			found := false
			for _, cmd := range capturedCmds {
				if strings.Contains(cmd, suffix) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("RemoveByStruct: expected %q in commands, got %v", suffix, capturedCmds)
			}
		}
	})

	t.Run("count>1 removes numbered disk copies", func(t *testing.T) {
		var capturedCmds []string
		withPS(t,
			func(c string) error {
				capturedCmds = append(capturedCmds, c)
				return nil
			},
			removeOutputMock,
		)
		var buf bytes.Buffer
		config := Summarize{
			Vms: map[string]VM{
				"web": {Count: 3, Disks: []string{"data"}},
			},
			Disks: map[string]Disk{
				"data": {Path: `C:\VMs\data.vhd`},
			},
		}
		_ = RemoveByStruct(config, &buf)
		for _, suffix := range []string{"data1.vhd", "data2.vhd", "data3.vhd"} {
			found := false
			for _, cmd := range capturedCmds {
				if strings.Contains(cmd, suffix) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("RemoveByStruct count>1: numbered disk %q not removed; commands: %v", suffix, capturedCmds)
			}
		}
		// Base path must not be removed directly.
		for _, cmd := range capturedCmds {
			if strings.Contains(cmd, "Remove-Item") && strings.Contains(cmd, `C:\VMs\data.vhd"`) {
				t.Errorf("RemoveByStruct count>1: base disk path removed directly: %q", cmd)
			}
		}
	})

	t.Run("not found resources are silently skipped", func(t *testing.T) {
		withPS(t,
			nil, // runPS not reached when resources don't exist
			func(cmd string) ([]byte, error) {
				if strings.Contains(cmd, "Test-Path") {
					return []byte("False\n"), nil // disk does not exist
				}
				return nil, errors.New("not found") // VM/switch not found
			},
		)
		var buf bytes.Buffer
		config := Summarize{
			Vms:      map[string]VM{"ghost-vm": {Count: 1}},
			Networks: map[string]VMSwitch{"ghost-sw": {Type: "internal"}},
			Disks:    map[string]Disk{"ghost-disk": {Path: `C:\ghost.vhd`}},
		}
		err := RemoveByStruct(config, &buf)
		if err != nil {
			t.Errorf("RemoveByStruct: expected nil for not-found resources, got %v", err)
		}
		if buf.Len() > 0 {
			t.Errorf("RemoveByStruct: expected no output for not-found resources, got: %s", buf.String())
		}
	})

	t.Run("genuine removal failure writes to writer and returns error", func(t *testing.T) {
		withPS(t,
			func(_ string) error {
				return errors.New("access denied")
			},
			func(_ string) ([]byte, error) {
				return stateOutput("Running"), nil // VM exists
			},
		)
		var buf bytes.Buffer
		config := Summarize{
			Vms: map[string]VM{"my-vm": {Count: 1}},
		}
		err := RemoveByStruct(config, &buf)
		if err == nil {
			t.Fatal("RemoveByStruct: expected error for removal failure, got nil")
		}
		if buf.Len() == 0 {
			t.Error("RemoveByStruct: expected error written to writer, got nothing")
		}
	})
}

// --- StartByStruct / StopByStruct ---

func TestStartByStruct(t *testing.T) {
	t.Run("off VM is started", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error { capturedCmd = c; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Off"), nil },
		)
		config := Summarize{Vms: map[string]VM{"router": {Count: 1}}}
		if err := StartByStruct(config, &bytes.Buffer{}); err != nil {
			t.Fatalf("StartByStruct: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "Start-VM") || !strings.Contains(capturedCmd, "router") {
			t.Errorf("StartByStruct: unexpected command %q", capturedCmd)
		}
	})

	t.Run("already running VM is skipped", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		config := Summarize{Vms: map[string]VM{"router": {Count: 1}}}
		if err := StartByStruct(config, &bytes.Buffer{}); err != nil {
			t.Fatalf("StartByStruct: expected nil, got %v", err)
		}
		if runCalled {
			t.Error("StartByStruct: runPS called for already-running VM")
		}
	})

	t.Run("not found VM is skipped", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return nil, errors.New("not found") },
		)
		config := Summarize{Vms: map[string]VM{"ghost": {Count: 1}}}
		if err := StartByStruct(config, &bytes.Buffer{}); err != nil {
			t.Fatalf("StartByStruct: expected nil for not-found VM, got %v", err)
		}
		if runCalled {
			t.Error("StartByStruct: runPS called for not-found VM")
		}
	})

	t.Run("count>1 starts all instances", func(t *testing.T) {
		var (
			mu           sync.Mutex
			capturedCmds []string
		)
		withPS(t,
			func(c string) error {
				mu.Lock()
				capturedCmds = append(capturedCmds, c)
				mu.Unlock()
				return nil
			},
			func(_ string) ([]byte, error) { return stateOutput("Off"), nil },
		)
		config := Summarize{Vms: map[string]VM{"router": {Count: 3}}}
		if err := StartByStruct(config, &bytes.Buffer{}); err != nil {
			t.Fatalf("StartByStruct count>1: expected nil, got %v", err)
		}
		mu.Lock()
		cmds := capturedCmds
		mu.Unlock()
		for _, name := range []string{"router1", "router2", "router3"} {
			found := false
			for _, cmd := range cmds {
				if strings.Contains(cmd, name) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("StartByStruct count>1: %q not started; commands: %v", name, cmds)
			}
		}
	})
}

func TestStopByStruct(t *testing.T) {
	t.Run("running VM is stopped", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error { capturedCmd = c; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		config := Summarize{Vms: map[string]VM{"router": {Count: 1}}}
		if err := StopByStruct(config, &bytes.Buffer{}); err != nil {
			t.Fatalf("StopByStruct: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "Stop-VM") || !strings.Contains(capturedCmd, "router") {
			t.Errorf("StopByStruct: unexpected command %q", capturedCmd)
		}
	})

	t.Run("non-running VM is skipped", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Off"), nil },
		)
		config := Summarize{Vms: map[string]VM{"router": {Count: 1}}}
		if err := StopByStruct(config, &bytes.Buffer{}); err != nil {
			t.Fatalf("StopByStruct: expected nil, got %v", err)
		}
		if runCalled {
			t.Error("StopByStruct: runPS called for non-running VM")
		}
	})
}

func TestRestartByStruct(t *testing.T) {
	t.Run("running VM is restarted", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error { capturedCmd = c; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		config := Summarize{Vms: map[string]VM{"router": {Count: 1}}}
		if err := RestartByStruct(config, &bytes.Buffer{}); err != nil {
			t.Fatalf("RestartByStruct: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "Restart-VM") || !strings.Contains(capturedCmd, "router") {
			t.Errorf("RestartByStruct: unexpected command %q", capturedCmd)
		}
	})

	t.Run("non-running VM is skipped", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Off"), nil },
		)
		config := Summarize{Vms: map[string]VM{"router": {Count: 1}}}
		if err := RestartByStruct(config, &bytes.Buffer{}); err != nil {
			t.Fatalf("RestartByStruct: expected nil, got %v", err)
		}
		if runCalled {
			t.Error("RestartByStruct: runPS called for non-running VM")
		}
	})
}

func TestSaveByStruct(t *testing.T) {
	t.Run("running VM is saved", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error { capturedCmd = c; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		config := Summarize{Vms: map[string]VM{"router": {Count: 1}}}
		if err := SaveByStruct(config, &bytes.Buffer{}); err != nil {
			t.Fatalf("SaveByStruct: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "Save-VM") || !strings.Contains(capturedCmd, "router") {
			t.Errorf("SaveByStruct: unexpected command %q", capturedCmd)
		}
	})

	t.Run("non-running VM is skipped", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Off"), nil },
		)
		config := Summarize{Vms: map[string]VM{"router": {Count: 1}}}
		if err := SaveByStruct(config, &bytes.Buffer{}); err != nil {
			t.Fatalf("SaveByStruct: expected nil, got %v", err)
		}
		if runCalled {
			t.Error("SaveByStruct: runPS called for non-running VM")
		}
	})
}

func TestResumeByStruct(t *testing.T) {
	t.Run("saved VM is resumed", func(t *testing.T) {
		var capturedCmd string
		withPS(t,
			func(c string) error { capturedCmd = c; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Saved"), nil },
		)
		config := Summarize{Vms: map[string]VM{"router": {Count: 1}}}
		if err := ResumeByStruct(config, &bytes.Buffer{}); err != nil {
			t.Fatalf("ResumeByStruct: expected nil, got %v", err)
		}
		if !strings.Contains(capturedCmd, "Start-VM") || !strings.Contains(capturedCmd, "router") {
			t.Errorf("ResumeByStruct: unexpected command %q", capturedCmd)
		}
	})

	t.Run("non-saved VM is skipped", func(t *testing.T) {
		runCalled := false
		withPS(t,
			func(_ string) error { runCalled = true; return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		config := Summarize{Vms: map[string]VM{"router": {Count: 1}}}
		if err := ResumeByStruct(config, &bytes.Buffer{}); err != nil {
			t.Fatalf("ResumeByStruct: expected nil, got %v", err)
		}
		if runCalled {
			t.Error("ResumeByStruct: runPS called for non-saved VM")
		}
	})
}

// TestByStructConcurrent verifies that the parallel ByStruct functions are race-free.
// Run with -race to detect data races.
func TestByStructConcurrent(t *testing.T) {
	multiConfig := Summarize{
		Vms: map[string]VM{
			"alpha": {Count: 1},
			"beta":  {Count: 3},
		},
	}

	t.Run("StartByStruct no data race", func(t *testing.T) {
		withPS(t,
			func(_ string) error { return nil },
			func(_ string) ([]byte, error) { return stateOutput("Off"), nil },
		)
		for i := 0; i < 20; i++ {
			if err := StartByStruct(multiConfig, &bytes.Buffer{}); err != nil {
				t.Fatalf("iteration %d: unexpected error: %v", i, err)
			}
		}
	})

	t.Run("StopByStruct no data race", func(t *testing.T) {
		withPS(t,
			func(_ string) error { return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		for i := 0; i < 20; i++ {
			if err := StopByStruct(multiConfig, &bytes.Buffer{}); err != nil {
				t.Fatalf("iteration %d: unexpected error: %v", i, err)
			}
		}
	})

	t.Run("RestartByStruct no data race", func(t *testing.T) {
		withPS(t,
			func(_ string) error { return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		for i := 0; i < 20; i++ {
			if err := RestartByStruct(multiConfig, &bytes.Buffer{}); err != nil {
				t.Fatalf("iteration %d: unexpected error: %v", i, err)
			}
		}
	})

	t.Run("SaveByStruct no data race", func(t *testing.T) {
		withPS(t,
			func(_ string) error { return nil },
			func(_ string) ([]byte, error) { return stateOutput("Running"), nil },
		)
		for i := 0; i < 20; i++ {
			if err := SaveByStruct(multiConfig, &bytes.Buffer{}); err != nil {
				t.Fatalf("iteration %d: unexpected error: %v", i, err)
			}
		}
	})

	t.Run("ResumeByStruct no data race", func(t *testing.T) {
		withPS(t,
			func(_ string) error { return nil },
			func(_ string) ([]byte, error) { return stateOutput("Saved"), nil },
		)
		for i := 0; i < 20; i++ {
			if err := ResumeByStruct(multiConfig, &bytes.Buffer{}); err != nil {
				t.Fatalf("iteration %d: unexpected error: %v", i, err)
			}
		}
	})

	t.Run("all instances execute despite concurrent errors", func(t *testing.T) {
		var (
			mu    sync.Mutex
			count int
		)
		withPS(t,
			func(_ string) error {
				mu.Lock()
				count++
				mu.Unlock()
				return errors.New("simulated error")
			},
			func(_ string) ([]byte, error) { return stateOutput("Off"), nil },
		)
		config := Summarize{Vms: map[string]VM{"router": {Count: 4}}}
		err := StartByStruct(config, &bytes.Buffer{})
		if err == nil {
			t.Fatal("StartByStruct: expected error, got nil")
		}
		mu.Lock()
		got := count
		mu.Unlock()
		if got != 4 {
			t.Errorf("StartByStruct: expected 4 PS calls, got %d", got)
		}
	})
}
