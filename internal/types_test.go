package internal

import (
	"testing"
	"encoding/xml"
	"io/ioutil"
	"os"
)

func TestParseXML(t *testing.T) {
	xmlFile, err := os.Open("tests/example.xml")
	if err != nil {
		t.Errorf("Unable to open tests/example.xml")
	}

	bytes, _ := ioutil.ReadAll(xmlFile)

	var collection Collection

	err := xml.Unmarshal(bytes, &collection)

	if err != nil {
		t.Errorf("Issue unmarshalling XML: ", err)
	}
}

