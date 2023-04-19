build: generate build-client build-server

build-client:
	go build -o bin/client cmd/client/main.go

build-server:
	go build -o bin/server cmd/server/main.go

generate:
	rm -r gen
	buf generate

build-server-nas: export GOOS=linux
build-server-nas: export GOARCH=amd64
build-server-nas:
	go build -o bin/server cmd/server/main.go
