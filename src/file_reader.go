package main

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

// watchFiles uses fsnotify filesystem change notifications to keep an eye on the
//   passwd and groups files, and update the database if they change.
// Errors are non-fatal as the watch functionality isn't critical.
func watchFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Failed to initialize file watcher", err)
	}
	defer watcher.Close()
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if event.Name == passwdFilePath {
					log.Println("Passwd file modified. Reloading...")
					err := readPasswdFile()
					if err != nil {
						log.Println("Passwd file parsing error: ", err)
					}
				}
				// if event.Name == groupsFilePath {
				// 	log.Println("Groups file modified. Reloading...")
				// 	err := readGroupsFile()
				// 	if err != nil {
				// 		log.Println("Groups file parsing error: ", err)
				// 	}
				// }
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("file watching error:", err)
		}
	}
}

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
