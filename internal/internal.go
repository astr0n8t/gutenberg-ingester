package internal

import (
	"log"
	"os"
	"os/signal"

	"github.com/astr0n8t/gutenberg-ingester/config"
)

// Runs gutenberg-ingester
func Run() {

	// Make sure we can load config
	config := config.Config()
	log.Printf("Loaded config file %v", config.ConfigFileUsed())

	// Don't exit until we receive stop from the OS
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+c to exit")
	<-stop
}
