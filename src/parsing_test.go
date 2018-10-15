package main

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestUserParsing(t *testing.T) {
	u1 := User{
		Name:    "bob",
		UID:     78,
		GID:     78,
		Comment: "Bob Jones",
		Home:    "/home/bob",
		Shell:   "/bin/bash",
	}
	users, _ := parsePasswd(bytes.NewBuffer([]byte("bob:*:78:78:Bob Jones:/home/bob:/bin/bash")))
	if !reflect.DeepEqual(users[0], u1) {
		t.Fail()
	}
	fmt.Println(users)
}
