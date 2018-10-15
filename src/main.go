package main

import (
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/healthcheck", healthCheck)
	e.GET("/users", getUsers)
	e.GET("/users/:id", getUserByUID)
	e.GET("/users/query", queryUsers)
	e.Logger.Fatal(e.Start(":1323"))
}
