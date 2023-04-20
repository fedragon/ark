package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fedragon/ark/internal/metrics"
	"github.com/fedragon/ark/migrations"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
)

type pgRepository struct {
	db *bun.DB
}

func NewPgRepository(address, user, password, database string) (Repository, error) {
	conn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(address),
		pgdriver.WithUser(user),
		pgdriver.WithInsecure(true),
		pgdriver.WithPassword(password),
		pgdriver.WithDatabase(database),
		pgdriver.WithTimeout(5*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(5*time.Second),
		pgdriver.WithWriteTimeout(5*time.Second),
	)

	sqldb := sql.OpenDB(conn)
	db := bun.NewDB(sqldb, pgdialect.New())
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	migrator := migrate.NewMigrator(db, migrations.Migrations)
	if err := migrator.Init(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize migrations: %w", err)
	}
	if _, err := migrator.Migrate(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return &pgRepository{db: db}, nil
}

func (r *pgRepository) Close() error {
	return r.db.Close()
}

func (r *pgRepository) Get(ctx context.Context, hash []byte) (*Media, error) {
	now := time.Now()
	defer func() {
		metrics.GetDurationMs.Observe(float64(time.Since(now).Milliseconds()))
	}()

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

func (r *pgRepository) Store(ctx context.Context, media Media) error {
	now := time.Now()
	defer func() {
		metrics.StoreDurationMs.Observe(float64(time.Since(now).Milliseconds()))
	}()

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
