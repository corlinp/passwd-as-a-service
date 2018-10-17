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
	if err := readPasswdFile(); err != nil {
		log.Fatal("Error reading passwd file: ", err.Error())
	}
	if err := readGroupsFile(); err != nil {
		log.Fatal("Error reading groups file: ", err.Error())
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

	e.GET("/users/:uid/groups", getGroupsByMember)
	e.GET("/groups", getGroups)
	e.GET("/groups/query", queryGroups)
	e.GET("/groups/:gid", getGroupByGID)

	e.File("/", "web/index.html")
	e.File("/jquery.min.js", "web/jquery.min.js")
	e.HideBanner = true
	e.Logger.Fatal(e.Start(":80"))
}

var autoTLSDomain = ""

func init() {
	// Read passwd and groups file path from env vars or command line args
	parsePath := func(path string) string {
		absPath, err := filepath.Abs(path)
		if err != nil {
			log.Fatal("Invalid file path:", path)
		}
		return absPath
	}

	if v := os.Getenv("PASSWD_PATH"); v != "" {
		passwdFilePath = parsePath(v)
	}
	if v := os.Getenv("GROUPS_PATH"); v != "" {
		groupsFilePath = parsePath(v)
	}
	autoTLSDomain = os.Getenv("TLS_DOMAIN")

	pathPtr := flag.String("passwd-path", "/etc/passwd", "path to the passwd file to host")
	if pathPtr != nil {
		passwdFilePath = parsePath(*pathPtr)
	}
	pathPtr = flag.String("groups-path", "/etc/group", "path to the groups file to host")
	if pathPtr != nil {
		groupsFilePath = parsePath(*pathPtr)
	}
	pathPtr = flag.String("tls-domain", "", "host whitelist for automatic TLS certification")
	if pathPtr != nil {
		autoTLSDomain = *pathPtr
	}
}
