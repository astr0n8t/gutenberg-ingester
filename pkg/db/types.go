package db

import (
	"sync"

	"github.com/astr0n8t/gutenberg-ingester/pkg/history"
)

type DB struct {
	Version      int             `json:"version"`
	LastFullSync string          `json:"last_full_sync"`
	Download     history.History `json:"download_history"`
	lock         sync.Mutex
}
