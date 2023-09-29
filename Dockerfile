FROM golang:1.21-alpine3.18 AS builder
WORKDIR $GOPATH/src/github.com/fedragon/ark
COPY . .
RUN CGO_ENABLED=0 go build -o bin/server cmd/server/main.go

FROM alpine:3.18
RUN adduser -D ark -u 5000
USER ark
COPY --from=builder /go/src/github.com/fedragon/ark/bin/server /bin/server
EXPOSE 9999
ENTRYPOINT ["/bin/server"]
