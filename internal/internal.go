package internal

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

// Runs gutenberg-ingester daemon
func Run() {
	var config ConfigStore

	// Make sure we can load config
	if getRunMode() == "development" {
		config = DevConfig()
		log.Printf("Running in development mode, will read config from environment")
	} else {
		config = Config()
		log.Printf("Loaded config file %v", config.ConfigFileUsed())
	}

	log.Printf("attempting to open or create db at location: %v", config.GetString("database_location"))
	db, dbErr := getDB(config, true)
	if dbErr != nil {
		log.Fatalf("issue initializing db: %v", dbErr)
	} else {
		log.Printf("successfuly opened db at location: %v", config.GetString("database_location"))
	}

	log.Printf("DB 0 downloaded: %v", db.GetDownloaded(0))
	db.SetDownloaded(0)
	log.Printf("DB 0 downloaded: %v", db.GetDownloaded(0))

	//	catalog, catalogErr := pullCatalog(config)
	//	if catalogErr != nil {
	//		log.Printf("error pulling catalog: %v", catalogErr)
	//	} else {
	//		log.Printf("pulled catalog and saved to: %v", catalog)
	//	}
	catalog := "/tmp/gutenberg-catalog.tar.zip"
	rdfNum, rdfRecordErr := getNumberOfRDFRecords(catalog)
	if rdfRecordErr != nil {
		log.Printf("issue reading rdf number from catalog: %v", rdfRecordErr)
	} else {
		log.Printf("Number of RDF records is %v", rdfNum)
	}

	rdfRecord, rdfErr := getRDFByID(1, catalog)
	if rdfErr != nil {
		log.Printf("issue reading rdf file from catalog: %v", rdfErr)
	}

	rdfName, rdfNameError := rdfRecord.Name()
	if rdfNameError != nil {
		log.Printf("issue reading rdf name from catalog: %v", rdfNameError)
	}
	log.Printf("Title of RDF record tile: %v", rdfName)

	// Don't exit until we receive stop from the OS
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+c to exit")
	<-stop
}

func getRunMode() string {
	runMode := "production"
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("cannot determine run status: %v", err)
	}

	// Check if we're running in a dev build
	dir := filepath.Dir(ex)
	if strings.Contains(dir, "go-build") {
		runMode = "development"
	}

	return runMode
}
