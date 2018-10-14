package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"
)


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
