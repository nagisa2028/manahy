package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize()
	RootCmd.AddCommand(
		versionCmd,
		vmCmd,
		switchCmd,
		diskCmd,
		storageCmd,
		memberCmd,
		checkpointCmd,
		hostCmd,
		build,
		remove,
	)
}

func init() {
	vmCmd.AddCommand(
		vmList,
		vmState,
		vmCreate,
		vmRemove,
		vmRename,
		vmStart,
		vmSave,
		vmShutdown,
		vmDestroy,
		vmSuspend,
		vmRestart,
		vmConnect,
		vmResume,
		vmExport,
		vmImport,
		vmMove,
		vmCopy,
		vmInfo,
		vmMeasure,
		vmIntegrationCmd,
		vmNicCmd,
		vmHardDiskCmd,
		vmDvdCmd,
	)

	vmIntegrationCmd.AddCommand(
		vmIntegrationList,
		vmIntegrationEnable,
		vmIntegrationDisable,
	)

	vmIntegrationEnable.Flags().StringVarP(&integrationServiceName, "name", "n", "", "integration service name")
	vmIntegrationDisable.Flags().StringVarP(&integrationServiceName, "name", "n", "", "integration service name")

	vmList.Flags().BoolVarP(&vmListOption.active, "active", "", true, "list active vm's")
	vmList.Flags().BoolVarP(&vmListOption.inactive, "inactive", "i", false, "list inactive vm's")
	vmList.Flags().BoolVarP(&vmListOption.saved, "saved", "s", false, "list saved vm's")
	vmList.Flags().BoolVarP(&vmListOption.paused, "paused", "p", false, "list paused vm's")
	vmList.Flags().BoolVarP(&vmListOption.all, "all", "a", false, "list all vm's")

	vmRename.Flags().StringVarP(&newVMName, "new-name", "n", "", "new vm name")
	vmExport.Flags().StringVarP(&exportPath, "path", "p", "", "export destination directory")
	vmMove.Flags().StringVarP(&exportPath, "path", "p", "", "destination storage directory")
	vmCopy.Flags().StringVarP(&newVMName, "new-name", "n", "", "new vm name")
	vmCopy.Flags().StringVarP(&exportPath, "path", "p", "", "directory for exported source files")

	vmCreate.Flags().StringVarP(&vm.Name, "name", "n", "", "new vm name")
	vmCreate.Flags().IntVarP(&vm.Generation, "generation", "g", 1, "set vm generation")
	vmCreate.Flags().IntVarP(&vm.CPU.Thread, "vcpus", "v", 1, "set vm vcpus")
	vmCreate.Flags().BoolVarP(&vm.CPU.Nested, "nested", "", false, "enable nested virtualization")
	vmCreate.Flags().StringVarP(&vm.Memory.Size, "memory", "m", "", "set vm memory")
	vmCreate.Flags().BoolVarP(&vm.Memory.Dynamic, "nodynamic", "", false, "disable dynamic memory")
	vmCreate.Flags().StringVarP(&vm.Path, "path", "p", "", "new vm path")
	vmCreate.Flags().StringVarP(&vm.Image, "image", "i", "", "image path")
	vmCreate.Flags().StringVarP(&vmDisk, "disk", "d", "", "disk path")
	vmCreate.Flags().StringVarP(&vmSwitch, "network", "s", "", "switch name")
}

func init() {
	switchCmd.AddCommand(
		switchList,
		switchCreate,
		switchRemove,
		switchRename,
		switchOptionCfgCmd,
	)

	switchList.Flags().BoolVarP(&switchListOption.external, "external", "e", false, "list external switches")
	switchList.Flags().BoolVarP(&switchListOption.internal, "internal", "i", false, "list internal switches")
	switchList.Flags().BoolVarP(&switchListOption.private, "private", "p", false, "list private switches")
	switchList.Flags().BoolVarP(&switchListOption.all, "all", "a", true, "list all switches")

	switchCreate.Flags().StringVarP(&switchCreateOption.Name, "name", "n", "", "set name")
	switchCreate.Flags().StringVarP(&switchCreateOption.Type, "type", "t", "", "set type")
	switchCreate.Flags().StringVarP(&switchCreateOption.ExternalInterface, "external-interface", "", "", "set external interface")
	switchCreate.Flags().BoolVarP(&switchCreateOption.AllowManagementOs, "allow-management-os", "", false, "set allow management os")

	switchRename.Flags().StringVarP(&newSwitchName, "new-name", "n", "", "rename switch")

	switchOptionCfgCmd.AddCommand(
		switchOptionCfgType,
		switchOptionCfgNetAdapter,
	)

	switchOptionCfgType.Flags().StringVarP(&switchType, "type", "t", "", "change switch type")
	switchOptionCfgNetAdapter.Flags().StringVarP(&netAdapter, "net-adapter", "n", "", "change network adapter")
}

