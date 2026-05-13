package hyperv

import (
	"fmt"
	"os/exec"
	"regexp"
)

// GetGroupMember returns the list of Hyper-V Administrators group members.
func GetGroupMember() ([]string, error) {
	res, err := exec.Command("powershell", "-NoProfile", "(Get-LocalGroupMember -Group 'Hyper-V Administrators' | Format-Table Name | Out-String).Trim()").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get Hyper-V Administrators members")
	}

	var members []string
	lines := regexp.MustCompile("\r\n|\n").Split(string(res), -1)
	for i, line := range lines {
		if i < 2 {
			continue
		}
		if line != "" {
			members = append(members, line)
		}
	}
	return members, nil
}

// AddGroupMember adds a user to the Hyper-V Administrators group.
func AddGroupMember(name string) error {
	_, err := exec.Command("powershell", "-NoProfile", "Add-LocalGroupMember -Group 'Hyper-V Administrators' -Member "+name).Output()
	if err != nil {
		return fmt.Errorf("failed to add %s to Hyper-V Administrators", name)
	}
	return nil
}

// RemoveGroupMember removes a user from the Hyper-V Administrators group.
func RemoveGroupMember(name string) error {
	_, err := exec.Command("powershell", "-NoProfile", "Remove-LocalGroupMember -Group 'Hyper-V Administrators' -Member "+name).Output()
	if err != nil {
		return fmt.Errorf("failed to remove %s from Hyper-V Administrators", name)
	}
	return nil
}
