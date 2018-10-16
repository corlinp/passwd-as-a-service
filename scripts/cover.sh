#!/bin/bash
t="/tmp/go-cover.tmp" 
export GOPATH=$PWD
cd src
go get -t
go test -coverprofile=$t
go tool cover -html=$t