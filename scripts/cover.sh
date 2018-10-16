#!/bin/bash
t="/tmp/go-cover.tmp" 
export GOPATH=$PWD
cd src
go get github.com/stretchr/testify/assert
go get
go test -coverprofile=$t
go tool cover -html=$t