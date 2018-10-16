package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
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

func TestUserEndpoints(t *testing.T) {
	passwdFilePath = "../test_files/passwd.txt"
	assert.NoError(t, readPasswdFile())
	code, body := mockRequest("/users/78", getUserByUID, "78")
	assert.Equal(t, http.StatusOK, code)
	testUser1JSON, _ := json.Marshal(testUser1)
	assert.Equal(t, testUser1JSON, body)

	parseUsers := func(body []byte) (users []User) {
		assert.NoError(t, json.Unmarshal(body, &users))
		return users
	}

	code, body = mockRequest("/users", getUsers)
	assert.Equal(t, http.StatusOK, code)
	users := parseUsers(body)
	assert.Len(t, users, 2)

	code, body = mockRequest("/users/query?uid=78", queryUsers)
	assert.Equal(t, http.StatusOK, code)
	users = parseUsers(body)
	assert.Len(t, users, 1)
	assert.Equal(t, testUser1, users[0])

	code, body = mockRequest("/users/search?q=bob", searchUsers)
	assert.Equal(t, http.StatusOK, code)
	users = parseUsers(body)
	assert.Len(t, users, 1)
	assert.Equal(t, testUser1, users[0])
}

func mockRequest(endpoint string, handler func(c echo.Context) error, uids ...string) (code int, body []byte) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, endpoint, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if len(uids) > 0 {
		c.SetPath("/users/:uid")
		c.SetParamNames("uid")
		c.SetParamValues(uids[0])
	}
	err := handler(c)
	if err != nil {
		fmt.Println(err)
	}
	code = rec.Code
	body = rec.Body.Bytes()
	return
}
