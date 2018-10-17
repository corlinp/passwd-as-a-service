package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testGroup1 = Group{
	Name:    "mygroup",
	GID:     24,
	Members: []string{"bob", "root"},
}

var testGroup2 = Group{
	Name:    "admin",
	GID:     80,
	Members: []string{"root"},
}

func TestGroupParsing(t *testing.T) {
	groups, _ := parseGroups(bytes.NewBuffer([]byte(
		`mygroup:*:24:bob,root
		admin:*:80:root`)))
	if !reflect.DeepEqual(groups[0], testGroup1) {
		t.Fail()
	}
	if !reflect.DeepEqual(groups[1], testGroup2) {
		t.Fail()
	}
}

func TestGroupDB(t *testing.T) {
	groupDB.SetGroupList(testGroup1, testGroup2)
	if len(groupDB.Query(nil)) != 2 {
		t.Fail()
	}
	q := map[string]interface{}{
		"gid":  24,
		"name": "mygroup",
	}
	if !reflect.DeepEqual(groupDB.Query(q)[0], testGroup1) {
		t.Fail()
	}
	q["name"] = "admin"
	if len(groupDB.Query(q)) != 0 {
		t.Fail()
	}
	q["gid"] = 80
	if !reflect.DeepEqual(groupDB.Query(q)[0], testGroup2) {
		t.Fail()
	}

	q = map[string]interface{}{
		"members": []string{"bob"},
	}
	assert.Len(t, groupDB.Query(q), 1)
	q = map[string]interface{}{
		"members": []string{"root"},
	}
	assert.Len(t, groupDB.Query(q), 2)
	q = map[string]interface{}{
		"members": []string{"bob", "root"},
	}
	assert.Len(t, groupDB.Query(q), 1)
	q = map[string]interface{}{
		"members": []string{"bob", "joe"},
	}
	assert.Len(t, groupDB.Query(q), 0)
}

func TestReadGroupFile(t *testing.T) {
	groupFilePath = "../sample_files/group.bad.txt"
	err := readGroupsFile()
	if err == nil {
		t.Error(err)
	}
	groupFilePath = "../sample_files/group.test.txt"
	err = readGroupsFile()
	if err != nil {
		t.Error(err)
	}
	groups := groupDB.Query(nil)
	if !reflect.DeepEqual(groups[0], testGroup1) {
		t.Fail()
	}
	if !reflect.DeepEqual(groups[1], testGroup2) {
		t.Fail()
	}
}

func TestGroupEndpoints(t *testing.T) {
	groupFilePath = "../sample_files/group.test.txt"
	assert.NoError(t, readGroupsFile())
	passwdFilePath = "../sample_files/passwd.test.txt"
	assert.NoError(t, readPasswdFile())

	code, body := mockParamRequest("/groups/24", "/groups/:gid", "gid", "24", getGroupByGID)
	assert.Equal(t, http.StatusOK, code)
	testGroup1JSON, _ := json.Marshal(testGroup1)
	assert.Equal(t, testGroup1JSON, body)

	parseGroups := func(body []byte) (groups []Group) {
		assert.NoError(t, json.Unmarshal(body, &groups))
		return groups
	}

	code, body = mockParamRequest("/users/78/groups", "/users/:uid/groups", "uid", "78", getGroupsByMember)
	assert.Equal(t, http.StatusOK, code)
	groups := parseGroups(body)
	assert.Len(t, groups, 1)
	assert.Equal(t, testGroup1, groups[0])

	code, body = mockRequest("/groups", getGroups)
	assert.Equal(t, http.StatusOK, code)
	groups = parseGroups(body)
	assert.Len(t, groups, 2)

	code, body = mockRequest("/groups/query?gid=24", queryGroups)
	assert.Equal(t, http.StatusOK, code)
	groups = parseGroups(body)
	assert.Len(t, groups, 1)
	assert.Equal(t, testGroup1, groups[0])

	code, body = mockRequest("/groups/query?member=bob&member=root", queryGroups)
	assert.Equal(t, http.StatusOK, code)
	groups = parseGroups(body)
	assert.Len(t, groups, 1)
	assert.Equal(t, testGroup1, groups[0])

	code, body = mockRequest("/groups/query?name=mygroup&gid=0", queryGroups)
	assert.Len(t, parseGroups(body), 0)
	code, body = mockRequest("/groups/query?gid=24&gid=123", queryGroups)
	assert.Equal(t, http.StatusBadRequest, code)
	code, body = mockRequest("/groups/query?gid=letter", queryGroups)
	assert.Equal(t, http.StatusBadRequest, code)
}
