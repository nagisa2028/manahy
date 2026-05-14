package hyperv

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const reSplitLine = "\r\n|\n"
const reBlankLine = "^[-\\s]*$"

func listingOfExecuteResults(res []byte, flag string) []string {
	var list []string
	split := regexp.MustCompile(reSplitLine).Split(string(res), -1)
	for _, line := range split {
		line = strings.Trim(line, " ")
		if line != flag && !regexp.MustCompile("^-*$").MatchString(line) && line != "" {
			list = append(list, line)
		}
	}
	return list
}

func vmListingOfExecuteResults(res []byte) (VMList, error) {
	var vmList VMList
	split := regexp.MustCompile(reSplitLine).Split(string(res), -1)
	for _, line := range split {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Name") || regexp.MustCompile(reBlankLine).MatchString(line) {
			continue
		}
		state := regexp.MustCompile("Running$|Saved$|Off$|Paused$").FindString(line)
		line = strings.TrimSpace(regexp.MustCompile("Running$|Saved$|Off$|Paused$").ReplaceAllString(line, ""))

		switch state {
		case "Running":
			vmList.Running = append(vmList.Running, line)
		case "Saved":
			vmList.Saved = append(vmList.Saved, line)
		case "Off":
			vmList.Off = append(vmList.Off, line)
		case "Paused":
			vmList.Paused = append(vmList.Paused, line)
		default:
			return vmList, fmt.Errorf("unknown VM state in output: %q", state)
		}
	}
	return vmList, nil
}

func switchListingOfExecuteResults(res []byte) (SwitchList, error) {
	var switchList SwitchList
	split := regexp.MustCompile(reSplitLine).Split(string(res), -1)
	for _, line := range split {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Name") || regexp.MustCompile(reBlankLine).MatchString(line) {
			continue
		}
		switchType := regexp.MustCompile("External$|Internal$|Private$").FindString(line)
		line = strings.TrimSpace(regexp.MustCompile("External$|Internal$|Private$").ReplaceAllString(line, ""))

		switch switchType {
		case "External":
			switchList.External = append(switchList.External, line)
		case "Internal":
			switchList.Internal = append(switchList.Internal, line)
		case "Private":
			switchList.Private = append(switchList.Private, line)
		default:
			return switchList, fmt.Errorf("unknown switch type in output: %q", switchType)
		}
	}
	return switchList, nil
}

func storageListingOfExecuteResults(res []byte) (StorageList, error) {
	var storageList StorageList
	split := regexp.MustCompile(reSplitLine).Split(string(res), -1)
	for _, line := range split {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Number") || regexp.MustCompile(reBlankLine).MatchString(line) {
			continue
		}
		storageList.Number = append(storageList.Number, regexp.MustCompile("^[0-9]+").FindString(line))
		line = regexp.MustCompile("^[0-9]+").ReplaceAllString(line, "")

		capacity, unit, err := computeCapacity(regexp.MustCompile("[0-9]+$").FindString(line))
		if err != nil {
			return storageList, err
		}
		storageList.Size = append(storageList.Size, capacity)
		storageList.SizeUnit = append(storageList.SizeUnit, unit)
		line = strings.TrimSpace(regexp.MustCompile("[0-9]+$").ReplaceAllString(line, ""))
		storageList.FriendlyName = append(storageList.FriendlyName, line)
	}
	return storageList, nil
}

func computeCapacity(raw string) (float64, string, error) {
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, "", fmt.Errorf("failed to parse capacity value: %s", raw)
	}
	unit := "B"
	for v >= 1024 {
		v /= 1024
		switch unit {
		case "B":
			unit = "KB"
		case "KB":
			unit = "MB"
		case "MB":
			unit = "GB"
		case "GB":
			unit = "TB"
		case "TB":
			unit = "PB"
		default:
			return 0, "", fmt.Errorf("unknown capacity unit: %s", unit)
		}
	}
	return v, unit, nil
}
