package rss

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseRSS(t *testing.T) {
	xmlFile, err := os.Open("../../tests/example.rss")
	if err != nil {
		t.Errorf("Unable to open tests/example.rss %v", err)
	}

	bytes, _ := ioutil.ReadAll(xmlFile)

	var rss RSS

	err = xml.Unmarshal(bytes, &rss)

	if err != nil {
		t.Errorf("Issue unmarshling XML: %v", err)
	}

	for _, i := range rss.Channel.Items {
		id, idErr := i.Id()
		if idErr != nil {
			t.Errorf("Issue unmarshling XML: %v", idErr)
		}
		fmt.Printf("ID is %v\n", id)
		title, titleErr := i.Title()
		if titleErr != nil {
			t.Errorf("Issue unmarshling XML: %v", titleErr)
		}
		fmt.Printf("Title is %v\n", title)
		url, urlErr := i.URL()
		if urlErr != nil {
			t.Errorf("Issue unmarshling XML: %v", urlErr)
		}
		fmt.Printf("URL is %v\n", url)
	}
}
