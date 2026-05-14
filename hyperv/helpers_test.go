package hyperv

import "testing"

// withPS sets psRunHook and psOutputHook for the duration of t and restores
// the originals via t.Cleanup. Place calls to this function at the start of
// each test (or sub-test) that needs to mock PowerShell execution.
func withPS(t *testing.T, runFn func(string) error, outputFn func(string) ([]byte, error)) {
	t.Helper()
	origRun, origOut := psRunHook, psOutputHook
	psRunHook = runFn
	psOutputHook = outputFn
	t.Cleanup(func() { psRunHook = origRun; psOutputHook = origOut })
}
