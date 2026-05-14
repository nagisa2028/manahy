package internal

import (
	"regexp"
	"strconv"
)

// ConversionBtoXB is conversion for unit [0-9]*B to [0-9]*[KMGTP]B
// ex) input: 2000B output: 2KB
func ConversionBtoXB(num string) (string, error) {
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return "", err
	}
	unit := "B"
	for f >= 1000 {
		f = f / 1000
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
		}
	}
	return strconv.FormatFloat(f, 'f', 0, 64) + unit, nil
}

// ConversionXBtoB is conversion for unit [0-9][KMGTP]B to [0-9]*B
// ex) input: 2KB output: 2000B
func ConversionXBtoB(value string) (string, error) {
	num := regexp.MustCompile("[0-9.]*").Split(value, -1)[0]
	unit := regexp.MustCompile("[KMGTP]B").Split(value, -1)[0]
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return "", err
	}

	for unit != "B" {
		f = f * 1000
		switch unit {
		case "PB":
			unit = "TB"
		case "TB":
			unit = "GB"
		case "GB":
			unit = "MB"
		case "MB":
			unit = "KB"
		case "KB":
			unit = "B"
		}
	}
	return strconv.FormatFloat(f, 'f', 0, 64) + "B", nil
}
