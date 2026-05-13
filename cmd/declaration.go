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
var newVmName string
var vm hyperv.Vm
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
