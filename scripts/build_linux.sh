#!/bin/bash
export GOPATH=$PWD
cd src
go get
GOOS=linux GOARCH=amd64 go build -o ../bin/passwd_linux