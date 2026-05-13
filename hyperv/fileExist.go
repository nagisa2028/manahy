// hyperv package is manage Hyper-V
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
		return false, fmt.Errorf("error: failed execute 'Test-Path'")
	}
	exist, _ := strconv.ParseBool(strings.Replace(string(res), "\r\n", "", -1))
	return exist, nil
}

func isFileExist(path string) error {
	exist, err := searchFilePath(path)
	if err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("%s is not exist", path)
	}
	return nil
}

func isNotFileExist(path string) error {
	exist, err := searchFilePath(path)
	if err != nil {
		return err
	} else if exist {
		return fmt.Errorf("%s is already exist", path)
	}
	return nil
}
