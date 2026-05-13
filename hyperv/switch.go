// hyperv package is manage Hyper-V
package hyperv

import (
	"fmt"
	"os/exec"
	"strconv"
)

// ---------- //
// Get Switch
// ---------- //

// GetSwitchList get a list of Switch
func GetSwitchList() (switchList SwitchList, err error) {
	res, err := exec.Command("powershell", "-NoProfile", "Get-VMSwitch | Sort-Object SwitchType | Format-Table Name, SwitchType").Output()
	if err != nil {
		return switchList, err
	}

	switchList, err = switchListingOfExecuteResults(res)
	if err != nil {
		return switchList, err
	}
	return switchList, nil
}

// GetSwitchType get a Switch Type
func GetSwitchType(name string) (state string) {
	res, err := exec.Command("powershell", "-NoProfile", "Get-VMSwitch '"+name+"' | Format-Table SwitchType").Output()
	if err != nil {
		state = "NotFound"
	} else {
		switchType := listingOfExecuteResults(res, "SwitchType")
		if len(switchType) == 1 {
			state = switchType[0]
		} else {
			state = "Unknown"
		}
	}
	return state
}

func IsSwitchExist(name string) error {
	state := GetSwitchType(name)
	if state == "Unknown" {
		return fmt.Errorf("failed get switch state")
	} else if state == "NotFound" {
		return fmt.Errorf("%s is not found", name)
	}
	return nil
}

func IsNotSwitchExist(name string) error {
	state := GetSwitchType(name)
	if state == "Unknown" {
		return fmt.Errorf("failed get switch state")
	} else if state != "NotFound" {
		return fmt.Errorf("%s is already found", name)
	}
	return nil
}

// ---------------------- //
// Create / Remove Switch
// ---------------------- //

// CreateSwitch create new switch
func CreateSwitch(newSwitch VMSwitch, output bool) error {
	var cmd string
	err := checkSwitchParam(newSwitch)
	if output {
		PrintError("Check Switch Param", err)
	}
	if err != nil {
		return err
	}

	cmd = "New-VMSwitch -name '" + newSwitch.Name + "'"
	if newSwitch.Type == "external" {
		cmd += " -NetAdapterName '" + newSwitch.ExternalInterface + "'"
		cmd += " -AllowManagementOS $" + strconv.FormatBool(newSwitch.AllowManagementOs)
	} else {
		cmd += " -SwitchType " + newSwitch.Type
	}

	err = exec.Command("powershell", "-NoProfile", cmd).Run()
	if output {
		PrintError("Create Switch", err)
	}
	if err != nil {
		return fmt.Errorf("failed create new switch")
	}
	return nil
}

// RemoveSwitch remove switch
func RemoveSwitch(name string) error {
	err := IsSwitchExist(name)
	if err != nil {
		return err
	}
	err = exec.Command("powershell", "-NoProfile", "Remove-VMSwitch '"+name+"' -Force").Run()
	if err != nil {
		return err
	}
	return nil
}

// -------------------- //
// Update Switch Option
// -------------------- //

// RenameSwitch rename switch
func RenameSwitch(name string, newName string) error {
	err := IsSwitchExist(name)
	if err != nil {
		return err
	}
	err = IsNotSwitchExist(newName)
	if err != nil {
		return err
	}

	err = exec.Command("powershell", "-NoProfile", "Rename-VMSwitch '"+name+"' -NewName '"+newName+"'").Run()
	if err != nil {
		return err
	}
	return nil
}

// ChangeSwitchType change switch type
func ChangeSwitchType(name string, switchType string) error {
	nameType := GetSwitchType(name)
	switch nameType {
	case "NotFound":
		return fmt.Errorf("%s is not exist", name)
	case "Unknown":
		return fmt.Errorf("unknown error")
	default:
		err := checkSwitchTypeParam(switchType)
		if err != nil {
			return err
		} else if nameType == switchType {
			return fmt.Errorf("%s's now type already %s", name, switchType)
		}

		err = exec.Command("powershell", "-NoProfile", "Set-VMSwitch '"+name+"' -SwitchType "+switchType).Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// ChangeSwitchNetAdapter change net adapter of external switch
func ChangeSwitchNetAdapter(name string, netAdapter string) error {
	nameType := GetSwitchType(name)
	switch nameType {
	case "NotFound":
		return fmt.Errorf("%s is not exist", name)
	case "Unknown":
		return fmt.Errorf("unknown error")
	default:
		err := exec.Command("powershell", "-NoProfile", "Set-VMSwitch '"+name+"' -NetAdapterName '"+netAdapter+"'").Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// ------------------- //
// Check Switch Option
// ------------------- //

func checkSwitchParam(newSwitch VMSwitch) error {
	if GetSwitchType(newSwitch.Name) != "NotFound" {
		return fmt.Errorf("%s is already exist", newSwitch.Name)
	}

	err := checkSwitchTypeParam(newSwitch.Type)
	if err != nil {
		return err
	}

	err = checkSwitchParamIntegrity(newSwitch.Type, newSwitch.ExternalInterface)
	if err != nil {
		return err
	}
	return nil
}

func checkSwitchTypeParam(switchType string) error {
	switch switchType {
	case "external":
	case "internal":
	case "private":
	default:
		return fmt.Errorf("undefined switch type")
	}
	return nil
}

func checkSwitchParamIntegrity(switchType string, externalInterface string) error {
	if switchType == "external" && externalInterface == "" {
		return fmt.Errorf("need external-interface option")
	} else if switchType != "external" && externalInterface != "" {
		return fmt.Errorf("do not need external-interface option")
	}
	return nil
}
