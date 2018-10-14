package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"
)

/*
passwd files contain lines of colon-delimited user info. Example:

bob:*:78:78:Bob Jones:/home/bob:/bin/bash

_mailman	username
* 			password - in this case nil
78			user ID (UID)
78			group ID (GID)
Bob Jones	user id info - a comment field
/home/bob	home directory for the user
/bin/bash	user's default shell (or a command)
*/
func parsePasswd(reader io.Reader) (users []User, err error) {
	// Read file line-by-line
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		// skip empty and commented lines
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			line = strings.TrimSpace(line)
			fields := strings.Split(line, ":")
			if len(fields) != 7 {
				err = errors.New("passwd parse error: incorrect field count")
				log.Println(err.Error(), line)
				return
			}
			var uid, gid int
			uid, err = strconv.Atoi(fields[2])
			if err != nil {
				log.Println("passwd parse error: uid", line)
				return
			}
			gid, err = strconv.Atoi(fields[3])
			if err != nil {
				log.Println("passwd parse error: gid", line)
				return
			}
			user := User{
				Name:    fields[0],
				UID:     uid,
				GID:     gid,
				Comment: fields[4],
				Home:    fields[5],
				Shell:   fields[6],
			}
			users = append(users, user)
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}
	return
}

/*
group files contain lines of colon-delimited group info. Example:

admin:*:80:root,other

admin		username
* 			password - in this case nil
80			group ID (GID)
root,other	comma-delimited list of group members
*/
func parseGroup(reader io.Reader) (groups []Group, err error) {
	// Read file line-by-line
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		// skip empty and commented lines
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			line = strings.TrimSpace(line)
			fields := strings.Split(line, ":")
			if len(fields) != 4 {
				err = errors.New("groups parse error: incorrect field count")
				log.Println(err.Error(), line)
				return
			}
			var gid int
			gid, err = strconv.Atoi(fields[2])
			if err != nil {
				log.Println("groups parse error: gid", line)
				return
			}
			group := Group{
				Name:    fields[0],
				GID:     gid,
				Members: strings.Split(fields[3], ","),
			}
			groups = append(groups, group)
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}
	return
}
