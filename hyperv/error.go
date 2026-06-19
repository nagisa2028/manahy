package hyperv

import "fmt"

func printError(msg string, err error, output bool) {
	if !output {
		return
	}
	if err != nil {
		fmt.Printf("%s  : [\x1b[31mFail\x1b[0m]\n  %s\n", msg, err)
	} else {
		fmt.Printf("%s  : [\x1b[32mPass\x1b[0m]\n", msg)
	}
}
