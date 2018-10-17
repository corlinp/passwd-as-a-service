package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

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

var passwdTestFile = "../sample_files/passwd.test.txt"

func TestReadPasswdFile(t *testing.T) {
	passwdFilePath = "../sample_files/passwd.bad.txt"
	err := readPasswdFile()
	if err == nil {
		t.Error(err)
	}
	passwdFilePath = passwdTestFile
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
	passwdFilePath = passwdTestFile
	assert.NoError(t, readPasswdFile())
	code, body := mockParamRequest("/users/78", "/users/:uid", "uid", "78", getUserByUID)
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

	code, body = mockRequest("/users/query?name=bob&uid=0", queryUsers)
	assert.Len(t, parseUsers(body), 0)
	code, body = mockRequest("/users/query?uid=78&uid=123", queryUsers)
	assert.Equal(t, http.StatusBadRequest, code)
	code, body = mockRequest("/users/query?uid=letter", queryUsers)
	assert.Equal(t, http.StatusBadRequest, code)

	code, body = mockRequest("/users/search?q=bob", searchUsers)
	assert.Equal(t, http.StatusOK, code)
	users = parseUsers(body)
	assert.Len(t, users, 1)
	assert.Equal(t, testUser1, users[0])
}

func mockRequest(endpoint string, handler func(c echo.Context) error) (code int, body []byte) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, endpoint, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := handler(c)
	if err != nil {
		fmt.Println(err)
	}
	code = rec.Code
	body = rec.Body.Bytes()
	return
}

func mockParamRequest(endpoint, path, paramname, paramvalue string, handler func(c echo.Context) error) (code int, body []byte) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, endpoint, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(path)
	c.SetParamNames(paramname)
	c.SetParamValues(paramvalue)
	err := handler(c)
	if err != nil {
		fmt.Println(err)
	}
	code = rec.Code
	body = rec.Body.Bytes()
	return
}

func TestFileReloading(t *testing.T) {
	reloadFile := passwdTestFile + ".reload"

	passwdText, err := ioutil.ReadFile(passwdTestFile)
	assert.NoError(t, err)

	err = ioutil.WriteFile(reloadFile, passwdText, os.ModePerm)
	assert.NoError(t, err)
	defer os.Remove(reloadFile)

	passwdFilePath = reloadFile
	err = readPasswdFile()
	assert.NoError(t, err)
	assert.Len(t, userDB.Query(nil), 2)

	// Start the file watcher, give it some time to initialize and read
	go watchFiles()
	time.Sleep(time.Millisecond * 300)

	passwdText = []byte(string(passwdText) + "\nnewuser:*:1:2:Just Added:/usr:/bin/bash")
	err = ioutil.WriteFile(reloadFile, passwdText, os.ModePerm)
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 300)
	assert.Len(t, userDB.Query(nil), 3)
}
