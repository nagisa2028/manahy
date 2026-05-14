package hyperv

import "fmt"

// GetGroupMember returns the list of Hyper-V Administrators group members.
func GetGroupMember() ([]string, error) {
	res, err := outputPS("(" + cmdGetLocalGroupMember + " -Group " + ps(hvAdminsGroup) + " | Format-Table Name | Out-String).Trim()")
	if err != nil {
		return nil, fmt.Errorf("failed to get Hyper-V Administrators members")
	}

	var members []string
	for i, line := range reSplit.Split(string(res), -1) {
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
	_, err := outputPS(cmdAddLocalGroupMember + " -Group " + ps(hvAdminsGroup) + " -Member " + ps(name))
	if err != nil {
		return fmt.Errorf("failed to add %s to Hyper-V Administrators", name)
	}
	return nil
}

// RemoveGroupMember removes a user from the Hyper-V Administrators group.
func RemoveGroupMember(name string) error {
	_, err := outputPS(cmdRemoveLocalGroupMember + " -Group " + ps(hvAdminsGroup) + " -Member " + ps(name))
	if err != nil {
		return fmt.Errorf("failed to remove %s from Hyper-V Administrators", name)
	}
	return nil
}
