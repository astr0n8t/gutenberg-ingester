package db

import (
	"sync"
	"time"

	"github.com/astr0n8t/gutenberg-ingester/pkg/history"
)

type DB struct {
	Version         int             `json:"version"`
	LastFullSync    time.Time       `json:"last_full_sync"`
	LastPartialSync time.Time       `json:"last_partial_sync"`
	Download        history.History `json:"download_history"`
	lock            sync.Mutex
}
