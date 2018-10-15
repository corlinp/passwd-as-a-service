package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
)

func main() {
	err := readPasswdFile()
	if err != nil {
		log.Fatal("Error reading passwd file: ", err.Error())
	}
	go watchFiles()

	e := echo.New()
	e.GET("/healthcheck", healthCheck)
	e.GET("/users", getUsers)
	e.GET("/users/:uid", getUserByUID)
	e.GET("/users/query", queryUsers)
	e.Logger.Fatal(e.Start(":1323"))
}

func init() {
	// Read passwd and groups file path from env vars
	parsePath := func(path string) string {
		absPath, err := filepath.Abs(path)
		if err != nil {
			log.Fatal("Invalid file path:", path)
		}
		return absPath
	}

	envPath := os.Getenv("PASSWD_PATH")
	if envPath != "" {
		passwdFilePath = parsePath(envPath)
	}
	envPath = os.Getenv("GROUPS_PATH")
	if envPath != "" {
		groupsFilePath = parsePath(envPath)
	}
}
