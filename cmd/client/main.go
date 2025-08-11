//go:build !windows

package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/auth"
	"github.com/fedragon/ark/internal/importer"

	"connectrpc.com/connect"
	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

const (
	fromFlag = "from"
)

type Config struct {
	FileTypes  []string `split_words:"true" default:"cr2,orc,jpg,jpeg,mp4,mov,avi,mpg,mpeg,wmv"`
	SigningKey string   `split_words:"true" required:"true"`
	Server     struct {
		Address  string `split_words:"true" default:"localhost:9999"`
		Protocol string `default:"http"`
	}
}

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	var cfg Config
	if err := envconfig.Process("ark_client", &cfg); err != nil {
		log.Fatal("Unable to process config", zap.Error(err))
	}

	app := &cli.App{
		Usage:           "Imports files to the Ark server",
		UsageText:       "ark [global options]",
		Version:         "0.1.0",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     fromFlag,
				Required: true,
				Usage:    "Absolute path of the directory containing the files to be imported.",
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		source, err := homedir.Expand(c.String(fromFlag))
		if err != nil {
			return err
		}

		now := time.Now()
		defer func() {
			log.Info("Import finished", zap.Duration("elapsed_time", time.Since(now)))
		}()

		serverURL := url.URL{
			Scheme: cfg.Server.Protocol,
			Host:   cfg.Server.Address,
		}

		log.Info("Importing files", zap.String("source_path", source), zap.String("server_url", serverURL.String()))

		interceptor, err := auth.NewInterceptor([]byte(cfg.SigningKey))
		if err != nil {
			return err
		}
		imp := importer.NewImporter(
			arkv1connect.NewArkApiClient(
				http.DefaultClient,
				serverURL.String(),
				connect.WithSendGzip(),
				connect.WithInterceptors(interceptor),
			),
			cfg.FileTypes,
			log,
		)

		return imp.Import(context.Background(), source)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("unable to run application", zap.Error(err))
	}
}
