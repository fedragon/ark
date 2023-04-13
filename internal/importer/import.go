package importer

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fedragon/ark/internal/db"
	"github.com/fedragon/ark/internal/fs"

	"github.com/natefinch/atomic"
)

type Importer interface {
	// Import imports all files in sourceDir into targetDir, skipping duplicates
	Import(ctx context.Context, sourceDir string, targetDir string) error
}

type ConcurrentImporter struct {
	Repo       db.Repository
	FileTypes  []string
	NumWorkers int
}

func (ci *ConcurrentImporter) Import(ctx context.Context, sourceDir string, targetDir string) error {
	media := fs.Walk(sourceDir, ci.FileTypes)

	for m := range media {
		if err := ci.Repo.Store(m); err != nil {
			return err
		}

		creationTime := m.ModTime

		year := creationTime.Format("2006")
		month := creationTime.Format("01")
		day := creationTime.Format("02")
		ymdDir := filepath.Join(targetDir, year, month, day)

		if err := os.MkdirAll(ymdDir, os.ModePerm); err != nil {
			return fmt.Errorf("unable to create archive subdirectory %v: %w", ymdDir, err)
		}

		data, err := os.Open(m.Path)
		if err != nil {
			log.Fatal(err.Error())
		}

		if err := atomic.WriteFile(filepath.Join(ymdDir, filepath.Base(m.Path)), data); err != nil {
			// TODO in case of error, I need to delete the table row. Holding a transaction open
			// doesn't sound sensible since copying big files could take several seconds.
			// Two options to consider: soft deletes (natively supported by Bun, although I should then probably
			// use it as a "real" ORM) or adding a "stored_at" column and only deleting the row if its value
			// the same as that of a "stored_at" variable that I would define before calling Store().
			data.Close()
			return err
		}

		data.Close()
	}

	return nil
}
