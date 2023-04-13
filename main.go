package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fedragon/ark/internal/db"
	"github.com/fedragon/ark/internal/importer"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

const dbPath = "ark.db"

func main() {
	now := time.Now()
	var appFs = afero.NewOsFs()

	repo, err := db.NewSqlite3Repository(appFs, dbPath)
	if err != nil {
		panic(err)
	}
	defer repo.Close()

	imp := importer.Imp{
		Repo:      repo,
		Fs:        appFs,
		FileTypes: []string{".cr2", ".jpg", ".jpeg"},
	}

	source, err := homedir.Expand(os.Args[1])
	if err != nil {
		panic(err)
	}
	target, err := homedir.Expand(os.Args[2])
	if err != nil {
		panic(err)
	}

	if err := imp.Import(context.Background(), source, target); err != nil {
		panic(err)
	}

	fmt.Println("Elapsed time:", time.Since(now))
}
