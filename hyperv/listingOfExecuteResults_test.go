package hyperv

import (
	"testing"
)

// ---------- computeCapacity ----------

func TestComputeCapacity(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		wantVal  float64
		wantUnit string
		wantErr  bool
	}{
		{name: "0 bytes", raw: "0", wantVal: 0, wantUnit: "B"},
		{name: "1023 bytes", raw: "1023", wantVal: 1023, wantUnit: "B"},
		{name: "1024 bytes = 1 KB", raw: "1024", wantVal: 1.0, wantUnit: "KB"},
		{name: "1 MB", raw: "1048576", wantVal: 1.0, wantUnit: "MB"},
		{name: "1 GB", raw: "1073741824", wantVal: 1.0, wantUnit: "GB"},
		{name: "1 TB", raw: "1099511627776", wantVal: 1.0, wantUnit: "TB"},
		{name: "non-numeric", raw: "abc", wantErr: true},
		{name: "empty string", raw: "", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val, unit, err := computeCapacity(tc.raw)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("computeCapacity(%q): expected error, got nil", tc.raw)
				}
				return
			}
			if err != nil {
				t.Fatalf("computeCapacity(%q): unexpected error: %v", tc.raw, err)
			}
			if val != tc.wantVal {
				t.Errorf("computeCapacity(%q): val = %v; want %v", tc.raw, val, tc.wantVal)
			}
			if unit != tc.wantUnit {
				t.Errorf("computeCapacity(%q): unit = %q; want %q", tc.raw, unit, tc.wantUnit)
			}
		})
	}
}

// ---------- listingOfExecuteResults ----------

