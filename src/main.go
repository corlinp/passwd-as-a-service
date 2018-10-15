package main

import (
	"log"

	"github.com/labstack/echo"
)

var passwdFilePath, groupsFilePath string

func main() {
	err := readPasswdFile()
	if err != nil {
		log.Fatal("Error reading passwd file: ", err.Error())
	}
	go watchFiles()

	e := echo.New()
	e.GET("/healthcheck", healthCheck)
	e.GET("/users", getUsers)
	e.GET("/users/:id", getUserByUID)
	e.GET("/users/query", queryUsers)
	e.Logger.Fatal(e.Start(":1323"))
}

func init() {
	passwdFilePath = "/etc/passwd"
	groupsFilePath = "/etc/groups"
}
