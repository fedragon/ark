package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/auth"
	"github.com/fedragon/ark/internal/db"
	"github.com/fedragon/ark/internal/server"

	"github.com/bufbuild/connect-go"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Config struct {
	DbPath      string   `split_words:"true" default:"./ark.db"`
	FileTypes   []string `split_words:"true" default:"cr2,orc,jpg,jpeg,mp4,mov,avi,mpg,mpeg,wmv"`
	ArchivePath string   `split_words:"true"`
	Server      struct {
		Address    string `split_words:"true" default:"localhost:9999"`
		SigningKey string `split_words:"true" default:"supersecret"`
	}
}

func main() {
	var cfg Config
	if err := envconfig.Process("ark", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	repo, err := db.NewSqlite3Repository(cfg.DbPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer repo.Close()

	handler := &server.Handler{
		Repo:        repo,
		FileTypes:   cfg.FileTypes,
		ArchivePath: cfg.ArchivePath,
	}

	mux := http.NewServeMux()
	interceptor, err := auth.NewInterceptor([]byte(cfg.Server.SigningKey))
	if err != nil {
		log.Fatal(err.Error())
	}
	mux.Handle(arkv1connect.NewArkApiHandler(
		handler,
		connect.WithInterceptors(interceptor),
	))

	fmt.Println("... Listening on", cfg.Server.Address)
	if err := http.ListenAndServe(
		cfg.Server.Address,
		h2c.NewHandler(mux, &http2.Server{}), // Use h2c so we can serve HTTP/2 without TLS.
	); err != nil {
		log.Fatal(err)
	}
}
