package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/DevelopNaoki/manahy/hyperv"
)

func TestWithStackProgress(t *testing.T) {
	data := hyperv.Summarize{
		Vms: map[string]hyperv.VM{
			"vm1": {Count: 1},
			"vm2": {Count: 2},
		},
	}

	t.Run("nil writer calls fn and returns its result", func(t *testing.T) {
		called := false
		err := withStackProgress(nil, "start", data, func() error {
			called = true
			return nil
		})
		if !called {
			t.Error("withStackProgress: fn was not called with nil writer")
		}
		if err != nil {
			t.Errorf("withStackProgress: expected nil, got %v", err)
		}
	})

	t.Run("non-nil writer prints processing message and done on success", func(t *testing.T) {
		var buf bytes.Buffer
		err := withStackProgress(&buf, "start", data, func() error { return nil })
		if err != nil {
			t.Errorf("withStackProgress: expected nil, got %v", err)
		}
		out := buf.String()
		if !strings.Contains(out, "start") {
			t.Errorf("expected output to contain 'start', got %q", out)
		}
		if !strings.Contains(out, "3") {
			// vm1 has count=1, vm2 has count=2 → 3 total
			t.Errorf("expected output to contain total count '3', got %q", out)
		}
		if !strings.Contains(out, "done") {
			t.Errorf("expected output to contain 'done', got %q", out)
		}
	})

	t.Run("non-nil writer prints error message on failure", func(t *testing.T) {
		var buf bytes.Buffer
		want := errors.New("stop failed")
		err := withStackProgress(&buf, "stop", data, func() error { return want })
		if err != want {
			t.Errorf("withStackProgress: expected %v, got %v", want, err)
		}
		out := buf.String()
		if !strings.Contains(out, "errors") {
			t.Errorf("expected output to mention 'errors' on failure, got %q", out)
		}
	})
}
