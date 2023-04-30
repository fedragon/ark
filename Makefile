build: generate build-client build-server

build-client:
	go build -o bin/client cmd/client/main.go

build-server:
	go build -o bin/server cmd/server/main.go

generate:
	rm -rf gen
	buf generate

build-server-nas: export GOOS=linux
build-server-nas: export GOARCH=amd64
build-server-nas:
	go build -o bin/server cmd/server/main.go

start-db:
	docker run -d -p 15432:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_DB=${POSTGRES_DB} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} postgres:15-alpine
