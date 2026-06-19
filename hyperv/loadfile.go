package hyperv

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2" //nolint:depguard
)

// UnmarshalYaml reads and parses a manahy YAML configuration file.
// Environment variables in the form ${VAR} or $VAR are expanded before parsing.
// All environment variables visible to the process are eligible for expansion,
// including those that may contain credentials. Avoid referencing secret env vars
// in fields that appear in PowerShell commands (path, notes) to prevent them from
// being exposed in verbose output or Windows PowerShell event logs.
func UnmarshalYaml(name string) (Summarize, error) {
	buf, err := loadFile(name)
	if err != nil {
		return Summarize{}, err
	}

	// Expand $VAR and ${VAR} references in the config before unmarshalling.
	expanded := []byte(os.ExpandEnv(string(buf)))

	var data Summarize
	if err = yaml.Unmarshal(expanded, &data); err != nil {
		return Summarize{}, err
	}
	return data, nil
}

func loadFile(name string) ([]byte, error) {
	buf, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", name, err)
	}
	return buf, nil
}
