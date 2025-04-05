package internal

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

// Runs gutenberg-ingester daemon
func Run() {
	var runtime Runtime

	// Make sure we can load config
	if getRunMode() == "development" {
		runtime.Config = DevConfig()
		log.Printf("Running in development mode, will read config from environment")
	} else {
		runtime.Config = Config()
		log.Printf("Loaded config file %v", runtime.Config.ConfigFileUsed())
	}

	log.Printf("attempting to open or create db at location: %v", runtime.Config.GetString("database_location"))
	var dbErr error
	runtime.DB, dbErr = getDB(runtime.Config, true)
	if dbErr != nil {
		log.Fatalf("issue initializing db: %v", dbErr)
	} else {
		log.Printf("successfuly opened db at location: %v", runtime.Config.GetString("database_location"))
	}

	// Start the main sync thread
	log.Printf("Starting sync scheduler")
	go runtime.startSyncSchedule()

	// Don't exit until we receive stop from the OS
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Printf("Writing DB before exiting\n")

	// Lock the DB exclusively and hold the lock
	err := runtime.DB.WriteDBToFileAndLock(runtime.Config.GetString("database_location"))
	if err != nil {
		log.Fatalf("issue saving db: %v", err)
	}

	log.Printf("DB written, exiting now.\n")
	// exit
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
