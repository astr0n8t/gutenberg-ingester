package rdf

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"testing"
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
}
