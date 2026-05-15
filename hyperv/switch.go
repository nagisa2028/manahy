package hyperv

import (
	"fmt"
	"strconv"
	"strings"
)

// GetSwitchInfo returns detailed information about a virtual switch.
func GetSwitchInfo(name string) (string, error) {
	if err := IsSwitchExist(name); err != nil {
		return "", err
	}
	out, err := outputPS(cmdGetVMSwitch + " -Name " + ps(name) +
		" | Format-List Name, SwitchType, NetAdapterInterfaceDescription, AllowManagementOS, Notes | Out-String")
	if err != nil {
		return "", fmt.Errorf("failed to get info for switch %s: %w", name, err)
	}
	return string(out), nil
}

// GetSwitchList returns a list of all virtual switches grouped by type.
func GetSwitchList() (SwitchList, error) {
	res, err := outputPS(cmdGetVMSwitch + " | Sort-Object SwitchType | Format-Table Name, SwitchType")
	if err != nil {
		return SwitchList{}, err
	}
	return switchListingOfExecuteResults(res)
}

// GetSwitchType returns the type of a virtual switch, or "NotFound" / "Unknown".
func GetSwitchType(name string) string {
	res, err := outputPS(cmdGetVMSwitch + " " + ps(name) + " | Format-Table SwitchType")
	if err != nil {
		return vmStateNotFound
	}
	switchType := listingOfExecuteResults(res, "SwitchType")
	if len(switchType) == 1 {
		return switchType[0]
	}
	return vmStateUnknown
}

// IsSwitchExist returns an error if the switch does not exist.
func IsSwitchExist(name string) error {
	switch GetSwitchType(name) {
	case vmStateUnknown:
		return fmt.Errorf("failed to get state of switch %s", name)
	case vmStateNotFound:
		return fmt.Errorf("switch %s does not exist", name)
	}
	return nil
}

// IsNotSwitchExist returns an error if the switch already exists.
func IsNotSwitchExist(name string) error {
	switch GetSwitchType(name) {
	case vmStateUnknown:
		return fmt.Errorf("failed to get state of switch %s", name)
	case vmStateNotFound:
		return nil
	}
	return fmt.Errorf("switch %s already exists", name)
}

// CreateSwitch creates a new virtual switch.
func CreateSwitch(newSwitch VMSwitch, output bool) error {
	err := checkSwitchParam(newSwitch)
	printError("Check Switch Param", err, output)
	if err != nil {
		return err
	}

	cmd := cmdNewVMSwitch + " -name " + ps(newSwitch.Name)
	if newSwitch.Type == switchTypeExternal {
		cmd += " -NetAdapterName " + ps(newSwitch.ExternalInterface)
		cmd += " -AllowManagementOS $" + strconv.FormatBool(newSwitch.AllowManagementOS)
	} else {
		cmd += " -SwitchType " + ps(newSwitch.Type)
	}

	err = runPS(cmd)
	printError("Create Switch", err, output)
	if err != nil {
		return fmt.Errorf("failed to create switch %s: %w", newSwitch.Name, err)
	}
	return nil
}

// RemoveSwitch deletes a virtual switch.
func RemoveSwitch(name string) error {
	if err := IsSwitchExist(name); err != nil {
		return err
	}
	return runPS(cmdRemoveVMSwitch + " " + ps(name) + " -Force")
}

// RenameSwitch renames a virtual switch.
func RenameSwitch(name string, newName string) error {
	if err := IsSwitchExist(name); err != nil {
		return err
	}
	if err := IsNotSwitchExist(newName); err != nil {
		return err
	}
	return runPS(cmdRenameVMSwitch + " " + ps(name) + " -NewName " + ps(newName))
}

// ChangeSwitchType changes the type of a virtual switch.
func ChangeSwitchType(name string, switchType string) error {
	nameType := GetSwitchType(name)
	switch nameType {
	case vmStateNotFound:
		return fmt.Errorf("switch %s does not exist", name)
	case vmStateUnknown:
		return fmt.Errorf("failed to get state of switch %s", name)
	}
	if err := checkSwitchType(switchType); err != nil {
		return err
	}
	if strings.EqualFold(nameType, switchType) {
		return fmt.Errorf("switch %s is already of type %s", name, switchType)
	}
	return runPS(cmdSetVMSwitch + " " + ps(name) + " -SwitchType " + ps(switchType))
}

// ChangeSwitchNetAdapter changes the net adapter of an external virtual switch.
func ChangeSwitchNetAdapter(name string, netAdapter string) error {
	switch GetSwitchType(name) {
	case vmStateNotFound:
		return fmt.Errorf("switch %s does not exist", name)
	case vmStateUnknown:
		return fmt.Errorf("failed to get state of switch %s", name)
	}
	return runPS(cmdSetVMSwitch + " " + ps(name) + " -NetAdapterName " + ps(netAdapter))
}

func checkSwitchParam(newSwitch VMSwitch) error {
	if GetSwitchType(newSwitch.Name) != vmStateNotFound {
		return fmt.Errorf("switch %s already exists", newSwitch.Name)
	}
	if err := checkSwitchType(newSwitch.Type); err != nil {
		return err
	}
	return checkSwitchIntegrity(newSwitch.Type, newSwitch.ExternalInterface)
}

func checkSwitchType(switchType string) error {
	switch switchType {
	case switchTypeExternal, switchTypeInternal, switchTypePrivate:
		return nil
	default:
		return fmt.Errorf("invalid switch type: %s", switchType)
	}
}

func checkSwitchIntegrity(switchType string, externalInterface string) error {
	if switchType == switchTypeExternal && externalInterface == "" {
		return fmt.Errorf("external interface is required for external switches")
	}
	if switchType != switchTypeExternal && externalInterface != "" {
		return fmt.Errorf("external interface is only valid for external switches")
	}
	return nil
}
