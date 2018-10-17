package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

var userDB UserDB
var groupDB GroupDB

func init() {
	userDB = &arrayUserStorage{}
	groupDB = &arrayGroupStorage{}
}

// arrayGroupStorage is a simple implementation of GroupDB that keeps all Groups in a slice
type arrayGroupStorage struct {
	// lock will prevent us from doing a query while the DB is being updated
	lock sync.RWMutex
	db   []Group
}

// SetGroupList stores Groups in the database - for simplicity, all Groups are set at once.
// If called again, old Group list will be rewritten
func (stor *arrayGroupStorage) SetGroupList(Groups ...Group) {
	stor.lock.Lock()
	defer stor.lock.Unlock()
	stor.db = Groups
	// stor.db = append(stor.db, Groups...)
}

// QueryGroups finds Groups in the DB that match parameters given in the 'query' map
func (stor *arrayGroupStorage) Query(query map[string]interface{}) (out []Group) {
	stor.lock.RLock()
	defer stor.lock.RUnlock()
	// nil means get all - we copy the slice so modifications don't disturb the DB
	if query == nil {
		out = make([]Group, len(stor.db))
		copy(out, stor.db)
		return
	}
	for _, Group := range stor.db {
		if matchesQuery(query, Group) {
			out = append(out, Group)
		}
	}
	return
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
		field := vals.Field(i)
		fieldName := vals.Type().Field(i).Tag.Get("json")
		if queryVal, ok := query[fieldName]; ok {
			// If we run into a slice, we make sure that all values in the query slice exist in the candidate slice
			// TODO: make this cleaner and not O(N^2) (sort of). We just hope members lists are short for now.
			if field.Kind() == reflect.Slice {
				cfield := reflect.ValueOf(queryVal)
			candidateLoop:
				for i := 0; i < cfield.Len(); i++ {
					for j := 0; j < field.Len(); j++ {
						if field.Index(j).Interface() == cfield.Index(i).Interface() {
							continue candidateLoop
						}
					}
					return false
				}
			} else if queryVal != field.Interface() {
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
