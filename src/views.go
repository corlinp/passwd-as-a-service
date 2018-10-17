package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

/***** USER ENDPOINTS *****/

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

/***** GROUP ENDPOINTS *****/

func getGroups(c echo.Context) error {
	return c.JSON(http.StatusOK, groupDB.Query(nil))
}

func queryGroups(c echo.Context) error {
	query, err := parseQueryParams(c.QueryParams())
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, groupDB.Query(query))
}

func getGroupsByMember(c echo.Context) error {
	query, err := parseQueryParams(paramsMap(c))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	memberResults := userDB.Query(query)
	if len(memberResults) == 0 {
		return c.String(http.StatusNotFound, "User not found")
	}
	query["members"] = []string{memberResults[0].Name}
	return c.JSON(http.StatusOK, groupDB.Query(query))
}

func getGroupByGID(c echo.Context) error {
	query, err := parseQueryParams(paramsMap(c))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	result := groupDB.Query(query)
	if len(result) == 0 {
		return c.String(http.StatusNotFound, "Group not found")
	}
	return c.JSON(http.StatusOK, result[0])
}
