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

