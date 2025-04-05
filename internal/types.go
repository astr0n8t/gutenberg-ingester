package internal

import "github.com/astr0n8t/gutenberg-ingester/pkg/db"

type Book interface {
	Id() (int, error)
	Name() (string, error)
	URL() (string, error)
	Language() (string, error)
}

type Runtime struct {
	DB      *db.DB
	Config  ConfigStore
	Catalog string
}
