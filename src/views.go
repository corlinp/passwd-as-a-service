package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func getUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, userDB.Query(nil))
}

func queryUsers(c echo.Context) error {
	query, err := parseQueryParams(c.QueryParams())
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, userDB.Query(query))
}

func searchUsers(c echo.Context) error {
	term := c.QueryParam("q")
	return c.JSON(http.StatusOK, userDB.Search(term))
}

func getUserByUID(c echo.Context) error {
	query, err := parseQueryParams(paramsMap(c))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	result := userDB.Query(query)
	if len(result) == 0 {
		return c.String(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, result[0])
}

func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
