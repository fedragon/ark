package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/auth"
	"github.com/fedragon/ark/internal/importer"

	"github.com/bufbuild/connect-go"
	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
)

const (
	fromFlag = "from"
)

type Config struct {
	FileTypes []string `split_words:"true" default:"cr2,orc,jpg,jpeg,mp4,mov,avi,mpg,mpeg,wmv"`
	Server    struct {
		Address    string `split_words:"true" default:"localhost:9999"`
		Protocol   string `default:"http"`
		SigningKey string `split_words:"true" default:"supersecret"`
	}
}

func main() {
	var cfg Config
	if err := envconfig.Process("ark", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	app := &cli.App{
		Usage:           "a CLI to manage an archive of media files",
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
			fmt.Println("Elapsed time:", time.Since(now))
		}()

		serverURL := url.URL{
			Scheme: cfg.Server.Protocol,
			Host:   cfg.Server.Address,
		}

		fmt.Println("Importing files from", source, "to", serverURL.String())

		interceptor, err := auth.NewInterceptor([]byte(cfg.Server.SigningKey))
		if err != nil {
			return err
		}
		imp := &importer.Imp{
			FileTypes: cfg.FileTypes,
			Client: arkv1connect.NewArkApiClient(
				http.DefaultClient,
				serverURL.String(),
				connect.WithSendGzip(),
				connect.WithInterceptors(interceptor),
			),
		}

		return imp.Import(context.Background(), source)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
