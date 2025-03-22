package db

import (
	"github.com/astr0n8t/gutenberg-ingester/internal/history"
)

func NewDB() *DB {
	return &DB{
		Version:  1,
		Download: *history.NewHistory(),
	}
}
