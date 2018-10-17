cover:
	./scripts/cover.sh

test:
	cd src && go test

build:
	./scripts/build.sh

linux:
	cd src && go get && GOOS=linux GOARCH=amd64 go build -o ../bin/passwd_linux

run:
	./bin/passwd

start: build run
