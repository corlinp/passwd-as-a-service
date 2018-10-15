package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/echo"
)

// Parses and validates query params into a valid query map
// Query params are in the URL like ?uid=123&name=root
func parseQueryParams(params map[string][]string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	for k, v := range params {
		if k == "member" {
			out["members"] = v
		} else if len(v) > 1 {
			return nil, fmt.Errorf("'%s' has too many query parameters", k)
		} else if k == "uid" || k == "gid" {
			intVal, err := strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf("'%s' must be an integer", k)
			}
			out[k] = intVal
		} else {
			out[k] = v
		}
	}
	if len(out) == 0 {
		return nil, errors.New("Query cannot be empty")
	}
	return out, nil
}

// Creates a map from URL parameters mimicing query params
// This allows us to use the parseQuery function for both
func paramsMap(c echo.Context) map[string][]string {
	out := make(map[string][]string)
	for _, param := range c.ParamNames() {
		out[param] = []string{c.Param(param)}
	}
	return out
}
