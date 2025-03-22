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


	rss, err := pullRSS()
	if err != nil {
		log.Printf("issue pulling rss feed: %v", err)
	}

	for _, i := range rss.Channel.Items {
		id, idErr := i.Id()
		if idErr != nil {
			log.Printf("Issue unmarshling XML: %v", idErr)
		}
		log.Printf("ID is %v\n", id)
		title, titleErr := i.Title()
		if titleErr != nil {
			log.Printf("Issue unmarshling XML: %v", titleErr)
		}
		log.Printf("Title is %v\n", title)
	}

	// Don't exit until we receive stop from the OS
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+c to exit")
	<-stop
}
