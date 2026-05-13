package hyperv

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func searchFilePath(path string) (bool, error) {
	res, e := exec.Command("powershell", "-NoProfile", "Test-Path \""+path+"\"").Output()
	if e != nil {
		return false, fmt.Errorf("failed to execute Test-Path")
	}
	exist, _ := strconv.ParseBool(strings.ReplaceAll(string(res), "\r\n", ""))
	return exist, nil
}

func isFileExist(path string) error {
	exist, err := searchFilePath(path)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("%s does not exist", path)
	}
	return nil
}

func isNotFileExist(path string) error {
	exist, err := searchFilePath(path)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("%s already exists", path)
	}
	return nil
}
