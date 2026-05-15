package hyperv

import (
	"fmt"
	"strconv"
	"strings"
)

func searchFilePath(path string) (bool, error) {
	res, e := outputPS(cmdTestPath + " " + ps(path))
	if e != nil {
		return false, fmt.Errorf("failed to execute Test-Path: %w", e)
	}
	exist, err := strconv.ParseBool(strings.TrimSpace(string(res)))
	if err != nil {
		return false, fmt.Errorf("unexpected Test-Path output for %s: %w", path, err)
	}
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
