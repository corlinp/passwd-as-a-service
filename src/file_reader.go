package main

import (
	"log"
	"os"
)

func readPasswdFile() error {
	passwdFile, err := os.Open(passwdFilePath)
	if err != nil {
		return err
	}
	users, err := parsePasswd(passwdFile)
	if err != nil {
		return err
	}
	userDB.SetUserList(users...)
	return nil
}

// func readGroupsFile() error {
// 	groupsFile, err := os.Open(groupsFilePath)
// 	if err != nil {
// 		return err
// 	}
// 	users, err := parseGroups(groupsFile)
// 	if err != nil {
// 		return err
// 	}
// 	groupDB.SetGroupList(users...)
// 	return nil
// }
