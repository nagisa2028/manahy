// Package internal provides shared utility functions for manahy.
package internal

import (
	"fmt"
)

// SizeAdjustment pads text with trailing spaces until it reaches the given size.
func SizeAdjustment(text string, size int) string {
	for len(text) < size {
		text += " "
	}
	return text
}

// PrintHeader prints a formatted table header with separator line.
func PrintHeader(header []string, headerSize []int) {
	if len(header) != len(headerSize) {
		return
	}

	for i := range header {
		header[i] = SizeAdjustment(header[i], headerSize[i])
		fmt.Printf("%s\t", header[i])
	}
	fmt.Printf("\n")

	for i := range header {
		for range header[i] {
			fmt.Printf("-")
		}
		fmt.Printf("\t")
	}
	fmt.Printf("\n")
}