func TestListingOfExecuteResults(t *testing.T) {
	tests := []struct {
		name  string
		input string
		flag  string
		want  []string
	}{
		{
			name:  "single value",
			input: "State\n-----\nRunning\n",
			flag:  "State",
			want:  []string{"Running"},
		},
		{
			name:  "multiple values",
			input: "Name\n----\nvm1\nvm2\n",
			flag:  "Name",
			want:  []string{"vm1", "vm2"},
		},
		{
			name:  "blank lines are skipped",
			input: "Name\n----\nvm1\n\nvm2\n",
			flag:  "Name",
			want:  []string{"vm1", "vm2"},
		},
		{
			name:  "dashes are skipped",
			input: "State\n---------\nOff\n",
			flag:  "State",
			want:  []string{"Off"},
		},
		{
			name:  "empty input returns nil",
			input: "",
			flag:  "State",
			want:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := listingOfExecuteResults([]byte(tc.input), tc.flag)
			if len(got) != len(tc.want) {
				t.Fatalf("listingOfExecuteResults: got %v; want %v", got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("listingOfExecuteResults[%d]: got %q; want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

// ---------- vmListingOfExecuteResults ----------

func TestVmListingOfExecuteResults(t *testing.T) {
	t.Run("running vm placed in Running list", func(t *testing.T) {
		input := "Name                   State\n----                   -----\nmy-vm                  Running\n"
		list, err := vmListingOfExecuteResults([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(list.Running) != 1 {
			t.Fatalf("expected 1 Running VM, got %d", len(list.Running))
		}
	})

	t.Run("saved vm placed in Saved list", func(t *testing.T) {
		input := "Name                   State\n----                   -----\nmy-vm                  Saved\n"
		list, err := vmListingOfExecuteResults([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(list.Saved) != 1 {
			t.Fatalf("expected 1 Saved VM, got %d", len(list.Saved))
		}
	})

	t.Run("off vm placed in Off list", func(t *testing.T) {
		input := "Name                   State\n----                   -----\nmy-vm                  Off\n"
		list, err := vmListingOfExecuteResults([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(list.Off) != 1 {
			t.Fatalf("expected 1 Off VM, got %d", len(list.Off))
		}
	})

	t.Run("paused vm placed in Paused list", func(t *testing.T) {
		input := "Name                   State\n----                   -----\nmy-vm                  Paused\n"
		list, err := vmListingOfExecuteResults([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(list.Paused) != 1 {
			t.Fatalf("expected 1 Paused VM, got %d", len(list.Paused))
		}
	})

	t.Run("unknown/transient state is skipped not errored", func(t *testing.T) {
		// VMs in transient states (Starting, Stopping, Saving, or any unrecognised
		// state) are silently skipped so that GetVMList does not fail while VMs are
		// transitioning.
		input := "Name                   State\n----                   -----\nmy-vm                  Crashed\ngood-vm                Running\n"
		list, err := vmListingOfExecuteResults([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error for unknown state: %v", err)
		}
		// "my-vm" with unknown state is skipped; "good-vm" is returned.
		if len(list.Running) != 1 || list.Running[0] != "good-vm" {
			t.Errorf("expected [good-vm] in Running, got %v", list.Running)
		}
	})

	t.Run("empty input returns empty VMList and nil", func(t *testing.T) {
		list, err := vmListingOfExecuteResults([]byte(""))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(list.Running)+len(list.Saved)+len(list.Off)+len(list.Paused) != 0 {
			t.Errorf("expected empty VMList, got %+v", list)
		}
	})
}

// ---------- switchListingOfExecuteResults ----------

func TestSwitchListingOfExecuteResults(t *testing.T) {
	t.Run("external switch placed in External list", func(t *testing.T) {
		input := "Name            SwitchType\n----            ----------\nmy-switch       External\n"
		list, err := switchListingOfExecuteResults([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(list.External) != 1 {
			t.Fatalf("expected 1 External switch, got %d", len(list.External))
		}
	})

	t.Run("internal switch placed in Internal list", func(t *testing.T) {
		input := "Name            SwitchType\n----            ----------\nmy-switch       Internal\n"
		list, err := switchListingOfExecuteResults([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(list.Internal) != 1 {
			t.Fatalf("expected 1 Internal switch, got %d", len(list.Internal))
		}
	})

	t.Run("private switch placed in Private list", func(t *testing.T) {
		input := "Name            SwitchType\n----            ----------\nmy-switch       Private\n"
		list, err := switchListingOfExecuteResults([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(list.Private) != 1 {
			t.Fatalf("expected 1 Private switch, got %d", len(list.Private))
		}
	})

	t.Run("unknown switch type is skipped not errored", func(t *testing.T) {
		// An unrecognised switch type is silently skipped so that GetSwitchList
		// does not fail if Hyper-V adds new types in future Windows versions.
		input := "Name            SwitchType\n----            ----------\nmy-switch       Bridged\ngood-switch     Internal\n"
		list, err := switchListingOfExecuteResults([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error for unknown switch type: %v", err)
		}
		// "my-switch" with unknown type is skipped; "good-switch" is returned.
		if len(list.Internal) != 1 || list.Internal[0] != "good-switch" {
			t.Errorf("expected [good-switch] in Internal, got %v", list.Internal)
		}
	})
}

// ---------- storageListingOfExecuteResults ----------

func TestStorageListingOfExecuteResults(t *testing.T) {
	t.Run("parses normal storage output", func(t *testing.T) {
		input := "Number FriendlyName Size\n------ ----------- ----\n0 WD 1073741824\n"
		list, err := storageListingOfExecuteResults([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(list.Number) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(list.Number))
		}
		if list.Number[0] != "0" {
			t.Errorf("Number[0] = %q; want %q", list.Number[0], "0")
		}
		if list.Size[0] != 1.0 {
			t.Errorf("Size[0] = %v; want 1.0", list.Size[0])
		}
		if list.SizeUnit[0] != "GB" {
			t.Errorf("SizeUnit[0] = %q; want %q", list.SizeUnit[0], "GB")
		}
	})

	t.Run("bad size value returns error", func(t *testing.T) {
		// reTrailNum.FindString on a line with no trailing number returns "".
		// computeCapacity("") returns an error.
		input := "Number FriendlyName Size\n------ ----------- ----\n0 NoSizeHere\n"
		_, err := storageListingOfExecuteResults([]byte(input))
		if err == nil {
			t.Fatal("expected error for missing size, got nil")
		}
	})
}
