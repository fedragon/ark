.PHONY: generate
generate:
	rm -rf gen
	buf generate

.PHONY: build
build: generate build-client build-server

.PHONY: build-client
build-client:
	go build -o bin/client cmd/client/main.go

.PHONY: build-server
build-server:
	go build -o bin/server cmd/server/main.go

.PHONY: build-server-nas
build-server-nas: export GOOS=linux
build-server-nas: export GOARCH=amd64
build-server-nas:
	go build -o bin/server cmd/server/main.go

.PHONY: test
test:
	go test -race -count=1 ./...