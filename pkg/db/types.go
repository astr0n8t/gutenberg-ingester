package db

import (
	"github.com/astr0n8t/gutenberg-ingester/pkg/history"
)

type DB struct {
	Version  int             `json:"version"`
	Download history.History `json:"download_history"`
}
