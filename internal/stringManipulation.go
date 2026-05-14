package internal

import (
	"fmt"
)

// Add trailing whitespace for size to text
func SizeAdjustment(text string, size int) string {
	for len(text) < size {
		text = text + " "
	}
	return text
}

// Display headers based on washed arguments
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
		for j := 0; j < len(header[i]); j++ {
			fmt.Printf("-")
		}
		fmt.Printf("\t")
	}
	fmt.Printf("\n")
}
