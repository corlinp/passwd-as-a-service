package main

import (
	"bytes"
	"reflect"
	"testing"
)

var testUser1 = User{
	Name:    "bob",
	UID:     78,
	GID:     78,
	Comment: "Bob Jones",
	Home:    "/home/bob",
	Shell:   "/bin/bash",
}

var testUser2 = User{
	Name:    "root",
	UID:     0,
	GID:     0,
	Comment: "Root User",
	Home:    "/root",
	Shell:   "/bin/bash",
}

func TestUserParsing(t *testing.T) {
	users, _ := parsePasswd(bytes.NewBuffer([]byte(
		`bob:*:78:78:Bob Jones:/home/bob:/bin/bash
		root:*:0:0:Root User:/root:/bin/bash`)))
	if !reflect.DeepEqual(users[0], testUser1) {
		t.Fail()
	}
	if !reflect.DeepEqual(users[1], testUser2) {
		t.Fail()
	}
}

func TestUserDB(t *testing.T) {
	userDB.SetUserList(testUser1, testUser2)
	if len(userDB.Query(nil)) != 2 {
		t.Fail()
	}
	q := map[string]interface{}{
		"uid":  78,
		"name": "bob",
	}
	if userDB.Query(q)[0] != testUser1 {
		t.Fail()
	}
	q["name"] = "root"
	if len(userDB.Query(q)) != 0 {
		t.Fail()
	}
	q["uid"] = 0
	if userDB.Query(q)[0] != testUser2 {
		t.Fail()
	}

	q = map[string]interface{}{
		"root": "/bin/bash",
	}
	if len(userDB.Query(q)) != 2 {
		t.Fail()
	}

	if userDB.Search("bob")[0] != testUser1 {
		t.Fail()
	}
}

func TestReadFile(t *testing.T) {
	passwdFilePath = "../test_files/passwd.bad.txt"
	err := readPasswdFile()
	if err == nil {
		t.Error(err)
	}
	passwdFilePath = "../test_files/passwd.txt"
	err = readPasswdFile()
	if err != nil {
		t.Error(err)
	}
	users := userDB.Query(nil)
	if !reflect.DeepEqual(users[0], testUser1) {
		t.Fail()
	}
	if !reflect.DeepEqual(users[1], testUser2) {
		t.Fail()
	}
}
