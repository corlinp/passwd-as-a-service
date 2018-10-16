package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

// UserDB is an interface to store and query Users
// Using an interface allows us to easily add new storage backends
type UserDB interface {
	SetUserList(...User)
	Query(map[string]interface{}) []User
	Search(term string) []User
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

// SetUserList stores users in the database - for simplicity, all users are set at once.
// If called again, old user list will be rewritten
func (stor *arrayUserStorage) SetUserList(users ...User) {
	stor.lock.Lock()
	defer stor.lock.Unlock()
	stor.db = users
	// stor.db = append(stor.db, users...)
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

func (stor *arrayUserStorage) Search(term string) (out []User) {
	stor.lock.RLock()
	defer stor.lock.RUnlock()
	// nil means get all - we copy the slice so modifications don't disturb the DB
	if term == "" {
		return
	}
	var results SearchResults
	term = strings.ToLower(term)
	for _, user := range stor.db {
		rel := matchesTerm(term, user)
		result := SearchResult{
			user:      user,
			relevance: rel,
		}
		results = append(results, result)
	}
	sort.Sort(results)
	// Return the top 3 results with any relevance
	for i, result := range results {
		if i < 3 && result.relevance > 0 {
			out = append(out, result.user)
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
		if queryVal, ok := query[fieldName]; ok {
			if queryVal != vals.Field(i).Interface() {
				return false
			}
		}
	}
	return true
}

type SearchResult struct {
	user      User
	relevance int
}

type SearchResults []SearchResult

func (p SearchResults) Len() int           { return len(p) }
func (p SearchResults) Less(i, j int) bool { return p[i].relevance > p[j].relevance }
func (p SearchResults) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Returns true if any stringified value in the candidate contains 'term'
func matchesTerm(term string, candidate interface{}) (relevance int) {
	vals := reflect.ValueOf(candidate)
	for i := 0; i < vals.NumField(); i++ {
		stringVal := strings.ToLower(fmt.Sprint(vals.Field(i).Interface()))
		// A basic method of ranking search relevance
		if stringVal == term {
			relevance += 5
		} else if strings.HasPrefix(stringVal, term) {
			relevance += 3
		} else if strings.Contains(stringVal, term) {
			relevance++
		}
	}
	return
}
