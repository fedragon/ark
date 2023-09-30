FROM golang:1.21-alpine3.18 AS builder
WORKDIR $GOPATH/src/github.com/fedragon/ark
COPY . .
RUN CGO_ENABLED=0 go build -o bin/server cmd/server/main.go

FROM alpine:3.18
RUN addgroup ark && adduser -D ark -G ark
RUN mkdir -p /ark/tmp && chown -R ark:ark /ark
USER ark
COPY --from=builder /go/src/github.com/fedragon/ark/bin/server /bin/server
EXPOSE 9999
ENTRYPOINT ["/bin/server"]
