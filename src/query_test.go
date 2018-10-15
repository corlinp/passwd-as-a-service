package main

import (
	"testing"
)

func TestUserDB(t *testing.T) {
	u1 := User{
		Name:    "bob",
		UID:     78,
		GID:     78,
		Comment: "Bob Jones",
		Home:    "/home/bob",
		Shell:   "/bin/bash",
	}

	userDB.Store(u1)

	q := map[string]interface{}{
		"uid":  78,
		"name": "bob",
	}
	if userDB.Query(q)[0] != u1 {
		t.Fail()
	}

	q["name"] = "notbob"
	if len(userDB.Query(q)) != 0 {
		t.Fail()
	}
}
