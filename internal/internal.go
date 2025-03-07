package internal

import (
	"log"
	"os"
	"os/signal"
	"encoding/xml"
	"io/ioutil"

	"github.com/astr0n8t/gutenberg-ingester/config"
)

// Runs gutenberg-ingester
func Run() {

	// Make sure we can load config
	config := config.Config()
	log.Printf("Loaded config file %v", config.ConfigFileUsed())

	// Insert main app code here
	xmlFile, _ := os.Open("tests/example.xml")

	bytes, _ := ioutil.ReadAll(xmlFile)

	var collection Collection

	xml.Unmarshal(bytes, &collection)

	url, _ := collection.Records[0].URL()
	log.Printf("URL is %v", url)
	url, _ = collection.Records[1].URL()
	log.Printf("URL is %v", url)

	// Don't exit until we receive stop from the OS
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+c to exit")
	<-stop
}
