package importer

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/afero"
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

type Imp struct {
	Repo      db.Repository
	Fs        afero.Fs
	FileTypes []string
}

func (imp *Imp) Import(ctx context.Context, sourceDir string, targetDir string) error {
	for m := range fs.Walk(imp.Fs, sourceDir, imp.FileTypes) {
		existing, err := imp.Repo.Get(ctx, m.Hash)
		if err != nil {
			return err
		}

		if existing == nil {
			newPath, err := copyFile(m, targetDir)
			if err != nil {
				return err
			}
			m.Path = newPath

			if err := imp.Repo.Store(ctx, m); err != nil {
				return err
			}
		} else {
			if _, err := os.Stat(existing.Path); err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return err
				}

				newPath, err := copyFile(m, targetDir)
				if err != nil {
					return err
				}
				m.Path = newPath

				if err := imp.Repo.Store(ctx, m); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func copyFile(m db.Media, targetDir string) (string, error) {
	year := m.CreatedAt.Format("2006")
	month := m.CreatedAt.Format("01")
	day := m.CreatedAt.Format("02")
	ymdDir := filepath.Join(targetDir, year, month, day)

	if err := os.MkdirAll(ymdDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("unable to create archive subdirectory %v: %w", ymdDir, err)
	}

	data, err := os.Open(m.Path)
	if err != nil {
		return "", err
	}
	defer data.Close()

	newPath := filepath.Join(ymdDir, filepath.Base(m.Path))
	return newPath, atomic.WriteFile(newPath, data)
}
