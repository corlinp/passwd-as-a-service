[![CircleCI](https://circleci.com/gh/corlinp/passwd-as-a-service.svg?style=svg)](https://circleci.com/gh/corlinp/passwd-as-a-service)

# Passwd as a Service

API to expose the user and group information on a UNIX system. Written in Go with the labstack/echo framework.

## Features

* User and Group enumeration and queries
* Text-based searches for users
* Live refresh of database when a passwd or group file changes
* Graphical front end for searching users
* Unit testing and code coverage maps
* CircleCI integration to run and report on unit tests
* Automatic TLS certification with Let's Encrypt

## Running PwaaS

### Using Hosted Service

Visit [passwd.corlin.io](http://passwd.corlin.io/) to test out PwaaS with linux and mac sample data.

### Pre-Compiled Binaries

Visit the [releases page](https://github.com/corlinp/passwd-as-a-service/releases/) for instructions.

### From Source

First: Make sure you have the [latest version of Go installed](https://golang.org/doc/install).

Clone the repo and run `make start`. A binary will be built in the `bin/` folder and executed.

`make cover` will run unit tests and open a browser page with coverage reports.

---

```
Usage of ./pwaas:
  -group-file   string
        path to the groups file to host (default "/etc/group")
  -passwd-file  string
        path to the passwd file to host (default "/etc/passwd")
  -port         int
        port to run server on (default 8000)
  -tls
        enable automatic TLS certification (default false)
```

## API Usage

### List Users

**GET** `/users`

Returns an array of all users. [Try it](http://passwd.corlin.io/users?pretty)

Example Response:

```json
[
{"name": "root", "uid": 0, "gid": 0, "comment": "root", "home": "/root", "shell": "/bin/bash"},
{"name": "dwoodlins", "uid": 1001, "gid": 1001, "comment": "", "home": "/home/dwoodlins", "shell": "/bin/false"}
]
```

### Get User by UID

**GET** `/users/<uid>`

Returns a single user. [Try it](http://passwd.corlin.io/users/33?pretty)

Example Response:

```json
{"name": "dwoodlins", "uid": 1001, "gid": 1001, "comment": "", "home": "/home/dwoodlins", "shell": "/bin/false"}
```

### Query Users by Field

**GET** `/users/query[?name=<nq>][&uid=<uq>][&gid=<gq>][&comment=<cq>][&home=<
hq>][&shell=<sq>]`

Queries users with exact matches to the given fields. [Try it](http://passwd.corlin.io/users/query?shell=%2Fbin%2Ffalse&pretty)

Example Query:
```
GET /users/query?shell=%2Fbin%2Ffalse
```

Example Response:
```json
[
{"name": "dwoodlins", "uid": 1001, "gid": 1001, "comment": "", "home": "/home/dwoodlins", "shell": "/bin/false"}
]
```

### Search Users by Text Matching

**GET** `/users/search?q=<term>`

Searches all properties of a user for full and partial matches, returns up to 3 results. [Try it](http://passwd.corlin.io/users/search?q=serv&pretty)

Example Query:
```
GET /users/search?q=dwoo
```

Example Response:
```json
[
{"name": "dwoodlins", "uid": 1001, "gid": 1001, "comment": "", "home": "/home/dwoodlins", "shell": "/bin/false"}
]
```

### Get User's Groups

**GET** `/users/<uid>/groups`

Returns all groups for a given user. [Try it](http://passwd.corlin.io/users/0/groups?pretty)

Example Query:
```
GET /users/1001/groups
```

Example Response:
```json
[
{"name": "docker", "gid": 1002, "members": ["dwoodlins"]}
]
```

### List Groups

**GET** `/groups`

Returns an array of all groups. [Try it](http://passwd.corlin.io/groups?pretty)

Example Response:

```json
[
{"name": "_analyticsusers", "gid": 250, "members":["_analyticsd" ,"_networkd", "_timed"]},
{"name": "docker", "gid": 1002, "members": []}
]
```

### Get Group by GID

**GET** `/groups/<gid>`

Returns a single group. [Try it](http://passwd.corlin.io/groups/29?pretty)

Example Response:

```json
{"name": "docker", "gid": 1002, "members": ["dwoodlins"]}
```

### Query Groups by Field

**GET** `/groups/query[?name=<nq>][&gid=<gq>][&member=<mq1>[&member=<mq2>][&...]]`

Queries groups with exact matches to the name and GID field, and that contain the members listed [Try it](http://passwd.corlin.io/groups/query?/groups/query?member=_analyticsd&member=_networkd&pretty)

Example Query:
```
GET /groups/query?member=_analyticsd&member=_networkd
```

Example Response:
```json
[
{"name": "_analyticsusers", "gid": 250, "members":["_analyticsd", "_networkd", "_timed"]}
]
```