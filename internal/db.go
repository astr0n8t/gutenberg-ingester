package internal

import (
	"fmt"
	"time"
	"log"
	"github.com/astr0n8t/gutenberg-ingester/pkg/db"
)

func getDB(config ConfigStore) (*db.DB, error) {
	dbLocation := config.GetString("database_location")

	db, dbOpenErr := db.OpenDBFromFile(dbLocation)
	if dbOpenErr != nil {
		return nil, fmt.Errorf("cannot open DB at location: %v with error: %v", dbLocation, dbOpenErr)
	}

	go startSaveDBThread(db, dbLocation)

	return db, nil
}

func startSaveDBThread(db *db.DB, dbSaveFile string) {
	for {
		// save state once a second
		time.Sleep(1 * time.Second)
		err := db.WriteDBToFile(dbSaveFile)
		if err != nil {
			log.Fatalf("issue saving db: %v", err)
		}
	}

}
