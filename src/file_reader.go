package main

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

var passwdFilePath string
var groupFilePath string

func readPasswdFile() error {
	passwdFile, err := os.Open(passwdFilePath)
	if err != nil {
		return err
	}
	defer passwdFile.Close()
	users, err := parsePasswd(passwdFile)
	if err != nil {
		return err
	}
	userDB.SetUserList(users...)
	log.Println("Parsed passwd file:", passwdFilePath)
	return nil
}

func readGroupFile() error {
	groupsFile, err := os.Open(groupFilePath)
	if err != nil {
		return err
	}
	defer groupsFile.Close()
	users, err := parseGroups(groupsFile)
	if err != nil {
		return err
	}
	groupDB.SetGroupList(users...)
	log.Println("Parsed groups file:", groupFilePath)
	return nil
}

// watchFiles uses fsnotify filesystem change notifications to keep an eye on the
//   passwd and groups files, and update the database if they change.
// Errors are non-fatal as the watch functionality isn't critical.
func watchFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Failed to initialize file watcher", err)
	}
	watcher.Add(passwdFilePath)
	watcher.Add(groupFilePath)
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
				if event.Name == groupFilePath {
					log.Println("Groups file modified. Reloading...")
					err := readGroupFile()
					if err != nil {
						log.Println("Groups file parsing error: ", err)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("file watching error:", err)
		}
	}
}
