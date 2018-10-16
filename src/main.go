package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	err := readPasswdFile()
	if err != nil {
		log.Fatal("Error reading passwd file: ", err.Error())
	}
	go watchFiles()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/healthcheck", healthCheck)
	e.GET("/users", getUsers)
	e.GET("/users/:uid", getUserByUID)
	e.GET("/users/query", queryUsers)
	e.GET("/users/search", searchUsers)
	e.File("/", "web/index.html")
	e.File("/jquery.min.js", "web/jquery.min.js")
	e.HideBanner = true
	e.Logger.Fatal(e.Start(":1323"))
}

func init() {
	// Read passwd and groups file path from env vars or command line args
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

	pathPtr := flag.String("passwd-path", "/etc/passwd", "path to the passwd file to host")
	if pathPtr != nil {
		passwdFilePath = parsePath(*pathPtr)
	}
	pathPtr = flag.String("groups-path", "/etc/groups", "path to the groups file to host")
	if pathPtr != nil {
		groupsFilePath = parsePath(*pathPtr)
	}
}
