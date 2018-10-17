#!/bin/bash
export GOPATH=$PWD
cd src
go get
go build -o ../bin/passwd