package hyperv

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// UnmarshalYaml reads and parses a manahy YAML configuration file.
func UnmarshalYaml(name string) (Summarize, error) {
	buf, err := loadFile(name)
	if err != nil {
		return Summarize{}, err
	}

	var data Summarize
	if err = yaml.Unmarshal(buf, &data); err != nil {
		return Summarize{}, err
	}
	return data, nil
}

func loadFile(name string) ([]byte, error) {
	if err := isFileExist(name); err != nil {
		return nil, fmt.Errorf("file does not exist: %s", name)
	}

	buf, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s", name)
	}
	return buf, nil
}
