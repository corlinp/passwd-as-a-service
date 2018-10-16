cover:
	./scripts/test.sh

test:
	cd src && go test

build:
	./scripts/build.sh

run:
	./bin/passwdaas

start: build run
