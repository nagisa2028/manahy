package hyperv

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

const psTimeout = 30 * time.Second

// psLongTimeout is used for operations that transfer or process large disk images.
const psLongTimeout = 10 * time.Minute

// ps returns s as a PowerShell single-quoted string literal.
// Single quotes within s are escaped by doubling them, preventing injection
// when user-controlled values are embedded in PowerShell commands.
func ps(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// runPS executes a PowerShell command string with a 30-second timeout.
func runPS(cmd string) error {
	ctx, cancel := context.WithTimeout(context.Background(), psTimeout)
	defer cancel()
	return exec.CommandContext(ctx, "powershell", "-NoProfile", cmd).Run()
}

// outputPS executes a PowerShell command string and returns its output.
func outputPS(cmd string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), psTimeout)
	defer cancel()
	return exec.CommandContext(ctx, "powershell", "-NoProfile", cmd).Output()
}

// runPSLong executes a PowerShell command with a 10-minute timeout for long-running
// operations such as exporting, importing, moving, or converting large disk images.
func runPSLong(cmd string) error {
	ctx, cancel := context.WithTimeout(context.Background(), psLongTimeout)
	defer cancel()
	return exec.CommandContext(ctx, "powershell", "-NoProfile", cmd).Run()
}
