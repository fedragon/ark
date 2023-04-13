package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/fedragon/ark/migrations"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/migrate"
)

const dbPath = "ark.db"

func main() {
	db, err := connect(dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	migrator := migrate.NewMigrator(db, migrations.Migrations)
	if err := migrator.Init(context.Background()); err != nil {
		panic(err)
	}
	if _, err := migrator.Migrate(context.Background()); err != nil {
		panic(err)
	}
}

func connect(dbPath string) (*bun.DB, error) {
	if _, err := os.Stat(dbPath); err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Create(dbPath); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	sqldb, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:%s?cache=shared&mode=rw", dbPath))
	if err != nil {
		return nil, err
	}

	return bun.NewDB(sqldb, sqlitedialect.New()), nil
}
