package db

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/astr0n8t/gutenberg-ingester/pkg/history"
)

func NewDB() *DB {
	return &DB{
		Version:         1,
		LastFullSync:    time.Unix(0, 0),
		LastPartialSync: time.Unix(0, 0),
		Download:        *history.NewHistory(),
	}
}

func OpenDBFromFile(filename string) (*DB, error) {
	var db DB
	db.lock.Lock()

	fileInfo, fileInfoErr := os.Stat(filename)
	if os.IsNotExist(fileInfoErr) {
		db.lock.Unlock()
		file, fileCreateErr := os.Create(filename)
		if fileCreateErr != nil {
			return nil, fmt.Errorf("failed to create download database file: %v with error: %v", filename, fileCreateErr)
		}
		// If we can create the file, hopefully we can read the file
		fileInfo, fileInfoErr = os.Stat(filename)
		if fileInfoErr != nil {
			return nil, fmt.Errorf("failed to read created download database file: %v with error: %v", filename, fileInfoErr)
		}
		file.Close()

		db = *NewDB()
		return &db, db.WriteDBToFile(filename)
	}

	defer db.lock.Unlock()

	file, fileOpenErr := os.Open(filename)
	if fileOpenErr != nil {
		return nil, fmt.Errorf("failed to open download database file: %v with error: %v", filename, fileOpenErr)
	}

	jsonData := make([]byte, fileInfo.Size())
	_, fileReadErr := file.Read(jsonData)
	if fileReadErr != nil {
		return nil, fmt.Errorf("failed to read download database file: %v with error: %v", filename, fileReadErr)
	}

	jsonErr := json.Unmarshal(jsonData, &db)
	if jsonErr != nil {
		return nil, fmt.Errorf("failed to unmarshal download database file: %v with error: %v", filename, fileReadErr)
	}

	if db.Version != 1 {
		return nil, fmt.Errorf("unknown database file version: %v when opening file: %v", db.Version, filename)
	}

	return &db, nil
}

// Write DB to file, releasing lock when done
func (d *DB) WriteDBToFile(filename string) error {
	return d.writeDBToFile(filename, true)
}

// Write DB to file, and don't release the lock when done
func (d *DB) WriteDBToFileAndLock(filename string) error {
	return d.writeDBToFile(filename, false)
}

func (d *DB) writeDBToFile(filename string, releaseLock bool) error {
	d.lock.Lock()
	if releaseLock {
		defer d.lock.Unlock()
	}

	fileInfo, fileInfoErr := os.Stat(filename)
	if os.IsNotExist(fileInfoErr) {
		file, fileCreateErr := os.Create(filename)
		if fileCreateErr != nil {
			return fmt.Errorf("failed to create download database file: %v with error: %v", filename, fileCreateErr)
		}
		// If we can create the file, hopefully we can read the file
		fileInfo, fileInfoErr = os.Stat(filename)
		if fileInfoErr != nil {
			return fmt.Errorf("failed to read created download database file: %v with error: %v", filename, fileInfoErr)
		}
		file.Close()
	}

	file, fileErr := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileInfo.Mode())
	if fileErr != nil {
		return fmt.Errorf("failed to open download database file: %v with error: %v", filename, fileErr)
	}

	jsonData, jsonErr := json.MarshalIndent(d, "", " ")
	if jsonErr != nil {
		return fmt.Errorf("unable to marshal DB to json %v", jsonErr)
	}

	_, writeErr := file.Write(jsonData)
	if writeErr != nil {
		return fmt.Errorf("unable to write DB to download database file: %v with error: %v", filename, writeErr)
	}

	file.Close()

	return nil
}

func (d *DB) SetDownloaded(id int) {
	d.lock.Lock()
	d.Download.SetHistory(id)
	d.lock.Unlock()
}

func (d *DB) UnsetDownloaded(id int) {
	d.lock.Lock()
	d.Download.UnsetHistory(id)
	d.lock.Unlock()
}

func (d *DB) GetDownloaded(id int) bool {
	return d.Download.GetHistory(id)
}

func (d *DB) GetLastFullSync() time.Time {
	return d.LastFullSync
}

func (d *DB) SetLastFullSync(date time.Time) {
	d.lock.Lock()
	d.LastFullSync = date
	d.lock.Unlock()
}

func (d *DB) GetLastPartialSync() time.Time {
	return d.LastPartialSync
}

func (d *DB) SetLastPartialSync(date time.Time) {
	d.lock.Lock()
	d.LastPartialSync = date
	d.lock.Unlock()
}