func init() {
	diskCmd.AddCommand(
		diskCreate,
		diskRemove,
		diskInfo,
		diskResize,
		diskOptimize,
		diskConvert,
		diskMount,
		diskDismount,
		diskMerge,
	)

	diskCreate.Flags().StringVarP(&diskCreateOption.Path, "path", "p", "", "set path")
	diskCreate.Flags().StringVarP(&diskCreateOption.Size, "size", "s", "", "set size")
	diskCreate.Flags().StringVarP(&diskCreateOption.Type, "type", "t", "dynamic", "set type")
	diskCreate.Flags().StringVarP(&diskCreateOption.ParentPath, "parent-path", "", "", "set parent path")
	diskCreate.Flags().IntVarP(&diskCreateOption.SourceDisk, "source-disk", "", 0, "set source disk")

	diskResize.Flags().StringVarP(&resizeSize, "size", "s", "", "new size (e.g. 20GB)")
	diskConvert.Flags().StringVarP(&convertDestPath, "dest", "d", "", "destination path")
	diskConvert.Flags().StringVarP(&convertDiskType, "type", "t", "", "disk type (Dynamic or Fixed)")
	diskMerge.Flags().StringVarP(&mergeDest, "dest", "d", "", "destination path (default: merge into parent)")

	storageCmd.AddCommand(
		storageList,
	)

	vmNicCmd.AddCommand(vmNicList, vmNicAdd, vmNicRemove, vmNicConnect)
	vmNicAdd.Flags().StringVarP(&vmNicName, "name", "n", "", "adapter name")
	vmNicAdd.Flags().StringVarP(&vmNicSwitch, "switch", "s", "", "switch name")
	vmNicRemove.Flags().StringVarP(&vmNicName, "name", "n", "", "adapter name")
	vmNicConnect.Flags().StringVarP(&vmNicName, "name", "n", "", "adapter name")
	vmNicConnect.Flags().StringVarP(&vmNicSwitch, "switch", "s", "", "switch name")

	vmHardDiskCmd.AddCommand(vmHardDiskList, vmHardDiskAdd, vmHardDiskRemove)
	vmHardDiskAdd.Flags().StringVarP(&vmHardDiskPath, "path", "p", "", "VHD path")
	vmHardDiskRemove.Flags().StringVarP(&vmHardDiskPath, "path", "p", "", "VHD path")

	vmDvdCmd.AddCommand(vmDvdList, vmDvdAdd, vmDvdRemove, vmDvdSet, vmDvdEject)
	vmDvdSet.Flags().StringVarP(&dvdImagePath, "image", "i", "", "ISO image path")
}

func init() {
	memberCmd.AddCommand(
		memberListCmd,
		memberAddCmd,
		memberRemoveCmd,
	)
}

func init() {
	checkpointCmd.AddCommand(
		checkpointCreate,
		checkpointList,
		checkpointRestore,
		checkpointRemove,
		checkpointRename,
		checkpointExport,
	)

	checkpointCreate.Flags().StringVarP(&checkpointName, "name", "n", "", "checkpoint name")
	checkpointRestore.Flags().StringVarP(&checkpointName, "name", "n", "", "checkpoint name")
	checkpointRemove.Flags().StringVarP(&checkpointName, "name", "n", "", "checkpoint name")
	checkpointRename.Flags().StringVarP(&checkpointName, "name", "n", "", "checkpoint name")
	checkpointRename.Flags().StringVarP(&newCheckpointName, "new-name", "", "", "new checkpoint name")
	checkpointExport.Flags().StringVarP(&checkpointName, "name", "n", "", "checkpoint name")
	checkpointExport.Flags().StringVarP(&exportPath, "path", "p", "", "export destination directory")

	hostCmd.AddCommand(hostShow)
}
