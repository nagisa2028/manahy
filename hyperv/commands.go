package hyperv

// psExecutable is the PowerShell executable.
// Change to "pwsh" for PowerShell 7+.
const psExecutable = "powershell"

// File system
const (
	cmdTestPath   = "Test-Path"
	cmdRemoveItem = "Remove-Item"
)

// VHD / disk
const (
	cmdNewVHD      = "New-VHD"
	cmdGetVHD      = "Get-VHD"
	cmdResizeVHD   = "Resize-VHD"
	cmdOptimizeVHD = "Optimize-VHD"
	cmdConvertVHD  = "Convert-VHD"
	cmdMountVHD    = "Mount-VHD"
	cmdDismountVHD = "Dismount-VHD"
	cmdMergeVHD    = "Merge-VHD"
)

// VM lifecycle
const (
	cmdGetVM         = "Get-VM"
	cmdNewVM         = "New-VM"
	cmdRemoveVM      = "Remove-VM"
	cmdRenameVM      = "Rename-VM"
	cmdStartVM       = "Start-VM"
	cmdStopVM        = "Stop-VM"
	cmdSaveVM        = "Save-VM"
	cmdSuspendVM     = "Suspend-VM"
	cmdRestartVM     = "Restart-VM"
	cmdResumeVM      = "Resume-VM"
	cmdExportVM      = "Export-VM"
	cmdImportVM      = "Import-VM"
	cmdMoveVMStorage = "Move-VMStorage"
	cmdVMConnect     = "vmconnect"
)

// VM configuration
const (
	cmdSetVMProcessor = "Set-VMProcessor"
	cmdSetVMMemory    = "Set-VMMemory"
)

// VM info / metrics
const (
	cmdGetVMProcessor              = "Get-VMProcessor"
	cmdGetVMMemory                 = "Get-VMMemory"
	cmdEnableVMResourceMetering    = "Enable-VMResourceMetering"
	cmdMeasureVM                   = "Measure-VM"
	cmdGetVMIntegrationService     = "Get-VMIntegrationService"
	cmdEnableVMIntegrationService  = "Enable-VMIntegrationService"
	cmdDisableVMIntegrationService = "Disable-VMIntegrationService"
	cmdGetVMHost                   = "Get-VMHost"
)

// Network adapter
const (
	cmdGetVMNetworkAdapter     = "Get-VMNetworkAdapter"
	cmdAddVMNetworkAdapter     = "Add-VMNetworkAdapter"
	cmdRemoveVMNetworkAdapter  = "Remove-VMNetworkAdapter"
	cmdConnectVMNetworkAdapter = "Connect-VMNetworkAdapter"
)

// Hard disk drive
const (
	cmdGetVMHardDiskDrive    = "Get-VMHardDiskDrive"
	cmdAddVMHardDiskDrive    = "Add-VMHardDiskDrive"
	cmdRemoveVMHardDiskDrive = "Remove-VMHardDiskDrive"
)

// DVD drive
const (
	cmdGetVMDvdDrive    = "Get-VMDvdDrive"
	cmdAddVMDvdDrive    = "Add-VMDvdDrive"
	cmdRemoveVMDvdDrive = "Remove-VMDvdDrive"
	cmdSetVMDvdDrive    = "Set-VMDvdDrive"
)

// Virtual switch
const (
	cmdGetVMSwitch    = "Get-VMSwitch"
	cmdNewVMSwitch    = "New-VMSwitch"
	cmdRemoveVMSwitch = "Remove-VMSwitch"
	cmdRenameVMSwitch = "Rename-VMSwitch"
	cmdSetVMSwitch    = "Set-VMSwitch"
)

// Checkpoint
const (
	cmdCheckpointVM        = "Checkpoint-VM"
	cmdGetVMCheckpoint     = "Get-VMCheckpoint"
	cmdRestoreVMCheckpoint = "Restore-VMCheckpoint"
	cmdRemoveVMCheckpoint  = "Remove-VMCheckpoint"
	cmdRenameVMCheckpoint  = "Rename-VMCheckpoint"
	cmdExportVMCheckpoint  = "Export-VMCheckpoint"
)

// Storage
const (
	cmdGetDisk = "Get-Disk"
)

// Authority
const (
	cmdGetLocalGroupMember    = "Get-LocalGroupMember"
	cmdAddLocalGroupMember    = "Add-LocalGroupMember"
	cmdRemoveLocalGroupMember = "Remove-LocalGroupMember"
	hvAdminsGroup             = "Hyper-V Administrators"
)
