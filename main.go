package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fedragon/ark/internal/db"
	"github.com/fedragon/ark/internal/importer"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

const (
	dbPathFlag    = "db"
	sourceFlag    = "source"
	archiveFlag   = "archive"
	fileTypesFlag = "file-types"
)

func main() {

	app := &cli.App{
		Usage:           "a CLI to manage an archive of media files",
		UsageText:       "ark [global options]",
		Version:         "0.1.0",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:  dbPathFlag,
				Value: "./ark.db",
				Usage: "Path to the SQLite3 database file. Will be created if it doesn't exist.",
			},
			&cli.PathFlag{
				Name:     sourceFlag,
				Required: true,
				Usage:    "Absolute path of the directory containing the files to be imported.",
			},
			&cli.PathFlag{
				Name:     archiveFlag,
				Required: true,
				Usage:    "Absolute path of the archive directory.",
			},
			&cli.StringSliceFlag{
				Name:  fileTypesFlag,
				Value: cli.NewStringSlice(".cr2", ".jpg", ".jpeg", ".mov", ".mp4", ".orf"),
				Usage: "File types to archive.",
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		dbPath, err := homedir.Expand(c.String(dbPathFlag))
		if err != nil {
			return err
		}
		source, err := homedir.Expand(c.String(sourceFlag))
		if err != nil {
			return err
		}
		archive, err := homedir.Expand(c.String(archiveFlag))
		if err != nil {
			return err
		}
		fileTypes := c.StringSlice(fileTypesFlag)

		now := time.Now()
		defer func() { fmt.Println("Elapsed time:", time.Since(now)) }()

		var appFs = afero.NewOsFs()

		repo, err := db.NewSqlite3Repository(appFs, dbPath)
		if err != nil {
			return err
		}
		defer repo.Close()

		imp := importer.Imp{
			Repo:      repo,
			Fs:        appFs,
			FileTypes: fileTypes,
		}

		return imp.Import(context.Background(), source, archive)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
