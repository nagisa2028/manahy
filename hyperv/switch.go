package hyperv

import (
	"fmt"
	"os/exec"
	"strconv"
)

// GetSwitchList returns a list of all virtual switches grouped by type.
func GetSwitchList() (SwitchList, error) {
	res, err := exec.Command("powershell", "-NoProfile", "Get-VMSwitch | Sort-Object SwitchType | Format-Table Name, SwitchType").Output()
	if err != nil {
		return SwitchList{}, err
	}
	return switchListingOfExecuteResults(res)
}

// GetSwitchType returns the type of a virtual switch, or "NotFound" / "Unknown".
func GetSwitchType(name string) string {
	res, err := exec.Command("powershell", "-NoProfile", "Get-VMSwitch '"+name+"' | Format-Table SwitchType").Output()
	if err != nil {
		return "NotFound"
	}
	switchType := listingOfExecuteResults(res, "SwitchType")
	if len(switchType) == 1 {
		return switchType[0]
	}
	return "Unknown"
}

// IsSwitchExist returns an error if the switch does not exist.
func IsSwitchExist(name string) error {
	switch GetSwitchType(name) {
	case "Unknown":
		return fmt.Errorf("failed to get state of switch %s", name)
	case "NotFound":
		return fmt.Errorf("switch %s does not exist", name)
	}
	return nil
}

// IsNotSwitchExist returns an error if the switch already exists.
func IsNotSwitchExist(name string) error {
	switch GetSwitchType(name) {
	case "Unknown":
		return fmt.Errorf("failed to get state of switch %s", name)
	case "NotFound":
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

	cmd := "New-VMSwitch -name '" + newSwitch.Name + "'"
	if newSwitch.Type == "external" {
		cmd += " -NetAdapterName '" + newSwitch.ExternalInterface + "'"
		cmd += " -AllowManagementOS $" + strconv.FormatBool(newSwitch.AllowManagementOs)
	} else {
		cmd += " -SwitchType " + newSwitch.Type
	}

	err = exec.Command("powershell", "-NoProfile", cmd).Run()
	printError("Create Switch", err, output)
	if err != nil {
		return fmt.Errorf("failed to create switch %s", newSwitch.Name)
	}
	return nil
}

// RemoveSwitch deletes a virtual switch.
func RemoveSwitch(name string) error {
	if err := IsSwitchExist(name); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile", "Remove-VMSwitch '"+name+"' -Force").Run()
}

// RenameSwitch renames a virtual switch.
func RenameSwitch(name string, newName string) error {
	if err := IsSwitchExist(name); err != nil {
		return err
	}
	if err := IsNotSwitchExist(newName); err != nil {
		return err
	}
	return exec.Command("powershell", "-NoProfile", "Rename-VMSwitch '"+name+"' -NewName '"+newName+"'").Run()
}

// ChangeSwitchType changes the type of a virtual switch.
func ChangeSwitchType(name string, switchType string) error {
	nameType := GetSwitchType(name)
	switch nameType {
	case "NotFound":
		return fmt.Errorf("switch %s does not exist", name)
	case "Unknown":
		return fmt.Errorf("failed to get state of switch %s", name)
	}
	if err := checkSwitchTypeParam(switchType); err != nil {
		return err
	}
	if nameType == switchType {
		return fmt.Errorf("switch %s is already of type %s", name, switchType)
	}
	return exec.Command("powershell", "-NoProfile", "Set-VMSwitch '"+name+"' -SwitchType "+switchType).Run()
}

// ChangeSwitchNetAdapter changes the net adapter of an external virtual switch.
func ChangeSwitchNetAdapter(name string, netAdapter string) error {
	switch GetSwitchType(name) {
	case "NotFound":
		return fmt.Errorf("switch %s does not exist", name)
	case "Unknown":
		return fmt.Errorf("failed to get state of switch %s", name)
	}
	return exec.Command("powershell", "-NoProfile", "Set-VMSwitch '"+name+"' -NetAdapterName '"+netAdapter+"'").Run()
}

func checkSwitchParam(newSwitch VMSwitch) error {
	if GetSwitchType(newSwitch.Name) != "NotFound" {
		return fmt.Errorf("switch %s already exists", newSwitch.Name)
	}
	if err := checkSwitchTypeParam(newSwitch.Type); err != nil {
		return err
	}
	return checkSwitchParamIntegrity(newSwitch.Type, newSwitch.ExternalInterface)
}

func checkSwitchTypeParam(switchType string) error {
	switch switchType {
	case "external", "internal", "private":
		return nil
	default:
		return fmt.Errorf("undefined switch type: %s", switchType)
	}
}

func checkSwitchParamIntegrity(switchType string, externalInterface string) error {
	if switchType == "external" && externalInterface == "" {
		return fmt.Errorf("--external-interface is required for external switches")
	}
	if switchType != "external" && externalInterface != "" {
		return fmt.Errorf("--external-interface is only valid for external switches")
	}
	return nil
}
