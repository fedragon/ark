package db

import (
	"context"
	"time"
)

type Media struct {
	Hash       []byte     `json:"hash"`
	Path       string     `json:"path"`
	CreatedAt  time.Time  `json:"created_at"`
	ImportedAt *time.Time `json:"imported_at,omitempty"`
	Err        error      `json:"-"`
}

type Repository interface {
	Close() error

	// Store stores a media in the database
	Store(ctx context.Context, media Media) error

	// Get returns a media from the database
	Get(ctx context.Context, hash []byte) (*Media, error)
}
