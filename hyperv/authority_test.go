package hyperv

import (
	"errors"
	"strings"
	"testing"
)

// ---------- GetGroupMember ----------

func TestGetGroupMember(t *testing.T) {
	t.Run("success: member list parsed correctly", func(t *testing.T) {
		// GetGroupMember skips the first two lines (header + dashes) then
		// collects non-empty lines.
		output := "Name\n----\nDOMAIN\\alice\nDOMAIN\\bob\n"
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte(output), nil
		})
		members, err := GetGroupMember()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(members) != 2 {
			t.Fatalf("expected 2 members, got %d: %v", len(members), members)
		}
		if !strings.Contains(members[0], "alice") {
			t.Errorf("members[0] = %q; expected to contain 'alice'", members[0])
		}
		if !strings.Contains(members[1], "bob") {
			t.Errorf("members[1] = %q; expected to contain 'bob'", members[1])
		}
	})

	t.Run("PS error → error returned", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("access denied")
		})
		_, err := GetGroupMember()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("empty member list", func(t *testing.T) {
		output := "Name\n----\n"
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte(output), nil
		})
		members, err := GetGroupMember()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(members) != 0 {
			t.Errorf("expected 0 members, got %d: %v", len(members), members)
		}
	})
}

// ---------- AddGroupMember ----------

func TestAddGroupMember(t *testing.T) {
	t.Run("success: nil returned", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte(""), nil
		})
		if err := AddGroupMember("alice"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("PS error → error returned", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("access denied")
		})
		if err := AddGroupMember("alice"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ---------- RemoveGroupMember ----------

func TestRemoveGroupMember(t *testing.T) {
	t.Run("success: nil returned", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return []byte(""), nil
		})
		if err := RemoveGroupMember("alice"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("PS error → error returned", func(t *testing.T) {
		withPS(t, nil, func(_ string) ([]byte, error) {
			return nil, errors.New("access denied")
		})
		if err := RemoveGroupMember("alice"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
