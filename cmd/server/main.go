//go:build !windows

package main

import (
	"net/http"

	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/auth"
	"github.com/fedragon/ark/internal/db"
	"github.com/fedragon/ark/internal/server"

	"github.com/bufbuild/connect-go"
	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/go-homedir"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Config struct {
	ArchivePath string `split_words:"true" required:"true"`
	SigningKey  string `split_words:"true" required:"true"`
	Address     string `split_words:"true" default:"0.0.0.0:9999"`
	Postgres    struct {
		Address  string `default:"localhost:15432"`
		User     string `required:"true"`
		Password string `required:"true"`
		Database string `required:"true"`
	}
}

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	var cfg Config
	if err := envconfig.Process("ark_server", &cfg); err != nil {
		log.Fatal("Unable to process config", zap.Error(err))
	}

	archivePath, err := homedir.Expand(cfg.ArchivePath)
	if err != nil {
		log.Fatal("Unable to expand home dir", zap.Error(err))
	}

	repo, err := db.NewPgRepository(
		cfg.Postgres.Address,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)
	if err != nil {
		log.Fatal("Unable to initialize repository", zap.Error(err))
	}
	defer repo.Close()

	handler := &server.Handler{
		Repo:        repo,
		ArchivePath: archivePath,
	}

	mux := http.NewServeMux()
	interceptor, err := auth.NewInterceptor([]byte(cfg.SigningKey))
	if err != nil {
		log.Fatal("Unable to initialize interceptor", zap.Error(err))
	}
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle(arkv1connect.NewArkApiHandler(
		handler,
		connect.WithInterceptors(interceptor),
	))

	log.Info("... Listening on", zap.String("address", cfg.Address))
	if err := http.ListenAndServe(
		cfg.Address,
		h2c.NewHandler(mux, &http2.Server{}), // Use h2c so we can serve HTTP/2 without TLS.
	); err != nil {
		log.Fatal("Unable to start HTTP server", zap.Error(err))
	}
}
