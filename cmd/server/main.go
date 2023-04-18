package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/auth"
	"github.com/fedragon/ark/internal/db"
	"github.com/fedragon/ark/internal/server"

	"github.com/bufbuild/connect-go"
	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/go-homedir"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Config struct {
	ArchivePath string `split_words:"true" required:"true"`
	SigningKey  string `split_words:"true" required:"true"`
	Address     string `split_words:"true" default:"0.0.0.0:9999"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("ark_server", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	dbPath, err := homedir.Expand(filepath.Join(cfg.ArchivePath, "ark.db"))
	if err != nil {
		log.Fatal(err.Error())
	}
	repo, err := db.NewSqlite3Repository(dbPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer repo.Close()

	handler := &server.Handler{
		Repo:        repo,
		ArchivePath: cfg.ArchivePath,
	}

	mux := http.NewServeMux()
	interceptor, err := auth.NewInterceptor([]byte(cfg.SigningKey))
	if err != nil {
		log.Fatal(err.Error())
	}
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle(arkv1connect.NewArkApiHandler(
		handler,
		connect.WithInterceptors(interceptor),
	))

	fmt.Println("... Listening on", cfg.Address)
	if err := http.ListenAndServe(
		cfg.Address,
		h2c.NewHandler(mux, &http2.Server{}), // Use h2c so we can serve HTTP/2 without TLS.
	); err != nil {
		log.Fatal(err)
	}
}
