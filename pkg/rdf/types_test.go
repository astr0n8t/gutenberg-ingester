package rdf

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"testing"
	"fmt"
)

func TestParseXML(t *testing.T) {
	xmlFile, err := os.Open("../../tests/example.rdf")
	if err != nil {
		t.Errorf("Unable to open tests/example.rdf %v", err)
	}

	bytes, _ := ioutil.ReadAll(xmlFile)

	var rdf RDF

	err = xml.Unmarshal(bytes, &rdf)

	if err != nil {
		t.Errorf("Issue unmarshling XML: %v", err)
	}

	id, idErr := rdf.Id()
	if idErr != nil {
		t.Errorf("Issue unmarshling XML: %v", idErr)
	}
	fmt.Printf("ID is %v\n", id)
	title, titleErr := rdf.Title()
	if titleErr != nil {
		t.Errorf("Issue unmarshling XML: %v", titleErr)
	}
	fmt.Printf("Title is %v\n", title)
	url, urlErr := rdf.URL()
	if urlErr != nil {
		t.Errorf("Issue unmarshling XML: %v", urlErr)
	}
	fmt.Printf("URL is %v\n", url)
}
