// hyperv package is manage Hyper-V
package hyperv

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// UnmarshalYaml reads and parses a manahy YAML file
func UnmarshalYaml(name string) (Summarize, error) {
	buf, err := loadFile(name)
	if err != nil {
		return Summarize{}, err
	}

	var data Summarize
	err = yaml.Unmarshal(buf, &data)
	if err != nil {
		return Summarize{}, err
	}
	return data, nil
}

func loadFile(name string) ([]byte, error) {
	err := isFileExist(name)
	if err != nil {
		return nil, fmt.Errorf("file does not exist: %s", name)
	}

	buf, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed read %s", name)
	}
	return buf, nil
}
