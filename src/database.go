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

func init() {
	userDB = &arrayUserStorage{}
}

// arrayUserStorage is a simple implementation of UserDB that keeps all users in a slice
type arrayUserStorage struct {
	// lock will prevent us from doing a query while the DB is being updated
	lock sync.RWMutex
	db   []User
}

func (stor *arrayUserStorage) Store(users ...User) {
	stor.lock.Lock()
	defer stor.lock.Unlock()
	stor.db = append(stor.db, users...)
}

// QueryUsers finds users in the DB that match parameters given in the 'query' map
func (stor *arrayUserStorage) Query(query map[string]interface{}) (out []User) {
	stor.lock.RLock()
	defer stor.lock.RUnlock()
	// nil means get all - we copy the slice so modifications don't disturb the DB
	if query == nil {
		out = make([]User, len(stor.db))
		copy(out, stor.db)
		return
	}
	for _, user := range stor.db {
		if matchesQuery(query, user) {
			out = append(out, user)
		}
	}
	return
}

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
