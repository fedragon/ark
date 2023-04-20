package db

import (
	"context"
	"time"

	"github.com/uptrace/bun"
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
