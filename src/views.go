package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func getUsers(c echo.Context) error {
	return nil
}

func queryUsers(c echo.Context) error {
	return nil
}

func getUserByUID(c echo.Context) error {
	return nil
}

func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
