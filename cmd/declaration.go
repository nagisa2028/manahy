package cmd

import (
	"github.com/DevelopNaoki/manahy/hyperv"
)

var vmListOption struct {
	active   bool
	saved    bool
	inactive bool
	paused   bool
	all      bool
}
var newVMName string
var vm hyperv.VM
var vmDisk string
var vmSwitch string

var diskCreateOption hyperv.Disk

var switchListOption struct {
	external bool
	internal bool
	private  bool
	all      bool
}
var switchCreateOption hyperv.VMSwitch
var newSwitchName string
var switchType string
var netAdapter string

var exportPath string
var checkpointName string
var newCheckpointName string
var integrationServiceName string

var resizeSize string
var convertDestPath string
var convertDiskType string
var mergeDest string
var vmNicName string
var vmNicSwitch string
var vmHardDiskPath string
var dvdImagePath string
