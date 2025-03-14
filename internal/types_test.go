package internal

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseXML(t *testing.T) {
	xmlFile, err := os.Open("../tests/example.xml")
	if err != nil {
		t.Errorf("Unable to open tests/example.xml %v", err)
	}

	bytes, _ := ioutil.ReadAll(xmlFile)

	var collection Collection

	err = xml.Unmarshal(bytes, &collection)

	if err != nil {
		t.Errorf("Issue unmarshalling XML: %v", err)
	}

	for _, r := range collection.Records {
		id, idErr := r.Id()
		if idErr != nil {
			t.Errorf("Issue unmarshalling XML: %v", idErr)
		}
		fmt.Printf("ID is %v\n", id)
		title, titleErr := r.Title()
		if titleErr != nil {
			t.Errorf("Issue unmarshalling XML: %v", titleErr)
		}
		fmt.Printf("Title is %v\n", title)
		url, urlErr := r.URL()
		if urlErr != nil {
			t.Errorf("Issue unmarshalling XML: %v", urlErr)
		}
		fmt.Printf("URL is %v\n", url)
	}
}

func TestParseRSS(t *testing.T) {
	xmlFile, err := os.Open("../tests/example.rss")
	if err != nil {
		t.Errorf("Unable to open tests/example.rss %v", err)
	}

	bytes, _ := ioutil.ReadAll(xmlFile)

	var rss RSS

	err = xml.Unmarshal(bytes, &rss)

	if err != nil {
		t.Errorf("Issue unmarshalling XML: %v", err)
	}

	for _, i := range rss.Channel.Items {
		id, idErr := i.Id()
		if idErr != nil {
			t.Errorf("Issue unmarshalling XML: %v", idErr)
		}
		fmt.Printf("ID is %v\n", id)
		title, titleErr := i.Title()
		if titleErr != nil {
			t.Errorf("Issue unmarshalling XML: %v", titleErr)
		}
		fmt.Printf("Title is %v\n", title)
		url, urlErr := i.URL()
		if urlErr != nil {
			t.Errorf("Issue unmarshalling XML: %v", urlErr)
		}
		fmt.Printf("URL is %v\n", url)
	}
}
