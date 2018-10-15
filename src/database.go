package main

import (
	"reflect"
	"sync"
)

// UserDB is an interface to store and query Users
// Using an interface allows us to easily add new storage backends
type UserDB interface {
	Store(...User)
	Query(map[string]interface{}) []User
}

// GroupDB is an interface to store and query Groups
type GroupDB interface {
	Store(...Group)
	Query(map[string]interface{}) []Group
}

var userDB UserDB
var groupDB GroupDB


// Returns true if values in query are equal to corresponding JSON values in candidate
func matchesQuery(query map[string]interface{}, candidate interface{}) bool {
	vals := reflect.ValueOf(candidate)
	for i := 0; i < vals.NumField(); i++ {
		field := vals.Type().Field(i)
		fieldName := field.Tag.Get("json")
		// fmt.Println("matching", fieldName)
		if queryVal, ok := query[fieldName]; ok {
			// fmt.Println("\t", queryVal, vals.Field(i).Interface())
			if queryVal != vals.Field(i).Interface() {
				return false
			}
		}
	}
	return true
}
