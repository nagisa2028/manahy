package hyperv

import "strings"

// CheckStatus represents the result of a single check.
type CheckStatus int

// Check status values.
const (
	CheckOK   CheckStatus = iota
	CheckWarn             // check ran but result is uncertain
	CheckFail
)

// CheckResult holds one check item and its outcome.
type CheckResult struct {
	Name    string
	Status  CheckStatus
	Message string
}

// CheckSystem verifies that PowerShell is available, Hyper-V is enabled,
// and the current user is a member of the Hyper-V Administrators group.
func CheckSystem() []CheckResult {
	return []CheckResult{
		checkPSAvailable(),
		checkHyperVEnabled(),
		checkHyperVPermission(),
	}
}

// CheckVMCommands verifies that VM-related PowerShell cmdlets are available.
func CheckVMCommands() []CheckResult {
	return checkCmdlets([]string{
		cmdGetVM, cmdNewVM, cmdRemoveVM, cmdStartVM, cmdStopVM,
		cmdSaveVM, cmdSuspendVM, cmdResumeVM, cmdRestartVM,
		cmdRenameVM, cmdExportVM, cmdImportVM, cmdMoveVMStorage,
		cmdSetVMProcessor, cmdSetVMMemory,
		cmdCheckpointVM, cmdGetVMCheckpoint, cmdRestoreVMCheckpoint, cmdRemoveVMCheckpoint,
		cmdGetVMDvdDrive, cmdAddVMDvdDrive, cmdRemoveVMDvdDrive, cmdSetVMDvdDrive,
		cmdGetVMHardDiskDrive, cmdAddVMHardDiskDrive, cmdRemoveVMHardDiskDrive,
		cmdGetVMIntegrationService, cmdEnableVMIntegrationService, cmdDisableVMIntegrationService,
		cmdGetVMProcessor, cmdGetVMMemory, cmdEnableVMResourceMetering, cmdMeasureVM,
		cmdVMConnect,
	})
}

// CheckDiskCommands verifies that VHD-related PowerShell cmdlets are available.
func CheckDiskCommands() []CheckResult {
	return checkCmdlets([]string{
		cmdNewVHD, cmdGetVHD, cmdResizeVHD, cmdOptimizeVHD,
		cmdConvertVHD, cmdMountVHD, cmdDismountVHD, cmdMergeVHD,
		cmdTestPath, cmdRemoveItem, cmdGetDisk,
	})
}

// CheckNetworkCommands verifies that network-related PowerShell cmdlets are available.
func CheckNetworkCommands() []CheckResult {
	return checkCmdlets([]string{
		cmdGetVMSwitch, cmdNewVMSwitch, cmdRemoveVMSwitch, cmdSetVMSwitch, cmdRenameVMSwitch,
		cmdGetVMNetworkAdapter, cmdAddVMNetworkAdapter, cmdRemoveVMNetworkAdapter, cmdConnectVMNetworkAdapter,
	})
}

func checkPSAvailable() CheckResult {
	out, err := outputPS("$PSVersionTable.PSVersion.ToString()")
	if err != nil {
		return CheckResult{Name: "PowerShell (" + psExecutable + ")", Status: CheckFail, Message: "not available"}
	}
	return CheckResult{Name: "PowerShell (" + psExecutable + ")", Status: CheckOK, Message: strings.TrimSpace(string(out))}
}

func checkHyperVEnabled() CheckResult {
	_, err := outputPS(cmdGetVM + " -ErrorAction Stop | Out-Null")
	if err != nil {
		return CheckResult{Name: "Hyper-V feature", Status: CheckFail, Message: "not enabled or not accessible"}
	}
	return CheckResult{Name: "Hyper-V feature", Status: CheckOK, Message: "enabled"}
}

func checkHyperVPermission() CheckResult {
	script := `$id = [System.Security.Principal.WindowsIdentity]::GetCurrent(); ` +
		`([System.Security.Principal.WindowsPrincipal]$id).IsInRole("Hyper-V Administrators")`
	out, err := outputPS(script)
	if err != nil {
		return CheckResult{Name: "Hyper-V Administrators", Status: CheckWarn, Message: "could not verify"}
	}
	if strings.TrimSpace(string(out)) == "True" {
		return CheckResult{Name: "Hyper-V Administrators", Status: CheckOK, Message: "member"}
	}
	return CheckResult{Name: "Hyper-V Administrators", Status: CheckFail, Message: "not a member"}
}

func checkCmdlets(cmdlets []string) []CheckResult {
	if len(cmdlets) == 0 {
		return nil
	}
	quoted := make([]string, len(cmdlets))
	for i, c := range cmdlets {
		quoted[i] = ps(c)
	}
	cmd := "Get-Command " + strings.Join(quoted, ",") +
		" -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Name"
	out, err := outputPS(cmd)
	if err != nil {
		// PS itself failed (not found, timeout, etc.); report all cmdlets as unavailable.
		results := make([]CheckResult, len(cmdlets))
		for i, c := range cmdlets {
			results[i] = CheckResult{Name: c, Status: CheckFail, Message: "PS unavailable"}
		}
		return results
	}
	found := make(map[string]bool, len(cmdlets))
	for _, line := range reSplit.Split(string(out), -1) {
		if line := strings.TrimSpace(line); line != "" {
			found[line] = true
		}
	}
	results := make([]CheckResult, len(cmdlets))
	for i, c := range cmdlets {
		if found[c] {
			results[i] = CheckResult{Name: c, Status: CheckOK, Message: "available"}
		} else {
			results[i] = CheckResult{Name: c, Status: CheckFail, Message: "not found"}
		}
	}
	return results
}
