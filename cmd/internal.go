package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/DevelopNaoki/manahy/hyperv"
)

// maxBulkArgs is the upper bound on resource names accepted by bulk operations.
const maxBulkArgs = 100

// confirmAction prompts the user to confirm a destructive action.
// Returns true immediately when force is set.
func confirmAction(prompt string, force bool) bool {
	if force {
		return true
	}
	fmt.Printf("%s [y/N]: ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.EqualFold(strings.TrimSpace(scanner.Text()), "y")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %s\n", err)
	}
	return false
}

func displayList(list []string, message string) {
	fmt.Println(message)
	for _, item := range list {
		fmt.Println("-", item)
	}
	fmt.Println()
}

func displayStorageList(storageList hyperv.StorageList) {
	fmt.Print("Storage\n")
	for i := range storageList.Number {
		fmt.Printf("- %s: %s: %.2f %s\n", storageList.Number[i], storageList.FriendlyName[i], storageList.Size[i], storageList.SizeUnit[i])
	}
	fmt.Print("\n")
	fmt.Print("More information, execute 'Get-Disk'\n")
}

// loadConfig loads a Summarize config from the given path.
// Returns an empty Summarize if path is empty.
// On read/parse error, prints a warning to stderr and returns an empty Summarize
// so the caller can still fall back to treating the argument as a literal value.
func loadConfig(path string) hyperv.Summarize {
	if path == "" {
		return hyperv.Summarize{}
	}
	config, err := hyperv.UnmarshalYaml(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not load config %s: %s\n", path, err)
		return hyperv.Summarize{}
	}
	return config
}

// runParallel runs fn for every name concurrently, writing each failure to w.
// Each failure is written immediately so the user sees partial results without
// waiting for the full batch to complete.
// Only the last error encountered is returned; callers that need to
// distinguish individual failures should inspect the w output.
func runParallel(w io.Writer, names []string, fn func(string) error) error {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		lastErr error
	)
	for _, name := range names {
		name := name
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fn(name); err != nil {
				mu.Lock()
				_, _ = fmt.Fprintf(w, "%s: %s\n", name, err)
				lastErr = err
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return lastErr
}

// resolveDisk resolves a disk alias to a path using the loaded config.
// If the arg matches a key in config.Disks, returns its path; otherwise returns arg as-is.
func resolveDisk(config hyperv.Summarize, arg string) string {
	if d, ok := config.Disks[arg]; ok {
		return d.Path
	}
	return arg
}

// withStackProgress wraps a stack operation with optional progress reporting.
// When w is non-nil it prints "label: processing N VM(s)..." before calling fn
// and "label: done" or "label: finished with errors" after.
// When w is nil it simply calls fn.
func withStackProgress(w io.Writer, label string, data hyperv.Summarize, fn func() error) error {
	if w == nil {
		return fn()
	}
	total := 0
	for _, vm := range data.Vms {
		c := vm.Count
		if c <= 0 {
			c = 1
		}
		total += c
	}
	fmt.Fprintf(w, "%s: processing %d VM(s)...\n", label, total)
	err := fn()
	if err != nil {
		fmt.Fprintf(w, "%s: finished with errors\n", label)
	} else {
		fmt.Fprintf(w, "%s: done\n", label)
	}
	return err
}
