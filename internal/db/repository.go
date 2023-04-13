package db

import "time"

type Media struct {
	Hash    []byte
	Path    string
	ModTime time.Time
	Err     error `json:"-"`
}

type Repository interface {
	// Store stores a media in the repository
	Store(media Media) error

	// List returns a channel of all media in the repository
	List() <-chan Media
}
