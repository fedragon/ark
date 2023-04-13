package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fedragon/ark/migrations"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/migrate"
)

const dbPath = "ark.db"

func main() {
	sqldb, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:%s?cache=shared&mode=rw", dbPath))
	if err != nil {
		panic(err)
	}
	defer sqldb.Close()

	db := bun.NewDB(sqldb, sqlitedialect.New())

	migrator := migrate.NewMigrator(db, migrations.Migrations)

	if err := migrator.Init(context.Background()); err != nil {
		panic(err)
	}

	if _, err := migrator.Migrate(context.Background()); err != nil {
		panic(err)
	}
}
