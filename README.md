[![CircleCI](https://circleci.com/gh/corlinp/passwd-as-a-service.svg?style=svg)](https://circleci.com/gh/corlinp/passwd-as-a-service)

# Passwd as a Service

Exposes the user and group information on a UNIX system. Written in Go with the labstack/echo framework.

## API Usage

### List Users

**GET** `/users`

Returns an array of all users.

Example Response:

```json
[
{"name": "root", "uid": 0, "gid": 0, "comment": "root", "home": "/root", "shell": "/bin/bash"},
{"name": "dwoodlins", "uid": 1001, "gid": 1001, "comment": "", "home": "/home/dwoodlins", "shell": "/bin/false"}
]
```

### Get User by UID

**GET** `/users/<uid>`

Returns a single user.

Example Response:

```json
{"name": "dwoodlins", "uid": 1001, "gid": 1001, "comment": "", "home": "/home/dwoodlins", "shell": "/bin/false"}
```

### Query Users by Field

**GET** `/users/query[?name=<nq>][&uid=<uq>][&gid=<gq>][&comment=<cq>][&home=<
hq>][&shell=<sq>]`

Queries users with exact matches to the given fields.

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

Searches all properties of a user for full and partial matches. 

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