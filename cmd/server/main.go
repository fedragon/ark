package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/db"
	"github.com/fedragon/ark/internal/server"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Config struct {
	DbPath        string   `split_words:"true" default:"./arkv1.db"`
	FileTypes     []string `split_words:"true" default:"cr2,orc,jpg,jpeg,mp4,mov,avi,mpg,mpeg,wmv"`
	ServerAddress string   `split_words:"true" default:"localhost:8080"`
	ArchivePath   string   `split_words:"true"`
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

	srv := &server.Ark{
		Repo:        repo,
		FileTypes:   cfg.FileTypes,
		ArchivePath: cfg.ArchivePath,
	}

	mux := http.NewServeMux()
	mux.Handle(arkv1connect.NewArkApiHandler(srv))

	fmt.Println("... Listening on", cfg.ServerAddress)
	if err := http.ListenAndServe(
		cfg.ServerAddress,
		h2c.NewHandler(mux, &http2.Server{}), // Use h2c so we can serve HTTP/2 without TLS.
	); err != nil {
		log.Fatal(err)
	}
}
