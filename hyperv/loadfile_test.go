package hyperv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUnmarshalYamlEnvExpansion(t *testing.T) {
	t.Run("env var in path is expanded", func(t *testing.T) {
		t.Setenv("TEST_VM_PATH", `C:\VMs\test`)
		t.Setenv("TEST_MEM", "2GB")

		yaml := `
vms:
  myvm:
    generation: 2
    path: ${TEST_VM_PATH}
    memory:
      size: ${TEST_MEM}
      dynamic: false
    cpu:
      thread: 1
`
		dir := t.TempDir()
		file := filepath.Join(dir, "manahy.yaml")
		if err := os.WriteFile(file, []byte(yaml), 0600); err != nil {
			t.Fatal(err)
		}

		data, err := UnmarshalYaml(file)
		if err != nil {
			t.Fatalf("UnmarshalYaml: %v", err)
		}
		vm := data.Vms["myvm"]
		if vm.Path != `C:\VMs\test` {
			t.Errorf("expected path %q, got %q", `C:\VMs\test`, vm.Path)
		}
		if vm.Memory.Size != "2GB" {
			t.Errorf("expected memory size %q, got %q", "2GB", vm.Memory.Size)
		}
	})

	t.Run("undefined env var expands to empty string", func(t *testing.T) {
		os.Unsetenv("UNDEFINED_VAR_XYZ")

		yaml := `
vms:
  myvm:
    generation: 1
    path: ${UNDEFINED_VAR_XYZ}
    memory:
      size: 1GB
      dynamic: false
    cpu:
      thread: 1
`
		dir := t.TempDir()
		file := filepath.Join(dir, "manahy.yaml")
		if err := os.WriteFile(file, []byte(yaml), 0600); err != nil {
			t.Fatal(err)
		}

		data, err := UnmarshalYaml(file)
		if err != nil {
			t.Fatalf("UnmarshalYaml: %v", err)
		}
		vm := data.Vms["myvm"]
		if vm.Path != "" {
			t.Errorf("expected empty path for undefined var, got %q", vm.Path)
		}
	})

	t.Run("file not found returns error", func(t *testing.T) {
		_, err := UnmarshalYaml("nonexistent_manahy.yaml")
		if err == nil {
			t.Fatal("UnmarshalYaml: expected error for missing file, got nil")
		}
	})
}
