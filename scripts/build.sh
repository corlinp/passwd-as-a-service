#!/bin/bash
export GOPATH=$PWD
cd src
go get
go build -o passwdaas
mv passwdaas ../bin/