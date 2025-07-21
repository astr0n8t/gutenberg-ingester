package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/astr0n8t/gutenberg-ingester/pkg/db"
)

func getDB(config ConfigStore, autoSave bool) (*db.DB, error) {
	dbLocation := config.GetString("database_location")

	db, dbOpenErr := db.OpenDBFromFile(dbLocation)
	if dbOpenErr != nil {
		return nil, fmt.Errorf("cannot open DB at location: %v with error: %v", dbLocation, dbOpenErr)
	}

	if autoSave {
		go startSaveDBThread(db, dbLocation)
		log.Printf("db will be saved to %v every 60 seconds", dbLocation)
	}

	return db, nil
}

func startSaveDBThread(db *db.DB, dbSaveFile string) {
	for {
		// save state once a second
		time.Sleep(60 * time.Second)
		err := db.WriteDBToFile(dbSaveFile)
		if err != nil {
			log.Fatalf("issue saving db: %v", err)
		}
	}
}
