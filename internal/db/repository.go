package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/fedragon/ark/migrations"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/migrate"
)

type Media struct {
	bun.BaseModel `bun:"table:media"`

	Hash       []byte     `bun:",pk" json:"hash"`
	Path       string     `json:"path"`
	CreatedAt  time.Time  `json:"created_at"`
	ImportedAt *time.Time `json:"imported_at,omitempty"`
	Err        error      `bun:"-" json:"-"`
}

type Repository interface {
	Close() error

	// Store stores a media in the database
	Store(ctx context.Context, media Media) error

	// Get returns a media from the database
	Get(ctx context.Context, hash []byte) (*Media, error)
}

type sqlite3Repository struct {
	db *bun.DB
}

func NewSqlite3Repository(dbPath string) (Repository, error) {
	if _, err := os.Stat(dbPath); err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Create(dbPath); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	sqldb, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:%s?cache=shared&mode=rw", dbPath))
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())
	migrator := migrate.NewMigrator(db, migrations.Migrations)
	if err := migrator.Init(context.Background()); err != nil {
		panic(err)
	}
	if _, err := migrator.Migrate(context.Background()); err != nil {
		panic(err)
	}

	return &sqlite3Repository{db: db}, nil
}

func (r *sqlite3Repository) Close() error {
	return r.db.Close()
}

func (r *sqlite3Repository) Get(ctx context.Context, hash []byte) (*Media, error) {
	media := make([]Media, 0)
	err := r.db.NewRaw("SELECT * FROM ? WHERE hash = ?", bun.Ident("media"), hash).Scan(ctx, &media)
	if err != nil {
		return nil, err
	}

	if len(media) == 0 {
		return nil, nil
	}

	return &media[0], err
}

func (r *sqlite3Repository) Store(ctx context.Context, media Media) error {
	now := time.Now()
	res, err := r.db.ExecContext(
		ctx,
		"INSERT INTO ? (hash, path, created_at, imported_at) VALUES (?, ?, ?, ?) ON CONFLICT(hash) DO UPDATE SET path = ?, imported_at = ?",
		bun.Ident("media"),
		// insert
		media.Hash,
		media.Path,
		media.CreatedAt,
		now,
		// update
		media.Path,
		now)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}
