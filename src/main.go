package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// Read the passwd and group files
	if err := readPasswdFile(); err != nil {
		log.Fatal("Error reading passwd file: ", err.Error())
	}
	if err := readGroupFile(); err != nil {
		log.Fatal("Error reading groups file: ", err.Error())
	}
	// Watch the files for changes in another goroutine, update the db if they change
	go watchFiles()

	// Build the echo server objects with endpoints, middleware, and TLS
	e := echo.New()
	if autoTLS {
		e.Pre(middleware.HTTPSRedirect())
		e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	}

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
	if autoTLS {
		e.Logger.Fatal(e.StartAutoTLS(":" + fmt.Sprint(port)))
	}
	e.Logger.Fatal(e.Start(":" + fmt.Sprint(port)))
}

var autoTLS bool
var port int

func init() {
	// Read passwd and group file path from command line args
	parsePath := func(path string) string {
		absPath, err := filepath.Abs(path)
		if err != nil {
			log.Fatal("Invalid file path:", path)
		}
		return absPath
	}

	passwdPathPtr := flag.String("passwd-file", "/etc/passwd", "path to the passwd file to host")
	groupsPathPtr := flag.String("group-file", "/etc/group", "path to the groups file to host")
	tlsPtr := flag.Bool("tls", false, "enable automatic TLS certification")
	portPtr := flag.Int("port", 8000, "port to run server on")
	flag.Parse()

	passwdFilePath = parsePath(*passwdPathPtr)
	groupFilePath = parsePath(*groupsPathPtr)
	autoTLS = *tlsPtr
	port = *portPtr
}
