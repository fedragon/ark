//go:build !windows

package main

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/auth"
	"github.com/fedragon/ark/internal/db"
	"github.com/fedragon/ark/internal/server"

	"connectrpc.com/connect"
	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/go-homedir"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Config struct {
	ArchivePath string `split_words:"true" required:"true"`
	SigningKey  string `split_words:"true" required:"true"`
	Address     string `split_words:"true" default:"0.0.0.0:9999"`
	Redis       struct {
		Address  string `default:"localhost:6379"`
		Password string `default:""`
		Database int    `default:"0"`
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

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})
	repo := db.NewRepository(client)
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

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Fatal("Unable to listen", zap.Error(err))
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdown
		log.Info("... Shutting down")
		close(shutdown)
		listener.Close()
		repo.Close()
		log.Sync()
		os.Exit(0)
	}()

	log.Info("... Listening on", zap.String("address", cfg.Address))
	if err := http.Serve(
		listener,
		h2c.NewHandler(mux, &http2.Server{}), // Use h2c so we can serve HTTP/2 without TLS.
	); err != nil {
		log.Fatal("Unable to start HTTP server", zap.Error(err))
	}
}
