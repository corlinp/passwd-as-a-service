cover:
	./scripts/cover.sh

test:
	cd src && go test

build:
	./scripts/build.sh

linux:
	./scripts/build_linux.sh

run:
	./bin/passwd

start: build run

connect:
	ssh root@138.68.13.212

provision: linux
	scp -r bin/passwd_linux root@138.68.13.212:/root/pwaas
	scp -r web root@138.68.13.212:/root/web
	scp -r sample_files/ root@138.68.13.212:/root/
	ssh root@138.68.13.212 ./pwaas --port 80 --passwd-file sample_files/passwd.txt --group-file sample_files/group.txt
