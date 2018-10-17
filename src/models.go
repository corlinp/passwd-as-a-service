package main

// User represents a UNIX user in a passwd file
type User struct {
	Name    string `json:"name"`
	UID     int    `json:"uid"`
	GID     int    `json:"gid"`
	Comment string `json:"comment"`
	Home    string `json:"home"`
	Shell   string `json:"shell"`
}

// Group represents a UNIX group in a groups file
type Group struct {
	Name    string   `json:"name"`
	GID     int      `json:"gid"`
	Members []string `json:"members"`
}

// UserDB is an interface to store and query Users
// Using an interface allows us to easily add new storage backends
type UserDB interface {
	SetUserList(...User)
	Query(map[string]interface{}) []User
	Search(term string) []User
}

// GroupDB is an interface to store and query Groups
type GroupDB interface {
	SetGroupList(...Group)
	Query(map[string]interface{}) []Group
}
